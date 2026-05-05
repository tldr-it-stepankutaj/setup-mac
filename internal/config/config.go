package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// autoDiscoverNames lists filenames searched in ~/.config/setup-mac/ when no
// --config is passed. `config.yaml` is the canonical name written by
// `setup-mac config init`; `default.yaml` is also accepted because the dist
// tarball ships its bundled config under that name and users have copied it
// in by hand.
var autoDiscoverNames = []string{"config.yaml", "default.yaml"}

// Load loads configuration from file or uses defaults. If configPath is empty,
// it auto-discovers a config file in ~/.config/setup-mac/. The returned path
// is the absolute path of the file actually loaded, or "" when only the
// embedded defaults were used.
func Load(configPath string) (*Config, string, error) {
	resolved := configPath
	if resolved == "" {
		resolved = autoDiscoverConfig()
	}
	return loadFromPath(resolved)
}

// LoadDefault loads ONLY the embedded default config — used by tests and any
// caller that explicitly does not want auto-discovery to leak the user's
// real config into the result.
func LoadDefault() (*Config, error) {
	cfg, _, err := loadFromPath("")
	return cfg, err
}

func loadFromPath(path string) (*Config, string, error) {
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))
	v.SetConfigType("yaml")

	if err := v.ReadConfig(bytes.NewBufferString(DefaultConfig)); err != nil {
		return nil, "", fmt.Errorf("failed to load default config: %w", err)
	}

	if path != "" {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil, "", fmt.Errorf("failed to resolve config path: %w", err)
		}
		if _, err := os.Stat(absPath); err != nil {
			return nil, "", fmt.Errorf("config file not found: %s", absPath)
		}

		v.SetConfigFile(absPath)
		if err := v.MergeInConfig(); err != nil {
			return nil, "", fmt.Errorf("failed to merge config: %w", err)
		}
		path = absPath
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, "", fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &cfg, path, nil
}

// autoDiscoverConfig returns the first existing config file under
// ~/.config/setup-mac/, or "" if none.
func autoDiscoverConfig() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	dir := filepath.Join(homeDir, ".config", "setup-mac")
	for _, name := range autoDiscoverNames {
		p := filepath.Join(dir, name)
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}
