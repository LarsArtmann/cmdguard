package main

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/samber/do/v2"

	"github.com/larsartmann/cmdguard/flightrecorder"
	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

// newTestRecorder returns a flight recorder scoped to a temp dir with a
// silent logger. Each test owns its recorder and stops it via t.Cleanup,
// keeping the process-wide singleton free for other tests in any order.
func newTestRecorder(t *testing.T) *flightrecorder.Recorder {
	t.Helper()

	rec := flightrecorder.New(flightrecorder.Config{
		MinAge:         1 * time.Second,
		MaxBytes:       1 << 20,
		OutputDir:      t.TempDir(),
		CaptureOnError: true,
		Log:            func(string, ...any) {},
	})
	t.Cleanup(rec.Stop)

	return rec
}

// runProduction mirrors main() end to end: full production composition
// (audit log, middleware, flight recorder, config file, validation),
// execution, audit export, recorder stop. Not parallel: the flight recorder
// is a process-wide singleton and WithGracefulShutdown installs signal
// handlers.
func runProduction(t *testing.T, args []string) error {
	t.Helper()

	t.Chdir(t.TempDir())

	return run(context.Background(), args)
}

//nolint:paralleltest // buildApp wires the process-wide flight recorder singleton
func TestBuildApp_Composition(t *testing.T) {
	cli, err := buildApp(newTestRecorder(t))
	if err != nil {
		t.Fatalf("buildApp: %v", err)
	}

	if cli.AuditLog() == nil {
		t.Error("expected audit log plugin to be wired")
	}

	if name := cli.Name(); name != "taskctl" {
		t.Errorf("cli name: want %q, got %q", "taskctl", name)
	}
}

//nolint:paralleltest // flight recorder singleton
func TestBuildApp_ExecuteList_ExportsAuditLog(t *testing.T) {
	err := runProduction(t, []string{"list", "--format", "json"})
	if err != nil {
		t.Fatalf("list: %v", err)
	}

	if _, statErr := os.Stat("taskctl-audit.html"); statErr != nil {
		t.Errorf("expected taskctl-audit.html after execution: %v", statErr)
	}
}

//nolint:paralleltest // t.Setenv
func TestExportAuditLog_InvalidFormat(t *testing.T) {
	t.Setenv("AUDIT_LOG_FORMAT", "bogus")

	cli, err := buildApp(newTestRecorder(t))
	if err != nil {
		t.Fatalf("buildApp: %v", err)
	}

	if execErr := cli.ExecuteWithArgs(context.Background(), []string{"list"}); execErr != nil {
		t.Fatalf("list: %v", execErr)
	}

	// Must not fail or write anything despite the bogus format.
	exportAuditLog(cli)

	if _, statErr := os.Stat("taskctl-audit.bogus"); statErr == nil {
		t.Error("no export file expected for invalid format")
	}
}

//nolint:paralleltest // t.Setenv
func TestExportAuditLog_JSONFormat(t *testing.T) {
	t.Setenv("AUDIT_LOG_FORMAT", "json")

	err := runProduction(t, []string{"stats"})
	if err != nil {
		t.Fatalf("stats: %v", err)
	}

	if _, statErr := os.Stat("taskctl-audit.json"); statErr != nil {
		t.Errorf("expected taskctl-audit.json: %v", statErr)
	}
}

//nolint:paralleltest // flight recorder singleton
func TestBuildApp_UnknownCommand(t *testing.T) {
	err := runProduction(t, []string{"definitely-not-a-command"})
	if err == nil {
		t.Fatal("expected error for unknown command")
	}
}

// TestBuildApp_StoreResolutionFailure covers the resolveStore error branches
// in every command handler: the overridden provider fails construction, so
// each handler must surface the DI error instead of panicking.
//
//nolint:paralleltest // flight recorder singleton; sequential by design
func TestBuildApp_StoreResolutionFailure(t *testing.T) {
	commands := [][]string{
		{"list"},
		{"add", "--title", "x"},
		{"done", "--id", "1"},
		{"stats"},
		{"inspect", "--id", "1"},
	}

	for _, args := range commands {
		cli, err := buildApp(newTestRecorder(t))
		if err != nil {
			t.Fatalf("buildApp: %v", err)
		}

		err = v4.Override(cli.Scope(), func(do.Injector) (*TaskStore, error) {
			return nil, errors.New("store construction failed")
		})
		if err != nil {
			t.Fatalf("override: %v", err)
		}

		execErr := cli.ExecuteWithArgs(context.Background(), args)
		if execErr == nil {
			t.Fatalf("expected error for %v with failing store", args)
		}
	}
}

//nolint:paralleltest // t.Setenv
func TestExportAuditLog_NoPluginIsNoOp(t *testing.T) {
	t.Setenv("AUDIT_LOG_FORMAT", "json")

	t.Chdir(t.TempDir())

	// A CLI without an audit plugin exports nothing and errors nowhere.
	cli, err := v4.NewCLI[AppConfig]("taskctl", "test", AppConfig{})
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	exportAuditLog(cli)

	if _, statErr := os.Stat("taskctl-audit.json"); statErr == nil {
		t.Error("no export expected without audit plugin")
	}
}

//nolint:paralleltest // t.Chdir
func TestExportAuditLog_WriteFailure(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root ignores directory permissions; read-only dir test is meaningless")
	}

	t.Setenv("AUDIT_LOG_FORMAT", "json")

	dir := t.TempDir()
	if err := os.Chmod(dir, 0o500); err != nil {
		t.Fatalf("chmod: %v", err)
	}

	t.Chdir(dir)

	cli, err := buildApp(newTestRecorder(t))
	if err != nil {
		t.Fatalf("buildApp: %v", err)
	}

	if execErr := cli.ExecuteWithArgs(context.Background(), []string{"list"}); execErr != nil {
		t.Fatalf("list: %v", execErr)
	}

	// Must report the failure without panicking.
	exportAuditLog(cli)
}

func TestNewProductionRecorder_Config(t *testing.T) {
	t.Parallel()

	rec := newProductionRecorder()
	cfg := rec.Config()

	if !cfg.CaptureOnSlow || !cfg.CaptureOnError {
		t.Errorf("expected slow+error capture, got %+v", cfg)
	}

	if cfg.SlowThreshold != 5*time.Second {
		t.Errorf("SlowThreshold: want 5s, got %s", cfg.SlowThreshold)
	}
}

func TestNewTaskStore_NilInjector(t *testing.T) {
	t.Parallel()

	store, err := NewTaskStore(nil)
	if err == nil {
		t.Fatal("expected error for nil injector")
	}

	if store != nil {
		t.Errorf("expected nil store, got %v", store)
	}
}
