package v4

import "errors"

// errorHint pairs a sentinel with one actionable sentence telling the caller
// how to fix the failure. Order matters: more specific sentinels must come
// before the generic wrappers they can be nested inside (notably
// ErrFlagParseFailed, which wraps most flag errors).
type errorHint struct {
	match error
	hint  string
}

// errorHints is the ordered hint table. Adding a sentinel without adding a
// hint here fails the classification coverage test.
func errorHints() []errorHint {
	return []errorHint{
		{ErrRequiredFlag, "pass the flag on the command line, or set its env var or config-file value"},
		{ErrInvalidEnum, "use one of the allowed values listed in this flag's --help text"},
		{ErrValueOutOfRange, "choose a value that fits the flag's declared type"},
		{ErrValueTooShort, "provide a longer value that satisfies the flag's validate rule"},
		{ErrValueTooLong, "provide a shorter value that satisfies the flag's validate rule"},
		{ErrValueTooSmall, "raise the value so it satisfies the flag's validate rule"},
		{ErrValueTooLarge, "lower the value so it satisfies the flag's validate rule"},
		{ErrValuePatternMismatch, "match the pattern required by the flag's validate rule"},
		{ErrValueEmpty, "provide a non-empty value for this flag"},
		{ErrUnknownValidator, "use a validator registered with RegisterValidator, or fix the validate tag spelling"},
		{ErrInvalidValidatorParam, "fix the parameter syntax in the flag's validate tag"},
		{ErrUnsupportedFormat, "run with --help to see the supported output formats"},
		{ErrFormatRequiresTypedData, "emit typed rows or values, or choose a text format instead"},
		{ErrUnsupportedAuditLogFormat, "choose one of the formats accepted by ParseAuditLogFormat"},
		{ErrConfigValidation, "fix the reported config fields at their source: flag, env var, or config file"},
		{ErrConfigFileParse, "check the config file's syntax and that its extension matches its JSON/YAML/TOML content"},
		{ErrConfigFileRead, "check that the config file exists and is readable by this user"},
		{ErrConfigFileNotFound, "pass --config with an existing path, or rely on defaults by omitting the file"},
		{ErrConfigFileLoad, "check the config file path, permissions, and syntax"},
		{ErrConfigNil, "pass the config value instead of nil"},
		{ErrConfigNotPointer, "pass a pointer to a struct as the config type"},
		{ErrMissingVersion, "pass WithCLIVersion to NewCLI before registering VersionCommand"},
		{ErrMissingShort, "add WithShort(...) so this command is discoverable in --help"},
		{ErrMissingLong, "add WithLong(...) describing what this parent command groups"},
		{ErrMissingExample, "add WithExample(...) or relax validation with WithStrictValidation"},
		{ErrMissingHandler, "give the command a RunE handler or at least one subcommand"},
		{ErrMissingName, "pass a non-empty command name"},
		{ErrDuplicateCommand, "rename one command; the first word of Use must be unique in the CLI"},
		{ErrInvalidCommand, "check the command definition for empty names or missing handlers"},
		{ErrCommandTimeout, "raise the timeout or make the operation faster; the context deadline expired"},
		{ErrCommandPanic, "this is a bug in the CLI: inspect the handler that panicked"},
		{ErrDoctorFailed, "run the failing checks individually to see which dependency is unhealthy"},
		{ErrInvalidExitCode, "use an exit code between 0 and 255"},
		{ErrServiceNotFound, "register the service with Provide or ProvideValue before invoking it"},
		{ErrServiceConstruction, "fix the provider's error; the service could not be constructed"},
		{ErrServiceRegistration, "check the service type and provider signature passed to Provide"},
		{ErrInvalidScope, "create the scope with NewScope before registering or invoking services"},
		{ErrFlagNotFound, "check the flag name against this command's --help output"},
		{ErrInvalidFlagType, "use a struct (or pointer to struct) for the command's flags type"},
		{ErrFlagInstance, "the flags value could not be copied; use a struct or pointer to struct"},
		{ErrFieldNotFound, "check the struct field name used by the flag tag"},
		{ErrFieldNotSettable, "export the struct field so it can receive a parsed flag value"},
		{ErrTypeConversion, "use a supported flag type or register a TypeHandler for this type"},
		{ErrUnsupportedConversion, "use a supported flag type or register a TypeHandler for this type"},
		{ErrNilValue, "provide a value; nil is not accepted here"},
		{ErrInvalidDuration, "use a duration like 30s, 5m, or 2h"},
		{ErrInvalidURL, "use a full URL including scheme, for example https://example.com"},
		{ErrInvalidEmail, "use an address like name@example.com"},
		{ErrInvalidPort, "use a port between 1 and 65535"},
		{ErrInvalidFilePath, "check the path syntax and that the file exists"},
		{ErrInvalidHostPort, "use host:port, for example localhost:8080"},
		{ErrLogLevel, "use one of: debug, info, warn, error"},
		{ErrLogFormat, "use one of: text, json"},
		{ErrNegativeArgCount, "fix the command registration: argument counts cannot be negative"},
		{ErrInvalidArgRange, "fix the command registration: minimum args must not exceed maximum args"},
		{ErrEditorTempFile, "check the temp directory is writable"},
		{ErrEditorWrite, "check the temp file is writable"},
		{ErrEditorRun, "check that $EDITOR points to an executable program"},
		{ErrEditorRead, "check the edited file is readable"},
		{ErrFlagParseFailed, "run this command with --help to see the accepted flags"},
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
	if err == nil {
		return ""
	}

	for _, entry := range errorHints() {
		if errors.Is(err, entry.match) {
			return entry.hint
		}
	}

	return ""
}
