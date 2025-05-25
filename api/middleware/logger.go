package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const (
	LoggerKey    = "logger"
	RequestIDKey = "request_id"
)

// Logger is a Gin middleware that injects a zerolog.Logger into the context.
// Each logger instance will have a unique request_id field.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a logger instance from the global logger
		logger := log.Logger

		// Generate a new request ID
		requestID := uuid.NewString()

		// Add request ID and other request details to the logger's context
		logger = logger.With().
			Str(RequestIDKey, requestID).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("ip", c.ClientIP()).
			Logger()

		// Store the logger in the Gin context
		c.Set(LoggerKey, logger)

		// Add request_id to response header
		c.Header(RequestIDKey, requestID)

		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Log request completion
		latency := time.Since(start)
		logger.Info().
			Int("status", c.Writer.Status()).
			Dur("latency", latency).
			Msg("Request completed")
	}
}
