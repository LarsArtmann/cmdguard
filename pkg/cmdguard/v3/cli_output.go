package v3

import (
	"fmt"
	"strings"

	output "github.com/larsartmann/go-output"
	"github.com/spf13/cobra"
)

// OutputFormat returns the resolved output format from the --output flag.
// If WithOutputFormat was not used, returns FormatTable.
func (cli *CLI[T]) OutputFormat() OutputFormat {
	if cli.spec.outputFormat == "" {
		return output.FormatTable
	}

	return cli.spec.outputFormat
}

// SetOutputFormat sets the output format at runtime.
func (cli *CLI[T]) SetOutputFormat(format OutputFormat) {
	cli.spec.outputFormat = format
}

// initOutputFlag sets up the --output flag with dynamic help from registered formats.
func (cli *CLI[T]) initOutputFlag() {
	if cli.spec.outputFormat == "" {
		return
	}

	formats := output.RegisteredTableMarshalFormats()

	names := make([]string, 0, len(formats))
	for _, f := range formats {
		names = append(names, string(f))
	}

	help := fmt.Sprintf("Output format (%s)", strings.Join(names, ", "))
	cli.AddGlobalFlag("output", "o", string(cli.spec.outputFormat), help)
}

// parseOutputFlag resolves the --output flag value after cobra parses flags.
func (cli *CLI[T]) parseOutputFlag(c *cobra.Command) error {
	if cli.spec.outputFormat == "" {
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

	cli.spec.outputFormat = format

	return nil
}
