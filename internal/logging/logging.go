package logging

import (
	"log/slog"
	"os"

	"github.com/merlinfuchs/blimp/internal/config"
)

func InitLogger() {
	logLevel := slog.LevelInfo
	if config.K.Bool("debug") {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))

	slog.SetDefault(logger)
}
