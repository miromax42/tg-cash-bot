package telegram

import (
	"context"
	"fmt"
	"time"

	tele "gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/util/logger"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/currency"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

type Server struct {
	logger       logger.Logger
	bot          *tele.Bot
	expense      repo.Expense
	userSettings repo.PersonalSettings
	exchange     currency.Exchange
}

func NewServer(
	ctx context.Context,
	cfg util.ConfigTelegram,
	log logger.Logger,
	expense repo.Expense,
	userSettings repo.PersonalSettings,
	exchange currency.Exchange,
) (*Server, error) {
	pref := tele.Settings{
		Token:  cfg.Token,
		Poller: &tele.LongPoller{Timeout: time.Second},
		OnError: func(err error, c tele.Context) {
			log.ErrorCtx(requestContext(c), fmt.Sprintf("%+v", err))
		},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	srv := &Server{
		logger:       log,
		bot:          bot,
		expense:      expense,
		userSettings: userSettings,
		exchange:     exchange,
	}

	srv.setupRoutes(ctx)

	return srv, nil
}

func (s *Server) setupRoutes(ctx context.Context) {
	s.bot.Use(s.WithContext(ctx))
	s.bot.Use(s.Logger)
	s.bot.Use(s.Authentication)

	s.bot.Handle("/ping", func(c tele.Context) error {
		return c.Send(fmt.Sprintf("pong! Your id: %d", c.Sender().ID))
	})
	s.bot.Handle("/start", StartHelp)

	s.bot.Handle("/exp", s.CreateExpense)
	s.bot.Handle("/all", s.ListExpenses)

	currencySelectorUI, anyButtonUI := getCurrencySelector()
	s.bot.Handle("/currency", s.SelectCurrency(currencySelectorUI))
	s.bot.Handle(anyButtonUI, s.SetCurrency)

	s.bot.Handle("/limit", s.SetLimit)
}

func (s *Server) Start() {
	s.bot.Start()
}

func (s *Server) Stop() {
	if s != nil && s.bot != nil {
		s.bot.Stop()
	}
}
