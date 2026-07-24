package v3_test

import (
	"context"
	"errors"
	"testing"

	v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

func TestWithOnError(t *testing.T) {
	t.Parallel()

	t.Run("fires on execution error", func(t *testing.T) {
		t.Parallel()

		var capturedErr error

		cli, err := v3.NewCLI("test", "Test CLI", testCLIConfig{},
			v3.WithFang(false),
			v3.WithOnError(func(err error) {
				capturedErr = err
			}),
		)
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		expectedErr := errors.New("boom")
		cmd, err := v3.NewCommand(
			"fail",
			v3.NoFlags{},
			func(_ context.Context, _ *testCLIConfig, _ v3.NoFlags) error {
				return expectedErr
			},
			v3.WithShort("Always fails"),
		)
		if err != nil {
			t.Fatalf("NewCommand failed: %v", err)
		}

		if err := v3.AddCommand(cli, cmd); err != nil {
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

		cli, err := v3.NewCLI("test", "Test CLI", testCLIConfig{},
			v3.WithFang(false),
			v3.WithOnError(func(_ error) {
				called = true
			}),
		)
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		cmd, err := v3.NewCommand(
			"ok",
			v3.NoFlags{},
			func(_ context.Context, _ *testCLIConfig, _ v3.NoFlags) error {
				return nil
			},
			v3.WithShort("Always succeeds"),
		)
		if err != nil {
			t.Fatalf("NewCommand failed: %v", err)
		}

		if err := v3.AddCommand(cli, cmd); err != nil {
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
