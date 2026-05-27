package logger

import (
	"log/slog"
	"os"
)

// New creates the process logger.
// For now logs go to stdout only; rotation, archives, and uploads belong to the
// later observability step.
func New(env string) *slog.Logger {
	options := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if env == "local" || env == "development" {
		return slog.New(slog.NewTextHandler(os.Stdout, options))
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, options))
}
