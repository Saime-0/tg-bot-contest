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
	usageErrPkg "github.com/Saime-0/tg-bot-contest/internal/tg/usageErr"
	"github.com/Saime-0/tg-bot-contest/internal/ue"
	chatTake "github.com/Saime-0/tg-bot-contest/internal/usecase/chat/take"
	contestCreate "github.com/Saime-0/tg-bot-contest/internal/usecase/contests/create"
	contestStop "github.com/Saime-0/tg-bot-contest/internal/usecase/contests/stop"
	"github.com/Saime-0/tg-bot-contest/internal/usecase/member/botStatusUpdate"
	memberStatusUpdate "github.com/Saime-0/tg-bot-contest/internal/usecase/member/statusUpdate"
	messageCreated "github.com/Saime-0/tg-bot-contest/internal/usecase/message/created"
)

type Controller struct {
	DB  *sqlx.DB
	Bot *gotgbot.Bot
}

type Middleware func(func(Request) error) func(Request) error

func onlyInPrivateChat(fn func(Request) error) func(Request) error {
	return func(r Request) error {
		if r.ctx.EffectiveChat != nil &&
			r.ctx.EffectiveChat.Type == gotgbot.ChatTypePrivate {
			return fn(r)
		}
		return nil
	}
}

func logDebug(fn func(Request) error) func(Request) error {
	return func(r Request) error {
		args := []any{
			slog.Int64("update_id", r.ctx.UpdateId),
		}
		if r.ctx.EffectiveChat != nil {
			args = append(args, slog.Int64("effective_chat_id", r.ctx.EffectiveChat.Id))
			args = append(args, slog.String("effective_chat_username", r.ctx.EffectiveChat.Username))
		}
		if r.ctx.EffectiveUser != nil {
			args = append(args, slog.Int64("effective_user_id", r.ctx.EffectiveUser.Id))
			args = append(args, slog.String("effective_user_username", r.ctx.EffectiveUser.Username))
		}
		if r.ctx.EffectiveMessage != nil {
			args = append(args, slog.Int64("effective_message_id", r.ctx.EffectiveMessage.MessageId))
		}

		slog.Debug("new update", args...)

		return fn(r)
	}
}

func (c *Controller) AddHandlers(dispatcher *ext.Dispatcher) error {
	baseChain := []Middleware{logDebug}
	privateChain := []Middleware{logDebug, onlyInPrivateChat}

	handlerGroup := []ext.Handler{
		c.NewCommand("contestConfigRun", contestConfigRun, privateChain),
		c.NewCommand("contestStop", contestStopHandler, privateChain),
		c.NewMyChatMember(newMyChatMember, baseChain),
		c.NewMessage(newMessage, baseChain),
		c.NewChatMember(newChatMember, baseChain),
	}

	for _, h := range handlerGroup {
		dispatcher.AddHandler(h)
	}

	return nil
}

func (c *Controller) NewCommand(com string, f func(Request) error, mws []Middleware) handlers.Command {
	return handlers.NewCommand(com, c.modulation(f, mws...))
}
func (c *Controller) NewMyChatMember(f func(Request) error, mws []Middleware) handlers.MyChatMember {
	return handlers.NewMyChatMember(nil, c.modulation(f, mws...))
}
func (c *Controller) NewMessage(f func(Request) error, mws []Middleware) handlers.Message {
	return handlers.NewMessage(nil, c.modulation(f, mws...))
}
func (c *Controller) NewChatMember(f func(Request) error, mws []Middleware) handlers.ChatMember {
	return handlers.NewChatMember(nil, c.modulation(f, mws...))
}

func newMyChatMember(r Request) (err error) {
	chat := silentUpdateChat(r)

	// Определить статус участия бота
	botMemberStatus := defineMemberStatus(
		r.ctx.MyChatMember.OldChatMember.GetStatus(), // OldStatus
		r.ctx.MyChatMember.NewChatMember.GetStatus(), // NewStatus
	)
	if botMemberStatus == 0 {
		return nil
	}

	err = inTransaction(r.DB, func(tx *sqlx.Tx) error {
		return (&botStatusUpdate.Params{
			TX:              tx,
			Chat:            chat,
			BotMemberStatus: botMemberStatus,
		}).Run()
	})
	if err != nil {
		slog.Debug("newMyChatMember: botStatusUpdate.Run: " + err.Error())
	}

	return nil
}

type Request struct {
	DB *sqlx.DB
	*gotgbot.Bot
	ctx *ext.Context
}

func (c *Controller) modulation(fn func(request Request) error, mws ...Middleware) handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		r := Request{DB: c.DB, Bot: b, ctx: ctx}
		// Собрать middleware в функцию
		for _, mw := range mws {
			fn = mw(fn)
		}
		return fn(r)
	}
}

