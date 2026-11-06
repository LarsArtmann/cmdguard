package v4

import (
	"errors"
	"fmt"
	"io"

	"charm.land/fang/v2"
)

// sentinelFact is the single canonical row describing a cmdguard sentinel:
// its Go identifier, the JSON error-envelope class, the stable machine error
// code, and one actionable hint. Everything that publishes sentinel facts
// (error classification, ErrorCode, ErrorHint, the generated error reference)
// reads this one table, so the facts cannot drift apart.
//
// Order matters: more specific sentinels must come before the generic
// wrappers they can be nested inside (notably ErrFlagParseFailed, which wraps
// most flag errors) — classification walks the table first match wins.
type sentinelFact struct {
	match error
	name  string
	class string
	code  string
	hint  string
}

// sentinelFacts is the ordered fact table. Adding a sentinel without a name,
// class, code, and hint here fails the classification coverage test.
//
// The code is part of the machine contract (the JSON error envelope's
// error_code field): it is stable across releases even if the Go identifier
// is ever renamed, so scripts and agents can switch on it. The name column is
// documentation-facing (godoc anchor), the class column feeds the envelope's
// type field, and both are derived here so adding a sentinel is a single
// table row.
func sentinelFacts() []sentinelFact {
	return []sentinelFact{
		{ErrRequiredFlag, "ErrRequiredFlag", "flag", "required_flag", "pass the flag on the command line, or set its env var or config-file value"},
		{ErrInvalidEnum, "ErrInvalidEnum", "flag", "invalid_enum", "use one of the allowed values listed in this flag's --help text"},
		{ErrValueOutOfRange, "ErrValueOutOfRange", "flag", "value_out_of_range", "choose a value that fits the flag's declared type"},
		{ErrValueTooShort, "ErrValueTooShort", "flag", "value_too_short", "provide a longer value that satisfies the flag's validate rule"},
		{ErrValueTooLong, "ErrValueTooLong", "flag", "value_too_long", "provide a shorter value that satisfies the flag's validate rule"},
		{ErrValueTooSmall, "ErrValueTooSmall", "flag", "value_too_small", "raise the value so it satisfies the flag's validate rule"},
		{ErrValueTooLarge, "ErrValueTooLarge", "flag", "value_too_large", "lower the value so it satisfies the flag's validate rule"},
		{ErrValuePatternMismatch, "ErrValuePatternMismatch", "flag", "value_pattern_mismatch", "match the pattern required by the flag's validate rule"},
		{ErrValueEmpty, "ErrValueEmpty", "flag", "value_empty", "provide a non-empty value for this flag"},
		{ErrInvalidDuration, "ErrInvalidDuration", "flag", "invalid_duration", "use a duration like 30s, 5m, or 2h"},
		{ErrInvalidURL, "ErrInvalidURL", "flag", "invalid_url", "use a full URL including scheme, for example https://example.com"},
		{ErrInvalidEmail, "ErrInvalidEmail", "flag", "invalid_email", "use an address like name@example.com"},
		{ErrInvalidPort, "ErrInvalidPort", "flag", "invalid_port", "use a port between 1 and 65535"},
		{ErrInvalidFilePath, "ErrInvalidFilePath", "flag", "invalid_file_path", "check the path syntax and that the file exists"},
		{ErrInvalidHostPort, "ErrInvalidHostPort", "flag", "invalid_host_port", "use host:port, for example localhost:8080"},
		{ErrLogLevel, "ErrLogLevel", "flag", "log_level", "use one of: debug, info, warn, error"},
		{ErrLogFormat, "ErrLogFormat", "flag", "log_format", "use one of: text, json"},
		{ErrUnknownValidator, "ErrUnknownValidator", "flag", "unknown_validator", "use a validator registered with RegisterValidator, or fix the validate tag spelling"},
		{ErrInvalidValidatorParam, "ErrInvalidValidatorParam", "flag", "invalid_validator_param", "fix the parameter syntax in the flag's validate tag"},
		{ErrInvalidFlagDefault, "ErrInvalidFlagDefault", "flag", "invalid_flag_default", "fix the flag's default so it satisfies the same tag's values and validate rules"},
		{ErrFlagNotFound, "ErrFlagNotFound", "flag", "flag_not_found", "check the flag name against this command's --help output"},
		{ErrInvalidFlagType, "ErrInvalidFlagType", "flag", "invalid_flag_type", "use a struct (or pointer to struct) for the command's flags type"},
		{ErrFlagInstance, "ErrFlagInstance", "flag", "flag_instance", "the flags value could not be copied; use a struct or pointer to struct"},
		{ErrTypeConversion, "ErrTypeConversion", "flag", "type_conversion", "use a supported flag type or register a TypeHandler for this type"},
		{ErrUnsupportedConversion, "ErrUnsupportedConversion", "flag", "unsupported_conversion", "use a supported flag type or register a TypeHandler for this type"},
		{ErrNilValue, "ErrNilValue", "flag", "nil_value", "provide a value; nil is not accepted here"},
		{ErrConfigValidation, "ErrConfigValidation", "config", "config_validation", "fix the reported config fields at their source: flag, env var, or config file"},
		{ErrConfigUnknownKey, "ErrConfigUnknownKey", "config", "config_unknown_key", "fix or remove the key; its spelling must match a struct field, a json tag, or a flag name"},
		{ErrConfigFileParse, "ErrConfigFileParse", "config", "config_file_parse", "check the config file's syntax and that its extension matches its JSON/YAML/TOML content"},
		{ErrConfigFileRead, "ErrConfigFileRead", "config", "config_file_read", "check that the config file exists and is readable by this user"},
		{ErrConfigFileNotFound, "ErrConfigFileNotFound", "config", "config_file_not_found", "pass --config with an existing path, or rely on defaults by omitting the file"},
		{ErrConfigFileLoad, "ErrConfigFileLoad", "config", "config_file_load", "check the config file path, permissions, and syntax"},
		{ErrConfigNil, "ErrConfigNil", "config", "config_nil", "pass the config value instead of nil"},
		{ErrConfigNotPointer, "ErrConfigNotPointer", "config", "config_not_pointer", "pass a pointer to a struct as the config type"},
		{ErrFieldNotFound, "ErrFieldNotFound", "config", "field_not_found", "check the struct field name used by the flag tag"},
		{ErrFieldNotSettable, "ErrFieldNotSettable", "config", "field_not_settable", "export the struct field so it can receive a parsed flag value"},
		{ErrMissingVersion, "ErrMissingVersion", "command", "missing_version", "pass WithCLIVersion to NewCLI before registering VersionCommand"},
		{ErrMissingShort, "ErrMissingShort", "command", "missing_short", "add WithShort(...) so this command is discoverable in --help"},
		{ErrMissingLong, "ErrMissingLong", "command", "missing_long", "add WithLong(...) describing what this parent command groups"},
		{ErrMissingExample, "ErrMissingExample", "command", "missing_example", "add WithExample(...) or relax validation with WithStrictValidation"},
		{ErrMissingHandler, "ErrMissingHandler", "command", "missing_handler", "give the command a RunE handler or at least one subcommand"},
		{ErrMissingName, "ErrMissingName", "command", "missing_name", "pass a non-empty command name"},
		{ErrDuplicateCommand, "ErrDuplicateCommand", "command", "duplicate_command", "rename one command; the first word of Use must be unique in the CLI"},
		{ErrInvalidCommand, "ErrInvalidCommand", "command", "invalid_command", "check the command definition for empty names or missing handlers"},
		{ErrCommandTimeout, "ErrCommandTimeout", "command", "command_timeout", "raise the timeout or make the operation faster; the context deadline expired"},
		{ErrCommandPanic, "ErrCommandPanic", "command", "command_panic", "this is a bug in the CLI: inspect the handler that panicked"},
		{ErrNegativeArgCount, "ErrNegativeArgCount", "command", "negative_arg_count", "fix the command registration: argument counts cannot be negative"},
		{ErrInvalidArgRange, "ErrInvalidArgRange", "command", "invalid_arg_range", "fix the command registration: minimum args must not exceed maximum args"},
		{ErrInvalidScope, "ErrInvalidScope", "service", "invalid_scope", "create the scope with NewScope before registering or invoking services"},
		{ErrServiceNotFound, "ErrServiceNotFound", "service", "service_not_found", "register the service with Provide or ProvideValue before invoking it"},
		{ErrServiceConstruction, "ErrServiceConstruction", "service", "service_construction", "fix the provider's error; the service could not be constructed"},
		{ErrServiceRegistration, "ErrServiceRegistration", "service", "service_registration", "check the service type and provider signature passed to Provide"},
		{ErrUnsupportedFormat, "ErrUnsupportedFormat", "output", "unsupported_format", "run with --help to see the supported output formats"},
		{ErrFormatRequiresTypedData, "ErrFormatRequiresTypedData", "output", "format_requires_typed_data", "emit typed rows or values, or choose a text format instead"},
		{ErrUnsupportedAuditLogFormat, "ErrUnsupportedAuditLogFormat", "audit", "unsupported_audit_log_format", "choose one of the formats accepted by ParseAuditLogFormat"},
		{ErrInvalidExitCode, "ErrInvalidExitCode", "exit_code", "invalid_exit_code", "use an exit code between 0 and 255"},
		{ErrDoctorFailed, "ErrDoctorFailed", "doctor", "doctor_failed", "run the failing checks individually to see which dependency is unhealthy"},
		{ErrEditorTempFile, "ErrEditorTempFile", "editor", "editor_temp_file", "check the temp directory is writable"},
		{ErrEditorWrite, "ErrEditorWrite", "editor", "editor_write", "check the temp file is writable"},
		{ErrEditorRun, "ErrEditorRun", "editor", "editor_run", "check that $EDITOR points to an executable program"},
		{ErrEditorRead, "ErrEditorRead", "editor", "editor_read", "check the edited file is readable"},
		{ErrFlagParseFailed, "ErrFlagParseFailed", "flag", "flag_parse_failed", "run this command with --help to see the accepted flags"},
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
