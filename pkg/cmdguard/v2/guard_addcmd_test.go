package v2

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestGuardedCommand_AddCommand(t *testing.T) {
	t.Run("adds valid command", func(t *testing.T) {
		g, err := New[testAppConfig, NoFlags]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = g.AddCommand(newTestCmd("greet"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		rootCmd := g.RootCommand()
		if len(rootCmd.Commands()) != 1 {
			t.Errorf("len(rootCmd.Commands()) = %d, want 1", len(rootCmd.Commands()))
		}
	})

	t.Run("error: invalid command", func(t *testing.T) {
		g, err := New[testAppConfig, NoFlags]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		cmd := Command[testAppConfig, NoFlags]{
			Use: "greet",
		}

		err = g.AddCommand(cmd)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrMissingHandler) {
			t.Errorf("expected ErrMissingHandler, got %v", err)
		}
	})

	t.Run("adds command with subcommands", func(t *testing.T) {
		g, err := New[testAppConfig, NoFlags]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		subCmd := newTestCmd("list")

		cmd := Command[testAppConfig, NoFlags]{
			Use:      "greet",
			Commands: []Command[testAppConfig, NoFlags]{subCmd},
		}

		err = g.AddCommand(cmd)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		rootCmd := g.RootCommand()
		if len(rootCmd.Commands()) != 1 {
			t.Errorf("len(rootCmd.Commands()) = %d, want 1", len(rootCmd.Commands()))
		}

		greetCmd := rootCmd.Commands()[0]
		if len(greetCmd.Commands()) != 1 {
			t.Errorf("len(greetCmd.Commands()) = %d, want 1", len(greetCmd.Commands()))
		}
	})

	t.Run("error: duplicate command", func(t *testing.T) {
		g, err := New[testAppConfig, NoFlags]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		cmd1 := newTestCmd("greet")
		cmd2 := newTestCmd("greet")

		err = g.AddCommand(cmd1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = g.AddCommand(cmd2)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrDuplicateCommand) {
			t.Errorf("expected ErrDuplicateCommand, got %v", err)
		}

		if !strings.Contains(err.Error(), "greet") {
			t.Errorf("error should contain 'greet', got %q", err.Error())
		}
	})
}

func TestAddAnyCommand(t *testing.T) {
	t.Run("adds command with different flag type", func(t *testing.T) {
		g, err := New[testAppConfig, NoFlags]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		type greetFlags struct {
			Name string
		}

		cmd := Command[testAppConfig, *greetFlags]{
			Use:   "greet",
			Short: "Greet someone",
			Flags: &greetFlags{},
			RunE: func(_ context.Context, _ *testAppConfig, _ *greetFlags) error {
				return nil
			},
		}

		err = AddAnyCommand(g, cmd)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		rootCmd := g.RootCommand()
		if len(rootCmd.Commands()) != 1 {
			t.Errorf("len(rootCmd.Commands()) = %d, want 1", len(rootCmd.Commands()))
		}
	})

	t.Run("error when command has no handler", func(t *testing.T) {
		g, err := New[testAppConfig, NoFlags]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		type otherFlags struct {
			Value string
		}

		cmd := Command[testAppConfig, *otherFlags]{
			Use:   "invalid",
			Flags: &otherFlags{},
		}

		err = AddAnyCommand(g, cmd)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrMissingHandler) {
			t.Errorf("expected ErrMissingHandler, got %v", err)
		}
	})

	t.Run("error on duplicate command", func(t *testing.T) {
		g, err := New[testAppConfig, NoFlags]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		type otherFlags struct{}

		cmd := Command[testAppConfig, *otherFlags]{
			Use: "test",
			RunE: func(_ context.Context, _ *testAppConfig, _ *otherFlags) error {
				return nil
			},
		}

		err = AddAnyCommand(g, cmd)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = AddAnyCommand(g, cmd)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrDuplicateCommand) {
			t.Errorf("expected ErrDuplicateCommand, got %v", err)
		}
	})
}
