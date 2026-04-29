package v2

import (
	"fmt"
	"io"
	"os"

	output "github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/table"
)

// OutputFormat is a type-safe output format enum.
// Supported formats: table, json, csv, tsv, markdown, xml, d2, yaml, html, tree, mermaid, dot.
type OutputFormat = output.Format

// Output format constants.
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
)

// ParseOutputFormat parses a string into an OutputFormat.
func ParseOutputFormat(s string) (OutputFormat, error) {
	return output.ParseFormat(s)
}

// OutputConfig holds the output formatting configuration.
type OutputConfig struct {
	Format OutputFormat
	Writer io.Writer
}

// DefaultOutputConfig returns the default output configuration (table format, stdout).
func DefaultOutputConfig() OutputConfig {
	return OutputConfig{
		Format: FormatTable,
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
		return outputTable(cfg, d)
	case output.TableData:
		return outputTable(cfg, &d)
	default:
		return outputAny(cfg, data)
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

func outputTable(cfg OutputConfig, data *output.TableData) error {
	switch cfg.Format {
	case FormatTable:
		t := table.New()
		t.SetHeaders(data.GetHeaders()...)
		for _, row := range data.GetRows() {
			t.AddRow(row...)
		}
		result, err := t.Render()
		if err != nil {
			return fmt.Errorf("rendering table: %w", err)
		}
		fmt.Fprintln(cfg.Writer, result)
		return nil
	case FormatJSON:
		result, err := output.MarshalJSON(data)
		if err != nil {
			return fmt.Errorf("marshaling JSON: %w", err)
		}
		fmt.Fprintln(cfg.Writer, string(result))
		return nil
	case FormatCSV:
		return outputCSV(cfg.Writer, data)
	case FormatTSV:
		return outputTSV(cfg.Writer, data)
	case FormatMarkdown:
		md := output.NewMarkdownTableFromData(data)
		result, err := md.Render()
		if err != nil {
			return fmt.Errorf("rendering markdown: %w", err)
		}
		fmt.Fprintln(cfg.Writer, result)
		return nil
	case FormatYAML:
		result, err := output.MarshalYAML(data)
		if err != nil {
			return fmt.Errorf("marshaling YAML: %w", err)
		}
		fmt.Fprintln(cfg.Writer, string(result))
		return nil
	case FormatXML:
		result, err := output.MarshalXMLFromTableData(data)
		if err != nil {
			return fmt.Errorf("marshaling XML: %w", err)
		}
		fmt.Fprintln(cfg.Writer, string(result))
		return nil
	case FormatHTML:
		r := output.NewHTMLRenderer()
		r.SetData(data)
		result, err := r.Render()
		if err != nil {
			return fmt.Errorf("rendering HTML: %w", err)
		}
		fmt.Fprintln(cfg.Writer, result)
		return nil
	case FormatTree:
		renderer := output.TreeRendererFromTableData(data)
		result, err := renderer.Render()
		if err != nil {
			return fmt.Errorf("rendering tree: %w", err)
		}
		fmt.Fprintln(cfg.Writer, result)
		return nil
	case FormatD2:
		result, err := output.D2FromTableData(data).Render()
		if err != nil {
			return fmt.Errorf("rendering D2: %w", err)
		}
		fmt.Fprintln(cfg.Writer, result)
		return nil
	case FormatMermaid:
		renderer := output.MermaidFlowchartRenderer(data)
		result, err := renderer.Render()
		if err != nil {
			return fmt.Errorf("rendering Mermaid: %w", err)
		}
		fmt.Fprintln(cfg.Writer, result)
		return nil
	case FormatDOT:
		renderer := output.DOTFromTableData(data)
		result, err := renderer.Render()
		if err != nil {
			return fmt.Errorf("rendering DOT: %w", err)
		}
		fmt.Fprintln(cfg.Writer, result)
		return nil
	default:
		return fmt.Errorf("unsupported output format: %s", cfg.Format)
	}
}

func outputAny(cfg OutputConfig, data any) error {
	switch cfg.Format {
	case FormatJSON:
		result, err := output.MarshalJSONIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling JSON: %w", err)
		}
		fmt.Fprintln(cfg.Writer, string(result))
		return nil
	case FormatYAML:
		result, err := output.MarshalYAML(data)
		if err != nil {
			return fmt.Errorf("marshaling YAML: %w", err)
		}
		fmt.Fprintln(cfg.Writer, string(result))
		return nil
	default:
		return fmt.Errorf("format %s requires typed data (TableData, GraphNode, etc.)", cfg.Format)
	}
}

func outputCSV(w io.Writer, data *output.TableData) error {
	cw := output.NewCSVWriter(w)
	cw.WriteHeader(data.GetHeaders())
	for _, row := range data.GetRows() {
		cw.WriteRow(row)
	}
	cw.Flush()
	return cw.Error()
}

func outputTSV(w io.Writer, data *output.TableData) error {
	result, err := output.MarshalTSV(data)
	if err != nil {
		return fmt.Errorf("marshaling TSV: %w", err)
	}
	fmt.Fprintln(w, string(result))
	return nil
}
