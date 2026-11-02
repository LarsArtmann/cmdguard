package v4

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"

	output "github.com/larsartmann/go-output"
)

type jsonErrorEnvelope struct {
	Error jsonErrorDetail `json:"error"`
}

type jsonErrorDetail struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Hint    string `json:"hint,omitempty"`
	Code    int    `json:"code,omitempty"`
}

type errorTypeEntry struct {
	match error
	name  string
}

func errorTypeNames() []errorTypeEntry {
	return []errorTypeEntry{
		{ErrFlagParseFailed, "flag"},
		{ErrRequiredFlag, "flag"},
		{ErrFlagNotFound, "flag"},
		{ErrInvalidFlagType, "flag"},
		{ErrInvalidEnum, "flag"},
		{ErrValueTooShort, "flag"},
		{ErrValueTooLong, "flag"},
		{ErrValueTooSmall, "flag"},
		{ErrValueTooLarge, "flag"},
		{ErrValueOutOfRange, "flag"},
		{ErrValuePatternMismatch, "flag"},
		{ErrValueEmpty, "flag"},
		{ErrConfigValidation, "config"},
		{ErrConfigNil, "config"},
		{ErrConfigNotPointer, "config"},
		{ErrConfigFileLoad, "config"},
		{ErrConfigFileRead, "config"},
		{ErrConfigFileParse, "config"},
		{ErrConfigFileNotFound, "config"},
		{ErrInvalidCommand, "command"},
		{ErrMissingHandler, "command"},
		{ErrMissingName, "command"},
		{ErrDuplicateCommand, "command"},
		{ErrCommandPanic, "command"},
		{ErrServiceNotFound, "service"},
		{ErrServiceConstruction, "service"},
		{ErrServiceRegistration, "service"},
		{ErrInvalidScope, "service"},
		{ErrUnsupportedFormat, "output"},
		{ErrFormatRequiresTypedData, "output"},
		{ErrUnsupportedAuditLogFormat, "audit"},
		{ErrInvalidExitCode, "exit_code"},
	}
}

func classifyError(err error) string {
	for _, entry := range errorTypeNames() {
		if errors.Is(err, entry.match) {
			return entry.name
		}
	}

	return "unknown"
}

func extractExitCode(err error) int {
	if exitCoder, ok := errors.AsType[ExitCoder](err); ok {
		// A 0 from an ExitCoder means "no override, fall back to the generic
		// failure code", not "success": err is non-nil here, so the process
		// failed and must not exit 0. Third-party error types (e.g.
		// go-error-family) use 0 as their "no explicit code set" sentinel while
		// still implementing ExitCoder; trusting it verbatim made failed CLIs
		// exit 0.
		if code := exitCoder.ExitCode(); code != 0 {
			return code
		}
	}

	return 1
}

// ExitCode extracts the process exit code from an error returned by Execute.
//
// If the error (or any error in its chain) implements ExitCoder with a
// NON-ZERO code, that code is used; a zero code from an ExitCoder is treated
// as "no override" (third-party error types use 0 as their "not set"
// sentinel) and falls back to the generic failure code 1. A nil error yields 0.
//
// cmdguard already prints the error exactly once (styled via fang when enabled,
// plain via cobra when disabled). The value returned by Execute is therefore for
// exit-code mapping and programmatic inspection only — do NOT print it again.
//
// Use ExecuteAndExit for the common case. Use ExitCode when you must run code
// after execution (flush logs, export an audit log, tear down resources) before
// the process exits:
//
//	err := cli.Execute(ctx)
//	// ...flush / export audit log / teardown...
//	os.Exit(v4.ExitCode(err))
func ExitCode(err error) int {
	if err == nil {
		return 0
	}

	return extractExitCode(err)
}

func writeJSONError(w io.Writer, err error) error {
	envelope := jsonErrorEnvelope{
		Error: jsonErrorDetail{
			Type:    classifyError(err),
			Message: err.Error(),
			Hint:    ErrorHint(err),
			Code:    extractExitCode(err),
		},
	}

	enc := jsontext.NewEncoder(w, jsontext.WithIndent("  "))

	encodeErr := json.MarshalEncode(enc, envelope)
	if encodeErr != nil {
		return fmt.Errorf("encoding JSON error envelope: %w", encodeErr)
	}

	return nil
}

func (cli *CLI[T]) writeFormattedError(err error) {
	if cli.spec.outputFormat == "" {
		return
	}

	// Formats in which errors are emitted as structured JSON envelopes
	// instead of plain text.
	machineReadable := []output.Format{
		output.FormatJSON,
		output.FormatJSONL,
		output.FormatYAML,
		output.FormatTOML,
	}

	if !slices.Contains(machineReadable, cli.spec.outputFormat) {
		return
	}

	writeErr := writeJSONError(os.Stderr, err)
	if writeErr != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}
}
