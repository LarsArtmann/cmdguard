package v4

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// VersionCommand creates a typed "version" subcommand that prints the CLI version.
// Returns an error if the CLI has no version set (use WithCLIVersion first).
//
// Usage:
//
//	cli, _ := v4.NewCLI[Config]("myapp", "My app", Config{},
//	    v4.WithCLIVersion("1.0.0"),
//	)
//	v4.AddCommand(cli, v4.VersionCommand(cli))
func VersionCommand[T any](cli *CLI[T]) (Command[T, NoFlags], error) {
	if cli.spec.version == "" {
		return Command[T, NoFlags]{}, fmt.Errorf(
			"%w: version command requires WithCLIVersion",
			ErrMissingVersion,
		)
	}

	appName := cli.spec.name
	appVersion := cli.spec.version

	return NewCommand(
		"version",
		NoFlags{},
		func(ctx context.Context, cfg *T, _ NoFlags) error {
			_, err := fmt.Fprintf(cli.rootCmd.OutOrStdout(), "%s %s\n", appName, appVersion)
			if err != nil {
				return fmt.Errorf("printing version: %w", err)
			}

			return nil
		},
		WithShort("Print version information"),
		WithLong(fmt.Sprintf("Print the version of %s.", appName)),
	)
}

// GenerateVersionCommand creates a raw cobra version command with custom formatting.
// This provides more control over output format and destination writer.
func GenerateVersionCommand[T any](cli *CLI[T], w io.Writer) (*cobra.Command, error) {
	if cli.spec.version == "" {
		return nil, fmt.Errorf("%w: version command requires WithCLIVersion", ErrMissingVersion)
	}

	appName := cli.spec.name
	appVersion := cli.spec.version

	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  fmt.Sprintf("Print the version of %s.", appName),
		RunE: func(_ *cobra.Command, _ []string) error {
			_, err := fmt.Fprintf(w, "%s %s\n", appName, appVersion)
			if err != nil {
				return fmt.Errorf("printing version: %w", err)
			}

			return nil
		},
	}, nil
}
