package log

import (
	"context"
	"log/slog"
)

type ContextLogKey struct{}

func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ContextLogKey{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}

func WithContext(ctx context.Context, l *slog.Logger) context.Context {
	if _, ok := ctx.Value(ContextLogKey{}).(*slog.Logger); ok {
		return ctx
	}
	return context.WithValue(ctx, ContextLogKey{}, l)
}
