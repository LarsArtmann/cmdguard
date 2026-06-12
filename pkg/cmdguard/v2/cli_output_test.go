package v2

import (
	"context"
	"testing"

	output "github.com/larsartmann/go-output"

	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

type outputTestConfig struct {
	Name string `flag:"name" short:"n" default:"world" help:"Name"`
}

func TestWithOutputFormat(t *testing.T) {
	t.Parallel()

	t.Run("default_format_when_not_enabled", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[outputTestConfig]("test", "test", outputTestConfig{})
		testutil.AssertNoError(t, err)

		format := cli.OutputFormat()
		testutil.AssertEqual(t, output.FormatTable, format)
	})

	t.Run("default_format_from_option", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[outputTestConfig](
			"test", "test", outputTestConfig{},
			WithOutputFormat[outputTestConfig](output.FormatJSON),
		)
		testutil.AssertNoError(t, err)

		format := cli.OutputFormat()
		testutil.AssertEqual(t, output.FormatJSON, format)
	})

	t.Run("set_format_at_runtime", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[outputTestConfig](
			"test", "test", outputTestConfig{},
			WithOutputFormat[outputTestConfig](output.FormatTable),
		)
		testutil.AssertNoError(t, err)

		cli.SetOutputFormat(output.FormatCSV)

		format := cli.OutputFormat()
		testutil.AssertEqual(t, output.FormatCSV, format)
	})

	t.Run("output_flag_parsed_from_args", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[outputTestConfig](
			"test", "test", outputTestConfig{},
			WithOutputFormat[outputTestConfig](output.FormatTable),
		)
		testutil.AssertNoError(t, err)

		resolved := ""

		cmd, err := NewCommand[outputTestConfig, NoFlags](
			"show",
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
