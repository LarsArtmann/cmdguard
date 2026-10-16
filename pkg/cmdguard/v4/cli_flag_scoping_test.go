package v4_test

import (
	"context"
	"strings"
	"testing"

	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

// scopingConfig defines one persistent flag and one root-local flag.
type scopingConfig struct {
	Global string `default:"" flag:"global" help:"global persistent flag"`
	Build  string `default:"full" flag:"build" help:"root-only flag" local:"true"`
}

// scopingScanFlags defines a flag that exists only on the scan command.
type scopingScanFlags struct {
	Depth int `default:"1" flag:"depth" help:"per-command flag"`
}

// newScopingCLI builds a CLI whose command tree exercises all three flag scopes:
// root persistent, root local, and per-command.
func newScopingCLI(t *testing.T) (*v4.CLI[scopingConfig], *scopingConfig, *scopingScanFlags) {
	t.Helper()

	cli, err := v4.NewCLI("scoping", "Scoping proof CLI", scopingConfig{}, v4.WithFang(false))
	if err != nil {
		t.Fatalf("NewCLI failed: %v", err)
	}

	var cfg scopingConfig

	var flags scopingScanFlags

	scan, err := v4.NewCommand("scan", scopingScanFlags{},
		func(_ context.Context, c *scopingConfig, f scopingScanFlags) error {
			cfg = *c
			flags = f

			return nil
		},
		v4.WithShort("Scan"),
	)
	if err != nil {
		t.Fatalf("NewCommand scan failed: %v", err)
	}

	noop, err := v4.NewCommand("noop", v4.NoFlags{},
		func(_ context.Context, _ *scopingConfig, _ v4.NoFlags) error { return nil },
		v4.WithShort("No-op"),
	)
	if err != nil {
		t.Fatalf("NewCommand noop failed: %v", err)
	}

	if err := v4.AddCommand(cli, scan); err != nil {
		t.Fatalf("AddCommand scan failed: %v", err)
	}

	if err := v4.AddCommand(cli, noop); err != nil {
		t.Fatalf("AddCommand noop failed: %v", err)
	}

	return cli, &cfg, &flags
}

// TestFlagScopingEndToEnd pins the three flag scopes through real execution:
// per-command F flags are accepted on their own command and rejected on siblings,
// root persistent flags are inherited everywhere, and root-local flags stay on root.
func TestFlagScopingEndToEnd(t *testing.T) {
	t.Parallel()

	t.Run("per-command flag accepted on its own command", func(t *testing.T) {
		t.Parallel()

		cli, cfg, flags := newScopingCLI(t)

		err := cli.ExecuteWithArgs(t.Context(), []string{"scan", "--depth", "7", "--global", "x"})
		if err != nil {
			t.Fatalf("scan --depth 7 --global x failed: %v", err)
		}

		if flags.Depth != 7 {
			t.Errorf("Depth = %d, want 7", flags.Depth)
		}

		if cfg.Global != "x" {
			t.Errorf("Global = %q, want %q (persistent flag inherited by scan)", cfg.Global, "x")
		}
	})

	t.Run("per-command flag rejected on sibling command", func(t *testing.T) {
		t.Parallel()

		cli, _, _ := newScopingCLI(t)

		err := cli.ExecuteWithArgs(t.Context(), []string{"noop", "--depth", "7"})
		if err == nil {
			t.Fatal("noop --depth must fail: depth lives in scan's F struct only")
		}

		if !strings.Contains(err.Error(), "--depth") {
			t.Errorf("error should name --depth, got: %v", err)
		}
	})

	t.Run("root-local flag rejected on subcommand", func(t *testing.T) {
		t.Parallel()

		cli, _, _ := newScopingCLI(t)

		err := cli.ExecuteWithArgs(t.Context(), []string{"scan", "--build", "fast"})
		if err == nil {
			t.Fatal("scan --build must fail: local:\"true\" stays root-only")
		}

		if !strings.Contains(err.Error(), "--build") {
			t.Errorf("error should name --build, got: %v", err)
		}
	})

	t.Run("root-local flag accepted on root", func(t *testing.T) {
		t.Parallel()

		cli, _, _ := newScopingCLI(t)

		// Root-local flags are registered on the root command's own flag set, so
		// the bare root run accepts them even though subcommands reject them.
		// (Values materialize in the config only when a command actually runs:
		// the help path short-circuits before PersistentPreRunE parses flags.)
		err := cli.ExecuteWithArgs(t.Context(), []string{"--build", "fast"})
		if err != nil {
			t.Fatalf("bare root --build fast failed: %v", err)
		}
	})
}
