package main

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func TestSignalsExample_Creation(t *testing.T) {
	t.Parallel()

	cli, err := v2.NewCLI[AppConfig](
		"signal-demo", "Signal handling demo", AppConfig{},
	)
	if err != nil {
		t.Fatalf("failed to create CLI: %v", err)
	}

	cmd, err := v2.NewCommand[AppConfig, *ServeFlags](
		"serve",
		func(_ context.Context, _ *AppConfig, _ *ServeFlags) error { return nil },
		v2.WithShort[AppConfig, *ServeFlags]("Start server"),
		v2.WithFlags[AppConfig, *ServeFlags](&ServeFlags{}),
	)
	if err != nil {
		t.Fatalf("failed to create command: %v", err)
	}

	if err := v2.AddCommand(cli, cmd); err != nil {
		t.Fatalf("failed to add command: %v", err)
	}
}

func TestSignalsExample_Help(t *testing.T) {
	t.Parallel()

	cli, _ := v2.NewCLI[AppConfig]("signal-demo", "Signal handling demo", AppConfig{})
	_ = cli.ExecuteWithArgs(context.Background(), []string{"--help"})
}
