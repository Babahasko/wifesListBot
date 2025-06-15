package shop

import "github.com/PaulSonOfLars/gotgbot/v2"

const (
	ButtonAddPurchase  = "🛒Добавить покупку"
	ButtonViewList     = "📋Посмотреть список покупок"
	ButtonAchivemenets = "🎉Свершения"
	ButtonCategories   = "📁Категории"
	ButtonCancel       = "❌Отмена"
	ButtonBack         = "⬅️Назад"
)

func getMainMenueKeyboard() *gotgbot.ReplyKeyboardMarkup {
	return &gotgbot.ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard: [][]gotgbot.KeyboardButton{
			{
				{Text: ButtonAddPurchase},
				{Text: ButtonViewList},
			},
			{
				{Text: ButtonAchivemenets},
				{Text: ButtonCategories},
			},
		},
	}
}

func getMenueKeyboard() *gotgbot.ReplyKeyboardMarkup {
	return &gotgbot.ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard: [][]gotgbot.KeyboardButton{
			{
				{Text: ButtonBack},
			},
			{
				{Text: ButtonCancel},
			},
		},
	}
}

func getShortMenueKeyboard() *gotgbot.ReplyKeyboardMarkup {
	return &gotgbot.ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard: [][]gotgbot.KeyboardButton{
			{
				{Text: ButtonCancel},
			},
		},
	}
}
