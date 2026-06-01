package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func newTestCLI(t *testing.T) *v2.CLI[AppConfig] {
	t.Helper()

	cli, err := v2.NewCLI[AppConfig](
		"taskctl", "A production-grade task manager CLI", AppConfig{},
		v2.WithCLIVersion[AppConfig]("1.0.0"),
		v2.WithStrictValidation[AppConfig](),
		v2.WithGroup[AppConfig]("tasks", "Task Management"),
		v2.WithGroup[AppConfig]("system", "System"),
	)
	if err != nil {
		t.Fatalf("failed to create CLI: %v", err)
	}

	if err := v2.Provide(cli.Scope(), NewTaskStore); err != nil {
		t.Fatalf("failed to register TaskStore: %v", err)
	}

	if err := buildCommands(cli); err != nil {
		t.Fatalf("failed to build commands: %v", err)
	}

	seedTasks(cli)

	return cli
}

// --- CLI Construction ---

func TestCLI_Construction(t *testing.T) {
	t.Parallel()

	cli, err := v2.NewCLI[AppConfig]("taskctl", "test", AppConfig{})
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}
	if cli == nil {
		t.Fatal("cli is nil")
	}
	if cli.Config() == nil {
		t.Fatal("config is nil")
	}
}

func TestCLI_WithOptions(t *testing.T) {
	t.Parallel()

	cli, err := v2.NewCLI[AppConfig](
		"taskctl", "test", AppConfig{},
		v2.WithCLIVersion[AppConfig]("2.0.0"),
		v2.WithEnvPrefix[AppConfig]("TASKCTL_"),
		v2.WithStrictValidation[AppConfig](),
		v2.WithGroup[AppConfig]("tasks", "Tasks"),
		v2.WithGroup[AppConfig]("system", "System"),
	)
	if err != nil {
		t.Fatalf("NewCLI with options: %v", err)
	}

	root := cli.RootCommand()
	if root.Use != "taskctl" {
		t.Errorf("Use = %q, want %q", root.Use, "taskctl")
	}
}

func TestCLI_ConfigDefaults(t *testing.T) {
	t.Parallel()

	cli, err := v2.NewCLI[AppConfig]("taskctl", "test", AppConfig{})
	if err != nil {
		t.Fatalf("NewCLI: %v", err)
	}

	cfg := cli.Config()
	// Defaults from struct tags are applied during flag parsing, not at construction.
	// At construction time, config has Go zero values.
	if cfg.DataDir != "" {
		t.Errorf("DataDir = %q at construction, want empty (zero value)", cfg.DataDir)
	}
	if cfg.Verbose != 0 {
		t.Errorf("Verbose = %d at construction, want 0", cfg.Verbose)
	}
}

// --- DI Services ---

func TestDI_TaskStoreRegistration(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)

	store, err := v2.Invoke[*TaskStore](cli.Scope())
	if err != nil {
		t.Fatalf("Invoke TaskStore: %v", err)
	}
	if store == nil {
		t.Fatal("store is nil")
	}
}

func TestDI_HealthCheck(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)

	if err := cli.HealthCheck(); err != nil {
		t.Errorf("HealthCheck: %v", err)
	}
}

func TestDI_HealthCheckWithContext(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)

	if err := cli.HealthCheckWithContext(context.Background()); err != nil {
		t.Errorf("HealthCheckWithContext: %v", err)
	}
}

func TestDI_Shutdown(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := cli.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown: %v", err)
	}
}

// --- List Command ---

func TestListCommand_Default(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"list"})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
}

func TestListCommand_JSON(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"list", "--format", "json"})
	if err != nil {
		t.Fatalf("list --format json: %v", err)
	}
}

func TestListCommand_CSV(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"list", "--format", "csv"})
	if err != nil {
		t.Fatalf("list --format csv: %v", err)
	}
}

func TestListCommand_YAML(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"list", "--format", "yaml"})
	if err != nil {
		t.Fatalf("list --format yaml: %v", err)
	}
}

