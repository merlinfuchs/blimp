package main

import (
	"fmt"

	"github.com/merlinfuchs/blimp/internal"
	"github.com/merlinfuchs/blimp/internal/config"
	"github.com/merlinfuchs/blimp/internal/logging"
	"github.com/rs/zerolog/log"
)

var outDatedConfigKeys = []string{
	"layout",
	"views",
}

func main() {
	config.InitConfig()
	logging.InitLogger()

	for _, key := range outDatedConfigKeys {
		if config.K.Get(key) != nil {
			fmt.Printf("Config key '%s' is no longer supported, please take a look at https://github.com/merlinfuchs/blimp#readme to see how to configure blimp!", key)
			return
		}
	}

	err := internal.AppEntry()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start app")
	}
}
