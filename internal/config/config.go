package config

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/pelletier/go-toml"
)

type Config struct {
	APIKey string `toml:"api_key"`
}

var configFile string

func init() {
	configFile = filepath.Join(xdg.DataHome, "gopud", "config.toml")
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(configFile)
	if os.IsNotExist(err) {
		return &Config{}, nil
	} else if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func SaveConfig(cfg *Config) error {
	if err := os.MkdirAll(filepath.Dir(configFile), 0700); err != nil {
		return err
	}

	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, data, 0600)
}

func DeleteConfig() error {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(configFile)
}
