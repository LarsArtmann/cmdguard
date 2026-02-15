package v2

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestAppConfig struct {
	Verbose bool   `flag:"verbose" short:"v" default:"false" help:"Enable verbose output"`
	Output  string `flag:"output" short:"o" default:"-" help:"Output file"`
}

func TestNew(t *testing.T) {
	t.Run("creates GuardedCommand", func(t *testing.T) {
		defaults := TestAppConfig{}
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI application", defaults)
		require.NoError(t, err)
		require.NotNil(t, g)

		assert.Equal(t, "myapp", g.Name())
		assert.Equal(t, "My CLI application", g.Short())
		assert.Equal(t, "", g.Long())
	})

	t.Run("error: empty name", func(t *testing.T) {
		defaults := TestAppConfig{}
		g, err := New[TestAppConfig, NoFlags]("", "My CLI", defaults)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidCommand)
		assert.Nil(t, g)
	})

	t.Run("registers config in scope", func(t *testing.T) {
		defaults := TestAppConfig{Verbose: true}
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", defaults)
		require.NoError(t, err)

		cfg := g.Config()
		require.NotNil(t, cfg)
		assert.True(t, cfg.Verbose)
	})

	t.Run("creates scope", func(t *testing.T) {
		defaults := TestAppConfig{}
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", defaults)
		require.NoError(t, err)

		scope := g.ScopeStruct()
		require.NotNil(t, scope)
		assert.Equal(t, "myapp", scope.Name())
	})
}

func TestNewWithLong(t *testing.T) {
	t.Run("creates GuardedCommand with long description", func(t *testing.T) {
		defaults := TestAppConfig{}
		g, err := NewWithLong[TestAppConfig, NoFlags]("myapp", "short", "long description", defaults)
		require.NoError(t, err)
		require.NotNil(t, g)

		assert.Equal(t, "myapp", g.Name())
		assert.Equal(t, "short", g.Short())
		assert.Equal(t, "long description", g.Long())
	})

	t.Run("error: empty name", func(t *testing.T) {
		defaults := TestAppConfig{}
		g, err := NewWithLong[TestAppConfig, NoFlags]("", "short", "long", defaults)
		require.Error(t, err)
		assert.Nil(t, g)
	})
}

func TestGuardedCommand_AddCommand(t *testing.T) {
	t.Run("adds valid command", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd := Command[TestAppConfig, NoFlags]{
			Use: "greet",
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				return nil
			},
		}

		err = g.AddCommand(cmd)
		require.NoError(t, err)

		rootCmd := g.RootCommand()
		require.Len(t, rootCmd.Commands(), 1)
	})

	t.Run("error: invalid command", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd := Command[TestAppConfig, NoFlags]{
			Use: "greet",
			// No RunE
		}

		err = g.AddCommand(cmd)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrMissingHandler)
	})

	t.Run("adds command with subcommands", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		subCmd := Command[TestAppConfig, NoFlags]{
			Use: "list",
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				return nil
			},
		}

		cmd := Command[TestAppConfig, NoFlags]{
			Use:      "greet",
			Commands: []Command[TestAppConfig, NoFlags]{subCmd},
		}

		err = g.AddCommand(cmd)
		require.NoError(t, err)

		rootCmd := g.RootCommand()
		require.Len(t, rootCmd.Commands(), 1)
		greetCmd := rootCmd.Commands()[0]
		require.Len(t, greetCmd.Commands(), 1)
	})
}

func TestGuardedCommand_AddCommandFunc(t *testing.T) {
	t.Run("adds command via function", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		err = g.AddCommandFunc(func() Command[TestAppConfig, NoFlags] {
			return Command[TestAppConfig, NoFlags]{
				Use: "greet",
				RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
					return nil
				},
			}
		})
		require.NoError(t, err)

		rootCmd := g.RootCommand()
		require.Len(t, rootCmd.Commands(), 1)
	})
}

