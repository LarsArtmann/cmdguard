package v4

import (
	"errors"
	"fmt"
	"io"

	"charm.land/fang/v2"
)

// sentinelFact pairs a sentinel with its stable machine error code and one
// actionable sentence telling the caller how to fix the failure. Order
// matters: more specific sentinels must come before the generic wrappers they
// can be nested inside (notably ErrFlagParseFailed, which wraps most flag
// errors).
type sentinelFact struct {
	match error
	code  string
	hint  string
}

// sentinelFacts is the ordered fact table. Adding a sentinel without a code
// and hint here fails the classification coverage test.
//
// The code is part of the machine contract (the JSON error envelope's
// error_code field): it is stable across releases even if the Go identifier
// is ever renamed, so scripts and agents can switch on it.
func sentinelFacts() []sentinelFact {
	return []sentinelFact{
		{ErrRequiredFlag, "required_flag", "pass the flag on the command line, or set its env var or config-file value"},
		{ErrInvalidEnum, "invalid_enum", "use one of the allowed values listed in this flag's --help text"},
		{ErrValueOutOfRange, "value_out_of_range", "choose a value that fits the flag's declared type"},
		{ErrValueTooShort, "value_too_short", "provide a longer value that satisfies the flag's validate rule"},
		{ErrValueTooLong, "value_too_long", "provide a shorter value that satisfies the flag's validate rule"},
		{ErrValueTooSmall, "value_too_small", "raise the value so it satisfies the flag's validate rule"},
		{ErrValueTooLarge, "value_too_large", "lower the value so it satisfies the flag's validate rule"},
		{ErrValuePatternMismatch, "value_pattern_mismatch", "match the pattern required by the flag's validate rule"},
		{ErrValueEmpty, "value_empty", "provide a non-empty value for this flag"},
		{ErrUnknownValidator, "unknown_validator", "use a validator registered with RegisterValidator, or fix the validate tag spelling"},
		{ErrInvalidValidatorParam, "invalid_validator_param", "fix the parameter syntax in the flag's validate tag"},
		{
			ErrInvalidFlagDefault,
			"invalid_flag_default",
			"fix the flag's default so it satisfies the same tag's values and validate rules",
		},
		{ErrUnsupportedFormat, "unsupported_format", "run with --help to see the supported output formats"},
		{ErrFormatRequiresTypedData, "format_requires_typed_data", "emit typed rows or values, or choose a text format instead"},
		{ErrUnsupportedAuditLogFormat, "unsupported_audit_log_format", "choose one of the formats accepted by ParseAuditLogFormat"},
		{ErrConfigValidation, "config_validation", "fix the reported config fields at their source: flag, env var, or config file"},
		{
			ErrConfigUnknownKey,
			"config_unknown_key",
			"fix or remove the key; its spelling must match a struct field, a json tag, or a flag name",
		},
		{
			ErrConfigFileParse,
			"config_file_parse",
			"check the config file's syntax and that its extension matches its JSON/YAML/TOML content",
		},
		{ErrConfigFileRead, "config_file_read", "check that the config file exists and is readable by this user"},
		{ErrConfigFileNotFound, "config_file_not_found", "pass --config with an existing path, or rely on defaults by omitting the file"},
		{ErrConfigFileLoad, "config_file_load", "check the config file path, permissions, and syntax"},
		{ErrConfigNil, "config_nil", "pass the config value instead of nil"},
		{ErrConfigNotPointer, "config_not_pointer", "pass a pointer to a struct as the config type"},
		{ErrMissingVersion, "missing_version", "pass WithCLIVersion to NewCLI before registering VersionCommand"},
		{ErrMissingShort, "missing_short", "add WithShort(...) so this command is discoverable in --help"},
		{ErrMissingLong, "missing_long", "add WithLong(...) describing what this parent command groups"},
		{ErrMissingExample, "missing_example", "add WithExample(...) or relax validation with WithStrictValidation"},
		{ErrMissingHandler, "missing_handler", "give the command a RunE handler or at least one subcommand"},
		{ErrMissingName, "missing_name", "pass a non-empty command name"},
		{ErrDuplicateCommand, "duplicate_command", "rename one command; the first word of Use must be unique in the CLI"},
		{ErrInvalidCommand, "invalid_command", "check the command definition for empty names or missing handlers"},
		{ErrCommandTimeout, "command_timeout", "raise the timeout or make the operation faster; the context deadline expired"},
		{ErrCommandPanic, "command_panic", "this is a bug in the CLI: inspect the handler that panicked"},
		{ErrDoctorFailed, "doctor_failed", "run the failing checks individually to see which dependency is unhealthy"},
		{ErrInvalidExitCode, "invalid_exit_code", "use an exit code between 0 and 255"},
		{ErrServiceNotFound, "service_not_found", "register the service with Provide or ProvideValue before invoking it"},
		{ErrServiceConstruction, "service_construction", "fix the provider's error; the service could not be constructed"},
		{ErrServiceRegistration, "service_registration", "check the service type and provider signature passed to Provide"},
		{ErrInvalidScope, "invalid_scope", "create the scope with NewScope before registering or invoking services"},
		{ErrFlagNotFound, "flag_not_found", "check the flag name against this command's --help output"},
		{ErrInvalidFlagType, "invalid_flag_type", "use a struct (or pointer to struct) for the command's flags type"},
		{ErrFlagInstance, "flag_instance", "the flags value could not be copied; use a struct or pointer to struct"},
		{ErrFieldNotFound, "field_not_found", "check the struct field name used by the flag tag"},
		{ErrFieldNotSettable, "field_not_settable", "export the struct field so it can receive a parsed flag value"},
		{ErrTypeConversion, "type_conversion", "use a supported flag type or register a TypeHandler for this type"},
		{ErrUnsupportedConversion, "unsupported_conversion", "use a supported flag type or register a TypeHandler for this type"},
		{ErrNilValue, "nil_value", "provide a value; nil is not accepted here"},
		{ErrInvalidDuration, "invalid_duration", "use a duration like 30s, 5m, or 2h"},
		{ErrInvalidURL, "invalid_url", "use a full URL including scheme, for example https://example.com"},
		{ErrInvalidEmail, "invalid_email", "use an address like name@example.com"},
		{ErrInvalidPort, "invalid_port", "use a port between 1 and 65535"},
		{ErrInvalidFilePath, "invalid_file_path", "check the path syntax and that the file exists"},
		{ErrInvalidHostPort, "invalid_host_port", "use host:port, for example localhost:8080"},
		{ErrLogLevel, "log_level", "use one of: debug, info, warn, error"},
		{ErrLogFormat, "log_format", "use one of: text, json"},
		{ErrNegativeArgCount, "negative_arg_count", "fix the command registration: argument counts cannot be negative"},
		{ErrInvalidArgRange, "invalid_arg_range", "fix the command registration: minimum args must not exceed maximum args"},
		{ErrEditorTempFile, "editor_temp_file", "check the temp directory is writable"},
		{ErrEditorWrite, "editor_write", "check the temp file is writable"},
		{ErrEditorRun, "editor_run", "check that $EDITOR points to an executable program"},
		{ErrEditorRead, "editor_read", "check the edited file is readable"},
		{ErrFlagParseFailed, "flag_parse_failed", "run this command with --help to see the accepted flags"},
	}
}

