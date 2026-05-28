package main

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func TestEnvTagsExample_Creation(t *testing.T) {
	t.Parallel()

	cli, err := v2.NewCLI[AppConfig](
		"env-demo", "Environment variable demo", AppConfig{},
		v2.WithEnvPrefix[AppConfig]("MYAPP_"),
	)
	if err != nil {
		t.Fatalf("failed to create CLI: %v", err)
	}

	cmd, err := v2.NewCommand[AppConfig, *DBFlags](
		"show",
		func(_ context.Context, _ *AppConfig, _ *DBFlags) error { return nil },
		v2.WithShort[AppConfig, *DBFlags]("Show DB config"),
		v2.WithFlags[AppConfig, *DBFlags](&DBFlags{}),
	)
	if err != nil {
		t.Fatalf("failed to create command: %v", err)
	}

	if err := v2.AddCommand(cli, cmd); err != nil {
		t.Fatalf("failed to add command: %v", err)
	}
}

func TestEnvTagsExample_Help(t *testing.T) {
	t.Parallel()

	cli, _ := v2.NewCLI[AppConfig]("env-demo", "Environment variable demo", AppConfig{})
	_ = cli.ExecuteWithArgs(context.Background(), []string{"--help"})
}
