package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"shopping_bot/configs"
	"shopping_bot/internal/shop"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// TODO: Заменить эти временные хранилища на нормальную бд
// var userStates = make(map[int64]shop.UserState)
// var currentPurchaseData = make(map[int64]*shop.Purchase)
// var purchases = make([]shop.Purchase, 0)

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

	// ==HANDLERS==
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/add", bot.MatchTypeExact, addPurchaseHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/replykb", bot.MatchTypeExact, handlerReplyKeyboard)

	if nil != err {
		// panics for the sake of simplicity.
		// you should handle this error properly in your code.
		panic(err)
	}

	// b.RegisterHandler(bot.HandlerTypeMessageText, "foo", bot.MatchTypeCommand, shop.FooHandler)
	// b.RegisterHandler(bot.HandlerTypeMessageText, "bar", bot.MatchTypeCommandStartOnly, shop.BarHandler)
	botInfo, err := b.GetMe(ctx)
	if err != nil {
		slog.Error("getMe",
			slog.String("error", err.Error()),
		)
	}
	slog.Info("buys bot started:%s",
		slog.String("botName", botInfo.Username))
	b.Start(ctx)
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Error("unhandeld update")
}

func handlerReplyKeyboard(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := shop.InitPurchaseKB(b)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Select example command from reply keyboard:",
		ReplyMarkup: kb,
	})
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: "Привет! Когда нибудь тут будет справка",
	})
}

func addPurchaseHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := shop.InitPurchaseKB(b)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: "Выберите действие",
		ReplyMarkup: kb,
	})
}
