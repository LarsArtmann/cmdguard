package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/samber/do/v2"
	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

// runProduction executes args against the full production composition
// (audit log, middleware, flight recorder, config file, validation) and
// exports the audit log afterwards, mirroring main() end to end.
// Not parallel: the flight recorder is a process-wide singleton and
// WithGracefulShutdown installs signal handlers.
func runProduction(t *testing.T, args ...string) error {
	t.Helper()

	t.Chdir(t.TempDir())

	cli, err := buildApp()
	if err != nil {
		t.Fatalf("buildApp: %v", err)
	}

	execErr := cli.ExecuteWithArgs(context.Background(), args)

	exportAuditLog(cli)

	return execErr
}

//nolint:paralleltest // buildApp wires the process-wide flight recorder singleton
func TestBuildApp_Composition(t *testing.T) {
	cli, err := buildApp()
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
	err := runProduction(t, "list", "--format", "json")
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

	cli, err := buildApp()
	if err != nil {
		t.Fatalf("buildApp: %v", err)
	}

	if execErr := cli.ExecuteWithArgs(context.Background(), "list"); execErr != nil {
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

	err := runProduction(t, "stats")
	if err != nil {
		t.Fatalf("stats: %v", err)
	}

	if _, statErr := os.Stat("taskctl-audit.json"); statErr != nil {
		t.Errorf("expected taskctl-audit.json: %v", statErr)
	}
}

//nolint:paralleltest // flight recorder singleton
func TestBuildApp_UnknownCommand(t *testing.T) {
	err := runProduction(t, "definitely-not-a-command")
	if err == nil {
		t.Fatal("expected error for unknown command")
	}
}

// TestBuildApp_StoreResolutionFailure covers the resolveStore error branches
// in every command handler: the overridden provider fails construction, so
// each handler must surface the DI error instead of panicking.
//
//nolint:paralleltest // buildApp wires the process-wide flight recorder singleton
func TestBuildApp_StoreResolutionFailure(t *testing.T) {
	commands := [][]string{
		{"list"},
		{"add", "--title", "x"},
		{"done", "--id", "1"},
		{"stats"},
		{"inspect", "--id", "1"},
		{"db", "status"},
	}

	for _, args := range commands {
		t.Run(args[0], func(t *testing.T) {
			t.Parallel()

			t.Chdir(t.TempDir())

			cli, err := buildApp()
			if err != nil {
				t.Fatalf("buildApp: %v", err)
			}

			err = v4.Override(cli.Scope(), func(i do.Injector) (*TaskStore, error) {
				return nil, errors.New("store construction failed")
			})
			if err != nil {
				t.Fatalf("override: %v", err)
			}

			execErr := cli.ExecuteWithArgs(context.Background(), args...)
			if execErr == nil {
				t.Fatalf("expected error for %v with failing store", args)
			}
		})
	}
}

//nolint:paralleltest // flight recorder singleton
func TestNewTaskStore_NilInjector(t *testing.T) {
	store, err := NewTaskStore(nil)
	if err == nil {
		t.Fatal("expected error for nil injector")
	}

	if store != nil {
		t.Errorf("expected nil store, got %v", store)
	}
}

//nolint:paralleltest // t.Setenv pins HOME so the config path is deterministic
func TestBuildApp_ConfigFileMissingIsSkipped(t *testing.T) {
	t.Setenv("HOME", filepath.Join(t.TempDir(), "home"))

	err := runProduction(t, "list")
	if err != nil {
		t.Fatalf("list with missing config file: %v", err)
	}
}
