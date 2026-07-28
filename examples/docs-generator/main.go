// Command docs-generator demonstrates cmdguard's GenerateDocs API.
// It builds a small CLI and writes markdown documentation for the full command tree.
package main

import (
	"context"
	"fmt"
	"os"

	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

type config struct {
	Output string `flag:"output" short:"o" default:"markdown" help:"Output format"`
}

type deployFlags struct {
	Environment string `flag:"environment" short:"e" default:"production" help:"Target environment"`
	DryRun      bool   `flag:"dry-run"                                    help:"Print actions without executing"`
}

func main() {
	cli, err := v4.NewCLI(
		"docs-generator",
		"Generate CLI documentation with cmdguard",
		config{},
		v4.WithCLIVersion("0.1.0"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	deployCmd, err := v4.NewCommand(
		"deploy",
		deployFlags{},
		func(_ context.Context, _ *config, flags deployFlags) error {
			fmt.Printf("Deploying to %s (dry-run: %v)\n", flags.Environment, flags.DryRun)
			return nil
		},
		v4.WithShort("Deploy the application"),
		v4.WithLong("Deploy the application to a target environment with optional dry-run mode."),
		v4.WithExample("docs-generator deploy -e staging --dry-run"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	statusCmd, err := v4.NewCommand(
		"status",
		v4.NoFlags{},
		func(_ context.Context, _ *config, _ v4.NoFlags) error {
			fmt.Println("All systems operational")
			return nil
		},
		v4.WithShort("Show deployment status"),
		v4.WithExample("docs-generator status"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := v4.AddCommand(cli, deployCmd); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := v4.AddCommand(cli, statusCmd); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := cli.GenerateDocs(os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error generating docs: %v\n", err)
		os.Exit(1)
	}
}
