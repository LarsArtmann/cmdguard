package v2

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestConfig struct {
	Name string
}

func TestCommand_Validate(t *testing.T) {
	t.Run("valid command with RunE", func(t *testing.T) {
		cmd := Command[TestConfig, NoFlags]{
			Use: "test",
			RunE: func(ctx context.Context, cfg *TestConfig, flags NoFlags) error {
				return nil
			},
		}
		err := cmd.Validate()
		require.NoError(t, err)
	})

	t.Run("valid command with subcommands", func(t *testing.T) {
		cmd := Command[TestConfig, NoFlags]{
			Use: "root",
			Commands: []Command[TestConfig, NoFlags]{
				{
					Use: "sub",
					RunE: func(ctx context.Context, cfg *TestConfig, flags NoFlags) error {
						return nil
					},
				},
			},
		}
		err := cmd.Validate()
		require.NoError(t, err)
	})

	t.Run("error: empty Use field", func(t *testing.T) {
		cmd := Command[TestConfig, NoFlags]{
			RunE: func(ctx context.Context, cfg *TestConfig, flags NoFlags) error {
				return nil
			},
		}
		err := cmd.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidCommand)
		assert.Contains(t, err.Error(), "no Use field")
	})

	t.Run("error: no RunE and no subcommands", func(t *testing.T) {
		cmd := Command[TestConfig, NoFlags]{
			Use: "test",
		}
		err := cmd.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrMissingHandler)
		assert.Contains(t, err.Error(), "no RunE and no subcommands")
	})

	t.Run("validates subcommands recursively", func(t *testing.T) {
		cmd := Command[TestConfig, NoFlags]{
			Use: "root",
			Commands: []Command[TestConfig, NoFlags]{
				{
					Use: "valid-sub",
					RunE: func(ctx context.Context, cfg *TestConfig, flags NoFlags) error {
						return nil
					},
				},
				{
					Use: "invalid-sub", // No RunE
				},
			},
		}
		err := cmd.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "subcommand 1")
		assert.Contains(t, err.Error(), "invalid-sub")
	})

	t.Run("error: duplicate subcommand names", func(t *testing.T) {
		cmd := Command[TestConfig, NoFlags]{
			Use: "root",
			Commands: []Command[TestConfig, NoFlags]{
				{
					Use: "duplicate",
					RunE: func(ctx context.Context, cfg *TestConfig, flags NoFlags) error {
						return nil
					},
				},
				{
					Use: "duplicate",
					RunE: func(ctx context.Context, cfg *TestConfig, flags NoFlags) error {
						return nil
					},
				},
			},
		}
		err := cmd.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrDuplicateCommand)
		assert.Contains(t, err.Error(), "duplicate")
	})

	t.Run("valid with flags", func(t *testing.T) {
		type Flags struct {
			Verbose bool `flag:"verbose" default:"false"`
		}

		cmd := Command[TestConfig, *Flags]{
			Use:   "test",
			Flags: &Flags{},
			RunE: func(ctx context.Context, cfg *TestConfig, flags *Flags) error {
				return nil
			},
		}
		err := cmd.Validate()
		require.NoError(t, err)
	})
}

func TestCommand_HasSubcommands(t *testing.T) {
	t.Run("returns true with subcommands", func(t *testing.T) {
		cmd := Command[TestConfig, NoFlags]{
			Use: "test",
			Commands: []Command[TestConfig, NoFlags]{
				{Use: "sub"},
			},
		}
		assert.True(t, cmd.HasSubcommands())
	})

	t.Run("returns false without subcommands", func(t *testing.T) {
		cmd := Command[TestConfig, NoFlags]{
			Use: "test",
			RunE: func(ctx context.Context, cfg *TestConfig, flags NoFlags) error {
				return nil
			},
		}
		assert.False(t, cmd.HasSubcommands())
	})
}

func TestCommand_HasHandler(t *testing.T) {
	t.Run("returns true with RunE", func(t *testing.T) {
		cmd := Command[TestConfig, NoFlags]{
			Use: "test",
			RunE: func(ctx context.Context, cfg *TestConfig, flags NoFlags) error {
				return nil
			},
		}
		assert.True(t, cmd.HasHandler())
	})

	t.Run("returns false without RunE", func(t *testing.T) {
		cmd := Command[TestConfig, NoFlags]{
			Use: "test",
		}
		assert.False(t, cmd.HasHandler())
	})
}

