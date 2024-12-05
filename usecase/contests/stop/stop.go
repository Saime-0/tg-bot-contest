package stop

import (
	"github.com/jmoiron/sqlx"

	"tgBotContest/l10n"
	"tgBotContest/ue"
)

type Params struct {
	TX *sqlx.Tx

	ChatID int
}

func (p *Params) Run() error {
	if res, err := p.TX.Exec(`
		update contests 
		set ended_at = current_timestamp
		where ended_at is null
			and chat_id=?
	`, p.ChatID); err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected == 0 {
		return ue.New(l10n.ContestStopNotFound)
	}

	return nil
}
