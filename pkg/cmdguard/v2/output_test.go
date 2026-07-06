package v2

import (
	"bytes"
	"strings"
	"testing"

	output "github.com/larsartmann/go-output"

	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

func TestDefaultOutputConfig(t *testing.T) {
	t.Parallel()

	cfg := DefaultOutputConfig()
	if cfg.Format != output.FormatTable {
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
		{output.FormatTable, "Alice"},
		{output.FormatJSON, "Alice"},
		{output.FormatCSV, "Alice"},
		{output.FormatTSV, "Alice"},
		{output.FormatMarkdown, "Alice"},
		{output.FormatXML, "Alice"},
		{output.FormatD2, "row0"},
		{output.FormatYAML, "Alice"},
		{output.FormatHTML, "Alice"},
		{output.FormatTree, "Alice"},
		{output.FormatMermaid, "row0"},
		{output.FormatDOT, "row0"},
		{output.FormatJSONL, "Alice"},
		{output.FormatAsciiDoc, "Alice"},
		{output.FormatTOML, "Alice"},
		{output.FormatPlantUML, "row0"},
	}

	for _, tt := range allFormats {
		t.Run(string(tt.format), func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			data := output.NewTable([]string{"Name"})
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

		cfg := OutputConfig{Format: output.FormatJSON, Writer: &buf}
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
		cfg := OutputConfig{Format: output.FormatTable, Writer: &buf}
		err := OutputResult(cfg, "just a string")
		testutil.AssertExpectedError(t, err)
	})
}

func TestOutputTable(t *testing.T) {
	t.Parallel()

	t.Run("renders table with headers and rows", func(t *testing.T) {
		t.Parallel()

		err := OutputTable(output.FormatTable, []string{"Name", "Value"}, [][]string{{"key", "val"}})
		testutil.AssertNoError(t, err)
	})

	t.Run("rejects mismatched row length", func(t *testing.T) {
		t.Parallel()

		err := OutputTable(output.FormatTable, []string{"Name", "Value"}, [][]string{{"only-one"}})
		testutil.AssertExpectedError(t, err)
	})
}

func TestOutputResult_NilData(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	cfg := OutputConfig{Format: output.FormatJSON, Writer: &buf}

	err := OutputResult(cfg, (*output.Table)(nil))
	testutil.AssertNoError(t, err)
}

func TestOutputResult_NilWriter(t *testing.T) {
	t.Parallel()

	t.Run("nil writer defaults to stdout", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		cfg := OutputConfig{Format: output.FormatJSON, Writer: &buf}
		data := output.NewTable([]string{"X"})
		data.AddRow([]string{"1"})

		err := OutputResult(cfg, data)
		testutil.AssertNoError(t, err)
	})
}

func TestOutputResult_UnsupportedFormat(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	cfg := OutputConfig{Format: OutputFormat("nonexistent"), Writer: &buf}

	err := OutputResult(cfg, output.NewTable([]string{"X"}))
	testutil.AssertExpectedError(t, err)
	testutil.AssertErrorIs(t, err, ErrUnsupportedFormat)
}

func TestOutputResult_AnyData_YAML(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	type Item struct {
		Name string `yaml:"name"`
	}

	cfg := OutputConfig{Format: output.FormatYAML, Writer: &buf}
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

	cfg := OutputConfig{Format: output.FormatTOML, Writer: &buf}
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

	cfg := OutputConfig{Format: output.FormatJSON, Writer: &buf}
	err := OutputResult(cfg, Item{Name: "JsonAny"})
	testutil.AssertNoError(t, err)
	testutil.AssertOutputContains(t, buf.String(), "JsonAny")
}

func TestOutputResult_TableOnlyFormats_RejectAnyData(t *testing.T) {
	t.Parallel()

	tableOnlyFormats := []OutputFormat{
		output.FormatTable, output.FormatCSV, output.FormatTSV, output.FormatXML, output.FormatMarkdown,
		output.FormatHTML, output.FormatTree, output.FormatD2, output.FormatMermaid, output.FormatDOT,
		output.FormatJSONL, output.FormatAsciiDoc, output.FormatPlantUML,
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

func TestRegisteredFormats(t *testing.T) {
	t.Parallel()

	formats := RegisteredFormats()
	if len(formats) == 0 {
		t.Fatal("RegisteredFormats() returned no formats")
	}

	hasTable := false
	for _, f := range formats {
		if f == output.FormatTable {
			hasTable = true
		}
	}
	if !hasTable {
		t.Error("RegisteredFormats() missing table format")
	}
}

func TestOutputResult_ShapeAwareError(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	cfg := OutputConfig{Format: OutputFormat("nonexistent"), Writer: &buf}

	err := OutputResult(cfg, output.NewTable([]string{"X"}))
	testutil.AssertExpectedError(t, err)

	errMsg := err.Error()
	if !strings.Contains(errMsg, "unsupported output format") {
		t.Errorf("error should mention unsupported format: %q", errMsg)
	}
}
