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
		insert into chats (id, title, username)
		values (:id, :title, :username)
		on conflict (id) do update set 
		   title = excluded.title,
		   username = excluded.username
    `, p.Chat)

	return err
}
