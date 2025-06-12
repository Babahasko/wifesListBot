package main

import (
	"shopping_bot/configs"
	"shopping_bot/pkg/logger"
)

// TODO: Заменить эти временные хранилища на нормальную бд
// var userStates = make(map[int64]shop.UserState)
// var currentPurchaseData = make(map[int64]*shop.Purchase)
// var purchases = make([]shop.Purchase, 0)

func main() {
	// Setup logger
	sugar := logger.NewSugarLogger()
	defer sugar.Sync()

	sugar.Infow("Starting the server")

	// Config Loading
	conf := configs.LoadConfig()
	sugar.Infow("Load configs",
		"token", conf.BotToken)
}
