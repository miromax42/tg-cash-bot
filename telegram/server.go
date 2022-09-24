package telegram

import (
	"time"

	tele "gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

type Server struct {
	logger util.Logger
	bot    *tele.Bot
}

func NewServer(cfg util.ConfigTelegram, log util.Logger) *Server {
	pref := tele.Settings{
		Token:  cfg.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	srv := &Server{
		logger: log,
		bot:    bot,
	}

	srv.setupRoutes()

	return srv
}

func (s *Server) setupRoutes() {
	s.bot.Handle("/ping", func(c tele.Context) error {
		return c.Send("pong!")
	})
}

func (s *Server) Start() {
	s.bot.Start()
}
