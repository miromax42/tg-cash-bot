package telegram

import (
	"context"
	"fmt"
	"time"

	tele "gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/util/logger"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util/metrics"

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
		return s.Send(c, fmt.Sprintf("pong! Your id: %d", c.Sender().ID))
	})
	s.bot.Handle("/start", instrumentedHandler("start", StartHelp))

	s.bot.Handle("/exp", instrumentedHandler("expense_add", s.CreateExpense))
	s.bot.Handle("/all", instrumentedHandler("expense_list", s.ListExpenses))

	currencySelectorUI, anyButtonUI := getCurrencySelector()
	s.bot.Handle("/currency", instrumentedHandler("currency_menu", s.SelectCurrency(currencySelectorUI)))
	s.bot.Handle(anyButtonUI, instrumentedHandler("currency_set", s.SetCurrency))

	s.bot.Handle("/limit", instrumentedHandler("limit_set", s.SetLimit))
}

func instrumentedHandler(label string, fn func(c tele.Context) error) func(c tele.Context) error {
	return func(c tele.Context) error {
		metrics.RequestOpsProcessed.WithLabelValues(label).Inc()

		err := fn(c)
		if err != nil {
			metrics.RequestOpsInternalError.WithLabelValues(label).Inc()
		}

		metrics.RequestDuration.WithLabelValues(label).Observe(
			time.Since(c.Message().Time()).Seconds(),
		)

		return err
	}
}

func (s *Server) Start() {
	s.bot.Start()
}

func (s *Server) Stop() {
	if s != nil && s.bot != nil {
		s.bot.Stop()
	}
}
