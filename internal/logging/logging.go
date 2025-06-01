package logging

import (
	"context"
	"log/slog"

	"github.com/evenlwanvik/smartsplit/internal/common"
)

const LoggerCtxKey common.ContextKey = "logger"

// WithLogger embeds a logger in the given context.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, LoggerCtxKey, logger)
}

// LoggerFromContext attempts to extract an embedded logger from the
// given context. If no context is found, it returns the default logger
// registered for the application.
func LoggerFromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(LoggerCtxKey).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return logger
}
