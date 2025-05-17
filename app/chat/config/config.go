package config

import (
	_ "embed"

	"github.com/monorepo/sdk/conf"
	"github.com/monorepo/sdk/must"
)

//go:embed config.yaml
var defaultConfig []byte

type Config struct {
	conf.Base  `mapstructure:",squash"`
	PostgreSQL *conf.PostgreSQL `yaml:"postgresql" mapstructure:"postgresql"`
}

func Load() *Config {
	return must.LoadConfig[Config](defaultConfig)
}
