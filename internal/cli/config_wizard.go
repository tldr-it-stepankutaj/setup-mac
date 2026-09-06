package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	yaml "go.yaml.in/yaml/v3"

	"github.com/tldr-it-stepankutaj/setup-mac/internal/config"
	"github.com/tldr-it-stepankutaj/setup-mac/internal/ui"
)

var (
	configWizardForce bool
	configWizardPath  string
)

var configWizardCmd = &cobra.Command{
	Use:   "wizard",
	Short: "Interactively build a config.yaml",
	Long: `Walk through an interactive wizard to build a tailored configuration
instead of hand-editing the default YAML.

The wizard controls which components are enabled (Homebrew, terminal setup,
shell aliases, macOS defaults, Git, SSH) and, for Homebrew, which packages to
install. Finer-grained settings (individual macOS tweaks, Git/shell aliases)
are left at their sane defaults — edit the resulting YAML by hand for those.

Examples:
  # Run the wizard, writing to the default location
  setup-mac config wizard

  # Overwrite an existing config
  setup-mac config wizard --force

  # Write to a custom path
  setup-mac config wizard --output ./my-setup.yaml`,
	// User-input errors (e.g. destination exists) shouldn't dump full usage —
	// the warning above the error already explains what to do.
	SilenceUsage: true,
	RunE:         runConfigWizard,
}

func init() {
	configCmd.AddCommand(configWizardCmd)

	configWizardCmd.Flags().BoolVarP(&configWizardForce, "force", "f", false, "overwrite existing file")
	configWizardCmd.Flags().StringVarP(&configWizardPath, "output", "o", "", "destination path (default ~/.config/setup-mac/config.yaml)")
}

func runConfigWizard(cmd *cobra.Command, args []string) error {
	dest, err := resolveConfigDest(configWizardPath)
	if err != nil {
		return err
	}

	if _, err := os.Stat(dest); err == nil && !configWizardForce {
		ui.PrintWarning(fmt.Sprintf("%s already exists — pass --force to overwrite", dest))
		return fmt.Errorf("destination exists")
	}

	cfg, err := config.LoadDefault()
	if err != nil {
		return fmt.Errorf("failed to load defaults: %w", err)
	}

	color.New(color.FgCyan, color.Bold).Println("setup-mac configuration wizard")
	fmt.Println("Press Enter to accept the default shown in brackets.")

	// This command only makes sense run interactively at a real terminal, so
	// the prompt is always interactive regardless of any config setting.
	prompt := ui.NewPrompt(true)

	ui.PrintHeader("General")
	if cfg.Settings.Interactive, err = prompt.Confirm("Ask for confirmation before installing?", cfg.Settings.Interactive); err != nil {
		return err
	}
	if cfg.Settings.BackupDotfiles, err = prompt.Confirm("Back up existing dotfiles before modifying them?", cfg.Settings.BackupDotfiles); err != nil {
		return err
	}

	ui.PrintHeader("Homebrew")
	if cfg.Homebrew.Install, err = prompt.Confirm("Install Homebrew and packages?", cfg.Homebrew.Install); err != nil {
		return err
	}
	if cfg.Homebrew.Install {
		useDefaults, err := prompt.Confirm(
			fmt.Sprintf("Use the default package list (%d formulae, %d casks, %d taps)?",
				len(cfg.Homebrew.Formulae), len(cfg.Homebrew.Casks), len(cfg.Homebrew.Taps)),
			true)
		if err != nil {
			return err
		}
		if !useDefaults {
			if cfg.Homebrew.Formulae, err = promptStringList(prompt, "Formulae (comma-separated)", cfg.Homebrew.Formulae); err != nil {
				return err
			}
			if cfg.Homebrew.Casks, err = promptStringList(prompt, "Casks (comma-separated)", cfg.Homebrew.Casks); err != nil {
				return err
			}
			if cfg.Homebrew.Taps, err = promptStringList(prompt, "Taps (comma-separated)", cfg.Homebrew.Taps); err != nil {
				return err
			}
		}
	}

	ui.PrintHeader("Terminal")
	installTerminal, err := prompt.Confirm("Install Oh My Zsh + Powerlevel10k?", cfg.Terminal.OhMyZsh.Install)
	if err != nil {
		return err
	}
	cfg.Terminal.OhMyZsh.Install = installTerminal
	cfg.Terminal.Powerlevel10k.Install = installTerminal

	ui.PrintHeader("Shell")
	configureShell, err := prompt.Confirm("Configure shell aliases & environment variables (~/.zshrc)?", true)
	if err != nil {
		return err
	}
	if !configureShell {
		cfg.Shell.Aliases = nil
		cfg.Shell.Environment = nil
		cfg.Shell.ZshrcExtras = nil
	}

	ui.PrintHeader("macOS defaults")
	if cfg.MacOS.Configure, err = prompt.Confirm("Apply opinionated macOS defaults (Dock, Finder, Keyboard)?", cfg.MacOS.Configure); err != nil {
		return err
	}

	ui.PrintHeader("Git")
	if cfg.Git.Configure, err = prompt.Confirm("Configure Git (identity, aliases, settings)?", cfg.Git.Configure); err != nil {
		return err
	}
	if cfg.Git.Configure {
		name := cfg.Git.User.Name
		if name == "" {
			name = existingGitConfig("user.name")
		}
		if cfg.Git.User.Name, err = prompt.Input("Git user name", name); err != nil {
			return err
		}

		email := cfg.Git.User.Email
		if email == "" {
			email = existingGitConfig("user.email")
		}
		if cfg.Git.User.Email, err = prompt.Input("Git user email", email); err != nil {
			return err
		}
	}

	ui.PrintHeader("SSH")
	if cfg.SSH.GenerateKey, err = prompt.Confirm("Generate an SSH key if one doesn't already exist?", cfg.SSH.GenerateKey); err != nil {
		return err
	}
	if cfg.SSH.GenerateKey {
		_, keyType, err := prompt.Select("SSH key type", []string{"ed25519", "rsa", "ecdsa"})
		if err != nil {
			return err
		}
		cfg.SSH.KeyType = keyType
		cfg.SSH.KeyFile = fmt.Sprintf("~/.ssh/id_%s", keyType)

		if cfg.SSH.Comment, err = prompt.Input("SSH key comment", cfg.Git.User.Email); err != nil {
			return err
		}
	}

	out, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to render config: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	if err := os.WriteFile(dest, out, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Println()
	color.New(color.FgGreen).Printf("✓ ")
	fmt.Printf("Wrote config to %s\n", dest)
	fmt.Println()
	fmt.Println("Review the file, then run:")
	fmt.Printf("  setup-mac install --all --config %s\n", dest)
	return nil
}

// promptStringList asks for a comma-separated list, pre-filled with the
// current values, and returns the trimmed, non-empty entries.
func promptStringList(prompt *ui.Prompt, label string, current []string) ([]string, error) {
	input, err := prompt.Input(label, strings.Join(current, ", "))
	if err != nil {
		return nil, err
	}
	var out []string
	for _, item := range strings.Split(input, ",") {
		if item = strings.TrimSpace(item); item != "" {
			out = append(out, item)
		}
	}
	return out, nil
}

// existingGitConfig returns a global git config value, or "" if unset or git
// isn't installed. Mirrors the pre-fill behavior the Git installer uses at
// install time, so the wizard doesn't ask for information already on disk.
func existingGitConfig(key string) string {
	out, err := exec.Command("git", "config", "--global", "--get", key).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