func TestGuardedCommand_Execute(t *testing.T) {
	t.Run("executes help command", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		err = g.ExecuteWithArgs(context.Background(), []string{"--help"})
		require.NoError(t, err)
	})

	t.Run("executes subcommand", func(t *testing.T) {
		executed := false
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd := Command[TestAppConfig, NoFlags]{
			Use: "greet",
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				executed = true
				return nil
			},
		}
		require.NoError(t, g.AddCommand(cmd))

		err = g.ExecuteWithArgs(context.Background(), []string{"greet"})
		require.NoError(t, err)
		assert.True(t, executed)
	})

	t.Run("error: unknown subcommand", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		// Add a valid subcommand first
		cmd := Command[TestAppConfig, NoFlags]{
			Use: "valid",
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				return nil
			},
		}
		require.NoError(t, g.AddCommand(cmd))

		// Now try an unknown subcommand - Cobra should return an error
		err = g.ExecuteWithArgs(context.Background(), []string{"unknown"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown")
	})

	t.Run("executes with flags", func(t *testing.T) {
		var receivedName string

		type GreetFlags struct {
			Name string `flag:"name" default:"World"`
		}

		g, err := New[TestAppConfig, *GreetFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd := Command[TestAppConfig, *GreetFlags]{
			Use:   "greet",
			Flags: &GreetFlags{},
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags *GreetFlags) error {
				receivedName = flags.Name
				return nil
			},
		}
		require.NoError(t, g.AddCommand(cmd))

		err = g.ExecuteWithArgs(context.Background(), []string{"greet", "--name", "Alice"})
		require.NoError(t, err)
		assert.Equal(t, "Alice", receivedName)
	})
}

func TestGuardedCommand_Scope(t *testing.T) {
	t.Run("returns injector", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		injector := g.Scope()
		require.NotNil(t, injector)
	})

	t.Run("returns scope struct", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		scope := g.ScopeStruct()
		require.NotNil(t, scope)
		assert.Equal(t, "myapp", scope.Name())
	})
}

func TestGuardedCommand_Config(t *testing.T) {
	t.Run("returns config", func(t *testing.T) {
		defaults := TestAppConfig{Verbose: true, Output: "/tmp/out"}
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", defaults)
		require.NoError(t, err)

		cfg := g.Config()
		require.NotNil(t, cfg)
		assert.True(t, cfg.Verbose)
		assert.Equal(t, "/tmp/out", cfg.Output)
	})

	t.Run("SetConfig updates config", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		newCfg := TestAppConfig{Verbose: true, Output: "/new/path"}
		g.SetConfig(newCfg)

		cfg := g.Config()
		assert.True(t, cfg.Verbose)
		assert.Equal(t, "/new/path", cfg.Output)
	})
}

func TestGuardedCommand_RootCommand(t *testing.T) {
	t.Run("returns cobra command", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		rootCmd := g.RootCommand()
		require.NotNil(t, rootCmd)
		assert.Equal(t, "myapp", rootCmd.Use)
		assert.Equal(t, "My CLI", rootCmd.Short)
	})
}

func TestGuardedCommand_Shutdown(t *testing.T) {
	t.Run("shutdown succeeds", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		err = g.Shutdown(context.Background())
		require.NoError(t, err)
	})
}

func TestGuardedCommand_HealthCheck(t *testing.T) {
	t.Run("health check succeeds", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		err = g.HealthCheck()
		require.NoError(t, err)
	})
}

func TestGuardedCommand_Metadata(t *testing.T) {
	t.Run("Name returns name", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)
		assert.Equal(t, "myapp", g.Name())
	})

	t.Run("Short returns short description", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)
		assert.Equal(t, "My CLI", g.Short())
	})

	t.Run("Long returns long description", func(t *testing.T) {
		g, err := NewWithLong[TestAppConfig, NoFlags]("myapp", "short", "long desc", TestAppConfig{})
		require.NoError(t, err)
		assert.Equal(t, "long desc", g.Long())
	})

	t.Run("SetLong updates long description", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		g.SetLong("new long description")
		assert.Equal(t, "new long description", g.Long())
		assert.Equal(t, "new long description", g.RootCommand().Long)
	})

	t.Run("SetVersion sets version", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		g.SetVersion("v1.0.0")
		assert.Equal(t, "v1.0.0", g.RootCommand().Version)
	})
}

func TestGuardedCommand_AddGlobalFlag(t *testing.T) {
	t.Run("adds global string flag", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		g.AddGlobalFlag("config", "c", "/etc/config.yaml", "Config file path")

		flag := g.RootCommand().PersistentFlags().Lookup("config")
		require.NotNil(t, flag)
		assert.Equal(t, "c", flag.Shorthand)
		assert.Equal(t, "/etc/config.yaml", flag.DefValue)
	})

	t.Run("adds global bool flag", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		g.AddGlobalBoolFlag("debug", "d", true, "Enable debug mode")

		flag := g.RootCommand().PersistentFlags().Lookup("debug")
		require.NotNil(t, flag)
		assert.Equal(t, "d", flag.Shorthand)
		assert.Equal(t, "true", flag.DefValue)
	})
}

