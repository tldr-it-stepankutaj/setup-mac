package executor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Executor handles command execution
type Executor struct {
	DryRun  bool
	Verbose bool
	Stdout  io.Writer
	Stderr  io.Writer
}

// New creates a new Executor
func New(dryRun, verbose bool) *Executor {
	return &Executor{
		DryRun:  dryRun,
		Verbose: verbose,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}
}

// Result contains command execution result
type Result struct {
	Command  string
	ExitCode int
	Stdout   string
	Stderr   string
	Duration time.Duration
	DryRun   bool
}

// Run executes a command and returns the result
func (e *Executor) Run(ctx context.Context, name string, args ...string) (*Result, error) {
	cmdStr := formatCommand(name, args)
	startTime := time.Now()

	if e.DryRun {
		color.New(color.FgYellow).Fprintf(e.Stdout, "[DRY-RUN] %s\n", cmdStr)
		return &Result{
			Command:  cmdStr,
			ExitCode: 0,
			DryRun:   true,
			Duration: time.Since(startTime),
		}, nil
	}

	if e.Verbose {
		color.New(color.FgCyan).Fprintf(e.Stdout, "[EXEC] %s\n", cmdStr)
	}

	cmd := exec.CommandContext(ctx, name, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	duration := time.Since(startTime)

	result := &Result{
		Command:  cmdStr,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Duration: duration,
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
		return result, fmt.Errorf("command failed: %w", err)
	}

	result.ExitCode = 0
	return result, nil
}

// RunShell executes a shell command
func (e *Executor) RunShell(ctx context.Context, command string) (*Result, error) {
	return e.Run(ctx, "sh", "-c", command)
}

// RunInteractive executes a command with interactive I/O
func (e *Executor) RunInteractive(ctx context.Context, name string, args ...string) error {
	cmdStr := formatCommand(name, args)

	if e.DryRun {
		color.New(color.FgYellow).Fprintf(e.Stdout, "[DRY-RUN] %s\n", cmdStr)
		return nil
	}

	if e.Verbose {
		color.New(color.FgCyan).Fprintf(e.Stdout, "[EXEC] %s\n", cmdStr)
	}

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Exists checks if a command exists
func (e *Executor) Exists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// Which returns the path to a command
func (e *Executor) Which(name string) (string, error) {
	return exec.LookPath(name)
}

func formatCommand(name string, args []string) string {
	if len(args) == 0 {
		return name
	}

	// Quote arguments with spaces
	quotedArgs := make([]string, len(args))
	for i, arg := range args {
		if strings.Contains(arg, " ") || strings.Contains(arg, "\"") {
			quotedArgs[i] = fmt.Sprintf("%q", arg)
		} else {
			quotedArgs[i] = arg
		}
	}

	return fmt.Sprintf("%s %s", name, strings.Join(quotedArgs, " "))
}
