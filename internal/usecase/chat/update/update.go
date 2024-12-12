package update

import (
	"database/sql"

	"github.com/Saime-0/tg-bot-contest/internal/model"
)

func Run(txdb txOrDB, chat model.Chat) error {
	_, err := txdb.NamedExec(`
		insert into chats (id, title, username)
		values (:id, :title, :username)
		on conflict (id) do update set 
		   title = excluded.title,
		   username = excluded.username
    `, chat)

	return err
}

type txOrDB interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
}
