package update

import (
	"database/sql"

	"github.com/Saime-0/tg-bot-contest/internal/model"
)

func Run(txdb txOrDB, chat model.Chat) error {
	_, err := txdb.NamedExec(`
		insert into chats (id, title, username, child_id, parent_id)
		values (:id, :title, :username, :child_id, :parent_id)
		on conflict (id) do update set 
		   title = excluded.title,
		   username = excluded.username,
		   child_id = excluded.child_id,
		   parent_id = excluded.parent_id
    `, chat)

	return err
}

type txOrDB interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
}
