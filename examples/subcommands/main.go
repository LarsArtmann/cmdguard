// Package main demonstrates NewParentCommand for command groups with subcommands.
//
// This example shows:
// - Using NewParentCommand to create a parent command with subcommands
// - Subcommands sharing the same flags type (DBFlags for all migrate commands)
// - Using AddCommand for commands with different flag types
// - Building a git-style CLI with nested command groups
//
// Note: NewParentCommand[T, F] requires all direct subcommands to share the
// same F type parameter. This is a Go generics constraint. For heterogeneous
// flag types, use v2.AddCommand to add each command individually.
//
// Usage:
//
//	go run examples/subcommands/main.go migrate up --env=production
//	go run examples/subcommands/main.go migrate down --env=staging --force
//	go run examples/subcommands/main.go migrate status
//	go run examples/subcommands/main.go version
package main

import (
	"context"
	"fmt"
	"os"

	examplesinternal "github.com/larsartmann/cmdguard/examples/internal"
	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// AppConfig is the application-level configuration.
type AppConfig struct {
	Verbose bool `default:"false" flag:"verbose" help:"Enable verbose output" short:"v"`
}

// DBFlags defines shared database flags for all migrate subcommands.
// Because NewParentCommand requires uniform flag types, sharing a single
// flags struct across siblings is the natural pattern.
type DBFlags struct {
	Env   string `flag:"env"   short:"e" help:"Environment (development, staging, production)" default:"development"`
	Force bool   `flag:"force" short:"f" help:"Skip confirmation prompts"                      default:"false"`
}

func main() {
	cli, err := v2.NewCLI[AppConfig]("subcmd-demo", "Subcommand hierarchy demo", AppConfig{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// --- migrate command group (all subcommands share DBFlags) ---
	migrateUpCmd, err := v2.NewCommand[AppConfig, *DBFlags](
		"up",
		func(_ context.Context, cfg *AppConfig, flags *DBFlags) error {
			fmt.Printf("Running migrations UP on %s", flags.Env)

			if flags.Force {
				fmt.Print(" (forced)")
			}

			fmt.Println()

			if cfg.Verbose {
				fmt.Println("[verbose] checking pending migrations...")
			}

			return nil
		},
		v2.WithShort[AppConfig, *DBFlags]("Apply pending migrations"),
		v2.WithFlags[AppConfig, *DBFlags](&DBFlags{}),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	migrateDownCmd, err := v2.NewCommand[AppConfig, *DBFlags](
		"down",
		func(_ context.Context, _ *AppConfig, flags *DBFlags) error {
			fmt.Printf("Rolling back last migration on %s", flags.Env)

			if flags.Force {
				fmt.Print(" (forced)")
			}

			fmt.Println()

			return nil
		},
		v2.WithShort[AppConfig, *DBFlags]("Rollback last migration"),
		v2.WithFlags[AppConfig, *DBFlags](&DBFlags{}),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	migrateStatusCmd, err := v2.NewCommand[AppConfig, *DBFlags](
		"status",
		func(_ context.Context, _ *AppConfig, flags *DBFlags) error {
			fmt.Printf("Migration status for %s: 3 applied, 2 pending\n", flags.Env)

			return nil
		},
		v2.WithShort[AppConfig, *DBFlags]("Show migration status"),
		v2.WithFlags[AppConfig, *DBFlags](&DBFlags{}),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// NewParentCommand creates a parent that groups the subcommands.
	// All subcommands in the slice must have the same [T, F] type parameters.
	migrateCmd, err := v2.NewParentCommand[AppConfig, *DBFlags](
		"migrate",
		"Database migration commands",
		[]v2.Command[AppConfig, *DBFlags]{migrateUpCmd, migrateDownCmd, migrateStatusCmd},
		v2.WithShort[AppConfig, *DBFlags]("Database migrations"),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := v2.AddCommand(cli, migrateCmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// --- standalone command (different flags type, added via AddCommand) ---
	versionCmd, err := examplesinternal.NewSimpleCommand[AppConfig](
		"version", "subcmd-demo v1.0.0", "Show version",
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := v2.AddCommand(cli, versionCmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	examplesinternal.Execute(context.Background(), cli)
}
