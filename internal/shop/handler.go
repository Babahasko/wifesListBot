package shop

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

type ShopHandler struct {
}

func NewShopHandler(router *ext.Dispatcher) {
	handler := &ShopHandler{}
	router.AddHandler(handlers.NewCommand("start", handler.Start))
}

func (handler *ShopHandler) Start(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Chat.SendMessage(b, MsgStart, nil) //&gotgbot.SendMessageOpts{ParseMode: "MarkdownV2"}
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	return nil
}