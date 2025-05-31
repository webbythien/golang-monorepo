package config

import (
	_ "embed"

	"github.com/monorepo/pkg/cache"
	"github.com/monorepo/sdk/conf"
	"github.com/monorepo/sdk/must"
)

//go:embed config.yaml
var defaultConfig []byte

type Config struct {
	conf.Base  `mapstructure:",squash"`
	PostgreSQL *conf.PostgreSQL   `yaml:"postgresql" mapstructure:"postgresql"`
	Redis      *cache.RedisConfig `yaml:"redis" mapstructure:"redis"`
}

func Load() *Config {
	return must.LoadConfig[Config](defaultConfig)
}
