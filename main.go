package main

import (
	"errors"
	"fmt"
	"shopping_bot/configs"
	"shopping_bot/internal/bot"
	"shopping_bot/internal/repository"
	"shopping_bot/internal/service"
	"shopping_bot/pkg/db"
	"shopping_bot/pkg/logger"
	"shopping_bot/pkg/middleware"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

// TODO: Заменить эти временные хранилища на нормальную бд

func main() {
	// == Setup LOGGER ==
	logger.InitLogger()
	defer logger.Sugar.Sync()

	// == Load CONFIGS ==
	conf := configs.LoadConfig()
	logger.Sugar.Infow("Load configs")

	// == Middlewares ==
	middlewareChain := middleware.Chain(
		middleware.NewLoggingMiddleware(),
		middleware.NewErrorMiddleware(),
	)

	// == Setup TELEGRAM BOT ==
	b, err := gotgbot.NewBot(conf.BotToken, &gotgbot.BotOpts{
		BotClient: middleware.NewMiddlewareClient(middlewareChain),
	})

	if err != nil {
		// panic("failed to create new bot: " + err.Error())
		logger.Sugar.Fatal(errors.New("invalid token"), err)
	}

	// сreate updater and dispatcher.
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			logger.Sugar.Errorf("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)

	//Repository
	db := db.NewDB(conf)
	shopRepo := repository.NewPostgresShoppingRepository(db.DB)

	// Service
	service := service.NewShopService(shopRepo)

	// Handlers
	bot.NewShopHandler(dispatcher, service)

	dispatcher.AddHandler(handlers.NewMessage(message.All, dafaultHandler))

	// Start receiving updates.
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})

	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	logger.Sugar.Infof("%s has been started...", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
	// Command Handlers
}

func dafaultHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, "Необработанное событие", nil)
	if err != nil {
		return fmt.Errorf("failed to handle message: %w", err)
	}
	return nil
}
