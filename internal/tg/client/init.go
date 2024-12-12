package client

import (
	"fmt"
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type Client struct {
	Bot        *gotgbot.Bot
	Dispatcher *ext.Dispatcher
	Updater    *ext.Updater
}

func NewBot(token string) (*gotgbot.Bot, error) {
	bot, err := gotgbot.NewBot(token, &gotgbot.BotOpts{
		BotClient: newLoggingClient(),
	})
	if err != nil {
		return nil, fmt.Errorf("gotgbot.NewBot: %w", err)
	}

	return bot, nil
}

func NewDispatcher() *ext.Dispatcher {
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	return dispatcher
}
