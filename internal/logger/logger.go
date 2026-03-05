package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

// Init initializes the global logger
func Init() {
	// Set log level from environment or default to info
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = time.RFC3339

	// Use pretty console output in development, JSON in production
	if os.Getenv("APP_ENV") == "development" {
		Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	} else {
		Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
}

// Info logs an info level message
func Info(msg string) *zerolog.Event {
	return Logger.Info().Str("msg", msg)
}

// Error logs an error level message
func Error(msg string) *zerolog.Event {
	return Logger.Error().Str("msg", msg)
}

// Debug logs a debug level message
func Debug(msg string) *zerolog.Event {
	return Logger.Debug().Str("msg", msg)
}

// Warn logs a warn level message
func Warn(msg string) *zerolog.Event {
	return Logger.Warn().Str("msg", msg)
}

// Fatal logs a fatal level message and exits
func Fatal(msg string) *zerolog.Event {
	return Logger.Fatal().Str("msg", msg)
}
