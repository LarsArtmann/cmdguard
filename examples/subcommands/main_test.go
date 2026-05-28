package main

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func TestSubcommandsExample_Creation(t *testing.T) {
	t.Parallel()

	cli, err := v2.NewCLI[AppConfig]("subcmd-demo", "Subcommand hierarchy demo", AppConfig{})
	if err != nil {
		t.Fatalf("failed to create CLI: %v", err)
	}

	upCmd, err := v2.NewCommand[AppConfig, *DBFlags](
		"up",
		func(_ context.Context, _ *AppConfig, _ *DBFlags) error { return nil },
		v2.WithShort[AppConfig, *DBFlags]("Apply pending migrations"),
		v2.WithFlags[AppConfig, *DBFlags](&DBFlags{}),
	)
	if err != nil {
		t.Fatalf("failed to create up command: %v", err)
	}

	downCmd, err := v2.NewCommand[AppConfig, *DBFlags](
		"down",
		func(_ context.Context, _ *AppConfig, _ *DBFlags) error { return nil },
		v2.WithShort[AppConfig, *DBFlags]("Rollback last migration"),
		v2.WithFlags[AppConfig, *DBFlags](&DBFlags{}),
	)
	if err != nil {
		t.Fatalf("failed to create down command: %v", err)
	}

	migrateCmd, err := v2.NewParentCommand[AppConfig, *DBFlags](
		"migrate",
		"Database migration commands",
		[]v2.Command[AppConfig, *DBFlags]{upCmd, downCmd},
		v2.WithShort[AppConfig, *DBFlags]("Database migrations"),
	)
	if err != nil {
		t.Fatalf("failed to create migrate parent: %v", err)
	}

	if err := v2.AddCommand(cli, migrateCmd); err != nil {
		t.Fatalf("failed to add migrate command: %v", err)
	}
}

func TestSubcommandsExample_Help(t *testing.T) {
	t.Parallel()

	cli, _ := v2.NewCLI[AppConfig]("subcmd-demo", "Subcommand hierarchy demo", AppConfig{})
	_ = cli.ExecuteWithArgs(context.Background(), []string{"--help"})
}
