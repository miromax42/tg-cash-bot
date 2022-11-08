package gapi

import (
	tele "gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/pb"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util/logger"
)

type Server struct {
	pb.UnimplementedBotSendServer
	t   *tele.Bot
	log logger.Logger
}

func New(t *tele.Bot, log logger.Logger) (*Server, error) {
	server := &Server{
		t:   t,
		log: log,
	}

	return server, nil
}
