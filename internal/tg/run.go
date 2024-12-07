package tg

import (
	"context"
	"fmt"
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/jmoiron/sqlx"

	client2 "github.com/Saime-0/tg-bot-contest/internal/tg/client"
	tgUpdatesController "github.com/Saime-0/tg-bot-contest/internal/tg/updatesController"
)

func Run(ctx context.Context, token string, db *sqlx.DB) (err error) {
	updatesController := &tgUpdatesController.Controller{
		DB: db,
	}

	client := &client2.Client{}
	if client.Bot, err = client2.NewBot(token); err != nil {
		return fmt.Errorf("can't create bot: %w", err)
	}

	client.Dispatcher = client2.NewDispatcher()
	if err = updatesController.AddHandlers(client.Dispatcher); err != nil {
		return fmt.Errorf("can't register updatesController: %w", err)
	}

	client.Updater = ext.NewUpdater(client.Dispatcher, nil)

	if err = client2.StartPolling(client.Updater, client.Bot); err != nil {
		return err
	}

	updaterCtx, cancel := context.WithCancel(context.Background())
	go func() {
		client.Updater.Idle()
		cancel()
	}()

	select {
	case <-ctx.Done():
		cancel()
		if err = client.Updater.Stop(); err != nil {
			log.Println("WARNING", err.Error())
		}
		return nil
	case <-updaterCtx.Done():
		return nil
	}
}
