// Command config-explain demonstrates cmdguard's configuration features:
// config files, strict key checking, value provenance (`config explain`),
// and the secret:"true" redaction marker.
//
// Run it with an explicit config file:
//
//	go run ./examples/config-explain --config config.example.yaml config explain
package main

import (
	"context"
	"fmt"
	"os"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
	"github.com/larsartmann/cmdguard/v5/pkg/version"
)

type databaseConfig struct {
	Host     string `flag:"db-host"     env:"DEMO_DB_HOST"     default:"localhost" help:"Database host"`
	Port     int    `flag:"db-port"     env:"DEMO_DB_PORT"     default:"5432"      help:"Database port"`
	Password string `flag:"db-password" env:"DEMO_DB_PASSWORD"                     help:"Database password" secret:"true"`
}

type config struct {
	LogLevel string         `flag:"log-level" env:"DEMO_LOG_LEVEL" default:"info" help:"Log level"`
	Database databaseConfig `                                                     help:"Database connection"`
	Config   string         `flag:"config"                                        help:"Path to the config file" local:"true"`
}

func main() {
	cli, err := v5.NewCLI[config]("config-explain",
		"Demo: config files, strict keys, and value provenance",
		config{},
		v5.WithCLIVersion(version.Version),
		v5.WithConfigFile("config.example.yaml"),
		v5.WithStrictConfigKeys(),
		v5.WithConfigExplainCommand(),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	greetCmd, err := v5.NewCommand(
		"greet",
		v5.NoFlags{},
		func(_ context.Context, cfg *config, _ v5.NoFlags) error {
			fmt.Printf("greeting with log level %s, database %s:%d\n",
				cfg.LogLevel, cfg.Database.Host, cfg.Database.Port)

			return nil
		},
		v5.WithShort("Greet using the resolved configuration"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := v5.AddCommand(cli, greetCmd); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cli.ExecuteAndExit(context.Background())
}
