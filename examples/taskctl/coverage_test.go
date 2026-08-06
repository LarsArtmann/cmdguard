package main

import (
	"testing"

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

	// seedTasks should silently return when Invoke fails (no TaskStore registered)
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
