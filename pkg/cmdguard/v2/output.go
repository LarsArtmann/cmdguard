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

	switch d := data.(type) {
	case *output.TableData:
		return renderTableData(cfg, d)
	case output.TableData:
		return renderTableData(cfg, &d)
	default:
		return renderAny(cfg, data)
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

// tableRenderer renders a TableData to a writer.
type tableRenderer func(w io.Writer, data *output.TableData) error

// renderMarshalFunc produces a string from any data (e.g., JSON, YAML, XML, TSV).
type renderMarshalFunc func(data any) ([]byte, error)

// renderStringFunc produces a rendered string from TableData (e.g., via .Render()).
type renderStringFunc func(data *output.TableData) (string, error)

// renderAndWrite calls fn to produce a string, then writes it to w.
func renderAndWrite(w io.Writer, label string, data *output.TableData, fn renderStringFunc) error {
	result, err := fn(data)
	if err != nil {
		return fmt.Errorf("rendering %s: %w", label, err)
	}

	fmt.Fprintln(w, result)

	return nil
}

// marshalAndWrite calls fn to produce bytes from data, then writes them to w.
func marshalAndWrite(w io.Writer, label string, data any, fn renderMarshalFunc) error {
	result, err := fn(data)
	if err != nil {
		return fmt.Errorf("marshaling %s: %w", label, err)
	}

	fmt.Fprintln(w, string(result))

	return nil
}

// tableFormatRegistry maps OutputFormat to rendering functions.
var tableFormatRegistry = map[OutputFormat]tableRenderer{
	output.FormatTable: renderTableStyled,
	output.FormatCSV:   renderTableCSV,
	output.FormatJSON: func(w io.Writer, data *output.TableData) error {
		return marshalAndWrite(w, "JSON", data, serialization.MarshalJSON)
	},
	output.FormatTSV: func(w io.Writer, data *output.TableData) error {
		return marshalAndWrite(w, "TSV", data, delimited.MarshalTSV)
	},
	output.FormatYAML: func(w io.Writer, data *output.TableData) error {
		return marshalAndWrite(w, "YAML", data, serialization.MarshalYAML)
	},
	output.FormatXML: func(w io.Writer, data *output.TableData) error {
		return renderAndWrite(w, "XML", data, func(d *output.TableData) (string, error) {
			b, err := markup.MarshalXMLFromTableData(d)

			return string(b), err
		})
	},
	output.FormatMarkdown: func(w io.Writer, data *output.TableData) error {
		return renderAndWrite(w, "markdown", data, func(d *output.TableData) (string, error) {
			return output.NewMarkdownTableFromData(d).Render()
		})
	},
	output.FormatHTML: func(w io.Writer, data *output.TableData) error {
		return renderAndWrite(w, "HTML", data, func(d *output.TableData) (string, error) {
			r := markup.NewHTMLRenderer()
			r.SetData(d)

			return r.Render()
		})
	},
	output.FormatTree: func(w io.Writer, data *output.TableData) error {
		return renderAndWrite(w, "tree", data, func(d *output.TableData) (string, error) {
			return output.TreeRendererFromTableData(d).Render()
		})
	},
	output.FormatD2: func(w io.Writer, data *output.TableData) error {
		return renderAndWrite(w, "D2", data, func(d *output.TableData) (string, error) {
			return d2.D2FromTableData(d).Render()
		})
	},
	output.FormatMermaid: func(w io.Writer, data *output.TableData) error {
		return renderAndWrite(w, "Mermaid", data, func(d *output.TableData) (string, error) {
			return graph.MermaidFromTableData(d).Render()
		})
	},
	output.FormatDOT: func(w io.Writer, data *output.TableData) error {
		return renderAndWrite(w, "DOT", data, func(d *output.TableData) (string, error) {
			return graph.DOTFromTableData(d).Render()
		})
	},
	output.FormatJSONL: func(w io.Writer, data *output.TableData) error {
		return renderAndWrite(w, "JSONL", data, func(d *output.TableData) (string, error) {
			b, err := serialization.MarshalJSONLFromTableData(d)

			return string(b), err
		})
	},
	output.FormatAsciiDoc: func(w io.Writer, data *output.TableData) error {
		return renderAndWrite(w, "AsciiDoc", data, func(d *output.TableData) (string, error) {
			b, err := markup.MarshalAsciiDocFromTableData(d)

			return string(b), err
		})
	},
	output.FormatTOML: func(w io.Writer, data *output.TableData) error {
		return renderAndWrite(w, "TOML", data, func(d *output.TableData) (string, error) {
			b, err := serialization.MarshalTOMLFromTableData(d)

			return string(b), err
		})
	},
	output.FormatPlantUML: func(w io.Writer, data *output.TableData) error {
		return renderAndWrite(w, "PlantUML", data, func(d *output.TableData) (string, error) {
			return plantuml.PlantUMLFromTableData(d).Render()
		})
	},
}

// anyRenderer renders arbitrary data to a writer.
type anyRenderer func(w io.Writer, data any) error

// anyFormatRegistry maps OutputFormat to generic rendering functions.
var anyFormatRegistry = map[OutputFormat]anyRenderer{
	output.FormatJSON: func(w io.Writer, data any) error {
		return marshalAndWrite(w, "JSON", data, func(v any) ([]byte, error) {
			return output.MarshalJSONIndent(v, "", "  ")
		})
	},
	output.FormatYAML: func(w io.Writer, data any) error {
		return marshalAndWrite(w, "YAML", data, serialization.MarshalYAML)
	},
	output.FormatTOML: func(w io.Writer, data any) error {
		return marshalAndWrite(w, "TOML", data, serialization.MarshalTOML)
	},
}

func renderTableData(cfg OutputConfig, data *output.TableData) error {
	renderer, ok := tableFormatRegistry[cfg.Format]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnsupportedFormat, cfg.Format)
	}

	return renderer(cfg.Writer, data)
}

func renderAny(cfg OutputConfig, data any) error {
	renderer, ok := anyFormatRegistry[cfg.Format]
	if !ok {
		return fmt.Errorf("%w: %s", ErrFormatRequiresTypedData, cfg.Format)
	}

	return renderer(cfg.Writer, data)
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
