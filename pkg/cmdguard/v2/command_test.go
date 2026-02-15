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
		cmd := Command[TestConfig]{
			Use: "test",
			RunE: func(ctx context.Context, cfg *TestConfig, flags any) error {
				return nil
			},
		}
		err := cmd.Validate()
		require.NoError(t, err)
	})

	t.Run("valid command with subcommands", func(t *testing.T) {
		cmd := Command[TestConfig]{
			Use: "root",
			Commands: []Command[TestConfig]{
				{
					Use: "sub",
					RunE: func(ctx context.Context, cfg *TestConfig, flags any) error {
						return nil
					},
				},
			},
		}
		err := cmd.Validate()
		require.NoError(t, err)
	})

	t.Run("error: empty Use field", func(t *testing.T) {
		cmd := Command[TestConfig]{
			RunE: func(ctx context.Context, cfg *TestConfig, flags any) error {
				return nil
			},
		}
		err := cmd.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidCommand)
		assert.Contains(t, err.Error(), "no Use field")
	})

	t.Run("error: no RunE and no subcommands", func(t *testing.T) {
		cmd := Command[TestConfig]{
			Use: "test",
		}
		err := cmd.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrMissingHandler)
		assert.Contains(t, err.Error(), "no RunE and no subcommands")
	})

	t.Run("validates subcommands recursively", func(t *testing.T) {
		cmd := Command[TestConfig]{
			Use: "root",
			Commands: []Command[TestConfig]{
				{
					Use: "valid-sub",
					RunE: func(ctx context.Context, cfg *TestConfig, flags any) error {
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

	t.Run("valid with flags", func(t *testing.T) {
		type Flags struct {
			Verbose bool `flag:"verbose" default:"false"`
		}

		cmd := Command[TestConfig]{
			Use:   "test",
			Flags: &Flags{},
			RunE: func(ctx context.Context, cfg *TestConfig, flags any) error {
				return nil
			},
		}
		err := cmd.Validate()
		require.NoError(t, err)
	})
}

func TestCommand_HasSubcommands(t *testing.T) {
	t.Run("returns true with subcommands", func(t *testing.T) {
		cmd := Command[TestConfig]{
			Use: "test",
			Commands: []Command[TestConfig]{
				{Use: "sub"},
			},
		}
		assert.True(t, cmd.HasSubcommands())
	})

	t.Run("returns false without subcommands", func(t *testing.T) {
		cmd := Command[TestConfig]{
			Use: "test",
			RunE: func(ctx context.Context, cfg *TestConfig, flags any) error {
				return nil
			},
		}
		assert.False(t, cmd.HasSubcommands())
	})
}

func TestCommand_HasHandler(t *testing.T) {
	t.Run("returns true with RunE", func(t *testing.T) {
		cmd := Command[TestConfig]{
			Use: "test",
			RunE: func(ctx context.Context, cfg *TestConfig, flags any) error {
				return nil
			},
		}
		assert.True(t, cmd.HasHandler())
	})

	t.Run("returns false without RunE", func(t *testing.T) {
		cmd := Command[TestConfig]{
			Use: "test",
		}
		assert.False(t, cmd.HasHandler())
	})
}

func TestCommand_IsExecutable(t *testing.T) {
	t.Run("returns true with RunE", func(t *testing.T) {
		cmd := Command[TestConfig]{
			Use: "test",
			RunE: func(ctx context.Context, cfg *TestConfig, flags any) error {
				return nil
			},
		}
		assert.True(t, cmd.IsExecutable())
	})

	t.Run("returns false without RunE", func(t *testing.T) {
		cmd := Command[TestConfig]{
			Use: "test",
		}
		assert.False(t, cmd.IsExecutable())
	})
}

func TestCommandOptions(t *testing.T) {
	t.Run("WithShort", func(t *testing.T) {
		cmd := Command[TestConfig]{Use: "test"}
		WithShort[TestConfig]("short description")(&cmd)
		assert.Equal(t, "short description", cmd.Short)
	})

	t.Run("WithLong", func(t *testing.T) {
		cmd := Command[TestConfig]{Use: "test"}
		WithLong[TestConfig]("long description")(&cmd)
		assert.Equal(t, "long description", cmd.Long)
	})

	t.Run("WithAliases", func(t *testing.T) {
		cmd := Command[TestConfig]{Use: "test"}
		WithAliases[TestConfig]("alias1", "alias2")(&cmd)
		assert.Equal(t, []string{"alias1", "alias2"}, cmd.Aliases)
	})

	t.Run("WithExample", func(t *testing.T) {
		cmd := Command[TestConfig]{Use: "test"}
		WithExample[TestConfig]("example usage")(&cmd)
		assert.Equal(t, "example usage", cmd.Example)
	})

	t.Run("WithFlags", func(t *testing.T) {
		type Flags struct {
			Verbose bool `flag:"verbose"`
		}
		flags := &Flags{}
		cmd := Command[TestConfig]{Use: "test"}
		WithFlags[TestConfig](flags)(&cmd)
		assert.Equal(t, flags, cmd.Flags)
	})

	t.Run("WithRunE", func(t *testing.T) {
		handler := func(ctx context.Context, cfg *TestConfig, flags any) error {
			return nil
		}
		cmd := Command[TestConfig]{Use: "test"}
		WithRunE[TestConfig](handler)(&cmd)
		assert.NotNil(t, cmd.RunE)
	})

	t.Run("WithPreRunE", func(t *testing.T) {
		preRun := func(ctx context.Context, cfg *TestConfig, flags any) error {
			return nil
		}
		cmd := Command[TestConfig]{Use: "test"}
		WithPreRunE[TestConfig](preRun)(&cmd)
		assert.NotNil(t, cmd.PreRunE)
	})

	t.Run("WithPostRunE", func(t *testing.T) {
		postRun := func(ctx context.Context, cfg *TestConfig, flags any) error {
			return nil
		}
		cmd := Command[TestConfig]{Use: "test"}
		WithPostRunE[TestConfig](postRun)(&cmd)
		assert.NotNil(t, cmd.PostRunE)
	})

	t.Run("WithSubcommands", func(t *testing.T) {
		subCmd := Command[TestConfig]{
			Use: "sub",
			RunE: func(ctx context.Context, cfg *TestConfig, flags any) error {
				return nil
			},
		}
		cmd := Command[TestConfig]{Use: "test"}
		WithSubcommands[TestConfig](subCmd)(&cmd)
		require.Len(t, cmd.Commands, 1)
		assert.Equal(t, "sub", cmd.Commands[0].Use)
	})

	t.Run("WithHidden", func(t *testing.T) {
		cmd := Command[TestConfig]{Use: "test"}
		WithHidden[TestConfig](true)(&cmd)
		assert.True(t, cmd.Hidden)
		WithHidden[TestConfig](false)(&cmd)
		assert.False(t, cmd.Hidden)
	})

	t.Run("WithDeprecated", func(t *testing.T) {
		cmd := Command[TestConfig]{Use: "test"}
		WithDeprecated[TestConfig]("use new-cmd instead")(&cmd)
		assert.Equal(t, "use new-cmd instead", cmd.Deprecated)
	})
}

func TestNewCommand(t *testing.T) {
	t.Run("creates valid command", func(t *testing.T) {
		cmd, err := NewCommand[TestConfig]("test",
			WithShort[TestConfig]("short description"),
			WithRunE[TestConfig](func(ctx context.Context, cfg *TestConfig, flags any) error {
				return nil
			}),
		)
		require.NoError(t, err)
		assert.Equal(t, "test", cmd.Use)
		assert.Equal(t, "short description", cmd.Short)
		assert.NotNil(t, cmd.RunE)
	})

	t.Run("error: empty use", func(t *testing.T) {
		cmd, err := NewCommand[TestConfig]("",
			WithRunE[TestConfig](func(ctx context.Context, cfg *TestConfig, flags any) error {
				return nil
			}),
		)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrMissingName)
		assert.Equal(t, Command[TestConfig]{}, cmd)
	})

	t.Run("error: validation fails", func(t *testing.T) {
		cmd, err := NewCommand[TestConfig]("test") // No RunE
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrMissingHandler)
		assert.Equal(t, Command[TestConfig]{}, cmd)
	})

	t.Run("applies all options", func(t *testing.T) {
		handler := func(ctx context.Context, cfg *TestConfig, flags any) error {
			return nil
		}
		cmd, err := NewCommand[TestConfig]("test",
			WithShort[TestConfig]("short"),
			WithLong[TestConfig]("long"),
			WithAliases[TestConfig]("alias1", "alias2"),
			WithExample[TestConfig]("example"),
			WithRunE[TestConfig](handler),
			WithHidden[TestConfig](true),
			WithDeprecated[TestConfig]("deprecated"),
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
		subCmd, err := NewCommand[TestConfig]("sub",
			WithRunE[TestConfig](func(ctx context.Context, cfg *TestConfig, flags any) error {
				return nil
			}),
		)
		require.NoError(t, err)

		root, err := NewCommand[TestConfig]("root",
			WithSubcommands[TestConfig](subCmd),
		)
		require.NoError(t, err)
		assert.Equal(t, "root", root.Use)
		assert.True(t, root.HasSubcommands())
		assert.False(t, root.HasHandler())
	})
}

func TestMustNewCommand(t *testing.T) {
	t.Run("creates valid command", func(t *testing.T) {
		cmd := MustNewCommand[TestConfig]("test",
			WithRunE[TestConfig](func(ctx context.Context, cfg *TestConfig, flags any) error {
				return nil
			}),
		)
		assert.Equal(t, "test", cmd.Use)
		assert.NotNil(t, cmd.RunE)
	})

	t.Run("panics on empty use", func(t *testing.T) {
		assert.Panics(t, func() {
			MustNewCommand[TestConfig]("",
				WithRunE[TestConfig](func(ctx context.Context, cfg *TestConfig, flags any) error {
					return nil
				}),
			)
		})
	})

	t.Run("panics on validation failure", func(t *testing.T) {
		assert.Panics(t, func() {
			MustNewCommand[TestConfig]("test") // No RunE
		})
	})
}

func TestCommand_CompleteStructure(t *testing.T) {
	t.Run("complete command definition", func(t *testing.T) {
		type AppFlags struct {
			Verbose bool   `flag:"verbose" short:"v" default:"false" help:"Enable verbose output"`
			Output  string `flag:"output" short:"o" default:"-" help:"Output file"`
		}

		type GreetFlags struct {
			Name  string `flag:"name" short:"n" default:"World" help:"Name to greet"`
			Shout bool   `flag:"shout" short:"s" default:"false" help:"Shout the greeting"`
		}

		greetCmd, err := NewCommand[TestConfig]("greet",
			WithShort[TestConfig]("Greet someone"),
			WithLong[TestConfig]("Send a greeting to the specified person."),
			WithExample[TestConfig]("greet --name Alice --shout"),
			WithFlags[TestConfig](&GreetFlags{}),
			WithRunE[TestConfig](func(ctx context.Context, cfg *TestConfig, flags any) error {
				f := flags.(*GreetFlags)
				greeting := "Hello, " + f.Name
				if f.Shout {
					greeting = greeting + "!"
				}
				return nil
			}),
		)
		require.NoError(t, err)

		root, err := NewCommand[TestConfig]("myapp",
			WithShort[TestConfig]("My CLI application"),
			WithFlags[TestConfig](&AppFlags{}),
			WithSubcommands[TestConfig](greetCmd),
		)
		require.NoError(t, err)

		assert.Equal(t, "myapp", root.Use)
		assert.True(t, root.HasSubcommands())
		assert.False(t, root.HasHandler())
		assert.NotNil(t, root.Flags)

		require.Len(t, root.Commands, 1)
		assert.Equal(t, "greet", root.Commands[0].Use)
		assert.True(t, root.Commands[0].HasHandler())
	})
}
