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
		RunE: func(ctx context.Context, cfg *testConfig, flags NoFlags) error {
			return nil
		},
	}
}

func TestCommand_Validate(t *testing.T) {
	t.Run("valid command with RunE", func(t *testing.T) {
		cmd := Command[testConfig, NoFlags]{
			Use: "test",
			RunE: func(ctx context.Context, cfg *testConfig, flags NoFlags) error {
				return nil
			},
		}

		err := cmd.Validate()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("valid command with subcommands", func(t *testing.T) {
		cmd := Command[testConfig, NoFlags]{
			Use: "root",
			Commands: []Command[testConfig, NoFlags]{
				{
					Use: "sub",
					RunE: func(ctx context.Context, cfg *testConfig, flags NoFlags) error {
						return nil
					},
				},
			},
		}

		err := cmd.Validate()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error: empty Use field", func(t *testing.T) {
		cmd := Command[testConfig, NoFlags]{
			RunE: func(ctx context.Context, cfg *testConfig, flags NoFlags) error {
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
		cmd := Command[testConfig, NoFlags]{
			Use: "root",
			Commands: []Command[testConfig, NoFlags]{
				{
					Use: "valid-sub",
					RunE: func(ctx context.Context, cfg *testConfig, flags NoFlags) error {
						return nil
					},
				},
				{
					Use: "invalid-sub",
				},
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
		cmd := Command[testConfig, NoFlags]{
			Use: "root",
			Commands: []Command[testConfig, NoFlags]{
				{
					Use: "duplicate",
					RunE: func(ctx context.Context, cfg *testConfig, flags NoFlags) error {
						return nil
					},
				},
				{
					Use: "duplicate",
					RunE: func(ctx context.Context, cfg *testConfig, flags NoFlags) error {
						return nil
					},
				},
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
		type flags struct {
			Verbose bool `default:"false" flag:"verbose"`
		}

		cmd := Command[testConfig, *flags]{
			Use:   "test",
			Flags: &flags{},
			RunE: func(ctx context.Context, cfg *testConfig, flags *flags) error {
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
	t.Run("returns true with subcommands", func(t *testing.T) {
		cmd := Command[testConfig, NoFlags]{
			Use: "test",
			Commands: []Command[testConfig, NoFlags]{
				{Use: "sub"},
			},
		}
		if !cmd.HasSubcommands() {
			t.Error("HasSubcommands() = false, want true")
		}
	})

	t.Run("returns false without subcommands", func(t *testing.T) {
		cmd := newTestCommand()
		if cmd.HasSubcommands() {
			t.Error("HasSubcommands() = true, want false")
		}
	})
}

func TestCommand_HasHandler(t *testing.T) {
	t.Run("returns true with RunE", func(t *testing.T) {
		cmd := newTestCommand()
		if !cmd.HasHandler() {
			t.Error("HasHandler() = false, want true")
		}
	})

	t.Run("returns false without RunE", func(t *testing.T) {
		cmd := Command[testConfig, NoFlags]{
			Use: "test",
		}
		if cmd.HasHandler() {
			t.Error("HasHandler() = true, want false")
		}
	})
}

func TestCommand_IsExecutable(t *testing.T) {
	t.Run("returns true with RunE", func(t *testing.T) {
		cmd := newTestCommand()
		if !cmd.IsExecutable() {
			t.Error("IsExecutable() = false, want true")
		}
	})

	t.Run("returns false without RunE", func(t *testing.T) {
		cmd := Command[testConfig, NoFlags]{
			Use: "test",
		}
		if cmd.IsExecutable() {
			t.Error("IsExecutable() = true, want false")
		}
	})
}

func TestCommandOptions(t *testing.T) {
	t.Run("WithShort", func(t *testing.T) {
		cmd := Command[testConfig, NoFlags]{Use: "test"}
		WithShort[testConfig, NoFlags]("short description")(&cmd)

		if cmd.Short != "short description" {
			t.Errorf("Short = %q, want %q", cmd.Short, "short description")
		}
	})

	t.Run("WithLong", func(t *testing.T) {
		cmd := Command[testConfig, NoFlags]{Use: "test"}
		WithLong[testConfig, NoFlags]("long description")(&cmd)

		if cmd.Long != "long description" {
			t.Errorf("Long = %q, want %q", cmd.Long, "long description")
		}
	})

	t.Run("WithAliases", func(t *testing.T) {
		cmd := Command[testConfig, NoFlags]{Use: "test"}
		WithAliases[testConfig, NoFlags]("alias1", "alias2")(&cmd)

		if len(cmd.Aliases) != 2 || cmd.Aliases[0] != "alias1" || cmd.Aliases[1] != "alias2" {
			t.Errorf("Aliases = %v, want [alias1 alias2]", cmd.Aliases)
		}
	})

	t.Run("WithExample", func(t *testing.T) {
		cmd := Command[testConfig, NoFlags]{Use: "test"}
		WithExample[testConfig, NoFlags]("example usage")(&cmd)

		if cmd.Example != "example usage" {
			t.Errorf("Example = %q, want %q", cmd.Example, "example usage")
		}
	})

	t.Run("WithFlags", func(t *testing.T) {
		type flags struct {
			Verbose bool `flag:"verbose"`
		}

		flagsInst := &flags{}
		cmd := Command[testConfig, *flags]{Use: "test"}
		WithFlags[testConfig, *flags](flagsInst)(&cmd)

		if cmd.Flags != flagsInst {
			t.Errorf("Flags = %p, want %p", cmd.Flags, flagsInst)
		}
	})

	t.Run("WithRunE", func(t *testing.T) {
		cmd := Command[testConfig, NoFlags]{Use: "test"}
		handler := func(ctx context.Context, cfg *testConfig, flags NoFlags) error {
			return nil
		}
		WithRunE(handler)(&cmd)

		if cmd.RunE == nil {
			t.Error("RunE = nil, want non-nil handler")
		}
	})

	t.Run("WithPreRunE", func(t *testing.T) {
		cmd := Command[testConfig, NoFlags]{Use: "test"}
		handler := func(ctx context.Context, cfg *testConfig, flags NoFlags) error {
			return nil
		}
		WithPreRunE(handler)(&cmd)

		if cmd.PreRunE == nil {
			t.Error("PreRunE = nil, want non-nil handler")
		}
	})

	t.Run("WithPostRunE", func(t *testing.T) {
		cmd := Command[testConfig, NoFlags]{Use: "test"}
		handler := func(ctx context.Context, cfg *testConfig, flags NoFlags) error {
			return nil
		}
		WithPostRunE(handler)(&cmd)

		if cmd.PostRunE == nil {
			t.Error("PostRunE = nil, want non-nil handler")
		}
	})

	t.Run("WithSubcommands", func(t *testing.T) {
		subCmd := Command[testConfig, NoFlags]{
			Use: "sub",
			RunE: func(ctx context.Context, cfg *testConfig, flags NoFlags) error {
				return nil
			},
		}
		cmd := Command[testConfig, NoFlags]{Use: "test"}
		WithSubcommands(subCmd)(&cmd)

		if len(cmd.Commands) != 1 {
			t.Errorf("len(Commands) = %d, want 1", len(cmd.Commands))
		}

		if cmd.Commands[0].Use != "sub" {
			t.Errorf("Commands[0].Use = %q, want %q", cmd.Commands[0].Use, "sub")
		}
	})

	t.Run("WithHidden", func(t *testing.T) {
		cmd := Command[testConfig, NoFlags]{Use: "test"}
		WithHidden[testConfig, NoFlags](true)(&cmd)

		if !cmd.Hidden {
			t.Error("Hidden = false, want true")
		}

		WithHidden[testConfig, NoFlags](false)(&cmd)

		if cmd.Hidden {
			t.Error("Hidden = true, want false")
		}
	})

	t.Run("WithDeprecated", func(t *testing.T) {
		cmd := Command[testConfig, NoFlags]{Use: "test"}
		WithDeprecated[testConfig, NoFlags]("use new-cmd instead")(&cmd)

		if cmd.Deprecated != "use new-cmd instead" {
			t.Errorf("Deprecated = %q, want %q", cmd.Deprecated, "use new-cmd instead")
		}
	})
}

func TestNewCommand(t *testing.T) {
	t.Run("creates valid command", func(t *testing.T) {
		cmd, err := NewCommand[testConfig, NoFlags](
			"test",
			WithShort[testConfig, NoFlags]("short description"),
			WithRunE(
				func(ctx context.Context, cfg *testConfig, flags NoFlags) error {
					return nil
				},
			),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cmd.Use != "test" {
			t.Errorf("Use = %q, want %q", cmd.Use, "test")
		}

		if cmd.Short != "short description" {
			t.Errorf("Short = %q, want %q", cmd.Short, "short description")
		}

		if cmd.RunE == nil {
			t.Error("RunE = nil, want non-nil")
		}
	})

	t.Run("error: empty use", func(t *testing.T) {
		cmd, err := NewCommand[testConfig, NoFlags](
			"",
			WithRunE(
				func(ctx context.Context, cfg *testConfig, flags NoFlags) error {
					return nil
				},
			),
		)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrMissingName) {
			t.Errorf("expected ErrMissingName, got %v", err)
		}

		if cmd.Use != "" {
			t.Errorf("expected empty command on error, got Use=%q", cmd.Use)
		}
	})

	t.Run("error: validation fails", func(t *testing.T) {
		cmd, err := NewCommand[testConfig, NoFlags]("test")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrMissingHandler) {
			t.Errorf("expected ErrMissingHandler, got %v", err)
		}

		if cmd.Use != "" {
			t.Errorf("expected empty command on error, got Use=%q", cmd.Use)
		}
	})

	t.Run("applies all options", func(t *testing.T) {
		handler := func(ctx context.Context, cfg *testConfig, flags NoFlags) error {
			return nil
		}

		cmd, err := NewCommand[testConfig, NoFlags]("test",
			WithShort[testConfig, NoFlags]("short"),
			WithLong[testConfig, NoFlags]("long"),
			WithAliases[testConfig, NoFlags]("alias1", "alias2"),
			WithExample[testConfig, NoFlags]("example"),
			WithRunE(handler),
			WithHidden[testConfig, NoFlags](true),
			WithDeprecated[testConfig, NoFlags]("deprecated"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cmd.Use != "test" {
			t.Errorf("Use = %q, want %q", cmd.Use, "test")
		}

		if cmd.Short != "short" {
			t.Errorf("Short = %q, want %q", cmd.Short, "short")
		}

		if cmd.Long != "long" {
			t.Errorf("Long = %q, want %q", cmd.Long, "long")
		}

		if len(cmd.Aliases) != 2 || cmd.Aliases[0] != "alias1" || cmd.Aliases[1] != "alias2" {
			t.Errorf("Aliases = %v, want [alias1 alias2]", cmd.Aliases)
		}

		if cmd.Example != "example" {
			t.Errorf("Example = %q, want %q", cmd.Example, "example")
		}

		if cmd.RunE == nil {
			t.Error("RunE = nil, want non-nil")
		}

		if !cmd.Hidden {
			t.Error("Hidden = false, want true")
		}

		if cmd.Deprecated != "deprecated" {
			t.Errorf("Deprecated = %q, want %q", cmd.Deprecated, "deprecated")
		}
	})

	t.Run("creates command with subcommands", func(t *testing.T) {
		subCmd, err := NewCommand[testConfig, NoFlags](
			"sub",
			WithRunE(
				func(ctx context.Context, cfg *testConfig, flags NoFlags) error {
					return nil
				},
			),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		root, err := NewCommand[testConfig, NoFlags]("root",
			WithSubcommands(subCmd),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if root.Use != "root" {
			t.Errorf("Use = %q, want %q", root.Use, "root")
		}

		if !root.HasSubcommands() {
			t.Error("HasSubcommands() = false, want true")
		}

		if root.HasHandler() {
			t.Error("HasHandler() = true, want false")
		}
	})
}

func TestCommand_CompleteStructure(t *testing.T) {
	t.Run("complete command definition", func(t *testing.T) {
		type appFlags struct {
			Verbose bool   `default:"false" flag:"verbose" help:"Enable verbose output" short:"v"`
			Output  string `default:"-"     flag:"output"  help:"Output file"           short:"o"`
		}

		subCmd, err := NewCommand[testConfig, *appFlags](
			"sub",
			WithShort[testConfig, *appFlags]("A subcommand"),
			WithRunE(
				func(ctx context.Context, cfg *testConfig, flags *appFlags) error {
					return nil
				},
			),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		root, err := NewCommand[testConfig, *appFlags]("myapp",
			WithShort[testConfig, *appFlags]("My CLI application"),
			WithFlags[testConfig, *appFlags](&appFlags{}),
			WithSubcommands[testConfig, *appFlags](subCmd),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if root.Use != "myapp" {
			t.Errorf("Use = %q, want %q", root.Use, "myapp")
		}

		if !root.HasSubcommands() {
			t.Error("HasSubcommands() = false, want true")
		}

		if root.HasHandler() {
			t.Error("HasHandler() = true, want false")
		}

		if root.Flags == nil {
			t.Error("Flags = nil, want non-nil")
		}
	})
}
