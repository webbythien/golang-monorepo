package must

import (
	"github.com/webbythien/monorepo/pkg/l"
	"github.com/webbythien/monorepo/sdk/conf"
)

func LoadConfig[Config any](defaultConfig []byte) *Config {
	var ll = l.New()
	cfg, err := conf.Load[Config](defaultConfig)
	if err != nil {
		ll.Fatal("Failed to load config", l.Error(err))
	}
	return cfg
}
