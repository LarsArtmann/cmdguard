// Basic example demonstrating simple cmdguard v2 usage.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/larsartmann/cmdguard/examples/internal"
	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

type AppConfig struct {
	Verbose bool `default:"false" flag:"verbose" help:"Enable verbose output" short:"v"`
}

func main() {
	cli, err := v2.NewCLI[AppConfig]("basic", "A basic CLI example", AppConfig{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	helloCmd, err := internal.NewSimpleCommand[AppConfig]("hello", "Hello, World!", "Say hello")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	goodbyeCmd, err := internal.NewSimpleCommand[AppConfig](
		"goodbye",
		"Goodbye, World!",
		"Say goodbye",
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := v2.AddCommand(cli, helloCmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := v2.AddCommand(cli, goodbyeCmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cli.ExecuteAndExit(context.Background())
}
