package shop

import "github.com/PaulSonOfLars/gotgbot/v2"

func getMainMenueKeyboard() *gotgbot.ReplyKeyboardMarkup {
	return &gotgbot.ReplyKeyboardMarkup{
		OneTimeKeyboard: true,
		ResizeKeyboard:  true,
		Keyboard: [][]gotgbot.KeyboardButton{
			{
				{Text: ButtonAddPurchase},
				{Text: ButtonViewList},
			},
		},
	}
}

func getCancelKeyboard() *gotgbot.ReplyKeyboardMarkup {
	return &gotgbot.ReplyKeyboardMarkup{
		IsPersistent:   true,
		ResizeKeyboard: true,
		Keyboard: [][]gotgbot.KeyboardButton{
			{
				{Text: ButtonCancel},
			},
		},
	}
}

func getFormListKeyboard() *gotgbot.ReplyKeyboardMarkup {
	return &gotgbot.ReplyKeyboardMarkup{
		IsPersistent:   true,
		ResizeKeyboard: true,
		Keyboard: [][]gotgbot.KeyboardButton{
			{
				{Text: ButtonFinishList},
				{Text: ButtonCancel},
			},
		},
	}
}
