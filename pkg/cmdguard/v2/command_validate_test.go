package v2

import (
	"context"
	"errors"
	"strings"
	"testing"
)

type testConfig struct {
	Name string
}

func newTestCommand() Command[testConfig, NoFlags] {
	return Command[testConfig, NoFlags]{
		Use: "test",
		RunE: func(_ context.Context, _ *testConfig, _ NoFlags) error {
			return nil
		},
	}
}

func newTestSubcommand(use string) Command[testConfig, NoFlags] {
	return Command[testConfig, NoFlags]{
		Use: use,
		RunE: func(_ context.Context, _ *testConfig, _ NoFlags) error {
			return nil
		},
	}
}

func TestCommand_Validate(t *testing.T) {
	t.Parallel()
	t.Run("valid command with RunE", func(t *testing.T) {
		t.Parallel()

		cmd := newTestCommand()

		err := cmd.Validate()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("valid command with subcommands", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{
			Use: "root",
			Commands: []Command[testConfig, NoFlags]{
				newTestSubcommand("sub"),
			},
		}

		err := cmd.Validate()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error: empty Use field", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{
			RunE: func(_ context.Context, _ *testConfig, _ NoFlags) error {
				return nil
			},
		}

		err := cmd.Validate()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrInvalidCommand) {
			t.Errorf("expected ErrInvalidCommand, got %v", err)
		}

		if !strings.Contains(err.Error(), "no Use field") {
			t.Errorf("error should contain 'no Use field', got %q", err.Error())
		}
	})

	t.Run("error: no RunE and no subcommands", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{
			Use: "test",
		}

		err := cmd.Validate()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrMissingHandler) {
			t.Errorf("expected ErrMissingHandler, got %v", err)
		}

		if !strings.Contains(err.Error(), "no RunE and no subcommands") {
			t.Errorf("error should contain 'no RunE and no subcommands', got %q", err.Error())
		}
	})

	t.Run("validates subcommands recursively", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{
			Use: "root",
			Commands: []Command[testConfig, NoFlags]{
				newTestSubcommand("valid-sub"),
				{Use: "invalid-sub"},
			},
		}

		err := cmd.Validate()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "subcommand 1") {
			t.Errorf("error should contain 'subcommand 1', got %q", err.Error())
		}

		if !strings.Contains(err.Error(), "invalid-sub") {
			t.Errorf("error should contain 'invalid-sub', got %q", err.Error())
		}
	})

	t.Run("error: duplicate subcommand names", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{
			Use: "root",
			Commands: []Command[testConfig, NoFlags]{
				newTestSubcommand("duplicate"),
				newTestSubcommand("duplicate"),
			},
		}

		err := cmd.Validate()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrDuplicateCommand) {
			t.Errorf("expected ErrDuplicateCommand, got %v", err)
		}

		if !strings.Contains(err.Error(), "duplicate") {
			t.Errorf("error should contain 'duplicate', got %q", err.Error())
		}
	})

	t.Run("valid with flags", func(t *testing.T) {
		t.Parallel()

		type flags struct {
			Verbose bool `default:"false" flag:"verbose"`
		}

		cmd := Command[testConfig, *flags]{
			Use:   "test",
			Flags: &flags{},
			RunE: func(_ context.Context, _ *testConfig, _ *flags) error {
				return nil
			},
		}

		err := cmd.Validate()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestCommand_HasSubcommands(t *testing.T) {
	t.Parallel()
	t.Run("returns true with subcommands", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{
			Use: "root",
			Commands: []Command[testConfig, NoFlags]{
				{Use: "sub1"},
				{Use: "sub2"},
			},
		}

		if !cmd.HasSubcommands() {
			t.Error("HasSubcommands() = false, want true")
		}
	})

	t.Run("returns false without subcommands", func(t *testing.T) {
		t.Parallel()

		cmd := newTestCommand()
		if cmd.HasSubcommands() {
			t.Error("HasSubcommands() = true, want false")
		}
	})
}

func TestCommand_HasHandler(t *testing.T) {
	t.Parallel()
	t.Run("returns true with RunE", func(t *testing.T) {
		t.Parallel()

		cmd := newTestCommand()
		if !cmd.HasHandler() {
			t.Error("HasHandler() = false, want true")
		}
	})

	t.Run("returns false without RunE", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{
			Use: "test",
		}
		if cmd.HasHandler() {
			t.Error("HasHandler() = true, want false")
		}
	})
}

func TestCommand_IsExecutable(t *testing.T) {
	t.Parallel()
	t.Run("returns true with RunE and no subcommands", func(t *testing.T) {
		t.Parallel()

		cmd := newTestCommand()
		if !cmd.IsExecutable() {
			t.Error("IsExecutable() = false, want true")
		}
	})

	t.Run("returns false without RunE", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{
			Use: "test",
		}
		if cmd.IsExecutable() {
			t.Error("IsExecutable() = true, want false")
		}
	})
}
