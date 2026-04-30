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

// tableFormatRegistry maps OutputFormat to rendering functions.
var tableFormatRegistry = map[OutputFormat]tableRenderer{
	output.FormatTable:    renderTableStyled,
	output.FormatJSON:     renderTableJSON,
	output.FormatCSV:      renderTableCSV,
	output.FormatTSV:      renderTableTSV,
	output.FormatMarkdown: renderTableMarkdown,
	output.FormatYAML:     renderTableYAML,
	output.FormatXML:      renderTableXML,
	output.FormatHTML:     renderTableHTML,
	output.FormatTree:     renderTableTree,
	output.FormatD2:       renderTableD2,
	output.FormatMermaid:  renderTableMermaid,
	output.FormatDOT:      renderTableDOT,
}

// anyRenderer renders arbitrary data to a writer.
type anyRenderer func(w io.Writer, data any) error

// anyFormatRegistry maps OutputFormat to generic rendering functions.
var anyFormatRegistry = map[OutputFormat]anyRenderer{
	output.FormatJSON: renderAnyJSON,
	output.FormatYAML: renderAnyYAML,
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

func renderTableJSON(w io.Writer, data *output.TableData) error {
	result, err := output.MarshalJSON(data)
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}

	fmt.Fprintln(w, string(result))

	return nil
}

func renderTableCSV(w io.Writer, data *output.TableData) error {
	cw := output.NewCSVWriter(w)

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

func renderTableTSV(w io.Writer, data *output.TableData) error {
	result, err := output.MarshalTSV(data)
	if err != nil {
		return fmt.Errorf("marshaling TSV: %w", err)
	}

	fmt.Fprintln(w, string(result))

	return nil
}

func renderTableMarkdown(w io.Writer, data *output.TableData) error {
	md := output.NewMarkdownTableFromData(data)

	result, err := md.Render()
	if err != nil {
		return fmt.Errorf("rendering markdown: %w", err)
	}

	fmt.Fprintln(w, result)

	return nil
}

func renderTableYAML(w io.Writer, data *output.TableData) error {
	result, err := output.MarshalYAML(data)
	if err != nil {
		return fmt.Errorf("marshaling YAML: %w", err)
	}

	fmt.Fprintln(w, string(result))

	return nil
}

func renderTableXML(w io.Writer, data *output.TableData) error {
	result, err := output.MarshalXMLFromTableData(data)
	if err != nil {
		return fmt.Errorf("marshaling XML: %w", err)
	}

	fmt.Fprintln(w, string(result))

	return nil
}

func renderTableHTML(w io.Writer, data *output.TableData) error {
	r := output.NewHTMLRenderer()
	r.SetData(data)

	result, err := r.Render()
	if err != nil {
		return fmt.Errorf("rendering HTML: %w", err)
	}

	fmt.Fprintln(w, result)

	return nil
}

func renderTableTree(w io.Writer, data *output.TableData) error {
	renderer := output.TreeRendererFromTableData(data)

	result, err := renderer.Render()
	if err != nil {
		return fmt.Errorf("rendering tree: %w", err)
	}

	fmt.Fprintln(w, result)

	return nil
}

func renderTableD2(w io.Writer, data *output.TableData) error {
	result, err := output.D2FromTableData(data).Render()
	if err != nil {
		return fmt.Errorf("rendering D2: %w", err)
	}

	fmt.Fprintln(w, result)

	return nil
}

func renderTableMermaid(w io.Writer, data *output.TableData) error {
	renderer := output.MermaidFlowchartRenderer(data)

	result, err := renderer.Render()
	if err != nil {
		return fmt.Errorf("rendering Mermaid: %w", err)
	}

	fmt.Fprintln(w, result)

	return nil
}

func renderTableDOT(w io.Writer, data *output.TableData) error {
	renderer := output.DOTFromTableData(data)

	result, err := renderer.Render()
	if err != nil {
		return fmt.Errorf("rendering DOT: %w", err)
	}

	fmt.Fprintln(w, result)

	return nil
}

func renderAnyJSON(w io.Writer, data any) error {
	result, err := output.MarshalJSONIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}

	fmt.Fprintln(w, string(result))

	return nil
}

func renderAnyYAML(w io.Writer, data any) error {
	result, err := output.MarshalYAML(data)
	if err != nil {
		return fmt.Errorf("marshaling YAML: %w", err)
	}

	fmt.Fprintln(w, string(result))

	return nil
}
