package v2

import (
	"context"
	"testing"
)

func TestGuardedCommand_Integration(t *testing.T) {
	t.Run("complete CLI workflow", func(t *testing.T) {
		type greetFlags struct {
			Name  string `default:"World" flag:"name"  help:"Name to greet"      short:"n"`
			Shout bool   `default:"false" flag:"shout" help:"Shout the greeting" short:"s"`
		}

		var greetResult struct {
			name  string
			shout bool
		}

		g, err := New[testAppConfig, *greetFlags]("greet-cli", "A greeting CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		greetCmd := Command[testAppConfig, *greetFlags]{
			Use:   "greet [name]",
			Short: "Greet someone",
			Long:  "Send a greeting to the specified person.",
			Flags: &greetFlags{},
			RunE: func(_ context.Context, _ *testAppConfig, flags *greetFlags) error {
				greetResult.name = flags.Name
				greetResult.shout = flags.Shout

				return nil
			},
		}
		if err := g.AddCommand(greetCmd); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = g.ExecuteWithArgs(
			t.Context(),
			[]string{"greet", "--name", "Alice", "--shout"},
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if greetResult.name != "Alice" {
			t.Errorf("greetResult.name = %q, want %q", greetResult.name, "Alice")
		}

		if !greetResult.shout {
			t.Error("greetResult.shout = false, want true")
		}

		if err := g.Shutdown(t.Context()); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
