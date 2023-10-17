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
