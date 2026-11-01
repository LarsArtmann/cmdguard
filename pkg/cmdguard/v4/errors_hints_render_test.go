package v4

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"

	"charm.land/fang/v2"
	output "github.com/larsartmann/go-output"

	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

// newFailingCLI builds a one-command CLI whose handler fails with a hinted
// sentinel error.
func newFailingCLI(t *testing.T, opts ...CLIOption) *CLI[testConfig] {
	t.Helper()

	cliOpts := append([]CLIOption{WithFang(false)}, opts...)

	cli, err := NewCLI("app", "App", testConfig{}, cliOpts...)
	testutil.AssertNoError(t, err)

	cmd, err := NewCommand(
		"fail",
		NoFlags{},
		func(_ context.Context, _ *testConfig, _ NoFlags) error {
			return fmt.Errorf("%w: name", ErrRequiredFlag)
		},
		WithShort("Fail on purpose"),
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, cmd))

	return cli
}

func TestExecute_PrintsHintWhenFangDisabled(t *testing.T) {
	t.Parallel()

	cli := newFailingCLI(t)

	var stderr bytes.Buffer

	cli.RootCommand().SetErr(&stderr)

	err := cli.ExecuteWithArgs(context.Background(), []string{"fail"})
	if err == nil {
		t.Fatal("expected the command error")
	}

	assertStringContains(t, stderr.String(), "Hint: "+ErrorHint(ErrRequiredFlag))
}

func TestExecute_SuppressesHintWhenDisabled(t *testing.T) {
	t.Parallel()

	cli := newFailingCLI(t, WithoutErrorHints())

	var stderr bytes.Buffer

	cli.RootCommand().SetErr(&stderr)

	err := cli.ExecuteWithArgs(context.Background(), []string{"fail"})
	if err == nil {
		t.Fatal("expected the command error")
	}

	assertStringContains(t, stderr.String(), "required flag not set")

	if bytes.Contains(stderr.Bytes(), []byte("Hint: ")) {
		t.Errorf("expected no hint line, got: %s", stderr.String())
	}
}

func TestWriteErrorWithHint_AppendsHint(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	writeErrorWithHint(&buf, fang.Styles{}, fmt.Errorf("%w: name", ErrRequiredFlag))

	assertStringContains(t, buf.String(), "required flag not set")
	assertStringContains(t, buf.String(), "Hint: "+ErrorHint(ErrRequiredFlag))
}

func TestWriteErrorWithHint_NoHintForUnknownError(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	writeErrorWithHint(&buf, fang.Styles{}, errors.New("something unexpected"))

	assertStringContains(t, buf.String(), "something unexpected")

	if bytes.Contains(buf.Bytes(), []byte("Hint: ")) {
		t.Errorf("expected no hint for an unknown error, got: %s", buf.String())
	}
}

func TestWritePlainErrorHint_SkipsMachineReadableFormat(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI(
		"app", "App", testConfig{},
		WithFang(false),
		WithOutputFormat(output.FormatJSON),
	)
	testutil.AssertNoError(t, err)

	var stderr bytes.Buffer

	cli.RootCommand().SetErr(&stderr)
	cli.writePlainErrorHint(fmt.Errorf("%w: name", ErrRequiredFlag))

	if stderr.Len() != 0 {
		t.Errorf("expected no plain hint for JSON output, got: %s", stderr.String())
	}
}

func TestIsMachineReadableFormat(t *testing.T) {
	t.Parallel()

	readable := []OutputFormat{
		output.FormatJSON,
		output.FormatJSONL,
		output.FormatYAML,
		output.FormatTOML,
	}

	for _, format := range readable {
		if !isMachineReadableFormat(format) {
			t.Errorf("expected %q to be machine-readable", format)
		}
	}

	notReadable := []OutputFormat{"", output.FormatTable, output.FormatCSV}

	for _, format := range notReadable {
		if isMachineReadableFormat(format) {
			t.Errorf("expected %q not to be machine-readable", format)
		}
	}
}
