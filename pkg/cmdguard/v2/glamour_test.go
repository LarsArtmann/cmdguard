package v2

import (
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
}

func TestApplyGlamourHelp(t *testing.T) {
	t.Parallel()

	original := "# Title\n\nSome **bold** text."
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Short desc",
		Long:  original,
	}

	applyGlamourHelp(cmd)

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

	applyGlamourHelp(cmd)

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

	applyGlamourHelp(parent)

	if child.Long == "" {
		t.Error("child Long should not be empty after glamour")
	}
}
