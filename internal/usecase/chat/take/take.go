// Package take
// так как поиск осуществляется по изменяемому полю username,
// перед этим необходимо данные синхронизировать, т.е внесни в БД
// актуальный username искомого чата.
package take

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/Saime-0/tg-bot-contest/internal/l10n"
	"github.com/Saime-0/tg-bot-contest/internal/model"
	"github.com/Saime-0/tg-bot-contest/internal/ue"
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
