package updatesController

import (
	"errors"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"tgBotContest/l10n"
	"tgBotContest/ue"
)

func (c *Controller) reactError(err error, b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.Message == nil {
		return err
	}

	var userErr *ue.Err
	if errors.As(err, &userErr) {
		text := l10n.ReactErrorPrefix + err.Error() + l10n.ReactErrorSuffix
		slog.Warn(text)
		if _, replErr := ctx.Message.Reply(b, text, nil); replErr != nil {
			slog.Error("reactError: reply to userErr: " + replErr.Error())
		}
		return nil
	}

	return err
}
