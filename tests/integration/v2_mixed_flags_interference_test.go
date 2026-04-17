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
](
	t *testing.T,
	cli *v2.CLI[T],
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

	cli, err := v2.NewCLI[RootConfig]("testapp", "Test application", RootConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cmdA, err := v2.NewCommand[RootConfig, *GreetFlags]("cmd-a",
		func(_ context.Context, _ *RootConfig, flags *GreetFlags) error {
			lastExecuted = "A"
			lastFlags = flags

			return nil
		},
		v2.WithShort[RootConfig, *GreetFlags]("Command A"),
		v2.WithFlags[RootConfig, *GreetFlags](&GreetFlags{Name: "default", Shout: false}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = v2.AddCommand(cli, cmdA)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cmdB, err := v2.NewCommand[RootConfig, *MathFlags]("cmd-b",
		func(_ context.Context, _ *RootConfig, flags *MathFlags) error {
			lastExecuted = "B"
			lastFlags = flags

			return nil
		},
		v2.WithShort[RootConfig, *MathFlags]("Command B"),
		v2.WithFlags[RootConfig, *MathFlags](&MathFlags{X: 0, Y: 0}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = v2.AddCommand(cli, cmdB)
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
	t.Parallel()
	ctx := t.Context()

	cli, err := v2.NewCLI[RootConfig]("testapp", "Test application", RootConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var executed bool

	simpleCmd, err := v2.NewCommand[RootConfig, v2.NoFlags]("simple",
		func(_ context.Context, _ *RootConfig, _ v2.NoFlags) error {
			executed = true

			return nil
		},
		v2.WithShort[RootConfig, v2.NoFlags]("Simple command"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = v2.AddCommand(cli, simpleCmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	greetCmd, err := v2.NewCommand[RootConfig, *GreetFlags]("greet",
		func(_ context.Context, _ *RootConfig, _ *GreetFlags) error {
			executed = true

			return nil
		},
		v2.WithShort[RootConfig, *GreetFlags]("Greet command"),
		v2.WithFlags[RootConfig, *GreetFlags](&GreetFlags{Name: "World", Shout: false}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = v2.AddCommand(cli, greetCmd)
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
