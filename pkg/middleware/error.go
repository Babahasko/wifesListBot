package middleware

import (
	"context"
	"encoding/json"
	"shopping_bot/pkg/logger"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type errorHandlingClient struct {
	gotgbot.BotClient
}

func NewErrorMiddleware() func(gotgbot.BotClient) gotgbot.BotClient {
	return func(next gotgbot.BotClient) gotgbot.BotClient {
		return &errorHandlingClient{
			BotClient: next,
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
		logger.Sugar.Infof("Error in method %s: %v", method, err)
	}
	return resp, err
}