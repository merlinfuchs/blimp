package config

import (
	"errors"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
)

var K = koanf.New(".")

const configName = "blimp.toml"
const envVarPrefix = "BLIMP__"

var CfgFile string

func fileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

func InitConfig() {
	setupDefaults()

	if err := K.Load(file.Provider(configName), toml.Parser()); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Panic().Err(err).Msgf("Failed to load config file %s", configName)
		}
	}

	if err := K.Load(env.Provider(envVarPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, envVarPrefix)), "_", ".", -1)
	}), nil); err != nil {
		log.Panic().Err(err).Msgf("Failed to load env vars")
	}
}
