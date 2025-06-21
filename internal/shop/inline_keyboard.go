package shop

import (
	"fmt"

	"shopping_bot/pkg/callback"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// TODO: Здесь будет клавиатура формирования списка покупок

func getListsKeyboard(lists []string) (*gotgbot.InlineKeyboardMarkup, error) {
	var ButtonsPerRow = 1
	var rows [][]gotgbot.InlineKeyboardButton
	var buttons []gotgbot.InlineKeyboardButton
	for _, list := range lists {
		//Callback
		callbackStr, err := callback.PackCallback(&ListCallback{
			Name: list,
		})
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

	return &gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}, nil
}

// TODO: обновить, чтобы возвращала указатель на клавиатуру
func getItemsKeyboard(items []*ShoppingItem) (gotgbot.InlineKeyboardMarkup, error) {
	if len(items) == 0 {
		return gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{
					{Text: ButtonEmptyList, CallbackData: "no_items"},
				},
				{
					{Text: ButtonBack, CallbackData: CallbackBackToList},
				},
			},
		}, nil
	}

	// Создаем кнопки для каждой покупки
	var keyboard [][]gotgbot.InlineKeyboardButton
	for _, item := range items {
		// Здесь можно добавить callback данные для каждой покупки
		// Например, чтобы можно было отметить как купленное
		callbackStr, err := callback.PackCallback(&ItemCallback{
			ListName: item.ListName,
			ItemName: item.Name,
		})
		if err != nil {
			return gotgbot.InlineKeyboardMarkup{}, fmt.Errorf("failed to pack callback")
		}

		text := item.Name
		if item.Checked{
			text = fmt.Sprintf("✅%s", item.Name)
		}

		row := []gotgbot.InlineKeyboardButton{
			{
				Text:        text,
				CallbackData: callbackStr,
			},
		}
		keyboard = append(keyboard, row)
	}

	// Добавляем кнопку "Назад"
	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{Text: ButtonBackToLists, CallbackData: CallbackBackToList},
		{Text: ButtonClearList, CallbackData: CallbackClearList},
	})
	keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{
		{Text: ButtonDeleteList, CallbackData: "delete_list"},
	})

	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}, nil
}
