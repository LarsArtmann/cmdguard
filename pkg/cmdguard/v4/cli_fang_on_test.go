package v4

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"fmt"
	"os"
	"strings"
	"testing"

	output "github.com/larsartmann/go-output"

	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

// captureOSOutput runs fn with os.Stdout and os.Stderr replaced by pipes and
// returns everything written to each stream. Fang renders to the process
// streams directly, so cmdguard's fang-on paths can only be observed here.
// Swapping the process streams forbids parallel test execution.
func captureOSOutput(t *testing.T, fn func()) (string, string) {
	t.Helper()

	oldStdout, oldStderr := os.Stdout, os.Stderr

	stdoutR, stdoutW, pipeErr := os.Pipe()
	if pipeErr != nil {
		t.Fatalf("stdout pipe: %v", pipeErr)
	}

	stderrR, stderrW, pipeErr := os.Pipe()
	if pipeErr != nil {
		t.Fatalf("stderr pipe: %v", pipeErr)
	}

	os.Stdout, os.Stderr = stdoutW, stderrW
	defer func() { os.Stdout, os.Stderr = oldStdout, oldStderr }()

	fn()

	stdoutW.Close()
	stderrW.Close()

	var stdoutBuf, stderrBuf bytes.Buffer

	_, _ = stdoutBuf.ReadFrom(stdoutR)
	_, _ = stderrBuf.ReadFrom(stderrR)

	return stdoutBuf.String(), stderrBuf.String()
}

// newFangOnFailingCLI builds a one-command CLI with fang left enabled (the
// default) whose handler fails with a hinted sentinel error.
func newFangOnFailingCLI(t *testing.T, opts ...CLIOption) *CLI[testConfig] {
	t.Helper()

	cli, err := NewCLI("app", "App", testConfig{}, opts...)
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

func TestExecuteFangOn_RendersHintThroughFangHandler(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	cli := newFangOnFailingCLI(t)

	var execErr error

	stdout, stderr := captureOSOutput(t, func() {
		execErr = cli.ExecuteWithArgs(context.Background(), []string{"fail"})
	})

	if execErr == nil {
		t.Fatal("expected the command error")
	}

	if strings.Contains(strings.ToLower(stdout), "required flag not set") {
		t.Errorf("errors must go to stderr, found in stdout: %s", stdout)
	}

	// fang's DefaultErrorHandler renders the message with its own casing, so
	// assertions on the message text compare case-insensitively.
	normalized := strings.ToLower(stderr)

	assertStringContains(t, normalized, "required flag not set")
	assertStringContains(t, stderr, "Hint: "+ErrorHint(ErrRequiredFlag))

	if got := strings.Count(normalized, "required flag not set"); got != 1 {
		t.Errorf("error must be printed exactly once under fang, found %d times: %s", got, stderr)
	}
}

func TestExecuteFangOn_BatteriesVersionAndDoctorRun(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	cli, err := NewCLI(
		"app", "App", testConfig{},
		WithCLIVersion("1.2.3"),
		WithBatteries(),
	)
	testutil.AssertNoError(t, err)

	stdout, _ := captureOSOutput(t, func() {
		err = cli.ExecuteWithArgs(context.Background(), []string{"version"})
	})
	testutil.AssertNoError(t, err)
	assertStringContains(t, stdout, "app 1.2.3")

	stdout, _ = captureOSOutput(t, func() {
		err = cli.ExecuteWithArgs(context.Background(), []string{"doctor"})
	})
	testutil.AssertNoError(t, err)
	assertStringContains(t, stdout, "0 passed, 0 failed")
}

func TestExecuteFangOn_WithoutErrorHintsSuppressesHintLine(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	cli := newFangOnFailingCLI(t, WithoutErrorHints())

	var execErr error

	_, stderr := captureOSOutput(t, func() {
		execErr = cli.ExecuteWithArgs(context.Background(), []string{"fail"})
	})

	if execErr == nil {
		t.Fatal("expected the command error")
	}

	assertStringContains(t, strings.ToLower(stderr), "required flag not set")

	if strings.Contains(stderr, "Hint: ") {
		t.Errorf("expected no hint line under WithoutErrorHints, got: %s", stderr)
	}
}

func TestExecuteFangOn_BatteriesDraconianBuiltInCommandsPass(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	cli, err := NewCLI(
		"app", "App", testConfig{},
		WithCLIVersion("1.2.3"),
		WithBatteries(),
		WithDraconianValidation(),
	)
	testutil.AssertNoError(t, err)

	for _, args := range [][]string{{"version"}, {"doctor"}, {"describe"}} {
		var execErr error

		captureOSOutput(t, func() {
			execErr = cli.ExecuteWithArgs(context.Background(), args)
		})

		if execErr != nil {
			t.Errorf("built-in command %v must pass Draconian validation, got: %v", args, execErr)
		}
	}
}

func TestExecuteFangOn_MachineFormatCarriesHintInEnvelopeOnly(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	cli := newFangOnFailingCLI(t, WithOutputFormat(output.FormatJSON))

	var execErr error

	_, stderr := captureOSOutput(t, func() {
		execErr = cli.ExecuteWithArgs(context.Background(), []string{"fail"})
	})

	if execErr == nil {
		t.Fatal("expected the command error")
	}

	var envelope jsonErrorEnvelope

	if decodeErr := json.Unmarshal([]byte(stderr), &envelope); decodeErr != nil {
		t.Fatalf("machine formats must emit a JSON envelope on stderr, got: %s", stderr)
	}

	if envelope.Error.Type != "flag" {
		t.Errorf("Error.Type = %q, want %q", envelope.Error.Type, "flag")
	}

	if envelope.Error.Hint != ErrorHint(ErrRequiredFlag) {
		t.Errorf("Error.Hint = %q, want %q", envelope.Error.Hint, ErrorHint(ErrRequiredFlag))
	}

	if strings.Contains(stderr, "Hint: ") {
		t.Errorf("machine formats must not print a standalone hint line, got: %s", stderr)
	}
}
