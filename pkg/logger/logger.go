package logger

import (
	"os"
	"sync"
	"time"

	slogformatter "github.com/samber/slog-formatter"
	C "github.com/suyog1pathak/services/pkg/config"
	"log/slog"
)

// Singleton implementation of logger slog
var once sync.Once
var logger *slog.Logger

// createLogger creates slog logger instance with custom config.
func createLogger(level string) {
	// set logger level from env var
	var hOptions slog.HandlerOptions
	switch level {
	case "DEBUG":
		hOptions = slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
	case "INFO":
		hOptions = slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
	case "ERROR":
		hOptions = slog.HandlerOptions{
			Level: slog.LevelError,
		}
	default:
		hOptions = slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
	}
	// create logger
	logger = slog.New(
		slogformatter.NewFormatterHandler(
			slogformatter.TimezoneConverter(time.Local),
			slogformatter.TimeFormatter(time.RFC3339, nil),
		)(
			slog.NewJSONHandler(os.Stdout, &hOptions),
		),
	)

}

func Info(format string, args ...interface{}) {
	Get().Info(format, args...)
}

func Error(format string, args ...interface{}) {
	Get().Error(format, args...)
}

func Debug(format string, args ...interface{}) {
	Get().Debug(format, args...)
}

func Warn(format string, args ...interface{}) {
	Get().Warn(format, args...)
}

// Get return pointer to already created logger via constructor.
func Get() *slog.Logger {
	once.Do(func() {
		config := C.GetConfig()
		createLogger(config.App.LogLevel)
	})
	return logger
}
