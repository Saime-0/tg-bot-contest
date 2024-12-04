package updatesController

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/jmoiron/sqlx"

	"tgBotCompetition/l10n"
	tgModel "tgBotCompetition/tg/model"
	"tgBotCompetition/usecase/chat/take"
	compCreate "tgBotCompetition/usecase/competitions/create"
	competitionStop "tgBotCompetition/usecase/competitions/stop"
	messageCreated "tgBotCompetition/usecase/message/created"
)

type Controller struct {
	DB *sqlx.DB
}

func (c *Controller) AddHandlers(dispatcher *ext.Dispatcher) error {
	dispatcher.AddHandler(handlers.NewMessage(func(msg *gotgbot.Message) bool {
		return msg.Chat.Type == gotgbot.ChatTypePrivate && strings.HasPrefix(msg.GetText(), "/competitionConfigRun")
	}, c.competitionConfigRun))
	dispatcher.AddHandler(handlers.NewMessage(func(msg *gotgbot.Message) bool {
		return msg.Chat.Type == gotgbot.ChatTypePrivate && strings.HasPrefix(msg.GetText(), "/competitionStop")
	}, c.competitionStop))
	dispatcher.AddHandler(handlers.NewMessage(nil, c.newMessage))
	dispatcher.AddHandler(handlers.NewChatMember(nil, c.newChatMember))
	//dispatcher.AddHandler(handlers.NewCommand("compr", c.competitionConfigRun))
	//dispatcher.AddHandler(handlers.NewCommand("competitionStop", c.competitionStop))

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

func (c *Controller) competitionConfigRun(b *gotgbot.Bot, ctx *ext.Context) error {
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
	if chat, err := (&take.Params{
		DB:       c.DB,
		Username: strings.TrimPrefix(kv[l10n.CfgChatUsername], "@"),
	}).Run(); err != nil {
		return err
	} else {
		params.ChatID = chat.ID
	}

	// Проверить наличие прав админа в чате
	if admins, err := b.GetChatAdministrators(int64(params.ChatID), nil); err != nil {
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
			return errors.New("user is not allowed to create competition")
		}
	}

	if multiplicity, err := strconv.ParseInt(kv[l10n.CfgMultiplicity], 10, 64); err != nil {
		return err
	} else {
		params.Multiplicity = int(multiplicity)
	}

	if topicID, err := strconv.ParseInt(kv[l10n.CfgTopic], 10, 64); err != nil {
		return err
	} else {
		params.TopicID = int(topicID)
	}

	if err := params.Run(); err != nil {
		return err
	}

	return nil
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
		return err
	}

	if len(messageCreatedOut.CreatedTickets) > 0 {
		if err = c.sendMessageAboutCreatedTickets(b, ctx, messageCreatedOut); err != nil {
			slog.Error("sendMessageAboutCreatedTickets: " + err.Error())
		}
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

func (c *Controller) competitionStop(b *gotgbot.Bot, ctx *ext.Context) error {
	words := strings.Fields(ctx.Message.GetText())
	if len(words) < 2 {
		return nil // todo: send msg
	}

	if err := (&competitionStop.Params{
		DB:           c.DB,
		ChatUsername: strings.TrimPrefix(words[1], "@"),
	}).Run(); err != nil {
		return err
	}

	return nil
}
