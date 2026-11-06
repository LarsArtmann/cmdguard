# Error Reference — cmdguard v4

**The sentinel table is generated from source** and drift-guarded by a test:
the facts (code, class, message, hint) come from the canonical sentinel fact
table in `pkg/cmdguard/v4/errors_facts.go`, and the doc cannot say anything
the code does not. Regenerate the bracketed section with:

```bash
go test ./pkg/cmdguard/v4 -run TestErrorReferenceDoc -update-golden
```

---

## Usage

All sentinel errors support `errors.Is()` for identification:

```go
if errors.Is(err, v4.ErrRequiredFlag) {
    // handle missing required flag
}
```

Prefer the stable machine code when the consumer is a script or an agent:
it survives identifier renames and reads from the JSON error envelope without
parsing messages.

```go
code := v4.ErrorCode(err) // "required_flag", "" when unknown
```

## Machine error codes

`ErrorCode(err)` returns one stable lowercase snake_case code per sentinel
(`""` when the error wraps no cmdguard sentinel). Machine-readable output
(`--output json|yaml|toml`) carries it in the JSON error envelope:

```json
{
  "error": {
    "type": "flag",
    "message": "required flag not set: name",
    "hint": "pass the flag on the command line, or set its env var or config-file value",
    "error_code": "required_flag",
    "code": 1
  }
}
```

The envelope fields mean: `type` is the error class, `message` the stable
error text, `hint` the actionable fix, `error_code` the stable machine code,
and `code` the process exit code.

## Actionable hints

`ErrorHint(err)` returns one sentence describing how to fix the failure, or
`""` when no hint is known:

```go
if hint := v4.ErrorHint(err); hint != "" {
    fmt.Fprintln(os.Stderr, hint)
}
```

Hints are separate from messages so error text stays stable for parsing and
tests. Every sentinel with a fact-table row has a hint; a test fails when a
new sentinel is added without one.

## Exit codes

`ExecuteAndExit` is the blessed entry point: exit `0` on success, `1` on a
failure with no explicit code, and the `ExitCoder` code when the error chain
carries one (0–255, validated by `NewExitError`). A zero code from a
third-party `ExitCoder` is treated as "not set" and falls back to `1`, so a
failed CLI can never exit `0`. `ExitCode(err)` exposes the same mapping for
post-execution teardown that needs the code before exiting.

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

## Sentinel Errors

