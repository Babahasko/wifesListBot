package shop

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Состояния FSM
type UserState int

const (
	WaitingForName UserState = iota
	WaitingForPrice
	WaitingForCategory
	WaitingForPriority
	Confirming
)

// TODO:
//Временное хранилище
// var userStates = make(map[int64]UserState)

// Временный слайс покупок
// var purchases = make([]Purchase, 0)


func AddItem(ctx context.Context, b *bot.Bot, update *models.Update) {

}

func FooHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "Caught *foo*",
		ParseMode: models.ParseModeMarkdown,
	})
}

func BarHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "Caught *bar*",
		ParseMode: models.ParseModeMarkdown,
	})
}