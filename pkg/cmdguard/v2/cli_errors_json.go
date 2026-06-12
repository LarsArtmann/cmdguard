package v2

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	output "github.com/larsartmann/go-output"
)

type jsonErrorEnvelope struct {
	Error jsonErrorDetail `json:"error"`
}

type jsonErrorDetail struct {
	Type    string `json:"type"`
	Message string `json:"message"`
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
		{ErrAuditLogNotEnabled, "audit"},
		{ErrInvalidOutputFormat, "audit"},
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
		return exitCoder.ExitCode()
	}

	return 1
}

func writeJSONError(w io.Writer, err error) error {
	envelope := jsonErrorEnvelope{
		Error: jsonErrorDetail{
			Type:    classifyError(err),
			Message: err.Error(),
			Code:    extractExitCode(err),
		},
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	if encodeErr := enc.Encode(envelope); encodeErr != nil {
		return fmt.Errorf("encoding JSON error envelope: %w", encodeErr)
	}

	return nil
}

func (cli *CLI[T]) writeFormattedError(err error) bool {
	if cli.outputFormat == "" {
		return false
	}

	switch cli.outputFormat {
	case output.FormatJSON, output.FormatJSONL, output.FormatYAML, output.FormatTOML:
		writeErr := writeJSONError(os.Stderr, err)
		if writeErr != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
		}

		return true
	default:
		return false
	}
}
