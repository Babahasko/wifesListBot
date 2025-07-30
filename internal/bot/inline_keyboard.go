package bot

import (
	"fmt"
	"shopping_bot/internal/models"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// TODO: Здесь будет клавиатура формирования списка покупок

func getListsKeyboard(lists []string, cbService *ListCallbackService) (*gotgbot.InlineKeyboardMarkup, error) {
	var ButtonsPerRow = 2
	var rows [][]gotgbot.InlineKeyboardButton
	var buttons []gotgbot.InlineKeyboardButton
	if len(lists) == 0 {
		return &gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{Text: ButtonAddList, CallbackData: CallbackAddList},
				},
			},
		}, nil
	}
	for _, list := range lists {
		//Callback
		callbackStr, err := cbService.Pack(list)
		if err != nil {
			return nil, fmt.Errorf("failed to pack callback: %v", err)
		}

		buttons = append(buttons, gotgbot.InlineKeyboardButton{
			Text:         list,
			CallbackData: callbackStr,
		})

		if len(buttons) == ButtonsPerRow {
			rows = append(rows, buttons)
			buttons = []gotgbot.InlineKeyboardButton{}
		}
	}

	if len(buttons) > 0 {
		rows = append(rows, buttons)
	}

	// Добавляем кнопку "Назад"
	rows = append(rows, []gotgbot.InlineKeyboardButton{
		{Text: ButtonAddList, CallbackData: CallbackAddList},
	})

	return &gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}, nil
}

// TODO: обновить, чтобы возвращала указатель на клавиатуру
func getItemsKeyboard(items []*models.ShoppingItem, cbService *ItemCallbackService) (*gotgbot.InlineKeyboardMarkup, error) {
	if len(items) == 0 {
		return &gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{Text: ButtonAddItems, CallbackData: CallbackAddItems},
				},
				{
					{Text: ButtonBack, CallbackData: CallbackBackToList},
				},
				{
					{Text: ButtonDeleteList, CallbackData: CallbackDeleteList},
				},
			},
		}, nil
	}

	// Создаем кнопки для каждой покупки
	var keyboard [][]gotgbot.InlineKeyboardButton
	for _, item := range items {
		// Здесь можно добавить callback данные для каждой покупки
		// Например, чтобы можно было отметить как купленное
		callbackStr, err := cbService.Pack(item.Name, item.ListName)
		if err != nil {
			return &gotgbot.InlineKeyboardMarkup{}, fmt.Errorf("failed to pack callback:%w", err)
		}

		text := item.Name
		if item.Checked {
			text = fmt.Sprintf("✅%s", item.Name)
		}

		row := []gotgbot.InlineKeyboardButton{
			{
				Text:         text,
				CallbackData: callbackStr,
			},
		}
		keyboard = append(keyboard, row)
	}

	// Добавляем кнопку "Назад"
	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{Text: ButtonAddItems, CallbackData: CallbackAddItems},
		{Text: ButtonClearList, CallbackData: CallbackClearList},
	})
	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{Text: ButtonBackToLists, CallbackData: CallbackBackToList},
		{Text: ButtonDeleteList, CallbackData: CallbackDeleteList},
	})

	return &gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}, nil
}
