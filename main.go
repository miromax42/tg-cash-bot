package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	_ "github.com/mattn/go-sqlite3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/currency"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/currency/exhange"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/ent"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo/database"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/telegram"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

func main() {
	var (
		db       *ent.Client
		exchange currency.Exchange
		srv      *telegram.Server
	)

	mainCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log := logrus.New()

	cfg, err := util.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	init, initCtx := errgroup.WithContext(mainCtx)

	init.Go(func() (err error) {
		db, err = ent.Open("sqlite3", "file:test.db?_fk=1")
		if err != nil {
			return fmt.Errorf("failed opening connection to sqlite: %w", err)
		}

		if err := db.Schema.Create(initCtx); err != nil {
			return fmt.Errorf("failed creating schema resources: %w", err)
		}

		return nil
	})

	init.Go(func() (err error) {
		if exchange, err = exhange.New(initCtx, cfg.Exchange); err != nil {
			return fmt.Errorf("init exchange: %w", err)
		}

		return nil
	})

	if err := init.Wait(); err != nil {
		log.Fatalf("exit reason: %s \n", err)
	}

	work, workCtx := errgroup.WithContext(mainCtx)

	work.Go(func() (err error) {
		expense := database.NewExpense(db)
		personalSettings := database.NewPersonalSettings(db)
		srv, err = telegram.NewServer(cfg.Telegram, log, expense, personalSettings, exchange)
		if err != nil {
			return err
		}

		srv.Start()

		return nil
	})

	work.Go(func() error {

		log.Info("bot started")

		<-workCtx.Done()
		srv.Stop()
		if err := db.Close(); err != nil {
			return fmt.Errorf("close db: %w", err)
		}

		log.Info("bot stopped")

		return nil
	})

	if err := work.Wait(); err != nil {
		log.Fatalf("gracefull stop: %s \n", err)
	} else {
		log.Info("gracefully stopped!")
	}
}
