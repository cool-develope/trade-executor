package config

import (
	"path/filepath"
	"strings"

	"github.com/cool-develope/trade-executor/internal/server"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// ExchangeConfig struct
type ExchangeConfig struct {
	Name    string
	Symbols []string
}

// Config struct
type Config struct {
	Exchange ExchangeConfig
	Server   server.Config
}

// Load loads the configuration.
func Load(configFilePath string) (*Config, error) {
	var cfg Config

	path, fullFile := filepath.Split(configFilePath)

	file := strings.Split(fullFile, ".")

	viper.AddConfigPath(path)
	viper.SetConfigName(file[0])
	viper.SetConfigType(file[1])

	err := viper.Unmarshal(&cfg, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc()))
	return &cfg, err
}
