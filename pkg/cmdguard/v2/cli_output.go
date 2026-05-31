package v2

import (
	"fmt"

	"github.com/spf13/cobra"
)

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
	}
}

// OutputFormat returns the resolved output format from the --output flag.
// If WithOutputFormat was not used, returns FormatTable.
func (cli *CLI[T]) OutputFormat() OutputFormat {
	if cli.outputFormat == "" {
		return FormatTable
	}

	return cli.outputFormat
}

// SetOutputFormat sets the output format at runtime.
func (cli *CLI[T]) SetOutputFormat(format OutputFormat) {
	cli.outputFormat = format
}

// initOutputFlag sets up the --output flag and hooks into flag parsing.
func (cli *CLI[T]) initOutputFlag() {
	if cli.outputFormat == "" {
		return
	}

	cli.AddGlobalFlag("output", "o", string(cli.outputFormat),
		"Output format (table, json, csv, yaml, markdown, xml)")
}

// parseOutputFlag resolves the --output flag value after cobra parses flags.
func (cli *CLI[T]) parseOutputFlag(c *cobra.Command) error {
	if cli.outputFormat == "" {
		return nil
	}

	formatStr, _ := c.Flags().GetString("output")

	if formatStr == "" {
		return nil
	}

	format, err := ParseOutputFormat(formatStr)
	if err != nil {
		return fmt.Errorf("%w: %q is not a valid output format", ErrUnsupportedFormat, formatStr)
	}

	cli.outputFormat = format

	return nil
}
