// Package internal provides shared utilities for cmdguard examples.
package internal

import (
	"context"
	"fmt"
	"os"
)

// CLIExecutor is an interface for CLI types that can be executed.
// This allows examples to share the execute function.
type CLIExecutor interface {
	Execute(context.Context) error
}

// Fatalf prints the formatted error to stderr and exits with code 1.
func Fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

// Execute runs the CLI and exits on error.
func Execute(ctx context.Context, cli CLIExecutor) {
	err := cli.Execute(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
