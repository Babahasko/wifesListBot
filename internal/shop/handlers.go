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
	LIST_NAME = "add_list_name"
	PURCHASES = "add_purchases"
	FINISH    = "finish_form_list"
	FINISH_ADDING = "finish_add_items"
	ADDING_ITEMS = "adding_items"
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

	// AddList conversation
	router.AddHandler(handlers.NewConversation(
		[]ext.Handler{
			handlers.NewCommand("add_list", handler.addList),
			handlers.NewCallback(callbackquery.Prefix(CallbackAddList), handler.addList),
		},
		map[string][]ext.Handler{
			LIST_NAME:      {handlers.NewMessage(message.Text, handler.addListName)},
		},
		&handlers.ConversationOpts{
			Exits:        []ext.Handler{handlers.NewMessage(message.Equal(ButtonCancel), handler.cancel)},
			StateStorage: conversation.NewInMemoryStorage(conversation.KeyStrategySenderAndChat),
			AllowReEntry: true,
		},
	))

	// Add items conversation
	router.AddHandler(handlers.NewConversation(
		[]ext.Handler{
			handlers.NewCallback(callbackquery.Prefix(CallbackNoItems), handler.startAddItems),
		},
		map[string][]ext.Handler{
			ADDING_ITEMS: {handlers.NewMessage(message.Text, handler.addItem)},
			FINISH_ADDING: {handlers.NewMessage(message.Equal(ButtonFinishList), handler.finishAddItem)},
		},
		&handlers.ConversationOpts{
			Exits: []ext.Handler{handlers.NewMessage(message.Equal(ButtonCancel), handler.cancel)},
			StateStorage: conversation.NewInMemoryStorage(conversation.KeyStrategySenderAndChat),
			AllowReEntry: true,
		},
	))

	// Lists Handlers
	router.AddHandler(handlers.NewMessage(message.Equal(ButtonViewList), handler.showLists))
	router.AddHandler(handlers.NewCallback(callbackquery.Prefix("list"), handler.showListItems))
	router.AddHandler(handlers.NewCallback(callbackquery.Prefix("item"), handler.markItem))
	router.AddHandler(handlers.NewCallback(callbackquery.Prefix(CallbackBackToList), handler.backToLists))
	router.AddHandler(handlers.NewCallback(callbackquery.Prefix(CallbackClearList), handler.clearList))
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
	_, err := ctx.EffectiveMessage.Chat.SendMessage(b, MsgWriteListName, &gotgbot.SendMessageOpts{
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
	_, err := ctx.EffectiveMessage.Chat.SendMessage(b, MsgWriteItemName, &gotgbot.SendMessageOpts{
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

	if itemName == ButtonFinishList {
		return handler.finish(b, ctx)
	}
	if itemName == "/end" {
		return handler.finish(b, ctx)
	}
	handler.Client.addItemToShoppingList(ctx, listName, itemName)
	logger.Sugar.Debugw("add item: %v to shopping list: %v", itemName, listName)

	return handlers.NextConversationState(PURCHASES)
}

func (handler *ShopHandler) finish(b *gotgbot.Bot, ctx *ext.Context) error {
	listName := handler.Client.getCurrentList(ctx)

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
				"У вас ещё нет ни одного списка! Создайте для начала командой /add_list",
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
	return nil
}

func (handler *ShopHandler) showListItems(b *gotgbot.Bot, ctx *ext.Context) error {
	// Получаем данные callback
	cbQuery := ctx.Update.CallbackQuery

	// Распаковываем callback с помощью реестра
	listCallback, err := callback.ParseCallback[*ListCallback](handler.CallbackRegistry, cbQuery.Data)
	if err != nil {
		logger.Sugar.Errorw("failed to parse callback data", "error", err, "data", cbQuery.Data)
		return fmt.Errorf("failed to parse callback data: %w", err)
	}

	// Теперь можно получить название списка
	listName := listCallback.Name
	// Устанавливаем текущий лист с которым работает пользователь
	handler.Client.setCurrentListName(ctx, listName)
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

	// Формируем клавиатуру со списком покупок
	itemsKeyboard, err := getItemsKeyboard(listItems)

	if err != nil {
		return fmt.Errorf("failed to get items keyboard:%w", err)
	}

	// Редактируем инлайн сообщение и отдаём клавиатуру со списком покупок
	_,_, err = cbQuery.Message.EditText(b, listName, &gotgbot.EditMessageTextOpts{
		ReplyMarkup: itemsKeyboard,
	})

	if err != nil {
		return fmt.Errorf("failed to send items keyboard")
	}

	return nil
}

func (handler *ShopHandler) markItem(b *gotgbot.Bot, ctx *ext.Context) error {
	cbQuery := ctx.Update.CallbackQuery
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
	_,_,err = cbQuery.Message.EditReplyMarkup(b, &gotgbot.EditMessageReplyMarkupOpts{
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

func(handler *ShopHandler) backToLists(b *gotgbot.Bot, ctx *ext.Context) error {
	cbQuery := ctx.Update.CallbackQuery
	userLists, err := handler.getUserLists(b, ctx)
	if err != nil {
		return fmt.Errorf("failed to get user lists in service:%w", err)
	}
	_,_, err = cbQuery.Message.EditText(b, "Ваши списки покупок", &gotgbot.EditMessageTextOpts{
		ReplyMarkup: *userLists,
	})

	if err != nil {
		return fmt.Errorf("failed to back to lists: %w", err)
	}
	return nil
}

func (handler *ShopHandler) addList(b *gotgbot.Bot, ctx *ext.Context) error {
	cbquery := ctx.Update.CallbackQuery
	if cbquery != nil {
		cbquery.Answer(b, nil)
	}
	_, err := ctx.EffectiveMessage.Chat.SendMessage(b, MsgWriteListName, &gotgbot.SendMessageOpts{
		ReplyMarkup: getCancelKeyboard(),
	})
	if err != nil {
		return fmt.Errorf("failed to send menue message: %w", err)
	}
	return handlers.NextConversationState(LIST_NAME)
}

func (handler *ShopHandler) addListName(b *gotgbot.Bot, ctx *ext.Context) error {
	//TODO: Добавить валидацию shopping list
	listName := ctx.EffectiveMessage.Text
	err := handler.Client.addShoppingList(ctx, listName)
	if err != nil {
		return fmt.Errorf("failed to add shopping list:%w", err)
	}
	listNames, err := handler.Client.getUserLists(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user lists: %w", err)
	}

	listsKeyboard, err := getListsKeyboard(listNames)
	if err != nil {
		return fmt.Errorf("failed to get lists keyboard")
	}
	_, err = ctx.EffectiveMessage.Chat.SendMessage(b, "Отлично, вот ваши списки покупок!", &gotgbot.SendMessageOpts{
		ReplyMarkup: getMainMenueKeyboard(),
	})
	if err != nil {
		return fmt.Errorf("failed to send successfull message")
	}
	_, err = ctx.EffectiveMessage.Chat.SendMessage(b, "Списки покупок:", &gotgbot.SendMessageOpts{
		ReplyMarkup: listsKeyboard,
	})
	if err != nil {
		return fmt.Errorf("failed to send list message")
	}
	return handlers.EndConversation()
}


func(handler *ShopHandler) clearList(b *gotgbot.Bot, ctx *ext.Context) error {
	cbQuery := ctx.Update.CallbackQuery
	// Получаем название текущего списка пользователя
	currentList := handler.Client.getCurrentList(ctx)
	// Удаляем все отмеченные покупки в листе(со статус checked)
	handler.Client.deleteMarkItems(ctx, currentList)
	// Формируем новую клавиатуру
	//		Получаем элементы списка
	listItems, err := handler.Client.getListItems(ctx, currentList)
	if err != nil {
		logger.Sugar.Errorw("failed to get list items", "list", currentList, "error", err)
		_, sendErr := cbQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
			Text: "Не удалось получить список покупок",
		})
		if sendErr != nil {
			return fmt.Errorf("failed to send error message: %w", sendErr)
		}
		return fmt.Errorf("failed to get list items: %w", err)
	}
	//		Формируем клавиатуру
	itemsKeyboard, err := getItemsKeyboard(listItems)

	if err != nil {
		return fmt.Errorf("failed to get items keyboard: %w", err)
	}
	// Редактируем InlineKeyboard
	_,_, err = cbQuery.Message.EditText(b, currentList, &gotgbot.EditMessageTextOpts{
		ReplyMarkup: itemsKeyboard,
	})

	if err != nil {
		return fmt.Errorf("failed to send items keyboard")
	}
	return nil
}

func (handler *ShopHandler) startAddItems(b *gotgbot.Bot, ctx *ext.Context) error {
	cbQuery := ctx.Update.CallbackQuery
	_, err := cbQuery.Message.Delete(b, nil)
	if err != nil {
		return fmt.Errorf("failed to delete query message")
	}
	_, err = ctx.EffectiveMessage.Chat.SendMessage(b, MsgWriteItemName, &gotgbot.SendMessageOpts{
		ReplyMarkup: getFormListKeyboard(),
	})
	if err != nil {
		return fmt.Errorf("failed to send menue message: %w", err)
	}
	return handlers.NextConversationState(ADDING_ITEMS)
}

// TODO: Дописать функцию добавления item-ов
func (handler *ShopHandler) addItem(b *gotgbot.Bot, ctx *ext.Context) error {
	itemName := ctx.EffectiveMessage.Text
	currentList := handler.Client.getCurrentList(ctx)
	logger.Sugar.Debugw("get current list from state:", "current_list", currentList)
	if itemName == ButtonFinishList {
		return handler.finishAddItem(b, ctx)
	}
	if itemName == "/end" {
		logger.Sugar.Debugw("/end case")
		return handler.finishAddItem(b, ctx)
	}
	handler.Client.addItemToShoppingList(ctx, currentList, itemName)
	logger.Sugar.Debugw("add item to shopping list","itemName", itemName, "listName", currentList)

	return handlers.NextConversationState(ADDING_ITEMS)
}

func (handler *ShopHandler) finishAddItem(b *gotgbot.Bot, ctx *ext.Context) error {
	logger.Sugar.Debugw("finishAddItem handler")
	listName := handler.Client.getCurrentList(ctx)
	
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

//TODO: Вынести логику поулчения клавиатуры юзера в Service
func (handler *ShopHandler) getUserLists(b *gotgbot.Bot, ctx *ext.Context) (*gotgbot.InlineKeyboardMarkup, error) {
	listNames, err := handler.Client.getUserLists(ctx)

	// Сначала проверяем, есть ли вообще ошибка
	if err != nil {
		// Проверяем, это наша специальная ошибка "нет списков"
		if err.Error() == ErrorNoLists {
			_, sendErr := ctx.EffectiveMessage.Reply(b,
				"У вас ещё нет ни одного списка! Создайте для начала командой /add",
				nil)
			if sendErr != nil {
				return nil, fmt.Errorf("failed to send no lists message: %w", sendErr)
			}
			return nil, nil
		}
		// Все другие ошибки
		return nil, fmt.Errorf("failed to get user lists: %w", err)
	}

	// Если ошибок нет, обрабатываем список
	logger.Sugar.Debugw("user lists", "listNames", listNames)

	if len(listNames) == 0 {
		_, err := ctx.EffectiveMessage.Reply(b,
			"У вас ещё нет ни одного списка! Создайте для начала командой /add",
			nil)
		if err != nil {
			return nil, fmt.Errorf("failed to send empty lists message: %w", err)
		}
		return nil, nil
	}

	listsKeyboard, err := getListsKeyboard(listNames)
	if err != nil {
		return nil, fmt.Errorf("failed to get lists keyboard")
	}
	return listsKeyboard, err
}
