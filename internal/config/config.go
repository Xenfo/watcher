package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Package struct {
	CurrentVersion string `json:"currentVersion"`
	TargetVersion  string `json:"targetVersion"`
	Betas          bool   `json:"includeBetas"`
	Notes          string `json:"notes"`
	Notify         bool   `json:"notify"`
}

type Config struct {
	Packages map[string]Package `json:"packages"`
}

func (c *Config) GetPackage(name string) Package {
	return c.Packages[name]
}

func Create() (*Config, error) {
	contents, err := os.ReadFile("config.json")
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %s", err)
	}

	var config *Config
	err = json.Unmarshal(contents, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %s", err)
	}

	return config, nil
}
