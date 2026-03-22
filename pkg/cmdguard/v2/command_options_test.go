package v2

import (
	"context"
	"errors"
	"testing"
)

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
		_, err := NewCommand[testConfig, NoFlags]("test")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrMissingHandler) {
			t.Errorf("expected ErrMissingHandler, got %v", err)
		}
	})
}

func TestCommand_CompleteStructure(t *testing.T) {
	t.Run("creates command with all fields", func(t *testing.T) {
		subCmd := Command[testConfig, NoFlags]{
			Use:   "sub",
			Short: "subcommand",
			RunE: func(ctx context.Context, cfg *testConfig, flags NoFlags) error {
				return nil
			},
		}

		cmd, err := NewCommand[testConfig, NoFlags](
			"root",
			WithShort[testConfig, NoFlags]("root command"),
			WithLong[testConfig, NoFlags]("root command long description"),
			WithAliases[testConfig, NoFlags]("r", "root-cmd"),
			WithExample[testConfig, NoFlags]("root sub"),
			WithRunE(
				func(ctx context.Context, cfg *testConfig, flags NoFlags) error {
					return nil
				},
			),
			WithPreRunE(
				func(ctx context.Context, cfg *testConfig, flags NoFlags) error {
					return nil
				},
			),
			WithPostRunE(
				func(ctx context.Context, cfg *testConfig, flags NoFlags) error {
					return nil
				},
			),
			WithSubcommands(subCmd),
			WithHidden[testConfig, NoFlags](false),
			WithDeprecated[testConfig, NoFlags](""),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cmd.Use != "root" {
			t.Errorf("Use = %q, want %q", cmd.Use, "root")
		}

		if cmd.Short != "root command" {
			t.Errorf("Short = %q, want %q", cmd.Short, "root command")
		}

		if cmd.Long != "root command long description" {
			t.Errorf("Long = %q, want %q", cmd.Long, "root command long description")
		}

		if len(cmd.Aliases) != 2 {
			t.Errorf("len(Aliases) = %d, want 2", len(cmd.Aliases))
		}

		if cmd.Example != "root sub" {
			t.Errorf("Example = %q, want %q", cmd.Example, "root sub")
		}

		if cmd.RunE == nil {
			t.Error("RunE = nil, want non-nil")
		}

		if cmd.PreRunE == nil {
			t.Error("PreRunE = nil, want non-nil")
		}

		if cmd.PostRunE == nil {
			t.Error("PostRunE = nil, want non-nil")
		}

		if len(cmd.Commands) != 1 {
			t.Errorf("len(Commands) = %d, want 1", len(cmd.Commands))
		}

		if cmd.Hidden {
			t.Error("Hidden = true, want false")
		}

		if cmd.Deprecated != "" {
			t.Errorf("Deprecated = %q, want empty", cmd.Deprecated)
		}
	})
}
