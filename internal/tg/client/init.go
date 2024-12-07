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

func Run(token string, regHand registerHandlers) (_ *Client, err error) {
	client := &Client{}

	if client.Bot, err = NewBot(token); err != nil {
		return nil, fmt.Errorf("can't create bot: %w", err)
	}

	client.Dispatcher = NewDispatcher()
	if err = regHand.AddHandlers(client.Dispatcher); err != nil {
		return nil, fmt.Errorf("can't register updatesController: %w", err)
	}

	client.Updater = ext.NewUpdater(client.Dispatcher, nil)

	if err = StartPolling(client.Updater, client.Bot); err != nil {
		return nil, err
	}

	return client, nil
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

type registerHandlers interface {
	AddHandlers(*ext.Dispatcher) error
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
