package take

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"tgBotContest/l10n"
	"tgBotContest/model"
	"tgBotContest/ue"
)

type Params struct {
	DB *sqlx.DB

	Username string
}

func (p *Params) Run() (model.Chat, error) {
	var chat model.Chat
	err := p.DB.Get(&chat, "select * from chats where username = ?", p.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return chat, errors.Join(err, ue.New(l10n.ChatTakeNotFound))
	}

	return chat, err
}

func Run(db *sqlx.DB, username string) (model.Chat, error) {
	return (&Params{
		DB:       db,
		Username: username,
	}).Run()
}
