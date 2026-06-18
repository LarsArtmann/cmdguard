package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	output "github.com/larsartmann/go-output"

	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

type jsonErrConfig struct {
	Verbose bool `default:"false" flag:"verbose" help:"Enable verbose output" short:"v"`
}

func TestClassifyError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		err      error
		wantType string
	}{
		{"flag parse", fmt.Errorf("%w: bad flag", ErrFlagParseFailed), "flag"},
		{"required flag", fmt.Errorf("%w: --name", ErrRequiredFlag), "flag"},
		{"config validation", fmt.Errorf("%w: missing field", ErrConfigValidation), "config"},
		{"config file load", fmt.Errorf("%w: file not found", ErrConfigFileLoad), "config"},
		{"invalid command", fmt.Errorf("%w: bad cmd", ErrInvalidCommand), "command"},
		{"service not found", fmt.Errorf("%w: missing svc", ErrServiceNotFound), "service"},
		{"value out of range", fmt.Errorf("%w: 999 into int8", ErrValueOutOfRange), "flag"},
		{"unsupported format", fmt.Errorf("%w: xml", ErrUnsupportedFormat), "output"},
		{"unknown error", errors.New("something unexpected"), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := classifyError(tt.err)
			if got != tt.wantType {
				t.Errorf("classifyError(%v) = %q, want %q", tt.err, got, tt.wantType)
			}
		})
	}
}

func TestWriteJSONError(t *testing.T) {
	t.Parallel()

	t.Run("produces valid JSON envelope", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		err := writeJSONError(&buf, fmt.Errorf("%w: --name is required", ErrRequiredFlag))
		testutil.AssertNoError(t, err)

		var envelope jsonErrorEnvelope
		decodeErr := json.Unmarshal(buf.Bytes(), &envelope)
		testutil.AssertNoError(t, decodeErr)

		if envelope.Error.Type != "flag" {
			t.Errorf("Error.Type = %q, want %q", envelope.Error.Type, "flag")
		}

		if envelope.Error.Code != 1 {
			t.Errorf("Error.Code = %d, want %d", envelope.Error.Code, 1)
		}

		if !strings.Contains(envelope.Error.Message, "--name is required") {
			t.Errorf("Error.Message = %q, want to contain '--name is required'", envelope.Error.Message)
		}
	})

	t.Run("includes exit code from ExitCoder", func(t *testing.T) {
		t.Parallel()

		exitErr, exitErrErr := NewExitError(42, fmt.Errorf("%w: bad", ErrInvalidCommand))
		testutil.AssertNoError(t, exitErrErr)

		var buf bytes.Buffer
		writeErr := writeJSONError(&buf, exitErr)
		testutil.AssertNoError(t, writeErr)

		var envelope jsonErrorEnvelope
		decodeErr := json.Unmarshal(buf.Bytes(), &envelope)
		testutil.AssertNoError(t, decodeErr)

		if envelope.Error.Code != 42 {
			t.Errorf("Error.Code = %d, want %d", envelope.Error.Code, 42)
		}

		if envelope.Error.Type != "command" {
			t.Errorf("Error.Type = %q, want %q", envelope.Error.Type, "command")
		}
	})
}

func TestExtractExitCode(t *testing.T) {
	t.Parallel()

	t.Run("returns 1 for plain errors", func(t *testing.T) {
		t.Parallel()

		code := extractExitCode(errors.New("plain"))
		if code != 1 {
			t.Errorf("extractExitCode(plain error) = %d, want 1", code)
		}
	})

	t.Run("returns code from ExitCoder", func(t *testing.T) {
		t.Parallel()

		exitErr, err := NewExitError(99, errors.New("bad"))
		testutil.AssertNoError(t, err)

		code := extractExitCode(exitErr)
		if code != 99 {
			t.Errorf("extractExitCode(ExitError) = %d, want 99", code)
		}
	})

	t.Run("returns code from wrapped ExitCoder", func(t *testing.T) {
		t.Parallel()

		exitErr, err := NewExitError(7, errors.New("wrapped"))
		testutil.AssertNoError(t, err)

		wrapped := fmt.Errorf("outer: %w", exitErr)
		code := extractExitCode(wrapped)
		if code != 7 {
			t.Errorf("extractExitCode(wrapped ExitError) = %d, want 7", code)
		}
	})
}

