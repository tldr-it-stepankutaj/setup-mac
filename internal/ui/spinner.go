package ui

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

// Spinner wraps a spinner for showing progress
type Spinner struct {
	s       *spinner.Spinner
	message string
	output  io.Writer
	enabled bool
}

// NewSpinner creates a new spinner
func NewSpinner(message string) *Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + message
	s.Writer = os.Stdout

	return &Spinner{
		s:       s,
		message: message,
		output:  os.Stdout,
		enabled: true,
	}
}

// Start starts the spinner
func (sp *Spinner) Start() {
	if sp.enabled {
		sp.s.Start()
	}
}

// Stop stops the spinner
func (sp *Spinner) Stop() {
	if sp.enabled {
		sp.s.Stop()
	}
}

// Success stops the spinner and shows success message
func (sp *Spinner) Success(msg string) {
	sp.Stop()
	if msg == "" {
		msg = sp.message
	}
	color.New(color.FgGreen).Fprint(sp.output, "✓ ")
	fmt.Fprintln(sp.output, msg)
}

// Fail stops the spinner and shows failure message
func (sp *Spinner) Fail(msg string) {
	sp.Stop()
	if msg == "" {
		msg = sp.message
	}
	color.New(color.FgRed).Fprint(sp.output, "✗ ")
	fmt.Fprintln(sp.output, msg)
}

// Info shows an info message
func (sp *Spinner) Info(msg string) {
	sp.Stop()
	color.New(color.FgCyan).Fprint(sp.output, "ℹ ")
	fmt.Fprintln(sp.output, msg)
}

// Warning shows a warning message
func (sp *Spinner) Warning(msg string) {
	sp.Stop()
	color.New(color.FgYellow).Fprint(sp.output, "⚠ ")
	fmt.Fprintln(sp.output, msg)
}

// UpdateMessage updates the spinner message
func (sp *Spinner) UpdateMessage(msg string) {
	sp.message = msg
	sp.s.Suffix = " " + msg
}

// SetEnabled enables or disables the spinner
func (sp *Spinner) SetEnabled(enabled bool) {
	sp.enabled = enabled
}

// PrintSuccess prints a success message
func PrintSuccess(msg string) {
	color.New(color.FgGreen).Print("✓ ")
	fmt.Println(msg)
}

// PrintError prints an error message
func PrintError(msg string) {
	color.New(color.FgRed).Print("✗ ")
	fmt.Println(msg)
}

// PrintInfo prints an info message
func PrintInfo(msg string) {
	color.New(color.FgCyan).Print("ℹ ")
	fmt.Println(msg)
}

// PrintWarning prints a warning message
func PrintWarning(msg string) {
	color.New(color.FgYellow).Print("⚠ ")
	fmt.Println(msg)
}

// PrintHeader prints a section header
func PrintHeader(msg string) {
	fmt.Println()
	color.New(color.FgMagenta, color.Bold).Println("═══════════════════════════════════════")
	color.New(color.FgMagenta, color.Bold).Printf("  %s\n", msg)
	color.New(color.FgMagenta, color.Bold).Println("═══════════════════════════════════════")
	fmt.Println()
}

// PrintHeaderWithProgress prints a section header with progress indicator [current/total]
func PrintHeaderWithProgress(msg string, current, total int) {
	fmt.Println()
	color.New(color.FgMagenta, color.Bold).Println("═══════════════════════════════════════")
	color.New(color.FgCyan, color.Bold).Printf("  [%d/%d] ", current, total)
	color.New(color.FgMagenta, color.Bold).Printf("%s\n", msg)
	color.New(color.FgMagenta, color.Bold).Println("═══════════════════════════════════════")
	fmt.Println()
}

// PrintStep prints a step message
func PrintStep(msg string) {
	color.New(color.FgBlue).Print("→ ")
	fmt.Println(msg)
}

// PrintDryRun prints a dry-run message
func PrintDryRun(msg string) {
	color.New(color.FgYellow).Print("[DRY-RUN] ")
	fmt.Println(msg)
}
