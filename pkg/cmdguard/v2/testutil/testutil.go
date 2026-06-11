// Package testutil provides testing utilities for cmdguard v2 consumers.
package testutil

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	v2 "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
)

// TestResult holds the outcome of executing a CLI command in tests.
type TestResult struct {
	// Stdout contains everything written to the command's output stream.
	Stdout string
	// Stderr contains everything written to the command's error stream.
	Stderr string
	// Error is the error returned by ExecuteWithArgs (nil on success).
	Error error
}

// ExitCode returns the exit code if the error implements v2.ExitCoder.
// Returns 0 if there was no error or the error does not carry an exit code.
func (r *TestResult) ExitCode() int {
	if r.Error == nil {
		return 0
	}

	if exitCoder, ok := errors.AsType[v2.ExitCoder](r.Error); ok {
		return exitCoder.ExitCode()
	}

	return 1 // default non-zero exit code for errors
}

// TestCLI wraps a cmdguard CLI for testing with captured output.
type TestCLI[T any] struct {
	cli    *v2.CLI[T]
	stdout *bytes.Buffer
	stderr *bytes.Buffer
}

// NewTestCLI creates a test wrapper around an existing CLI.
// It captures stdout and stderr automatically.
func NewTestCLI[T any](cli *v2.CLI[T]) *TestCLI[T] {
	return &TestCLI[T]{
		cli:    cli,
		stdout: new(bytes.Buffer),
		stderr: new(bytes.Buffer),
	}
}

// CLI returns the underlying cmdguard CLI.
func (tc *TestCLI[T]) CLI() *v2.CLI[T] {
	return tc.cli
}

// Stdout returns the current contents of the captured stdout buffer.
func (tc *TestCLI[T]) Stdout() string {
	return tc.stdout.String()
}

// Stderr returns the current contents of the captured stderr buffer.
func (tc *TestCLI[T]) Stderr() string {
	return tc.stderr.String()
}

// ExecuteWithArgs runs the CLI with the given arguments and captures output.
// The result contains stdout, stderr, and any error returned.
func (tc *TestCLI[T]) ExecuteWithArgs(ctx context.Context, args []string) *TestResult {
	// Reset buffers for a clean capture.
	tc.stdout.Reset()
	tc.stderr.Reset()

	// Wire captured output into the root cobra command.
	root := tc.cli.RootCommand()
	root.SetOut(tc.stdout)
	root.SetErr(tc.stderr)

	err := tc.cli.ExecuteWithArgs(ctx, args)

	return &TestResult{
		Stdout: tc.stdout.String(),
		Stderr: tc.stderr.String(),
		Error:  err,
	}
}

// AddCommand registers a Command on a CLI and fails the test on error.
// Centralizes the AddCommand fatal pattern used across test files.
func AddCommand[T, F any](t *testing.T, cli *v2.CLI[T], cmd v2.Command[T, F]) {
	t.Helper()

	if err := v2.AddCommand(cli, cmd); err != nil {
		t.Fatalf("failed to add command: %v", err)
	}
}

// AssertOutputContains fails the test if output does not contain substring.
// Use for captured CLI command output buffers.
func AssertOutputContains(t *testing.T, output, substr string) {
	t.Helper()

	if !strings.Contains(output, substr) {
		t.Errorf("output should contain %q, got: %s", substr, output)
	}
}
