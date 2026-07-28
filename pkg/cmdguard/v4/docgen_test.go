package v4

import (
	"context"
	"strings"
	"testing"

	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

func TestGenerateDocs_ProducesMarkdown(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Name string `flag:"name" default:"world" help:"name to greet"`
	}

	cli, err := NewCLI("myapp", "1.0", cfg{})
	testutil.AssertNoError(t, err)

	cmd, err := NewCommand(
		"greet",
		NoFlags{},
		func(_ context.Context, _ *cfg, _ NoFlags) error {
			return nil
		},
		WithShort("greets someone"),
		WithLong("A longer description of the greet command that spans a sentence."),
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, cmd))

	var buf strings.Builder

	err = cli.GenerateDocs(&buf)
	testutil.AssertNoError(t, err)

	output := buf.String()

	if !strings.Contains(output, "myapp") {
		t.Error("expected docs to contain root command 'myapp'")
	}

	if !strings.Contains(output, "greet") {
		t.Error("expected docs to contain 'greet' command")
	}

	if !strings.Contains(output, "### Usage") {
		t.Error("expected docs to contain Usage section")
	}

	if !strings.Contains(output, "### Flags") {
		t.Error("expected docs to contain Flags section")
	}

	if !strings.Contains(output, "--name") {
		t.Error("expected docs to contain '--name' flag")
	}
}

func TestGenerateDocs_HiddenCommandSkipped(t *testing.T) {
	t.Parallel()

	type cfg struct{}

	cli, err := NewCLI("root", "1.0", cfg{})
	testutil.AssertNoError(t, err)

	visible, err := NewCommand(
		"visible",
		NoFlags{},
		func(_ context.Context, _ *cfg, _ NoFlags) error {
			return nil
		},
		WithShort("visible command"),
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, visible))

	hidden, err := NewCommand(
		"secret",
		NoFlags{},
		func(_ context.Context, _ *cfg, _ NoFlags) error {
			return nil
		},
		WithShort("hidden command"),
		WithHidden(true),
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, hidden))

	var buf strings.Builder

	err = cli.GenerateDocs(&buf)
	testutil.AssertNoError(t, err)

	output := buf.String()

	if !strings.Contains(output, "visible") {
		t.Error("expected docs to contain 'visible' command")
	}

	if strings.Contains(output, "secret") {
		t.Error("expected docs to NOT contain hidden 'secret' command")
	}
}
