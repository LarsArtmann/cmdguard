package v2

import (
	"testing"
)

func TestGuardedCommand_Scope(t *testing.T) {
	t.Parallel()
	t.Run("returns injector", func(t *testing.T) {
		t.Parallel()
		g, err := New[testAppConfig, NoFlags]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		injector := g.Scope()
		if injector == nil {
			t.Fatal("expected non-nil injector")
		}
	})

	t.Run("returns scope struct", func(t *testing.T) {
		t.Parallel()
		g, err := New[testAppConfig, NoFlags]("myapp", "My CLI", testAppConfig{})
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

func TestGuardedCommand_Config(t *testing.T) {
	t.Run("returns config", func(t *testing.T) {
		t.Parallel()
		defaults := testAppConfig{Verbose: true, Output: "/tmp/out"}

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

		if cfg.Output != "/tmp/out" {
			t.Errorf("cfg.Output = %q, want %q", cfg.Output, "/tmp/out")
		}
	})

	t.Run("SetConfig updates config", func(t *testing.T) {
		t.Parallel()
		g, err := New[testAppConfig, NoFlags]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		newCfg := testAppConfig{Verbose: true, Output: "/new/path"}
		g.SetConfig(newCfg)

		cfg := g.Config()
		if !cfg.Verbose {
			t.Error("cfg.Verbose = false, want true")
		}

		if cfg.Output != "/new/path" {
			t.Errorf("cfg.Output = %q, want %q", cfg.Output, "/new/path")
		}
	})
}

func TestGuardedCommand_RootCommand(t *testing.T) {
	t.Run("returns cobra command", func(t *testing.T) {
		t.Parallel()
		g, err := New[testAppConfig, NoFlags]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		rootCmd := g.RootCommand()
		if rootCmd == nil {
			t.Fatal("expected non-nil root command")
		}

		if rootCmd.Use != "myapp" {
			t.Errorf("rootCmd.Use = %q, want %q", rootCmd.Use, "myapp")
		}

		if rootCmd.Short != "My CLI" {
			t.Errorf("rootCmd.Short = %q, want %q", rootCmd.Short, "My CLI")
		}
	})
}

func TestGuardedCommand_Metadata(t *testing.T) {
	t.Run("Name returns name", func(t *testing.T) {
		t.Parallel()
		g, err := New[testAppConfig, NoFlags]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if g.Name() != "myapp" {
			t.Errorf("Name() = %q, want %q", g.Name(), "myapp")
		}
	})

	t.Run("Short returns short description", func(t *testing.T) {
		t.Parallel()
		g, err := New[testAppConfig, NoFlags]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if g.Short() != "My CLI" {
			t.Errorf("Short() = %q, want %q", g.Short(), "My CLI")
		}
	})

	t.Run("Long returns long description", func(t *testing.T) {
		t.Parallel()
		g, err := NewWithLong[testAppConfig, NoFlags](
			"myapp",
			"short",
			"long desc",
			testAppConfig{},
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if g.Long() != "long desc" {
			t.Errorf("Long() = %q, want %q", g.Long(), "long desc")
		}
	})

	t.Run("SetLong updates long description", func(t *testing.T) {
		t.Parallel()
		g, err := New[testAppConfig, NoFlags]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		g.SetLong("new long description")

		if g.Long() != "new long description" {
			t.Errorf("Long() = %q, want %q", g.Long(), "new long description")
		}

		if g.RootCommand().Long != "new long description" {
			t.Errorf("rootCmd.Long = %q, want %q", g.RootCommand().Long, "new long description")
		}
	})

	t.Run("SetVersion sets version", func(t *testing.T) {
		t.Parallel()
		g, err := New[testAppConfig, NoFlags]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		g.SetVersion("v1.0.0")

		if g.RootCommand().Version != "v1.0.0" {
			t.Errorf("rootCmd.Version = %q, want %q", g.RootCommand().Version, "v1.0.0")
		}
	})
}

func TestGuardedCommand_AddGlobalFlag(t *testing.T) {
	t.Run("adds global string flag", func(t *testing.T) {
		t.Parallel()
		g, err := New[testAppConfig, NoFlags]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		g.AddGlobalFlag("config", "c", "/etc/config.yaml", "Config file path")

		flag := g.RootCommand().PersistentFlags().Lookup("config")
		if flag == nil {
			t.Fatal("expected non-nil flag")
		}

		if flag.Shorthand != "c" {
			t.Errorf("flag.Shorthand = %q, want %q", flag.Shorthand, "c")
		}

		if flag.DefValue != "/etc/config.yaml" {
			t.Errorf("flag.DefValue = %q, want %q", flag.DefValue, "/etc/config.yaml")
		}
	})

	t.Run("adds global bool flag", func(t *testing.T) {
		t.Parallel()
		g, err := New[testAppConfig, NoFlags]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		g.AddGlobalBoolFlag("debug", "d", true, "Enable debug mode")

		flag := g.RootCommand().PersistentFlags().Lookup("debug")
		if flag == nil {
			t.Fatal("expected non-nil flag")
		}

		if flag.Shorthand != "d" {
			t.Errorf("flag.Shorthand = %q, want %q", flag.Shorthand, "d")
		}

		if flag.DefValue != "true" {
			t.Errorf("flag.DefValue = %q, want %q", flag.DefValue, "true")
		}
	})
}
