package client

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// loggingClient это обертка для стандартного BotClient, которая добавляет функциональность ведения журнала.
// Он регистрирует ответы на вызовы Telegram API для целей отладки.
type loggingClient struct {
	gotgbot.BotClient // Embedded interface for delegation
}

// RequestWithContext переопределяет стандартный метод RequestWithContext, чтобы добавить ведение журнала.
// Он пересылает запрос базовому BotClient и регистрирует ответ.
func (b loggingClient) RequestWithContext(ctx context.Context, token string, method string, params map[string]string, data map[string]gotgbot.FileReader, opts *gotgbot.RequestOpts) (json.RawMessage, error) {
	rm, err := b.BotClient.RequestWithContext(ctx, token, method, params, data, opts)
	b1, _ := json.Marshal(rm)
	if string(b1) != ("[]") {
		slog.Debug(string(b1))
	}
	return rm, err
}

// newLoggingClient создает новый экземпляр клиента ведения журнала.
// Он инициализирует клиент с настройками по умолчанию и добавляет в него функциональность ведения журнала.
func newLoggingClient() loggingClient {
	return loggingClient{
		BotClient: &gotgbot.BaseBotClient{
			Client:             http.Client{},
			UseTestEnvironment: false,
			DefaultRequestOpts: &gotgbot.RequestOpts{
				Timeout: gotgbot.DefaultTimeout,
				APIURL:  gotgbot.DefaultAPIURL,
			},
		},
	}
}
