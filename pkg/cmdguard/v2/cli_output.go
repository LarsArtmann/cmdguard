package v2

import (
	"fmt"
	"strings"

	output "github.com/larsartmann/go-output"
	"github.com/spf13/cobra"
)

// WithOutputFormat adds a global --output flag for format selection.
// When used, the CLI gets a persistent --output flag (short -o) that
// controls the output format. Access the resolved format via cli.OutputFormat().
//
// Usage:
//
//	cli, _ := v2.NewCLI[Config]("app", "My app", Config{},
//	    v2.WithOutputFormat[Config](output.FormatTable),
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
		return output.FormatTable
	}

	return cli.outputFormat
}

// SetOutputFormat sets the output format at runtime.
func (cli *CLI[T]) SetOutputFormat(format OutputFormat) {
	cli.outputFormat = format
}

// initOutputFlag sets up the --output flag with dynamic help from registered formats.
func (cli *CLI[T]) initOutputFlag() {
	if cli.outputFormat == "" {
		return
	}

	formats := output.RegisteredTableDataFormats()

	names := make([]string, len(formats), len(formats))
	for i, f := range formats {
		names[i] = string(f)
	}

	help := fmt.Sprintf("Output format (%s)", strings.Join(names, ", "))
	cli.AddGlobalFlag("output", "o", string(cli.outputFormat), help)
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

	format, err := output.ParseFormat(formatStr)
	if err != nil {
		return fmt.Errorf("%w: %q is not a valid output format", ErrUnsupportedFormat, formatStr)
	}

	cli.outputFormat = format

	return nil
}
