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

func (c *Controller) AddHandlers(dispatcher *ext.Dispatcher) error {
	dispatcher.AddHandler(handlers.NewMessage(func(msg *gotgbot.Message) bool {
		return msg.Chat.Type == gotgbot.ChatTypePrivate && strings.HasPrefix(msg.GetText(), "/contestConfigRun")
	}, c.contestConfigRun))
	dispatcher.AddHandler(handlers.NewMessage(func(msg *gotgbot.Message) bool {
		return msg.Chat.Type == gotgbot.ChatTypePrivate && strings.HasPrefix(msg.GetText(), "/contestStop")
	}, c.contestStop))
	dispatcher.AddHandler(handlers.NewMessage(nil, c.newMessage))
	dispatcher.AddHandler(handlers.NewChatMember(nil, c.newChatMember))
	//dispatcher.AddHandler(handlers.NewCommand("compr", c.contestConfigRun))
	//dispatcher.AddHandler(handlers.NewCommand("contestStop", c.contestStop))

	return nil
}

func (c *Controller) SetMyCommands(b *gotgbot.Bot) error {
	//_, err := b.SetMyCommands([]gotgbot.BotCommand{
	//	{
	//		Command:     "/start",
	//		Description: l10n.CmdDescStart,
	//	},
	//}, &gotgbot.SetMyCommandsOpts{
	//	Scope: gotgbot.BotCommandScopeAllPrivateChats{},
	//})

	return nil
}

func (c *Controller) contestConfigRun(b *gotgbot.Bot, ctx *ext.Context) error {
	// Разобрать сообщение конфига на параметры
	lines := strings.Split(ctx.Message.GetText(), "\n")
	kv := make(map[string]string, len(lines))
	for _, line := range lines {
		parts := strings.Split(strings.TrimSpace(line), l10n.CfgDelimiter)
		if len(parts) == 2 {
			kv[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	// Параметры для создания конкурса
	params := compCreate.Params{
		DB:        c.DB,
		Keyword:   kv[l10n.CfgKeyword],
		CreatorID: int(ctx.EffectiveSender.Id()),
	}

	// Найти чат
	chatUsername := strings.TrimPrefix(kv[l10n.CfgChatUsername], "@")
	if chat, err := chatTake.Run(c.DB, chatUsername); err != nil {
		return c.reactError(ue.Sql(err), b, ctx)
	} else {
		params.ChatID = chat.ID
	}

	// Проверить наличие прав админа в чате
	if err := c.checkAdminRights(b, ctx, int64(params.ChatID)); err != nil {
		return c.reactError(err, b, ctx)
	}

	if multiplicity, err := strconv.ParseInt(kv[l10n.CfgMultiplicity], 10, 64); err != nil {
		_, err = ctx.Message.Reply(b, l10n.ContestConfigRunUsage, nil)
		return err
	} else {
		params.Multiplicity = int(multiplicity)
	}

	if topicID, err := strconv.ParseInt(kv[l10n.CfgTopic], 10, 64); err != nil {
		_, err = ctx.Message.Reply(b, l10n.ContestConfigRunUsage, nil)
		return err
	} else {
		params.TopicID = int(topicID)
	}

	if err := params.Run(); err != nil {
		return c.reactError(err, b, ctx)
	}

	_, err := ctx.Message.Reply(b, l10n.ContestConfigRunSuccess, nil)
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

func (c *Controller) newMessage(b *gotgbot.Bot, ctx *ext.Context) (err error) {
	msg := ctx.Message
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
		DB:      c.DB,
		Chat:    tgModel.ChatDomain(msg.Chat),
		User:    tgModel.UserDomain(*msg.From),
		Text:    msg.GetText(),
		TopicID: topicID,
	}).Run(); err != nil {
		return c.reactError(err, b, ctx)
	}

	if len(messageCreatedOut.CreatedTickets) > 0 {
		err = c.sendMessageAboutCreatedTickets(b, ctx, messageCreatedOut)
		return err
	} else if messageCreatedOut.CalculationWasStarted {
		_, err = ctx.Message.Reply(b, l10n.DintGetRightNumberOfInvitations, nil)
		return err
	}

	return nil
}

func (c *Controller) sendMessageAboutCreatedTickets(b *gotgbot.Bot, ctx *ext.Context, o messageCreated.Out) error {
	numbers := make([]string, len(o.CreatedTickets))
	for i := range o.CreatedTickets {
		numbers[i] = strconv.Itoa(o.CreatedTickets[i].Number)
	}

	text := l10n.YourTicketNumbers + strings.Join(numbers, l10n.YourTicketNumbersDelimiter)
	if _, err := ctx.Message.Reply(b, text, nil); err != nil {
		return err
	}

	return nil
}

func (c *Controller) contestStop(b *gotgbot.Bot, ctx *ext.Context) (err error) {
	words := strings.Fields(ctx.Message.GetText())
	if len(words) < 2 {
		_, err = ctx.Message.Reply(b, l10n.ContestStopUsage, nil)
		return err
	}

	// Найти чат
	var chat model.Chat
	if chat, err = chatTake.Run(c.DB, strings.TrimPrefix(words[1], "@")); err != nil {
		return c.reactError(ue.Sql(err), b, ctx)
	}

	// Проверить наличие прав админа в чате
	if err = c.checkAdminRights(b, ctx, int64(chat.ID)); err != nil {
		return c.reactError(err, b, ctx)
	}

	if err = (&contestStop.Params{
		DB:     c.DB,
		ChatID: chat.ID,
	}).Run(); err != nil {
		return c.reactError(err, b, ctx)
	}

	_, err = ctx.Message.Reply(b, l10n.ContestStopSuccess, nil)
	return err
}

func (c *Controller) checkAdminRights(b *gotgbot.Bot, ctx *ext.Context, chatID int64) error {
	if admins, err := b.GetChatAdministrators(chatID, nil); err != nil {
		return err
	} else {
		var allowed bool
		for _, admin := range admins {
			if admin.GetUser().Id == ctx.EffectiveSender.Id() {
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
