package v2

import (
	"fmt"

	"charm.land/glamour/v2"
	"github.com/spf13/cobra"
)

// WithGlamourHelp enables markdown rendering for command help text.
// When enabled, the Long and Example fields of all commands are passed
// through glamour for styled terminal output. The theme is determined by
// the GLAMOUR_STYLE environment variable, defaulting to "dark".
//
// Usage:
//
//	cli, _ := v2.NewCLI[Config]("app", "My app", Config{},
//	    v2.WithGlamourHelp[Config](),
//	)
func WithGlamourHelp[T any]() CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.glamourHelp = true
		cli.glamourTheme = ""
	}
}

// WithGlamourHelpTheme enables markdown rendering with a specific glamour theme.
// Supported themes: "ascii", "dark", "dracula", "light", "notty",
// "pink", "tokyo-night".
//
// Usage:
//
//	cli, _ := v2.NewCLI[Config]("app", "My app", Config{},
//	    v2.WithGlamourHelpTheme[Config]("dark"),
//	)
func WithGlamourHelpTheme[T any](theme string) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.glamourHelp = true
		cli.glamourTheme = theme
	}
}

// applyGlamourHelp recursively renders command Long and Example descriptions
// with glamour. It mutates the cobra command fields in place.
// An empty theme uses environment-based detection (GLAMOUR_STYLE env var).
func applyGlamourHelp(cmd *cobra.Command, theme string) {
	render := func(markdown string) (string, error) {
		if theme == "" {
			return glamour.RenderWithEnvironmentConfig(markdown)
		}

		return glamour.Render(markdown, theme)
	}

	if cmd.Long != "" {
		rendered, err := render(cmd.Long)
		if err == nil {
			cmd.Long = rendered
		}
	}

	if cmd.Example != "" {
		rendered, err := render(cmd.Example)
		if err == nil {
			cmd.Example = rendered
		}
	}

	for _, sub := range cmd.Commands() {
		applyGlamourHelp(sub, theme)
	}
}

// renderGlamourOrFallback renders markdown with glamour, falling back to the
// raw text on error. Uses environment-based theme detection.
func renderGlamourOrFallback(markdown string) string {
	rendered, err := glamour.RenderWithEnvironmentConfig(markdown)
	if err != nil {
		return markdown
	}

	return rendered
}

// RenderMarkdown is a convenience helper that renders markdown to styled
// terminal output using glamour. Uses environment-based theme detection.
// On error, the original markdown is returned.
func RenderMarkdown(markdown string) string {
	return renderGlamourOrFallback(markdown)
}

// RenderMarkdownWithTheme renders markdown with a specific glamour theme.
// Supported themes: "ascii", "dark", "dracula", "light", "notty",
// "pink", "tokyo-night".
func RenderMarkdownWithTheme(markdown, theme string) string {
	rendered, err := glamour.Render(markdown, theme)
	if err != nil {
		return fmt.Sprintf("render error: %v\n%s", err, markdown)
	}

	return rendered
}
