package v2

import (
	"fmt"
	"io"
	"os"

	output "github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/d2"
	"github.com/larsartmann/go-output/delimited"
	"github.com/larsartmann/go-output/graph"
	"github.com/larsartmann/go-output/markup"
	"github.com/larsartmann/go-output/plantuml"
	"github.com/larsartmann/go-output/serialization"
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
func OutputResult(cfg OutputConfig, data any) error {
	if cfg.Writer == nil {
		cfg.Writer = os.Stdout
	}

	strategy, ok := formatRegistry[cfg.Format]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnsupportedFormat, cfg.Format)
	}

	return strategy.Render(cfg.Writer, data)
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

// FormatStrategy renders data to a writer in a specific output format.
type FormatStrategy interface {
	Render(w io.Writer, data any) error
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

// tableRenderStrategy renders TableData to a writer via a string-producing function.
// Returns ErrFormatRequiresTypedData for non-TableData input.
type tableRenderStrategy struct {
	label  string
	render func(*output.TableData) (string, error)
}

func (s *tableRenderStrategy) Render(w io.Writer, data any) error {
	td := unwrapTableData(data)
	if td == nil {
		return fmt.Errorf("%w: %s", ErrFormatRequiresTypedData, s.label)
	}

	result, err := s.render(td)
	if err != nil {
		return fmt.Errorf("rendering %s: %w", s.label, err)
	}

	fmt.Fprintln(w, result)

	return nil
}

// marshalStrategy renders any data (including TableData) to a writer via a marshal function.
type marshalStrategy struct {
	label   string
	marshal func(any) ([]byte, error)
}

func (s *marshalStrategy) Render(w io.Writer, data any) error {
	result, err := s.marshal(data)
	if err != nil {
		return fmt.Errorf("marshaling %s: %w", s.label, err)
	}

	fmt.Fprintln(w, string(result))

	return nil
}

// dualStrategy delegates to different strategies for TableData vs arbitrary data.
type dualStrategy struct {
	table FormatStrategy
	any   FormatStrategy
}

func (s *dualStrategy) Render(w io.Writer, data any) error {
	if unwrapTableData(data) != nil {
		return s.table.Render(w, data)
	}

	return s.any.Render(w, data)
}

// styledTableStrategy renders TableData as a styled terminal table via go-output/table.
type styledTableStrategy struct{}

func (s *styledTableStrategy) Render(w io.Writer, data any) error {
	td := unwrapTableData(data)
	if td == nil {
		return fmt.Errorf("%w: table", ErrFormatRequiresTypedData)
	}

	return renderTableStyled(w, td)
}

// csvStrategy renders TableData as CSV via streaming writer.
type csvStrategy struct{}

func (s *csvStrategy) Render(w io.Writer, data any) error {
	td := unwrapTableData(data)
	if td == nil {
		return fmt.Errorf("%w: csv", ErrFormatRequiresTypedData)
	}

	return renderTableCSV(w, td)
}

// formatRegistry maps OutputFormat to FormatStrategy implementations.
var formatRegistry = map[OutputFormat]FormatStrategy{
	output.FormatTable: &styledTableStrategy{},
	output.FormatCSV:   &csvStrategy{},
	output.FormatJSON: &dualStrategy{
		table: &tableRenderStrategy{label: "JSON", render: func(d *output.TableData) (string, error) {
			b, err := serialization.MarshalJSON(d)

			return string(b), err
		}},
		any: &marshalStrategy{label: "JSON", marshal: func(v any) ([]byte, error) {
			return output.MarshalJSONIndent(v, "", "  ")
		}},
	},
	output.FormatTSV: &tableRenderStrategy{label: "TSV", render: func(d *output.TableData) (string, error) {
		b, err := delimited.MarshalTSV(d)

		return string(b), err
	}},
	output.FormatYAML: &marshalStrategy{label: "YAML", marshal: serialization.MarshalYAML},
	output.FormatXML: &tableRenderStrategy{label: "XML", render: func(d *output.TableData) (string, error) {
		b, err := markup.MarshalXMLFromTableData(d)

		return string(b), err
	}},
	output.FormatMarkdown: &tableRenderStrategy{label: "markdown", render: func(d *output.TableData) (string, error) {
		return output.NewMarkdownTableFromData(d).Render()
	}},
	output.FormatHTML: &tableRenderStrategy{label: "HTML", render: func(d *output.TableData) (string, error) {
		r := markup.NewHTMLRenderer()
		r.SetData(d)

		return r.Render()
	}},
	output.FormatTree: &tableRenderStrategy{label: "tree", render: func(d *output.TableData) (string, error) {
		return output.TreeRendererFromTableData(d).Render()
	}},
	output.FormatD2: &tableRenderStrategy{label: "D2", render: func(d *output.TableData) (string, error) {
		return d2.D2FromTableData(d).Render()
	}},
	output.FormatMermaid: &tableRenderStrategy{label: "Mermaid", render: func(d *output.TableData) (string, error) {
		return graph.MermaidFromTableData(d).Render()
	}},
	output.FormatDOT: &tableRenderStrategy{label: "DOT", render: func(d *output.TableData) (string, error) {
		return graph.DOTFromTableData(d).Render()
	}},
	output.FormatJSONL: &tableRenderStrategy{label: "JSONL", render: func(d *output.TableData) (string, error) {
		b, err := serialization.MarshalJSONLFromTableData(d)

		return string(b), err
	}},
	output.FormatAsciiDoc: &tableRenderStrategy{label: "AsciiDoc", render: func(d *output.TableData) (string, error) {
		b, err := markup.MarshalAsciiDocFromTableData(d)

		return string(b), err
	}},
	output.FormatTOML: &dualStrategy{
		table: &tableRenderStrategy{label: "TOML", render: func(d *output.TableData) (string, error) {
			b, err := serialization.MarshalTOMLFromTableData(d)

			return string(b), err
		}},
		any: &marshalStrategy{label: "TOML", marshal: serialization.MarshalTOML},
	},
	output.FormatPlantUML: &tableRenderStrategy{label: "PlantUML", render: func(d *output.TableData) (string, error) {
		return plantuml.PlantUMLFromTableData(d).Render()
	}},
}

func renderTableStyled(w io.Writer, data *output.TableData) error {
	t := table.New()
	t.SetHeaders(data.GetHeaders()...)

	for _, row := range data.GetRows() {
		t.AddRow(row...)
	}

	result, err := t.Render()
	if err != nil {
		return fmt.Errorf("rendering table: %w", err)
	}

	fmt.Fprintln(w, result)

	return nil
}

func renderTableCSV(w io.Writer, data *output.TableData) error {
	cw := delimited.NewCSVWriter(w)

	if err := cw.WriteHeader(data.GetHeaders()); err != nil {
		return fmt.Errorf("writing CSV header: %w", err)
	}

	for _, row := range data.GetRows() {
		if err := cw.WriteRow(row); err != nil {
			return fmt.Errorf("writing CSV row: %w", err)
		}
	}

	cw.Flush()

	return cw.Error()
}
