package take

import (
	"github.com/jmoiron/sqlx"

	"tgBotContest/model"
)

type Params struct {
	TX *sqlx.Tx

	Username string
}

func (p *Params) Run() (model.Chat, error) {
	var chat model.Chat
	return chat, p.TX.Get(&chat, "select * from chats where username = ?", p.Username)
}

func Run(tx *sqlx.Tx, username string) (model.Chat, error) {
	return (&Params{
		TX:       tx,
		Username: username,
	}).Run()
}
