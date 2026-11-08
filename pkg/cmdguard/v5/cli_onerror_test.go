package v5_test

import (
	"context"
	"errors"
	"testing"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

func TestWithOnError(t *testing.T) {
	t.Parallel()

	t.Run("fires on execution error", func(t *testing.T) {
		t.Parallel()

		var capturedErr error

		cli, err := v5.NewCLI(
			"test", "Test CLI", testCLIConfig{},
			v5.WithFang(false),
			v5.WithOnError(func(err error) {
				capturedErr = err
			}),
		)
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		expectedErr := errors.New("boom")
		cmd, err := v5.NewCommand(
			"fail",
			v5.NoFlags{},
			func(_ context.Context, _ *testCLIConfig, _ v5.NoFlags) error {
				return expectedErr
			},
			v5.WithShort("Always fails"),
		)
		if err != nil {
			t.Fatalf("NewCommand failed: %v", err)
		}

		if err := v5.AddCommand(cli, cmd); err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}

		ctx := t.Context()

		_ = cli.ExecuteWithArgs(ctx, []string{"fail"})

		if capturedErr == nil {
			t.Fatal("WithOnError callback was not invoked")
		}

		if !errors.Is(capturedErr, expectedErr) {
			t.Errorf("callback received wrong error: got %v, want %v", capturedErr, expectedErr)
		}
	})

	t.Run("does not fire on success", func(t *testing.T) {
		t.Parallel()

		called := false

		cli, err := v5.NewCLI(
			"test", "Test CLI", testCLIConfig{},
			v5.WithFang(false),
			v5.WithOnError(func(_ error) {
				called = true
			}),
		)
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cmd, err := v5.NewCommand(
			"ok",
			v5.NoFlags{},
			func(_ context.Context, _ *testCLIConfig, _ v5.NoFlags) error {
				return nil
			},
			v5.WithShort("Always succeeds"),
		)
		if err != nil {
			t.Fatalf("NewCommand failed: %v", err)
		}

		if err := v5.AddCommand(cli, cmd); err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}

		ctx := t.Context()

		err = cli.ExecuteWithArgs(ctx, []string{"ok"})
		if err != nil {
			t.Fatalf("ExecuteWithArgs failed: %v", err)
		}

		if called {
			t.Error("WithOnError should not fire on successful execution")
		}
	})
}
