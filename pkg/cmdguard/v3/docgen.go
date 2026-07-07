package v3

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

const (
	docHeadingMin = 1
	docHeadingMax = 6
)

// GenerateDocs writes markdown documentation for the CLI's full command tree
// to w. Each command (including subcommands) is rendered with its synopsis,
// usage, flags, and examples.
func (cli *CLI[T]) GenerateDocs(w io.Writer) error {
	var b strings.Builder

	err := writeCommandDocs(&b, cli.rootCmd, docHeadingMin)
	if err != nil {
		return err
	}

	if _, err := io.WriteString(w, b.String()); err != nil {
		return fmt.Errorf("writing generated docs: %w", err)
	}

	return nil
}

// writeCommandDocs recursively writes markdown for cmd and all available
// subcommands into b.
func writeCommandDocs(b *strings.Builder, cmd *cobra.Command, depth int) error {
	if !cmd.IsAvailableCommand() {
		return nil
	}

	heading := strings.Repeat("#", clampDepth(depth))

	synopsis := cmd.Short
	if cmd.Long != "" {
		synopsis = cmd.Long
	}

	fmt.Fprintf(b, "%s %s\n\n%s\n\n", heading, cmd.CommandPath(), synopsis)

	if cmd.Example != "" {
		fmt.Fprintf(b, "### Examples\n\n```\n%s\n```\n\n", cmd.Example)
	}

	fmt.Fprintf(b, "### Usage\n\n```\n%s\n```\n\n", cmd.UseLine())

	if flags := cmd.LocalFlags().FlagUsages(); flags != "" {
		fmt.Fprintf(b, "### Flags\n\n```\n%s```\n\n", flags)
	}

	for _, sub := range cmd.Commands() {
		if sub.IsAdditionalHelpTopicCommand() || sub.Hidden {
			continue
		}

		err := writeCommandDocs(b, sub, depth+1)
		if err != nil {
			return err
		}
	}

	return nil
}

// clampDepth keeps markdown headings within the valid 1-6 range.
func clampDepth(d int) int {
	if d < docHeadingMin {
		return docHeadingMin
	}

	if d > docHeadingMax {
		return docHeadingMax
	}

	return d
}
