package main

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func newTestCLI(t *testing.T) *v2.CLI[AppConfig] {
	t.Helper()

	cli, err := v2.NewCLI[AppConfig]("counting-demo", "Counting flag demo", AppConfig{})
	if err != nil {
		t.Fatalf("failed to create CLI: %v", err)
	}

	return cli
}

func TestCountingExample_Creation(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)

	cmd, err := v2.NewCommand[AppConfig, *GreetFlags](
		"greet",
		func(_ context.Context, _ *AppConfig, _ *GreetFlags) error { return nil },
		v2.WithShort[AppConfig, *GreetFlags]("Greet someone"),
		v2.WithFlags[AppConfig, *GreetFlags](&GreetFlags{}),
	)
	if err != nil {
		t.Fatalf("failed to create greet command: %v", err)
	}

	if err := v2.AddCommand(cli, cmd); err != nil {
		t.Fatalf("failed to add command: %v", err)
	}
}

func TestCountingExample_ExecuteHelp(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	_ = cli.ExecuteWithArgs(context.Background(), []string{"--help"})
}
