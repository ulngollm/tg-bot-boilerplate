package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	botToken, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		log.Fatal("BOT_TOKEN is not set")
	}

	return &Config{
		BotToken: botToken,
	}, nil
}
