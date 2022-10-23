package logger

import (
	"context"

	"github.com/cockroachdb/logtags"
	"github.com/sirupsen/logrus"
)

type Adapter struct {
	l *logrus.Logger
}

func New() *Adapter {
	return &Adapter{l: logrus.New()}
}

func (a *Adapter) DebugfCtx(ctx context.Context, format string, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Debugf(format, args...)
}

func (a *Adapter) InfofCtx(ctx context.Context, format string, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Infof(format, args...)
}

func (a *Adapter) PrintfCtx(ctx context.Context, format string, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Printf(format, args...)
}

func (a *Adapter) WarnfCtx(ctx context.Context, format string, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Warnf(format, args...)
}

func (a *Adapter) ErrorfCtx(ctx context.Context, format string, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Errorf(format, args...)
}

func (a *Adapter) FatalfCtx(ctx context.Context, format string, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Fatalf(format, args...)
}

func (a *Adapter) PanicfCtx(ctx context.Context, format string, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Panicf(format, args...)
}

func (a *Adapter) DebugCtx(ctx context.Context, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Debug(args...)
}

func (a *Adapter) InfoCtx(ctx context.Context, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Info(args...)
}

func (a *Adapter) PrintCtx(ctx context.Context, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Print(args...)
}

func (a *Adapter) WarnCtx(ctx context.Context, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Warn(args...)
}

func (a *Adapter) ErrorCtx(ctx context.Context, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Error(args...)
}

func (a *Adapter) FatalCtx(ctx context.Context, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Fatal(args...)
}

func (a *Adapter) PanicCtx(ctx context.Context, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Panic(args...)
}

func mapFromContext(ctx context.Context) map[string]interface{} {
	m := make(map[string]interface{})

	buffer := logtags.FromContext(ctx)
	if buffer == nil {
		return m
	}

	for _, v := range buffer.Get() {
		m[v.Key()] = v.Value()
	}

	return m
}
