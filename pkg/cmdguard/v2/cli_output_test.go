package v2

import (
	"context"
	"testing"

	"github.com/larsartmann/cmdguard/pkg/testutil"
)

type outputTestConfig struct {
	Name string `flag:"name" short:"n" default:"world" help:"Name"`
}

//nolint:fatcontext // context in closures
func TestWithOutputFormat(t *testing.T) {
	t.Parallel()

	t.Run("default_format_when_not_enabled", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[outputTestConfig]("test", "test", outputTestConfig{})
		testutil.AssertNoError(t, err)

		format := cli.OutputFormat()
		testutil.AssertEqual(t, FormatTable, format)
	})

	t.Run("default_format_from_option", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[outputTestConfig]("test", "test", outputTestConfig{},
			WithOutputFormat[outputTestConfig](FormatJSON),
		)
		testutil.AssertNoError(t, err)

		format := cli.OutputFormat()
		testutil.AssertEqual(t, FormatJSON, format)
	})

	t.Run("set_format_at_runtime", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[outputTestConfig]("test", "test", outputTestConfig{},
			WithOutputFormat[outputTestConfig](FormatTable),
		)
		testutil.AssertNoError(t, err)

		cli.SetOutputFormat(FormatCSV)

		format := cli.OutputFormat()
		testutil.AssertEqual(t, FormatCSV, format)
	})

	t.Run("output_flag_parsed_from_args", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[outputTestConfig]("test", "test", outputTestConfig{},
			WithOutputFormat[outputTestConfig](FormatTable),
		)
		testutil.AssertNoError(t, err)

		resolved := ""

		cmd, err := NewCommand[outputTestConfig, NoFlags]("show",
			func(_ context.Context, _ *outputTestConfig, _ NoFlags) error {
				resolved = string(cli.OutputFormat())

				return nil
			},
			WithShort[outputTestConfig, NoFlags]("show something"),
		)
		testutil.AssertNoError(t, err)

		err = AddCommand(cli, cmd)
		testutil.AssertNoError(t, err)

		err = cli.ExecuteWithArgs(context.Background(), []string{"show", "--output", "json"})
		testutil.AssertNoError(t, err)

		testutil.AssertEqual(t, "json", resolved)
	})
}
