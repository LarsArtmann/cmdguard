package main

import (
	"context"
	"fmt"
	"log"

	"github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
	"github.com/samber/do/v2"
)

type AppConfig struct {
	DSN string `flag:"dsn" env:"APP_DSN" default:"postgres://localhost/app" help:"Database DSN"`
}

type Database struct{ dsn string }

func (db *Database) HealthCheck(_ context.Context) error { return nil }
func (db *Database) Shutdown() error                     { return nil }

func main() {
	cli, err := v5.NewCLI[AppConfig]("app", "DI example", AppConfig{})
	if err != nil {
		log.Fatal(err)
	}

	err = v5.Provide(cli.Scope(), func(_ do.Injector) (*Database, error) {
		cfg, err := v5.Invoke[*AppConfig](cli.Scope())
		if err != nil {
			return nil, err
		}
		return &Database{dsn: cfg.DSN}, nil
	})
	if err != nil {
		log.Fatal(err)
	}

	b := v5.NewBuilder(cli)
	v5.Add(b, "info", v5.NoFlags{},
		func(_ context.Context, _ *AppConfig, _ v5.NoFlags) error {
			db, err := v5.Invoke[*Database](cli.Scope())
			if err != nil {
				return err
			}
			fmt.Println("dsn:", db.dsn)
			return nil
		},
		v5.WithShort("Show the resolved database DSN"),
		v5.WithNoArgs(),
	)
	if err := b.Build(); err != nil {
		log.Fatal(err)
	}

	cli.ExecuteAndExit(context.Background())
}
