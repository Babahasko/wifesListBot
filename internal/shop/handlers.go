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
	handler.CallbackRegistry.Register(NewListCallback)
	handler.CallbackRegistry.Register(NewItemsCallback)

	router.AddHandler(handlers.NewCommand("start", handler.start))

	// Form list conversation handlers
	router.AddHandler(handlers.NewConversation(
		[]ext.Handler{
			handlers.NewMessage(message.Equal(ButtonAddPurchase), handler.formList),
			handlers.NewCommand("add", handler.formList),
		},
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

	// Lists Handlers
	router.AddHandler(handlers.NewMessage(message.Equal(ButtonViewList), handler.showLists))
	router.AddHandler(handlers.NewCallback(callbackquery.Prefix("list"), handler.showListItems))
	router.AddHandler(handlers.NewCallback(callbackquery.Prefix("item"), handler.markItem))
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
	handler.Client.setCurrentListName(ctx, listName) // это у нас в кэш летит состояние пользователя
	logger.Sugar.Debugw("current user list: %v", handler.Client.getCurrentList(ctx))
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
	listName := handler.Client.getCurrentList(ctx)
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
	handler.Client.addItemToShoppingList(ctx, listName, itemName)
	logger.Sugar.Debugw("add item: %v to shopping list: %v", itemName, listName)

	return handlers.NextConversationState(PURCHASES)
}

func (handler *ShopHandler) finish(b *gotgbot.Bot, ctx *ext.Context) error {
	listName := handler.Client.getCurrentList(ctx)
	if listName == "" {
		logger.Sugar.Errorw("Current list name not found in user states", "userID", ctx.EffectiveUser.Id)
		_, err := ctx.EffectiveMessage.Chat.SendMessage(b, "Ошибка: список покупок не задан. Начните заново", nil)
		return fmt.Errorf("failed to send message: %w", err)
	}
	listItems, err := handler.Client.getListItems(ctx, listName)
	if err != nil {
		return fmt.Errorf("failed to get list items")
	}
	listMessage := formListMessage(listName, listItems)
	_, err = ctx.EffectiveMessage.Chat.SendMessage(b, listMessage, &gotgbot.SendMessageOpts{
		ReplyMarkup: getMainMenueKeyboard(),
	})
	if err != nil {
		return fmt.Errorf("failed to send finish message")
	}
	return handlers.EndConversation()
}

func (handler *ShopHandler) showLists(b *gotgbot.Bot, ctx *ext.Context) error {
	listNames, err := handler.Client.getUserLists(ctx)

	// Сначала проверяем, есть ли вообще ошибка
	if err != nil {
		// Проверяем, это наша специальная ошибка "нет списков"
		if err.Error() == ErrorNoLists {
			_, sendErr := ctx.EffectiveMessage.Reply(b,
				"У вас ещё нет ни одного списка! Создайте для начала командой /add",
				nil)
			if sendErr != nil {
				return fmt.Errorf("failed to send no lists message: %w", sendErr)
			}
			return nil
		}
		// Все другие ошибки
		return fmt.Errorf("failed to get user lists: %w", err)
	}

	// Если ошибок нет, обрабатываем список
	logger.Sugar.Debugw("user lists", "listNames", listNames)

	// Здесь должна быть логика формирования и отправки клавиатуры
	// Например:
	if len(listNames) == 0 {
		_, err := ctx.EffectiveMessage.Reply(b,
			"У вас ещё нет ни одного списка! Создайте для начала командой /add",
			nil)
		if err != nil {
			return fmt.Errorf("failed to send empty lists message: %w", err)
		}
		return nil
	}

	listsKeyboard, err := getListsKeyboard(listNames)
	if err != nil {
		return fmt.Errorf("failed to get lists keyboard")
	}

	_, err = ctx.EffectiveMessage.Chat.SendMessage(b, "Ваши списки", &gotgbot.SendMessageOpts{
		ReplyMarkup: listsKeyboard,
	})
	if err != nil {
		return fmt.Errorf("failed to send keyboard message")
	}

	// TODO: Добавьте здесь код для создания и отправки inline-клавиатуры
	// с использованием listNames

	return nil
}

func (handler *ShopHandler) showListItems(b *gotgbot.Bot, ctx *ext.Context) error {
	// Получаем данные callback
	cbQuery := ctx.Update.CallbackQuery
	if cbQuery == nil || cbQuery.Data == "" {
		return fmt.Errorf("empty callback data")
	}

	// Распаковываем callback с помощью реестра
	listCallback, err := callback.ParseCallback[*ListCallback](handler.CallbackRegistry, cbQuery.Data)
	if err != nil {
		logger.Sugar.Errorw("failed to parse callback data", "error", err, "data", cbQuery.Data)
		return fmt.Errorf("failed to parse callback data: %w", err)
	}

	// Теперь можно получить название списка
	listName := listCallback.Name
	logger.Sugar.Debugw("processing list", "name", listName)

	// Получаем элементы списка
	listItems, err := handler.Client.getListItems(ctx, listName)
	if err != nil {
		logger.Sugar.Errorw("failed to get list items", "list", listName, "error", err)
		_, sendErr := cbQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Text: "Не удалось получить список покупок",
		})
		if sendErr != nil {
			return fmt.Errorf("failed to send error message: %w", sendErr)
		}
		return fmt.Errorf("failed to get list items: %w", err)
	}

	// Отправляем сообщение с выбраным списком
	_, err = cbQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text: listName,
	})
	if err != nil {
		return fmt.Errorf("failed to answer callback message: %w", err)
	}

	// Удаляем сообщение с клавиатурой листа покупок
	_, err = cbQuery.Message.Delete(b, nil)

	if err != nil {
		return fmt.Errorf("failed to send add name message %w", err)
	}

	itemsKeyboard, err := getItemsKeyboard(listItems)
	if err != nil {
		return fmt.Errorf("failed to get items keyboard")
	}
	// Отправляем новое сообщение с клавиатурой покупок
	_, err = ctx.EffectiveMessage.Chat.SendMessage(b, listName, &gotgbot.SendMessageOpts{
		ReplyMarkup: itemsKeyboard,
		ParseMode: "MarkdownV2",
	})
	if err != nil {
		return fmt.Errorf("failed to send items keyboard")
	}

	return nil
}

