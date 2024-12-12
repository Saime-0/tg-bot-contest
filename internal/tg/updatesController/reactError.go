package updatesController

import (
	"errors"
	"log/slog"

	"github.com/Saime-0/tg-bot-contest/internal/l10n"
	usageErrPkg "github.com/Saime-0/tg-bot-contest/internal/tg/usageErr"
	"github.com/Saime-0/tg-bot-contest/internal/ue"
)

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
