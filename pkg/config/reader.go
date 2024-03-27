package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func FromFile(configPath string) *Config {
	file, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("error opening config file %q", err))
	}

	var c Config
	if err := json.Unmarshal(file, &c); err != nil {
		panic(fmt.Sprintf("error parsing config %q", err))
	}

	return &c
}
