package telegram

import (
	"fmt"
	"time"

	tele "gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/currency"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

type Server struct {
	logger       util.Logger
	bot          *tele.Bot
	expense      repo.Expense
	userSettings repo.PersonalSettings
	exchange     currency.Exchange
}

func NewServer(
	cfg util.ConfigTelegram,
	log util.Logger,
	expense repo.Expense,
	userSettings repo.PersonalSettings,
	exchange currency.Exchange,
) (*Server, error) {
	pref := tele.Settings{
		Token:  cfg.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
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

	srv.setupRoutes()

	return srv, nil
}

func (s *Server) setupRoutes() {
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
