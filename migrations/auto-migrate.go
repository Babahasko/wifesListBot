package main

import (
	"log"
	"os"
	"shopping_bot/internal/repository"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(`D:\Programming\6_Private_Project\BuysBot\.env`)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")),  &gorm.Config{
	})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&repository.GormShoppingList{}, &repository.GormShoppingItem{}, &repository.GormUserState{})
	if err != nil {
        log.Printf("Ошибка при выполнении миграций: %v", err)
		return
    }
	log.Println("Миграции успешно выполнены")
}