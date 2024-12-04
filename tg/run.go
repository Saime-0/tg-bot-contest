package tg

import (
	"context"
	"fmt"
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/jmoiron/sqlx"

	tgClient "tgBotContest/tg/client"
	tgUpdatesController "tgBotContest/tg/updatesController"
)

func Run(ctx context.Context, token string, db *sqlx.DB) (err error) {
	updatesController := &tgUpdatesController.Controller{
		DB: db,
	}

	client := &tgClient.Client{}
	if client.Bot, err = tgClient.NewBot(token); err != nil {
		return fmt.Errorf("can't create bot: %w", err)
	}
	if err = updatesController.SetMyCommands(client.Bot); err != nil {
		return fmt.Errorf("can't set commands: %w", err)
	}

	client.Dispatcher = tgClient.NewDispatcher()
	if err = updatesController.AddHandlers(client.Dispatcher); err != nil {
		return fmt.Errorf("can't register updatesController: %w", err)
	}

	client.Updater = ext.NewUpdater(client.Dispatcher, nil)

	if err = tgClient.StartPolling(client.Updater, client.Bot); err != nil {
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
