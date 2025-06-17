package middleware

import (
	"context"
	"encoding/json"
	"shopping_bot/pkg/logger"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type loggingClient struct {
	gotgbot.BotClient
}

func NewLoggingMiddleware() func(gotgbot.BotClient) gotgbot.BotClient {
	return func(next gotgbot.BotClient) gotgbot.BotClient {
		return &loggingClient{
			BotClient: next,
		}
	}
}

func (l *loggingClient) RequestWithContext(
	ctx context.Context,
	token string,
	method string,
	params map[string]string,
	data map[string]gotgbot.FileReader,
	opts *gotgbot.RequestOpts,
) (json.RawMessage, error) {
	if method != "getUpdates" {
		logger.Sugar.Infof("Call method: %s, Parameters: %+v", method, params)
	}
	return l.BotClient.RequestWithContext(ctx, token, method, params, data, opts)
}
