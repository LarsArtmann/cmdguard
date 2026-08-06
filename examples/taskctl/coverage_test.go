package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/samber/do/v2"

	"github.com/larsartmann/cmdguard/flightrecorder"
	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

func TestTaskStatusLabel(t *testing.T) {
	t.Parallel()

	if got := taskStatusLabel(true); got != "done" {
		t.Errorf("taskStatusLabel(true) = %q, want %q", got, "done")
	}
	if got := taskStatusLabel(false); got != "pending" {
		t.Errorf("taskStatusLabel(false) = %q, want %q", got, "pending")
	}
}

func TestTaskRow_DoneTrue(t *testing.T) {
	t.Parallel()

	task := Task{ID: 5, Title: "Done task", Priority: "high", Done: true}
	row := task.Row()
	if row[3] != "done" {
		t.Errorf("Row()[3] = %q, want %q for done task", row[3], "done")
	}
}

func TestTaskStore_GetNotFound(t *testing.T) {
	t.Parallel()

	store := &TaskStore{tasks: []Task{}, next: 1}
	store.Add("Task 1", "low")

	if _, ok := store.Get(999); ok {
		t.Error("Get(999) should return false for non-existent task")
	}
}

func TestTaskStore_GetFound(t *testing.T) {
	t.Parallel()

	store := &TaskStore{tasks: []Task{}, next: 1}
	store.Add("Task 1", "low")

	if task, ok := store.Get(1); !ok {
		t.Error("Get(1) should return true for existing task")
	} else if task.Title != "Task 1" {
		t.Errorf("Get(1).Title = %q, want %q", task.Title, "Task 1")
	}
}

func TestTaskStore_ListHideDone(t *testing.T) {
	t.Parallel()

	store := &TaskStore{tasks: []Task{}, next: 1}
	store.Add("Pending task", "low")
	store.Add("Done task", "low")
	_, _ = store.Done(2)

	visible := store.List("", false)
	if len(visible) != 1 {
		t.Fatalf("List(\"\", false) = %d tasks, want 1 (done hidden)", len(visible))
	}
	if visible[0].Title != "Pending task" {
		t.Errorf("visible task = %q, want %q", visible[0].Title, "Pending task")
	}
}

func TestTaskStore_ListFilterPriorityAndHideDone(t *testing.T) {
	t.Parallel()

	store := &TaskStore{tasks: []Task{}, next: 1}
	store.Add("Pending high", "high")
	store.Add("Done high", "high")
	store.Add("Pending low", "low")
	_, _ = store.Done(2)

	highVisible := store.List("high", false)
	if len(highVisible) != 1 {
		t.Errorf("List(\"high\", false) = %d, want 1", len(highVisible))
	}
}

func TestResolveStore_Error(t *testing.T) {
	t.Parallel()

	scope := v4.NewScope("empty")

	_, err := resolveStore(scope)
	if err == nil {
		t.Error("resolveStore on empty scope should return error")
	}
}

func TestSeedTasks_NoStore(t *testing.T) {
	t.Parallel()

	cli, err := v4.NewCLI[AppConfig](
		"taskctl", "test", AppConfig{},
		v4.WithCLIVersion("1.0.0"),
	)
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	seedTasks(cli)
}

func TestNewTaskStore_MissingConfig(t *testing.T) {
	t.Parallel()

	scope := v4.NewScope("test-no-config")

	_, err := NewTaskStore(scope.Injector())
	if err == nil {
		t.Error("NewTaskStore without config should return error")
	}
}

func TestInspectCommand_NoTaskFound(t *testing.T) {
	t.Parallel()

	cli := newEmptyTestCLI(t)

	err := cli.ExecuteWithArgs(context.Background(), []string{"inspect", "1"})
	if err != nil {
		t.Fatalf("inspect on empty store should not error: %v", err)
	}
}

func TestInspectCommand_WithMetadata(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)

	err := cli.ExecuteWithArgs(context.Background(), []string{"inspect", "1", "--metadata"})
	if err != nil {
		t.Fatalf("inspect --metadata should not error: %v", err)
	}
}

