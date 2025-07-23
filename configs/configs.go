package configs

import (
	"os"

	"github.com/joho/godotenv"
	"shopping_bot/pkg/logger"
)

type Config struct {
	BotToken string
	DB DBConfig
}

type DBConfig struct {
	DSN string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Sugar.Error("Error loading .env file, using default config")
		os.Exit(1)
	}
	return &Config{
		BotToken: os.Getenv("BOT_TOKEN"),
		DB: DBConfig {
			DSN: os.Getenv("DSN"),
		},
	}
}