func (handler *ShopHandler) markItem(b *gotgbot.Bot, ctx *ext.Context) error {
	cbQuery := ctx.Update.CallbackQuery
	if cbQuery == nil || cbQuery.Data == "" {
		return fmt.Errorf("empty callback data")
	}

	// Распаковываем callback с помощью реестра
	itemCallback, err := callback.ParseCallback[*ItemCallback](handler.CallbackRegistry, cbQuery.Data)
	if err != nil {
		logger.Sugar.Errorw("failed to parse callback data", "error", err, "data", cbQuery.Data)
		return fmt.Errorf("failed to parse callback data: %w", err)
	}
	listName := itemCallback.ListName
	itemName := itemCallback.ItemName

	handler.Client.markItem(ctx, listName, itemName)
	cbQuery.Answer(b,&gotgbot.AnswerCallbackQueryOpts{
		Text: fmt.Sprintf("%s:%s", listName, itemName),
	})

	// Удаляем сообщение с клавиатурой листа покупок
	_, err = cbQuery.Message.Delete(b, nil)

	if err != nil {
		return fmt.Errorf("failed to send add name message %w", err)
	}

	listItems, err := handler.Client.getListItems(ctx, listName)
	if err != nil {
		logger.Sugar.Errorw("failed to get list items", "list", listName, "error", err)
		_, sendErr := cbQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Text: "Не удалось получить список покупок",
		})
		if sendErr != nil {
			return fmt.Errorf("failed to send error message: %w", sendErr)
		}
		return fmt.Errorf("failed to get list items: %w", err)
	}

	itemsKeyboard, err := getItemsKeyboard(listItems)
	if err != nil {
		return fmt.Errorf("failed to get items keyboard: %w", err)
	}
	// Отправляем новое сообщение с клавиатурой покупок
	_, err = ctx.EffectiveMessage.Chat.SendMessage(b, listName, &gotgbot.SendMessageOpts{
		ReplyMarkup: itemsKeyboard,
	})
	if err != nil {
		return fmt.Errorf("failed to send items keyboard")
	}

	

	return nil
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
