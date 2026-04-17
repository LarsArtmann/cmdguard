// Integration test for basic example
package main

import (
	"bytes"
	"context"
	"os"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func newTestCLI(t *testing.T) *v2.CLI[AppConfig] {
	t.Helper()

	cli, err := v2.NewCLI[AppConfig]("basic", "A basic CLI example", AppConfig{})
	if err != nil {
		t.Fatalf("failed to create CLI: %v", err)
	}

	return cli
}

func addHelloCommand(t *testing.T, cli *v2.CLI[AppConfig], output *bytes.Buffer) {
	t.Helper()

	cmd, err := v2.NewCommand[AppConfig, v2.NoFlags](
		"hello",
		func(_ context.Context, _ *AppConfig, _ v2.NoFlags) error {
			if output != nil {
				output.WriteString("Hello, World!\n")
			}

			return nil
		},
		v2.WithShort[AppConfig, v2.NoFlags]("Say hello"),
	)
	if err != nil {
		t.Fatalf("failed to create hello command: %v", err)
	}

	v2.AddCommand(cli, cmd)
}

func addGoodbyeCommand(t *testing.T, cli *v2.CLI[AppConfig], output *bytes.Buffer) {
	t.Helper()

	cmd, err := v2.NewCommand[AppConfig, v2.NoFlags](
		"goodbye",
		func(_ context.Context, _ *AppConfig, _ v2.NoFlags) error {
			if output != nil {
				output.WriteString("Goodbye, World!\n")
			}

			return nil
		},
		v2.WithShort[AppConfig, v2.NoFlags]("Say goodbye"),
	)
	if err != nil {
		t.Fatalf("failed to create goodbye command: %v", err)
	}

	v2.AddCommand(cli, cmd)
}

func TestBasicExample_HelloCommand(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)

	var output bytes.Buffer

	addHelloCommand(t, cli, &output)
	addGoodbyeCommand(t, cli, &output)

	if err := cli.ExecuteWithArgs(context.Background(), []string{"hello"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.String() != "Hello, World!\n" {
		t.Errorf("output = %q, want %q", output.String(), "Hello, World!\n")
	}

	output.Reset()

	if err := cli.ExecuteWithArgs(context.Background(), []string{"goodbye"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.String() != "Goodbye, World!\n" {
		t.Errorf("output = %q, want %q", output.String(), "Goodbye, World!\n")
	}
}

func TestBasicExample_RootHasSubcommands(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	addHelloCommand(t, cli, nil)
	addGoodbyeCommand(t, cli, nil)

	root := cli.RootCommand()
	if root.Use != "basic" {
		t.Errorf("cmd.Use = %q, want %q", root.Use, "basic")
	}

	if root.Short != "A basic CLI example" {
		t.Errorf("cmd.Short = %q, want %q", root.Short, "A basic CLI example")
	}

	subcmds := root.Commands()
	if len(subcmds) < 2 {
		t.Errorf("len(cmd.Commands()) = %d, want at least 2", len(subcmds))
	}
}

func TestBasicExample_HelpOutput(t *testing.T) {
	t.Parallel()

	if os.Getenv("CI") == "true" {
		t.Skip("Skipping help output test in CI")
	}

	cli := newTestCLI(t)
	addHelloCommand(t, cli, nil)
	addGoodbyeCommand(t, cli, nil)

	_ = cli.ExecuteWithArgs(context.Background(), []string{"--help"})
}
