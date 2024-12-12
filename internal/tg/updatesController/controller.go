package updatesController

import (
	"log/slog"
	"slices"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/jmoiron/sqlx"

	"github.com/Saime-0/tg-bot-contest/internal/l10n"
	"github.com/Saime-0/tg-bot-contest/internal/model"
	tgModel "github.com/Saime-0/tg-bot-contest/internal/tg/model"
	"github.com/Saime-0/tg-bot-contest/internal/ue"
	chatTake "github.com/Saime-0/tg-bot-contest/internal/usecase/chat/take"
	chatUpdate "github.com/Saime-0/tg-bot-contest/internal/usecase/chat/update"
	contestCreate "github.com/Saime-0/tg-bot-contest/internal/usecase/contests/create"
	contestStop "github.com/Saime-0/tg-bot-contest/internal/usecase/contests/stop"
	memberStatusUpdate "github.com/Saime-0/tg-bot-contest/internal/usecase/member/statusUpdate"
	messageCreated "github.com/Saime-0/tg-bot-contest/internal/usecase/message/created"
)

type Controller struct {
	DB  *sqlx.DB
	Bot *gotgbot.Bot
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
		handlers.NewMessage(c.inviteBotFilter, c.modulation(inviteBotMessage)),
		handlers.NewMessage(c.leftBotFilter, c.modulation(leftBotMessage)),
		handlers.NewMessage(nil, c.modulation(newMessage)),
		handlers.NewChatMember(nil, c.modulation(newChatMember)),
	}

	for _, h := range handlerGroup {
		dispatcher.AddHandler(h)
	}

	return nil
}

func leftBotMessage(r Request) error {
	chat := tgModel.ChatDomain(r.ctx.Message.Chat)
	if err := chatUpdate.Run(r.DB, chat); err != nil {
		slog.Warn("leftBotMessage: chatUpdate.Run: " + err.Error())
	}

	if err := InTransaction(r.DB, func(tx *sqlx.Tx) error {
		return (&contestStop.Params{TX: tx, ChatID: chat.ID}).Run()
	}); err != nil {
		slog.Debug("leftBotMessage: " + err.Error())
	}

	return nil
}

func inviteBotMessage(r Request) error {
	chat := tgModel.ChatDomain(r.ctx.Message.Chat)
	return chatUpdate.Run(r.DB, chat)
}

func (c *Controller) inviteBotFilter(msg *gotgbot.Message) bool {
	for _, member := range msg.NewChatMembers {
		if member.Id == c.Bot.Id {
			return true
		}
	}

	return false
}

type Request struct {
	DB *sqlx.DB
	*gotgbot.Bot
	ctx *ext.Context
}

func (c *Controller) modulation(fn func(request Request) error) func(b *gotgbot.Bot, ctx *ext.Context) error {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		return fn(Request{
			DB:  c.DB,
			Bot: b,
			ctx: ctx,
		})
	}
}

func (c *Controller) leftBotFilter(msg *gotgbot.Message) bool {
	return msg.LeftChatMember.Id == c.Bot.Id
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

	chatUsername := clearAt(kv[l10n.CfgChatUsername])

	// Параметры для создания конкурса
	params := contestCreate.Params{
		TX:        nil,
		Keyword:   kv[l10n.CfgKeyword],
		CreatorID: int(r.ctx.EffectiveSender.Id()),
	}

	if chatID, err := strconv.ParseInt(kv[l10n.CfgChatID], 10, 64); err != nil && kv[l10n.CfgChatID] != "" {
		return right(r.ctx.Message.Reply(r.Bot, l10n.ContestConfigRunUsage, &gotgbot.SendMessageOpts{
			ParseMode: gotgbot.ParseModeMarkdownV2,
		}))
	} else if chatID != 0 {
		params.ChatID = int(chatID)
	}
	if params.ChatID == 0 {
		// Найти чат
		if chat, err := chatTake.Run(r.DB, chatUsername); err != nil {
			return r.reactError(err)
		} else {
			params.ChatID = chat.ID
		}
	}

	// Проверить наличие прав админа в чате
	if err := r.checkAdminRights(int64(params.ChatID)); err != nil {
		return r.reactError(err)
	}

	if multiplicity, err := strconv.ParseInt(kv[l10n.CfgMultiplicity], 10, 64); err != nil {
		return right(r.ctx.Message.Reply(r.Bot, l10n.ContestConfigRunUsage, &gotgbot.SendMessageOpts{
			ParseMode: gotgbot.ParseModeMarkdownV2,
		}))
	} else {
		params.Multiplicity = int(multiplicity)
	}

	if topicID, err := strconv.ParseInt(kv[l10n.CfgTopic], 10, 64); err != nil {
		return right(r.ctx.Message.Reply(r.Bot, l10n.ContestConfigRunUsage, &gotgbot.SendMessageOpts{
			ParseMode: gotgbot.ParseModeMarkdownV2,
		}))
	} else {
		params.TopicID = int(topicID)
	}

	if err := InTransaction(r.DB, func(tx *sqlx.Tx) error {
		params.TX = tx
		return params.Run()
	}); err != nil {
		return r.reactError(err)
	}

	return right(r.ctx.Message.Reply(r.Bot, l10n.ContestConfigRunSuccess, nil))
}

