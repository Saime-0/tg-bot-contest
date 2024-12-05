package update

import (
	"database/sql"

	"tgBotContest/model"
)

type Params struct {
	TXDB txOrDB

	Chat model.Chat
}

func (p *Params) Run() error {
	_, err := p.TXDB.NamedExec(`
		insert into chats (id, title, username)
		values (:id, :title, :username)
		on conflict (id) do update set 
		   title = excluded.title,
		   username = excluded.username
    `, p.Chat)

	return err
}

func Run(db txOrDB, chat model.Chat) error {
	return (&Params{TXDB: db, Chat: chat}).Run()
}

type txOrDB interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
}
