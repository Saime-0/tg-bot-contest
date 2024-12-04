package take

import (
	"github.com/jmoiron/sqlx"

	"tgBotCompetition/model"
)

type Params struct {
	DB *sqlx.DB

	Username string
}

func (p *Params) Run() (model.Chat, error) {
	var chat model.Chat
	return chat, p.DB.Get(&chat, "select * from chats where username = ?", p.Username)
}
