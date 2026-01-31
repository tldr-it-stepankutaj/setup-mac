package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tldr-it-stepankutaj/setup-mac/internal/config"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration file",
	Long: `Validate a configuration file without running any installations.

Examples:
  # Validate the default embedded config
  setup-mac validate

  # Validate a custom config file
  setup-mac validate --config my-config.yaml`,
	RunE: runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

// ValidationResult contains the result of config validation
type ValidationResult struct {
	Valid    bool
	Errors   []string
	Warnings []string
}

func runValidate(cmd *cobra.Command, args []string) error {
	if cfgFile != "" {
		fmt.Printf("Validating config: %s\n", cfgFile)
	} else {
		fmt.Println("Validating embedded default config")
	}
	fmt.Println()

	// Check if config file exists (for custom configs)
	if cfgFile != "" {
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			color.New(color.FgRed).Printf("✗ Config file not found: %s\n", cfgFile)
			return fmt.Errorf("validation failed")
		}
	}

	// Try to load the config
	cfg, err := config.Load(cfgFile)
	if err != nil {
		color.New(color.FgRed).Printf("✗ Configuration invalid: %v\n", err)
		return fmt.Errorf("validation failed")
	}

	// Perform detailed validation
	result := validateConfig(cfg)

	// Print results
	printValidationResult(result, cfg)

	if !result.Valid {
		return fmt.Errorf("validation failed with %d error(s)", len(result.Errors))
	}

	return nil
}

func validateConfig(cfg *config.Config) ValidationResult {
	result := ValidationResult{Valid: true}

	// Validate Homebrew config
	if cfg.Homebrew.Install {
		// Check for duplicate formulae
		seen := make(map[string]bool)
		for _, formula := range cfg.Homebrew.Formulae {
			if seen[formula] {
				result.Warnings = append(result.Warnings, fmt.Sprintf("Duplicate formula: %s", formula))
			}
			seen[formula] = true
		}

		// Check for duplicate casks
		seen = make(map[string]bool)
		for _, cask := range cfg.Homebrew.Casks {
			if seen[cask] {
				result.Warnings = append(result.Warnings, fmt.Sprintf("Duplicate cask: %s", cask))
			}
			seen[cask] = true
		}

		// Check for duplicate taps
		seen = make(map[string]bool)
		for _, tap := range cfg.Homebrew.Taps {
			if seen[tap] {
				result.Warnings = append(result.Warnings, fmt.Sprintf("Duplicate tap: %s", tap))
			}
			seen[tap] = true
		}
	}

	// Validate Git config
	if cfg.Git.Configure {
		if cfg.Git.User.Name == "" {
			result.Warnings = append(result.Warnings, "Git user.name is not set")
		}
		if cfg.Git.User.Email == "" {
			result.Warnings = append(result.Warnings, "Git user.email is not set")
		}
		if cfg.Git.User.Email != "" && !strings.Contains(cfg.Git.User.Email, "@") {
			result.Errors = append(result.Errors, fmt.Sprintf("Git user.email appears invalid: %s", cfg.Git.User.Email))
			result.Valid = false
		}
	}

	// Validate SSH config
	if cfg.SSH.GenerateKey {
		validKeyTypes := map[string]bool{"ed25519": true, "rsa": true, "ecdsa": true, "": true}
		if !validKeyTypes[cfg.SSH.KeyType] {
			result.Errors = append(result.Errors, fmt.Sprintf("Invalid SSH key type: %s (valid: ed25519, rsa, ecdsa)", cfg.SSH.KeyType))
			result.Valid = false
		}
	}

	// Validate Shell config
	for name, cmd := range cfg.Shell.Aliases {
		if cmd == "" {
			result.Errors = append(result.Errors, fmt.Sprintf("Empty alias command for: %s", name))
			result.Valid = false
		}
	}

	return result
}

func printValidationResult(result ValidationResult, cfg *config.Config) {
	// Print config summary
	color.New(color.FgCyan, color.Bold).Println("Configuration Summary")
	fmt.Println("──────────────────────────────────────")

	// Homebrew
	if cfg.Homebrew.Install {
		fmt.Printf("  Homebrew:     %s (%d formulae, %d casks, %d taps)\n",
			color.GreenString("enabled"),
			len(cfg.Homebrew.Formulae),
			len(cfg.Homebrew.Casks),
			len(cfg.Homebrew.Taps))
	} else {
		fmt.Printf("  Homebrew:     %s\n", color.YellowString("disabled"))
	}

	// Oh My Zsh
	if cfg.Terminal.OhMyZsh.Install {
		fmt.Printf("  Oh My Zsh:    %s (%d plugins)\n",
			color.GreenString("enabled"),
			len(cfg.Terminal.OhMyZsh.Plugins))
	} else {
		fmt.Printf("  Oh My Zsh:    %s\n", color.YellowString("disabled"))
	}

	// Powerlevel10k
	if cfg.Terminal.Powerlevel10k.Install {
		fmt.Printf("  Powerlevel10k: %s\n", color.GreenString("enabled"))
	} else {
		fmt.Printf("  Powerlevel10k: %s\n", color.YellowString("disabled"))
	}

	// Git
	if cfg.Git.Configure {
		fmt.Printf("  Git:          %s (user: %s)\n",
			color.GreenString("enabled"),
			cfg.Git.User.Name)
	} else {
		fmt.Printf("  Git:          %s\n", color.YellowString("disabled"))
	}

	// SSH
	if cfg.SSH.GenerateKey {
		keyType := cfg.SSH.KeyType
		if keyType == "" {
			keyType = "ed25519"
		}
		fmt.Printf("  SSH:          %s (type: %s)\n",
			color.GreenString("enabled"),
			keyType)
	} else {
		fmt.Printf("  SSH:          %s\n", color.YellowString("disabled"))
	}

	// Shell
	fmt.Printf("  Shell:        %d aliases, %d env vars\n",
		len(cfg.Shell.Aliases),
		len(cfg.Shell.Environment))

	// macOS
	if cfg.MacOS.Configure {
		fmt.Printf("  macOS:        %s\n", color.GreenString("enabled"))
	} else {
		fmt.Printf("  macOS:        %s\n", color.YellowString("disabled"))
	}

	fmt.Println()

	// Print errors
	if len(result.Errors) > 0 {
		color.New(color.FgRed, color.Bold).Println("Errors")
		fmt.Println("──────────────────────────────────────")
		for _, err := range result.Errors {
			color.New(color.FgRed).Printf("  ✗ %s\n", err)
		}
		fmt.Println()
	}

	// Print warnings
	if len(result.Warnings) > 0 {
		color.New(color.FgYellow, color.Bold).Println("Warnings")
		fmt.Println("──────────────────────────────────────")
		for _, warn := range result.Warnings {
			color.New(color.FgYellow).Printf("  ⚠ %s\n", warn)
		}
		fmt.Println()
	}

	// Print final status
	if result.Valid {
		color.New(color.FgGreen, color.Bold).Println("✓ Configuration is valid")
	} else {
		color.New(color.FgRed, color.Bold).Println("✗ Configuration has errors")
	}
}
