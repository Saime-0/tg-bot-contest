package updatesController

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/jmoiron/sqlx"
)

type Controller struct {
	DB *sqlx.DB
}

func (c *Controller) AddHandlers(dispatcher *ext.Dispatcher) error {
	dispatcher.AddHandler(handlers.NewChatMember(nil, c.newChatMember))

	return nil
}
