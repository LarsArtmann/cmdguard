# SUPERB CLI Enforcement Gap Analysis

**Date:** 2026-05-16
**Status:** Post-v2.2 audit of enforcement capabilities

---

## Current Enforcement (What cmdguard already does)

| Enforcement | How |
|---|---|
| Name required on all commands | `Validate()` rejects empty `use` |
| Handler required on leaf commands | `Validate()` rejects no `runE` and no subcommands |
| Long description required on parent commands | `Validate()` rejects subcommands without `long` |
| No duplicate commands | `AddCommand` checks `registeredCmds` map |
| No duplicate subcommands | `Validate()` checks `seen` map |
| Short description required (opt-in) | `WithStrictValidation[T]()` + `ValidateStrict()` |
| Config validation before commands run | `WithConfigValidation[T](fn)` |
| Positional args bounds | `WithExactArgs`, `WithMinimumArgs`, `WithMaximumArgs`, `WithRangeArgs`, `WithNoArgs` |
| Custom exit codes | `ExitCoder` interface + `NewExitError(code, err)` |
| Version command helper | `VersionCommand[T](cli)` / `MustVersionCommand[T](cli)` |
| Typed errors with context | `CommandError`, `FlagError`, `ServiceError`, `ConfigError`, `EnumError` |
| Flag typo suggestions | Levenshtein-based `SuggestFlag` / `SuggestCommand` |
| Typed flag values | 9 custom types: Duration, Enum, LogLevel, LogFormat, URL, Email, Port, FilePath, HostPort |
| Flag validators | `validate:"email,min=5,max=100,regex=^..."` struct tags |
| Required flags | `required:"true"` struct tag |
| Env var binding | `env:"VAR_NAME"` + `WithEnvPrefix[T]()` |
| Signal handling | `WithSignalHandling[T]()` for SIGINT/SIGTERM |
| Panic recovery | `RecoveryMiddleware[T]()` with stack traces |
| Shell completion | `WithCompletion[T,F](fn)` / `WithValidArgs[T,F](args...)` |
| Man page generation | `cli.ManPage(section)`, `GenerateManPageCommand[T](cli)` |
| Styled help output | fang integration (enabled by default) |
| No panics in library API | All functions return errors |

---

## Critical Gaps (Ranked by Impact)

### 1. No CLI Test Harness

**Impact:** Users don't test → users ship broken CLIs.

Every example test manually constructs `CLI[T]`, calls `ExecuteWithArgs`, captures output. This boilerplate repeats 100+ times. No `NewTestCLI`, no `CaptureOutput`, no `ExecuteForTest`.

**What top frameworks provide:**
- urfave/cli v2: `cli.NewContext()` + `cli_test` helpers for isolated command testing
- kong: `kong.Parse()` with `kong.Exit(func(s int))` for capture, `kong.Vars` for injection
- cobra: `rootCmd.SetArgs()` + `bytes.Buffer` stdout/stderr capture (manual)

**Proposal:**

```go
// pkg/testutil/cli_harness.go (public, importable by users)

type TestResult struct {
    Stdout   string
    Stderr   string
    ExitCode int
    Err      error
}

// CaptureOutput runs a function capturing stdout and stderr.
func CaptureOutput(fn func() error) (stdout, stderr string, err error)

// NewTestCLI creates a CLI for testing with fang disabled and output captured.
func NewTestCLI[T any](name string, defaults T, opts ...v2.CLIOption[T]) (*v2.CLI[T], error)
```

---

### 2. No `WithRequireExamples` in Strict Mode

**Impact:** A superb CLI without examples is just "good." Users copy-paste from help.

`WithStrictValidation[T]()` should optionally enforce `WithExample[T,F](...)` on all leaf commands. Without examples, the `--help` output is technically valid but practically useless for discovery.

**Proposal:**

```go
// In strict mode, require examples on leaf commands
type CLI[T] struct {
    // ...
    requireExamples bool // set by WithRequireExamples
}

// Option: make it part of strict mode or a separate option
WithRequireExamples[T]() CLIOption[T]
```

Or simpler: add `requireExample bool` to the strict validation check, so `ValidateStrict()` also checks `c.example == ""` for leaf commands.

---

### 3. `help` Tag Not Enforced on Flags

**Impact:** `--help` output has gaps. Users see flags without descriptions.

The `help` struct tag is optional. A flag like:

```go
Verbose bool `flag:"verbose" default:"false"`
```

produces a help entry with **no description**. Strict mode should enforce it.

**Proposal:**

Add a strict validation pass in `FlagRegistry` that rejects empty `help` tags when strict mode is enabled. Wire through a `strict` flag from `CLI` to `FlagRegistry`.

---

### 4. NO_COLOR Not Explicitly Supported

**Impact:** CLIs break in CI, pipes, and CI logs.

