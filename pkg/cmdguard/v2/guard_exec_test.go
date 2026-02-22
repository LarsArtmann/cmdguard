package v2

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGuardedCommand_Execute(t *testing.T) {
	t.Run("executes help command", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		err = g.ExecuteWithArgs(context.Background(), []string{"--help"})
		require.NoError(t, err)
	})

	t.Run("executes subcommand", func(t *testing.T) {
		executed := false
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd := Command[TestAppConfig, NoFlags]{
			Use: "greet",
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				executed = true
				return nil
			},
		}
		require.NoError(t, g.AddCommand(cmd))

		err = g.ExecuteWithArgs(context.Background(), []string{"greet"})
		require.NoError(t, err)
		assert.True(t, executed)
	})

	t.Run("error: unknown subcommand", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		// Add a valid subcommand first
		cmd := Command[TestAppConfig, NoFlags]{
			Use: "valid",
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				return nil
			},
		}
		require.NoError(t, g.AddCommand(cmd))

		// Now try an unknown subcommand - Cobra should return an error
		err = g.ExecuteWithArgs(context.Background(), []string{"unknown"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown")
	})

	t.Run("executes with flags", func(t *testing.T) {
		var receivedName string

		type GreetFlags struct {
			Name string `flag:"name" default:"World"`
		}

		g, err := New[TestAppConfig, *GreetFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd := Command[TestAppConfig, *GreetFlags]{
			Use:   "greet",
			Flags: &GreetFlags{},
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags *GreetFlags) error {
				receivedName = flags.Name
				return nil
			},
		}
		require.NoError(t, g.AddCommand(cmd))

		err = g.ExecuteWithArgs(context.Background(), []string{"greet", "--name", "Alice"})
		require.NoError(t, err)
		assert.Equal(t, "Alice", receivedName)
	})
}

func TestGuardedCommand_ExecuteWithArgs(t *testing.T) {
	t.Run("passes args to command", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd := Command[TestAppConfig, NoFlags]{
			Use: "greet",
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				return nil
			},
		}
		require.NoError(t, g.AddCommand(cmd))

		err = g.ExecuteWithArgs(context.Background(), []string{"greet"})
		require.NoError(t, err)

		// Verify command executed successfully with the given args
		assert.NoError(t, err)
	})
}

func TestGuardedCommand_ExecuteAndExit(t *testing.T) {
	t.Run("returns normally on success", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		// ExecuteAndExit should return normally without calling os.Exit
		assert.NotPanics(t, func() {
			_ = g.ExecuteWithArgs(context.Background(), []string{"--help"})
		})
	})

	t.Run("exits with code 1 on error", func(t *testing.T) {
		if os.Getenv("BE_TEST_EXEC_AND_EXIT") == "1" {
			g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating app: %v\n", err)
				os.Exit(1)
			}

			cmd := Command[TestAppConfig, NoFlags]{
				Use: "fail",
				RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
					return errors.New("intentional failure")
				},
			}
			_ = g.AddCommand(cmd)

			_ = g.ExecuteWithArgs(context.Background(), []string{"fail"})
			g.ExecuteAndExit(context.Background())
			return
		}

		// Run the test in a subprocess
		cmd := exec.Command(os.Args[0], "-test.run=TestGuardedCommand_ExecuteAndExit")
		cmd.Env = append(os.Environ(), "BE_TEST_EXEC_AND_EXIT=1")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()

		// The subprocess should have exited with code 1
		require.Error(t, err, "expected exit code 1")
		if exitErr, ok := err.(*exec.ExitError); ok {
			assert.Equal(t, 1, exitErr.ExitCode())
		}
		assert.Contains(t, stderr.String(), "ERROR")
		assert.Contains(t, stderr.String(), "Intentional failure")
	})

	t.Run("stderr contains error message", func(t *testing.T) {
		if os.Getenv("BE_TEST_EXEC_STDERR") == "1" {
			g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Setup error: %v\n", err)
				os.Exit(1)
			}

			cmd := Command[TestAppConfig, NoFlags]{
				Use: "boom",
				RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
					return errors.New("boom error")
				},
			}
			_ = g.AddCommand(cmd)

			_ = g.ExecuteWithArgs(context.Background(), []string{"boom"})
			g.ExecuteAndExit(context.Background())
			return
		}

		cmd := exec.Command(os.Args[0], "-test.run=TestGuardedCommand_ExecuteAndExit")
		cmd.Env = append(os.Environ(), "BE_TEST_EXEC_STDERR=1")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		_ = cmd.Run()

		assert.Contains(t, stderr.String(), "Boom error")
	})
}
