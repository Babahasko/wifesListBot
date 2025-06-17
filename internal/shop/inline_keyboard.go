package shop

import (
	"log"

	"shopping_bot/pkg/callback"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// TODO: Заменить на нормальную бд

var userCategories = map[int64][]string{
	123: {"Электроника", "Одежда", "Книги"},
	456: {"Мебель", "Спорт", "Продукты"},
}

// Получение категорий пользователя
func getUserCategories(userID int64) []string {
	if cats, ok := userCategories[userID]; ok {
		return cats
	}
	return []string{"Общие товары"} // Категории по умолчанию
}

func getCategoriesKeyboard(categories []string) *gotgbot.InlineKeyboardMarkup {
	var ButtonsPerRow = 3
	var rows [][]gotgbot.InlineKeyboardButton
	var buttons []gotgbot.InlineKeyboardButton
	for _, category := range categories {

		//Callback 
		callbackStr, err := callback.PackCallback(&CategoryCallback{
			Name: category,
		})
		if err != nil {
			log.Printf("Failed to pack callback: %v", err)
		}
		
		buttons = append(buttons, gotgbot.InlineKeyboardButton{
			Text:         category,
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
	}
}
