package telegram

import (
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
		return c.Send("pong!")
	})
	s.bot.Handle("/start", s.StartHelp)

	s.bot.Handle("/exp", s.CreateExpense)
	s.bot.Handle("/all", s.ListExpenses)

	s.bot.Handle("/currency", s.SelectCurrency)
	for _, b := range currencyButtonsUI {
		s.bot.Handle(&b, s.SetCurrency)
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
