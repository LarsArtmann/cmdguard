# cmdguard — Comprehensive Status Report

**Date:** 2026-04-15 23:31 CEST
**Session:** Multi-session audit-and-improve sprint
**Branch:** master (clean working tree)
**Last Commit:** `0d6fc58` — feat: add command groups, middleware, Result type, and flag validation
**Go Version:** 1.26 | **Tests:** 14 packages, all PASS with `-race` | **Coverage:** v2 at 81.3%

---

## A. FULLY DONE ✅ (8 items)

### 1. Result[T] Type — `types_result.go` + `types_result_test.go`

- Rust-style `Ok(value)` / `Err(error)` discriminated union
- Methods: `Map`, `MapErr`, `MapOr`, `And`, `Or`, `IfOk`, `IfErr`, `ToOption`, `ResultFrom`, `ToPair`
- `Unwrap`, `UnwrapOr`, `UnwrapOrElse`, `UnwrapErr`, `Expect`, `ExpectErr`
- `MarshalJSON`, `String` for serialization
- **639 lines of tests** covering all methods, edge cases, panic paths

### 2. NoFlags Type Safety Fix — `command.go`

- Changed from `type NoFlags = struct{}` (alias) to `type NoFlags struct{}` (defined type)
- Prevents accidental comparison with bare `struct{}{}`, makes it a proper named type
- Technically a breaking change but acceptable since `NoFlags` is the documented API entry point

### 3. Silent Error Swallowing Fix — `config_parsing.go`

- `parseBoolDefault`, `parseIntDefault`, `parseUintDefault`, `parseFloat64Default` all changed from void/error-ignoring to `(value, error)` return tuples
- `registerFlag` in `flags.go` now propagates errors from these functions instead of silently using zero values
- Fixes a correctness bug where invalid default tags like `default:"not-a-number"` on int fields silently produced `0`

### 4. MergeConfigs Mutation Fix — `config.go` + `config_merge_test.go`

- `MergeConfigs[T]()` now deep-copies `configs[0]` before merging
- Previously mutated the first input in place — a silent data corruption bug
- Added regression tests: "does not mutate input configs" and "returned config is independent"

### 5. BranchingFlowContext Shared Map Fix — `flow_context.go` + `flow_context_value_test.go`

- `newChild()` now uses `maps.Clone(b.values)` instead of sharing the reference
- Children get their own map copy, preventing parent contamination via `SetValueLocal`
- Added `TestBranchingFlowContext_ChildValueIsolation` regression test

### 6. LogLevel/LogFormat Deduplication — `types_log.go`

- Extracted `logLevelAllowed` and `logFormatAllowed` package-level slices
- All 6 references (Parse/IsValid/Validate functions) now share the canonical slices
- Eliminates 4 duplicate slice literals that could drift apart

### 7. Middleware/Interceptor Chain — `middleware.go` + `middleware_test.go`

- `Middleware[T]` type: `func(ctx, cfg, info CommandInfo, next func() error) error`
- `buildChain` builds right-to-left so first middleware wraps outermost
- Built-in: `TimingMiddleware[T]` (logs duration), `RecoveryMiddleware[T]` (recovers panics to errors)
- `WithMiddleware[T]` CLI option registers middleware on CLI
- `wireHandlerWithMiddleware` wires chain into cobra's RunE/PreRunE/PostRunE
- `CommandInfo` carries command metadata (name, phase, hasRunE) to middleware
- **438 lines of tests**: chaining, error propagation, short-circuit, timing, recovery, subcommands, flags integration

### 8. Command Groups — `cli_command.go` + `cli_options.go` + `cli_groups_test.go`

- `Group string` field on `Command[T, F]` struct
- Wired to cobra's `GroupID` in `cliToCobraCommand`
- `WithGroup[T](id, title string)` CLI option registers `cobra.Group` on root command
- **235 lines of tests**: grouping, help output verification, no-group fallback, subcommands, multiple groups

---

## B. PARTIALLY DONE ⚠️ (1 item)

### 9. Flag Validators via `validate` Tag — `flags_validate.go` + `flags_validate_test.go`

