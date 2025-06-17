package shop

import (
	"fmt"
	"shopping_bot/pkg/callback"
	"shopping_bot/pkg/logger"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/conversation"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

type ShopHandler struct {
	CallbackRegistry *callback.Registry
	Client *ShopClient
}

const (
	CATEGORY = "category"
	NAME = "purchase_name"
	PRIORITY = "priority"
	PRICE = "price"
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
		[]ext.Handler{handlers.NewMessage(message.Equal(ButtonAddPurchase), handler.addPurchase)},
		map[string][]ext.Handler{
			CATEGORY: {handlers.NewCallback(callbackquery.Prefix("cat_"), handler.category)},
			NAME: {handlers.NewMessage(message.Text, handler.addName)},
			PRIORITY: {handlers.NewMessage(message.Text, handler.addPriority)},
			PRICE: {handlers.NewMessage(message.Text, handler.addPrice)},
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

func (handler *ShopHandler) addPurchase(b *gotgbot.Bot, ctx *ext.Context) error {
	categories := getUserCategories(123)
	catKeyboard, err := getCategoriesKeyboard(categories)
	if err != nil {
		return fmt.Errorf("failed to get categories keyboard: %w", err)
	}

	_, err = ctx.EffectiveMessage.Chat.SendMessage(b, "Выберите категорию", &gotgbot.SendMessageOpts{
		ReplyMarkup: catKeyboard,
	})
	if err != nil {
		return fmt.Errorf("failed to send menue message: %w", err)
	}
	return handlers.NextConversationState(CATEGORY)
}



func (handler *ShopHandler) category(b *gotgbot.Bot, ctx *ext.Context) error {
	cbQuery := ctx.Update.CallbackQuery
	cbData := ctx.Update.CallbackQuery.Data

	// Unpack CallBack Data
	data, err := handler.CallbackRegistry.Parse(cbData)
	if err != nil {
		return fmt.Errorf("failed to parse callback: %w", err)
	}

	categoryData, ok := data.(*CategoryCallback)
	if !ok {
		return fmt.Errorf("unexpected callback type")
	}


	_, err = cbQuery.Answer(b,  &gotgbot.AnswerCallbackQueryOpts{
		Text: fmt.Sprintf("Выбрана категория %s", categoryData.Name),
		ShowAlert: false,
	})

	// set Category in UserData
	handler.Client.setUserData(ctx, "category", categoryData.Name)
	sessionData, result := handler.Client.getUserData(ctx, "category")
	logger.Sugar.Debugw("get user data", "data", sessionData, "result", result)

	if err != nil {
		return fmt.Errorf("failed to send category message %w", err)
	}

	_, err = cbQuery.Message.Delete(b,nil)

	if err != nil {
		return fmt.Errorf("failed to delete cb message %w", err)
	}

	_, err = ctx.EffectiveMessage.Chat.SendMessage(b, fmt.Sprintf("Выбрана категория: %s. Введите название покупки", categoryData.Name), &gotgbot.SendMessageOpts{
		ReplyMarkup: getMenueKeyboard(),
	})

	if err != nil {
		return fmt.Errorf("failed to send add name message %w", err)
	}
	
	// }
	return handlers.NextConversationState(NAME)
}

func (handler *ShopHandler) addName(b *gotgbot.Bot, ctx *ext.Context) error {
	purchaseName := ctx.EffectiveMessage.Text
	handler.Client.setUserData(ctx, "purchaseName", purchaseName)
	_, err := ctx.EffectiveMessage.Chat.SendMessage(b, "Введите приоритет от 1 до 10", &gotgbot.SendMessageOpts{
		ReplyMarkup: getMenueKeyboard(),
	})
	if err != nil {
		return fmt.Errorf("failed to send priority message: %w", err)
	}
	return handlers.NextConversationState(PRIORITY)
}

func (handler *ShopHandler) addPriority(b *gotgbot.Bot, ctx *ext.Context) error {
	purchasePriority := ctx.EffectiveMessage.Text
	handler.Client.setUserData(ctx, "purchasePriority", purchasePriority)
	_, err := ctx.EffectiveMessage.Chat.SendMessage(b, "Введите цену покупки", &gotgbot.SendMessageOpts{
		ReplyMarkup: getMenueKeyboard(),
	})
	if err != nil {
		return fmt.Errorf("failed to send price message: %w", err)
	}
	return handlers.NextConversationState(PRICE)
}

func (handler *ShopHandler) addPrice(b *gotgbot.Bot, ctx *ext.Context) error {
	purchasePrice := ctx.EffectiveMessage.Text
	handler.Client.setUserData(ctx, "purchasePrice", purchasePrice)
	name, _ := handler.Client.getUserData(ctx, "purchaseName")
	priority, _ := handler.Client.getUserData(ctx, "purchasePriority")
	price, _ := handler.Client.getUserData(ctx, "purchasePrice")
	_, err := ctx.EffectiveMessage.Chat.SendMessage(b, fmt.Sprintf("Итоговая покупка = название: %v, приоритет: %v, цена: %v.",name, priority, price), nil)
	if err != nil {
		return fmt.Errorf("failed to send price message: %w", err)
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
