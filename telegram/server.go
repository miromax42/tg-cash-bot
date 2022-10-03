package telegram

import (
	"time"

	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

type Server struct {
	logger  util.Logger
	bot     *tele.Bot
	expense repo.Expense
}

func NewServer(cfg util.ConfigTelegram, log util.Logger, expense repo.Expense) *Server {
	pref := tele.Settings{
		Token:  cfg.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	srv := &Server{
		logger:  log,
		bot:     bot,
		expense: expense,
	}

	srv.setupRoutes()

	return srv
}

func (s *Server) setupRoutes() {
	s.bot.Use(middleware.Logger())

	s.bot.Handle("/ping", func(c tele.Context) error {
		return c.Send("pong!")
	})
	s.bot.Handle("/start", s.StartHelp)

	s.bot.Handle("/exp", s.CreateExpense)
	s.bot.Handle("/all", s.ListExpenses)
}

func (s *Server) Start() {
	s.bot.Start()
}
