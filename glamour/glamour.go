// Package glamour provides markdown help rendering for cmdguard CLIs.
// It is an optional module — import it only when you want styled markdown
// help output, to avoid pulling in chroma, goldmark, and bluemonday.
//
// Usage:
//
//	import (
//	    v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
//	    "github.com/larsartmann/cmdguard/glamour"
//	)
//
//	cli, _ := v3.NewCLI[Config]("app", "My app", Config{},
//	    glamour.WithHelp[T](),
//	)
package glamour

import (
	"fmt"

	glamourlib "charm.land/glamour/v2"
	"github.com/spf13/cobra"

	v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

// WithHelp enables markdown rendering for command help text.
// When enabled, the Long and Example fields of all commands are rendered
// through glamour for styled terminal output. The theme is determined by
// the GLAMOUR_STYLE environment variable, defaulting to "dark".
func WithHelp() v3.CLIOption {
	return v3.WithHelpTransform(func(cmd *cobra.Command) {
		applyToTree(cmd, "")
	})
}

// WithHelpTheme enables markdown rendering with a specific glamour theme.
// Supported themes: "ascii", "dark", "dracula", "light", "notty",
// "pink", "tokyo-night".
func WithHelpTheme(theme string) v3.CLIOption {
	return v3.WithHelpTransform(func(cmd *cobra.Command) {
		applyToTree(cmd, theme)
	})
}

// applyToTree recursively renders command Long and Example descriptions.
// An empty theme uses environment-based detection (GLAMOUR_STYLE env var).
func applyToTree(cmd *cobra.Command, theme string) {
	render := func(markdown string) (string, error) {
		if theme == "" {
			return glamourlib.RenderWithEnvironmentConfig(markdown)
		}

		return glamourlib.Render(markdown, theme)
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
		applyToTree(sub, theme)
	}
}

// RenderMarkdown renders markdown to styled terminal output using glamour.
// Uses environment-based theme detection. On error, the original markdown
// is returned.
func RenderMarkdown(markdown string) string {
	rendered, err := glamourlib.RenderWithEnvironmentConfig(markdown)
	if err != nil {
		return markdown
	}

	return rendered
}

// RenderMarkdownWithTheme renders markdown with a specific glamour theme.
func RenderMarkdownWithTheme(markdown, theme string) string {
	rendered, err := glamourlib.Render(markdown, theme)
	if err != nil {
		return fmt.Sprintf("render error: %v\n%s", err, markdown)
	}

	return rendered
}
