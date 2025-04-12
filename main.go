package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/ulngollm/quizbot/flow/feedback"
	"github.com/ulngollm/teleflow"
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
	if err := run(t); err != nil {
		log.Printf("run: %s", err)
	}
}

func run(token string) error {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: time.Second},
	}
	bot, err := tele.NewBot(pref)
	if err != nil {
		return fmt.Errorf("tele.NewBot: %v", err)
	}

	pool := teleflow.NewMemoryPool()
	flowManager := teleflow.NewFlowManager(pool)
	router := teleflow.NewFlowRouter(flowManager)

	newUserCtrl := feedback.New(flowManager)

	ng := router.Group(feedback.FlowName)
	//Куда удобнее направлние стейтов задавать все-таки снаружи
	//иначк одна задача размывается по разным файлам, между ними приходится прынать и переключаться
	ng.AddHandler(feedback.StateUndefined, feedback.StateAskedCategory, newUserCtrl.AskCategory)
	ng.AddHandler(feedback.StateAskedCategory, feedback.StateAskedProduct, newUserCtrl.AskProduct)
	ng.AddHandler(feedback.StateAskedProduct, feedback.StateAskedDetails, newUserCtrl.AskDetails)
	ng.AddHandler(feedback.StateAskedDetails, feedback.StateAskedScreenshot, newUserCtrl.AskScreenshot)
	ng.AddHandler(feedback.StateAskedScreenshot, feedback.StateComplete, newUserCtrl.Thank)

	bot.Handle("/feedback", func(c tele.Context) error {
		return bot.Trigger(tele.OnText, c)
	}, newUserCtrl.Init)
	bot.Handle(tele.OnText, handle, router.Middleware())

	bot.Start()

	return nil
}

func handle(c tele.Context) error {
	return c.Send("выберите команду")
}
