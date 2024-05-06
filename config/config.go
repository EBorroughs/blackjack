package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	StorageBackend        string `json:"storageBackend"`
	RequestTimeoutSeconds int    `json:"requestTimeoutSeconds"`
	SessionMaxAgeSeconds  int    `json:"sessionMaxAgeSeconds"`
}

func LoadConfig() (Config, error) {
	var config Config
	configFile, err := os.Open("config.json")
	if err != nil {
		return config, fmt.Errorf("unable to open config file: %w", err)
	}
	defer configFile.Close()

	parser := json.NewDecoder(configFile)
	err = parser.Decode(&config)
	if err != nil {
		return config, fmt.Errorf("failed to decode config file: %w", err)
	}

	return config, nil
}