func contestConfigRun(r Request) (err error) {
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
	params := contestCreate.Params{
		TX:        nil,
		Keyword:   kv[l10n.CfgKeyword],
		CreatorID: int(r.ctx.EffectiveSender.Id()),
	}

	// Чат (дочерний чат)
	chatUsername := clearAt(kv[l10n.CfgChatUsername])
	if chatUsername == "" {
		if params.KeywordChatID, err = getIntParameter(kv, l10n.CfgChatID, true, 0); err != nil {
			return r.reactError(err)
		}
	} else {
		// Найти чат по username
		var chat model.Chat
		if chat, err = chatTake.Run(r.DB, chatUsername); err != nil {
			return r.reactError(err)
		}
		params.KeywordChatID = chat.ID // сохранить ID в параметры
	}

	// Канал (родительский чат)
	chanUsername := clearAt(kv[l10n.CfgChannelUsername])
	if chanUsername == "" {
		if params.CompetitiveChatID, err = getIntParameter(kv, l10n.CfgChannelID, false, params.KeywordChatID); err != nil {
			return r.reactError(err)
		}
	} else {
		// Найти чат по username
		var chat model.Chat
		if chat, err = chatTake.Run(r.DB, chanUsername); err != nil {
			return r.reactError(err)
		}
		params.CompetitiveChatID = chat.ID // сохранить ID в параметры
	}

	// Проверить наличие прав админа в чате
	if err := r.checkAdminRights(int64(params.KeywordChatID)); err != nil {
		return r.reactError(err)
	}

	// Достать кратность из параметров
	if params.Multiplicity, err = getIntParameter(kv, l10n.CfgMultiplicity, true, 0); err != nil {
		return r.reactError(err)
	}

	// Достать ID топика
	if params.KeywordTopicID, err = getIntParameter(kv, l10n.CfgTopic, false, 0); err != nil {
		return r.reactError(err)
	}

	// Проверить доступность писать в чат
	if err = checkChatAvailability(r, params.KeywordChatID, params.KeywordTopicID); err != nil {
		return r.reactError(err)
	}

	if err := inTransaction(r.DB, func(tx *sqlx.Tx) error {
		params.TX = tx
		return params.Run()
	}); err != nil {
		return r.reactError(err)
	}

	return right(fastReply(r, l10n.ContestConfigRunSuccess))
}

func checkChatAvailability(r Request, chatID, topicID int) error {
	pingMsg, err := r.SendMessage(int64(chatID), "ping", &gotgbot.SendMessageOpts{
		MessageThreadId: int64(topicID),
	})
	if err != nil {
		return r.reactError(ue.New(l10n.ContestConfigBotCannotSendMsg))
	}
	if _, err = pingMsg.Delete(r.Bot, nil); err != nil {
		slog.Warn("contestConfigRun: pingMsg.Delete: " + err.Error())
	}

	return nil
}

func getIntParameter(kv map[string]string, name string, isRequired bool, defaultValue int) (int, error) {
	if isRequired && kv[name] == "" {
		return 0, &usageErrPkg.UsageErr{
			Err:   ue.New(l10n.ParameterNotProvided + ": " + name),
			Usage: l10n.ContestConfigRunUsage,
		}
	}

	val, err := strconv.ParseInt(kv[name], 10, 64)
	if err != nil && isRequired {
		return 0, &usageErrPkg.UsageErr{
			Err:   nil,
			Usage: l10n.ContestConfigRunUsage,
		}
	} else if err != nil {
		return defaultValue, nil
	}

	return int(val), nil
}

func newMessage(r Request) (err error) {
	chat := silentUpdateChat(r)

	msg := r.ctx.Message
	if !isGroup(msg.Chat) || // Выйти, если сообщение не из группы ...
		!msg.GetSender().IsUser() || // или не отправлено пользователем ...
		msg.GetText() == "" { // или оно пустое
		return nil
	}
	var topicID int
	if msg.IsTopicMessage {
		topicID = int(msg.MessageThreadId)
	}

	var messageCreatedOut messageCreated.Out
	if messageCreatedOut, err = (&messageCreated.Params{
		DB:      r.DB,
		Chat:    chat,
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

func chatIDFromChatProperty(r Request, property string) (int, error) {
	val, err := strconv.ParseInt(property, 10, 64)
	if err == nil {
		return int(val), nil
	}

	// Найти чат
	var chat model.Chat
	if chat, err = chatTake.Run(r.DB, clearAt(property)); err != nil {
		return 0, err
	}

	return chat.ID, nil
}

func contestStopHandler(r Request) (err error) {
	words := strings.Fields(r.ctx.Message.GetText())
	if len(words) < 2 {
		return right(r.ctx.Message.Reply(r.Bot, l10n.ContestStopUsage, nil))
	}

	// Определить ID чата
	var chatID int
	if chatID, err = chatIDFromChatProperty(r, words[1]); err != nil {
		return r.reactError(err)
	}

	// Проверить наличие прав админа в чате
	if err = r.checkAdminRights(int64(chatID)); err != nil {
		return r.reactError(err)
	}

	if err = inTransaction(r.DB, func(tx *sqlx.Tx) error {
		return (&contestStop.Params{TX: tx, ChatID: chatID}).Run()
	}); err != nil {
		return r.reactError(err)
	}

	return right(r.ctx.Message.Reply(r.Bot, l10n.ContestStopSuccess, nil))
}

func (r Request) checkAdminRights(chatID int64) error {
	if admins, err := r.GetChatAdministrators(chatID, nil); err != nil {
		slog.Warn("checkAdminRights: " + err.Error())
		return ue.New(l10n.CreateContestCantVerifyAdminRights)
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
	chat := silentUpdateChat(r)

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
		Chat:         chat,
		MemberStatus: memberStatus,
		Participant:  tgModel.UserDomain(participant),
		Initiator:    tgModel.UserDomain(initiator),
		ViaLink:      viaLink,
	}

	if err := inTransaction(r.DB, func(tx *sqlx.Tx) error {
		memberStatusUpdateParams.TX = tx
		return memberStatusUpdateParams.Run()
	}); err != nil {
		return err
	}

	return nil
}
