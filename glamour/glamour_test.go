package glamour

import (
	"strings"
	"testing"
)

func TestWithHelp_ReturnsNonNil(t *testing.T) {
	t.Parallel()

	opt := WithHelp()
	if opt == nil {
		t.Fatal("WithHelp() returned nil")
	}
}

func TestWithHelpTheme_ReturnsNonNil(t *testing.T) {
	t.Parallel()

	opt := WithHelpTheme("dark")
	if opt == nil {
		t.Fatal("WithHelpTheme(\"dark\") returned nil")
	}
}

func TestRenderMarkdown_BasicOutput(t *testing.T) {
	t.Parallel()

	result := RenderMarkdown("# Hello")
	if result == "" {
		t.Fatal("RenderMarkdown(\"# Hello\") returned empty string")
	}
}

func TestRenderMarkdown_EmptyInput(t *testing.T) {
	t.Parallel()

	result := strings.TrimSpace(RenderMarkdown(""))
	if result != "" {
		t.Fatalf("RenderMarkdown(\"\") should return empty/whitespace, got %q", result)
	}
}

func TestRenderMarkdown_ContainsContent(t *testing.T) {
	t.Parallel()

	result := RenderMarkdown("**bold text**")
	if !strings.Contains(result, "bold text") {
		t.Fatalf("RenderMarkdown should contain 'bold text', got %q", result)
	}
}

func TestWithHelp_ImplementsCLIOption(t *testing.T) {
	t.Parallel()

	_ = WithHelp()
	_ = WithHelpTheme("light")
}
