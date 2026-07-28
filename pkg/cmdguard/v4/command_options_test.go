package v4

import (
	"context"
	"testing"

	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

func TestCommandOptions(t *testing.T) {
	t.Parallel()
	t.Run("WithShort", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{spec: commandSpec{use: "test"}}
		WithShort("short description")(&cmd.spec)

		if cmd.Short() != "short description" {
			t.Errorf("Short() = %q, want %q", cmd.Short(), "short description")
		}
	})

	t.Run("WithLong", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{spec: commandSpec{use: "test"}}
		WithLong("long description")(&cmd.spec)

		if cmd.Long() != "long description" {
			testutil.AssertFieldEqString(t, cmd.Long(), "long description", "Long")
		}
	})

	t.Run("WithAliases", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{spec: commandSpec{use: "test"}}
		WithAliases("alias1", "alias2")(&cmd.spec)

		if len(cmd.Aliases()) != 2 || cmd.Aliases()[0] != "alias1" || cmd.Aliases()[1] != "alias2" {
			t.Errorf("Aliases() = %v, want [alias1 alias2]", cmd.Aliases())
		}
	})

	t.Run("WithExample", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{spec: commandSpec{use: "test"}}
		WithExample("example usage")(&cmd.spec)

		if cmd.Example() != "example usage" {
			testutil.AssertFieldEqString(t, cmd.Example(), "example usage", "Example")
		}
	})

	t.Run("WithPreRunE", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{spec: commandSpec{use: "test"}}
		WithPreRunE(noOpHandler())(&cmd.spec)

		if cmd.PreRunE() == nil {
			t.Error("PreRunE() = nil, want non-nil handler")
		}
	})

	t.Run("WithPostRunE", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{spec: commandSpec{use: "test"}}
		WithPostRunE(noOpHandler())(&cmd.spec)

		if cmd.PostRunE() == nil {
			t.Error("PostRunE() = nil, want non-nil handler")
		}
	})

	t.Run("WithSubcommands", func(t *testing.T) {
		t.Parallel()

		subCmd := Command[testConfig, NoFlags]{
			spec: commandSpec{use: "sub"},
			runE: noOpHandler(),
		}
		cmd := Command[testConfig, NoFlags]{spec: commandSpec{use: "test"}}
		WithSubcommands(subCmd)(&cmd.spec)

		testutil.AssertFieldLen(t, cmd.Commands(), 1, "Commands")

		if cmd.Commands()[0].Use() != "sub" {
			testutil.AssertFieldEqString(t, cmd.Commands()[0].Use(), "sub", "Commands[0].Use")
		}
	})

	t.Run("WithHidden", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{spec: commandSpec{use: "test"}}
		WithHidden(true)(&cmd.spec)

		if !cmd.Hidden() {
			t.Error("Hidden() = false, want true")
		}

		WithHidden(false)(&cmd.spec)

		if cmd.Hidden() {
			t.Error("Hidden() = true, want false")
		}
	})

	t.Run("WithDeprecated", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{spec: commandSpec{use: "test"}}
		WithDeprecated("use new-cmd instead")(&cmd.spec)

		if cmd.Deprecated() != "use new-cmd instead" {
			testutil.AssertFieldEqString(t, cmd.Deprecated(), "use new-cmd instead", "Deprecated")
		}
	})
}

func TestNewCommand(t *testing.T) {
	t.Parallel()
	t.Run("creates valid command", func(t *testing.T) {
		t.Parallel()

		cmd, err := NewCommand(
			"test",
			NoFlags{},
			noOpRunE[testConfig, NoFlags],
			WithShort("short description"),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cmd.Use() != "test" {
			testutil.AssertFieldEqString(t, cmd.Use(), "test", "Use")
		}

		if cmd.Short() != "short description" {
			testutil.AssertFieldEqString(t, cmd.Short(), "short description", "Short")
		}

		if cmd.RunE() == nil {
			t.Error("RunE() = nil, want non-nil")
		}
	})

	t.Run("error: empty use", func(t *testing.T) {
		t.Parallel()

		cmd, err := NewCommand(
			"",
			NoFlags{},
			noOpHandler(),
		)
		testutil.AssertExpectedError(t, err)

		testutil.AssertErrorIs(t, err, ErrMissingName)

		if cmd.Use() != "" {
			t.Errorf("expected empty command on error, got Use()=%q", cmd.Use())
		}
	})

	t.Run("error: nil runE", func(t *testing.T) {
		t.Parallel()

		var nilHandler func(context.Context, *testConfig, NoFlags) error
		_, err := NewCommand("test", NoFlags{}, nilHandler)
		testutil.AssertExpectedError(t, err)

		testutil.AssertErrorIs(t, err, ErrMissingHandler)
	})
}

func TestCommand_CompleteStructure(t *testing.T) {
	t.Parallel()
	t.Run("creates command with all fields", func(t *testing.T) {
		t.Parallel()

		subCmd := Command[testConfig, NoFlags]{
			spec: commandSpec{use: "sub", short: "subcommand"},
			runE: noOpHandler(),
		}

		cmd, err := NewCommand(
			"root",
			NoFlags{},
			noOpHandler(),
			WithShort("root command"),
			WithLong("root command long description"),
			WithAliases("r", "root-cmd"),
			WithExample("root sub"),
			WithPreRunE(noOpHandler()),
			WithPostRunE(noOpHandler()),
			WithSubcommands(subCmd),
			WithHidden(false),
			WithDeprecated(""),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cmd.Use() != "root" {
			testutil.AssertFieldEqString(t, cmd.Use(), "root", "Use")
		}

		if cmd.Short() != "root command" {
			testutil.AssertFieldEqString(t, cmd.Short(), "root command", "Short")
		}

		if cmd.Long() != "root command long description" {
			testutil.AssertFieldEqString(t, cmd.Long(), "root command long description", "Long")
		}

		testutil.AssertFieldLen(t, cmd.Aliases(), 2, "Aliases")

		if cmd.Example() != "root sub" {
			testutil.AssertFieldEqString(t, cmd.Example(), "root sub", "Example")
		}

		if cmd.RunE() == nil {
			t.Error("RunE() = nil, want non-nil")
		}

		if cmd.PreRunE() == nil {
			t.Error("PreRunE() = nil, want non-nil")
		}

		if cmd.PostRunE() == nil {
			t.Error("PostRunE() = nil, want non-nil")
		}

		testutil.AssertFieldLen(t, cmd.Commands(), 1, "Commands")

		if cmd.Hidden() {
			t.Error("Hidden() = true, want false")
		}

		if cmd.Deprecated() != "" {
			t.Errorf("Deprecated() = %q, want empty", cmd.Deprecated())
		}
	})
}
