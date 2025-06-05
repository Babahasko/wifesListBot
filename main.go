package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"shopping_bot/configs"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// TODO: Заменить эти временные хранилища на нормальную бд
// var userStates = make(map[int64]shop.UserState)
// var currentPurchaseData = make(map[int64]*shop.Purchase)
// var purchases = make([]shop.Purchase, 0)
// Хранилище покупок


func main() {
	// Config Loading

	conf := configs.LoadConfig()

	// Bot setup
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}

	b, err := bot.New(conf.BotToken, opts...)
	if nil != err {
		// panics for the sake of simplicity.
		// you should handle this error properly in your code.
		panic(err)
	}

	// b.RegisterHandler(bot.HandlerTypeMessageText, "foo", bot.MatchTypeCommand, shop.FooHandler)
	// b.RegisterHandler(bot.HandlerTypeMessageText, "bar", bot.MatchTypeCommandStartOnly, shop.BarHandler)

	slog.Info("buys bot started...")
	b.Start(ctx)
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	
}
