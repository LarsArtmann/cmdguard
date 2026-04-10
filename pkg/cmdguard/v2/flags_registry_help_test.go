package v2

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestFlagRegistry_GenerateHelp(t *testing.T) {
	t.Parallel()
	t.Run("generates help for all flags", func(t *testing.T) {
		t.Parallel()

		type testConfig struct {
			Name    string `default:"default" flag:"name,n"    help:"The name to use"`
			Verbose bool   `                  flag:"verbose,v" help:"Enable verbose output"`
		}

		registry, err := NewFlagRegistry(testConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		help := registry.GenerateHelp()
		if !strings.Contains(help, "--name") {
			t.Errorf("help should contain '--name', got %q", help)
		}

		if !strings.Contains(help, "-n") {
			t.Errorf("help should contain '-n', got %q", help)
		}

		if !strings.Contains(help, "The name to use") {
			t.Errorf("help should contain 'The name to use', got %q", help)
		}

		if !strings.Contains(help, "default: default") {
			t.Errorf("help should contain 'default: default', got %q", help)
		}

		if !strings.Contains(help, "--verbose") {
			t.Errorf("help should contain '--verbose', got %q", help)
		}

		if !strings.Contains(help, "-v") {
			t.Errorf("help should contain '-v', got %q", help)
		}

		if !strings.Contains(help, "Enable verbose output") {
			t.Errorf("help should contain 'Enable verbose output', got %q", help)
		}
	})

	t.Run("help formatting without optional elements", func(t *testing.T) {
		t.Parallel()

		type testConfig struct {
			Name string `flag:"name" help:"The name"`
		}

		registry, err := NewFlagRegistry(testConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		help := registry.GenerateHelp()
		if !strings.Contains(help, "--name") {
			t.Errorf("help should contain '--name', got %q", help)
		}
	})
}

func TestFlagRegistry_ShortFlags(t *testing.T) {
	t.Parallel()
	t.Run("register short flags", func(t *testing.T) {
		t.Parallel()

		type testConfig struct {
			Name    string `default:""      flag:"name"    short:"n"`
			Count   int    `default:"0"     flag:"count"   short:"c"`
			Verbose bool   `default:"false" flag:"verbose" short:"v"`
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

		flags := cmd.Flags()

		nameFlag := flags.Lookup("name")
		if nameFlag == nil {
			t.Fatal("expected 'name' flag to be registered")
		}

		if nameFlag.Shorthand != "n" {
			t.Errorf("nameFlag.Shorthand = %q, want %q", nameFlag.Shorthand, "n")
		}

		countFlag := flags.Lookup("count")
		if countFlag == nil {
			t.Fatal("expected 'count' flag to be registered")
		}

		if countFlag.Shorthand != "c" {
			t.Errorf("countFlag.Shorthand = %q, want %q", countFlag.Shorthand, "c")
		}

		verboseFlag := flags.Lookup("verbose")
		if verboseFlag == nil {
			t.Fatal("expected 'verbose' flag to be registered")
		}

		if verboseFlag.Shorthand != "v" {
			t.Errorf("verboseFlag.Shorthand = %q, want %q", verboseFlag.Shorthand, "v")
		}
	})

	t.Run("register uint short flags", func(t *testing.T) {
		t.Parallel()

		type testConfig struct {
			Workers uint `default:"4" flag:"workers" short:"w"`
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

		flag := cmd.Flags().Lookup("workers")
		if flag == nil {
			t.Fatal("expected 'workers' flag to be registered")
		}

		if flag.Shorthand != "w" {
			t.Errorf("flag.Shorthand = %q, want %q", flag.Shorthand, "w")
		}
	})

	t.Run("register uint64 short flags", func(t *testing.T) {
		t.Parallel()

		type testConfig struct {
			Bytes uint64 `default:"1024" flag:"bytes" short:"b"`
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

		flag := cmd.Flags().Lookup("bytes")
		if flag == nil {
			t.Fatal("expected 'bytes' flag to be registered")
		}

		if flag.Shorthand != "b" {
			t.Errorf("flag.Shorthand = %q, want %q", flag.Shorthand, "b")
		}
	})
}
