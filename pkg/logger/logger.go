package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// NewLogger creates and configures a new zerolog.Logger instance.
// It's configured to provide human-readable "pretty" logs
// if the env is "development", and structured JSON logs otherwise.
func NewLogger(env string) zerolog.Logger {
	var logger zerolog.Logger

	if env == "development" {
		// Use console writer for human-readable logs in development
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339, // A more readable time format
		}
		logger = zerolog.New(consoleWriter).With().Timestamp().Logger()
		// Set global level to Debug for development
		logger = logger.Level(zerolog.DebugLevel)
	} else {
		// Use standard JSON output for production/staging
		// This is the default zerolog behavior
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix // Use Unix timestamp
		logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
		// Set global level to Info for production
		logger = logger.Level(zerolog.InfoLevel)
	}

	return logger
}