func TestListCommand_All(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"list", "--all"})
	if err != nil {
		t.Fatalf("list --all: %v", err)
	}
}

func TestListCommand_FilterPriority(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"list", "--priority", "high"})
	if err != nil {
		t.Fatalf("list --priority high: %v", err)
	}
}

func TestListCommand_InvalidFormat(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"list", "--format", "invalid"})
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestListCommand_AliasLS(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"ls"})
	if err != nil {
		t.Fatalf("ls alias: %v", err)
	}
}

func TestListCommand_Empty(t *testing.T) {
	t.Parallel()

	cli, _ := v2.NewCLI[AppConfig](
		"taskctl", "test", AppConfig{},
		v2.WithCLIVersion[AppConfig]("1.0.0"),
		v2.WithStrictValidation[AppConfig](),
		v2.WithGroup[AppConfig]("tasks", "Task Management"),
		v2.WithGroup[AppConfig]("system", "System"),
	)
	_ = v2.Provide(cli.Scope(), NewTaskStore)
	_ = buildCommands(cli)

	err := cli.ExecuteWithArgs(context.Background(), []string{"list"})
	if err != nil {
		t.Fatalf("list empty: %v", err)
	}
}

// --- Add Command ---

func TestAddCommand_Valid(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"add", "--title", "Write tests"})
	if err != nil {
		t.Fatalf("add: %v", err)
	}
}

func TestAddCommand_WithPriority(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"add", "--title", "Fix bug", "--priority", "high"})
	if err != nil {
		t.Fatalf("add with priority: %v", err)
	}
}

func TestAddCommand_InvalidPriority(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"add", "--title", "Valid title", "--priority", "urgent"})
	if err == nil {
		t.Fatal("expected error for invalid priority")
	}
}

// --- Done Command ---

func TestDoneCommand_Valid(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"done", "--id", "1"})
	if err != nil {
		t.Fatalf("done: %v", err)
	}
}

func TestDoneCommand_NotFound(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"done", "--id", "999"})
	if err == nil {
		t.Fatal("expected error for not-found task")
	}
}

func TestDoneCommand_MissingID(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"done"})
	if err == nil {
		t.Fatal("expected error for missing id")
	}
}

// --- Stats Command ---

func TestStatsCommand_Default(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"stats"})
	if err != nil {
		t.Fatalf("stats: %v", err)
	}
}

func TestStatsCommand_JSON(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"stats", "--format", "json"})
	if err != nil {
		t.Fatalf("stats --format json: %v", err)
	}
}

// --- DB Subcommands (NewParentCommand) ---

func TestDBCommand_Migrate(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"db", "migrate", "--env", "production"})
	if err != nil {
		t.Fatalf("db migrate: %v", err)
	}
}

func TestDBCommand_Seed(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"db", "seed"})
	if err != nil {
		t.Fatalf("db seed: %v", err)
	}
}

func TestDBCommand_Status(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"db", "status"})
	if err != nil {
		t.Fatalf("db status: %v", err)
	}
}

// --- Health Command (MustNewCommand) ---

func TestHealthCommand(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"health"})
	if err != nil {
		t.Fatalf("health: %v", err)
	}
}

// --- Version Command ---

func TestVersionCommand(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"version"})
	if err != nil {
		t.Fatalf("version: %v", err)
	}
}

// --- Config Command ---

func TestConfigShowCommand(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"config", "show"})
	if err != nil {
		t.Fatalf("config show: %v", err)
	}
}

// --- Hidden Command ---

func TestHiddenCommand(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)

	// Hidden command should still execute
	err := cli.ExecuteWithArgs(context.Background(), []string{"secret"})
	if err != nil {
		t.Fatalf("secret: %v", err)
	}

	// But not show in help
	root := cli.RootCommand()
	for _, cmd := range root.Commands() {
		if cmd.Name() == "secret" && !cmd.Hidden {
			t.Error("secret command should be hidden")
		}
	}
}

// --- Deprecated Command ---

