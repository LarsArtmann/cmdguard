package v2

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
//	cli, _ := v2.NewCLI[Config]("myapp", "My app", Config{},
//	    v2.WithCLIVersion[Config]("1.0.0"),
//	)
//	v2.AddCommand(cli, v2.VersionCommand[Config](cli))
func VersionCommand[T any](cli *CLI[T]) (Command[T, NoFlags], error) {
	if cli.version == "" {
		return Command[T, NoFlags]{}, fmt.Errorf(
			"%w: version command requires WithCLIVersion",
			ErrMissingVersion,
		)
	}

	appName := cli.name
	appVersion := cli.version

	return NewCommand[T, NoFlags](
		"version",
		func(ctx context.Context, cfg *T, _ NoFlags) error {
			_, err := fmt.Fprintf(cli.rootCmd.OutOrStdout(), "%s %s\n", appName, appVersion)
			if err != nil {
				return fmt.Errorf("printing version: %w", err)
			}

			return nil
		},
		WithShort[T, NoFlags]("Print version information"),
		WithLong[T, NoFlags](fmt.Sprintf("Print the version of %s.", appName)),
	)
}

// GenerateVersionCommand creates a raw cobra version command with custom formatting.
// This provides more control over output format and destination writer.
func GenerateVersionCommand[T any](cli *CLI[T], w io.Writer) (*cobra.Command, error) {
	if cli.version == "" {
		return nil, fmt.Errorf("%w: version command requires WithCLIVersion", ErrMissingVersion)
	}

	appName := cli.name
	appVersion := cli.version

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
