package tg

import (
	"context"
	"fmt"

	tgClient "tgBotCompetition/tg/client"
	tgUpdatesController "tgBotCompetition/tg/updatesController"
)

func Run(ctx context.Context, token string) (err error) {
	var client *tgClient.Client
	updatesController := &tgUpdatesController.Controller{}

	if client, err = tgClient.Run(token, updatesController); err != nil {
		return fmt.Errorf("tgClient.Run: %w", err)
	}

	updaterCtx, cancel := context.WithCancel(context.Background())
	go func() {
		client.Updater.Idle()
		cancel()
	}()

	select {
	case <-ctx.Done():
		client.Updater.Stop()
		return nil
	case <-updaterCtx.Done():
		return nil
	}
}