func TestNewTaskStore_Success(t *testing.T) {
	t.Parallel()

	scope := v4.NewScope("test-success")

	if err := v4.Provide(scope, func(i do.Injector) (*AppConfig, error) {
		return &AppConfig{}, nil
	}); err != nil {
		t.Fatalf("provide config: %v", err)
	}

	store, err := NewTaskStore(scope.Injector())
	if err != nil {
		t.Fatalf("NewTaskStore with config should succeed: %v", err)
	}
	if store == nil {
		t.Fatal("store should not be nil")
	}
}

//nolint:paralleltest // FR is a process-wide singleton
func TestFR_Integration_CaptureOnCommandError(t *testing.T) {
	dir := t.TempDir()

	rec := flightrecorder.New(flightrecorder.Config{
		MinAge:         1 * time.Second,
		MaxBytes:       1 << 20,
		OutputDir:      dir,
		CaptureOnError: true,
		Log:            func(string, ...any) {},
	})

	defer rec.Stop()

	if err := rec.Start(); err != nil {
		t.Fatalf("rec.Start: %v", err)
	}

	cli, err := v4.NewCLI[AppConfig](
		"taskctl", "test", AppConfig{},
		v4.WithCLIVersion("1.0.0"),
		v4.WithGroup("tasks", "Task Management"),
		v4.WithGroup("system", "System"),
		flightrecorder.WithFlightRecorderRecorder[AppConfig](rec),
	)
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	if err := v4.Provide(cli.Scope(), NewTaskStore); err != nil {
		t.Fatalf("provide TaskStore: %v", err)
	}

	if err := buildCommands(cli); err != nil {
		t.Fatalf("buildCommands: %v", err)
	}

	_ = cli.ExecuteWithArgs(context.Background(), []string{"done", "--id=999"})

	path := waitTrace(t, dir)
	if path == "" {
		t.Fatal("expected a trace snapshot file after command error")
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat snapshot: %v", err)
	}

	if info.Size() == 0 {
		t.Error("snapshot file is empty")
	}

	base := filepath.Base(path)
	if !strings.Contains(base, "-error-") {
		t.Errorf("filename %q should contain '-error-' reason", base)
	}
}

//nolint:paralleltest // FR is a process-wide singleton
func TestFR_Integration_NoCaptureOnSuccess(t *testing.T) {
	dir := t.TempDir()

	rec := flightrecorder.New(flightrecorder.Config{
		MinAge:         1 * time.Second,
		MaxBytes:       1 << 20,
		OutputDir:      dir,
		CaptureOnSlow:  true,
		SlowThreshold:  1 * time.Hour,
		CaptureOnError: true,
		Log:            func(string, ...any) {},
	})

	defer rec.Stop()

	cli, err := v4.NewCLI[AppConfig](
		"taskctl", "test", AppConfig{},
		v4.WithCLIVersion("1.0.0"),
		v4.WithGroup("tasks", "Task Management"),
		v4.WithGroup("system", "System"),
		flightrecorder.WithFlightRecorderRecorder[AppConfig](rec),
	)
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	if err := v4.Provide(cli.Scope(), NewTaskStore); err != nil {
		t.Fatalf("provide TaskStore: %v", err)
	}

	if err := buildCommands(cli); err != nil {
		t.Fatalf("buildCommands: %v", err)
	}

	if err := cli.ExecuteWithArgs(context.Background(), []string{"list"}); err != nil {
		t.Fatalf("ExecuteWithArgs list: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Errorf("expected no snapshot files, got %d", len(entries))
	}
}

func waitTrace(t *testing.T, dir string) string {
	t.Helper()

	deadline := time.Now().Add(500 * time.Millisecond)

	for time.Now().Before(deadline) {
		entries, _ := os.ReadDir(dir)
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".trace") {
				path := filepath.Join(dir, e.Name())
				if info, err := os.Stat(path); err == nil && info.Size() > 0 {
					return path
				}
			}
		}

		time.Sleep(10 * time.Millisecond)
	}

	return ""
}