func TestCommand_IsExecutable(t *testing.T) {
	t.Run("returns true with RunE", func(t *testing.T) {
		cmd := Command[TestConfig, NoFlags]{
			Use: "test",
			RunE: func(ctx context.Context, cfg *TestConfig, flags NoFlags) error {
				return nil
			},
		}
		assert.True(t, cmd.IsExecutable())
	})

	t.Run("returns false without RunE", func(t *testing.T) {
		cmd := Command[TestConfig, NoFlags]{
			Use: "test",
		}
		assert.False(t, cmd.IsExecutable())
	})
}

func TestCommandOptions(t *testing.T) {
	t.Run("WithShort", func(t *testing.T) {
		cmd := Command[TestConfig, NoFlags]{Use: "test"}
		WithShort[TestConfig, NoFlags]("short description")(&cmd)
		assert.Equal(t, "short description", cmd.Short)
	})

	t.Run("WithLong", func(t *testing.T) {
		cmd := Command[TestConfig, NoFlags]{Use: "test"}
		WithLong[TestConfig, NoFlags]("long description")(&cmd)
		assert.Equal(t, "long description", cmd.Long)
	})

	t.Run("WithAliases", func(t *testing.T) {
		cmd := Command[TestConfig, NoFlags]{Use: "test"}
		WithAliases[TestConfig, NoFlags]("alias1", "alias2")(&cmd)
		assert.Equal(t, []string{"alias1", "alias2"}, cmd.Aliases)
	})

	t.Run("WithExample", func(t *testing.T) {
		cmd := Command[TestConfig, NoFlags]{Use: "test"}
		WithExample[TestConfig, NoFlags]("example usage")(&cmd)
		assert.Equal(t, "example usage", cmd.Example)
	})

	t.Run("WithFlags", func(t *testing.T) {
		type Flags struct {
			Verbose bool `flag:"verbose"`
		}
		flags := &Flags{}
		cmd := Command[TestConfig, *Flags]{Use: "test"}
		WithFlags[TestConfig, *Flags](flags)(&cmd)
		assert.Equal(t, flags, cmd.Flags)
	})

	t.Run("WithRunE", func(t *testing.T) {
		handler := func(ctx context.Context, cfg *TestConfig, flags NoFlags) error {
			return nil
		}
		cmd := Command[TestConfig, NoFlags]{Use: "test"}
		WithRunE[TestConfig, NoFlags](handler)(&cmd)
		assert.NotNil(t, cmd.RunE)
	})

	t.Run("WithPreRunE", func(t *testing.T) {
		preRun := func(ctx context.Context, cfg *TestConfig, flags NoFlags) error {
			return nil
		}
		cmd := Command[TestConfig, NoFlags]{Use: "test"}
		WithPreRunE[TestConfig, NoFlags](preRun)(&cmd)
		assert.NotNil(t, cmd.PreRunE)
	})

	t.Run("WithPostRunE", func(t *testing.T) {
		postRun := func(ctx context.Context, cfg *TestConfig, flags NoFlags) error {
			return nil
		}
		cmd := Command[TestConfig, NoFlags]{Use: "test"}
		WithPostRunE[TestConfig, NoFlags](postRun)(&cmd)
		assert.NotNil(t, cmd.PostRunE)
	})

	t.Run("WithSubcommands", func(t *testing.T) {
		subCmd := Command[TestConfig, NoFlags]{
			Use: "sub",
			RunE: func(ctx context.Context, cfg *TestConfig, flags NoFlags) error {
				return nil
			},
		}
		cmd := Command[TestConfig, NoFlags]{Use: "test"}
		WithSubcommands[TestConfig, NoFlags](subCmd)(&cmd)
		require.Len(t, cmd.Commands, 1)
		assert.Equal(t, "sub", cmd.Commands[0].Use)
	})

	t.Run("WithHidden", func(t *testing.T) {
		cmd := Command[TestConfig, NoFlags]{Use: "test"}
		WithHidden[TestConfig, NoFlags](true)(&cmd)
		assert.True(t, cmd.Hidden)
		WithHidden[TestConfig, NoFlags](false)(&cmd)
		assert.False(t, cmd.Hidden)
	})

	t.Run("WithDeprecated", func(t *testing.T) {
		cmd := Command[TestConfig, NoFlags]{Use: "test"}
		WithDeprecated[TestConfig, NoFlags]("use new-cmd instead")(&cmd)
		assert.Equal(t, "use new-cmd instead", cmd.Deprecated)
	})
}

