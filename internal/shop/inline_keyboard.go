package shop

// import (
// 	"fmt"

// 	"shopping_bot/pkg/callback"

// 	"github.com/PaulSonOfLars/gotgbot/v2"
// )

// TODO: Здесь будет клавиатура формирования списка покупок

// func getCategoriesKeyboard(categories []string) (*gotgbot.InlineKeyboardMarkup, error) {
// 	var ButtonsPerRow = 3
// 	var rows [][]gotgbot.InlineKeyboardButton
// 	var buttons []gotgbot.InlineKeyboardButton
// 	for _, category := range categories {

// 		//Callback
// 		callbackStr, err := callback.PackCallback(&CategoryCallback{
// 			Name: category,
// 		})
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to pack callback: %v", err)
// 		}

// 		buttons = append(buttons, gotgbot.InlineKeyboardButton{
// 			Text:         category,
// 			CallbackData: callbackStr,
// 		})

// 		if len(buttons) == ButtonsPerRow {
// 			rows = append(rows, buttons)
// 			buttons = []gotgbot.InlineKeyboardButton{}
// 		}
// 	}

// 	if len(buttons) > 0 {
// 		rows = append(rows, buttons)
// 	}

// 	return &gotgbot.InlineKeyboardMarkup{
// 		InlineKeyboard: rows,
// 	}, nil
// }
