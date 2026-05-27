package main

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func TestKitchenSinkExample_Creation(t *testing.T) {
	t.Parallel()

	cli, err := v2.NewCLI[AppConfig](
		"taskctl", "A production-grade task manager CLI", AppConfig{},
		v2.WithCLIVersion[AppConfig]("1.0.0"),
	)
	if err != nil {
		t.Fatalf("failed to create CLI: %v", err)
	}

	cmd, err := v2.NewCommand[AppConfig, *ListFlags]("list",
		func(_ context.Context, _ *AppConfig, _ *ListFlags) error { return nil },
		v2.WithShort[AppConfig, *ListFlags]("List tasks"),
		v2.WithFlags[AppConfig, *ListFlags](&ListFlags{}),
	)
	if err != nil {
		t.Fatalf("failed to create command: %v", err)
	}

	if err := v2.AddCommand(cli, cmd); err != nil {
		t.Fatalf("failed to add command: %v", err)
	}
}

func TestKitchenSinkExample_Help(t *testing.T) {
	t.Parallel()

	cli, _ := v2.NewCLI[AppConfig](
		"taskctl", "A production-grade task manager CLI", AppConfig{},
	)
	_ = cli.ExecuteWithArgs(context.Background(), []string{"--help"})
}
