package v2

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// GenerateDocs writes markdown documentation for the CLI's full command tree
// to w. Each command (including subcommands) is rendered with its synopsis,
// usage, flags, and examples. This is useful for generating a docs/ reference
// or a static site.
func (cli *CLI[T]) GenerateDocs(w io.Writer) error {
	return writeCommandDocs(w, cli.rootCmd)
}

// writeCommandDocs recursively writes markdown for cmd and all available
// subcommands. Help-topic and hidden commands are skipped.
func writeCommandDocs(w io.Writer, cmd *cobra.Command) error {
	if !cmd.IsAvailableCommand() {
		return nil
	}

	if err := cmd.GenMarkdown(w); err != nil {
		return fmt.Errorf("generating markdown for %q: %w", cmd.CommandPath(), err)
	}

	for _, sub := range cmd.Commands() {
		if sub.IsAdditionalHelpTopicCommand() || sub.Hidden {
			continue
		}

		if err := writeCommandDocs(w, sub); err != nil {
			return err
		}
	}

	return nil
}