// ErrorHint returns one actionable sentence describing how to fix err, or an
// empty string when no hint is known.
//
// The hint is deliberately separate from the error message: messages stay
// stable (and machine-parseable), while hints can improve without breaking
// callers. The JSON error envelope emitted for machine-readable output formats
// includes the same hint under "hint".
func ErrorHint(err error) string {
	fact, ok := classifySentinel(err)
	if !ok {
		return ""
	}

	return fact.hint
}

// ErrorCode returns the stable machine error code for err, or an empty string
// when the error is not a cmdguard sentinel (or a wrapper around one).
//
// Codes are lowercase snake_case identifiers that stay stable across releases
// even if the Go identifier is renamed; scripts and agents switch on them via
// the JSON error envelope's error_code field instead of parsing messages.
func ErrorCode(err error) string {
	fact, ok := classifySentinel(err)
	if !ok {
		return ""
	}

	return fact.code
}

// classifySentinel walks the facts table most-specific-first and returns the
// first entry whose sentinel is in err's chain.
func classifySentinel(err error) (sentinelFact, bool) {
	if err == nil {
		return sentinelFact{}, false
	}

	for _, fact := range sentinelFacts() {
		if errors.Is(err, fact.match) {
			return fact, true
		}
	}

	return sentinelFact{}, false
}

// writeErrorWithHint is the default fang error handler. It renders the error
// exactly as fang normally would, then appends one actionable hint line when
// cmdguard knows how to fix the failure.
func writeErrorWithHint(w io.Writer, styles fang.Styles, err error) {
	fang.DefaultErrorHandler(w, styles, err)

	hint := ErrorHint(err)
	if hint == "" {
		return
	}

	_, _ = fmt.Fprintf(w, "Hint: %s\n", hint)
}
