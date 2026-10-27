package logger

import (
	"log/slog"
	"os"
	"strings"
	"sync"
)

var (
	appLogger *slog.Logger
	initOnce  sync.Once
)

// Init creates the shared structured logger using the configured log level.
// Call this once when each service starts.
func Init(level string) {
	initOnce.Do(func() {
		var logLevel slog.Level

		switch strings.ToUpper(strings.TrimSpace(level)) {
		case "DEBUG":
			logLevel = slog.LevelDebug
		case "INFO":
			logLevel = slog.LevelInfo
		case "WARN", "WARNING":
			logLevel = slog.LevelWarn
		case "ERROR":
			logLevel = slog.LevelError
		default:
			logLevel = slog.LevelInfo
		}

		// Write structured logs to the console.
		handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})

		appLogger = slog.New(handler)
	})
}

// Get returns the shared logger.
// If Init has not been called, it initializes an INFO-level logger.
func Get() *slog.Logger {
	if appLogger == nil {
		Init("INFO")
	}

	return appLogger
}
