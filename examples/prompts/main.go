package main

import (
	"context"
	"fmt"
	"os"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

type AppConfig struct {
	Verbose bool `flag:"verbose" short:"v" default:"false" help:"Enable verbose output"`
}

type DeployFlags struct {
	Environment string `flag:"env"     prompt:"Which environment?" values:"dev,staging,prod" default:"dev"    help:"Target environment"`
	Confirm     bool   `flag:"confirm" prompt:"Are you sure?"                                default:"false"  help:"Confirm deployment"`
	Version     string `flag:"version" prompt:"Version to deploy"                            default:"latest" help:"Version tag"`
}

func main() {
	cli, err := v2.NewCLI[AppConfig]("deployer", "Deployment CLI", AppConfig{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	cmd, err := v2.NewCommand[AppConfig, *DeployFlags]("deploy",
		func(_ context.Context, _ *AppConfig, flags *DeployFlags) error {
			fmt.Printf("Deploying %s to %s (confirmed=%v)\n", flags.Version, flags.Environment, flags.Confirm)

			return nil
		},
		v2.WithShort[AppConfig, *DeployFlags]("Deploy an application"),
		v2.WithFlags[AppConfig, *DeployFlags](&DeployFlags{}),
		v2.WithPromptOnMissing[AppConfig, *DeployFlags](),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := v2.AddCommand(cli, cmd); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if err := cli.Execute(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
