// Package internal provides shared utilities for cmdguard examples.
package internal

import (
	"context"
	"fmt"
	"os"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// CLIExecutor is an interface for CLI types that can be executed.
// This allows examples to share the execute function.
type CLIExecutor interface {
	Execute(ctx context.Context) error
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

// NewSimpleCommand creates a simple leaf command that prints a message.
func NewSimpleCommand[Config any](
	name, message, short string,
) (v2.Command[Config, v2.NoFlags], error) {
	return v2.NewCommand[Config, v2.NoFlags](name,
		func(_ context.Context, _ *Config, _ v2.NoFlags) error {
			fmt.Println(message)

			return nil
		},
		v2.WithShort[Config, v2.NoFlags](short),
	)
}
