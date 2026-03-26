// Package integration provides end-to-end tests for cmdguard.
package integration

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// assertCommandExecution runs a command multiple times and verifies the execution state.
func assertCommandExecution[
	T any,
	F any,
](
	t *testing.T,
	cli *v2.GuardedCommand[T, F],
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
	ctx := t.Context()

	cli, err := v2.New[RootConfig, *GreetFlags]("testapp", "Test application", RootConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = cli.AddCommand(v2.Command[RootConfig, *GreetFlags]{
		Use:   "cmd-a",
		Short: "Command A",
		Flags: &GreetFlags{Name: "default", Shout: false},
		RunE: func(_ context.Context, _ *RootConfig, flags *GreetFlags) error {
			lastExecuted = "A"
			lastFlags = flags

			return nil
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = v2.AddAnyCommand(cli,
		v2.Command[RootConfig, *MathFlags]{
			Use:   "cmd-b",
			Short: "Command B",
			Flags: &MathFlags{X: 0, Y: 0},
			RunE: func(_ context.Context, _ *RootConfig, flags *MathFlags) error {
				lastExecuted = "B"
				lastFlags = flags

				return nil
			},
		},
	)
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
		err = cli.ExecuteWithArgs(ctx, []string{"cmd-a", "--shout"})
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
	ctx := t.Context()

	cli, err := v2.New[RootConfig, v2.NoFlags]("testapp", "Test application", RootConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var executed bool

	err = v2.AddAnyCommand(cli,
		v2.Command[RootConfig, v2.NoFlags]{
			Use:   "simple",
			Short: "Simple command",
			RunE: func(_ context.Context, _ *RootConfig, _ v2.NoFlags) error {
				executed = true

				return nil
			},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = v2.AddAnyCommand(cli,
		v2.Command[RootConfig, *GreetFlags]{
			Use:   "greet",
			Short: "Greet command",
			Flags: &GreetFlags{Name: "World", Shout: false},
			RunE: func(_ context.Context, _ *RootConfig, _ *GreetFlags) error {
				executed = true

				return nil
			},
		},
	)
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
