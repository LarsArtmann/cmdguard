// Package main demonstrates rich output formatting with cmdguard v2.
//
// This example shows:
// - OutputTable for quick table rendering
// - OutputResult with typed TableData
// - Multiple output formats: table, json, csv, yaml
// - OutputStyledTable for terminal-pretty output
//
// Usage:
//
//	go run examples/output/main.go users
//	go run examples/output/main.go users --format json
//	go run examples/output/main.go users --format csv
//	go run examples/output/main.go users --format yaml
//	go run examples/output/main.go styled
package main

import (
	"context"
	"fmt"
	"os"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// AppConfig is the application-level configuration.
type AppConfig struct {
	Debug bool `default:"false" flag:"debug" help:"Enable debug mode" short:"d"`
}

// ListFlags contains flags for the list command.
type ListFlags struct {
	Format string `flag:"format" short:"f" default:"table" help:"Output format (table, json, csv, yaml)"`
}

func main() {
	cli, err := v2.NewCLI[AppConfig]("output-demo", "Output formatting demo", AppConfig{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	usersCmd, err := v2.NewCommand[AppConfig, *ListFlags]("users",
		func(_ context.Context, _ *AppConfig, flags *ListFlags) error {
			format, err := v2.ParseOutputFormat(flags.Format)
			if err != nil {
				return fmt.Errorf("invalid format %q: %w", flags.Format, err)
			}

			headers := []string{"ID", "Name", "Email", "Role"}
			rows := [][]string{
				{"1", "Alice", "alice@example.com", "admin"},
				{"2", "Bob", "bob@example.com", "editor"},
				{"3", "Charlie", "charlie@example.com", "viewer"},
				{"4", "Diana", "diana@example.com", "editor"},
			}

			return v2.OutputTable(format, headers, rows)
		},
		v2.WithShort[AppConfig, *ListFlags]("List users in various formats"),
		v2.WithFlags[AppConfig, *ListFlags](&ListFlags{}),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := v2.AddCommand(cli, usersCmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	styledCmd, err := v2.NewCommand[AppConfig, v2.NoFlags]("styled",
		func(_ context.Context, _ *AppConfig, _ v2.NoFlags) error {
			return v2.OutputStyledTable(
				[]string{"Name", "Status", "Uptime"},
				[][]string{
					{"web-01", "healthy", "99.9%"},
					{"web-02", "healthy", "99.7%"},
					{"db-01", "degraded", "95.2%"},
				},
			)
		},
		v2.WithShort[AppConfig, v2.NoFlags]("Show styled table output"),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := v2.AddCommand(cli, styledCmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cli.ExecuteAndExit(context.Background())
}
