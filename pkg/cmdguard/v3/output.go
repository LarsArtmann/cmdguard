package v3

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	output "github.com/larsartmann/go-output"
)

// OutputFormat is a type-safe output format enum, aliased from go-output.
// Use output.FormatTable, output.FormatJSON, etc. for format values.
type OutputFormat = output.Format

// OutputConfig holds the output formatting configuration.
type OutputConfig struct {
	Format OutputFormat
	Writer io.Writer
}

// DefaultOutputConfig returns the default output configuration (table format, stdout).
func DefaultOutputConfig() OutputConfig {
	return OutputConfig{
		Format: output.FormatTable,
		Writer: os.Stdout,
	}
}

// OutputResult renders data in the configured output format.
// For *Table, delegates to go-output's RenderTable registry.
// For arbitrary data, delegates to go-output's RenderUnknown registry (JSON, YAML, TOML).
// Returns shape-aware errors when a format does not support the provided data type.
func OutputResult(cfg OutputConfig, data any) error {
	if cfg.Writer == nil {
		cfg.Writer = os.Stdout
	}

	opts := output.RenderOptions{Writer: cfg.Writer}

	if td, ok := data.(*output.Table); ok {
		err := output.RenderTable(td, cfg.Format, opts)
		if _, unsupported := errors.AsType[*output.UnsupportedFormatError](err); unsupported {
			return fmt.Errorf("%w: %s (format supports %s, not table data)",
				ErrUnsupportedFormat, cfg.Format, formatShapes(cfg.Format))
		}

		return err
	}

	err := output.RenderUnknown(data, cfg.Format, opts)
	if _, unsupported := errors.AsType[*output.UnsupportedFormatError](err); unsupported {
		return fmt.Errorf("%w: %s (format does not support arbitrary data)",
			ErrFormatRequiresTypedData, cfg.Format)
	}

	return err
}

// OutputTable is a convenience function to output table data with headers and rows.
// Uses AddRowChecked for fail-fast row validation.
func OutputTable(format OutputFormat, headers []string, rows [][]string) error {
	data := output.NewTable(headers)

	for _, row := range rows {
		err := data.AddRowChecked(row)
		if err != nil {
			return err
		}
	}

	return OutputResult(OutputConfig{Format: format}, data)
}

// RegisteredFormats returns all output formats with registered Table marshalers.
// Use this to dynamically discover available formats based on imported sub-modules.
func RegisteredFormats() []OutputFormat {
	return output.RegisteredTableMarshalFormats()
}

// formatShapes returns a human-readable description of what shapes a format supports.
func formatShapes(f output.Format) string {
	shapes := f.Shapes()
	if len(shapes) > 0 {
		names := make([]string, len(shapes))
		for i, s := range shapes {
			names[i] = string(s)
		}

		return strings.Join(names, ", ")
	}

	return "unknown"
}
