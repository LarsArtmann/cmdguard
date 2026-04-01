package v2

import (
	"slices"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func containsString(slice []string, s string) bool {
	return slices.Contains(slice, s)
}

func TestNewFlagRegistry(t *testing.T) {
	t.Parallel()
	t.Run("valid config", func(t *testing.T) {
		type testConfig struct {
			Name  string `default:"default-name" flag:"name"  help:"name help"`
			Count int    `default:"10"           flag:"count" help:"count help"`
		}

		cfg := testConfig{}

		registry, err := NewFlagRegistry(cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if registry == nil {
			t.Fatal("expected non-nil registry")
		}

		if len(registry.Tags()) != 2 {
			t.Errorf("len(registry.Tags()) = %d, want 2", len(registry.Tags()))
		}
	})

	t.Run("non-struct config", func(t *testing.T) {
		registry, err := NewFlagRegistry("not a struct")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if registry != nil {
			t.Error("expected nil registry")
		}

		if !strings.Contains(err.Error(), "expected struct") {
			t.Errorf("error should contain 'expected struct', got %q", err.Error())
		}
	})

	t.Run("config with short flags", func(t *testing.T) {
		type testConfig struct {
			Name string `flag:"name" help:"name help" short:"n"`
		}

		registry, err := NewFlagRegistry(testConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		tags := registry.Tags()
		if len(tags) != 1 {
			t.Errorf("len(tags) = %d, want 1", len(tags))
		}

		if len(tags) > 0 && tags[0].Short != "n" {
			t.Errorf("tags[0].Short = %q, want %q", tags[0].Short, "n")
		}
	})
}

func TestFlagRegistry_RegisterFlags(t *testing.T) {
	t.Run("registers all flag types", func(t *testing.T) {
		type testConfig struct {
			String  string   `default:"str"   flag:"string"`
			Bool    bool     `default:"true"  flag:"bool"`
			Int     int      `default:"42"    flag:"int"`
			Uint    uint     `default:"10"    flag:"uint"`
			Uint64  uint64   `default:"1024"  flag:"uint64"`
			Float   float64  `default:"3.14"  flag:"float"`
			Strings []string `default:"a,b,c" flag:"strings"`
		}

		registry, err := NewFlagRegistry(testConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		cmd := &cobra.Command{Use: "test"}

		err = registry.RegisterFlags(cmd)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		flags := cmd.PersistentFlags()
		if flags.Lookup("string") == nil {
			t.Error("expected 'string' flag to be registered")
		}

		if flags.Lookup("bool") == nil {
			t.Error("expected 'bool' flag to be registered")
		}

		if flags.Lookup("int") == nil {
			t.Error("expected 'int' flag to be registered")
		}

		if flags.Lookup("uint") == nil {
			t.Error("expected 'uint' flag to be registered")
		}

		if flags.Lookup("uint64") == nil {
			t.Error("expected 'uint64' flag to be registered")
		}

		if flags.Lookup("float") == nil {
			t.Error("expected 'float' flag to be registered")
		}

		if flags.Lookup("strings") == nil {
			t.Error("expected 'strings' flag to be registered")
		}
	})

	t.Run("registers custom types", func(t *testing.T) {
		type testConfig struct {
			Level  LogLevel  `default:"info" flag:"level"`
			Format LogFormat `default:"json" flag:"format"`
		}

		registry, err := NewFlagRegistry(testConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		cmd := &cobra.Command{Use: "test"}

		err = registry.RegisterFlags(cmd)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cmd.PersistentFlags().Lookup("level") == nil {
			t.Error("expected 'level' flag to be registered")
		}

		if cmd.PersistentFlags().Lookup("format") == nil {
			t.Error("expected 'format' flag to be registered")
		}
	})

	t.Run("registers Duration type", func(t *testing.T) {
		type testConfig struct {
			Timeout Duration `default:"30s" flag:"timeout"`
		}

		registry, err := NewFlagRegistry(testConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		cmd := &cobra.Command{Use: "test"}

		err = registry.RegisterFlags(cmd)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		flag := cmd.PersistentFlags().Lookup("timeout")
		if flag == nil {
			t.Fatal("expected 'timeout' flag to be registered")
		}

		if flag.DefValue != "30s" {
			t.Errorf("flag.DefValue = %q, want %q", flag.DefValue, "30s")
		}
	})

	t.Run("registers enum with values", func(t *testing.T) {
		type testConfig struct {
			Mode Enum `default:"dev" flag:"mode" values:"dev,staging,prod"`
		}

		registry, err := NewFlagRegistry(testConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		cmd := &cobra.Command{Use: "test"}

		err = registry.RegisterFlags(cmd)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		flag := cmd.PersistentFlags().Lookup("mode")
		if flag == nil {
			t.Fatal("expected 'mode' flag to be registered")
		}

		if !strings.Contains(flag.Usage, "one of: dev, staging, prod") {
			t.Errorf("flag.Usage should contain 'one of: dev, staging, prod', got %q", flag.Usage)
		}
	})
}

func TestFlagRegistry_Tags(t *testing.T) {
	t.Run("returns all tags", func(t *testing.T) {
		type testConfig struct {
			Name  string `flag:"name"  help:"name help"`
			Count int    `flag:"count" help:"count help"`
		}

		registry, err := NewFlagRegistry(testConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		tags := registry.Tags()
		if len(tags) != 2 {
			t.Errorf("len(tags) = %d, want 2", len(tags))
		}

		names := make([]string, len(tags))
		for i, tag := range tags {
			names[i] = tag.Name
		}

		if !containsString(names, "name") {
			t.Errorf("names should contain 'name', got %v", names)
		}

		if !containsString(names, "count") {
			t.Errorf("names should contain 'count', got %v", names)
		}
	})
}

func TestFlagRegistry_FlagNames(t *testing.T) {
	t.Run("returns all flag names", func(t *testing.T) {
		type testConfig struct {
			Verbose bool   `flag:"verbose" short:"v"`
			Config  string `flag:"config"  short:"c"`
			Output  string `flag:"output"  short:"o"`
		}

		registry, err := NewFlagRegistry(testConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		names := registry.FlagNames()
		if len(names) != 3 {
			t.Errorf("len(names) = %d, want 3", len(names))
		}

		if !containsString(names, "verbose") {
			t.Errorf("names should contain 'verbose', got %v", names)
		}

		if !containsString(names, "config") {
			t.Errorf("names should contain 'config', got %v", names)
		}

		if !containsString(names, "output") {
			t.Errorf("names should contain 'output', got %v", names)
		}
	})

	t.Run("empty registry returns empty slice", func(t *testing.T) {
		type emptyConfig struct{}

		registry, err := NewFlagRegistry(emptyConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		names := registry.FlagNames()
		if len(names) != 0 {
			t.Errorf("len(names) = %d, want 0", len(names))
		}
	})
}