func TestGuardedCommand_PreRunE_PostRunE(t *testing.T) {
	t.Run("calls PreRunE before RunE", func(t *testing.T) {
		var order []string

		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd := Command[TestAppConfig, NoFlags]{
			Use: "test",
			PreRunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				order = append(order, "pre")
				return nil
			},
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				order = append(order, "run")
				return nil
			},
		}
		require.NoError(t, g.AddCommand(cmd))

		err = g.ExecuteWithArgs(context.Background(), []string{"test"})
		require.NoError(t, err)
		assert.Equal(t, []string{"pre", "run"}, order)
	})

	t.Run("calls PostRunE after RunE", func(t *testing.T) {
		var order []string

		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd := Command[TestAppConfig, NoFlags]{
			Use: "test",
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				order = append(order, "run")
				return nil
			},
			PostRunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				order = append(order, "post")
				return nil
			},
		}
		require.NoError(t, g.AddCommand(cmd))

		err = g.ExecuteWithArgs(context.Background(), []string{"test"})
		require.NoError(t, err)
		assert.Equal(t, []string{"run", "post"}, order)
	})

	t.Run("PreRunE error stops execution", func(t *testing.T) {
		called := false

		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd := Command[TestAppConfig, NoFlags]{
			Use: "test",
			PreRunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				return errors.New("pre-run error")
			},
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				called = true
				return nil
			},
		}
		require.NoError(t, g.AddCommand(cmd))

		err = g.ExecuteWithArgs(context.Background(), []string{"test"})
		require.Error(t, err)
		assert.False(t, called)
	})
}

func TestGuardedCommand_CommandOptions(t *testing.T) {
	t.Run("hidden command", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd := Command[TestAppConfig, NoFlags]{
			Use:    "secret",
			Hidden: true,
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				return nil
			},
		}
		require.NoError(t, g.AddCommand(cmd))

		cobraCmd := g.RootCommand().Commands()[0]
		assert.True(t, cobraCmd.Hidden)
	})

	t.Run("deprecated command", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd := Command[TestAppConfig, NoFlags]{
			Use:        "old",
			Deprecated: "use new-cmd instead",
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				return nil
			},
		}
		require.NoError(t, g.AddCommand(cmd))

		cobraCmd := g.RootCommand().Commands()[0]
		assert.Equal(t, "use new-cmd instead", cobraCmd.Deprecated)
	})

	t.Run("command with aliases", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd := Command[TestAppConfig, NoFlags]{
			Use:      "list",
			Aliases:  []string{"ls", "l"},
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				return nil
			},
		}
		require.NoError(t, g.AddCommand(cmd))

		cobraCmd := g.RootCommand().Commands()[0]
		assert.Equal(t, []string{"ls", "l"}, cobraCmd.Aliases)
	})

	t.Run("command with version", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		cmd := Command[TestAppConfig, NoFlags]{
			Use:     "versioned",
			Version: "v1.2.3",
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags NoFlags) error {
				return nil
			},
		}
		require.NoError(t, g.AddCommand(cmd))

		cobraCmd := g.RootCommand().Commands()[0]
		assert.Equal(t, "v1.2.3", cobraCmd.Version)
	})
}

func TestGuardedCommand_Integration(t *testing.T) {
	t.Run("complete CLI workflow", func(t *testing.T) {
		type GreetFlags struct {
			Name  string `flag:"name" short:"n" default:"World" help:"Name to greet"`
			Shout bool   `flag:"shout" short:"s" default:"false" help:"Shout the greeting"`
		}

		var greetResult struct {
			name  string
			shout bool
		}

		g, err := New[TestAppConfig, *GreetFlags]("greet-cli", "A greeting CLI", TestAppConfig{})
		require.NoError(t, err)

		greetCmd := Command[TestAppConfig, *GreetFlags]{
			Use:   "greet [name]",
			Short: "Greet someone",
			Long:  "Send a greeting to the specified person.",
			Flags: &GreetFlags{},
			RunE: func(ctx context.Context, cfg *TestAppConfig, flags *GreetFlags) error {
				greetResult.name = flags.Name
				greetResult.shout = flags.Shout
				return nil
			},
		}
		require.NoError(t, g.AddCommand(greetCmd))

		err = g.ExecuteWithArgs(context.Background(), []string{"greet", "--name", "Alice", "--shout"})
		require.NoError(t, err)
		assert.Equal(t, "Alice", greetResult.name)
		assert.True(t, greetResult.shout)

		require.NoError(t, g.Shutdown(context.Background()))
	})
}

