package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AppConfig struct {
	APIAddress string
	DnsDBPath  string
	Secret     string
}

var Config *AppConfig

func LoadConfig(configFilePath string) (*AppConfig, error) {
	viper.SetConfigFile(configFilePath)
	viper.SetDefault("APIAddress", "127.0.0.1:7777")
	viper.SetDefault("DnsDBPath", "./db/dns")
	viper.SetDefault("Secret", "!%Secret")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config:   %w", err)
	}
	Config = &config
	return &config, nil
}