func newMessage(r Request) (err error) {
	chat := tgModel.ChatDomain(r.ctx.Message.Chat)
	if err = chatUpdate.Run(r.DB, chat); err != nil {
		slog.Warn("newMessage: chatUpdate.Run: " + err.Error())
	}

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
		DB:      r.DB,
		Chat:    tgModel.ChatDomain(msg.Chat),
		User:    tgModel.UserDomain(*msg.From),
		Text:    msg.GetText(),
		TopicID: topicID,
	}).Run(); err != nil {
		return r.reactError(err)
	}

	if len(messageCreatedOut.CreatedTickets) > 0 {
		return r.sendMessageAboutCreatedTickets(messageCreatedOut)
	} else if messageCreatedOut.CalculationWasStarted {
		return right(r.ctx.Message.Reply(r.Bot, l10n.DintGetRightNumberOfInvitations, nil))
	}

	return nil
}

func (r Request) sendMessageAboutCreatedTickets(o messageCreated.Out) error {
	numbers := make([]string, len(o.CreatedTickets))
	for i := range o.CreatedTickets {
		numbers[i] = strconv.Itoa(o.CreatedTickets[i].Number)
	}

	text := l10n.YourTicketNumbers + strings.Join(numbers, l10n.YourTicketNumbersDelimiter)
	return right(r.ctx.Message.Reply(r.Bot, text, nil))
}

func contestStopHandler(r Request) (err error) {
	words := strings.Fields(r.ctx.Message.GetText())
	if len(words) < 2 {
		return right(r.ctx.Message.Reply(r.Bot, l10n.ContestStopUsage, nil))
	}

	// Найти чат
	var chat model.Chat
	if chat, err = chatTake.Run(r.DB, clearAt(words[1])); err != nil {
		return r.reactError(err)
	}

	// Проверить наличие прав админа в чате
	if err = r.checkAdminRights(int64(chat.ID)); err != nil {
		return r.reactError(err)
	}

	if err = InTransaction(r.DB, func(tx *sqlx.Tx) error {
		return (&contestStop.Params{TX: tx, ChatID: chat.ID}).Run()
	}); err != nil {
		return r.reactError(err)
	}

	return right(r.ctx.Message.Reply(r.Bot, l10n.ContestStopSuccess, nil))
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

func defineMemberStatus(old, new string) uint {
	oldStatus := tgModel.MemberStatusID[old]
	newStatus := tgModel.MemberStatusID[new]
	switch {
	case slices.Contains(tgModel.AlienStatus, oldStatus) && slices.Contains(tgModel.ParticipantStatus, newStatus):
		return model.MemberStatusJoin
	case slices.Contains(tgModel.ParticipantStatus, oldStatus) && slices.Contains(tgModel.AlienStatus, newStatus):
		return model.MemberStatusLeave
	default:
		return 0
	}
}

func newChatMember(r Request) error {
	chat := tgModel.ChatDomain(r.ctx.Message.Chat)
	if err := chatUpdate.Run(r.DB, chat); err != nil {
		slog.Warn("newChatMember: chatUpdate.Run: " + err.Error())
	}

	oldStatus := r.ctx.ChatMember.OldChatMember.GetStatus()
	newStatus := r.ctx.ChatMember.NewChatMember.GetStatus()

	memberStatus := defineMemberStatus(oldStatus, newStatus)
	if memberStatus == 0 {
		return nil
	}
	initiator := r.ctx.ChatMember.From
	participant := r.ctx.ChatMember.NewChatMember.GetUser()
	viaLink := r.ctx.ChatMember.InviteLink != nil ||
		r.ctx.ChatMember.IsJoinRequest() ||
		r.ctx.ChatMember.ViaChatFolderInviteLink

	memberStatusUpdateParams := &memberStatusUpdate.Params{
		TX:           nil,
		Chat:         tgModel.ChatDomain(r.ctx.ChatMember.Chat),
		MemberStatus: memberStatus,
		Participant:  tgModel.UserDomain(participant),
		Initiator:    tgModel.UserDomain(initiator),
		ViaLink:      viaLink,
	}

	if err := InTransaction(r.DB, func(tx *sqlx.Tx) error {
		memberStatusUpdateParams.TX = tx
		return memberStatusUpdateParams.Run()
	}); err != nil {
		return err
	}

	return nil
}
