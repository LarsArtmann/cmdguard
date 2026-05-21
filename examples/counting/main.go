// Package main demonstrates counting flags with cmdguard v2.
//
// This example shows:
// - Using count:"true" tag for counting flags (-v, -vv, -vvv)
// - Mapping verbosity levels to behavior
// - Counting flags as an alternative to enum-based levels
//
// Usage:
//
//	go run examples/counting/main.go greet --name=Alice
//	go run examples/counting/main.go greet --name=Alice -v
//	go run examples/counting/main.go greet --name=Alice -vvv
package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// AppConfig is the application-level configuration.
type AppConfig struct {
	Output string `default:"text" flag:"output" help:"Output format" short:"o"`
}

// GreetFlags demonstrates a counting flag for verbosity.
type GreetFlags struct {
	Name string `flag:"name"    short:"n" default:"World" help:"Name to greet"`
	Verb int    `flag:"verbose" short:"v" default:"0"     help:"Verbosity level (-v, -vv, -vvv)" count:"true"`
}

func main() {
	cli, err := v2.NewCLI[AppConfig]("counting-demo", "Counting flag demo", AppConfig{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	greetCmd, err := v2.NewCommand[AppConfig, *GreetFlags](
		"greet",
		func(_ context.Context, _ *AppConfig, flags *GreetFlags) error {
			msg := fmt.Sprintf("Hello, %s!", flags.Name)
			fmt.Println(msg)

			switch flags.Verb {
			case 0:
				// silent
			case 1:
				fmt.Println("  [info] greeting delivered")
			case 2:
				fmt.Println("  [info] greeting delivered")
				fmt.Printf("  [debug] name=%q verb=%d\n", flags.Name, flags.Verb)
			default:
				fmt.Println("  [info] greeting delivered")
				fmt.Printf("  [debug] name=%q verb=%d\n", flags.Name, flags.Verb)
				fmt.Printf("  [trace] output_format=%s\n", "text")
				fmt.Printf("  [trace] command=greet args=%v\n", []string{"--name", flags.Name})
			}

			fmt.Printf("\nVerbosity level: %s (%d)\n",
				strings.Repeat("v", flags.Verb), flags.Verb)

			return nil
		},
		v2.WithShort[AppConfig, *GreetFlags]("Greet someone with adjustable verbosity"),
		v2.WithLong[AppConfig, *GreetFlags](
			`Greet someone with adjustable verbosity using counting flags.

Use -v for info, -vv for debug, -vvv for trace-level output.`,
		),
		v2.WithFlags[AppConfig, *GreetFlags](&GreetFlags{}),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := v2.AddCommand(cli, greetCmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cli.ExecuteAndExit(context.Background())
}
