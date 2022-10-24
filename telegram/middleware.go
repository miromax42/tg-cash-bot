package telegram

import (
	"context"

	"github.com/cockroachdb/logtags"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	tele "gopkg.in/telebot.v3"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/telegram/tools"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
)

// Authentication automatically sets user settings to context
func (s *Server) Authentication(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		settings, err := s.userSettings.Get(requestContext(c), c.Sender().ID)
		if err != nil || settings == nil {
			return s.SendError(err, c, tools.ErrInternal)
		}

		c.Set(SettingsKey.String(), settings)

		return next(c)
	}
}

func (s *Server) WithContext(ctx context.Context) func(next tele.HandlerFunc) tele.HandlerFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			rctx, span := otel.Tracer(util.RequestTrace).Start(ctx, "telegram_handler")
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
		s.logger.InfoCtx(requestContext(c), "request msg")

		return next(c)
	}

}
