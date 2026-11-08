package main

import (
	"context"
	"fmt"
	"os"

	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

// Fixture v4 consumer exercising every mechanical codemod rule:
// module path, import alias, v4 qualifier, and the TypeHandler rename
// (registering a custom codec for a domain type).
//
// Lives under an underscore directory so the Go toolchain ignores it in the
// cmdguard workspace (the v4 module no longer exists there); roundtrip.sh
// copies it into a temp module, proves it builds against published v4.2.1,
// runs the codemod, and asserts the result builds clean on v5.0.0.

type AppConfig struct {
	Timeout v4.Duration `flag:"timeout" default:"30s" help:"Request timeout"`
	Verbose bool        `flag:"verbose" short:"v" help:"Chatty output"`
}

type RunFlags struct {
	Name string `flag:"name" short:"n" default:"world" help:"Who to run for"`
}

func main() {
	cli, err := v4.NewCLI("fixture", "Fixture consumer", AppConfig{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	runCmd, err := v4.NewCommand("run", RunFlags{},
		func(ctx context.Context, cfg *AppConfig, f RunFlags) error {
			fmt.Printf("run %s (timeout %s)\n", f.Name, cfg.Timeout.String())
			return nil
		},
		v4.WithShort("Run the fixture"),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := v4.AddCommands(cli, runCmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cli.ExecuteAndExit(context.Background())
}
