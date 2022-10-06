package telegram

import (
	"context"

	tele "gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/telegram/errors"
)

// Authentication automatically sets user settings to context
func (s *Server) Authentication(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		settings, err := s.userSettings.Get(context.TODO(), c.Sender().ID)
		if err != nil || settings == nil {
			errors.SendError(c, errors.ErrInternal)
			return err
		}

		c.Set(SettingsKey.String(), settings)

		return next(c)
	}
}