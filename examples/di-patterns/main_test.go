package main

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func TestDIPatternsExample_Creation(t *testing.T) {
	t.Parallel()

	cli, err := v2.NewCLI[AppConfig](
		"di-patterns", "DI service patterns", AppConfig{},
	)
	if err != nil {
		t.Fatalf("failed to create CLI: %v", err)
	}

	cmd, err := v2.NewCommand[AppConfig, v2.NoFlags]("list",
		func(_ context.Context, _ *AppConfig, _ v2.NoFlags) error { return nil },
		v2.WithShort[AppConfig, v2.NoFlags]("List items"),
	)
	if err != nil {
		t.Fatalf("failed to create command: %v", err)
	}

	if err := v2.AddCommand(cli, cmd); err != nil {
		t.Fatalf("failed to add command: %v", err)
	}
}

func TestDIPatternsExample_Help(t *testing.T) {
	t.Parallel()

	cli, _ := v2.NewCLI[AppConfig]("di-patterns", "DI service patterns", AppConfig{})
	_ = cli.ExecuteWithArgs(context.Background(), []string{"--help"})
}
