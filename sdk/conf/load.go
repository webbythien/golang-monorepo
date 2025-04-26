package conf

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func Load[Config any](defaultConfig []byte) (*Config, error) {
	var cfg Config

	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to read viper config: %w", err)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
