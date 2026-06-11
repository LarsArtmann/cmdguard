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

	allFormats := []struct {
		format   output.Format
		mustHave string
	}{
		{FormatTable, "Alice"},
		{FormatJSON, "Alice"},
		{FormatCSV, "Alice"},
		{FormatTSV, "Alice"},
		{FormatMarkdown, "Alice"},
		{FormatXML, "Alice"},
		{FormatD2, "row0"},
		{FormatYAML, "Alice"},
		{FormatHTML, "Alice"},
		{FormatTree, "Alice"},
		{FormatMermaid, "row0"},
		{FormatDOT, "row0"},
		{FormatJSONL, "Alice"},
		{FormatAsciiDoc, "Alice"},
		{FormatTOML, "Alice"},
		{FormatPlantUML, "row0"},
	}

	for _, tt := range allFormats {
		t.Run(string(tt.format), func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			data := output.NewTableData([]string{"Name"})
			data.AddRow([]string{"Alice"})

			cfg := OutputConfig{Format: tt.format, Writer: &buf}
			err := OutputResult(cfg, data)
			testutil.AssertNoError(t, err)

			if !strings.Contains(buf.String(), tt.mustHave) {
				t.Errorf("%s output missing %q: %q", tt.format, tt.mustHave, buf.String())
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

func TestOutputResult_NilData(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	cfg := OutputConfig{Format: FormatJSON, Writer: &buf}

	err := OutputResult(cfg, (*output.TableData)(nil))
	testutil.AssertNoError(t, err)
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

func TestUnwrapTableData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		data any
		want bool
	}{
		{"pointer table data", output.NewTableData([]string{"X"}), true},
		{"value table data", *output.NewTableData([]string{"X"}), true},
		{"string", "hello", false},
		{"nil", nil, false},
		{"int", 42, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := unwrapTableData(tt.data)
			if (result != nil) != tt.want {
				t.Errorf("unwrapTableData(%T) returned nil=%v, want non-nil=%v", tt.data, result == nil, tt.want)
			}
		})
	}
}

func TestOutputResult_UnsupportedFormat(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	cfg := OutputConfig{Format: OutputFormat("nonexistent"), Writer: &buf}

	err := OutputResult(cfg, output.NewTableData([]string{"X"}))
	testutil.AssertExpectedError(t, err)
	testutil.AssertErrorIs(t, err, ErrUnsupportedFormat)
}

func TestOutputResult_AnyData_YAML(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	type Item struct {
		Name string `yaml:"name"`
	}

	cfg := OutputConfig{Format: FormatYAML, Writer: &buf}
	err := OutputResult(cfg, Item{Name: "YamlAny"})
	testutil.AssertNoError(t, err)
	testutil.AssertOutputContains(t, buf.String(), "YamlAny")
}

func TestOutputResult_AnyData_TOML(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	type Item struct {
		Name string `toml:"name"`
	}

	cfg := OutputConfig{Format: FormatTOML, Writer: &buf}
	err := OutputResult(cfg, Item{Name: "TomlAny"})
	testutil.AssertNoError(t, err)
	testutil.AssertOutputContains(t, buf.String(), "TomlAny")
}

func TestOutputResult_AnyData_JSON(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	type Item struct {
		Name string `json:"name"`
	}

	cfg := OutputConfig{Format: FormatJSON, Writer: &buf}
	err := OutputResult(cfg, Item{Name: "JsonAny"})
	testutil.AssertNoError(t, err)
	testutil.AssertOutputContains(t, buf.String(), "JsonAny")
}

func TestOutputResult_TableOnlyFormats_RejectAnyData(t *testing.T) {
	t.Parallel()

	tableOnlyFormats := []OutputFormat{
		FormatTable, FormatCSV, FormatTSV, FormatXML, FormatMarkdown,
		FormatHTML, FormatTree, FormatD2, FormatMermaid, FormatDOT,
		FormatJSONL, FormatAsciiDoc, FormatPlantUML,
	}

	for _, f := range tableOnlyFormats {
		t.Run(string(f), func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			cfg := OutputConfig{Format: f, Writer: &buf}

			err := OutputResult(cfg, "not table data")
			testutil.AssertExpectedError(t, err)
			testutil.AssertErrorIs(t, err, ErrFormatRequiresTypedData)
		})
	}
}

func TestSupportedFormats(t *testing.T) {
	t.Parallel()

	formats := SupportedFormats()
	if len(formats) != 16 {
		t.Errorf("SupportedFormats() returned %d formats, want 16", len(formats))
	}
}

func TestIsFormatSupported(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		format OutputFormat
		want   bool
	}{
		{"table", FormatTable, true},
		{"json", FormatJSON, true},
		{"invalid", OutputFormat("nonexistent"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := IsFormatSupported(tt.format)
			if got != tt.want {
				t.Errorf("IsFormatSupported(%q) = %v, want %v", tt.format, got, tt.want)
			}
		})
	}
}
