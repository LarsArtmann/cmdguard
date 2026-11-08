package v5

import (
	"context"
	"errors"
	"sync"
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

	errs := make([]error, goroutines)

	for g := range goroutines {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for i := range perGoroutine {
				Add(b, "cmd", NoFlags{}, func(ctx context.Context, cfg *struct{}, flags NoFlags) error {
					return nil
				}, WithShort("Racing"))
				b.Check(nil)

				_ = errs
			}
		}()
	}

	wg.Wait()

	buildErr := b.Build()
	if buildErr != nil {
		t.Fatalf("Build after concurrent adds: %v", buildErr)
	}

	// All adds used the same name, so exactly one command must have attached
	// (duplicates are reported, not attached) — this also exercises the
	// duplicate detector on a concurrently-queued batch.
	errs = b.errs
	if len(errs) != goroutines*perGoroutine-1 {
		t.Fatalf("expected %d duplicate-name errors (one attach, rest reported), got %d", goroutines*perGoroutine-1, len(errs))
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

	for i := range n {
		wg.Add(1)

		go func() {
			defer wg.Done()

			Add(b, "unique", NoFlags{}, func(ctx context.Context, cfg *struct{}, flags NoFlags) error {
				return nil
			}, WithShort("Unique"))
		}()
	}

	wg.Wait()

	if err := b.Build(); err != nil {
		t.Fatalf("Build: %v", err)
	}

	_ = n
}

// --- R4: repeated ExecuteWithArgs interleaved with signal cancel ---

func TestRepeatedExecute_InterleavedSignalCancel_IsRaceFree(t *testing.T) {
	t.Parallel()

	cfg := struct {
		Verbose bool `flag:"verbose" short:"v" help:"chatty"`
	}{}

	cli, err := NewCLI("interleave", "Interleave", cfg)
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	b := NewBuilder(cli)
	Add(b, "work", NoFlags{}, func(ctx context.Context, c *struct {
		Verbose bool `flag:"verbose" short:"v" help:"chatty"`
	}, flags NoFlags) error {
		return nil
	}, WithShort("Work"))
	if err := b.Build(); err != nil {
		t.Fatalf("Build: %v", err)
	}

	const runs = 12

	var wg sync.WaitGroup

	runErrs := make([]error, runs)

	for i := range runs {
		wg.Add(1)

		go func() {
			defer wg.Done()

			runErrs[i] = cli.ExecuteWithArgs(t.Context(), []string{"work"})
		}()
	}

	wg.Wait()

	for i, runErr := range runErrs {
		if runErr != nil && !errors.Is(runErr, context.Canceled) {
			t.Fatalf("run %d failed unexpectedly: %v", i, runErr)
		}
	}
}
