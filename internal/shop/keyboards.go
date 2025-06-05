package shop

import "github.com/go-telegram/bot/models"

func GetMainKeyboard() models.ReplyKeyboardMarkup {
    return models.ReplyKeyboardMarkup{
        Keyboard: [][]models.KeyboardButton{
            {
                {Text: "🛒 Добавить покупку"},
                {Text: "📋 Показать список"},
            },
            {
                {Text: "🗑 Очистить список"},
            },
        },
        ResizeKeyboard: true,
    }
}