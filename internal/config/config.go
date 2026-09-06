package config

import (
	"fmt"
	"os"
	"path/filepath"

	yaml "go.yaml.in/yaml/v3"
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
	merged := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(DefaultConfig), &merged); err != nil {
		return nil, "", fmt.Errorf("failed to load default config: %w", err)
	}

	if path != "" {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil, "", fmt.Errorf("failed to resolve config path: %w", err)
		}

		data, err := os.ReadFile(absPath)
		if err != nil {
			return nil, "", fmt.Errorf("config file not found: %s", absPath)
		}

		var override map[string]interface{}
		if err := yaml.Unmarshal(data, &override); err != nil {
			return nil, "", fmt.Errorf("failed to parse config: %w", err)
		}
		merged = deepMergeMaps(merged, override)
		path = absPath
	}

	// Round-trip through YAML rather than a generic decoder (e.g.
	// mapstructure) so struct fields resolve via their `yaml` tags — the
	// same tags the merged data was itself produced from.
	mergedYAML, err := yaml.Marshal(merged)
	if err != nil {
		return nil, "", fmt.Errorf("failed to render merged config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(mergedYAML, &cfg); err != nil {
		return nil, "", fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &cfg, path, nil
}

// deepMergeMaps merges override on top of base and returns the result,
// recursing into nested maps and replacing everything else (scalars,
// slices) wholesale with override's value.
//
// This intentionally avoids viper's Load/MergeInConfig, which normalizes
// every map key to lowercase (via its case-insensitive internal store) even
// for values the user cares about verbatim — silently turning
// `shell.environment.EDITOR` into `editor`, which real shells never read.
// Keys here are kept byte-for-byte as written.
func deepMergeMaps(base, override map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(base))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range override {
		if baseVal, ok := out[k]; ok {
			baseMap, baseIsMap := baseVal.(map[string]interface{})
			overrideMap, overrideIsMap := v.(map[string]interface{})
			if baseIsMap && overrideIsMap {
				out[k] = deepMergeMaps(baseMap, overrideMap)
				continue
			}
		}
		out[k] = v
	}
	return out
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
