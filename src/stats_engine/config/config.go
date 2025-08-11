package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type APIKeysConfig struct {
	AlphaVantageAPIKey string `toml:"alphavantage"`
}

func LoadConfig(path string) (*APIKeysConfig, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found at path: %s", path)
	}

	var config APIKeysConfig
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, fmt.Errorf("error decoding TOML file: %w", err)
	}

	if config.AlphaVantageAPIKey == "" {
		return nil, fmt.Errorf("alphavantage API key is missing from config file")
	}

	fmt.Printf("Successfully loaded configuration from: %s\n", path)
	return &config, nil
}
