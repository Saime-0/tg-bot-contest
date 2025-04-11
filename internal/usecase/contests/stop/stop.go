package stop

import (
	"github.com/jmoiron/sqlx"

	"github.com/Saime-0/tg-bot-contest/internal/l10n"
	"github.com/Saime-0/tg-bot-contest/internal/ue"
)

type Params struct {
	TX *sqlx.Tx

	ChatID int
}

var (
	ErrContestNotFound = ue.New(l10n.ContestStopNotFound)
)

func (p *Params) Run() error {
	if res, err := p.TX.Exec(`
		update contests 
		set ended_at = current_timestamp
		where ended_at is null
			and competitive_chat_id=?
	`, p.ChatID); err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected == 0 {
		return ErrContestNotFound
	}

	return nil
}
