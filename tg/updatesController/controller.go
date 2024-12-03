package updatesController

import (
	"errors"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/jmoiron/sqlx"

	"tgBotCompetition/l10n"
	"tgBotCompetition/usecase/chats/take"
	compCreate "tgBotCompetition/usecase/competitions/create"
)

type Controller struct {
	DB *sqlx.DB
}

func (c *Controller) AddHandlers(dispatcher *ext.Dispatcher) error {
	dispatcher.AddHandler(handlers.NewChatMember(nil, c.newChatMember))
	dispatcher.AddHandler(handlers.NewCommand("competitionConfigRun", c.competitionConfigRun))
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
	var params compCreate.Params

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

	// Заполнить остальные параметры
	params.Keyword = kv[l10n.CfgKeyword]

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
