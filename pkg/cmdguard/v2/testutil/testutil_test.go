package testutil

import (
	"context"
	"errors"
	"testing"

	v2 "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
)

type testConfig struct {
	Verbose bool `flag:"verbose" short:"v" default:"false" help:"Enable verbose output"`
}

func TestNewTestCLI(t *testing.T) {
	t.Parallel()

	cli, err := v2.NewCLI[testConfig]("testapp", "Test application", testConfig{})
	if err != nil {
		t.Fatalf("failed to create CLI: %v", err)
	}

	tc := NewTestCLI(cli)
	if tc == nil {
		t.Fatal("NewTestCLI returned nil")
	}

	if tc.CLI() != cli {
		t.Error("CLI() returned unexpected value")
	}
}

func TestTestCLI_ExecuteWithArgs(t *testing.T) {
	t.Parallel()

	cli, err := v2.NewCLI[testConfig]("testapp", "Test application", testConfig{})
	if err != nil {
		t.Fatalf("failed to create CLI: %v", err)
	}

	cmd, err := v2.NewCommand[testConfig, v2.NoFlags](
		"hello",
		func(_ context.Context, _ *testConfig, _ v2.NoFlags) error {
			return nil
		},
		v2.WithShort[testConfig, v2.NoFlags]("Say hello"),
	)
	if err != nil {
		t.Fatalf("failed to create command: %v", err)
	}

	if err := v2.AddCommand(cli, cmd); err != nil {
		t.Fatalf("failed to add command: %v", err)
	}

	tc := NewTestCLI(cli)
	result := tc.ExecuteWithArgs(context.Background(), []string{"hello"})

	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}
}

func TestTestCLI_ExecuteWithArgs_Error(t *testing.T) {
	t.Parallel()

	cli, err := v2.NewCLI[testConfig]("testapp", "Test application", testConfig{})
	if err != nil {
		t.Fatalf("failed to create CLI: %v", err)
	}

	cmd, err := v2.NewCommand[testConfig, v2.NoFlags](
		"fail",
		func(_ context.Context, _ *testConfig, _ v2.NoFlags) error {
			return errors.New("intentional failure")
		},
		v2.WithShort[testConfig, v2.NoFlags]("Fail on purpose"),
	)
	if err != nil {
		t.Fatalf("failed to create command: %v", err)
	}

	if err := v2.AddCommand(cli, cmd); err != nil {
		t.Fatalf("failed to add command: %v", err)
	}

	tc := NewTestCLI(cli)
	result := tc.ExecuteWithArgs(context.Background(), []string{"fail"})

	if result.Error == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestTestResult_ExitCode(t *testing.T) {
	t.Parallel()

	t.Run("no error returns 0", func(t *testing.T) {
		t.Parallel()
		result := &TestResult{Error: nil}
		if result.ExitCode() != 0 {
			t.Errorf("ExitCode() = %d, want 0", result.ExitCode())
		}
	})

	t.Run("plain error returns 1", func(t *testing.T) {
		t.Parallel()
		result := &TestResult{Error: errors.New("plain error")}
		if result.ExitCode() != 1 {
			t.Errorf("ExitCode() = %d, want 1", result.ExitCode())
		}
	})

	t.Run("exit error returns custom code", func(t *testing.T) {
		t.Parallel()
		exitErr, _ := v2.NewExitError(42, errors.New("custom exit"))
		result := &TestResult{Error: exitErr}
		if result.ExitCode() != 42 {
			t.Errorf("ExitCode() = %d, want 42", result.ExitCode())
		}
	})
}

func TestTestCLI_HelpCapture(t *testing.T) {
	t.Parallel()

	cli, err := v2.NewCLI[testConfig]("testapp", "Test application", testConfig{})
	if err != nil {
		t.Fatalf("failed to create CLI: %v", err)
	}

	cmd, err := v2.NewCommand[testConfig, v2.NoFlags](
		"cmd",
		func(_ context.Context, _ *testConfig, _ v2.NoFlags) error {
			return nil
		},
		v2.WithShort[testConfig, v2.NoFlags]("A test command"),
	)
	if err != nil {
		t.Fatalf("failed to create command: %v", err)
	}

	if err := v2.AddCommand(cli, cmd); err != nil {
		t.Fatalf("failed to add command: %v", err)
	}

	tc := NewTestCLI(cli)
	result := tc.ExecuteWithArgs(context.Background(), []string{"--help"})

	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}

	if result.Stdout == "" {
		t.Error("expected help output in stdout, got empty string")
	}
}
