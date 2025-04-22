package client

import (
	"fmt"
	"log/slog"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// Client представляет полноценный клиент Telegram бота со всеми необходимыми компонентами.
// Он инкапсулирует экземпляр бота, диспетчер для обработки обновлений и обновлятель
// для получения обновлений от Telegram API.
type Client struct {
	Bot        *gotgbot.Bot    // Экземпляр Telegram бота
	Dispatcher *ext.Dispatcher // Обрабатывает маршрутизацию обновлений к соответствующим обработчикам
	Updater    *ext.Updater    // Управляет процессом получения обновлений от Telegram
}

// NewBot создает новый экземпляр Telegram бота с предоставленным токеном.
// Он использует пользовательский клиент логирования для включения отладочного логирования ответов API.
func NewBot(token string) (*gotgbot.Bot, error) {
	bot, err := gotgbot.NewBot(token, &gotgbot.BotOpts{
		BotClient: newLoggingClient(),
	})
	if err != nil {
		return nil, fmt.Errorf("gotgbot.NewBot: %w", err)
	}

	return bot, nil
}

// NewDispatcher создает новый диспетчер для обработки обновлений бота.
// Он настраивает обработку ошибок и устанавливает максимальное количество одновременных подпрограмм.
func NewDispatcher() *ext.Dispatcher {
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// Обработчик ошибок логирует ошибки и продолжает обработку обновлений
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			slog.Error("произошла ошибка при обработке обновления: " + err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines, // Использовать настройки параллелизма по умолчанию
	})

	return dispatcher
}
