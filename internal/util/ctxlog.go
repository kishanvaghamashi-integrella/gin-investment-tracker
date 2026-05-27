package util

import (
	"context"

	"go.uber.org/zap"
)

type ctxLogKey struct{}

func WithLogger(ctx context.Context, log *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, ctxLogKey{}, log)
}

func FromContext(ctx context.Context) *zap.SugaredLogger {
	if log, ok := ctx.Value(ctxLogKey{}).(*zap.SugaredLogger); ok {
		return log
	}
	return Logger
}