- **Implementation (DONE):**
  - `FlagValidator func(value string) error` type
  - `validatorRegistry` with `sync.RWMutex` for goroutine safety
  - `RegisterValidator(name, validator)` public API for custom validators
  - 8 built-in validators: `email`, `url`, `minlen`, `maxlen`, `min`, `max`, `regex`, `nonempty`
  - `validate` tag parsing in `parseFieldFlag`
  - `Validate` field added to `FlagTag` struct
  - Wired into `validateTag` in `config.go` (for `ValidateConfig`)
  - Wired into `validateTagRules` in `flags.go` (for `ValidateFlags`)
  - `formatFieldValue` converts reflected values to strings for validation
  - **296 lines of implementation**

- **Tests (INCOMPLETE):**
  - The existing `flags_validate_test.go` only has tests for the pre-existing enum/required flag validation — **zero tests for the new `validate` tag feature**
  - Need comprehensive tests for: email, url, minlen, maxlen, min, max, regex, nonempty, combined rules, custom RegisterValidator, empty/missing tags, integration with ValidateConfig and CLI execution

- **Lint Warnings (23 warnings in `flags_validate.go`):**
  - `nilerr`: Silent nil return on parse errors in min/max/minlen/maxlen validators (intentional — treat malformed params as no-op, but should be documented or logged)
  - `err113`: Dynamic errors instead of sentinel errors (validation error messages include values, so dynamic is appropriate, but should wrap a sentinel)

---

## C. NOT STARTED 📝 (see Section F for full list)

### 10. Custom Type Registry Refactor

- Currently `flags.go`, `flags_parse.go`, `config_setfield.go`, `config_parsing.go` each have independent switch statements matching `reflect.Type` for Duration, Enum, LogLevel, LogFormat, URL, Email, Port, FilePath, HostPort
- Planned: `TypeRegistry` that maps `reflect.Type` → parser/flag-factory, with `RegisterCustomType` for user extensibility
- **Not started.**

---

## D. TOTALLY FUCKED UP 💥 (0 items)

Nothing is broken. All 14 packages pass `go test ./... -count=1 -timeout 120s -race`.
Build is clean (`go build ./...` exits 0).

---

## E. WHAT WE SHOULD IMPROVE

### Critical

1. **validate tag tests** — 296 lines of untested validator code. This is the biggest gap.
2. **nilerr lint warnings** — The validators silently swallow malformed parameters (e.g., `min=abc`). Either return an error or document the behavior.
3. **err113 lint warnings** — Validation errors should wrap sentinel errors for `errors.Is()` support.

### Important

4. **wireHandler is unused** — `cli_command.go:122` has a dead function since all calls go through `wireHandlerWithMiddleware`. Should be removed.
5. **FEATURES.md / TODO_LIST.md outdated** — Don't reflect the new features (Result, middleware, groups, validate tags). Needs update.
6. **AGENTS.md outdated** — Doesn't document middleware, groups, validate tags, or Result type.

### Nice-to-Have

7. **Coverage regression** — v2 dropped from 87.9% → 81.3% with the new code. New files are dragging it down.
8. **validate tag documentation** — No docs for the `validate` struct tag format or available validators.
9. **v2 coverage gap in flags_validate.go** — Zero test coverage for the validation registry code.

---

## F. Top #25 Things We Should Get Done Next

### P0 — Correctness & Quality (must-do)

| #   | Task                                       | Estimate | Why                                      |
| --- | ------------------------------------------ | -------- | ---------------------------------------- |
| 1   | Write comprehensive validate tag tests     | 30min    | 296 lines of untested production code    |
| 2   | Fix nilerr warnings in validators          | 10min    | Silent error swallowing is a bug smell   |
| 3   | Fix err113 warnings — wrap sentinel errors | 10min    | Enables errors.Is() for callers          |
| 4   | Remove dead `wireHandler` function         | 2min     | Dead code is confusing                   |
| 5   | Run golangci-lint and fix all issues       | 20min    | 154 warnings currently (mostly existing) |

### P1 — Documentation & Maintenance (should-do)

