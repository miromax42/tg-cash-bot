package telegram

import (
	"github.com/cockroachdb/errors"
	tele "gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/telegram/tools"
)

func (s *Server) SendError(err error, c tele.Context, e tools.UserError) error {
	terr := c.Send(e.With(err).Error())

	if terr != nil {
		err = errors.CombineErrors(err, errors.WithHint(terr, "during handling main error"))
	} else if !e.IsNotExpected {
		s.logger.WarnCtx(requestContext(c), e.Title)

		return nil
	}

	return err
}
