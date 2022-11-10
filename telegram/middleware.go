package telegram

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/logtags"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	tele "gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo/cache"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/telegram/tools"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

// Authentication automatically sets user settings to context
func (s *Server) Authentication(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		var (
			ctx       = requestContext(c)
			userID    = c.Sender().ID
			userToken = cache.UserSettingsToken(userID)

			userSettings repo.PersonalSettingsResp
		)

		err := s.dbCache.Get(ctx, userToken, &userSettings)
		if err != nil {
			if errors.Is(err, cache.ErrMiss) {
				settings, err := s.userSettings.Get(ctx, userID) //nolint:govet
				if err != nil || settings == nil {
					return s.SendError(err, c, tools.ErrInternal)
				}

				if err := s.dbCache.Set(ctx, userToken, settings); err != nil {
					return s.SendError(err, c, tools.ErrInternal)
				}

				userSettings = *settings
			} else {
				_ = s.SendError(err, c, tools.ErrInternal)

				return errors.Wrapf(err, "cache")
			}

		}

		c.Set(SettingsKey.String(), &userSettings)

		return next(c)
	}
}

func (s *Server) WithContext(ctx context.Context) func(next tele.HandlerFunc) tele.HandlerFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			rctx, span := otel.Tracer(util.RequestTrace).Start(ctx, "Telegram.Handler")
			defer span.End()

			rctx = logtags.AddTag(rctx, "trace.id", span.SpanContext().TraceID())
			rctx = logtags.AddTag(rctx, "request.user.id", c.Sender().ID)
			rctx = logtags.AddTag(rctx, "request.user.name", c.Sender().Username)
			rctx = logtags.AddTag(rctx, "request.message", c.Message().Text)

			c.Set(RequestContextKey.String(), rctx)

			err := next(c)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())

				return err
			}

			return nil
		}
	}
}

func (s *Server) Logger(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		s.logger.Info(requestContext(c), "request msg")

		return next(c)
	}

}

func (s *Server) DropUserSettingsCache(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		var (
			ctx       = requestContext(c)
			userToken = cache.UserSettingsToken(c.Sender().ID)
		)

		err := next(c)
		if err != nil {
			return err
		}

		err = s.dbCache.Del(ctx, userToken)

		return errors.Wrapf(err, "delete key %q", string(userToken))
	}
}