//nolint:paralleltest // os.Stderr manipulation
func TestJSONErrorIntegration(t *testing.T) {
	tests := []struct {
		name           string
		cliOpts        []CLIOption[jsonErrConfig]
		handlerErr     error
		assertEnvelope func(t *testing.T, stderrOutput string, execErr error)
	}{
		{
			name:       "no JSON output when output format not set",
			cliOpts:    nil,
			handlerErr: errors.New("plain error"),
			assertEnvelope: func(t *testing.T, stderrOutput string, execErr error) {
				t.Helper()
				if execErr == nil {
					t.Fatal("expected error from failing command")
				}
				var envelope jsonErrorEnvelope
				if decodeErr := json.Unmarshal([]byte(stderrOutput), &envelope); decodeErr == nil {
					t.Error("expected no JSON envelope when output format is not set, but got valid JSON")
				}
			},
		},
		{
			name:       "silences cobra errors when JSON output is active",
			cliOpts:    []CLIOption[jsonErrConfig]{WithOutputFormat[jsonErrConfig](output.FormatJSON)},
			handlerErr: errors.New("plain error"),
			assertEnvelope: func(t *testing.T, stderrOutput string, execErr error) {
				t.Helper()
				if execErr == nil {
					t.Fatal("expected error from failing command")
				}
				var envelope jsonErrorEnvelope
				decodeErr := json.Unmarshal([]byte(stderrOutput), &envelope)
				testutil.AssertNoError(t, decodeErr)
				testutil.AssertBoolTrue(t, envelope.Error.Code == 1, "exit code should be 1")
			},
		},
		{
			name:       "YAML output also triggers structured errors",
			cliOpts:    []CLIOption[jsonErrConfig]{WithOutputFormat[jsonErrConfig](output.FormatYAML)},
			handlerErr: fmt.Errorf("%w: missing", ErrConfigValidation),
			assertEnvelope: func(t *testing.T, stderrOutput string, execErr error) {
				t.Helper()
				if execErr == nil {
					t.Fatal("expected error from failing command")
				}
				var envelope jsonErrorEnvelope
				decodeErr := json.Unmarshal([]byte(stderrOutput), &envelope)
				testutil.AssertNoError(t, decodeErr)
				if envelope.Error.Type != "config" {
					t.Errorf("Error.Type = %q, want %q", envelope.Error.Type, "config")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli, cliErr := NewCLI[jsonErrConfig]("test", "test", jsonErrConfig{}, tt.cliOpts...)
			testutil.AssertNoError(t, cliErr)

			cmd, cmdErr := NewCommand[jsonErrConfig, *NoFlags](
				"fail",
				func(_ context.Context, _ *jsonErrConfig, _ *NoFlags) error {
					return tt.handlerErr
				},
				WithShort[jsonErrConfig, *NoFlags]("fails"),
			)
			testutil.AssertNoError(t, cmdErr)
			testutil.AssertNoError(t, AddCommand(cli, cmd))

			oldStderr := os.Stderr
			r, w, pipeErr := os.Pipe()
			if pipeErr != nil {
				t.Fatalf("pipe: %v", pipeErr)
			}

			os.Stderr = w

			execErr := cli.ExecuteWithArgs(context.Background(), []string{"fail"})

			w.Close()
			os.Stderr = oldStderr

			var buf bytes.Buffer
			_, _ = buf.ReadFrom(r)

			tt.assertEnvelope(t, buf.String(), execErr)
		})
	}
}
