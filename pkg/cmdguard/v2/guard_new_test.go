package v2

import (
	"context"
	"errors"
	"testing"
)

type testAppConfig struct {
	Verbose bool   `default:"false" flag:"verbose" help:"Enable verbose output" short:"v"`
	Output  string `default:"-"     flag:"output"  help:"Output file"           short:"o"`
}

// newTestCmd creates a Command with a RunE handler for testing.
// This helper reduces duplication across test files that need simple valid commands.
// If err is provided, the RunE handler returns it; otherwise returns nil.
func newTestCmd(use string, err ...error) Command[testAppConfig, NoFlags] {
	var runErr error
	if len(err) > 0 {
		runErr = err[0]
	}

	return Command[testAppConfig, NoFlags]{
		Use: use,
		RunE: func(_ context.Context, _ *testAppConfig, _ NoFlags) error {
			return runErr
		},
	}
}

func TestVersion(t *testing.T) {
	t.Parallel()
	t.Run("version is set", func(t *testing.T) {
		t.Parallel()
		if Version == "" {
			t.Error("Version is empty")
		}

		if Version != "2.1.0" {
			t.Errorf("Version = %q, want %q", Version, "2.1.0")
		}
	})
}

func TestNew(t *testing.T) {
	t.Parallel()
	t.Run("creates GuardedCommand", func(t *testing.T) {
		t.Parallel()
		defaults := testAppConfig{}

		g, err := New[testAppConfig, NoFlags]("myapp", "My CLI application", defaults)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if g == nil {
			t.Fatal("expected non-nil GuardedCommand")
		}

		if g.Name() != "myapp" {
			t.Errorf("Name() = %q, want %q", g.Name(), "myapp")
		}

		if g.Short() != "My CLI application" {
			t.Errorf("Short() = %q, want %q", g.Short(), "My CLI application")
		}

		if g.Long() != "" {
			t.Errorf("Long() = %q, want empty", g.Long())
		}
	})

	t.Run("error: empty name", func(t *testing.T) {
		t.Parallel()
		defaults := testAppConfig{}

		g, err := New[testAppConfig, NoFlags]("", "My CLI", defaults)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrInvalidCommand) {
			t.Errorf("expected ErrInvalidCommand, got %v", err)
		}

		if g != nil {
			t.Errorf("expected nil GuardedCommand, got %v", g)
		}
	})

	t.Run("registers config in scope", func(t *testing.T) {
		t.Parallel()
		defaults := testAppConfig{Verbose: true}

		g, err := New[testAppConfig, NoFlags]("myapp", "My CLI", defaults)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		cfg := g.Config()
		if cfg == nil {
			t.Fatal("expected non-nil config")
		}

		if !cfg.Verbose {
			t.Error("cfg.Verbose = false, want true")
		}
	})

	t.Run("creates scope", func(t *testing.T) {
		t.Parallel()
		defaults := testAppConfig{}

		g, err := New[testAppConfig, NoFlags]("myapp", "My CLI", defaults)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		scope := g.ScopeStruct()
		if scope == nil {
			t.Fatal("expected non-nil scope")
		}

		if scope.Name() != "myapp" {
			t.Errorf("scope.Name() = %q, want %q", scope.Name(), "myapp")
		}
	})
}

func TestNewWithLong(t *testing.T) {
	t.Parallel()
	t.Run("creates GuardedCommand with long description", func(t *testing.T) {
		t.Parallel()
		defaults := testAppConfig{}

		g, err := NewWithLong[testAppConfig, NoFlags](
			"myapp",
			"short",
			"long description",
			defaults,
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if g == nil {
			t.Fatal("expected non-nil GuardedCommand")
		}

		if g.Name() != "myapp" {
			t.Errorf("Name() = %q, want %q", g.Name(), "myapp")
		}

		if g.Short() != "short" {
			t.Errorf("Short() = %q, want %q", g.Short(), "short")
		}

		if g.Long() != "long description" {
			t.Errorf("Long() = %q, want %q", g.Long(), "long description")
		}
	})

	t.Run("error: empty name", func(t *testing.T) {
		t.Parallel()
		defaults := testAppConfig{}

		g, err := NewWithLong[testAppConfig, NoFlags]("", "short", "long", defaults)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if g != nil {
			t.Errorf("expected nil GuardedCommand, got %v", g)
		}
	})
}

func TestNew_FlagTypeValidation(t *testing.T) {
	t.Parallel()
	t.Run("rejects invalid flag type in New", func(t *testing.T) {
		t.Parallel()
		g, err := New[testAppConfig, int]("myapp", "My CLI", testAppConfig{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrInvalidFlagType) {
			t.Errorf("expected ErrInvalidFlagType, got %v", err)
		}

		if g != nil {
			t.Errorf("expected nil GuardedCommand, got %v", g)
		}
	})

	t.Run("accepts NoFlags in New", func(t *testing.T) {
		t.Parallel()
		g, err := New[testAppConfig, NoFlags]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if g == nil {
			t.Fatal("expected non-nil GuardedCommand")
		}
	})

	t.Run("accepts pointer to struct in New", func(t *testing.T) {
		t.Parallel()
		type cmdFlags struct {
			Name string `flag:"name"`
		}

		g, err := New[testAppConfig, *cmdFlags]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if g == nil {
			t.Fatal("expected non-nil GuardedCommand")
		}
	})
}

func TestNewSimple(t *testing.T) {
	t.Parallel()
	t.Run("creates SimpleCLI with defaults", func(t *testing.T) {
		t.Parallel()
		g, err := NewSimple("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if g == nil {
			t.Fatal("expected non-nil SimpleCLI")
		}

		if g.Name() != "myapp" {
			t.Errorf("Name() = %q, want %q", g.Name(), "myapp")
		}

		if g.Short() != "My CLI" {
			t.Errorf("Short() = %q, want %q", g.Short(), "My CLI")
		}
	})

	t.Run("returns error for empty name", func(t *testing.T) {
		t.Parallel()
		g, err := NewSimple("", "My CLI", testAppConfig{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if g != nil {
			t.Errorf("expected nil SimpleCLI, got %v", g)
		}
	})
}

func TestNewSimpleWithLong(t *testing.T) {
	t.Parallel()
	t.Run("creates SimpleCLI with long description", func(t *testing.T) {
		t.Parallel()
		g, err := NewSimpleWithLong(
			"myapp",
			"short",
			"long description",
			testAppConfig{},
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if g == nil {
			t.Fatal("expected non-nil SimpleCLI")
		}

		if g.Name() != "myapp" {
			t.Errorf("Name() = %q, want %q", g.Name(), "myapp")
		}

		if g.Short() != "short" {
			t.Errorf("Short() = %q, want %q", g.Short(), "short")
		}

		if g.Long() != "long description" {
			t.Errorf("Long() = %q, want %q", g.Long(), "long description")
		}
	})

	t.Run("returns error for empty name", func(t *testing.T) {
		t.Parallel()
		g, err := NewSimpleWithLong("", "short", "long", testAppConfig{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if g != nil {
			t.Errorf("expected nil SimpleCLI, got %v", g)
		}
	})
}
