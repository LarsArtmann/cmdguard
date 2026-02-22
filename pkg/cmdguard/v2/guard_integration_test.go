package v2

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
