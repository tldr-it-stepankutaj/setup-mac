package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Load loads configuration from file or uses defaults
func Load(configPath string) (*Config, error) {
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))
	v.SetConfigType("yaml")

	// Load defaults first
	if err := v.ReadConfig(bytes.NewBufferString(DefaultConfig)); err != nil {
		return nil, fmt.Errorf("failed to load default config: %w", err)
	}

	// If custom config provided, merge it
	if configPath != "" {
		absPath, err := filepath.Abs(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve config path: %w", err)
		}

		if _, err := os.Stat(absPath); err != nil {
			return nil, fmt.Errorf("config file not found: %s", absPath)
		}

		v.SetConfigFile(absPath)
		if err := v.MergeInConfig(); err != nil {
			return nil, fmt.Errorf("failed to merge config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// LoadDefault loads the default configuration
func LoadDefault() (*Config, error) {
	return Load("")
}
