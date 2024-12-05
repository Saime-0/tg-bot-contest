package update

import (
	"github.com/jmoiron/sqlx"

	"tgBotContest/model"
)

type Params struct {
	TX *sqlx.Tx

	User model.User
}

func (p *Params) Run() error {
	_, err := p.TX.NamedExec(`
		insert into users (id,is_bot,first_name,username)
		values (:id,:is_bot,:first_name,:username)
		on conflict(id) do update set
			is_bot = excluded.is_bot,
			first_name = excluded.first_name,
			username = excluded.username,
			updated_at = current_timestamp
    `, p.User)

	return err
}

func Run(tx *sqlx.Tx, initiator model.User) error {
	return (&Params{TX: tx, User: initiator}).Run()
}
