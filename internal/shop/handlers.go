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
	Client *ShopClient
}

const (
	NAME = "list_name"
	PURCHASES = "add_purchases"
	FINISH = "finish_form_list"
)

func NewShopHandler(router *ext.Dispatcher) {
	// Handlers for shop
	handler := &ShopHandler{
		//Add callback registry
		CallbackRegistry: callback.NewRegistry(),
		Client: &ShopClient{},
	}
	// Register Callbacks
	handler.CallbackRegistry.Register(NewCategoryCallback)

	router.AddHandler(handlers.NewCommand("start", handler.start))
	// router.AddHandler(handlers.NewMessage(message.Equal(ButtonAddPurchase), handler.AddPurchase))
	router.AddHandler(handlers.NewConversation(
		[]ext.Handler{handlers.NewMessage(message.Equal(ButtonAddPurchase), handler.formList)},
		map[string][]ext.Handler{
			NAME: {handlers.NewMessage(message.Text, handler.addName)},
			PURCHASES: {handlers.NewMessage(message.Text, handler.addPurchase)},
			FINISH: {handlers.NewMessage(message.Equal(ButtonFinishList), handler.finish)},
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
		ReplyMarkup: getMenueKeyboard(),
	})
	if err != nil {
		return fmt.Errorf("failed to send menue message: %w", err)
	}
	return handlers.NextConversationState(PURCHASES)
}
func (handler *ShopHandler) addName(b *gotgbot.Bot, ctx *ext.Context) error {
	listName := ctx.EffectiveMessage.Text
	handler.Client.addShoppingList(ctx, listName)
	handler.Client.setUserData(ctx, "current_list", listName)
	_, err := ctx.EffectiveMessage.Chat.SendMessage(b, "Введите название покупки", &gotgbot.SendMessageOpts{
		ReplyMarkup: getMenueKeyboard(),
	})
	if err != nil {
		return fmt.Errorf("failed to send addPurchase message")
	}
	return handlers.NextConversationState(PURCHASES)
}


//TODO: Add validation for item name`s`
// clear user list when finished
func (handler *ShopHandler) addPurchase(b *gotgbot.Bot, ctx *ext.Context) error {
	itemName := ctx.EffectiveMessage.Text
	// Получаем имя списка и проверяем тип
    listNameIface, ok := handler.Client.getUserData(ctx, "current_list")
    if !ok || listNameIface == nil {
        logger.Sugar.Errorw("Current list name not found in user data", "userID", ctx.EffectiveUser.Id)
        _, err := ctx.EffectiveMessage.Chat.SendMessage(b, "Ошибка: список покупок не задан. Начните заново.", nil)
        return err
    }

    listName, ok := listNameIface.(string)
    if !ok {
        logger.Sugar.Errorw("Expected current_list to be a string", "value", listNameIface, "type", fmt.Sprintf("%T", listNameIface))
        _, err := ctx.EffectiveMessage.Chat.SendMessage(b, "Ошибка: некорректный формат списка. Начните заново.", nil)
        return err
    }

	handler.Client.AddItemToShoppingList(ctx, listName, itemName)
	return handlers.NextConversationState(PURCHASES)
}

func (handler *ShopHandler) finish(b *gotgbot.Bot, ctx *ext.Context) error {
	listNameIface, ok := handler.Client.getUserData(ctx, "current_list")
    if !ok || listNameIface == nil {
        logger.Sugar.Errorw("Current list name not found in user data", "userID", ctx.EffectiveUser.Id)
        _, err := ctx.EffectiveMessage.Chat.SendMessage(b, "Ошибка: список покупок не задан. Начните заново.", nil)
        return err
    }

    listName, ok := listNameIface.(string)
    if !ok {
        logger.Sugar.Errorw("Expected current_list to be a string", "value", listNameIface, "type", fmt.Sprintf("%T", listNameIface))
        _, err := ctx.EffectiveMessage.Chat.SendMessage(b, "Ошибка: некорректный формат списка. Начните заново.", nil)
        return err
    }
	listItems, err := handler.Client.GetShoppingList(ctx, listName)
	if err != nil {
		return fmt.Errorf("failed to get list items")
	}
	_, err = ctx.EffectiveMessage.Chat.SendMessage(b, fmt.Sprintf("Отлично вот ваш список покупоЖ %v", listItems), nil)
	if err != nil {
		return fmt.Errorf("failed to send finish message")
	}
	return handlers.EndConversation()
}


func (handler *ShopHandler) cancel(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Chat.SendMessage(b, "Жаль, что вы прервались!", &gotgbot.SendMessageOpts{
		ParseMode: "html",
		ReplyMarkup: getMainMenueKeyboard(),
	})
	if err != nil {
		return fmt.Errorf("failed to send cancel message: %w", err)
	}
	return handlers.EndConversation()
}
