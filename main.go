package main

import (
	"github.com/sirupsen/logrus"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/telegram"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

func main() {
	var log util.Logger = logrus.New()

	cfg, err := util.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	srv := telegram.NewServer(cfg.Telegram, log)
	srv.Start()
}
