package update

import (
	"github.com/jmoiron/sqlx"

	"tgBotContest/model"
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

func Run(db *sqlx.DB, chat model.Chat) error {
	return (&Params{DB: db, Chat: chat}).Run()
}
