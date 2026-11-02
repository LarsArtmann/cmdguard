package v4

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

const (
	docHeadingMin = 1
	docHeadingMax = 6

	// docEmptyCell keeps table cells non-empty so markdown renderers do not
	// collapse a column when every value in it is absent.
	docEmptyCell = "-"
)

// GenerateDocs writes markdown documentation for the CLI's full command tree
// to w. Each command (including subcommands) is rendered with its synopsis,
// usage, flags, and examples.
//
// The flag tables are derived from the same contract [CLI.Describe] emits, not
// from pflag's help text. That makes the human documentation and the
// machine-readable contract one source of truth: a fact added to a struct tag
// appears in both, and [TestGenerateDocs_MatchesDescribeContract] fails when it
// stops appearing in either.
func (cli *CLI[T]) GenerateDocs(w io.Writer) error {
	var b strings.Builder

	err := writeCommandDocs(&b, cli.rootCmd, docHeadingMin)
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, b.String())
	if err != nil {
		return fmt.Errorf("writing generated docs: %w", err)
	}

	return nil
}

// writeCommandDocs recursively writes markdown for cmd and all available
// subcommands into b.
func writeCommandDocs(b *strings.Builder, cmd *cobra.Command, depth int) error {
	// The root is the entry point and is always documented, even when it has
	// no subcommands and no Run of its own; a root-only CLI otherwise produced
	// an empty document. Descendants follow cobra's availability rules.
	if cmd.HasParent() && !cmd.IsAvailableCommand() {
		return nil
	}

	// Document the flags cobra synthesizes too: --help exists on every command
	// and the machine contract lists it, so hiding it here would make the two
	// representations disagree.
	cmd.InitDefaultHelpFlag()

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

	local, inherited := splitInheritedFlags(flagsFor(cmd))
	writeFlagTable(b, "### Flags", visibleFlags(local))
	writeFlagTable(b, "### Inherited Flags", visibleFlags(inherited))

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

// splitInheritedFlags separates a command's own flags from the ones it accepts
// from its parents. Inherited flags get their own section: repeating the root's
// flags under every command buries the flags that command actually adds.
func splitInheritedFlags(flags []FlagDescription) ([]FlagDescription, []FlagDescription) {
	var local, inherited []FlagDescription

	for _, flag := range flags {
		if flag.Inherited {
			inherited = append(inherited, flag)

			continue
		}

		local = append(local, flag)
	}

	return local, inherited
}

// visibleFlags drops hidden flags. They are part of the machine contract (an
// agent may still call them) but explicitly not part of the human docs.
func visibleFlags(flags []FlagDescription) []FlagDescription {
	visible := make([]FlagDescription, 0, len(flags))

	for _, flag := range flags {
		if !flag.Hidden {
			visible = append(visible, flag)
		}
	}

	return visible
}

// writeFlagTable renders one markdown table of flags. Nothing is written when
// there are no flags, so a parent command without its own flags stays short.
func writeFlagTable(b *strings.Builder, heading string, flags []FlagDescription) {
	if len(flags) == 0 {
		return
	}

	fmt.Fprintf(b, "%s\n\n", heading)
	b.WriteString("| Flag | Type | Default | Env | Values | Required | Secret | Description |\n")
	b.WriteString("| ---- | ---- | ------- | --- | ------ | -------- | ------ | ----------- |\n")

	for _, flag := range flags {
		fmt.Fprintf(
			b,
			"| %s | %s | %s | %s | %s | %s | %s | %s |\n",
			docFlagCell(flag),
			docOrDash(flag.Type),
			docCodeOrDash(flag.Default),
			docCodeOrDash(flag.Env),
			docValuesCell(flag.Values),
			docRequiredCell(flag.Required),
			docRequiredCell(flag.Secret),
			docTextCell(flag.Usage),
		)
	}

	b.WriteString("\n")
}

// docFlagCell renders the flag name and shorthand as inline code.
func docFlagCell(flag FlagDescription) string {
	cell := "`--" + flag.Name + "`"
	if flag.Shorthand != "" {
		cell += ", `-" + flag.Shorthand + "`"
	}

	return cell
}

// docOrDash renders a plain value, or a dash when it is absent.
func docOrDash(value string) string {
	if value == "" {
		return docEmptyCell
	}

	return docCellText(value)
}

// docCodeOrDash renders a value as inline code, or a dash when it is absent.
func docCodeOrDash(value string) string {
	if value == "" {
		return docEmptyCell
	}

	return "`" + docCellText(value) + "`"
}

// docValuesCell renders the accepted enum values.
func docValuesCell(values []string) string {
	if len(values) == 0 {
		return docEmptyCell
	}

	rendered := make([]string, 0, len(values))
	for _, value := range values {
		rendered = append(rendered, "`"+docCellText(value)+"`")
	}

	return strings.Join(rendered, ", ")
}

// docRequiredCell renders whether the flag must be supplied.
func docRequiredCell(required bool) string {
	if !required {
		return docEmptyCell
	}

	return "yes"
}

// docTextCell sanitizes prose for a markdown table cell.
func docTextCell(text string) string {
	text = docCellText(text)
	if text == "" {
		return docEmptyCell
	}

	return text
}

// docCellText flattens text so it cannot break out of a table row.
func docCellText(text string) string {
	text = strings.ReplaceAll(text, "\r\n", " ")
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\\", `\\`)
	text = strings.ReplaceAll(text, "|", `\|`)

	return strings.TrimSpace(text)
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
