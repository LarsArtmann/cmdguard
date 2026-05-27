// Package main demonstrates a production-grade CLI using every major cmdguard feature.
//
// This "task manager" CLI shows:
//   - Type-safe config with env var bindings
//   - Dependency injection with lifecycle hooks (HealthCheck, Shutdown)
//   - Per-command typed flags (WithFlags)
//   - PreRunE validation and PostRunE cleanup
//   - Middleware (timing + recovery)
//   - Rich output in multiple formats (OutputTable)
//   - Command groups (WithGroup)
//   - Error handling with typed errors (NewCommandError, NewFlagError)
//   - Exit codes (NewExitError)
//   - Strict validation (WithStrictValidation)
//   - Version command (MustVersionCommand)
//   - Signal handling (WithSignalHandling)
//
// Usage:
//
//	go run examples/kitchen-sink/main.go list
//	go run examples/kitchen-sink/main.go list --format json --all
//	go run examples/kitchen-sink/main.go add --title "Buy groceries" --priority high
//	go run examples/kitchen-sink/main.go done --id 1
//	go run examples/kitchen-sink/main.go stats
//	go run examples/kitchen-sink/main.go version
//	go run examples/kitchen-sink/main.go health
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/samber/do/v2"

	examplesinternal "github.com/larsartmann/cmdguard/examples/internal"
	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// --- Config ---

type AppConfig struct {
	LogLevel v2.LogLevel `flag:"log-level" short:"l" default:"info"   help:"Log level (debug, info, warn, error)" env:"TASK_LOG_LEVEL"`
	DataDir  string      `flag:"data-dir"  short:"d" default:"./data" help:"Directory for task storage"           env:"TASK_DATA_DIR"`
}

// --- Flags ---

type ListFlags struct {
	Format   string `flag:"format"   short:"f" default:"table" help:"Output format (table, json, csv, yaml)"`
	All      bool   `flag:"all"      short:"a" default:"false" help:"Show completed tasks too"`
	Priority string `flag:"priority" short:"p" default:""      help:"Filter by priority (low, medium, high)"`
}

type AddFlags struct {
	Title    string `flag:"title"    short:"t" required:"true"  help:"Task title"`
	Priority string `flag:"priority" short:"p" default:"medium" help:"Task priority (low, medium, high)"`
}

type DoneFlags struct {
	ID uint `flag:"id" short:"i" required:"true" help:"Task ID to mark as done"`
}

type StatsFlags struct {
	Format string `flag:"format" short:"f" default:"table" help:"Output format (table, json)"`
}

// --- Domain ---

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type Task struct {
	ID        uint
	Title     string
	Priority  Priority
	Done      bool
	CreatedAt time.Time
}

func (t Task) Row() []string {
	status := "pending"

	if t.Done {
		status = "done"
	}

	return []string{
		fmt.Sprintf("%d", t.ID),
		t.Title,
		string(t.Priority),
		status,
		t.CreatedAt.Format("2006-01-02"),
	}
}

// --- TaskStore (DI Service) ---

type TaskStore struct {
	mu    sync.Mutex
	tasks []Task
	next  uint
}

var (
	_ do.ShutdownerWithError      = (*TaskStore)(nil)
	_ do.HealthcheckerWithContext = (*TaskStore)(nil)
)

func NewTaskStore(i do.Injector) (*TaskStore, error) {
	scope, err := v2.NewScopeFromInjector(i, "provider")

	if err != nil {
		return nil, fmt.Errorf("create scope: %w", err)
	}

	cfg, err := v2.Invoke[*AppConfig](scope)
	if err != nil {
		return nil, fmt.Errorf("resolve config: %w", err)
	}

	fmt.Printf("[store] initialized with data-dir=%s\n", cfg.DataDir)

	return &TaskStore{
		tasks: []Task{},
		next:  1,
	}, nil
}

func (s *TaskStore) Add(title string, priority Priority) Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := Task{
		ID:        s.next,
		Title:     title,
		Priority:  priority,
		Done:      false,
		CreatedAt: time.Now(),
	}

	s.tasks = append(s.tasks, t)
	s.next++

	return t
}

func (s *TaskStore) Done(id uint) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.tasks {
		if s.tasks[i].ID == id {
			s.tasks[i].Done = true

			return s.tasks[i], nil
		}
	}

	return Task{}, fmt.Errorf("task %d not found", id)
}

func (s *TaskStore) List(filterPriority string, showAll bool) []Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []Task

	for _, t := range s.tasks {
		if !showAll && t.Done {
			continue
		}

		if filterPriority != "" && string(t.Priority) != filterPriority {
			continue
		}

		result = append(result, t)
	}

	return result
}

