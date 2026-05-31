package v2

import (
	"fmt"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

// WithGlamourHelp enables markdown rendering for command help text.
// When enabled, the Long and Example fields of all commands are passed
// through glamour for styled terminal output.
//
// Usage:
//
//	cli, _ := v2.NewCLI[Config]("app", "My app", Config{},
//	    v2.WithGlamourHelp[Config](),
//	)
func WithGlamourHelp[T any]() CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.glamourHelp = true
	}
}

// applyGlamourHelp recursively renders command Long and Example descriptions
// with glamour. It mutates the cobra command fields in place.
func applyGlamourHelp(cmd *cobra.Command) {
	if cmd.Long != "" {
		rendered, err := glamour.Render(cmd.Long, "auto")
		if err == nil {
			cmd.Long = rendered
		}
	}

	if cmd.Example != "" {
		rendered, err := glamour.Render(cmd.Example, "auto")
		if err == nil {
			cmd.Example = rendered
		}
	}

	for _, sub := range cmd.Commands() {
		applyGlamourHelp(sub)
	}
}

// renderGlamourOrFallback renders markdown with glamour, falling back to the
// raw text on error. Useful for one-off rendering outside of the CLI option.
func renderGlamourOrFallback(markdown string) string {
	rendered, err := glamour.Render(markdown, "auto")
	if err != nil {
		return markdown
	}

	return rendered
}

// RenderMarkdown is a convenience helper that renders markdown to styled
// terminal output using glamour. On error, the original markdown is returned.
func RenderMarkdown(markdown string) string {
	return renderGlamourOrFallback(markdown)
}

// RenderMarkdownWithTheme renders markdown with a specific glamour theme.
// Supported themes: "ascii", "auto", "dark", "dracula", "light", "notty",
// "pink", "tokyo-night".
func RenderMarkdownWithTheme(markdown, theme string) string {
	rendered, err := glamour.Render(markdown, theme)
	if err != nil {
		return fmt.Sprintf("render error: %v\n%s", err, markdown)
	}

	return rendered
}
