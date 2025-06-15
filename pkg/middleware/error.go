package middleware

import (
	"context"
	"encoding/json"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"go.uber.org/zap"
)

type errorHandlingClient struct {
	gotgbot.BotClient
	sugar *zap.SugaredLogger
}

func NewErrorMiddleware(sugar *zap.SugaredLogger) func(gotgbot.BotClient) gotgbot.BotClient {
	return func(next gotgbot.BotClient) gotgbot.BotClient {
		return &errorHandlingClient{
			BotClient: next,
			sugar: sugar,
		}
	}
}

func (e *errorHandlingClient) RequestWithContext(
	ctx context.Context,
	token string,
	method string,
	params map[string]string,
	data map[string]gotgbot.FileReader,
	opts *gotgbot.RequestOpts,
) (json.RawMessage, error) {
	resp, err := e.BotClient.RequestWithContext(ctx, token, method, params, data, opts)
	if err != nil {
		e.sugar.Infof("Error in method %s: %v", method, err)
	}
	return resp, err
}