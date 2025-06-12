package main

import (
	"fmt"
	"shopping_bot/configs"
)

// TODO: Заменить эти временные хранилища на нормальную бд
// var userStates = make(map[int64]shop.UserState)
// var currentPurchaseData = make(map[int64]*shop.Purchase)
// var purchases = make([]shop.Purchase, 0)

func main() {
	// Config Loading
	conf := configs.LoadConfig()
	fmt.Println(conf.BotToken)
}