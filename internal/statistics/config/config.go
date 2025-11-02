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
	// Intended to be a production config setup using an env variable
	apiKeyFromEnv := os.Getenv("ALPHAVANTAGE_API_KEY")
	if apiKeyFromEnv != "" {
		return &APIKeysConfig{AlphaVantageAPIKey: apiKeyFromEnv}, nil
	}
	// If no env var use the toml config file, intended to be for development
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found at path: %s", path)
	}

	var config APIKeysConfig
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, fmt.Errorf("error decoding TOML file: %w", err)
	}

	if config.AlphaVantageAPIKey == "" {
		errorMsg := "alphavantage API key is missing. \n"
		errorMsg += "You can provide it either via the ALPHAVANTAGE_API_KEY environment variable (recommended for production),\n"
		errorMsg += fmt.Sprintf("or by adding it to your config file at '%s' like this:\n\n", path)
		errorMsg += `alphavantage = "YOUR_API_KEY_HERE"`
		return nil, fmt.Errorf("%s", errorMsg)
	}

	return &config, nil
}
