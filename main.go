package main

import (
	"fmt"

	"log/slog"

	"github.com/merlinfuchs/blimp/internal"
	"github.com/merlinfuchs/blimp/internal/config"
	"github.com/merlinfuchs/blimp/internal/logging"
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
		slog.With("error", err).Error("Failed to start blimp")
	}
}