func TestDeprecatedCommand(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"complete"})
	if err != nil {
		t.Fatalf("complete: %v", err)
	}
}

// --- Error Handling ---

func TestErrorHandling_FlagError(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"list", "--format", "invalid"})

	if err == nil {
		t.Fatal("expected error")
	}

	if flagErr, ok := errors.AsType[*v2.FlagError](err); ok {
		if flagErr.FlagName != "format" {
			t.Errorf("FlagName = %q, want %q", flagErr.FlagName, "format")
		}
	}
}

func TestErrorHandling_ExitCode(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"done", "--id", "999"})
	if err == nil {
		t.Fatal("expected error")
	}

	if exitCoder, ok := errors.AsType[*v2.ExitError](err); ok {
		if exitCoder.Code != 2 {
			t.Errorf("ExitCode = %d, want 2", exitCoder.Code)
		}
	}
}

// --- Help Output ---

func TestHelpOutput(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"--help"})
	// --help returns a special error in cobra, that's expected
	_ = err
}

// --- Priority Parsing ---

func TestParsePriority(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  string
		err   bool
	}{
		{"low", PriorityLow, false},
		{"medium", PriorityMedium, false},
		{"high", PriorityHigh, false},
		{"urgent", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			got, err := v2.ParseEnum(tt.input, strings.Split(allowedPriorities, ","))
			if tt.err {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.String() != tt.want {
				t.Errorf("ParseEnum(%q) = %q, want %q", tt.input, got.String(), tt.want)
			}
		})
	}
}

// --- ValueOrDefault helper ---

func TestValueOrDefault(t *testing.T) {
	t.Parallel()

	s := "hello"
	if got := v2.ValueOrDefault(&s, "default"); got != "hello" {
		t.Errorf("ValueOrDefault = %q, want %q", got, "hello")
	}
	if got := v2.ValueOrDefault[string](nil, "default"); got != "default" {
		t.Errorf("ValueOrDefault(nil) = %q, want %q", got, "default")
	}
}

// --- EnsureValid helper ---

func TestEnsureValid(t *testing.T) {
	t.Parallel()

	s := "value"
	if err := v2.EnsureValid(&s, "test"); err != nil {
		t.Errorf("EnsureValid with non-nil: %v", err)
	}
	if err := v2.EnsureValid[string](nil, "test"); err == nil {
		t.Error("EnsureValid with nil should error")
	}
}

// --- Integration: Full Workflow ---

func TestIntegration_FullWorkflow(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)

	// List seeded tasks
	if err := cli.ExecuteWithArgs(context.Background(), []string{"list"}); err != nil {
		t.Fatalf("list: %v", err)
	}

	// Add a new task
	if err := cli.ExecuteWithArgs(
		context.Background(),
		[]string{"add", "--title", "Integration test task"},
	); err != nil {
		t.Fatalf("add: %v", err)
	}

	// Complete task #1
	if err := cli.ExecuteWithArgs(context.Background(), []string{"done", "--id", "1"}); err != nil {
		t.Fatalf("done: %v", err)
	}

	// Stats
	if err := cli.ExecuteWithArgs(context.Background(), []string{"stats"}); err != nil {
		t.Fatalf("stats: %v", err)
	}

	// Health check
	if err := cli.ExecuteWithArgs(context.Background(), []string{"health"}); err != nil {
		t.Fatalf("health: %v", err)
	}
}

// --- Inspect Command ---

func TestInspectCommand(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"inspect", "1"})
	if err != nil {
		t.Fatalf("inspect 1: %v", err)
	}
}

func TestInspectCommand_NoArgs(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"inspect"})
	if err == nil {
		t.Fatal("expected error for inspect without args")
	}
}

func TestInspectCommand_TooManyArgs(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"inspect", "1", "2"})
	if err == nil {
		t.Fatal("expected error for inspect with too many args")
	}
}

// --- Root command structure ---

