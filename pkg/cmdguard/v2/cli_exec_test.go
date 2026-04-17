package v2

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

const testTimeout = 5 * time.Second

func TestCLI_Execute(t *testing.T) {
	t.Parallel()
	t.Run("executes help command", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[testAppConfig]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = cli.ExecuteWithArgs(t.Context(), []string{"--help"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("executes subcommand", func(t *testing.T) {
		t.Parallel()

		executed := false

		cli, err := NewCLI[testAppConfig]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		cmd := Command[testAppConfig, NoFlags]{
			use: "greet",
			runE: func(_ context.Context, _ *testAppConfig, _ NoFlags) error {
				executed = true

				return nil
			},
		}
		addCommand(t, cli, cmd)

		err = cli.ExecuteWithArgs(t.Context(), []string{"greet"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !executed {
			t.Error("command was not executed")
		}
	})

	t.Run("error: unknown subcommand", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[testAppConfig]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		addCommand(t, cli, newTestCmd("valid"))

		err = cli.ExecuteWithArgs(t.Context(), []string{"unknown"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		assertErrorContains(t, err, "unknown")
	})

	t.Run("executes with flags", func(t *testing.T) {
		t.Parallel()

		var receivedName string

		type greetFlags struct {
			Name string `default:"World" flag:"name"`
		}

		cli, err := NewCLI[testAppConfig]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		cmd := Command[testAppConfig, *greetFlags]{
			use:   "greet",
			flags: &greetFlags{},
			runE: func(_ context.Context, _ *testAppConfig, flags *greetFlags) error {
				receivedName = flags.Name

				return nil
			},
		}
		addCommand(t, cli, cmd)

		err = cli.ExecuteWithArgs(t.Context(), []string{"greet", "--name", "Alice"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if receivedName != "Alice" {
			t.Errorf("receivedName = %q, want %q", receivedName, "Alice")
		}
	})
}

func TestCLI_ExecuteWithArgs(t *testing.T) {
	t.Parallel()
	t.Run("passes args to command", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[testAppConfig]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		addCommand(t, cli, newTestCmd("greet"))

		err = cli.ExecuteWithArgs(t.Context(), []string{"greet"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

var errTestIntentionalFailure = errors.New("intentional failure")

func runExecuteAndExitSubprocess(envVar, use string, testErr error) bool {
	if os.Getenv(envVar) == "1" {
		cli, err := NewCLI[testAppConfig]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Setup error: %v\n", err)
			os.Exit(1)
		}

		cmd := newTestCmd(use, testErr)
		_ = AddCommand(cli, cmd)

		_ = cli.ExecuteWithArgs(context.Background(), []string{use})
		cli.ExecuteAndExit(context.Background())

		return true
	}

	return false
}

func TestCLI_ExecuteAndExit(t *testing.T) {
	t.Parallel()
	t.Run("returns normally on success", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[testAppConfig]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		func() {
			didPanic := false

			func() {
				defer func() {
					if r := recover(); r != nil {
						didPanic = true
					}
				}()

				_ = cli.ExecuteWithArgs(t.Context(), []string{"--help"})
			}()

			if didPanic {
				t.Error("ExecuteAndExit should not panic")
			}
		}()
	})

	t.Run("exits with code 1 on error", func(t *testing.T) {
		t.Parallel()

		if runExecuteAndExitSubprocess("BE_TEST_EXEC_AND_EXIT", "fail", errTestIntentionalFailure) {
			return
		}

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestCLI_ExecuteAndExit")

		cmd.Env = append(os.Environ(), "BE_TEST_EXEC_AND_EXIT=1")

		var stderr bytes.Buffer

		cmd.Stderr = &stderr

		err := cmd.Run()
		if err == nil {
			t.Fatal("expected error (exit code 1)")
		}

		exitErr := &exec.ExitError{}
		if errors.As(err, &exitErr) {
			if exitErr.ExitCode() != 1 {
				t.Errorf("exit code = %d, want 1", exitErr.ExitCode())
			}
		}

		if !strings.Contains(stderr.String(), "ERROR") {
			t.Errorf("stderr should contain 'ERROR', got %q", stderr.String())
		}

		assertStderrContains(t, stderr.String(), "intentional failure")
	})

	t.Run("stderr contains error message", func(t *testing.T) {
		t.Parallel()

		if runExecuteAndExitSubprocess("BE_TEST_EXEC_STDERR", "boom", errors.New("boom error")) {
			return
		}

		ctx, cancel := context.WithTimeout(t.Context(), testTimeout)
		defer cancel()

		cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=TestCLI_ExecuteAndExit")

		cmd.Env = append(os.Environ(), "BE_TEST_EXEC_STDERR=1")

		var stderr bytes.Buffer

		cmd.Stderr = &stderr
		_ = cmd.Run()

		assertStderrContains(t, stderr.String(), "boom error")
	})
}
