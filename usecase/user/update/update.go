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
		INSERT OR REPLACE INTO users (id,is_bot,first_name,username)
		VALUES (:id,:is_bot,:first_name,:username)
    `, p.User)

	return err
}
