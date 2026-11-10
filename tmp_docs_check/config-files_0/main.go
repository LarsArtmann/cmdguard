package main

import (
	"context"
	"fmt"
	"log"

	"github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

type DatabaseConfig struct {
	Host     string `flag:"db-host"     env:"APP_DB_HOST"     default:"localhost" help:"Database host"`
	Port     int    `flag:"db-port"     default:"5432"        help:"Database port" validate:"min=1,max=65535"`
	Password string `flag:"db-password" env:"APP_DB_PASSWORD" secret:"true"       help:"Database password"`
}

type AppConfig struct {
	LogLevel string         `flag:"log-level" env:"APP_LOG_LEVEL" default:"info" help:"Log level"`
	Database DatabaseConfig `help:"Database settings"`
	Config   string         `flag:"config" help:"Path to the config file" local:"true"`
}

func main() {
	cli, err := v5.NewCLI[AppConfig]("app", "Config file example", AppConfig{},
		v5.WithConfigFile("$HOME/.config/app/config.yaml", "app.yaml"),
		v5.WithConfigExplainCommand(),
	)
	if err != nil {
		log.Fatal(err)
	}

	b := v5.NewBuilder(cli)
	v5.Add(b, "status", v5.NoFlags{},
		func(_ context.Context, cfg *AppConfig, _ v5.NoFlags) error {
			fmt.Printf("%s:%d at %s\n", cfg.Database.Host, cfg.Database.Port, cfg.LogLevel)
			return nil
		},
		v5.WithShort("Show the resolved database endpoint"),
	)
	if err := b.Build(); err != nil {
		log.Fatal(err)
	}

	cli.ExecuteAndExit(context.Background())
}
