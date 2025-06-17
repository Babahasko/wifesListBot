package shop

import "github.com/PaulSonOfLars/gotgbot/v2"

const (
	ButtonAddPurchase  = "✍️Сформировать список покупок"
	ButtonViewList     = "📋Посмотреть списки покупок"
	ButtonCancel       = "❌Отмена"
	ButtonBack         = "⬅️Назад"
)

func getMainMenueKeyboard() *gotgbot.ReplyKeyboardMarkup {
	return &gotgbot.ReplyKeyboardMarkup{
		OneTimeKeyboard: true,
		ResizeKeyboard: true,
		Keyboard: [][]gotgbot.KeyboardButton{
			{
				{Text: ButtonAddPurchase},
				{Text: ButtonViewList},
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