func TestRootCommand_Structure(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	root := cli.RootCommand()

	names := make(map[string]bool)
	for _, cmd := range root.Commands() {
		names[cmd.Name()] = true
	}

	expected := []string{
		"list",
		"add",
		"done",
		"stats",
		"inspect",
		"db",
		"health",
		"config",
		"version",
		"secret",
		"complete",
	}
	for _, name := range expected {
		if !names[name] {
			t.Errorf("missing command: %s", name)
		}
	}
}

// --- TaskStore domain logic ---

func TestTaskStore_AddAndList(t *testing.T) {
	t.Parallel()

	store := &TaskStore{mu: sync.Mutex{}, tasks: []Task{}, next: 1}

	task := store.Add("Test task", PriorityHigh)
	if task.ID != 1 {
		t.Errorf("task.ID = %d, want 1", task.ID)
	}
	if task.Title != "Test task" {
		t.Errorf("task.Title = %q, want %q", task.Title, "Test task")
	}
	if task.Done {
		t.Error("new task should not be done")
	}

	tasks := store.List("", true)
	if len(tasks) != 1 {
		t.Fatalf("len(tasks) = %d, want 1", len(tasks))
	}
}

func TestTaskStore_Done(t *testing.T) {
	t.Parallel()

	store := &TaskStore{mu: sync.Mutex{}, tasks: []Task{}, next: 1}
	store.Add("Task 1", PriorityMedium)

	task, err := store.Done(1)
	if err != nil {
		t.Fatalf("Done: %v", err)
	}
	if !task.Done {
		t.Error("task should be done")
	}
}

func TestTaskStore_DoneNotFound(t *testing.T) {
	t.Parallel()

	store := &TaskStore{mu: sync.Mutex{}, tasks: []Task{}, next: 1}
	_, err := store.Done(999)
	if err == nil {
		t.Fatal("expected error for not-found task")
	}
}

func TestTaskStore_FilterPriority(t *testing.T) {
	t.Parallel()

	store := &TaskStore{mu: sync.Mutex{}, tasks: []Task{}, next: 1}
	store.Add("Low task", PriorityLow)
	store.Add("High task", PriorityHigh)

	high := store.List("high", true)
	if len(high) != 1 {
		t.Errorf("len(high) = %d, want 1", len(high))
	}

	low := store.List("low", true)
	if len(low) != 1 {
		t.Errorf("len(low) = %d, want 1", len(low))
	}
}

func TestTaskStore_Stats(t *testing.T) {
	t.Parallel()

	store := &TaskStore{mu: sync.Mutex{}, tasks: []Task{}, next: 1}
	store.Add("Task 1", PriorityHigh)
	store.Add("Task 2", PriorityLow)
	_, _ = store.Done(1)

	total, pending, done, byPriority := store.Stats()
	if total != 2 {
		t.Errorf("total = %d, want 2", total)
	}
	if pending != 1 {
		t.Errorf("pending = %d, want 1", pending)
	}
	if done != 1 {
		t.Errorf("done = %d, want 1", done)
	}
	if byPriority[PriorityHigh] != 1 {
		t.Errorf("high = %d, want 1", byPriority[PriorityHigh])
	}
}

func TestTaskStore_IDs(t *testing.T) {
	t.Parallel()

	store := &TaskStore{mu: sync.Mutex{}, tasks: []Task{}, next: 1}
	store.Add("Task 1", PriorityMedium)
	store.Add("Task 2", PriorityMedium)

	ids := store.IDs()
	if len(ids) != 2 {
		t.Fatalf("len(ids) = %d, want 2", len(ids))
	}
	if ids[0] != "1" || ids[1] != "2" {
		t.Errorf("ids = %v, want [1 2]", ids)
	}
}

// --- MustParse helper ---

func TestMustParse(t *testing.T) {
	t.Parallel()

	got := v2.MustParse("duration", "5s", v2.ParseDuration)
	if got.Duration().String() != "5s" {
		t.Errorf("MustParse duration = %v, want 5s", got)
	}
}

// --- Duration type ---

