package shop

import (
	"fmt"
	// "html"
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/conversation"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
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
			CATEGORY: {handlers.NewCallback(callbackquery.Prefix("cat_"), category)},
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
	categories := getUserCategories(123)
	catKeyboard := getCategoriesKeyboard(categories)

	_, err := ctx.EffectiveMessage.Chat.SendMessage(b, "Выберите категорию", &gotgbot.SendMessageOpts{
		ReplyMarkup: catKeyboard,
	})
	if err != nil {
		return fmt.Errorf("failed to send menue message: %w", err)
	}
	return handlers.NextConversationState(CATEGORY)
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
	log.Println("зашли в обработчик категории")
	cb := ctx.Update.CallbackQuery
	text := ctx.Update.CallbackQuery.Data
	_, err := cb.Answer(b,  &gotgbot.AnswerCallbackQueryOpts{
		Text: text,
		ShowAlert: false,
	})
	if err != nil {
		return fmt.Errorf("failed to send category message %w", err)
	}
	_, _, err = cb.Message.EditText(b, "Введите название", &gotgbot.EditMessageTextOpts{
		ParseMode: "html",
	})

	if err != nil {
		return fmt.Errorf("failed to send add name message %w", err)
	}

	_, err = ctx.EffectiveMessage.Chat.SendMessage(b, "", &gotgbot.SendMessageOpts{
		ReplyMarkup: getMenueKeyboard(),
	})

	if err != nil {
		return fmt.Errorf("failed to send keyboard %w", err)
	}
	
	// _, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Отлично, вы выбрали категорию %s!\n\n", html.EscapeString(inputCategory)), &gotgbot.SendMessageOpts{
	// 	ParseMode: "html",
	// })
	// if err != nil {
	// 	return fmt.Errorf("failed to send category message: %w", err)
	// }
	return handlers.EndConversation()
}

