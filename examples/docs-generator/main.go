// Command docs-generator demonstrates cmdguard's GenerateDocs API.
// It builds a small CLI and writes markdown documentation for the full command tree.
package main

import (
	"context"
	"fmt"
	"os"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
	"github.com/larsartmann/cmdguard/v5/pkg/version"
)

type config struct {
	Output string `flag:"output" short:"o" default:"markdown" help:"Output format"`
}

type deployFlags struct {
	Environment string `flag:"environment" short:"e" default:"production" help:"Target environment"`
	DryRun      bool   `flag:"dry-run"                                    help:"Print actions without executing"`
}

// newContractCLI builds the full command tree. main() and the describe
// contract snapshot test share it so the golden always matches the shipped
// binary's contract.
func newContractCLI() (*v5.CLI[config], error) {
	cli, err := v5.NewCLI(
		"docs-generator",
		"Generate CLI documentation with cmdguard",
		config{},
		v5.WithCLIVersion(version.Version),
	)
	if err != nil {
		return nil, err
	}

	deployCmd, err := v5.NewCommand(
		"deploy",
		deployFlags{},
		func(_ context.Context, _ *config, flags deployFlags) error {
			fmt.Printf("Deploying to %s (dry-run: %v)\n", flags.Environment, flags.DryRun)
			return nil
		},
		v5.WithShort("Deploy the application"),
		v5.WithLong("Deploy the application to a target environment with optional dry-run mode."),
		v5.WithExample("docs-generator deploy -e staging --dry-run"),
	)
	if err != nil {
		return nil, err
	}

	statusCmd, err := v5.NewCommand(
		"status",
		v5.NoFlags{},
		func(_ context.Context, _ *config, _ v5.NoFlags) error {
			fmt.Println("All systems operational")
			return nil
		},
		v5.WithShort("Show deployment status"),
		v5.WithExample("docs-generator status"),
	)
	if err != nil {
		return nil, err
	}

	if err := v5.AddCommand(cli, deployCmd); err != nil {
		return nil, err
	}

	if err := v5.AddCommand(cli, statusCmd); err != nil {
		return nil, err
	}

	return cli, nil
}

func main() {
	cli, err := newContractCLI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := cli.GenerateDocs(os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "error generating docs: %v\n", err)
		os.Exit(1)
	}
}
