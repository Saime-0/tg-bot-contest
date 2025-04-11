package update

import (
	"database/sql"

	"github.com/Saime-0/tg-bot-contest/internal/model"
)

type Params struct {
	TX interface {
		NamedExec(query string, arg interface{}) (sql.Result, error)
	}
	Linking model.ChatLinking
}

func (p *Params) Run() error {
	_, err := p.TX.NamedExec(`
		insert into chats (parent_id, child_id)
		values (:parent_id, :child_id)
		on conflict (id) do update set
			parent_id = excluded.parent_id
			child_id = excluded.child_id
    `, p.Linking)

	return err
}
