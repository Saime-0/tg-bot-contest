package take

import (
	"github.com/jmoiron/sqlx"

	"tgBotContest/model"
)

type Params struct {
	DB *sqlx.DB

	Username string
}

func (p *Params) Run() (model.Chat, error) {
	var chat model.Chat
	return chat, p.DB.Get(&chat, "select * from chats where username = ?", p.Username)
}

func Run(db *sqlx.DB, username string) (model.Chat, error) {
	return (&Params{
		DB:       db,
		Username: username,
	}).Run()
}
