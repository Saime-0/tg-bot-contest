package tg

import (
	"context"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/jmoiron/sqlx"

	tgClient "tgBotCompetition/tg/client"
	tgUpdatesController "tgBotCompetition/tg/updatesController"
)

func Run(ctx context.Context, token string, db *sqlx.DB) (err error) {
	client := &tgClient.Client{}

	if client.Bot, err = tgClient.NewBot(token); err != nil {
		return fmt.Errorf("can't create bot: %w", err)
	}
	updatesController := &tgUpdatesController.Controller{
		DB: db,
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
		client.Updater.Stop()
		return nil
	case <-updaterCtx.Done():
		return nil
	}
}
