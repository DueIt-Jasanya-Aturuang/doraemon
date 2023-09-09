package config

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LogInit() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logFile, err := os.Create("app.log")
	if err != nil {
		panic(err)
	}

	multi := zerolog.MultiLevelWriter(logFile, os.Stdout)

	log.Logger = zerolog.New(multi).With().Caller().Timestamp().Logger()

	log.Info().Msg("Logger initialization successfully")
}
