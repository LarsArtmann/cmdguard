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

func TestFormatRegistry_Completeness(t *testing.T) {
	t.Parallel()

	expected := []OutputFormat{
		FormatTable, FormatJSON, FormatCSV, FormatTSV, FormatMarkdown,
		FormatXML, FormatD2, FormatYAML, FormatHTML, FormatTree,
		FormatMermaid, FormatDOT, FormatJSONL, FormatAsciiDoc, FormatTOML,
		FormatPlantUML,
	}

	for _, f := range expected {
		t.Run(string(f), func(t *testing.T) {
			t.Parallel()

			_, ok := formatRegistry[f]
			if !ok {
				t.Errorf("formatRegistry missing entry for %q", f)
			}
		})
	}
}

func TestTableRenderStrategy_RejectsNonTableData(t *testing.T) {
	t.Parallel()

	s := &tableRenderStrategy{label: "test", render: func(_ *output.TableData) (string, error) {
		return "", nil
	}}

	var buf bytes.Buffer
	err := s.Render(&buf, "not table data")
	testutil.AssertExpectedError(t, err)
	testutil.AssertErrorIs(t, err, ErrFormatRequiresTypedData)
}

func TestMarshalStrategy_RendersAnyData(t *testing.T) {
	t.Parallel()

	s := &marshalStrategy{label: "test", marshal: func(v any) ([]byte, error) {
		return []byte(`{"ok":true}`), nil
	}}

	var buf bytes.Buffer
	err := s.Render(&buf, struct{ X int }{X: 1})
	testutil.AssertNoError(t, err)
	testutil.AssertOutputContains(t, buf.String(), `"ok"`)
}

func TestMarshalStrategy_RendersTableData(t *testing.T) {
	t.Parallel()

	s := &marshalStrategy{label: "test", marshal: func(v any) ([]byte, error) {
		return []byte(`marshaled`), nil
	}}

	var buf bytes.Buffer
	data := output.NewTableData([]string{"H"})
	data.AddRow([]string{"V"})

	err := s.Render(&buf, data)
	testutil.AssertNoError(t, err)
	testutil.AssertOutputContains(t, buf.String(), "marshaled")
}

func TestDualStrategy_DelegatesToTable(t *testing.T) {
	t.Parallel()

	var tableCalled bool
	tableStrategy := &tableRenderStrategy{label: "inner", render: func(_ *output.TableData) (string, error) {
		tableCalled = true

		return "table output", nil
	}}
	anyStrategy := &marshalStrategy{label: "inner", marshal: func(_ any) ([]byte, error) {
		return []byte("any output"), nil
	}}

	s := &dualStrategy{table: tableStrategy, any: anyStrategy}
	var buf bytes.Buffer
	data := output.NewTableData([]string{"H"})

	err := s.Render(&buf, data)
	testutil.AssertNoError(t, err)
	testutil.AssertBoolTrue(t, tableCalled, "tableCalled")
}

func TestDualStrategy_DelegatesToAny(t *testing.T) {
	t.Parallel()

	var anyCalled bool
	tableStrategy := &tableRenderStrategy{label: "inner", render: func(_ *output.TableData) (string, error) {
		return "table output", nil
	}}
	anyStrategy := &marshalStrategy{label: "inner", marshal: func(_ any) ([]byte, error) {
		anyCalled = true

		return []byte("any output"), nil
	}}

	s := &dualStrategy{table: tableStrategy, any: anyStrategy}
	var buf bytes.Buffer

	err := s.Render(&buf, "arbitrary data")
	testutil.AssertNoError(t, err)
	testutil.AssertBoolTrue(t, anyCalled, "anyCalled")
}

func TestStyledTableStrategy_RejectsNonTableData(t *testing.T) {
	t.Parallel()

	s := &styledTableStrategy{}
	var buf bytes.Buffer

	err := s.Render(&buf, "not table data")
	testutil.AssertExpectedError(t, err)
	testutil.AssertErrorIs(t, err, ErrFormatRequiresTypedData)
}

func TestCSVStrategy_RejectsNonTableData(t *testing.T) {
	t.Parallel()

	s := &csvStrategy{}
	var buf bytes.Buffer

	err := s.Render(&buf, 42)
	testutil.AssertExpectedError(t, err)
	testutil.AssertErrorIs(t, err, ErrFormatRequiresTypedData)
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

func TestFormatStrategy_Interface(t *testing.T) {
	t.Parallel()

	var _ FormatStrategy = &tableRenderStrategy{}
	var _ FormatStrategy = &marshalStrategy{}
	var _ FormatStrategy = &dualStrategy{}
	var _ FormatStrategy = &styledTableStrategy{}
	var _ FormatStrategy = &csvStrategy{}
}
