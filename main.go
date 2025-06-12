package main

import (
	"shopping_bot/configs"
	"shopping_bot/internal/shop"
	"shopping_bot/pkg/logger"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// TODO: Заменить эти временные хранилища на нормальную бд
// var userStates = make(map[int64]shop.UserState)
// var currentPurchaseData = make(map[int64]*shop.Purchase)
// var purchases = make([]shop.Purchase, 0)

func main() {
	// == Setup LOGGER ==
	sugar := logger.NewSugarLogger()
	defer sugar.Sync()

	// == Load CONFIGS ==
	conf := configs.LoadConfig()
	sugar.Infow("Load configs")

	// == Setup TELEGRAM BOT ==
	b, err := gotgbot.NewBot(conf.BotToken, nil)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// сreate updater and dispatcher.
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			sugar.Errorf("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)

	// Handlers
	shop.NewShopHandler(dispatcher)
	
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
	sugar.Infof("%s has been started...", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
	// Command Handlers
}
