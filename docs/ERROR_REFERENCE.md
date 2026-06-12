# Error Reference — cmdguard v2

**Auto-generated from source:** `pkg/cmdguard/v2/errors_*.go`, `types_*.go`
**Last updated:** 2026-06-10

---

## Usage

All sentinel errors support `errors.Is()` for identification:

```go
if errors.Is(err, v2.ErrRequiredFlag) {
    // handle missing required flag
}
```

## Typed Errors

| Type            | Constructor                                         | Purpose                               |
| --------------- | --------------------------------------------------- | ------------------------------------- |
| `CommandError`  | `NewCommandError(name, err)`                        | Wraps error with command name context |
| `FlagError`     | `NewFlagError(name, err)`                           | Wraps error with flag name context    |
| `FlagError`     | `NewFlagErrorWithSuggestion(name, err, suggestion)` | Includes typo suggestion              |
| `ConfigError`   | `NewConfigError(field, err)`                        | Wraps error with config field context |
| `EnumError`     | `NewEnumError(value, allowed)`                      | Invalid enum value with allowed list  |
| `DurationError` | `NewDurationError(value, err)`                      | Invalid duration format               |
| `ServiceError`  | `NewServiceError(type, err)`                        | DI service construction failure       |
| `ExitError`     | `NewExitError(code, err)`                           | Error with custom exit code (0–255)   |

All typed errors implement `Unwrap()` for `errors.Is()`/`errors.As()` chaining.

## Sentinel Errors by Category

### Command Errors

| Sentinel              | Message                                            | Source File       |
| --------------------- | -------------------------------------------------- | ----------------- |
| `ErrInvalidCommand`   | `"invalid command"`                                | errors_command.go |
| `ErrMissingHandler`   | `"command has no handler"`                         | errors_command.go |
| `ErrMissingName`      | `"command has no name"`                            | errors_command.go |
| `ErrMissingShort`     | `"command has no short description"`               | errors_command.go |
| `ErrMissingLong`      | `"command has no long description"`                | errors_command.go |
| `ErrMissingExample`   | `"command has no example"`                         | errors_command.go |
| `ErrMissingVersion`   | `"version is required but not set"`                | errors.go         |
| `ErrDuplicateCommand` | `"duplicate command"`                              | errors_command.go |
| `ErrCommandPanic`     | `"command panicked"`                               | errors_command.go |
| `ErrNegativeArgCount` | `"argument count must not be negative"`            | errors_command.go |
| `ErrInvalidArgRange`  | `"minimum argument count must not exceed maximum"` | errors_command.go |

### Flag Errors

| Sentinel                   | Message                            | Source File     |
| -------------------------- | ---------------------------------- | --------------- |
| `ErrFlagParseFailed`       | `"failed to parse flags"`          | errors_flags.go |
| `ErrFlagNotFound`          | `"flag not found"`                 | errors_flags.go |
| `ErrFlagInstance`          | `"failed to create flag instance"` | errors_flags.go |
| `ErrRequiredFlag`          | `"required flag not set"`          | errors_flags.go |
| `ErrInvalidFlagType`       | `"invalid flag type"`              | errors_flags.go |
| `ErrUnknownValidator`      | `"unknown validator"`              | errors_flags.go |
| `ErrInvalidValidatorParam` | `"invalid validator parameter"`    | errors_flags.go |

### Config Errors

| Sentinel                | Message                                | Source File      |
| ----------------------- | -------------------------------------- | ---------------- |
| `ErrConfigNil`          | `"config must not be nil"`             | errors_config.go |
| `ErrConfigNotPointer`   | `"config must be a pointer to struct"` | errors_config.go |
| `ErrConfigValidation`   | `"config validation failed"`           | errors_config.go |
| `ErrConfigFileLoad`     | `"failed to load config file"`         | errors_config.go |
| `ErrConfigFileNotFound` | `"config file not found"`              | errors_config.go |
| `ErrConfigFileParse`    | `"failed to parse config file"`        | errors_config.go |
| `ErrConfigFileRead`     | `"failed to read config file"`         | errors_config.go |
| `ErrFieldNotFound`      | `"field not found"`                    | errors_config.go |
| `ErrFieldNotSettable`   | `"field is not settable"`              | errors_config.go |

### DI / Scope Errors

| Sentinel                 | Message                         | Source File  |
| ------------------------ | ------------------------------- | ------------ |
| `ErrInvalidScope`        | `"invalid scope"`               | errors_di.go |
| `ErrServiceNotFound`     | `"service not found"`           | errors_di.go |
| `ErrServiceConstruction` | `"service construction failed"` | errors_di.go |
| `ErrServiceRegistration` | `"service registration failed"` | errors_di.go |

### Type & Value Errors

| Sentinel                   | Message                           | Source File |
| -------------------------- | --------------------------------- | ----------- |
| `ErrInvalidEnum`           | `"invalid enum value"`            | errors.go   |
| `ErrInvalidDuration`       | `"invalid duration"`              | errors.go   |
| `ErrInvalidURL`            | `"invalid URL"`                   | errors.go   |
| `ErrInvalidEmail`          | `"invalid email address"`         | errors.go   |
| `ErrInvalidPort`           | `"invalid port"`                  | errors.go   |
| `ErrInvalidFilePath`       | `"invalid file path"`             | errors.go   |
| `ErrInvalidHostPort`       | `"invalid host:port"`             | errors.go   |
| `ErrLogLevel`              | `"invalid log level"`             | errors.go   |
| `ErrLogFormat`             | `"invalid log format"`            | errors.go   |
| `ErrTypeConversion`        | `"type conversion failed"`        | errors.go   |
| `ErrUnsupportedConversion` | `"unsupported string conversion"` | errors.go   |
| `ErrNilValue`              | `"value must not be nil"`         | errors.go   |

### Validation Errors

| Sentinel                  | Message                          | Source File |
| ------------------------- | -------------------------------- | ----------- |
| `ErrValueEmpty`           | `"value is empty"`               | errors.go   |
| `ErrValueTooShort`        | `"value too short"`              | errors.go   |
| `ErrValueTooLong`         | `"value too long"`               | errors.go   |
| `ErrValueTooSmall`        | `"value too small"`              | errors.go   |
| `ErrValueTooLarge`        | `"value too large"`              | errors.go   |
| `ErrValuePatternMismatch` | `"value does not match pattern"` | errors.go   |

### Output Errors

| Sentinel                     | Message                        | Source File |
| ---------------------------- | ------------------------------ | ----------- |
| `ErrUnsupportedFormat`       | `"unsupported output format"`  | errors.go   |
| `ErrFormatRequiresTypedData` | `"format requires typed data"` | errors.go   |

### Editor Errors

| Sentinel            | Message                                   | Source File |
| ------------------- | ----------------------------------------- | ----------- |
| `ErrEditorTempFile` | `"failed to create temp file for editor"` | errors.go   |
| `ErrEditorWrite`    | `"failed to write to temp file"`          | errors.go   |
| `ErrEditorRun`      | `"editor execution failed"`               | errors.go   |
| `ErrEditorRead`     | `"failed to read edited file"`            | errors.go   |

### Exit & Doctor Errors

| Sentinel             | Message                                 | Source File |
| -------------------- | --------------------------------------- | ----------- |
| `ErrInvalidExitCode` | `"exit code must be between 0 and 255"` | errors.go   |
| `ErrDoctorFailed`    | `"doctor checks failed"`                | errors.go   |

### Unused (reserved for future use)

| Sentinel | Message | Source File |
| -------- | ------- | ----------- |

---

**Total: 62 sentinel errors, 8 typed error constructors**
