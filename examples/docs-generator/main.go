// Command docs-generator demonstrates cmdguard's GenerateDocs API.
// It builds a small CLI and writes markdown documentation for the full command tree.
package main

import (
	"context"
	"fmt"
	"os"

	v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

type config struct {
	Output string `flag:"output" short:"o" default:"markdown" help:"Output format"`
}

type deployFlags struct {
	Environment string `flag:"environment" short:"e" default:"production" help:"Target environment"`
	DryRun      bool   `flag:"dry-run"                                    help:"Print actions without executing"`
}

func main() {
	cli, err := v3.NewCLI(
		"docs-generator",
		"Generate CLI documentation with cmdguard",
		config{},
		v3.WithCLIVersion("0.1.0"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	deployCmd, err := v3.NewCommand(
		"deploy",
		deployFlags{},
		func(_ context.Context, _ *config, flags deployFlags) error {
			fmt.Printf("Deploying to %s (dry-run: %v)\n", flags.Environment, flags.DryRun)
			return nil
		},
		v3.WithShort("Deploy the application"),
		v3.WithLong("Deploy the application to a target environment with optional dry-run mode."),
		v3.WithExample("docs-generator deploy -e staging --dry-run"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	statusCmd, err := v3.NewCommand(
		"status",
		v3.NoFlags{},
		func(_ context.Context, _ *config, _ v3.NoFlags) error {
			fmt.Println("All systems operational")
			return nil
		},
		v3.WithShort("Show deployment status"),
		v3.WithExample("docs-generator status"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := v3.AddCommand(cli, deployCmd); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := v3.AddCommand(cli, statusCmd); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := cli.GenerateDocs(os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error generating docs: %v\n", err)
		os.Exit(1)
	}
}
