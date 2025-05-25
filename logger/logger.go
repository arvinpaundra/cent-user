package logger

import (
	"context"
	// Assuming middleware constants are used
	"github.com/arvinpaundra/cent/user/api/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Get retrieves a zerolog.Logger from the given context.
// If no logger is found in the context, it returns the global zerolog.Logger.
func Get(ctx context.Context) zerolog.Logger {
	if ctx == nil {
		return log.Logger // Return global logger
	}
	logger, ok := ctx.Value(middleware.LoggerKey).(zerolog.Logger)
	if !ok {
		// Fallback to global logger if not found or type assertion fails
		return log.Logger
	}
	return logger
}

// CtxWithLogger returns a new context with the provided zerolog.Logger stored in it.
func CtxWithLogger(ctx context.Context, logger zerolog.Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, middleware.LoggerKey, logger)
}
