package client

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// StartPolling инициализирует и запускает механизм получения обновлений бота.
// Он настраивает бота на прослушивание определенных типов обновлений с соответствующими таймаутами.
func StartPolling(updater *ext.Updater, bot *gotgbot.Bot) error {
	err := updater.StartPolling(bot, &ext.PollingOpts{
		DropPendingUpdates: true, // Игнорировать любые обновления, которые произошли, пока бот был офлайн
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			// Слушать только эти конкретные типы обновлений:
			AllowedUpdates: []string{
				"chat_member",    // Обновления об изменениях статуса участника чата
				"message",        // Новые сообщения
				"my_chat_member", // Обновления о статусе самого бота в чате
			},
			Timeout: 5, // Ожидать до 5 секунд для получения обновлений в режиме long polling
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10, // Таймаут HTTP запроса
			},
		},
	})
	if err != nil {
		return fmt.Errorf("updater.StartPolling: %w", err)
	}

	slog.Info(bot.Username + " has been started...")

	return nil
}
