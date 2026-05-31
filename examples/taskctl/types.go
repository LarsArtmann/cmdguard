package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/samber/do/v2"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// --- Config ---

type AppConfig struct {
	LogLevel v2.LogLevel `flag:"log-level" short:"l" default:"info"   help:"Log level (debug, info, warn, error)" env:"TASK_LOG_LEVEL"`
	DataDir  string      `flag:"data-dir"  short:"d" default:"./data" help:"Directory for task storage"           env:"TASK_DATA_DIR"`
	Timeout  v2.Duration `flag:"timeout"             default:"30s"    help:"Default operation timeout"`
	Port     v2.Port     `flag:"port"                default:"8080"   help:"API port"`
	Verbose  int         `flag:"verbose"   short:"v" default:"0"      help:"Verbosity (-v, -vv, -vvv)"                                 count:"true"`
}

// --- Priority Enum ---

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

func parsePriority(s string) (Priority, error) {
	switch strings.ToLower(s) {
	case string(PriorityLow):
		return PriorityLow, nil
	case string(PriorityMedium), "":
		return PriorityMedium, nil
	case string(PriorityHigh):
		return PriorityHigh, nil
	default:
		return "", fmt.Errorf("invalid priority %q: use low, medium, or high", s)
	}
}

// --- Domain ---

const (
	taskStatusPending = "pending"
	taskStatusDone    = "done"
)

type Task struct {
	ID        uint
	Title     string
	Priority  Priority
	Done      bool
	CreatedAt time.Time
}

func (t Task) Row() []string {
	status := taskStatusPending
	if t.Done {
		status = taskStatusDone
	}
	return []string{
		strconv.FormatUint(uint64(t.ID), 10),
		t.Title,
		string(t.Priority),
		status,
		t.CreatedAt.Format("2006-01-02"),
	}
}

// --- Flags ---

type ListFlags struct {
	Format   string `flag:"format"   short:"f" default:"table" help:"Output format (table, json, csv, yaml)"`
	All      bool   `flag:"all"      short:"a" default:"false" help:"Show completed tasks too"`
	Priority string `flag:"priority"           default:""      help:"Filter by priority (low, medium, high)"`
}

type AddFlags struct {
	Title    string `flag:"title"    short:"t" required:"true" help:"Task title"                   prompt:"Task title?"`
	Priority string `flag:"priority" short:"P"                 help:"Priority (low, medium, high)"                      default:"medium"`
}

type DoneFlags struct {
	ID uint `flag:"id" short:"i" required:"true" help:"Task ID to complete"`
}

type StatsFlags struct {
	Format string `flag:"format" short:"f" default:"table" help:"Output format (table, json)"`
}

type DBFlags struct {
	Env   string `flag:"env"   short:"e" help:"Environment (dev, staging, prod)" default:"dev"`
	Force bool   `flag:"force" short:"f" help:"Skip confirmation prompts"        default:"false"`
}

type InspectFlags struct {
	ShowMetadata bool `flag:"metadata" short:"m" default:"false" help:"Show metadata"`
}

// --- TaskStore (DI Service) ---

type TaskStore struct {
	mu    chan struct{}
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
		mu:    make(chan struct{}, 1),
		tasks: []Task{},
		next:  1,
	}, nil
}

func (s *TaskStore) lock()   { s.mu <- struct{}{} }
func (s *TaskStore) unlock() { <-s.mu }

func (s *TaskStore) Add(title string, priority Priority) Task {
	s.lock()
	defer s.unlock()

	t := Task{ID: s.next, Title: title, Priority: priority, Done: false, CreatedAt: time.Now()}
	s.tasks = append(s.tasks, t)
	s.next++
	return t
}

func (s *TaskStore) Done(id uint) (Task, error) {
	s.lock()
	defer s.unlock()

	for i := range s.tasks {
		if s.tasks[i].ID == id {
			s.tasks[i].Done = true
			return s.tasks[i], nil
		}
	}
	return Task{}, fmt.Errorf("task %d not found", id)
}

func (s *TaskStore) List(filterPriority string, showAll bool) []Task {
	s.lock()
	defer s.unlock()

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

func (s *TaskStore) Get(id uint) (Task, bool) {
	s.lock()
	defer s.unlock()

	for _, t := range s.tasks {
		if t.ID == id {
			return t, true
		}
	}
	return Task{}, false
}

func (s *TaskStore) IDs() []string {
	s.lock()
	defer s.unlock()

	ids := make([]string, len(s.tasks))
	for i, t := range s.tasks {
		ids[i] = strconv.FormatUint(uint64(t.ID), 10)
	}
	return ids
}

func (s *TaskStore) Stats() (total, pending, done int, byPriority map[Priority]int) {
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
