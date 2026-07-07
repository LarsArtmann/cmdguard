// Package manpage provides roff man page generation for cmdguard CLIs.
// It is an optional module — import it only when you need man page output,
// to keep your dependency tree lean.
//
// Usage:
//
//	import (
//	    v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
//	    "github.com/larsartmann/cmdguard/manpage"
//	)
//
//	content, _ := manpage.Generate(cli, 1)
//	fmt.Println(content)
package manpage

import (
	"fmt"
	"io"

	mango "github.com/muesli/mango"
	mcobra "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"

	v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

// Generate produces a roff man page for the CLI.
// Section is typically 1 for user commands or 8 for system commands.
func Generate[T any](cli *v3.CLI[T], section uint) (string, error) {
	mp, err := mcobra.NewManPage(section, cli.RootCommand())
	if err != nil {
		return "", fmt.Errorf("section=%d: %w", section, err)
	}

	return mp.Build(roff.NewDocument()), nil
}

// Write generates and writes a roff man page to the given writer.
func Write[T any](cli *v3.CLI[T], w io.Writer, section uint) error {
	content, err := Generate[T](cli, section)
	if err != nil {
		return fmt.Errorf("section=%d: %w", section, err)
	}

	_, err = fmt.Fprint(w, content)
	if err != nil {
		return fmt.Errorf("writing man page: %w", err)
	}

	return nil
}

// GenerateCommand creates a cobra command that generates man pages.
// Add this as a subcommand to provide `myapp man` functionality.
func GenerateCommand[T any](cli *v3.CLI[T]) (*cobra.Command, error) {
	return &cobra.Command{
		Use:   "man [section]",
		Short: "Generate man page",
		Long:  `Generate a roff man page for this CLI. Default section is 1.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			section := uint(1)
			if len(args) == 1 {
				_, err := fmt.Sscanf(args[0], "%d", &section)
				if err != nil {
					return fmt.Errorf("invalid section number: %w", err)
				}
			}

			return Write[T](cli, cmd.OutOrStdout(), section)
		},
	}, nil
}

// NewManPage creates a mango man page from a cobra command.
// Useful for custom man page generation pipelines.
func NewManPage(section uint, cmd *cobra.Command) (*mango.ManPage, error) {
	return mcobra.NewManPage(section, cmd)
}
