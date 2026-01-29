package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefault(t *testing.T) {
	cfg, err := LoadDefault()
	if err != nil {
		t.Fatalf("failed to load default config: %v", err)
	}

	if cfg.Version != "1.0" {
		t.Errorf("expected version 1.0, got %s", cfg.Version)
	}

	if !cfg.Homebrew.Install {
		t.Error("expected homebrew.install to be true")
	}

	if len(cfg.Homebrew.Formulae) == 0 {
		t.Error("expected homebrew.formulae to have items")
	}

	if !cfg.Terminal.OhMyZsh.Install {
		t.Error("expected terminal.oh_my_zsh.install to be true")
	}

	if !cfg.Terminal.Powerlevel10k.Install {
		t.Error("expected terminal.powerlevel10k.install to be true")
	}

	if len(cfg.Shell.Aliases) == 0 {
		t.Error("expected shell.aliases to have items")
	}

	if !cfg.MacOS.Configure {
		t.Error("expected macos.configure to be true")
	}

	if !cfg.Git.Configure {
		t.Error("expected git.configure to be true")
	}

	if !cfg.SSH.GenerateKey {
		t.Error("expected ssh.generate_key to be true")
	}
}

func TestLoadCustomConfig(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	configContent := `
version: "2.0"
settings:
  dry_run: true
  interactive: false
homebrew:
  install: false
  formulae:
    - custom-formula
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("failed to load custom config: %v", err)
	}

	if cfg.Version != "2.0" {
		t.Errorf("expected version 2.0, got %s", cfg.Version)
	}

	if !cfg.Settings.DryRun {
		t.Error("expected settings.dry_run to be true")
	}

	if cfg.Settings.Interactive {
		t.Error("expected settings.interactive to be false")
	}

	if cfg.Homebrew.Install {
		t.Error("expected homebrew.install to be false")
	}

	if len(cfg.Homebrew.Formulae) != 1 || cfg.Homebrew.Formulae[0] != "custom-formula" {
		t.Errorf("expected formulae [custom-formula], got %v", cfg.Homebrew.Formulae)
	}
}

func TestLoadNonExistentConfig(t *testing.T) {
	_, err := Load("/nonexistent/config.yaml")
	if err == nil {
		t.Error("expected error for non-existent config")
	}
}
