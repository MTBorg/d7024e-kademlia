package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger(level string) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(zerolog.InfoLevel) //Default to info level
	logLevel, err := zerolog.ParseLevel(level)
	if err == nil {
		zerolog.SetGlobalLevel(logLevel)
		log.Info().Str("Level", level).Msg("Log level set")
	} else {
		log.Error().Str("Level", level).Msg("Failed to parse log level, defaulting to info level...")
	}
}
