package config

import (
	_ "embed"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/rs/zerolog/log"
)

//go:embed default.config.toml
var defaultConfigTomlBytes []byte

func setupDefaults() {
	if err := K.Load(rawbytes.Provider(defaultConfigTomlBytes), toml.Parser()); err != nil {
		log.Panic().Err(err).Msgf("Failed to load config file %s", configName)
	}
}
