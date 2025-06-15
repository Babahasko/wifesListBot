package middleware

import (
	"net/http"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type Middleware func(gotgbot.BotClient) gotgbot.BotClient

func Chain(middlewares ...Middleware) Middleware {
	return func(next gotgbot.BotClient) gotgbot.BotClient {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

func NewMiddlewareClient(chain Middleware) gotgbot.BotClient{
	baseClient := &gotgbot.BaseBotClient{
		Client:             http.Client{},
		UseTestEnvironment: false,
		DefaultRequestOpts: &gotgbot.RequestOpts{
			Timeout: gotgbot.DefaultTimeout,
			APIURL:  gotgbot.DefaultAPIURL,
		},
	}

	return chain(baseClient)
}