package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"log/slog"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
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
			slog.With("error", err).Error(fmt.Sprintf("Failed to load config file %s", configName))
			panic(err)
		}
	}

	if err := K.Load(env.Provider(envVarPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, envVarPrefix)), "_", ".", -1)
	}), nil); err != nil {
		slog.With("error", err).Error("Failed to load config from environment")
		panic(err)
	}
}
