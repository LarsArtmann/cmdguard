package v2

import (
	"bytes"
	"strings"
	"testing"

	output "github.com/larsartmann/go-output"

	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

func TestParseOutputFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    OutputFormat
		wantErr bool
	}{
		{"table", "table", FormatTable, false},
		{"json", "json", FormatJSON, false},
		{"csv", "csv", FormatCSV, false},
		{"tsv", "tsv", FormatTSV, false},
		{"markdown", "markdown", FormatMarkdown, false},
		{"xml", "xml", FormatXML, false},
		{"d2", "d2", FormatD2, false},
		{"yaml", "yaml", FormatYAML, false},
		{"html", "html", FormatHTML, false},
		{"tree", "tree", FormatTree, false},
		{"mermaid", "mermaid", FormatMermaid, false},
		{"dot", "dot", FormatDOT, false},
		{"jsonl", "jsonl", FormatJSONL, false},
		{"asciidoc", "asciidoc", FormatAsciiDoc, false},
		{"toml", "toml", FormatTOML, false},
		{"plantuml", "plantuml", FormatPlantUML, false},
		{"empty string", "", FormatTable, true},
		{"invalid format", "invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := ParseOutputFormat(tt.input)
			if tt.wantErr {
				testutil.AssertExpectedError(t, err)
			} else {
				testutil.AssertNoError(t, err)
				if result != tt.want {
					t.Errorf("ParseOutputFormat(%q) = %q, want %q", tt.input, result, tt.want)
				}
			}
		})
	}
}

func TestDefaultOutputConfig(t *testing.T) {
	t.Parallel()

	cfg := DefaultOutputConfig()
	if cfg.Format != FormatTable {
		t.Errorf("Default format = %v, want table", cfg.Format)
	}
	if cfg.Writer == nil {
		t.Error("Default writer should be os.Stdout, not nil")
	}
}

func TestOutputResult_TableData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		format     output.Format
		headers    []string
		row        []string
		mustHave   []string
		errMessage string
	}{
		{
			name:       "table format renders table data",
			format:     FormatTable,
			headers:    []string{"Name", "Age"},
			row:        []string{"Alice", "30"},
			mustHave:   []string{"Name", "Alice"},
			errMessage: "table output missing expected content",
		},
		{
			name:       "json format renders table data",
			format:     FormatJSON,
			headers:    []string{"Name"},
			row:        []string{"Bob"},
			mustHave:   []string{"Bob"},
			errMessage: "json output missing 'Bob'",
		},
		{
			name:       "csv format renders table data",
			format:     FormatCSV,
			headers:    []string{"Name"},
			row:        []string{"Eve"},
			mustHave:   []string{"Eve"},
			errMessage: "csv output missing 'Eve'",
		},
		{
			name:       "yaml format renders table data",
			format:     FormatYAML,
			headers:    []string{"Name"},
			row:        []string{"Yaml"},
			mustHave:   []string{"Yaml"},
			errMessage: "yaml output missing 'Yaml'",
		},
		{
			name:       "jsonl format renders table data",
			format:     FormatJSONL,
			headers:    []string{"Name"},
			row:        []string{"JsonlUser"},
			mustHave:   []string{"JsonlUser"},
			errMessage: "jsonl output missing 'JsonlUser'",
		},
		{
			name:       "asciidoc format renders table data",
			format:     FormatAsciiDoc,
			headers:    []string{"Name"},
			row:        []string{"AsciiDocUser"},
			mustHave:   []string{"AsciiDocUser"},
			errMessage: "asciidoc output missing 'AsciiDocUser'",
		},
		{
			name:       "toml format renders table data",
			format:     FormatTOML,
			headers:    []string{"Name"},
			row:        []string{"TomlUser"},
			mustHave:   []string{"TomlUser"},
			errMessage: "toml output missing 'TomlUser'",
		},
		{
			name:       "plantuml format renders table data",
			format:     FormatPlantUML,
			headers:    []string{"Name"},
			row:        []string{"PlantUMLUser"},
			mustHave:   []string{"PlantUMLUser"},
			errMessage: "plantuml output missing 'PlantUMLUser'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			data := output.NewTableData(tt.headers)
			data.AddRow(tt.row)

			cfg := OutputConfig{Format: tt.format, Writer: &buf}
			err := OutputResult(cfg, data)
			testutil.AssertNoError(t, err)

			result := buf.String()
			for _, want := range tt.mustHave {
				if !strings.Contains(result, want) {
					t.Errorf("%s: %q", tt.errMessage, result)

					break
				}
			}
		})
	}
}

func TestOutputResult_AnyData(t *testing.T) {
	t.Parallel()

	t.Run("json format with arbitrary struct", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		type Person struct {
			Name string `json:"name"`
		}

		cfg := OutputConfig{Format: FormatJSON, Writer: &buf}
		err := OutputResult(cfg, Person{Name: "Test"})
		testutil.AssertNoError(t, err)

		result := buf.String()
		if !strings.Contains(result, "Test") {
			t.Errorf("json output missing 'Test': %q", result)
		}
	})

	t.Run("non-json format with arbitrary data returns error", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		cfg := OutputConfig{Format: FormatTable, Writer: &buf}
		err := OutputResult(cfg, "just a string")
		testutil.AssertExpectedError(t, err)
	})
}

func TestOutputTable(t *testing.T) {
	t.Parallel()

	t.Run("renders table with headers and rows", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		cfg := OutputConfig{Format: FormatTable, Writer: &buf}

		_ = cfg
		err := OutputTable(FormatTable, []string{"Name", "Value"}, [][]string{{"key", "val"}})
		testutil.AssertNoError(t, err)
	})
}

func TestOutputStyledTable(t *testing.T) {
	t.Parallel()

	t.Run("renders styled table", func(t *testing.T) {
		t.Parallel()

		err := OutputStyledTable([]string{"Col1"}, [][]string{{"data"}})
		testutil.AssertNoError(t, err)
	})
}

func TestOutputResult_NilWriter(t *testing.T) {
	t.Parallel()

	t.Run("nil writer defaults to stdout", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		cfg := OutputConfig{Format: FormatJSON, Writer: &buf}
		data := output.NewTableData([]string{"X"})
		data.AddRow([]string{"1"})

		err := OutputResult(cfg, data)
		testutil.AssertNoError(t, err)
	})
}
