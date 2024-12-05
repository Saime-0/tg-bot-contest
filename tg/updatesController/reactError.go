package updatesController

import (
	"errors"
	"log/slog"

	"tgBotContest/l10n"
	"tgBotContest/ue"
)

func (r Request) reactError(err error) error {
	if r.ctx.Message == nil {
		return err
	}

	var userErr *ue.Err
	if errors.As(err, &userErr) {
		text := l10n.ReactErrorPrefix + err.Error() + l10n.ReactErrorSuffix
		slog.Warn(text)
		if _, replErr := r.ctx.Message.Reply(r.Bot, text, nil); replErr != nil {
			slog.Error("reactError: reply to userErr: " + replErr.Error())
		}
		return nil
	}

	return err
}
