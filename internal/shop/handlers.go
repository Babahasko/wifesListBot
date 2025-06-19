package shop

import (
	"fmt"
	"shopping_bot/pkg/callback"
	"shopping_bot/pkg/logger"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/conversation"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

type ShopHandler struct {
	CallbackRegistry *callback.Registry
	Client           *ShopClient
}

const (
	NAME      = "list_name"
	PURCHASES = "add_purchases"
	FINISH    = "finish_form_list"
)

func NewShopHandler(router *ext.Dispatcher) {
	// Handlers for shop
	handler := &ShopHandler{
		//Add callback registry
		CallbackRegistry: callback.NewRegistry(),
		Client:           &ShopClient{},
	}
	// Register Callbacks
	handler.CallbackRegistry.Register(NewCategoryCallback)

	router.AddHandler(handlers.NewCommand("start", handler.start))
	// router.AddHandler(handlers.NewMessage(message.Equal(ButtonAddPurchase), handler.AddPurchase))
	router.AddHandler(handlers.NewConversation(
		[]ext.Handler{handlers.NewMessage(message.Equal(ButtonAddPurchase), handler.formList)},
		map[string][]ext.Handler{
			NAME:      {handlers.NewMessage(message.Text, handler.addName)},
			PURCHASES: {handlers.NewMessage(message.Text, handler.addPurchase)},
			FINISH:    {handlers.NewMessage(message.Equal(ButtonFinishList), handler.finish)},
		},
		&handlers.ConversationOpts{
			Exits:        []ext.Handler{handlers.NewMessage(message.Equal(ButtonCancel), handler.cancel)},
			StateStorage: conversation.NewInMemoryStorage(conversation.KeyStrategySenderAndChat),
			AllowReEntry: true,
		},
	))
}

func (handler *ShopHandler) start(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Chat.SendMessage(b, MsgStart, &gotgbot.SendMessageOpts{
		ReplyMarkup: getMainMenueKeyboard(),
	}) //&gotgbot.SendMessageOpts{ParseMode: "MarkdownV2"}
	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}
	return nil
}

func (handler *ShopHandler) formList(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Chat.SendMessage(b, "Введите название списка", &gotgbot.SendMessageOpts{
		ReplyMarkup: getCancelKeyboard(),
	})
	if err != nil {
		return fmt.Errorf("failed to send menue message: %w", err)
	}
	return handlers.NextConversationState(NAME)
}
func (handler *ShopHandler) addName(b *gotgbot.Bot, ctx *ext.Context) error {
	listName := ctx.EffectiveMessage.Text
	handler.Client.addShoppingList(ctx, listName)    // это у нас в базу летит shopping list
	handler.Client.SetCurrentListName(ctx, listName) // это у нас в кэш летит состояние пользователя
	logger.Sugar.Debugw("текущий список пользователя: %v", handler.Client.GetCurrentList(ctx))
	_, err := ctx.EffectiveMessage.Chat.SendMessage(b, "Введите название покупки", &gotgbot.SendMessageOpts{
		ReplyMarkup: getFormListKeyboard(),
	})
	if err != nil {
		return fmt.Errorf("failed to send addPurchase message")
	}
	return handlers.NextConversationState(PURCHASES)
}

// TODO: Add validation for item name`s`
// clear user list when finished
func (handler *ShopHandler) addPurchase(b *gotgbot.Bot, ctx *ext.Context) error {
	itemName := ctx.EffectiveMessage.Text
	// Получаем имя текущего списка юзера
	listName := handler.Client.GetCurrentList(ctx)
	if listName == "" {
		logger.Sugar.Errorw("current list name not found in user state", "userID", ctx.EffectiveUser.Id)
		_, err := ctx.EffectiveMessage.Chat.SendMessage(b, "Ошибка не задан текущий список, начните заново", nil)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
	}
	if itemName == ButtonFinishList {
		return handler.finish(b, ctx)
	}
	handler.Client.AddItemToShoppingList(ctx, listName, itemName)
	logger.Sugar.Debugw("add item: %v to shopping list: %v", itemName, listName)

	return handlers.NextConversationState(PURCHASES)
}

func (handler *ShopHandler) finish(b *gotgbot.Bot, ctx *ext.Context) error {
	listName := handler.Client.GetCurrentList(ctx)
	if listName == "" {
		logger.Sugar.Errorw("Current list name not found in user states", "userID", ctx.EffectiveUser.Id)
		_, err := ctx.EffectiveMessage.Chat.SendMessage(b, "Ошибка: список покупок не задан. Начните заново", nil)
		return fmt.Errorf("failed to send message: %w", err)
	}
	listItems, err := handler.Client.GetListItems(ctx, listName)
	if err != nil {
		return fmt.Errorf("failed to get list items")
	}
	listMessage := formListMessage(listItems)
	_, err = ctx.EffectiveMessage.Chat.SendMessage(b, listMessage, &gotgbot.SendMessageOpts{
		ReplyMarkup: getMainMenueKeyboard(),
	})
	if err != nil {
		return fmt.Errorf("failed to send finish message")
	}
	return handlers.EndConversation()
}

func (handler *ShopHandler) cancel(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Chat.SendMessage(b, "Жаль, что вы прервались!", &gotgbot.SendMessageOpts{
		ParseMode:   "html",
		ReplyMarkup: getMainMenueKeyboard(),
	})
	if err != nil {
		return fmt.Errorf("failed to send cancel message: %w", err)
	}
	return handlers.EndConversation()
}