Every production CLI respects `NO_COLOR` (https://no-color.org/) and/or provides `--no-color`. Fang/lipgloss handles it implicitly, but:
- No `--no-color` flag
- No documentation about it
- No `WithNoColor[T]()` option
- No `TERM=dumb` handling

**What top frameworks do:**
- urfave/cli v2: Built-in `NO_COLOR` detection
- kong: `term.IsTerminal` + `NO_COLOR` convention
- Most production CLIs: Explicit `--no-color` flag + env var check + `TERM=dumb`

**Proposal:**

```go
// Respect NO_COLOR automatically, add --no-color flag
WithNoColor[T]() CLIOption[T]

// Or: document that fang already respects NO_COLOR via lipgloss
// and add --no-color as a convenience flag
```

---

### 5. No Structured Error Output

**Impact:** CLIs are not automation-friendly. Errors are plain text only.

When `--output=json` is set, data output is JSON but errors are still plain text. The typed error hierarchy (`CommandError`, `FlagError`, `ServiceError`, `ExitError`) would serialize beautifully but isn't wired for it.

**What top frameworks do:**
- urfave/cli v2: `cli.ErrWriter` + JSON error formatting
- kong: `--format=json` for machine-readable output including errors
- Terraform/Pulumi: Structured error output for CI/automation

**Proposal:**

```go
// MarshalError serializes an error to structured JSON
func MarshalError(err error) ([]byte, error)

// ErrorSchema represents a structured error for machine consumption
type ErrorSchema struct {
    Type       string `json:"type"`       // "flag", "command", "config", "service"
    Code       string `json:"code"`       // sentinel error name
    Message    string `json:"message"`    // human-readable message
    Field      string `json:"field,omitempty"`    // flag/command/field name
    Suggestion string `json:"suggestion,omitempty"` // typo fix
    ExitCode   int    `json:"exit_code,omitempty"`
}
```

---

## Important Gaps

### 6. Examples Don't Cover 12+ Features

| Feature | Has Example? |
|---|---|
| `WithStrictValidation` | No |
| `WithConfigValidation` | No |
| `VersionCommand` / `MustVersionCommand` | No |
| Positional args (`WithExactArgs`, etc.) | No |
| `WithOutputFormat` | No |
| Middleware (`TimingMiddleware`, `RecoveryMiddleware`) | No |
| Man page generation | No |
| Shell completion (`WithCompletion`) | No |
| `EditInEditor` | No |
| `MustNewCommand` / `MustNewParentCommand` | No |
| `ExitCoder` / `ExitError` | No |
| `WithGroup` / `WithGroupID` | No |

**Proposal:** Create `examples/superb/` demonstrating all enforcement features together.

---

### 7. README is 25+ Public APIs Behind

**CLI options not in README:**
- `WithConfigValidation[T](fn)` — cross-field config validation
- `WithStrictValidation[T]()` — require short descriptions
- `WithOutputFormat[T]()` — auto `--output` flag
- `WithMiddleware[T](mw...)` — command middleware
- `WithGroup[T](id, title)` — command groups
- `WithFangOptions[T](opts...)` — custom fang styling

**Command options not in README:**
- `WithCompletion[T,F](fn)` — dynamic shell completion
- `WithValidArgs[T,F](args...)` — static valid args
- `WithExactArgs[T,F](n)` / `WithMinimumArgs` / `WithMaximumArgs` / `WithRangeArgs` / `WithNoArgs`
- `WithArgs[T,F](fn)` — custom args validator

**Entire API sections not in README:**
- Version command: `VersionCommand[T]`, `MustVersionCommand[T]`, `GenerateVersionCommand[T]`
- Man pages: `cli.ManPage(section)`, `cli.WriteManPage(w, section)`, `GenerateManPageCommand[T]`
- Error types: `ExitCoder`, `ExitError`, `NewExitError(code, err)` + exit codes in `ExecuteAndExit`
- Middleware: `TimingMiddleware`, `RecoveryMiddleware`, custom `Middleware[T]`
- Helpers: `EditInEditor`, `Ptr[T]`, `ValueOrDefault`, `MustParse[T]`, `MergeConfigs`
- Value types: 9 custom types (Duration, Enum, LogLevel, LogFormat, URL, Email, Port, FilePath, HostPort)
- Shell completion: `CompletionFunc`
- Flow context: `BranchingFlowContext`, `ArgsFromContext(ctx)`

---

### 8. No Flag Deprecation

Commands can be deprecated (`WithDeprecated`) but individual flags cannot.

**Proposal:**

```go
type FlagTag struct {
    // ...
    Deprecated string // `deprecated:"use --verbose instead"` tag
}
```

Add `deprecated` struct tag support that warns (not errors) when the flag is used.

---

## Prioritized Execution Plan

| # | What | Impact | Effort | Priority |
|---|------|--------|--------|----------|
| 1 | CLI Test Harness | Critical | Medium | P0 |
| 2 | `help` tag enforcement in strict mode | High | Small | P0 |
| 3 | `WithRequireExamples` in strict mode | High | Small | P0 |
| 4 | Explicit NO_COLOR support | High | Small | P1 |
| 5 | Structured error output (JSON errors) | High | Medium | P1 |
| 6 | Comprehensive example (`examples/superb/`) | High | Medium | P1 |
| 7 | README update for 25+ missing APIs | High | Medium | P1 |
| 8 | Flag deprecation tag | Medium | Small | P2 |

---

## Design Principle: Enforcement Spectrum

cmdguard should operate on a spectrum:

```
Lenient (default)          Strict                draconian
    |                        |                      |
    |  name required         |  + short required    |  + example required
    |  handler required      |  + help tags required|  + long required on ALL
    |  long for parents      |  + config validated  |  + flags must have short
    |                        |                      |  + no deprecated commands
```

Current state: lenient is solid. strict mode exists but is minimal. The gap is making strict mode genuinely enforce superb CLI quality.
