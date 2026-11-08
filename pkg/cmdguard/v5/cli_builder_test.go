package v5

import (
	"errors"
	"strings"
	"testing"

	"github.com/larsartmann/cmdguard/v5/pkg/testutil"
)

func TestBuilder_CollectsAndBuildsAtomically(t *testing.T) {
	t.Parallel()

	t.Run("Add success attaches every command", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI("app", "App", testConfig{})
		testutil.AssertNoError(t, err)

		b := NewBuilder(cli)
		Add(b, "one", NoFlags{}, noOpHandler(), WithShort("one"))
		Add(b, "two", NoFlags{}, noOpHandler(), WithShort("two"))
		testutil.AssertNoError(t, b.Build())

		got := strings.Join(commandNames(cli), " ")
		testutil.AssertFieldEqString(t, got, "one two", "attached commands")
	})

	t.Run("AddParent queues a parent command", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI("app", "App", testConfig{})
		testutil.AssertNoError(t, err)

		child, err := NewCommand("leaf", NoFlags{}, noOpHandler(), WithShort("leaf"))
		testutil.AssertNoError(t, err)

		b := NewBuilder(cli)
		AddParent(b, "group", "Group of commands", NoFlags{}, WithShort("group"), WithSubcommands(child))
		testutil.AssertNoError(t, b.Build())

		got := strings.Join(commandNames(cli), " ")
		testutil.AssertFieldEqString(t, got, "group", "attached commands")
	})

	t.Run("construction errors are joined and nothing attaches", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI("app", "App", testConfig{})
		testutil.AssertNoError(t, err)

		b := NewBuilder(cli)
		Add(b, "fine", NoFlags{}, noOpHandler(), WithShort("fine"))
		Add(b, "bad", NoFlags{}, noOpHandler())                       // missing short under strict
		Add(b, "worse", NoFlags{}, nil, WithShort("worse"))           // missing handler
		Try(b, Command[testConfig, NoFlags]{}, errors.New("pair err")) // pre-bound pair

		err = b.Build()
		testutil.AssertExpectedError(t, err)

		if !errors.Is(err, ErrMissingShort) ||
			!errors.Is(err, ErrMissingHandler) ||
			!strings.Contains(err.Error(), "pair err") {
			t.Errorf("joined error must carry every failure, got: %v", err)
		}

		if len(commandNames(cli)) != 0 {
			t.Error("a failed Build must leave the CLI unchanged")
		}
	})

	t.Run("Try queues a pre-bound pair", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI("app", "App", testConfig{})
		testutil.AssertNoError(t, err)

		cmd, err := NewCommand("held", NoFlags{}, noOpHandler(), WithShort("held"))
		testutil.AssertNoError(t, err)

		b := NewBuilder(cli)
		Try(b, cmd, err)
		testutil.AssertNoError(t, b.Build())

		got := strings.Join(commandNames(cli), " ")
		testutil.AssertFieldEqString(t, got, "held", "attached commands")
	})

	t.Run("Build reports duplicate names via AddCommands", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI("app", "App", testConfig{})
		testutil.AssertNoError(t, err)

		b := NewBuilder(cli)
		Add(b, "dupe", NoFlags{}, noOpHandler(), WithShort("dupe"))
		Add(b, "dupe", NoFlags{}, noOpHandler(), WithShort("dupe"))

		err = b.Build()
		testutil.AssertExpectedError(t, err)

		if !errors.Is(err, ErrDuplicateCommand) {
			t.Errorf("expected ErrDuplicateCommand, got: %v", err)
		}
	})

	t.Run("nil Builder calls are no-ops", func(t *testing.T) {
		t.Parallel()

		var b *Builder[testConfig]

		Add(b, "x", NoFlags{}, noOpHandler(), WithShort("x"))
		AddParent(b, "y", "long", NoFlags{})
		Try(b, Command[testConfig, NoFlags]{}, nil)
	})
}
