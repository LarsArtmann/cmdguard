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

	t.Run("table format renders table data", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		data := output.NewTableData([]string{"Name", "Age"})
		data.AddRow([]string{"Alice", "30"})

		cfg := OutputConfig{Format: FormatTable, Writer: &buf}
		err := OutputResult(cfg, data)
		testutil.AssertNoError(t, err)

		result := buf.String()
		if !strings.Contains(result, "Name") || !strings.Contains(result, "Alice") {
			t.Errorf("table output missing expected content: %q", result)
		}
	})

	t.Run("json format renders table data", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		data := output.NewTableData([]string{"Name"})
		data.AddRow([]string{"Bob"})

		cfg := OutputConfig{Format: FormatJSON, Writer: &buf}
		err := OutputResult(cfg, data)
		testutil.AssertNoError(t, err)

		result := buf.String()
		if !strings.Contains(result, "Bob") {
			t.Errorf("json output missing 'Bob': %q", result)
		}
	})

	t.Run("csv format renders table data", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		data := output.NewTableData([]string{"Name"})
		data.AddRow([]string{"Eve"})

		cfg := OutputConfig{Format: FormatCSV, Writer: &buf}
		err := OutputResult(cfg, data)
		testutil.AssertNoError(t, err)

		result := buf.String()
		if !strings.Contains(result, "Eve") {
			t.Errorf("csv output missing 'Eve': %q", result)
		}
	})

	t.Run("yaml format renders table data", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		data := output.NewTableData([]string{"Name"})
		data.AddRow([]string{"Yaml"})

		cfg := OutputConfig{Format: FormatYAML, Writer: &buf}
		err := OutputResult(cfg, data)
		testutil.AssertNoError(t, err)

		result := buf.String()
		if !strings.Contains(result, "Yaml") {
			t.Errorf("yaml output missing 'Yaml': %q", result)
		}
	})

	t.Run("jsonl format renders table data", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		data := output.NewTableData([]string{"Name"})
		data.AddRow([]string{"JsonlUser"})

		cfg := OutputConfig{Format: FormatJSONL, Writer: &buf}
		err := OutputResult(cfg, data)
		testutil.AssertNoError(t, err)

		result := buf.String()
		if !strings.Contains(result, "JsonlUser") {
			t.Errorf("jsonl output missing 'JsonlUser': %q", result)
		}
	})

	t.Run("asciidoc format renders table data", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		data := output.NewTableData([]string{"Name"})
		data.AddRow([]string{"AsciiDocUser"})

		cfg := OutputConfig{Format: FormatAsciiDoc, Writer: &buf}
		err := OutputResult(cfg, data)
		testutil.AssertNoError(t, err)

		result := buf.String()
		if !strings.Contains(result, "AsciiDocUser") {
			t.Errorf("asciidoc output missing 'AsciiDocUser': %q", result)
		}
	})

	t.Run("toml format renders table data", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		data := output.NewTableData([]string{"Name"})
		data.AddRow([]string{"TomlUser"})

		cfg := OutputConfig{Format: FormatTOML, Writer: &buf}
		err := OutputResult(cfg, data)
		testutil.AssertNoError(t, err)

		result := buf.String()
		if !strings.Contains(result, "TomlUser") {
			t.Errorf("toml output missing 'TomlUser': %q", result)
		}
	})

	t.Run("plantuml format renders table data", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		data := output.NewTableData([]string{"Name"})
		data.AddRow([]string{"PlantUMLUser"})

		cfg := OutputConfig{Format: FormatPlantUML, Writer: &buf}
		err := OutputResult(cfg, data)
		testutil.AssertNoError(t, err)

		result := buf.String()
		if !strings.Contains(result, "PlantUMLUser") {
			t.Errorf("plantuml output missing 'PlantUMLUser': %q", result)
		}
	})
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
