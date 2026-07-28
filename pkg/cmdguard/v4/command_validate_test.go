package v4

import (
	"context"
	"errors"
	"testing"
)

type testConfig struct {
	Name string
}

func newTestCommand() Command[testConfig, NoFlags] {
	return Command[testConfig, NoFlags]{
		spec: commandSpec{use: "test"},
		runE: noOpHandler(),
	}
}

func newTestSubcommand(use string) Command[testConfig, NoFlags] {
	return Command[testConfig, NoFlags]{
		spec: commandSpec{use: use},
		runE: noOpHandler(),
	}
}

func TestCommand_Validate(t *testing.T) {
	t.Parallel()
	t.Run("valid command with RunE", func(t *testing.T) {
		t.Parallel()

		cmd := newTestCommand()

		err := cmd.Validate()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("valid command with subcommands", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{
			spec: commandSpec{use: "root", long: "Root command with subcommands"},
			commands: []Command[testConfig, NoFlags]{
				newTestSubcommand("sub"),
			},
		}

		err := cmd.Validate()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error: empty Use field", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{
			runE: noOpHandler(),
		}

		err := cmd.Validate()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrInvalidCommand) {
			t.Errorf("expected ErrInvalidCommand, got %v", err)
		}

		assertErrorContains(t, err, "no Use field")
	})

	t.Run("error: no RunE and no subcommands", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{
			spec: commandSpec{use: "test"},
		}

		err := cmd.Validate()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrMissingHandler) {
			t.Errorf("expected ErrMissingHandler, got %v", err)
		}

		assertErrorContains(t, err, "no RunE and no subcommands")
	})

	t.Run("validates subcommands recursively", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{
			spec: commandSpec{use: "root", long: "Root command"},
			commands: []Command[testConfig, NoFlags]{
				newTestSubcommand("valid-sub"),
				{spec: commandSpec{use: "invalid-sub"}},
			},
		}

		err := cmd.Validate()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		assertErrorContains(t, err, "subcommand 1", "invalid-sub")
	})

	t.Run("error: duplicate subcommand names", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{
			spec: commandSpec{use: "root", long: "Root command"},
			commands: []Command[testConfig, NoFlags]{
				newTestSubcommand("duplicate"),
				newTestSubcommand("duplicate"),
			},
		}

		err := cmd.Validate()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrDuplicateCommand) {
			t.Errorf("expected ErrDuplicateCommand, got %v", err)
		}

		assertErrorContains(t, err, "duplicate")
	})

	t.Run("error: parent command without Long", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{
			spec: commandSpec{use: "parent"},
			commands: []Command[testConfig, NoFlags]{
				newTestSubcommand("child"),
			},
		}

		err := cmd.Validate()
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !errors.Is(err, ErrMissingLong) {
			t.Errorf("expected ErrMissingLong, got %v", err)
		}

		assertErrorContains(t, err, "no Long description")
	})

	t.Run("valid with flags", func(t *testing.T) {
		t.Parallel()

		type flags struct {
			Verbose bool `default:"false" flag:"verbose"`
		}

		noOpRunE := func(_ context.Context, _ *testConfig, _ *flags) error {
			return nil
		}

		cmd := Command[testConfig, *flags]{
			spec:  commandSpec{use: "test"},
			flags: &flags{},
			runE:  noOpRunE,
		}

		err := cmd.Validate()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestCommand_HasSubcommands(t *testing.T) {
	t.Parallel()
	t.Run("returns true with subcommands", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{
			spec: commandSpec{use: "root", long: "Root command"},
			commands: []Command[testConfig, NoFlags]{
				{spec: commandSpec{use: "sub1"}, runE: noOpHandler()},
				{spec: commandSpec{use: "sub2"}, runE: noOpHandler()},
			},
		}

		if !cmd.HasSubcommands() {
			t.Error("HasSubcommands() = false, want true")
		}
	})

	t.Run("returns false without subcommands", func(t *testing.T) {
		t.Parallel()

		cmd := newTestCommand()

		if cmd.HasSubcommands() {
			t.Error("HasSubcommands() = true, want false")
		}
	})
}

func TestCommand_HasHandler(t *testing.T) {
	t.Parallel()
	t.Run("returns true with handler", func(t *testing.T) {
		t.Parallel()

		cmd := newTestCommand()

		if !cmd.HasHandler() {
			t.Error("HasHandler() = false, want true")
		}
	})

	t.Run("returns false without handler", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{
			spec: commandSpec{use: "test"},
		}

		if cmd.HasHandler() {
			t.Error("HasHandler() = true, want false")
		}
	})
}
