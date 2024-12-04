package update

import (
	"github.com/jmoiron/sqlx"

	"tgBotCompetition/model"
)

type Params struct {
	DB *sqlx.DB

	Chat model.Chat
}

func (p *Params) Run() error {
	_, err := p.DB.NamedExec(`
		insert or replace into chats (id, title, username)
		values (:id, :title, :username)
    `, p.Chat)

	return err
}
