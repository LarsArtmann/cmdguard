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
			Use: "greet",
			RunE: func(_ context.Context, _ *testAppConfig, _ NoFlags) error {
				executed = true

				return nil
			},
		}
		if err := AddCommand(cli, cmd); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

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

		if err := AddCommand(cli, newTestCmd("valid")); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = cli.ExecuteWithArgs(t.Context(), []string{"unknown"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "unknown") {
			t.Errorf("error should contain 'unknown', got %q", err.Error())
		}
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
			Use:   "greet",
			Flags: &greetFlags{},
			RunE: func(_ context.Context, _ *testAppConfig, flags *greetFlags) error {
				receivedName = flags.Name

				return nil
			},
		}
		if err := AddCommand(cli, cmd); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

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

		if err := AddCommand(cli, newTestCmd("greet")); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

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
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("ExecuteAndExit panicked: %v", r)
				}
			}()

			_ = cli.ExecuteWithArgs(t.Context(), []string{"--help"})
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

		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() != 1 {
				t.Errorf("exit code = %d, want 1", exitErr.ExitCode())
			}
		}

		if !strings.Contains(stderr.String(), "ERROR") {
			t.Errorf("stderr should contain 'ERROR', got %q", stderr.String())
		}

		if !strings.Contains(strings.ToLower(stderr.String()), "intentional failure") {
			t.Errorf("stderr should contain 'intentional failure', got %q", stderr.String())
		}
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

		if !strings.Contains(strings.ToLower(stderr.String()), "boom error") {
			t.Errorf("stderr should contain 'boom error', got %q", stderr.String())
		}
	})
}
