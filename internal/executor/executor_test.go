package executor

import (
	"bytes"
	"context"
	"testing"
)

func TestExecutorDryRun(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exec := New(true, false)
	exec.Stdout = &stdout
	exec.Stderr = &stderr

	ctx := context.Background()
	result, err := exec.Run(ctx, "echo", "hello")

	if err != nil {
		t.Fatalf("dry-run should not return error: %v", err)
	}

	if !result.DryRun {
		t.Error("expected result.DryRun to be true")
	}

	if result.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", result.ExitCode)
	}
}

func TestExecutorRun(t *testing.T) {
	exec := New(false, false)
	ctx := context.Background()

	result, err := exec.Run(ctx, "echo", "hello")
	if err != nil {
		t.Fatalf("failed to run command: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", result.ExitCode)
	}

	expected := "hello\n"
	if result.Stdout != expected {
		t.Errorf("expected stdout %q, got %q", expected, result.Stdout)
	}
}

func TestExecutorRunShell(t *testing.T) {
	exec := New(false, false)
	ctx := context.Background()

	result, err := exec.RunShell(ctx, "echo hello && echo world")
	if err != nil {
		t.Fatalf("failed to run shell command: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", result.ExitCode)
	}

	expected := "hello\nworld\n"
	if result.Stdout != expected {
		t.Errorf("expected stdout %q, got %q", expected, result.Stdout)
	}
}

func TestExecutorExists(t *testing.T) {
	exec := New(false, false)

	// 'echo' should exist on all systems
	if !exec.Exists("echo") {
		t.Error("expected 'echo' command to exist")
	}

	// Nonexistent command
	if exec.Exists("nonexistent-command-12345") {
		t.Error("expected nonexistent command to not exist")
	}
}

func TestExecutorFailingCommand(t *testing.T) {
	exec := New(false, false)
	ctx := context.Background()

	result, err := exec.Run(ctx, "sh", "-c", "exit 1")

	if err == nil {
		t.Error("expected error for failing command")
	}

	if result.ExitCode != 1 {
		t.Errorf("expected exit code 1, got %d", result.ExitCode)
	}
}
