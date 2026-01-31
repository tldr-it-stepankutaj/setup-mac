package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tldr-it-stepankutaj/setup-mac/internal/config"
	"github.com/tldr-it-stepankutaj/setup-mac/internal/installer"
)

var jsonOutput bool

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show installation status of all components",
	Long: `Show which components are installed and which are not.

Examples:
  # Show status of all components
  setup-mac status

  # Output as JSON (for scripting)
  setup-mac status --json`,
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().BoolVar(&jsonOutput, "json", false, "output as JSON")
}

// ComponentStatus represents the status of a single component
type ComponentStatus struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Installed   bool   `json:"installed"`
}

// SystemStatus represents the overall system status
type SystemStatus struct {
	System     SystemInfo        `json:"system"`
	Components []ComponentStatus `json:"components"`
}

// SystemInfo contains system information
type SystemInfo struct {
	OS           string `json:"os"`
	Arch         string `json:"arch"`
	AppleSilicon bool   `json:"apple_silicon"`
	MacOSVersion string `json:"macos_version,omitempty"`
}

func runStatus(cmd *cobra.Command, args []string) error {
	// Load configuration (to get proper context)
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create installer context (dry-run doesn't matter for status)
	ictx := installer.NewContext(cfg, false, verbose)
	ctx := context.Background()

	// Get system info
	sysInfo := getSystemInfo(ctx, ictx)

	// Check all installers
	installers := []installer.Installer{
		installer.NewXcodeInstaller(ictx),
		installer.NewRosettaInstaller(ictx),
		installer.NewHomebrewInstaller(ictx),
		installer.NewOhMyZshInstaller(ictx),
		installer.NewPowerlevel10kInstaller(ictx),
		installer.NewShellInstaller(ictx),
		installer.NewMacOSInstaller(ictx),
		installer.NewGitInstaller(ictx),
		installer.NewSSHInstaller(ictx),
	}

	var components []ComponentStatus
	for _, inst := range installers {
		status := ComponentStatus{
			Name:        inst.Name(),
			Description: inst.Description(),
			Installed:   inst.IsInstalled(ctx),
		}
		components = append(components, status)
	}

	// Build full status
	status := SystemStatus{
		System:     sysInfo,
		Components: components,
	}

	if jsonOutput {
		return outputJSON(status)
	}

	return outputHuman(status)
}

func getSystemInfo(ctx context.Context, ictx *installer.Context) SystemInfo {
	info := SystemInfo{
		OS:           runtime.GOOS,
		Arch:         runtime.GOARCH,
		AppleSilicon: runtime.GOARCH == "arm64",
	}

	// Get macOS version
	result, err := ictx.Executor.Run(ctx, "sw_vers", "-productVersion")
	if err == nil {
		info.MacOSVersion = strings.TrimSpace(result.Stdout)
	}

	return info
}

func outputJSON(status SystemStatus) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(status)
}

func outputHuman(status SystemStatus) error {
	// Print header
	color.New(color.FgCyan, color.Bold).Println("System Information")
	fmt.Println("──────────────────────────────────────")
	fmt.Printf("  OS:            %s\n", status.System.OS)
	fmt.Printf("  Architecture:  %s\n", status.System.Arch)
	if status.System.AppleSilicon {
		fmt.Printf("  Apple Silicon: %s\n", color.GreenString("Yes"))
	}
	if status.System.MacOSVersion != "" {
		fmt.Printf("  macOS Version: %s", status.System.MacOSVersion)
	}
	fmt.Println()

	// Print components
	color.New(color.FgCyan, color.Bold).Println("Components")
	fmt.Println("──────────────────────────────────────")

	installed := 0
	for _, comp := range status.Components {
		var statusIcon string
		var statusColor *color.Color
		if comp.Installed {
			statusIcon = "✓"
			statusColor = color.New(color.FgGreen)
			installed++
		} else {
			statusIcon = "✗"
			statusColor = color.New(color.FgRed)
		}

		statusColor.Printf("  %s ", statusIcon)
		fmt.Printf("%-20s %s\n", comp.Name, color.New(color.Faint).Sprint(comp.Description))
	}

	fmt.Println()
	fmt.Printf("  %d/%d components installed\n", installed, len(status.Components))

	return nil
}
