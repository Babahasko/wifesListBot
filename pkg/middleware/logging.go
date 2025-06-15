package middleware

import (
	"context"
	"encoding/json"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"go.uber.org/zap"
)

type loggingClient struct {
	gotgbot.BotClient
	sugar *zap.SugaredLogger
}

func NewLoggingMiddleware(sugar *zap.SugaredLogger) func(gotgbot.BotClient) gotgbot.BotClient {
	return func(next gotgbot.BotClient) gotgbot.BotClient {
		return &loggingClient{
			BotClient: next,
			sugar: sugar,
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
		l.sugar.Infof("Call method: %s, Parameters: %+v", method, params)
	}
	return l.BotClient.RequestWithContext(ctx, token, method, params, data, opts)
}
