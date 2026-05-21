// Package main demonstrates environment variable configuration with cmdguard v2.
//
// This example shows:
// - Using env tags to read defaults from environment variables
// - Priority chain: flag > env > default
// - Using WithEnvPrefix to namespace environment variables
//
// Usage:
//
//	go run examples/env-tags/main.go show
//	DB_HOST=postgres.example.com go run examples/env-tags/main.go show
//	MYAPP_DB_HOST=postgres.example.com go run examples/env-tags/main.go show --host=custom
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

// DBFlags demonstrates env tag usage.
// Each field can be set via: flag (--host), env var (DB_HOST), or default (localhost).
type DBFlags struct {
	Host     string `flag:"host"     env:"DB_HOST"     default:"localhost" help:"Database host"`
	Port     int    `flag:"port"     env:"DB_PORT"     default:"5432"      help:"Database port"`
	User     string `flag:"user"     env:"DB_USER"     default:"postgres"  help:"Database user"`
	Password string `flag:"password" env:"DB_PASSWORD" default:""          help:"Database password"`
	Name     string `flag:"name"     env:"DB_NAME"     default:"mydb"      help:"Database name"`
}

func main() {
	cli, err := v2.NewCLI[AppConfig](
		"env-demo", "Environment variable demo", AppConfig{},
		v2.WithEnvPrefix[AppConfig]("MYAPP_"),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	showCmd, err := v2.NewCommand[AppConfig, *DBFlags](
		"show",
		func(_ context.Context, cfg *AppConfig, flags *DBFlags) error {
			fmt.Println("Database Configuration:")
			fmt.Printf("  Host:     %s\n", flags.Host)
			fmt.Printf("  Port:     %d\n", flags.Port)
			fmt.Printf("  User:     %s\n", flags.User)
			fmt.Printf("  Password: %s\n", maskPassword(flags.Password))
			fmt.Printf("  Name:     %s\n", flags.Name)
			fmt.Printf("  Debug:    %v\n", cfg.Debug)

			return nil
		},
		v2.WithShort[AppConfig, *DBFlags]("Show database config"),
		v2.WithLong[AppConfig, *DBFlags](
			`Shows the resolved database configuration.

Priority order (highest wins):
  1. Explicit flag:     --host=custom
  2. Environment variable: DB_HOST=postgres.example.com
                        (or MYAPP_DB_HOST with env prefix)
  3. Struct tag default:  localhost`,
		),
		v2.WithFlags[AppConfig, *DBFlags](&DBFlags{}),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := v2.AddCommand(cli, showCmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cli.ExecuteAndExit(context.Background())
}

func maskPassword(p string) string {
	if p == "" {
		return "(not set)"
	}

	return "********"
}
