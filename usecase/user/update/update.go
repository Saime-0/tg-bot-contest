package update

import (
	"github.com/jmoiron/sqlx"

	"tgBotCompetition/model"
)

type Params struct {
	DB *sqlx.DB

	User model.User
}

func (p *Params) Run() error {
	_, err := p.DB.NamedExec(`
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

func Run(db *sqlx.DB, initiator model.User) error {
	return (&Params{DB: db, User: initiator}).Run()
}
