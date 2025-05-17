package must

import (
	"github.com/monorepo/pkg/l"
	"github.com/monorepo/sdk/conf"
)

func LoadConfig[Config any](defaultConfig []byte) *Config {
	ll := l.New()
	cfg, err := conf.Load[Config](defaultConfig)
	if err != nil {
		ll.Fatal("Failed to load config", l.Error(err))
	}
	return cfg
}
