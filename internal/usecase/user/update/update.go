package update

import (
	"database/sql"

	"tgBotContest/internal/model"
)

type Params struct {
	TXDB txOrDB

	User model.User
}

type txOrDB interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
}

func (p *Params) Run() error {
	_, err := p.TXDB.NamedExec(`
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

func Run(db txOrDB, initiator model.User) error {
	return (&Params{TXDB: db, User: initiator}).Run()
}
