package telegram

import (
	"github.com/cockroachdb/errors"
	"go.opentelemetry.io/otel"
	tele "gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/telegram/tools"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

func (s *Server) SendError(err error, c tele.Context, e tools.UserError) error {
	terr := s.Send(c, e.With(err).Error())

	if terr != nil {
		err = errors.CombineErrors(err, errors.WithHint(terr, "during handling main error"))
	} else if !e.IsNotExpected {
		s.logger.WarnCtx(requestContext(c), e.Title)

		return nil
	}

	return err
}

func (s *Server) Send(c tele.Context, what interface{}, opts ...interface{}) error {
	_, span := otel.Tracer(util.RequestTrace).Start(requestContext(c), "Telegram.Send")
	defer span.End()

	return c.Send(what, opts...)
}
