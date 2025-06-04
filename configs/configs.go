package configs

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("Error loading .env file, using default config")
		os.Exit(1)
	}
	return &Config{
		BotToken: os.Getenv("BOT_TOKEN"),
	}
}