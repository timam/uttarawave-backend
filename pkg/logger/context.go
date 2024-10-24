package logger

import (
	"context"
	"go.uber.org/zap"
)

type ctxKey struct{}

// WithLogger returns a new context with the given logger attached
func WithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

// FromContext returns the logger associated with the context,
// or the global logger if no logger is associated with the context
func FromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return l
	}
	return GetLogger()
}
