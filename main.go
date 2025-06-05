package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"shopping_bot/configs"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
    "github.com/go-telegram/ui/keyboard/reply"
)

// TODO: Заменить эти временные хранилища на нормальную бд
// var userStates = make(map[int64]shop.UserState)
// var currentPurchaseData = make(map[int64]*shop.Purchase)
// var purchases = make([]shop.Purchase, 0)

// Keyboards

var AddPurchaseKeyboard *reply.ReplyKeyboard

func main() {
	// Config Loading
	conf := configs.LoadConfig()

    // Initkeyboards

	// Bot setup
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}

	b, err := bot.New(conf.BotToken, opts...)

    // Init keyboards
    initReplyKeyboard(b)

    //Register handlers
    b.RegisterHandler(bot.HandlerTypeMessageText, "/replykb", bot.MatchTypeExact, handlerReplyKeyboard)

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

func initReplyKeyboard(b *bot.Bot) {
	AddPurchaseKeyboard = reply.New(
		reply.WithPrefix("reply_keyboard"),
		reply.IsSelective(),
		reply.IsOneTimeKeyboard(),
	).
		Button("Button", b, bot.MatchTypeExact, onReplyKeyboardSelect).
		Row().
		Button("Cancel", b, bot.MatchTypeExact, onReplyKeyboardSelect)
}

func handlerReplyKeyboard(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Select example command from reply keyboard:",
		ReplyMarkup: AddPurchaseKeyboard,
	})
}

func onReplyKeyboardSelect(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "You selected: " + string(update.Message.Text),
	})
}
