package main

import (
	"context"

	"github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/ent"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo/database"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/telegram"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

func main() {
	var log util.Logger = logrus.New()

	cfg, err := util.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := ent.Open("sqlite3", "file:test.db?_fk=1")
	if err != nil {
		log.Fatal("failed opening connection to sqlite: ", err)
	}
	defer db.Close()

	if err := db.Schema.Create(context.Background()); err != nil {
		log.Fatal("failed creating schema resources: ", err)
	}

	expense := database.NewExpense(db)
	personalSettings := database.NewPersonalSettings(db)
	srv := telegram.NewServer(cfg.Telegram, log, expense, personalSettings)

	log.Info("Started")
	srv.Start()
}
