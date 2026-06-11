package v2

import (
	"errors"
	"fmt"
	"io"
	"os"

	output "github.com/larsartmann/go-output"
	_ "github.com/larsartmann/go-output/d2"
	_ "github.com/larsartmann/go-output/delimited"
	_ "github.com/larsartmann/go-output/graph"
	_ "github.com/larsartmann/go-output/markup"
	_ "github.com/larsartmann/go-output/plantuml"
	_ "github.com/larsartmann/go-output/serialization"
	"github.com/larsartmann/go-output/table"
)

// OutputFormat is a type-safe output format enum.
// Supported formats: table, json, csv, tsv, markdown, xml, d2, yaml, html, tree, mermaid, dot, jsonl, asciidoc, toml, plantuml.
type OutputFormat = output.Format

// Output format constants re-exported from go-output for convenience.
var (
	FormatTable    = output.FormatTable
	FormatJSON     = output.FormatJSON
	FormatCSV      = output.FormatCSV
	FormatTSV      = output.FormatTSV
	FormatMarkdown = output.FormatMarkdown
	FormatXML      = output.FormatXML
	FormatD2       = output.FormatD2
	FormatYAML     = output.FormatYAML
	FormatHTML     = output.FormatHTML
	FormatTree     = output.FormatTree
	FormatMermaid  = output.FormatMermaid
	FormatDOT      = output.FormatDOT
	FormatJSONL    = output.FormatJSONL
	FormatAsciiDoc = output.FormatAsciiDoc
	FormatTOML     = output.FormatTOML
	FormatPlantUML = output.FormatPlantUML
)

// ParseOutputFormat parses a string into an OutputFormat.
func ParseOutputFormat(s string) (OutputFormat, error) {
	f, err := output.ParseFormat(s)
	if err != nil {
		return f, fmt.Errorf("parsing output format %q: %w", s, err)
	}

	return f, nil
}

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
// For TableData, delegates to go-output's RenderTableData registry.
// For arbitrary data, delegates to go-output's RenderAnyData registry (JSON, YAML, TOML).
func OutputResult(cfg OutputConfig, data any) error {
	if cfg.Writer == nil {
		cfg.Writer = os.Stdout
	}

	opts := output.RenderOptions{Writer: cfg.Writer}

	if td := unwrapTableData(data); td != nil {
		err := output.RenderTableData(td, cfg.Format, opts)

		var unsupported *output.UnsupportedFormatError
		if errors.As(err, &unsupported) {
			return fmt.Errorf("%w: %s", ErrUnsupportedFormat, cfg.Format)
		}

		return err
	}

	err := output.RenderAnyData(data, cfg.Format, opts)

	var unsupported *output.UnsupportedFormatError
	if errors.As(err, &unsupported) {
		return fmt.Errorf("%w: %s", ErrFormatRequiresTypedData, cfg.Format)
	}

	return err
}

// unwrapTableData extracts *output.TableData from any, handling both pointer and value types.
// Returns nil if data is not a TableData.
func unwrapTableData(data any) *output.TableData {
	switch d := data.(type) {
	case *output.TableData:
		return d
	case output.TableData:
		return &d
	default:
		return nil
	}
}

// OutputTable is a convenience function to output table data with headers and rows.
func OutputTable(format OutputFormat, headers []string, rows [][]string) error {
	data := output.NewTableData(headers)

	for _, row := range rows {
		data.AddRow(row)
	}

	return OutputResult(OutputConfig{Format: format}, data)
}

// OutputStyledTable renders a styled terminal table using lipgloss.
//
// Deprecated: use OutputResult(OutputConfig{Format: FormatTable}, data) instead.
func OutputStyledTable(headers []string, rows [][]string) error {
	t := table.New()
	t.SetHeaders(headers...)

	for _, row := range rows {
		t.AddRow(row...)
	}

	result, err := t.Render()
	if err != nil {
		return fmt.Errorf("rendering styled table: %w", err)
	}

	fmt.Fprintln(os.Stdout, result)

	return nil
}

// SupportedFormats returns all output formats supported by the current configuration.
func SupportedFormats() []OutputFormat {
	return output.AllFormats
}

// IsFormatSupported returns true if the format is a valid, registered output format.
func IsFormatSupported(f OutputFormat) bool {
	return f.IsValid()
}
