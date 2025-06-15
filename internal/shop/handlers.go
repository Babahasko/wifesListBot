package shop

import (
	"fmt"
	"html"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/conversation"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

type ShopHandler struct {
}

const (
	CATEGORY = "category"
)

func NewShopHandler(router *ext.Dispatcher) {
	handler := &ShopHandler{}
	router.AddHandler(handlers.NewCommand("start", handler.Start))
	// router.AddHandler(handlers.NewMessage(message.Equal(ButtonAddPurchase), handler.AddPurchase))
	router.AddHandler(handlers.NewConversation(
		[]ext.Handler{handlers.NewMessage(message.Equal(ButtonAddPurchase), handler.AddPurchase)},
		map[string][]ext.Handler{
			CATEGORY: {handlers.NewMessage(isCategory, category)},
		},
		&handlers.ConversationOpts{
			Exits:        []ext.Handler{handlers.NewMessage(message.Equal(ButtonCancel), cancel)},
			StateStorage: conversation.NewInMemoryStorage(conversation.KeyStrategySenderAndChat),
			AllowReEntry: true,
		},
	))
}

func (handler *ShopHandler) Start(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Chat.SendMessage(b, MsgStart, &gotgbot.SendMessageOpts{
		ReplyMarkup: getMainMenueKeyboard(),
	}) //&gotgbot.SendMessageOpts{ParseMode: "MarkdownV2"}
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	return nil
}

func (handler *ShopHandler) AddPurchase(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Chat.SendMessage(b, "Пошла возня добавления покупки", &gotgbot.SendMessageOpts{
		ReplyMarkup: getShortMenueKeyboard(),
	})
	if err != nil {
		return fmt.Errorf("failed to send menue message: %w", err)
	}
	return handlers.NextConversationState(CATEGORY)
}

//TODO: 1. Refactor IsCategory function
func isCategory(msg *gotgbot.Message) bool {
	return message.Text(msg) && !message.Command(msg)
}

func cancel(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Chat.SendMessage(b, "Жаль, что вы прервались!", &gotgbot.SendMessageOpts{
		ParseMode: "html",
		ReplyMarkup: getMainMenueKeyboard(),
	})
	if err != nil {
		return fmt.Errorf("failed to send cancel message: %w", err)
	}
	return handlers.EndConversation()
}

func category(b *gotgbot.Bot, ctx *ext.Context) error {
	inputCategory := ctx.EffectiveMessage.Text
	_, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Отлично, вы выбрали категорию %s!\n\n", html.EscapeString(inputCategory)), &gotgbot.SendMessageOpts{
		ParseMode: "html",
	})
	if err != nil {
		return fmt.Errorf("failed to send category message: %w", err)
	}
	return handlers.EndConversation()
}

