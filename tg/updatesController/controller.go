package updatesController

import (
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/jmoiron/sqlx"

	"tgBotContest/l10n"
	"tgBotContest/model"
	tgModel "tgBotContest/tg/model"
	"tgBotContest/ue"
	chatTake "tgBotContest/usecase/chat/take"
	compCreate "tgBotContest/usecase/contests/create"
	contestStop "tgBotContest/usecase/contests/stop"
	messageCreated "tgBotContest/usecase/message/created"
)

type Controller struct {
	DB *sqlx.DB
}

func onlyInPrivateChat(fn func(b *gotgbot.Bot, ctx *ext.Context) error) func(b *gotgbot.Bot, ctx *ext.Context) error {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		if ctx.EffectiveChat != nil &&
			ctx.EffectiveChat.Type == gotgbot.ChatTypePrivate {
			return fn(b, ctx)
		}

		return nil
	}
}

func (c *Controller) AddHandlers(dispatcher *ext.Dispatcher) error {
	handlerGroup := []ext.Handler{
		handlers.NewCommand("contestConfigRun", onlyInPrivateChat(c.modulation(contestConfigRun))),
		handlers.NewCommand("contestStop", onlyInPrivateChat(c.modulation(contestStopHandler))),
		handlers.NewMessage(nil, c.modulation(newMessage)),
		handlers.NewChatMember(nil, c.modulation(newChatMember)),
	}

	for _, h := range handlerGroup {
		dispatcher.AddHandler(h)
	}

	return nil
}

type Request struct {
	TX func() *sqlx.Tx
	*gotgbot.Bot
	ctx *ext.Context
}

func (c *Controller) modulation(fn func(request Request) error) func(b *gotgbot.Bot, ctx *ext.Context) error {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		return InLazyTransaction(c.DB, func(tx func() *sqlx.Tx) error {
			r := Request{
				TX:  tx,
				Bot: b,
				ctx: ctx,
			}
			return fn(r)
		})
	}
}

func contestConfigRun(r Request) error {
	// Разобрать сообщение конфига на параметры
	lines := strings.Split(r.ctx.Message.GetText(), "\n")
	kv := make(map[string]string, len(lines))
	for _, line := range lines {
		parts := strings.Split(strings.TrimSpace(line), l10n.CfgDelimiter)
		if len(parts) == 2 {
			kv[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	// Параметры для создания конкурса
	params := compCreate.Params{
		TX:        r.TX(),
		Keyword:   kv[l10n.CfgKeyword],
		CreatorID: int(r.ctx.EffectiveSender.Id()),
	}

	// Найти чат
	chatUsername := strings.TrimPrefix(kv[l10n.CfgChatUsername], "@")
	if chat, err := chatTake.Run(r.TX(), chatUsername); err != nil {
		return r.reactError(ue.Sql(err))
	} else {
		params.ChatID = chat.ID
	}

	// Проверить наличие прав админа в чате
	if err := r.checkAdminRights(int64(params.ChatID)); err != nil {
		return r.reactError(err)
	}

	if multiplicity, err := strconv.ParseInt(kv[l10n.CfgMultiplicity], 10, 64); err != nil {
		_, err = r.ctx.Message.Reply(r.Bot, l10n.ContestConfigRunUsage, nil)
		return err
	} else {
		params.Multiplicity = int(multiplicity)
	}

	if topicID, err := strconv.ParseInt(kv[l10n.CfgTopic], 10, 64); err != nil {
		_, err = r.ctx.Message.Reply(r.Bot, l10n.ContestConfigRunUsage, nil)
		return err
	} else {
		params.TopicID = int(topicID)
	}

	if err := params.Run(); err != nil {
		return r.reactError(err)
	}

	_, err := r.ctx.Message.Reply(r.Bot, l10n.ContestConfigRunSuccess, nil)
	return err
}

func (c *Controller) cmdStart(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := b.SendMessage(
		ctx.EffectiveSender.Id(),
		"",
		&gotgbot.SendMessageOpts{
			ReplyMarkup: &gotgbot.ReplyKeyboardMarkup{
				Keyboard: [][]gotgbot.KeyboardButton{
					{
						gotgbot.KeyboardButton{
							Text: l10n.KBBtnRunConfig,
						},
					},
				},
				IsPersistent:    false,
				ResizeKeyboard:  true,
				OneTimeKeyboard: true,
			},
		},
	)

	return err
}

func newMessage(r Request) (err error) {
	msg := r.ctx.Message
	if (msg.Chat.Type != gotgbot.ChatTypeSupergroup && msg.Chat.Type != gotgbot.ChatTypeGroup) ||
		!msg.GetSender().IsUser() ||
		msg.From == nil ||
		msg.GetText() == "" {
		return nil
	}
	var topicID int
	if msg.IsTopicMessage {
		topicID = int(msg.MessageThreadId)
	}

	var messageCreatedOut messageCreated.Out
	if messageCreatedOut, err = (&messageCreated.Params{
		TX:      r.TX(),
		Chat:    tgModel.ChatDomain(msg.Chat),
		User:    tgModel.UserDomain(*msg.From),
		Text:    msg.GetText(),
		TopicID: topicID,
	}).Run(); err != nil {
		return r.reactError(err)
	}

	if len(messageCreatedOut.CreatedTickets) > 0 {
		err = r.sendMessageAboutCreatedTickets(messageCreatedOut)
		return err
	} else if messageCreatedOut.CalculationWasStarted {
		_, err = r.ctx.Message.Reply(r.Bot, l10n.DintGetRightNumberOfInvitations, nil)
		return err
	}

	return nil
}

func (r Request) sendMessageAboutCreatedTickets(o messageCreated.Out) error {
	numbers := make([]string, len(o.CreatedTickets))
	for i := range o.CreatedTickets {
		numbers[i] = strconv.Itoa(o.CreatedTickets[i].Number)
	}

	text := l10n.YourTicketNumbers + strings.Join(numbers, l10n.YourTicketNumbersDelimiter)
	if _, err := r.ctx.Message.Reply(r.Bot, text, nil); err != nil {
		return err
	}

	return nil
}

func contestStopHandler(r Request) (err error) {
	words := strings.Fields(r.ctx.Message.GetText())
	if len(words) < 2 {
		_, err = r.ctx.Message.Reply(r.Bot, l10n.ContestStopUsage, nil)
		return err
	}

	// Найти чат
	var chat model.Chat
	if chat, err = chatTake.Run(r.TX(), strings.TrimPrefix(words[1], "@")); err != nil {
		return r.reactError(ue.Sql(err))
	}

	// Проверить наличие прав админа в чате
	if err = r.checkAdminRights(int64(chat.ID)); err != nil {
		return r.reactError(err)
	}

	if err = (&contestStop.Params{
		TX:     r.TX(),
		ChatID: chat.ID,
	}).Run(); err != nil {
		return r.reactError(err)
	}

	_, err = r.ctx.Message.Reply(r.Bot, l10n.ContestStopSuccess, nil)
	return err
}

func (r Request) checkAdminRights(chatID int64) error {
	if admins, err := r.GetChatAdministrators(chatID, nil); err != nil {
		return err
	} else {
		var allowed bool
		for _, admin := range admins {
			if admin.GetUser().Id == r.ctx.EffectiveSender.Id() {
				allowed = true
				break
			}
		}
		if !allowed {
			return ue.New(l10n.CreateContestNoAdminRights)
		}
	}

	return nil
}
