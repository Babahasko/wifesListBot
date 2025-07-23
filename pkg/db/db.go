package db

import (
	"shopping_bot/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct{
	*gorm.DB
}

func NewDB(conf *configs.Config) *Db {
	db, err := gorm.Open(postgres.Open(conf.DB.DSN),  &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &Db{db}
}