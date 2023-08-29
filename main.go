package main

import (
	"github.com/merlinfuchs/blimp/internal"
	"github.com/merlinfuchs/blimp/internal/config"
	"github.com/rs/zerolog/log"
)

func main() {
	config.InitConfig()

	err := internal.AppEntry()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start app")
	}
}
