// Package integration provides end-to-end tests for cmdguard.
package integration

import (
	"context"
	"testing"

	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

// assertCommandExecution runs a command multiple times and verifies the execution state.
func assertCommandExecution[
	T any,
](
	t *testing.T,
	cli *v4.CLI[T],
	args []string,
	wantExecuted string,
	assertFlags func(t *testing.T, flags any),
) {
	t.Helper()

	ctx := t.Context()

	for i := range 3 {
		err := cli.ExecuteWithArgs(ctx, args)
		if err != nil {
			t.Fatalf("unexpected error on iteration %d: %v", i+1, err)
		}

		if wantExecuted != lastExecuted {
			t.Errorf("iteration %d: lastExecuted = %q, want %q", i, lastExecuted, wantExecuted)
		}

		assertFlags(t, lastFlags)
	}
}

var (
	lastExecuted string
	lastFlags    any
)

func TestV2_MixedFlagTypes_NoInterference(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	cli, err := v4.NewCLI[RootConfig]("testapp", "Test application", RootConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cmdA, err := v4.NewCommand(
		"cmd-a",
		&GreetFlags{},
		func(_ context.Context, _ *RootConfig, flags *GreetFlags) error {
			lastExecuted = "A"
			lastFlags = flags

			return nil
		},
		v4.WithShort("Command A"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = v4.AddCommand(cli, cmdA)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cmdB, err := v4.NewCommand(
		"cmd-b",
		&MathFlags{},
		func(_ context.Context, _ *RootConfig, flags *MathFlags) error {
			lastExecuted = "B"
			lastFlags = flags

			return nil
		},
		v4.WithShort("Command B"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = v4.AddCommand(cli, cmdB)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertCommandExecution(
		t,
		cli,
		[]string{"cmd-a", "--name=test"},
		"A",
		func(t *testing.T, flags any) {
			t.Helper()

			gf, ok := flags.(*GreetFlags)
			if !ok {
				t.Fatalf("expected *GreetFlags, got %T", flags)
			}

			if gf.Name != "test" {
				t.Errorf("gf.Name = %q, want %q", gf.Name, "test")
			}
		},
	)

	assertCommandExecution(
		t,
		cli,
		[]string{"cmd-b", "--x=42"},
		"B",
		func(t *testing.T, flags any) {
			t.Helper()

			mf, ok := flags.(*MathFlags)
			if !ok {
				t.Fatalf("expected *MathFlags, got %T", flags)
			}

			if mf.X != 42 {
				t.Errorf("mf.X = %d, want %d", mf.X, 42)
			}
		},
	)

	for range 5 {
		err = cli.ExecuteWithArgs(ctx, []string{"cmd-a", flagShout})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if lastExecuted != "A" {
			t.Errorf("lastExecuted = %q, want %q", lastExecuted, "A")
		}

		err = cli.ExecuteWithArgs(ctx, []string{"cmd-b", "--y=99"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if lastExecuted != "B" {
			t.Errorf("lastExecuted = %q, want %q", lastExecuted, "B")
		}
	}
}

func TestV2_MixedFlagTypes_WithNoFlags(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	cli, err := v4.NewCLI[RootConfig]("testapp", "Test application", RootConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var executed bool

	simpleCmd, err := v4.NewCommand(
		"simple",
		v4.NoFlags{},
		func(_ context.Context, _ *RootConfig, _ v4.NoFlags) error {
			executed = true

			return nil
		},
		v4.WithShort("Simple command"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = v4.AddCommand(cli, simpleCmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	greetCmd, err := v4.NewCommand(
		"greet",
		&GreetFlags{},
		func(_ context.Context, _ *RootConfig, _ *GreetFlags) error {
			executed = true

			return nil
		},
		v4.WithShort("Greet command"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = v4.AddCommand(cli, greetCmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	executed = false

	err = cli.ExecuteWithArgs(ctx, []string{"simple"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !executed {
		t.Error("executed should be true")
	}

	executed = false

	err = cli.ExecuteWithArgs(ctx, []string{"greet", "--name=Bob"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !executed {
		t.Error("executed should be true")
	}
}