| #   | Task                                        | Estimate | Why                                                      |
| --- | ------------------------------------------- | -------- | -------------------------------------------------------- |
| 6   | Update FEATURES.md with new features        | 15min    | Currently lists middleware/groups/validate as "planned"  |
| 7   | Update TODO_LIST.md — mark completed items  | 10min    | Result[T], command groups, middleware, validate all done |
| 8   | Update AGENTS.md with new API surface       | 15min    | Missing: middleware, groups, validate, Result[T]         |
| 9   | Document validate tag format in README/docs | 10min    | Users don't know this exists                             |
| 10  | Add godoc examples for new public APIs      | 20min    | RegisterValidator, WithMiddleware, WithGroup, Result     |

### P2 — Feature Completion (good-to-do)

| #   | Task                                               | Estimate | Why                                         |
| --- | -------------------------------------------------- | -------- | ------------------------------------------- |
| 11  | Refactor custom type parsing to TypeRegistry       | 60min    | Eliminates 4+ duplicated switch statements  |
| 12  | Allow user `RegisterCustomType[T]()` extensibility | 30min    | Registry pattern enables this               |
| 13  | Add `WithGroup` for help/completion command groups | 15min    | Cobra supports `SetHelpCommandGroupID` etc. |
| 14  | Add `WithMiddlewareCommandFilter` option           | 20min    | Allow middleware only on specific commands  |
| 15  | Integration test: validate tags + CLI execution    | 15min    | End-to-end validation in real CLI flow      |

### P3 — Polish & DX (stretch)

| #   | Task                                            | Estimate | Why                                |
| --- | ----------------------------------------------- | -------- | ---------------------------------- |
| 16  | Add fuzz tests to config_parsing.go             | 20min    | Currently in TODO_LIST             |
| 17  | Add performance benchmarks for middleware chain | 15min    | Ensure zero-overhead promise holds |
| 18  | Add `Result.Try(f)` for wrapping panics         | 10min    | Natural extension of Result type   |
| 19  | Shell completion helpers                        | 60min    | Listed in TODO_LIST v3.0           |
| 20  | Config file auto-loading with koanf             | 60min    | Listed in TODO_LIST v3.0           |

### P4 — Release & Infrastructure (eventually)

| #   | Task                                         | Estimate | Why                            |
| --- | -------------------------------------------- | -------- | ------------------------------ |
| 21  | Fix pre-commit hooks (5 pre-existing errors) | 30min    | Currently requires --no-verify |
| 22  | Set up CI pipeline (lint + test + coverage)  | 30min    | No CI exists                   |
| 23  | Add codecov integration                      | 15min    | Track coverage trends          |
| 24  | Create v2.2.0 release tag and notes          | 20min    | Ship the new features          |
| 25  | Add benchmark regression detection           | 15min    | Performance guardrail          |

---

## G. Top #1 Question I Cannot Figure Out Myself

**Should the validate tag feature use the same "parameter:value" encoding scheme (e.g., `minlen=3` becomes `validateMinLen("3:hello")`) or should we switch to a cleaner API like `ValidateFunc func(reflect.Value) error` that receives the actual typed value instead of a string?**

The current approach converts everything to strings and back, which:

- Works for all types (string, int, float) via `formatFieldValue`
- But loses type information and requires re-parsing in validators
- Makes `min`/`max` validators parse floats from strings twice

A typed validator API would be cleaner but would require users to handle type assertions. The current approach is simpler for the common case (string validators).

---

## Project Stats

| Metric                 | Value                                                                 |
| ---------------------- | --------------------------------------------------------------------- |
| v2 source files        | 28                                                                    |
| v2 test files          | 55                                                                    |
| v2 source lines        | 4,896                                                                 |
| v2 test lines          | 9,797                                                                 |
| Total packages         | 14                                                                    |
| Test pass rate         | 100% (all 14 packages)                                                |
| Race detection         | Clean                                                                 |
| v2 coverage            | 81.3%                                                                 |
| Lint warnings          | 154 (0 errors)                                                        |
| New code this session  | +2,242 lines across 17 files                                          |
| Correctness bugs fixed | 4 (default parsing, MergeConfigs mutation, shared map, NoFlags alias) |
| New features added     | 4 (Result[T], middleware, command groups, validate tags)              |

---

_This report was generated by a comprehensive codebase audit across 83 source and test files._
