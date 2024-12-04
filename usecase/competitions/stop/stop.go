package stop

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type Params struct {
	DB *sqlx.DB

	ChatUsername string
}

func (p *Params) Run() error {
	if res, err := p.DB.Exec(`
		update competitions 
		set ended_at = current_timestamp
		where ended_at is null
			and exists (select 1 from chats where username=?)
	`, p.ChatUsername); err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected == 0 {
		return errors.New("competition not found")
	}

	return nil
}
