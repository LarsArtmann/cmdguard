package v4

import (
	"context"
	"testing"
)

func TestCLI_Integration(t *testing.T) {
	t.Parallel()
	t.Run("complete CLI workflow", func(t *testing.T) {
		t.Parallel()
		type greetFlags struct {
			Name  string `default:"World" flag:"name"  help:"Name to greet"      short:"n"`
			Shout bool   `default:"false" flag:"shout" help:"Shout the greeting" short:"s"`
		}

		var greetResult struct {
			name  string
			shout bool
		}

		cli, err := NewCLI("greet-cli", "A greeting CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		greetCmd := Command[testAppConfig, *greetFlags]{
			spec: commandSpec{
				use:   "greet [name]",
				short: "Greet someone",
				long:  "Send a greeting to the specified person.",
			},
			flags: &greetFlags{},
			runE: func(_ context.Context, _ *testAppConfig, flags *greetFlags) error {
				greetResult.name = flags.Name
				greetResult.shout = flags.Shout

				return nil
			},
		}
		if err := AddCommand(cli, greetCmd); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = cli.ExecuteWithArgs(
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

		if err := cli.Shutdown(t.Context()); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
