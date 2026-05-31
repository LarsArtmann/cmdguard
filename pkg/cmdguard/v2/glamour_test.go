package v2

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRenderMarkdown(t *testing.T) {
	t.Parallel()

	result := RenderMarkdown("# Hello\n\nThis is **bold**.")

	if result == "" {
		t.Error("RenderMarkdown should return non-empty output")
	}

	if result == "# Hello\n\nThis is **bold**." {
		t.Error("RenderMarkdown should transform the input, not return it verbatim")
	}
}

func TestRenderMarkdownWithTheme(t *testing.T) {
	t.Parallel()

	result := RenderMarkdownWithTheme("# Hello", "dark")

	if result == "" {
		t.Error("RenderMarkdownWithTheme should return non-empty output")
	}

	if strings.Contains(result, "# Hello") {
		t.Error("RenderMarkdownWithTheme(dark) should render markdown, not return raw input")
	}
}

func TestRenderMarkdownWithTheme_InvalidTheme(t *testing.T) {
	t.Parallel()

	result := RenderMarkdownWithTheme("# Hello", "nonexistent-theme-xyz")

	if result == "" {
		t.Error("RenderMarkdownWithTheme should fall back to raw markdown on error")
	}

	if !strings.Contains(result, "# Hello") {
		t.Error("fallback should contain original markdown")
	}
}

func TestWithGlamourHelp(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	cli, err := NewCLI[testConfig](
		"test", "Test CLI", testConfig{},
		WithGlamourHelp[testConfig](),
		WithFang[testConfig](false),
	)
	if err != nil {
		t.Fatalf("NewCLI failed: %v", err)
	}

	if !cli.glamourHelp {
		t.Error("glamourHelp should be true")
	}

	if cli.glamourTheme != "auto" {
		t.Errorf("expected theme %q, got %q", "auto", cli.glamourTheme)
	}
}

func TestWithGlamourHelpTheme(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	cli, err := NewCLI[testConfig](
		"test", "Test CLI", testConfig{},
		WithGlamourHelpTheme[testConfig]("dark"),
		WithFang[testConfig](false),
	)
	if err != nil {
		t.Fatalf("NewCLI failed: %v", err)
	}

	if !cli.glamourHelp {
		t.Error("glamourHelp should be true")
	}

	if cli.glamourTheme != "dark" {
		t.Errorf("expected theme %q, got %q", "dark", cli.glamourTheme)
	}
}

func TestApplyGlamourHelp(t *testing.T) {
	t.Parallel()

	original := "# Title\n\nSome **bold** text."
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Short desc",
		Long:  original,
	}

	applyGlamourHelp(cmd, "auto")

	if cmd.Long == "" {
		t.Error("Long should not be empty after glamour")
	}

	if cmd.Long == original {
		t.Error("Long should be transformed by glamour")
	}
}

func TestApplyGlamourHelp_Example(t *testing.T) {
	t.Parallel()

	original := "```bash\necho hello\n```"
	cmd := &cobra.Command{
		Use:     "test",
		Short:   "Short desc",
		Example: original,
	}

	applyGlamourHelp(cmd, "auto")

	if cmd.Example == "" {
		t.Error("Example should not be empty after glamour")
	}

	if cmd.Example == original {
		t.Error("Example should be transformed by glamour")
	}
}

func TestApplyGlamourHelp_Subcommands(t *testing.T) {
	t.Parallel()

	child := &cobra.Command{
		Use:   "child",
		Short: "Child command",
		Long:  "# Child\n\nA child description.",
	}

	parent := &cobra.Command{
		Use:   "parent",
		Short: "Parent command",
		Long:  "# Parent\n\nA parent description.",
	}

	parent.AddCommand(child)

	applyGlamourHelp(parent, "auto")

	if child.Long == "" {
		t.Error("child Long should not be empty after glamour")
	}
}

//nolint:paralleltest // captures os.Stdout
func TestGlamourHelp_E2ERendering(t *testing.T) {
	type testConfig struct{}

	markdownLong := "# Overview\n\nThis is **bold** and *italic* text.\n\n## Details\n\n- Item 1\n- Item 2"

	cli, err := NewCLI[testConfig](
		"testapp", "Test Application", testConfig{},
		WithGlamourHelpTheme[testConfig]("dark"),
		WithFang[testConfig](false),
	)
	if err != nil {
		t.Fatalf("NewCLI failed: %v", err)
	}

	cmd, err := NewCommand[testConfig, NoFlags](
		"greet",
		func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
		WithShort[testConfig, NoFlags]("Greet someone"),
		WithLong[testConfig, NoFlags](markdownLong),
	)
	if err != nil {
		t.Fatalf("NewCommand failed: %v", err)
	}

	if err := AddCommand(cli, cmd); err != nil {
		t.Fatalf("AddCommand failed: %v", err)
	}

	var buf bytes.Buffer
	cli.rootCmd.SetOut(&buf)
	cli.rootCmd.SetArgs([]string{"greet", "--help"})

	_ = cli.Execute(context.Background())

	output := buf.String()

	if output == "" {
		t.Fatal("help output should not be empty")
	}

	if strings.Contains(output, "# Overview") {
		t.Error("raw markdown heading should be rendered, not kept verbatim")
	}

	if strings.Contains(output, "**bold**") {
		t.Error("raw bold markdown syntax should be rendered, not kept verbatim")
	}

	if !strings.Contains(output, "bold") {
		t.Error("rendered output should contain the word 'bold'")
	}

	if !strings.Contains(output, "italic") {
		t.Error("rendered output should contain the word 'italic'")
	}
}
