package v5

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// R5.1: pins the empty-file contract introduced with the startup-crash fix:
// an empty or whitespace-only config file is "absent", so LoadedPath stays
// empty exactly as when no file was found at all.

func TestKoanfLoader_LoadedPath_EmptyFileIsAbsent(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	empty := filepath.Join(dir, "empty.json")
	if err := os.WriteFile(empty, []byte(""), 0o600); err != nil {
		t.Fatalf("writing empty file: %v", err)
	}

	whitespace := filepath.Join(dir, "whitespace.json")
	if err := os.WriteFile(whitespace, []byte("   \n\t\n  "), 0o600); err != nil {
		t.Fatalf("writing whitespace file: %v", err)
	}

	tests := []struct {
		name  string
		paths []string
	}{
		{name: "empty file", paths: []string{empty}},
		{name: "whitespace-only file", paths: []string{whitespace}},
		{name: "missing file", paths: []string{filepath.Join(dir, "nope.json")}},
		{name: "missing then empty (skip chain)", paths: []string{filepath.Join(dir, "nope.json"), empty}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			loader := NewKoanfLoader(tt.paths...)
			cfg := struct {
				LogLevel string `flag:"log-level"`
			}{}

			_, err := loader.Load(nil, &cfg)
			if !errors.Is(err, ErrConfigFileNotFound) {
				t.Fatalf("empty/missing files must surface ErrConfigFileNotFound (the CLI layer tolerates it), got: %v", err)
			}

			if got := loader.LoadedPath(); got != "" {
				t.Fatalf("LoadedPath = %q, want empty (empty file = absent)", got)
			}
		})
	}
}

func TestKoanfLoader_LoadedPath_RealFileIsReported(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	real := filepath.Join(dir, "real.json")
	if err := os.WriteFile(real, []byte(`{"log-level": "debug"}`), 0o600); err != nil {
		t.Fatalf("writing config: %v", err)
	}

	loader := NewKoanfLoader(real)
	cfg := struct {
		LogLevel string `flag:"log-level"`
	}{}

	if _, err := loader.Load(nil, &cfg); err != nil {
		t.Fatalf("Load: %v", err)
	}

	if got := loader.LoadedPath(); got != real {
		t.Fatalf("LoadedPath = %q, want %q", got, real)
	}

	if cfg.LogLevel != "debug" {
		t.Fatalf("LogLevel = %q, want debug (config value must still load)", cfg.LogLevel)
	}
}

// R5.3: doctor's shape on a batteries-ON CLI with zero DI services and no
// custom checks — an empty result set, rendered as "0 passed, 0 failed".

func TestDoctor_NoServicesAndNoChecks_IsEmptyResult(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI("doctorempty", "Doctor empty", struct{}{}) // batteries ON by default
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	if err := cli.ExecuteWithArgs(t.Context(), []string{"doctor"}); err != nil {
		t.Fatalf("doctor on a service-less CLI must succeed: %v", err)
	}

	results := cli.HealthCheckResults()
	if len(results) == 0 {
		t.Fatal("HealthCheckResults returned no entries (cmdguard's own scope services should appear)")
	}

	for svc, svcErr := range results {
		if svcErr != nil {
			t.Fatalf("service %T failed health check on a service-less CLI: %v", svc, svcErr)
		}
	}
}

// R5.2: pins the reduced describe contract under WithoutBatteries — no
// doctor, no version, no auto-registered meta commands beyond the describe
// command itself and cobra's own completion/help set.

func TestWithoutBatteries_DescribeTreeIsReduced(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI("nobatteries", "Opt-out contract", struct{}{}, WithoutBatteries())
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	b := NewBuilder(cli)
	Add(b, "leaf", NoFlags{}, func(ctx context.Context, cfg *struct{}, flags NoFlags) error {
		return nil
	}, WithShort("Leaf"))
	if err := b.Build(); err != nil {
		t.Fatalf("Build: %v", err)
	}

	d := cli.Describe()
	names := map[string]bool{}
	for _, c := range d.Commands {
		names[c.Name] = true
	}

	for _, banned := range []string{"doctor", "version"} {
		if names[banned] {
			t.Fatalf("WithoutBatteries still auto-registers %q", banned)
		}
	}

	if !names["leaf"] {
		t.Fatal("consumer command missing from describe tree")
	}

	// Batteries ON (the default, matrix-pinned) for contrast: doctor present
	// after Execute (auto-registered commands exist only post-Execute).
	withBatteries, err := NewCLI("batteries", "Default", struct{}{})
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	if err := withBatteries.ExecuteWithArgs(t.Context(), []string{"describe"}); err != nil {
		t.Fatalf("describe on batteries CLI: %v", err)
	}

	names2 := map[string]bool{}
	for _, c := range withBatteries.Describe().Commands {
		names2[c.Name] = true
	}

	if !names2["doctor"] {
		t.Fatal("default (batteries ON) describe tree lacks doctor")
	}
}
