// Package main is the fixture CLI for DescribeContractFromBinary: a minimal
// cmdguard application whose hidden describe subcommand emits the machine
// contract, exercising the process-boundary capture path end to end.
package main

import (
	"context"
	"fmt"
	"os"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

type config struct {
	Verbose int `count:"true" default:"0" env:"VERBOSE" flag:"verbose" help:"Verbosity level"`
}

type greetFlags struct {
	Name string `default:"world" flag:"name" help:"Who to greet"`
}

func main() {
	app, err := v5.NewApp(
		"fixtureapp",
		"Fixture CLI for the binary contract capture test",
		config{Verbose: 0},
		func(cli *v5.CLI[config]) error {
			leaf, cmdErr := v5.NewCommand(
				"greet",
				greetFlags{Name: ""},
				func(_ context.Context, _ *config, _ greetFlags) error { return nil },
				v5.WithShort("Greet someone"),
			)
			if cmdErr != nil {
				return fmt.Errorf("building greet command: %w", cmdErr)
			}

			return v5.AddCommand(cli, leaf)
		},
		v5.WithEnvPrefix("FIXTURE_APP_"),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	app.ExecuteAndExit(context.Background())
}
