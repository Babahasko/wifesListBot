package bot

import (
	"errors"
	"fmt"
	"shopping_bot/internal/service"
	"shopping_bot/pkg/logger"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/conversation"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

type ShopHandler struct {
	ListCallbackService *ListCallbackService
	ItemCallbackService *ItemCallbackService
	Service             *service.ShoppingService
}

const (
	NAME          = "list_name"
	LIST_NAME     = "add_list_name"
	PURCHASES     = "add_purchases"
	FINISH        = "finish_form_list"
	FINISH_ADDING = "finish_add_items"
	ADDING_ITEMS  = "adding_items"
)

func NewShopHandler(router *ext.Dispatcher, service *service.ShoppingService) {
	// Handlers for shop
	handler := &ShopHandler{
		ListCallbackService: NewListCallbackService(),
		ItemCallbackService: NewItemCallbackService(),
		Service:             service,
	}

	router.AddHandler(handlers.NewCommand("start", handler.start))

	// Form list conversation
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

	// Add List conversation
	router.AddHandler(handlers.NewConversation(
		[]ext.Handler{
			handlers.NewCommand("add_list", handler.addList),
			handlers.NewCallback(callbackquery.Prefix(CallbackAddList), handler.addList),
		},
		map[string][]ext.Handler{
			LIST_NAME: {handlers.NewMessage(message.Text, handler.addListName)},
		},
		&handlers.ConversationOpts{
			Exits:        []ext.Handler{handlers.NewMessage(message.Equal(ButtonCancel), handler.cancel)},
			StateStorage: conversation.NewInMemoryStorage(conversation.KeyStrategySenderAndChat),
			AllowReEntry: true,
		},
	))

	// Add Items conversation
	router.AddHandler(handlers.NewConversation(
		[]ext.Handler{
			handlers.NewCallback(callbackquery.Prefix(CallbackAddItems), handler.startAddItems),
		},
		map[string][]ext.Handler{
			ADDING_ITEMS:  {handlers.NewMessage(message.Text, handler.addItem)},
			FINISH_ADDING: {handlers.NewMessage(message.Equal(ButtonFinishList), handler.finishAddItem)},
		},
		&handlers.ConversationOpts{
			Exits:        []ext.Handler{handlers.NewMessage(message.Equal(ButtonCancel), handler.cancel)},
			StateStorage: conversation.NewInMemoryStorage(conversation.KeyStrategySenderAndChat),
			AllowReEntry: true,
		},
	))

	// Lists Handlers
	router.AddHandler(handlers.NewMessage(message.Equal(ButtonViewList), handler.showLists))
	router.AddHandler(handlers.NewCallback(callbackquery.Prefix(handler.ListCallbackService.Prefix), handler.showListItems))
	router.AddHandler(handlers.NewCallback(callbackquery.Prefix(handler.ItemCallbackService.Prefix), handler.markItem))
	router.AddHandler(handlers.NewCallback(callbackquery.Prefix(CallbackBackToList), handler.backToLists))
	router.AddHandler(handlers.NewCallback(callbackquery.Prefix(CallbackClearList), handler.clearList))
	router.AddHandler(handlers.NewCallback(callbackquery.Prefix(CallbackDeleteList), handler.deleteList))
}

func (handler *ShopHandler) start(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Chat.SendMessage(b, MsgStart, &gotgbot.SendMessageOpts{
		ReplyMarkup: getMainMenueKeyboard(),
	})
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
	userID := ctx.EffectiveUser.Id

	handler.Service.SetCurrentList(userID, listName)
	handler.Service.AddShoppingList(userID, listName)
	// handler.Client.addShoppingList(ctx, listName)    // это у нас в базу летит shopping list
	// handler.Client.setCurrentListName(ctx, listName) // это у нас в кэш летит состояние пользователя
	currentlist, err := handler.Service.GetCurrentList(userID)
	if err != nil {
		return fmt.Errorf("failed to get current list namae for user")
	}
	logger.Sugar.Debugw("current user list", "current_list", currentlist)

	_, err = ctx.EffectiveMessage.Chat.SendMessage(b, MsgWriteItemName, &gotgbot.SendMessageOpts{
		ReplyMarkup: getFormListKeyboard(),
	})
	if err != nil {
		return fmt.Errorf("failed to send addPurchase message")
	}
	return handlers.NextConversationState(PURCHASES)
}


// clear user list when finished
func (handler *ShopHandler) addPurchase(b *gotgbot.Bot, ctx *ext.Context) error {
	itemName := ctx.EffectiveMessage.Text
	// Получаем имя текущего списка из кэша
	userID := ctx.EffectiveUser.Id
	listName, err := handler.Service.GetCurrentList(userID)
	if err != nil {
		return fmt.Errorf("failed to get current list namae for user")
	}
	
	if itemName == ButtonFinishList {
		return handler.finish(b, ctx)
	}
	if itemName == fmt.Sprintf("/%s", CommandEnd) {
		return handler.finish(b, ctx)
	}
	handler.Service.AddItemToShoppingList(userID, listName, itemName)
	// handler.Client.addItemToShoppingList(ctx, listName, itemName)
	logger.Sugar.Debugw("add item: to shopping list","itemName", itemName, "listName", listName)

	return handlers.NextConversationState(PURCHASES)
}

func (handler *ShopHandler) finish(b *gotgbot.Bot, ctx *ext.Context) error {
	// listName := handler.Client.getCurrentList(ctx)
	userID := ctx.EffectiveUser.Id
	listName, err := handler.Service.GetCurrentList(userID)
	if err != nil {
		return fmt.Errorf("failed to get current list namae for user")
	}

	// listItems, err := handler.Client.getListItems(ctx, listName)
	listItems, err := handler.Service.GetListItems(userID, listName)
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
	userID := ctx.EffectiveUser.Id

	// listNames, err := handler.Client.getUserLists(ctx)
	listNames, err := handler.Service.GetUserLists(userID)

	// Сначала проверяем, есть ли вообще ошибка
	switch {
	case errors.Is(err, ErrNoLists):
		_, sendErr := ctx.EffectiveMessage.Reply(b,
			MsgNoLists,
			nil)
		if sendErr != nil {
			return fmt.Errorf("failed to send no lists message: %w", sendErr)
		}
		return nil
	case err != nil:
		return fmt.Errorf("failed to get user lists: %w", err)
	}
	
	// Если ошибок нет, обрабатываем список
	logger.Sugar.Debugw("user lists", "listNames", listNames)

	listsKeyboard, err := getListsKeyboard(listNames, handler.ListCallbackService)
	if err != nil {
		return fmt.Errorf("failed to get lists keyboard")
	}

	_, err = ctx.EffectiveMessage.Chat.SendMessage(b, MsgYourLists, &gotgbot.SendMessageOpts{
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
	userID := ctx.EffectiveUser.Id

	// Распаковываем callback с помощью реестра
	listCallback := handler.ListCallbackService.Unpack(cbQuery.Data)

	// Теперь можно получить название списка
	listName := listCallback.Name

	// Устанавливаем текущий лист с которым работает пользователь
	// handler.Client.setCurrentListName(ctx, listName)
	handler.Service.SetCurrentList(userID, listName)

	logger.Sugar.Debugw("processing list", "name", listName)

	// Получаем элементы списка
	// listItems, err := handler.Client.getListItems(ctx, listName)
	listItems, err := handler.Service.GetListItems(userID, listName)

	if err != nil {
		return fmt.Errorf("failed to get list items: %w", err)
	}

	// Формируем клавиатуру со списком покупок
	itemsKeyboard, err := getItemsKeyboard(listItems, handler.ItemCallbackService)

	if err != nil {
		return fmt.Errorf("failed to get items keyboard:%w", err)
	}

	// Редактируем инлайн сообщение и отдаём клавиатуру со списком покупок
	_, _, err = cbQuery.Message.EditText(b, listName, &gotgbot.EditMessageTextOpts{
		ReplyMarkup: *itemsKeyboard,
	})

	if err != nil {
		return fmt.Errorf("failed to send items keyboard")
	}

	return nil
}

func (handler *ShopHandler) markItem(b *gotgbot.Bot, ctx *ext.Context) error {
	cbQuery := ctx.Update.CallbackQuery
	userID := ctx.EffectiveUser.Id

	// Распаковываем callback с помощью реестра
	itemCallback := handler.ItemCallbackService.Unpack(cbQuery.Data)
	listName := itemCallback.ListName
	itemName := itemCallback.ItemName

	// handler.Client.markItem(ctx, listName, itemName)
	handler.Service.MarkItem(userID, listName, itemName)
	cbQuery.Answer(b, &gotgbot.AnswerCallbackQueryOpts{
		Text: fmt.Sprintf("%s:%s", listName, itemName),
	})

	// listItems, err := handler.Client.getListItems(ctx, listName)
	listItems, err := handler.Service.GetListItems(userID, listName)
	if err != nil {
		return fmt.Errorf("failed to get list items: %w", err)
	}

	// Формируем новую клавиатуру
	itemsKeyboard, err := getItemsKeyboard(listItems, handler.ItemCallbackService)
	if err != nil {
		return fmt.Errorf("failed to get items keyboard: %w", err)
	}

	// Отправляем новое сообщение с клавиатурой покупок
	_, _, err = cbQuery.Message.EditReplyMarkup(b, &gotgbot.EditMessageReplyMarkupOpts{
		ReplyMarkup: *itemsKeyboard,
	})
	if err != nil {
		return fmt.Errorf("failed to send items keyboard")
	}
	return nil
}

func (handler *ShopHandler) cancel(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Chat.SendMessage(b, MsgPitty, &gotgbot.SendMessageOpts{
		ReplyMarkup: getMainMenueKeyboard(),
	})
	if err != nil {
		return fmt.Errorf("failed to send cancel message: %w", err)
	}
	return handlers.EndConversation()
}

func (handler *ShopHandler) backToLists(b *gotgbot.Bot, ctx *ext.Context) error {
	cbQuery := ctx.Update.CallbackQuery
	userID := ctx.EffectiveUser.Id
	// userLists, err := handler.getUserLists(b, ctx)
	userLists, err := handler.Service.GetUserLists(userID)
	if err != nil {
		return fmt.Errorf("failed to get user lists:%w", err)
	}
	listsKeyboard, err := getListsKeyboard(userLists, handler.ListCallbackService)
	if err != nil {
		return fmt.Errorf("failed to form lists keyboard: %w", err)
	}
	_, _, err = cbQuery.Message.EditText(b, "Ваши списки покупок", &gotgbot.EditMessageTextOpts{
		ReplyMarkup: *listsKeyboard,
	})

	if err != nil {
		return fmt.Errorf("failed to send lists keyboard: %w", err)
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
	userID := ctx.EffectiveUser.Id

	// err := handler.Client.addShoppingList(ctx, listName)
	err := handler.Service.AddShoppingList(userID, listName)
	if err != nil {
		return fmt.Errorf("failed to add shopping list:%w", err)
	}
	// listNames, err := handler.Client.getUserLists(ctx)
	listNames, err := handler.Service.GetUserLists(userID)
	if err != nil {
		return fmt.Errorf("failed to get user lists: %w", err)
	}
	listsKeyboard, err := getListsKeyboard(listNames, handler.ListCallbackService)
	if err != nil {
		return fmt.Errorf("failed to get lists keyboard")
	}
	_, err = ctx.EffectiveMessage.Chat.SendMessage(b, "Списки покупок:", &gotgbot.SendMessageOpts{
		ReplyMarkup: listsKeyboard,
	})
	if err != nil {
		return fmt.Errorf("failed to send list message")
	}
	return handlers.EndConversation()
}

func (handler *ShopHandler) clearList(b *gotgbot.Bot, ctx *ext.Context) error {
	cbQuery := ctx.Update.CallbackQuery
	userID := ctx.EffectiveUser.Id
	// Получаем название текущего списка пользователя
	// currentList := handler.Client.getCurrentList(ctx)
	currentList, err := handler.Service.GetCurrentList(userID)
	if err != nil {
		return fmt.Errorf("failed to get current list:%w", err)
	}
	// Удаляем все отмеченные покупки в листе(со статус checked)
	// handler.Client.deleteMarkItems(ctx, currentList)
	handler.Service.DeleteMarkedItems(userID, currentList)
	// Формируем новую клавиатуру
	// Получаем элементы списка
	// listItems, err := handler.Client.getListItems(ctx, currentList)
	listItems, err := handler.Service.GetListItems(userID, currentList)
	if err != nil {
		return fmt.Errorf("failed to get list items: %w", err)
	}
	// Формируем клавиатуру
	itemsKeyboard, err := getItemsKeyboard(listItems, handler.ItemCallbackService)

	if err != nil {
		return fmt.Errorf("failed to get items keyboard: %w", err)
	}
	// Редактируем InlineKeyboard
	_, _, err = cbQuery.Message.EditText(b, currentList, &gotgbot.EditMessageTextOpts{
		ReplyMarkup: *itemsKeyboard,
	})

	if err != nil {
		return fmt.Errorf("failed to send items keyboard")
	}
	return nil
}

func (handler *ShopHandler) deleteList(b *gotgbot.Bot, ctx *ext.Context) error {
	cbQuery := ctx.Update.CallbackQuery
	userID := ctx.EffectiveUser.Id
	// current_list := handler.Client.getCurrentList(ctx)
	current_list, err := handler.Service.GetCurrentList(userID)
	if err != nil {
		return fmt.Errorf("failed to get current list:%w", err)
	}
	// handler.Client.deleteList(ctx, current_list)
	handler.Service.DeleteList(userID, current_list)

	// userLists, err := handler.Client.getUserLists(ctx)
	userLists, err := handler.Service.GetUserLists(userID)
	if err != nil {
		return fmt.Errorf("failed to get user lists: %w", err)
	}
	listsKeyboard, err := getListsKeyboard(userLists, handler.ListCallbackService)
	if err != nil {
		return fmt.Errorf("failed to get user lists: %w", err)
	}
	_, _, err = cbQuery.Message.EditText(b, "Ваши списки покупок", &gotgbot.EditMessageTextOpts{
		ReplyMarkup: *listsKeyboard,
	})

	if err != nil {
		return fmt.Errorf("failed to edit callback query message")
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
	userID := ctx.EffectiveUser.Id
	// currentList := handler.Client.getCurrentList(ctx)
	currentList, err := handler.Service.GetCurrentList(userID) 
	if err != nil {
		return fmt.Errorf("failed to get current list name")
	}

	logger.Sugar.Debugw("get current list from state:", "current_list", currentList)
	if itemName == ButtonFinishList {
		return handler.finishAddItem(b, ctx)
	}
	if itemName == fmt.Sprintf("/%s", CommandEnd) {
		return handler.finishAddItem(b, ctx)
	}
	// handler.Client.addItemToShoppingList(ctx, currentList, itemName)
	handler.Service.AddItemToShoppingList(userID, currentList, itemName)
	logger.Sugar.Debugw("add item to shopping list", "itemName", itemName, "listName", currentList)
	return handlers.NextConversationState(ADDING_ITEMS)
}

func (handler *ShopHandler) finishAddItem(b *gotgbot.Bot, ctx *ext.Context) error {
	logger.Sugar.Debugw("finishAddItem handler")
	// listName := handler.Client.getCurrentList(ctx)
	userID:= ctx.EffectiveUser.Id
	currentList, err := handler.Service.GetCurrentList(userID)
	if err != nil {
		return fmt.Errorf("failed to get current list name: %w", err)
	}
	// listItems, err := handler.Client.getListItems(ctx, listName)
	listItems, err := handler.Service.GetListItems(userID, currentList)
	if err != nil {
		return fmt.Errorf("failed to get list items:%w", err)
	}
	listMessage := formListMessage(currentList, listItems)
	_, err = ctx.EffectiveMessage.Chat.SendMessage(b, listMessage, &gotgbot.SendMessageOpts{
		ReplyMarkup: getMainMenueKeyboard(),
	})
	if err != nil {
		return fmt.Errorf("failed to send finish message")
	}
	return handlers.EndConversation()
}