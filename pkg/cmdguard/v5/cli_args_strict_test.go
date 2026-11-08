package v5

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

// --- R4: strict positional args (opt-in today; default flip is owner-gated, see ROADMAP) ---

func TestWithStrictPositionalArgs_RejectsUndeclaredPositionals(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI("strictargs", "Strict args", struct{}{})
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	b := NewBuilder(cli)
	Add(b, "leaf", NoFlags{}, func(ctx context.Context, cfg *struct{}, flags NoFlags) error {
		return nil
	}, WithShort("Leaf"), WithStrictPositionalArgs())
	if err := b.Build(); err != nil {
		t.Fatalf("Build: %v", err)
	}

	err = cli.ExecuteWithArgs(t.Context(), []string{"leaf", "unexpected-positional"})
	if err == nil {
		t.Fatal("leaf command with WithStrictPositionalArgs accepted an undeclared positional, want error")
	}
}

func TestWithStrictPositionalArgs_DoesNotOverrideExplicitValidator(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI("strictargs2", "Strict args two", struct{}{})
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	b := NewBuilder(cli)
	// WithExactArgs(1) wins: WithStrictPositionalArgs (applied first) must be inert.
	Add(b, "one", NoFlags{}, func(ctx context.Context, cfg *struct{}, flags NoFlags) error {
		return nil
	}, WithShort("One"), WithStrictPositionalArgs(), WithExactArgs(1))
	if err := b.Build(); err != nil {
		t.Fatalf("Build: %v", err)
	}

	if err := cli.ExecuteWithArgs(t.Context(), []string{"one", "declared"}); err != nil {
		t.Fatalf("declared positional rejected despite WithExactArgs(1): %v", err)
	}
}

func TestWithArgsPassthrough_AcceptsExtras(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI("passthrough", "Passthrough", struct{}{})
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	b := NewBuilder(cli)
	Add(b, "free", NoFlags{}, func(ctx context.Context, cfg *struct{}, flags NoFlags) error {
		return nil
	}, WithShort("Free"), WithArgsPassthrough())
	if err := b.Build(); err != nil {
		t.Fatalf("Build: %v", err)
	}

	if err := cli.ExecuteWithArgs(t.Context(), []string{"free", "a", "b", "c"}); err != nil {
		t.Fatalf("passthrough command rejected free-form positionals: %v", err)
	}
}

func TestDefaultLeafCommand_StaysPassthrough(t *testing.T) {
	t.Parallel()

	// Pins the CURRENT default (owner-gated flip tracked in ROADMAP): a leaf
	// command without an args validator still accepts extras. When the strict
	// positional default lands, THIS test flips to expecting an error and the
	// matrix test gains the new row.
	cli, err := NewCLI("defaultargs", "Default args", struct{}{})
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	b := NewBuilder(cli)
	Add(b, "quiet", NoFlags{}, func(ctx context.Context, cfg *struct{}, flags NoFlags) error {
		return nil
	}, WithShort("Quiet"))
	if err := b.Build(); err != nil {
		t.Fatalf("Build: %v", err)
	}

	if err := cli.ExecuteWithArgs(t.Context(), []string{"quiet", "extra"}); err != nil {
		t.Fatalf("default leaf command no longer passthrough (unexpected default flip): %v", err)
	}
}

// --- R4: Builder concurrency (parallel Add/AddParent/Try/Check, single Build) ---

func TestBuilder_ConcurrentAdds_AreRaceFree(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI("racebuilder", "Race builder", struct{}{})
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	const goroutines = 16
	const perGoroutine = 8

	b := NewBuilder(cli)

	var wg sync.WaitGroup

	for range goroutines {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for range perGoroutine {
				Add(b, "cmd", NoFlags{}, func(ctx context.Context, cfg *struct{}, flags NoFlags) error {
					return nil
				}, WithShort("Racing"))
				b.Check(nil)
			}
		}()
	}

	wg.Wait()

	// All adds share one name, so exactly one command attaches and Build
	// reports the rest as duplicates — this exercises the duplicate detector
	// on a concurrently-queued batch.
	buildErr := b.Build()
	if buildErr == nil {
		t.Fatal("duplicate-name batch silently built, want duplicate errors")
	}

	if !errors.Is(buildErr, ErrDuplicateCommand) {
		t.Fatalf("Build error does not wrap ErrDuplicateCommand: %v", buildErr)
	}
}

func TestBuilder_ConcurrentAdds_UniqueNames_AllAttach(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI("racebuilder2", "Race builder two", struct{}{})
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	const n = 32

	b := NewBuilder(cli)

	var wg sync.WaitGroup

	var counter atomic.Int32

	for range n {
		wg.Add(1)

		go func() {
			defer wg.Done()

			name := fmt.Sprintf("cmd-%d", counter.Add(1))

			Add(b, name, NoFlags{}, func(ctx context.Context, cfg *struct{}, flags NoFlags) error {
				return nil
			}, WithShort("Unique"))
		}()
	}

	wg.Wait()

	if err := b.Build(); err != nil {
		t.Fatalf("Build: %v", err)
	}

	got := 0

	for _, sub := range cli.rootCmd.Commands() {
		if strings.HasPrefix(sub.Name(), "cmd-") {
			got++
		}
	}

	if got != n {
		t.Fatalf("expected %d attached commands, got %d", n, got)
	}
}

// --- R4: repeated ExecuteWithArgs interleaved with signal cancel ---

func TestRepeatedExecute_InterleavedSignalCancel_IsRaceFree(t *testing.T) {
	t.Parallel()

	type appCfg struct {
		Verbose bool `flag:"verbose" short:"v" help:"chatty"`
	}

	cli, err := NewCLI("interleave", "Interleave", appCfg{})
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	b := NewBuilder(cli)
	Add(b, "work", NoFlags{}, func(ctx context.Context, c *appCfg, flags NoFlags) error {
		return nil
	}, WithShort("Work"))
	if err := b.Build(); err != nil {
		t.Fatalf("Build: %v", err)
	}

	// Contract note: one CLI is single-flight — concurrent Execute on the
	// SAME instance mutates shared flag state by design (resetFlagStateForReuse).
	// The supported reuse pattern is SEQUENTIAL ExecuteWithArgs, interleaved
	// with cancelled runs (the SIGINT path). This test races the detector over
	// that pattern; parallel Executes belong on separate CLI instances (every
	// t.Parallel test here does exactly that).
	const runs = 12

	cancelled := 0

	for i := range runs {
		ctx, cancel := context.WithCancel(t.Context())
		if i%3 == 0 {
			cancel() // simulate the SIGINT path feeding a cancelled run
		}

		runErr := cli.ExecuteWithArgs(ctx, []string{"work"})
		cancel()

		if ctx.Err() != nil {
			cancelled++

			continue // signal-cancelled run: cancellation error is expected
		}

		if runErr != nil {
			t.Fatalf("run %d failed unexpectedly: %v", i, runErr)
		}
	}

	if cancelled == 0 {
		t.Fatal("no cancelled runs exercised the signal path")
	}
}