func TestNewCommand(t *testing.T) {
	t.Run("creates valid command", func(t *testing.T) {
		cmd, err := NewCommand[TestConfig, NoFlags]("test",
			WithShort[TestConfig, NoFlags]("short description"),
			WithRunE[TestConfig, NoFlags](func(ctx context.Context, cfg *TestConfig, flags NoFlags) error {
				return nil
			}),
		)
		require.NoError(t, err)
		assert.Equal(t, "test", cmd.Use)
		assert.Equal(t, "short description", cmd.Short)
		assert.NotNil(t, cmd.RunE)
	})

	t.Run("error: empty use", func(t *testing.T) {
		cmd, err := NewCommand[TestConfig, NoFlags]("",
			WithRunE[TestConfig, NoFlags](func(ctx context.Context, cfg *TestConfig, flags NoFlags) error {
				return nil
			}),
		)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrMissingName)
		assert.Equal(t, Command[TestConfig, NoFlags]{}, cmd)
	})

	t.Run("error: validation fails", func(t *testing.T) {
		cmd, err := NewCommand[TestConfig, NoFlags]("test") // No RunE
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrMissingHandler)
		assert.Equal(t, Command[TestConfig, NoFlags]{}, cmd)
	})

	t.Run("applies all options", func(t *testing.T) {
		handler := func(ctx context.Context, cfg *TestConfig, flags NoFlags) error {
			return nil
		}
		cmd, err := NewCommand[TestConfig, NoFlags]("test",
			WithShort[TestConfig, NoFlags]("short"),
			WithLong[TestConfig, NoFlags]("long"),
			WithAliases[TestConfig, NoFlags]("alias1", "alias2"),
			WithExample[TestConfig, NoFlags]("example"),
			WithRunE[TestConfig, NoFlags](handler),
			WithHidden[TestConfig, NoFlags](true),
			WithDeprecated[TestConfig, NoFlags]("deprecated"),
		)
		require.NoError(t, err)
		assert.Equal(t, "test", cmd.Use)
		assert.Equal(t, "short", cmd.Short)
		assert.Equal(t, "long", cmd.Long)
		assert.Equal(t, []string{"alias1", "alias2"}, cmd.Aliases)
		assert.Equal(t, "example", cmd.Example)
		assert.NotNil(t, cmd.RunE)
		assert.True(t, cmd.Hidden)
		assert.Equal(t, "deprecated", cmd.Deprecated)
	})

	t.Run("creates command with subcommands", func(t *testing.T) {
		subCmd, err := NewCommand[TestConfig, NoFlags]("sub",
			WithRunE[TestConfig, NoFlags](func(ctx context.Context, cfg *TestConfig, flags NoFlags) error {
				return nil
			}),
		)
		require.NoError(t, err)

		root, err := NewCommand[TestConfig, NoFlags]("root",
			WithSubcommands[TestConfig, NoFlags](subCmd),
		)
		require.NoError(t, err)
		assert.Equal(t, "root", root.Use)
		assert.True(t, root.HasSubcommands())
		assert.False(t, root.HasHandler())
	})
}

func TestNewCommand_ErrorCases(t *testing.T) {
	t.Run("returns error on empty use", func(t *testing.T) {
		_, err := NewCommand[TestConfig, NoFlags]("",
			WithRunE[TestConfig, NoFlags](func(ctx context.Context, cfg *TestConfig, flags NoFlags) error {
				return nil
			}),
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "use is required")
	})

	t.Run("returns error on validation failure", func(t *testing.T) {
		_, err := NewCommand[TestConfig, NoFlags]("test") // No RunE
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no handler")
	})
}

func TestCommand_CompleteStructure(t *testing.T) {
	t.Run("complete command definition", func(t *testing.T) {
		type AppFlags struct {
			Verbose bool   `flag:"verbose" short:"v" default:"false" help:"Enable verbose output"`
			Output  string `flag:"output" short:"o" default:"-" help:"Output file"`
		}

		subCmd, err := NewCommand[TestConfig, *AppFlags]("sub",
			WithShort[TestConfig, *AppFlags]("A subcommand"),
			WithRunE[TestConfig, *AppFlags](func(ctx context.Context, cfg *TestConfig, flags *AppFlags) error {
				return nil
			}),
		)
		require.NoError(t, err)

		root, err := NewCommand[TestConfig, *AppFlags]("myapp",
			WithShort[TestConfig, *AppFlags]("My CLI application"),
			WithFlags[TestConfig, *AppFlags](&AppFlags{}),
			WithSubcommands[TestConfig, *AppFlags](subCmd),
		)
		require.NoError(t, err)

		assert.Equal(t, "myapp", root.Use)
		assert.True(t, root.HasSubcommands())
		assert.False(t, root.HasHandler())
		assert.NotNil(t, root.Flags)
	})
}
