package v2

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"
)

// outputState holds the resolved output format for the CLI.
type outputState struct {
	mu     sync.Mutex
	format OutputFormat
}

// WithOutputFormat adds a global --output flag for format selection.
// When used, the CLI gets a persistent --output flag (short -o) that
// controls the output format. Access the resolved format via cli.OutputFormat().
//
// Usage:
//
//	cli, _ := v2.NewCLI[Config]("app", "My app", Config{},
//	    v2.WithOutputFormat[Config](v2.FormatTable),
//	)
//	// CLI now has --output/-o flag
//	// In handlers: format := cli.OutputFormat()
func WithOutputFormat[T any](defaultFormat OutputFormat) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.outputFormat = defaultFormat
		cli.outputEnabled = true
	}
}

// OutputFormat returns the resolved output format from the --output flag.
// If WithOutputFormat was not used, returns FormatTable.
func (cli *CLI[T]) OutputFormat() OutputFormat {
	if cli.outputState == nil {
		return FormatTable
	}

	cli.outputState.mu.Lock()
	defer cli.outputState.mu.Unlock()

	return cli.outputState.format
}

// SetOutputFormat sets the output format at runtime.
func (cli *CLI[T]) SetOutputFormat(format OutputFormat) {
	if cli.outputState == nil {
		cli.outputState = &outputState{format: format}

		return
	}

	cli.outputState.mu.Lock()
	cli.outputState.format = format
	cli.outputState.mu.Unlock()
}

// initOutputFlag sets up the --output flag and hooks into flag parsing.
func (cli *CLI[T]) initOutputFlag() {
	if !cli.outputEnabled {
		return
	}

	cli.outputState = &outputState{format: cli.outputFormat}

	cli.AddGlobalFlag("output", "o", string(cli.outputFormat),
		"Output format (table, json, csv, yaml, markdown, xml)")
}

// parseOutputFlag resolves the --output flag value after cobra parses flags.
func (cli *CLI[T]) parseOutputFlag(c *cobra.Command) error {
	if !cli.outputEnabled {
		return nil
	}

	formatStr, err := c.Flags().GetString("output")
	if err != nil {
		return nil // flag not found is ok
	}

	if formatStr == "" {
		return nil
	}

	format, err := ParseOutputFormat(formatStr)
	if err != nil {
		return fmt.Errorf("invalid output format %q: %w", formatStr, err)
	}

	cli.outputState.mu.Lock()
	cli.outputState.format = format
	cli.outputState.mu.Unlock()

	return nil
}