func TestDurationType(t *testing.T) {
	t.Parallel()

	d, err := v2.ParseDuration("1h30m")
	if err != nil {
		t.Fatalf("ParseDuration: %v", err)
	}
	if d.Duration().Minutes() != 90 {
		t.Errorf("Duration = %v, want 90m", d.Duration())
	}
}

// --- LogLevel type ---

func TestLogLevelType(t *testing.T) {
	t.Parallel()

	ll, err := v2.ParseLogLevel("debug")
	if err != nil {
		t.Fatalf("ParseLogLevel: %v", err)
	}
	if ll.String() != "debug" {
		t.Errorf("LogLevel = %q, want %q", ll.String(), "debug")
	}
}

// --- Enum type ---

func TestEnumType(t *testing.T) {
	t.Parallel()

	e, err := v2.ParseEnum("high", []string{"low", "medium", "high"})
	if err != nil {
		t.Fatalf("ParseEnum: %v", err)
	}
	if e.String() != "high" {
		t.Errorf("Enum = %q, want %q", e.String(), "high")
	}
}

func TestEnumType_Invalid(t *testing.T) {
	t.Parallel()

	_, err := v2.ParseEnum("urgent", []string{"low", "medium", "high"})
	if err == nil {
		t.Fatal("expected error for invalid enum value")
	}
}

// --- Port type ---

func TestPortType(t *testing.T) {
	t.Parallel()

	p, err := v2.ParsePort("8080")
	if err != nil {
		t.Fatalf("ParsePort: %v", err)
	}
	if p.String() != "8080" {
		t.Errorf("Port = %q, want %q", p.String(), "8080")
	}
}

// --- BranchingFlowContext ---

func TestBranchingFlowContext(t *testing.T) {
	t.Parallel()

	cli := newTestCLI(t)
	err := cli.ExecuteWithArgs(context.Background(), []string{"inspect", "1"})
	if err != nil {
		t.Fatalf("inspect: %v", err)
	}
}

// --- NewCommandError / NewServiceError ---

func TestCommandError(t *testing.T) {
	t.Parallel()

	inner := fmt.Errorf("inner error")
	cmdErr := v2.NewCommandError("test", inner)
	if !strings.Contains(cmdErr.Error(), "test") {
		t.Errorf("CommandError should contain command name: %v", cmdErr)
	}
}

func TestServiceError(t *testing.T) {
	t.Parallel()

	inner := fmt.Errorf("inner error")
	svcErr := v2.NewServiceError("*TaskStore", inner)
	if !strings.Contains(svcErr.Error(), "*TaskStore") {
		t.Errorf("ServiceError should contain service type: %v", svcErr)
	}
}

// --- NewExitError ---

func TestNewExitError(t *testing.T) {
	t.Parallel()

	exitErr, err := v2.NewExitError(42, fmt.Errorf("test"))
	if err != nil {
		t.Fatalf("NewExitError: %v", err)
	}
	if exitErr.Code != 42 {
		t.Errorf("ExitCode = %d, want 42", exitErr.Code)
	}
}

func TestNewExitError_InvalidCode(t *testing.T) {
	t.Parallel()

	_, err := v2.NewExitError(300, fmt.Errorf("test"))
	if err == nil {
		t.Fatal("expected error for invalid exit code")
	}
}

// --- FlagError ---

func TestFlagError(t *testing.T) {
	t.Parallel()

	flagErr := v2.NewFlagError("port", fmt.Errorf("out of range"))
	if !strings.Contains(flagErr.Error(), "port") {
		t.Errorf("FlagError should contain flag name: %v", flagErr)
	}
}

func TestFlagError_WithSuggestion(t *testing.T) {
	t.Parallel()

	flagErr := v2.NewFlagErrorWithSuggestion("prot", fmt.Errorf("unknown flag"), "port")
	if flagErr.Suggestion != "port" {
		t.Errorf("Suggestion = %q, want %q", flagErr.Suggestion, "port")
	}
	if !strings.Contains(flagErr.Error(), "prot") {
		t.Errorf("FlagError should contain flag name: %v", flagErr)
	}
}
