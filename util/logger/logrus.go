package logger

import (
	"context"

	"github.com/cockroachdb/logtags"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

type Adapter struct {
	l *logrus.Logger
}

func New() *Adapter {
	return &Adapter{l: logrus.New()}
}

func NewTest() *Adapter {
	l, _ := test.NewNullLogger()

	return &Adapter{l}
}

func (a *Adapter) Debugf(ctx context.Context, format string, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Debugf(format, args...)
}

func (a *Adapter) Infof(ctx context.Context, format string, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Infof(format, args...)
}

func (a *Adapter) Printf(ctx context.Context, format string, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Printf(format, args...)
}

func (a *Adapter) Warnf(ctx context.Context, format string, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Warnf(format, args...)
}

func (a *Adapter) Errorf(ctx context.Context, format string, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Errorf(format, args...)
}

func (a *Adapter) Fatalf(ctx context.Context, format string, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Fatalf(format, args...)
}

func (a *Adapter) Panicf(ctx context.Context, format string, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Panicf(format, args...)
}

func (a *Adapter) Debug(ctx context.Context, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Debug(args...)
}

func (a *Adapter) Info(ctx context.Context, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Info(args...)
}

func (a *Adapter) Print(ctx context.Context, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Print(args...)
}

func (a *Adapter) Warn(ctx context.Context, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Warn(args...)
}

func (a *Adapter) Error(ctx context.Context, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Error(args...)
}

func (a *Adapter) Fatal(ctx context.Context, args ...interface{}) {
	a.l.WithFields(mapFromContext(ctx)).Fatal(args...)
}

func (a *Adapter) Panic(ctx context.Context, args ...interface{}) {
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
