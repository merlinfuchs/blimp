package config

import (
	_ "embed"
	"log/slog"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/rawbytes"
)

//go:embed default.config.toml
var defaultConfigTomlBytes []byte

func setupDefaults() {
	if err := K.Load(rawbytes.Provider(defaultConfigTomlBytes), toml.Parser()); err != nil {
		slog.With("error", err).Error("Failed to load default config")
		panic(err)
	}
}
