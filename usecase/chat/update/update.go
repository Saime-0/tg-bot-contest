package update

import (
	"github.com/jmoiron/sqlx"

	"tgBotContest/model"
)

type Params struct {
	TX *sqlx.Tx

	Chat model.Chat
}

func (p *Params) Run() error {
	_, err := p.TX.NamedExec(`
		insert into chats (id, title, username)
		values (:id, :title, :username)
		on conflict (id) do update set 
		   title = excluded.title,
		   username = excluded.username
    `, p.Chat)

	return err
}

func Run(tx *sqlx.Tx, chat model.Chat) error {
	return (&Params{TX: tx, Chat: chat}).Run()
}
