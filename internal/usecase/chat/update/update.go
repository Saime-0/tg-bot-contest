package update

import (
	"database/sql"

	"github.com/Saime-0/tg-bot-contest/internal/model"
)

func Run(txdb txOrDB, chat model.Chat) error {
	_, err := txdb.NamedExec(`
		insert into chats (id, title, username, linked_chat_id)
		values (:id, :title, :username, :linked_chat_id)
		on conflict (id) do update set 
		   title = excluded.title,
		   username = excluded.username,
		   linked_chat_id = excluded.linked_chat_id
    `, chat)

	return err
}

type txOrDB interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
}
