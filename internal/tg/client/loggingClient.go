package client

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type loggingClient struct {
	gotgbot.BotClient
}

func (b loggingClient) RequestWithContext(ctx context.Context, token string, method string, params map[string]string, data map[string]gotgbot.FileReader, opts *gotgbot.RequestOpts) (json.RawMessage, error) {
	rm, err := b.BotClient.RequestWithContext(ctx, token, method, params, data, opts)
	b1, _ := json.MarshalIndent(rm, "", "\t")
	if string(b1) != ("[]") {
		slog.Debug(string(b1))
	}
	return rm, err
}

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