func TestCloneFlags(t *testing.T) {
	type TestFlags struct {
		Name  string
		Count int
	}

	t.Run("clones struct", func(t *testing.T) {
		original := TestFlags{Name: "test", Count: 42}
		cloned := cloneFlags(original)

		require.NotNil(t, cloned)
		// cloned is already TestFlags, no type assertion needed
		assert.Equal(t, original.Name, cloned.Name)
		assert.Equal(t, original.Count, cloned.Count)

		// Verify it's a copy (modifying clone doesn't affect original)
		cloned.Name = "modified"
		assert.Equal(t, "test", original.Name)
	})

	t.Run("clones pointer to struct", func(t *testing.T) {
		original := &TestFlags{Name: "test", Count: 42}
		cloned := cloneFlags(original)

		require.NotNil(t, cloned)
		// cloned is already *TestFlags, no type assertion needed
		assert.Equal(t, original.Name, cloned.Name)
		assert.Equal(t, original.Count, cloned.Count)

		// Verify it's a different pointer
		assert.NotSame(t, original, cloned)
	})

	t.Run("returns nil for nil pointer", func(t *testing.T) {
		var original *TestFlags // nil
		cloned := cloneFlags[*TestFlags](original)
		assert.Nil(t, cloned)
	})

	t.Run("returns as-is for non-struct", func(t *testing.T) {
		original := "string value"
		cloned := cloneFlags(original)
		assert.Equal(t, original, cloned)
	})
}

func TestFlagTypeConstraint(t *testing.T) {
	t.Run("accepts NoFlags (struct{})", func(t *testing.T) {
		err := FlagTypeConstraint[NoFlags]()
		assert.NoError(t, err)
	})

	t.Run("accepts pointer to struct", func(t *testing.T) {
		err := FlagTypeConstraint[*TestFlags]()
		assert.NoError(t, err)
	})

	t.Run("accepts empty struct", func(t *testing.T) {
		type EmptyFlags struct{}
		err := FlagTypeConstraint[EmptyFlags]()
		assert.NoError(t, err)
	})

	t.Run("accepts struct with fields", func(t *testing.T) {
		err := FlagTypeConstraint[TestFlags]()
		assert.NoError(t, err)
	})

	t.Run("rejects pointer to non-struct", func(t *testing.T) {
		err := FlagTypeConstraint[*string]()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidFlagType)
		assert.Contains(t, err.Error(), "*string")
	})

	t.Run("rejects int", func(t *testing.T) {
		err := FlagTypeConstraint[int]()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidFlagType)
		assert.Contains(t, err.Error(), "int")
	})

	t.Run("rejects string", func(t *testing.T) {
		err := FlagTypeConstraint[string]()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidFlagType)
		assert.Contains(t, err.Error(), "string")
	})

	t.Run("rejects slice", func(t *testing.T) {
		err := FlagTypeConstraint[[]string]()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidFlagType)
		assert.Contains(t, err.Error(), "[]string")
	})

	t.Run("rejects map", func(t *testing.T) {
		err := FlagTypeConstraint[map[string]string]()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidFlagType)
		assert.Contains(t, err.Error(), "map[string]string")
	})
}

func TestNew_FlagTypeValidation(t *testing.T) {
	t.Run("rejects invalid flag type in New", func(t *testing.T) {
		g, err := New[TestAppConfig, int]("myapp", "My CLI", TestAppConfig{})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidFlagType)
		assert.Nil(t, g)
	})

	t.Run("accepts NoFlags in New", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)
		require.NotNil(t, g)
	})

	t.Run("accepts pointer to struct in New", func(t *testing.T) {
		type CmdFlags struct {
			Name string `flag:"name"`
		}
		g, err := New[TestAppConfig, *CmdFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)
		require.NotNil(t, g)
	})
}
