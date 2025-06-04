package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"shopping_bot/configs"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Send any text message to the bot after the bot has been started

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
	if update.Message == nil {
		return
	}
	kb := models.ReplyKeyboardMarkup{
		Keyboard:              [][]models.KeyboardButton{{models.KeyboardButton{Text: "hide"}}},
		IsPersistent:          false,
		ResizeKeyboard:        true,
		OneTimeKeyboard:       false,
		InputFieldPlaceholder: "",
		Selective:             false,
	}
	_, errSend := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: `Приветствую! Я позволю тебе хранить список покупок.
		Ты можешь выбрать для них категорию, приоритет и цену.
		Получить список покупок, отфильтровать их. И конечно же посмотреть список своих больших свершений.`,
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: &kb,
	})

	if errSend != nil {
		fmt.Printf("error sending message: %v\n", errSend)
		return
	}
}
