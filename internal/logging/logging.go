package logging

import (
	"io"
	"os"

	"github.com/merlinfuchs/blimp/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger() {
	logContext := log.With()

	logWriters := make([]io.Writer, 0)
	logWriters = append(logWriters, &lumberjack.Logger{
		Filename:   config.K.String("logging.filename"),
		MaxSize:    config.K.Int("logging.max_size"),
		MaxAge:     config.K.Int("logging.max_age"),
		MaxBackups: config.K.Int("logging.max_backups"),
	})
	if config.K.Bool("debug") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		logWriters = append(logWriters, zerolog.ConsoleWriter{Out: os.Stdout})
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	writer := io.MultiWriter(logWriters...)
	log.Logger = logContext.Logger().Output(writer)
}