func (s *TaskStore) Stats() (total, pending, done int, byPriority map[Priority]int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	byPriority = map[Priority]int{}

	for _, t := range s.tasks {
		total++

		if t.Done {
			done++
		} else {
			pending++
		}

		byPriority[t.Priority]++
	}

	return total, pending, done, byPriority
}

func (s *TaskStore) Shutdown() error {
	fmt.Println("[store] shutting down — flushing data")

	return nil
}

func (s *TaskStore) HealthCheck(_ context.Context) error {
	return nil
}

// --- Helpers ---

func parsePriority(s string) (Priority, error) {
	switch strings.ToLower(s) {
	case "low":
		return PriorityLow, nil
	case "medium", "":
		return PriorityMedium, nil
	case "high":
		return PriorityHigh, nil
	default:
		return "", fmt.Errorf("invalid priority %q: use low, medium, or high", s)
	}
}

func resolveStore(scope *v2.Scope) (*TaskStore, error) {
	store, err := v2.Invoke[*TaskStore](scope)
	if err != nil {
		return nil, v2.NewCommandError("task", fmt.Errorf("resolve store: %w", err))
	}

	return store, nil
}

// --- Main ---

func main() {
	ctx := context.Background()

	cli, err := v2.NewCLI[AppConfig](
		"taskctl", "A production-grade task manager CLI", AppConfig{},
		v2.WithCLIVersion[AppConfig]("1.0.0"),
		v2.WithEnvPrefix[AppConfig]("TASKCTL_"),
		v2.WithSignalHandling[AppConfig](),
		v2.WithStrictValidation[AppConfig](),
		v2.WithMiddleware[AppConfig](
			v2.TimingMiddleware[AppConfig](func(name string, d time.Duration) {
				fmt.Fprintf(os.Stderr, "[timing] %s took %v\n", name, d)
			}),
			v2.RecoveryMiddleware[AppConfig](),
		),
		v2.WithGroup[AppConfig]("tasks", "Task Management"),
		v2.WithGroup[AppConfig]("system", "System"),
	)
	if err != nil {
		examplesinternal.Fatalf("Error: %v\n", err)
	}

	// Register DI services
	if err := v2.Provide(cli.Scope(), NewTaskStore); err != nil {
		examplesinternal.Fatalf("Error: %v\n", err)
	}

	// Seed demo data
	seedTasks(cli)

	// Build commands
	scope := cli.Scope()

	// list — shows tasks in multiple output formats
	listCmd, err := v2.NewCommand[AppConfig, *ListFlags]("list",
		func(_ context.Context, _ *AppConfig, flags *ListFlags) error {
			store, err := resolveStore(scope)
			if err != nil {
				return err
			}

			tasks := store.List(flags.Priority, flags.All)

			if len(tasks) == 0 {
				fmt.Println("No tasks found.")

				return nil
			}

			format, err := v2.ParseOutputFormat(flags.Format)
			if err != nil {
				return v2.NewFlagError("format", err)
			}

			headers := []string{"ID", "Title", "Priority", "Status", "Created"}

			rows := make([][]string, len(tasks))
			for i, t := range tasks {
				rows[i] = t.Row()
			}

			return v2.OutputTable(format, headers, rows)
		},
		v2.WithShort[AppConfig, *ListFlags]("List tasks"),
		v2.WithExample[AppConfig, *ListFlags](
			"taskctl list --format json --all",
		),
		v2.WithFlags[AppConfig, *ListFlags](&ListFlags{}),
		v2.WithGroupID[AppConfig, *ListFlags]("tasks"),
	)
	if err != nil {
		examplesinternal.Fatalf("Error: %v\n", err)
	}

	if err := v2.AddCommand(cli, listCmd); err != nil {
		examplesinternal.Fatalf("Error: %v\n", err)
	}

	// add — with PreRunE validation
	addCmd, err := v2.NewCommand[AppConfig, *AddFlags]("add",
		func(_ context.Context, _ *AppConfig, flags *AddFlags) error {
			store, err := resolveStore(scope)
			if err != nil {
				return err
			}

			priority, err := parsePriority(flags.Priority)
			if err != nil {
				return v2.NewFlagError("priority", err)
			}

			task := store.Add(flags.Title, priority)
			fmt.Printf("Created task #%d: %s [%s]\n", task.ID, task.Title, task.Priority)

			return nil
		},
		v2.WithShort[AppConfig, *AddFlags]("Add a new task"),
		v2.WithExample[AppConfig, *AddFlags](
			"taskctl add --title \"Fix bug\" --priority high",
		),
		v2.WithFlags[AppConfig, *AddFlags](&AddFlags{}),
		v2.WithPreRunE[AppConfig, *AddFlags](
			func(_ context.Context, _ *AppConfig, flags *AddFlags) error {
				if len(strings.TrimSpace(flags.Title)) < 3 {
					return v2.NewFlagError(
						"title",
						errors.New("title must be at least 3 characters"),
					)
				}

				return nil
			},
		),
		v2.WithGroupID[AppConfig, *AddFlags]("tasks"),
	)
	if err != nil {
		examplesinternal.Fatalf("Error: %v\n", err)
	}

	if err := v2.AddCommand(cli, addCmd); err != nil {
		examplesinternal.Fatalf("Error: %v\n", err)
	}

	// done — with exit codes and PostRunE cleanup
	doneCmd, err := v2.NewCommand[AppConfig, *DoneFlags]("done",
		func(_ context.Context, _ *AppConfig, flags *DoneFlags) error {
			store, err := resolveStore(scope)
			if err != nil {
				return err
			}

			task, err := store.Done(flags.ID)
			if err != nil {
				exitErr, _ := v2.NewExitError(2, v2.NewCommandError("done", err))
				return exitErr
			}

			fmt.Printf("Completed task #%d: %s\n", task.ID, task.Title)

			return nil
		},
		v2.WithShort[AppConfig, *DoneFlags]("Mark a task as done"),
		v2.WithExample[AppConfig, *DoneFlags]("taskctl done --id 1"),
		v2.WithFlags[AppConfig, *DoneFlags](&DoneFlags{}),
		v2.WithPostRunE[AppConfig, *DoneFlags](
			func(_ context.Context, _ *AppConfig, _ *DoneFlags) error {
				fmt.Println("[cleanup] syncing state")

				return nil
			},
		),
		v2.WithGroupID[AppConfig, *DoneFlags]("tasks"),
	)
	if err != nil {
		examplesinternal.Fatalf("Error: %v\n", err)
	}

	if err := v2.AddCommand(cli, doneCmd); err != nil {
		examplesinternal.Fatalf("Error: %v\n", err)
	}

	// stats — statistics with multi-format output
	statsCmd, err := v2.NewCommand[AppConfig, *StatsFlags]("stats",
		func(_ context.Context, _ *AppConfig, flags *StatsFlags) error {
			store, err := resolveStore(scope)
			if err != nil {
				return err
			}

			total, pending, done, byPriority := store.Stats()

			format, err := v2.ParseOutputFormat(flags.Format)
			if err != nil {
				return v2.NewFlagError("format", err)
			}

			headers := []string{"Metric", "Value"}

			rows := [][]string{
				{"Total", fmt.Sprintf("%d", total)},
				{"Pending", fmt.Sprintf("%d", pending)},
				{"Done", fmt.Sprintf("%d", done)},
				{"High Priority", fmt.Sprintf("%d", byPriority[PriorityHigh])},
				{"Medium Priority", fmt.Sprintf("%d", byPriority[PriorityMedium])},
				{"Low Priority", fmt.Sprintf("%d", byPriority[PriorityLow])},
			}

			return v2.OutputTable(format, headers, rows)
		},
		v2.WithShort[AppConfig, *StatsFlags]("Show task statistics"),
		v2.WithFlags[AppConfig, *StatsFlags](&StatsFlags{}),
		v2.WithGroupID[AppConfig, *StatsFlags]("tasks"),
	)
	if err != nil {
		examplesinternal.Fatalf("Error: %v\n", err)
	}

	if err := v2.AddCommand(cli, statsCmd); err != nil {
		examplesinternal.Fatalf("Error: %v\n", err)
	}

	// health — system command with health check
	healthCmd, err := v2.NewCommand[AppConfig, v2.NoFlags]("health",
		func(ctx context.Context, _ *AppConfig, _ v2.NoFlags) error {
			if err := cli.HealthCheckWithContext(ctx); err != nil {
				fmt.Printf("UNHEALTHY: %v\n", err)

				exitErr, _ := v2.NewExitError(1, err)
				return exitErr
			}

			fmt.Println("All systems healthy")

			return nil
		},
		v2.WithShort[AppConfig, v2.NoFlags]("Check system health"),
		v2.WithGroupID[AppConfig, v2.NoFlags]("system"),
	)
	if err != nil {
		examplesinternal.Fatalf("Error: %v\n", err)
	}

	if err := v2.AddCommand(cli, healthCmd); err != nil {
		examplesinternal.Fatalf("Error: %v\n", err)
	}

	// version
	versionCmd := v2.MustVersionCommand[AppConfig](cli)
	if err := v2.AddCommand(cli, versionCmd); err != nil {
		examplesinternal.Fatalf("Error: %v\n", err)
	}

	examplesinternal.Execute(ctx, cli)
}

func seedTasks(cli *v2.CLI[AppConfig]) {
	store, err := v2.Invoke[*TaskStore](cli.Scope())
	if err != nil {
		return
	}

	store.Add("Set up CI pipeline", PriorityHigh)
	store.Add("Write API documentation", PriorityMedium)
	store.Add("Review pull request #42", PriorityLow)
}
