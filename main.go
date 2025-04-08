package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("godotenv.Load: %s", err)
		return
	}
}

func main() {
	t, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		log.Fatalf("bot token is empty")
		return
	}

	pref := tele.Settings{
		Token:  t,
		Poller: &tele.LongPoller{Timeout: time.Second},
	}
	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("tele.NewBot: %s", err)
		return
	}

	bot.Handle("/start", handle)

	bot.Start()
}

func handle(c tele.Context) error {
	return c.Send(c.Message())
}
