// Package main demonstrates dependency injection patterns with cmdguard v2.
//
// This example shows:
// - Registering services with v2.Provide
// - Injecting services into command handlers with v2.Invoke
// - Service lifecycle: health checks and shutdown
// - Layered architecture: config -> service -> handler
//
// Usage:
//
//	go run examples/di-patterns/main.go list
//	go run examples/di-patterns/main.go add --title "Buy groceries" --priority high
//	go run examples/di-patterns/main.go check
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	id "github.com/larsartmann/go-branded-id"
	"github.com/samber/do/v2"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// AppConfig is the application configuration.
type AppConfig struct {
	StorePath string `default:"/tmp/todos.json" flag:"store" help:"Path to task store" short:"s"`
}

// TaskBrand is the phantom type for TaskID.
type TaskBrand struct{}

// TaskID is a strongly-typed identifier for a task.
type TaskID = id.ID[TaskBrand, int]

// Task represents a single work item.
type Task struct {
	ID       TaskID `json:"id"`
	Title    string `json:"title"`
	Priority string `json:"priority"`
	Done     bool   `json:"done"`
}

// NextID is the next TaskID to be assigned.
type NextID = TaskID

// TaskStore is an in-memory task store.
type TaskStore struct {
	mu     sync.Mutex
	tasks  []Task
	nextID NextID
}

var (
	_ do.ShutdownerWithError      = (*TaskStore)(nil)
	_ do.HealthcheckerWithContext = (*TaskStore)(nil)
)

// NewTaskStore creates a new task store.
func NewTaskStore(i do.Injector) (*TaskStore, error) {
	return &TaskStore{
		tasks:  []Task{},
		nextID: id.NewID[TaskBrand](1),
	}, nil
}

// Shutdown implements do.ShutdownerWithError.
func (s *TaskStore) Shutdown() error {
	fmt.Println("TaskStore: flushing data...")

	time.Sleep(10 * time.Millisecond)

	fmt.Println("TaskStore: shutdown complete")

	return nil
}

// HealthCheck implements do.HealthcheckerWithContext.
func (s *TaskStore) HealthCheck(_ context.Context) error {
	if s == nil {
		return errors.New("task store is nil")
	}

	return nil
}

// Add creates a new task.
func (s *TaskStore) Add(title, priority string) Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := Task{
		ID:       s.nextID,
		Title:    title,
		Priority: priority,
		Done:     false,
	}

	s.tasks = append(s.tasks, task)
	s.nextID = id.NewID[TaskBrand](s.nextID.Get() + 1)

	return task
}

// List returns all tasks.
func (s *TaskStore) List() []Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]Task, len(s.tasks))
	copy(result, s.tasks)

	return result
}

// AddFlags defines flags for the add command.
type AddFlags struct {
	Title    string `flag:"title"    short:"t" help:"Task title"                   required:"true"`
	Priority string `flag:"priority" short:"p" help:"Priority (low, medium, high)"                 default:"medium"`
}

func main() {
	cli, err := v2.NewCLI[AppConfig]("di-patterns", "DI patterns demo", AppConfig{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := v2.Provide(cli.Scope(), NewTaskStore); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	listCmd, err := v2.NewCommand[AppConfig, v2.NoFlags](
		"list",
		func(_ context.Context, _ *AppConfig, _ v2.NoFlags) error {
			store, err := v2.Invoke[*TaskStore](cli.Scope())
			if err != nil {
				return fmt.Errorf("resolving task store: %w", err)
			}

			tasks := store.List()

			if len(tasks) == 0 {
				fmt.Println("No tasks found. Add one with 'add --title \"My task\"'")

				return nil
			}

			fmt.Println("Tasks:")

			for _, t := range tasks {
				status := " "
				if t.Done {
					status = "x"
				}

				fmt.Printf("  [%s] #%d %s (%s)\n", status, t.ID.Get(), t.Title, t.Priority)
			}

			return nil
		},
		v2.WithShort[AppConfig, v2.NoFlags]("List all tasks"),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := v2.AddCommand(cli, listCmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	addCmd, err := v2.NewCommand[AppConfig, *AddFlags](
		"add",
		func(_ context.Context, _ *AppConfig, flags *AddFlags) error {
			store, err := v2.Invoke[*TaskStore](cli.Scope())
			if err != nil {
				return fmt.Errorf("resolving task store: %w", err)
			}

			task := store.Add(flags.Title, flags.Priority)

			fmt.Printf("Added task #%d: %s (%s)\n", task.ID.Get(), task.Title, task.Priority)

			return nil
		},
		v2.WithShort[AppConfig, *AddFlags]("Add a new task"),
		v2.WithFlags[AppConfig, *AddFlags](&AddFlags{}),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := v2.AddCommand(cli, addCmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	checkCmd, err := v2.NewCommand[AppConfig, v2.NoFlags](
		"check",
		func(ctx context.Context, _ *AppConfig, _ v2.NoFlags) error {
			fmt.Println("Running health checks...")

			if err := cli.HealthCheckWithContext(ctx); err != nil {
				fmt.Printf("Health check FAILED: %v\n", err)

				return err
			}

			fmt.Println("All health checks PASSED!")

			return nil
		},
		v2.WithShort[AppConfig, v2.NoFlags]("Run health checks"),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := v2.AddCommand(cli, checkCmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cli.ExecuteAndExit(context.Background())
}
