package telegram

import (
	"context"

	"github.com/cockroachdb/logtags"
	"github.com/google/uuid"
	tele "gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/telegram/tools"
)

// Authentication automatically sets user settings to context
func (s *Server) Authentication(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		settings, err := s.userSettings.Get(requestContext(c), c.Sender().ID)
		if err != nil || settings == nil {
			return tools.SendError(err, c, tools.ErrInternal)
		}

		c.Set(SettingsKey.String(), settings)

		return next(c)
	}
}

func (s *Server) WithContext(ctx context.Context) func(next tele.HandlerFunc) tele.HandlerFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			// request level context
			rctx, cancel := context.WithCancel(ctx)
			defer cancel()

			rctx = logtags.AddTag(rctx, "request.id", uuid.New().String())
			rctx = logtags.AddTag(rctx, "request.user.id", c.Sender().ID)
			rctx = logtags.AddTag(rctx, "request.user.name", c.Sender().Username)

			c.Set(RequestContextKey.String(), rctx)

			return next(c)
		}
	}
}
