package main

import (
	"log"

	"tg-bot-boilerplate/internal/bot"
	"tg-bot-boilerplate/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	b, err := bot.New(cfg.BotToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	b.RegisterHandlers()
	b.Start()
}
