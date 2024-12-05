package updatesController

import (
	"slices"
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
	memberStatusUpdate "tgBotContest/usecase/member/statusUpdate"
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
		TX:        nil,
		Keyword:   kv[l10n.CfgKeyword],
		CreatorID: int(r.ctx.EffectiveSender.Id()),
	}

	// Найти чат
	if chat, err := chatTake.Run(r.DB, clearAt(kv[l10n.CfgChatUsername])); err != nil {
		return r.reactError(err)
	} else {
		params.ChatID = chat.ID
	}

	// Проверить наличие прав админа в чате
	if err := r.checkAdminRights(int64(params.ChatID)); err != nil {
		return r.reactError(err)
	}

	if multiplicity, err := strconv.ParseInt(kv[l10n.CfgMultiplicity], 10, 64); err != nil {
		return right(r.ctx.Message.Reply(r.Bot, l10n.ContestConfigRunUsage, nil))
	} else {
		params.Multiplicity = int(multiplicity)
	}

	if topicID, err := strconv.ParseInt(kv[l10n.CfgTopic], 10, 64); err != nil {
		return right(r.ctx.Message.Reply(r.Bot, l10n.ContestConfigRunUsage, nil))
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
