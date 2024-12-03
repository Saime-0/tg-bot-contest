package added

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
		INSERT OR REPLACE INTO chats (id, title, username)
		VALUES (:id, :title, :username)
    `, p.Chat)

	return err
}
