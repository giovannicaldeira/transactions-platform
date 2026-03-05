package logger

import (
	"time"

	"github.com/gin-gonic/gin"
)

// GinLogger returns a middleware that logs HTTP requests
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		statusCode := c.Writer.Status()

		// Build log entry
		event := Logger.Info()

		if statusCode >= 500 {
			event = Logger.Error()
		} else if statusCode >= 400 {
			event = Logger.Warn()
		}

		event.
			Str("method", c.Request.Method).
			Str("path", path).
			Str("query", raw).
			Int("status", statusCode).
			Dur("latency", latency).
			Str("ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Int("body_size", c.Writer.Size())

		// Add error if present
		if len(c.Errors) > 0 {
			event.Str("error", c.Errors.String())
		}

		event.Msg("HTTP request")
	}
}

// GinRecovery returns a middleware that recovers from panics
func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				Logger.Error().
					Interface("error", err).
					Str("path", c.Request.URL.Path).
					Str("method", c.Request.Method).
					Msg("Panic recovered")
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
