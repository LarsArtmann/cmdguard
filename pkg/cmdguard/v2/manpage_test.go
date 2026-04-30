package v2

import (
	"context"
	"strings"
	"testing"

	"github.com/larsartmann/cmdguard/pkg/testutil"
)

type manTestConfig struct {
	Verbose bool `flag:"verbose" short:"v" default:"false" help:"Verbose output"`
}

func TestManPage(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI[manTestConfig]("testcli", "A test CLI for man pages", manTestConfig{})
	testutil.AssertNoError(t, err)

	cmd, err := NewCommand[manTestConfig, NoFlags]("hello",
		func(_ context.Context, _ *manTestConfig, _ NoFlags) error { return nil },
		WithShort[manTestConfig, NoFlags]("Say hello"),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, cmd)
	testutil.AssertNoError(t, err)

	t.Run("generates_roff_man_page", func(t *testing.T) {
		t.Parallel()

		content, err := cli.ManPage(1)
		testutil.AssertNoError(t, err)

		if !strings.Contains(content, "testcli") {
			t.Error("man page should contain CLI name")
		}

		if !strings.Contains(content, "hello") {
			t.Error("man page should contain command name")
		}
	})

	t.Run("write_man_page", func(t *testing.T) {
		t.Parallel()

		var buf strings.Builder
		err := cli.WriteManPage(&buf, 1)
		testutil.AssertNoError(t, err)

		if buf.Len() == 0 {
			t.Error("WriteManPage should write content")
		}
	})

	t.Run("generate_man_page_command", func(t *testing.T) {
		t.Parallel()

		manCmd, err := GenerateManPageCommand(cli)
		testutil.AssertNoError(t, err)

		if manCmd.Use != "man [section]" {
			t.Errorf("expected 'man [section]', got %q", manCmd.Use)
		}
	})
}
