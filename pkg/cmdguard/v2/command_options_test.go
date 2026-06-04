package v2

import (
	"testing"

	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

func TestCommandOptions(t *testing.T) {
	t.Parallel()
	t.Run("WithShort", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{use: "test"}
		WithShort[testConfig, NoFlags]("short description")(&cmd)

		if cmd.Short() != "short description" {
			t.Errorf("Short() = %q, want %q", cmd.Short(), "short description")
		}
	})

	t.Run("WithLong", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{use: "test"}
		WithLong[testConfig, NoFlags]("long description")(&cmd)

		if cmd.Long() != "long description" {
			testutil.AssertFieldEqString(t, cmd.Long(), "long description", "Long")
		}
	})

	t.Run("WithAliases", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{use: "test"}
		WithAliases[testConfig, NoFlags]("alias1", "alias2")(&cmd)

		if len(cmd.Aliases()) != 2 || cmd.Aliases()[0] != "alias1" || cmd.Aliases()[1] != "alias2" {
			t.Errorf("Aliases() = %v, want [alias1 alias2]", cmd.Aliases())
		}
	})

	t.Run("WithExample", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{use: "test"}
		WithExample[testConfig, NoFlags]("example usage")(&cmd)

		if cmd.Example() != "example usage" {
			testutil.AssertFieldEqString(t, cmd.Example(), "example usage", "Example")
		}
	})

	t.Run("WithFlags", func(t *testing.T) {
		t.Parallel()

		type flags struct {
			Verbose bool `flag:"verbose"`
		}

		flagsInst := &flags{}
		cmd := Command[testConfig, *flags]{use: "test"}
		WithFlags[testConfig, *flags](flagsInst)(&cmd)

		if cmd.Flags() != flagsInst {
			testutil.AssertPointerEq(t, cmd.Flags(), flagsInst)
		}
	})

	t.Run("WithRunE", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{use: "test"}
		WithRunE(noOpHandler())(&cmd)

		if cmd.RunE() == nil {
			t.Error("RunE() = nil, want non-nil handler")
		}
	})

	t.Run("WithPreRunE", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{use: "test"}
		WithPreRunE(noOpHandler())(&cmd)

		if cmd.PreRunE() == nil {
			t.Error("PreRunE() = nil, want non-nil handler")
		}
	})

	t.Run("WithPostRunE", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{use: "test"}
		WithPostRunE(noOpHandler())(&cmd)

		if cmd.PostRunE() == nil {
			t.Error("PostRunE() = nil, want non-nil handler")
		}
	})

	t.Run("WithSubcommands", func(t *testing.T) {
		t.Parallel()

		subCmd := newTestCommand()
		subCmd.use = "sub"
		cmd := Command[testConfig, NoFlags]{use: "test"}
		WithSubcommands(subCmd)(&cmd)

		testutil.AssertFieldLen(t, cmd.Commands(), 1, "Commands")

		if cmd.Commands()[0].Use() != "sub" {
			testutil.AssertFieldEqString(t, cmd.Commands()[0].Use(), "sub", "Commands[0].Use")
		}
	})

	t.Run("WithHidden", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{use: "test"}
		WithHidden[testConfig, NoFlags](true)(&cmd)

		if !cmd.Hidden() {
			t.Error("Hidden() = false, want true")
		}

		WithHidden[testConfig, NoFlags](false)(&cmd)

		if cmd.Hidden() {
			t.Error("Hidden() = true, want false")
		}
	})

	t.Run("WithDeprecated", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{use: "test"}
		WithDeprecated[testConfig, NoFlags]("use new-cmd instead")(&cmd)

		if cmd.Deprecated() != "use new-cmd instead" {
			testutil.AssertFieldEqString(t, cmd.Deprecated(), "use new-cmd instead", "Deprecated")
		}
	})
}

func TestNewCommand(t *testing.T) {
	t.Parallel()
	t.Run("creates valid command", func(t *testing.T) {
		t.Parallel()

		cmd, err := NewCommand[testConfig, NoFlags](
			"test",
			noOpRunE[testConfig, NoFlags],
			WithShort[testConfig, NoFlags]("short description"),
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

		cmd, err := NewCommand[testConfig, NoFlags](
			"",
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

		_, err := NewCommand[testConfig, NoFlags]("test", nil)
		testutil.AssertExpectedError(t, err)

		testutil.AssertErrorIs(t, err, ErrMissingHandler)
	})
}

func TestCommand_CompleteStructure(t *testing.T) {
	t.Parallel()
	t.Run("creates command with all fields", func(t *testing.T) {
		t.Parallel()

		subCmd := Command[testConfig, NoFlags]{
			use:   "sub",
			short: "subcommand",
			runE:  noOpHandler(),
		}

		cmd, err := NewCommand[testConfig, NoFlags](
			"root",
			noOpHandler(),
			WithShort[testConfig, NoFlags]("root command"),
			WithLong[testConfig, NoFlags]("root command long description"),
			WithAliases[testConfig, NoFlags]("r", "root-cmd"),
			WithExample[testConfig, NoFlags]("root sub"),
			WithPreRunE(noOpHandler()),
			WithPostRunE(noOpHandler()),
			WithSubcommands(subCmd),
			WithHidden[testConfig, NoFlags](false),
			WithDeprecated[testConfig, NoFlags](""),
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
