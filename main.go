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

	// client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	// if err != nil {
	// 	log.Fatalf("failed opening connection to sqlite: %v", err)
	// }
	// defer client.Close()
	// // Run the auto migration tool.
	// if err := client.Schema.Create(context.Background()); err != nil {
	// 	log.Fatalf("failed creating schema resources: %v", err)
	// }

	srv := telegram.NewServer(cfg.Telegram, log)
	srv.Start()
}
