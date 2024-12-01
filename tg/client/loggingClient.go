package client

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type loggingClient struct {
	gotgbot.BotClient
}

func (b loggingClient) RequestWithContext(ctx context.Context, token string, method string, params map[string]string, data map[string]gotgbot.FileReader, opts *gotgbot.RequestOpts) (json.RawMessage, error) {
	rawMessage, err := b.BotClient.RequestWithContext(ctx, token, method, params, data, opts)
	return rawMessage, err
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
