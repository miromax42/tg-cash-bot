package logger

import "context"

type Logger interface {
	DebugfCtx(ctx context.Context, format string, args ...interface{})
	InfofCtx(ctx context.Context, format string, args ...interface{})
	PrintfCtx(ctx context.Context, format string, args ...interface{})
	WarnfCtx(ctx context.Context, format string, args ...interface{})
	ErrorfCtx(ctx context.Context, format string, args ...interface{})
	FatalfCtx(ctx context.Context, format string, args ...interface{})
	PanicfCtx(ctx context.Context, format string, args ...interface{})

	DebugCtx(ctx context.Context, args ...interface{})
	InfoCtx(ctx context.Context, args ...interface{})
	PrintCtx(ctx context.Context, args ...interface{})
	WarnCtx(ctx context.Context, args ...interface{})
	ErrorCtx(ctx context.Context, args ...interface{})
	FatalCtx(ctx context.Context, args ...interface{})
	PanicCtx(ctx context.Context, args ...interface{})
}
