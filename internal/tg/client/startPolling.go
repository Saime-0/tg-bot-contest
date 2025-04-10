package client

import (
	"fmt"
	"log"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func StartPolling(updater *ext.Updater, bot *gotgbot.Bot) error {
	err := updater.StartPolling(bot, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			AllowedUpdates: []string{
				"chat_member",
				"message",
				"my_chat_member",
			},
			Timeout: 5,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to start polling: %w", err)
	}

	log.Printf("%s has been started...\n", bot.Username)

	return nil
}
