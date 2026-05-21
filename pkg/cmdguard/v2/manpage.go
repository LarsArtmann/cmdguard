package v2

import (
	"fmt"
	"io"

	mango "github.com/muesli/mango"
	mcobra "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"
)

// ManPage generates a roff man page for the CLI.
// Section is typically 1 for user commands or 8 for system commands.
func (cli *CLI[T]) ManPage(section uint) (string, error) {
	mp, err := mcobra.NewManPage(section, cli.rootCmd)
	if err != nil {
		return "", fmt.Errorf("section=%d: %w", section, err)
	}

	return mp.Build(roff.NewDocument()), nil
}

// WriteManPage generates and writes a roff man page to the given writer.
func (cli *CLI[T]) WriteManPage(w io.Writer, section uint) error {
	content, err := cli.ManPage(section)
	if err != nil {
		return fmt.Errorf("section=%d: %w", section, err)
	}

	_, err = fmt.Fprint(w, content)

	return err
}

// GenerateManPageCommand creates a cobra command that generates man pages.
// Add this as a subcommand to provide `myapp man` functionality.
func GenerateManPageCommand[T any](cli *CLI[T]) (*cobra.Command, error) {
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

			return cli.WriteManPage(cmd.OutOrStdout(), section)
		},
	}, nil
}

// NewManPage creates a mango man page from a cobra command.
// Useful for custom man page generation pipelines.
func NewManPage(section uint, cmd *cobra.Command) (*mango.ManPage, error) {
	return mcobra.NewManPage(section, cmd)
}
