package main

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func TestErrorHandlingExample_Creation(t *testing.T) {
	t.Parallel()

	cli, err := v2.NewCLI[AppConfig](
		"error-demo", "Error handling demo", AppConfig{},
	)
	if err != nil {
		t.Fatalf("failed to create CLI: %v", err)
	}

	cmd, err := v2.NewCommand[AppConfig, *FetchFlags](
		"fetch",
		func(_ context.Context, _ *AppConfig, _ *FetchFlags) error { return nil },
		v2.WithShort[AppConfig, *FetchFlags]("Fetch a URL"),
		v2.WithFlags[AppConfig, *FetchFlags](&FetchFlags{}),
	)
	if err != nil {
		t.Fatalf("failed to create command: %v", err)
	}

	if err := v2.AddCommand(cli, cmd); err != nil {
		t.Fatalf("failed to add command: %v", err)
	}
}

func TestErrorHandlingExample_Help(t *testing.T) {
	t.Parallel()

	cli, _ := v2.NewCLI[AppConfig]("error-demo", "Error handling demo", AppConfig{})
	_ = cli.ExecuteWithArgs(context.Background(), []string{"--help"})
}
