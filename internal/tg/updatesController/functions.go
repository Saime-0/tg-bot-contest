package updatesController

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/jmoiron/sqlx"

	"github.com/Saime-0/tg-bot-contest/internal/l10n"
	tgModel "github.com/Saime-0/tg-bot-contest/internal/tg/model"
	usageErrPkg "github.com/Saime-0/tg-bot-contest/internal/tg/usageErr"
	"github.com/Saime-0/tg-bot-contest/internal/ue"
	chatUpdate "github.com/Saime-0/tg-bot-contest/internal/usecase/chat/update"
)

func clearAt(s string) string {
	return strings.TrimPrefix(s, "@")
}

func right[A, B any](_ A, b B) B {
	return b
}

func inTransaction(db *sqlx.DB, f func(tx *sqlx.Tx) error) (err error) {
	var tx *sqlx.Tx
	if tx, err = db.BeginTxx(context.Background(), nil); err != nil {
		return err
	}

	defer func() { // если во время выполнения действий случилась паника, то откатываем транзакцию
		if r := recover(); r != nil {
			if err := tx.Rollback(); err != nil {
				slog.Error("tx.Rollback: " + err.Error())
			}
		}
	}()

	if err = f(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			slog.Error("tx.Rollback: " + err.Error())
		}
		return err
	}

	return tx.Commit()
}

// silentUpdateChat втихую обновляет чат
func silentUpdateChat(r Request) {
	if r.ctx.EffectiveChat == nil {
		return
	}
	chat := tgModel.ChatDomain(*r.ctx.EffectiveChat)
	if err := chatUpdate.Run(r.DB, chat); err != nil {
		slog.Warn("newMessage: chatUpdate.Run: " + err.Error())
	}
}

func fastMDReply(r Request, msg string) (*gotgbot.Message, error) {
	return r.ctx.Message.Reply(r.Bot, msg, &gotgbot.SendMessageOpts{
		ParseMode: gotgbot.ParseModeMarkdownV2,
	})
}

func fastReply(r Request, msg string) (*gotgbot.Message, error) {
	return r.ctx.Message.Reply(r.Bot, msg, nil)
}

func isGroup(chat gotgbot.Chat) bool {
	return chat.Type == gotgbot.ChatTypeSupergroup ||
		chat.Type == gotgbot.ChatTypeGroup
}

func (r Request) reactError(err error) error {
	if r.ctx.Message == nil {
		return err
	}
	if err == nil {
		return nil
	}

	var userErr *ue.Err
	if errors.As(err, &userErr) {
		text := l10n.ReactErrorPrefix + userErr.Error() + l10n.ReactErrorSuffix
		slog.Warn(text)

		if _, replErr := fastReply(r, text); replErr != nil {
			slog.Error("reactError: reply to userErr: " + replErr.Error())
		}
		return nil
	}

	var usageErr *usageErrPkg.UsageErr
	if errors.As(err, &usageErr) {
		if _, replErr := fastMDReply(r, usageErr.Usage); replErr != nil {
			slog.Error("reactError: reply with usage: " + replErr.Error())
		}
		return r.reactError(usageErr.Err)
	}

	return err
}
