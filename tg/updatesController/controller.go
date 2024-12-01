package updatesController

import (
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

type Controller struct{}

//type requiredUsecases interface {
//	MemberLeft() error
//}

func (c *Controller) AddHandlers(dispatcher *ext.Dispatcher) error {
	//dispatcher.AddHandler(handlers.NewMessage(nil, func(b *gotgbot.Bot, ctx *ext.Context) error {
	//	if len(ctx.Message.NewChatMembers) > 0 {
	//		return c.newMessageNewChatMembers(b, ctx)
	//	}
	//	return nil
	//}))

	dispatcher.AddHandler(handlers.NewChatMember(nil, c.newChatMember))

	return nil
}
