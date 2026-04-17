# cmdguard — Full Comprehensive Status Report

**Date:** 2026-04-15 23:36 CEST
**Branch:** master (clean working tree)
**Last Commit:** `0d6fc58` — feat: add command groups, middleware, Result type, and flag validation
**Go Version:** 1.26 | **Tests:** 14/14 PASS with `-race` | **Coverage:** v2 at 81.3%
**Total Lines:** 14,693 (28 source files @ 4,896 lines, 55 test files @ 9,797 lines)

---

## A. FULLY DONE ✅ (9 items)

### 1. Result[T] Type

- **Files:** `types_result.go` (237 lines) + `types_result_test.go` (639 lines)
- Rust-style `Ok(value)` / `Err(error)` discriminated union
- Full method suite: `Map`, `MapErr`, `MapOr`, `And`, `Or`, `IfOk`, `IfErr`, `ToOption`, `ResultFrom`, `ToPair`, `Unwrap`, `UnwrapOr`, `UnwrapOrElse`, `UnwrapErr`, `Expect`, `ExpectErr`, `MarshalJSON`, `String`
- 639 lines of tests covering all methods, edge cases, panic paths

### 2. NoFlags Type Safety Fix

- **File:** `command.go`
- Changed `type NoFlags = struct{}` (alias) → `type NoFlags struct{}` (defined type)
- Prevents accidental comparison with bare `struct{}{}`, makes it a proper named type

### 3. Silent Error Swallowing Fix

- **Files:** `config_parsing.go`, `flags.go`
- `parseBoolDefault`, `parseIntDefault`, `parseUintDefault`, `parseFloat64Default` changed from void/error-ignoring to `(value, error)` return tuples
- `registerFlag` now propagates errors instead of silently using zero values
- **Bug fixed:** invalid `default:"not-a-number"` on int fields silently produced `0`

### 4. MergeConfigs Mutation Fix

- **Files:** `config.go`, `config_merge_test.go`
- `MergeConfigs[T]()` now deep-copies `configs[0]` before merging
- **Bug fixed:** first input was mutated in place (silent data corruption)
- Regression tests added

### 5. BranchingFlowContext Shared Map Fix

- **Files:** `flow_context.go`, `flow_context_value_test.go`
- `newChild()` now uses `maps.Clone(b.values)` instead of sharing reference
- **Bug fixed:** child `SetValueLocal` could contaminate parent values
- Regression test added

### 6. LogLevel/LogFormat Deduplication

- **File:** `types_log.go`
- Extracted `logLevelAllowed` and `logFormatAllowed` package-level slices
- Eliminated 4 duplicate slice literals that could drift apart

### 7. Middleware/Interceptor Chain

- **Files:** `middleware.go` (77 lines) + `middleware_test.go` (438 lines)
- `Middleware[T]` type: `func(ctx, cfg, info CommandInfo, next func() error) error`
- `buildChain` builds right-to-left so first middleware wraps outermost
- Built-in: `TimingMiddleware[T]`, `RecoveryMiddleware[T]`
- `WithMiddleware[T]` CLI option
- `CommandInfo` carries command metadata (name, phase, hasRunE)
- 438 lines of tests: chaining, error propagation, short-circuit, timing, recovery, subcommands

### 8. Command Groups

- **Files:** `cli_command.go`, `cli_options.go`, `cli_groups_test.go` (235 lines)
- `Group string` field on `Command[T, F]`
- `WithGroup[T](id, title)` CLI option registers `cobra.Group`
- 235 lines of tests: grouping, help output, no-group fallback, subcommands

### 9. Flag Validators Implementation (code only)

- **File:** `flags_validate.go` (296 lines)
- `FlagValidator func(value string) error` type
- `validatorRegistry` with `sync.RWMutex` for goroutine safety
- `RegisterValidator(name, validator)` public API
- 8 built-in validators: `email`, `url`, `minlen`, `maxlen`, `min`, `max`, `regex`, `nonempty`
- Wired into `validateTag` (config.go) and `validateTagRules` (flags.go)
- `FlagTag.Validate` field added, `parseFieldFlag` parses `validate` struct tag
- **⚠️ SEE SECTION B — tests are missing**

---

## B. PARTIALLY DONE ⚠️ (1 item)

### Flag Validators — Tests Missing

- **Implementation:** DONE (296 lines in `flags_validate.go`)
- **Tests:** ❌ ZERO tests for the new `validate` tag feature
  - Existing `flags_validate_test.go` only has pre-existing enum/required validation tests
  - Missing coverage for: email, url, minlen, maxlen, min, max, regex, nonempty, combined rules (`minlen=2,maxlen=10`), custom `RegisterValidator`, empty/missing tags, `ValidateConfig` integration, CLI execution flow integration, `parseValidateRules` unit tests
- **Lint Warnings (23 warnings in `flags_validate.go`):**
  - `nilerr`: Silent nil return on parse errors in min/max/minlen/maxlen validators
  - `err113`: Dynamic errors instead of sentinel-wrapped errors

---

## C. NOT STARTED 📝 (7 items)

| #   | Task                          | Scope                             | Notes                                              |
| --- | ----------------------------- | --------------------------------- | -------------------------------------------------- |
| 1   | **Validate tag tests**        | `validate_tag_test.go` (new file) | 296 lines of untested production code              |
| 2   | **Fix nilerr warnings**       | `flags_validate.go`               | Validators silently swallow malformed params       |
| 3   | **Fix err113 warnings**       | `flags_validate.go`               | Wrap sentinel errors for `errors.Is()`             |
| 4   | **Remove dead `wireHandler`** | `cli_command.go:122`              | All calls go through `wireHandlerWithMiddleware`   |
| 5   | **Update FEATURES.md**        | docs                              | Missing: Result, middleware, groups, validate tags |
| 6   | **Update TODO_LIST.md**       | docs                              | Mark completed items, update remaining             |
| 7   | **Update AGENTS.md**          | docs                              | Missing new API surface docs                       |

---

## D. TOTALLY FUCKED UP 💥 (0 items)

**Nothing is broken.** All 14 packages pass `go test ./... -count=1 -timeout 120s -race`.
Build is clean: `go build ./...` exits 0.

---

## E. WHAT WE SHOULD IMPROVE

### Critical (blocks release quality)

1. **validate tag tests** — 296 lines of untested code is the single biggest quality gap
2. **nilerr lint warnings** — Silent error swallowing in validators is a correctness smell
3. **err113 lint warnings** — Callers can't use `errors.Is()` on validation errors

### Important (blocks v2.2.0 release)

4. **Dead `wireHandler`** — Confusing dead code at `cli_command.go:122`
5. **Docs outdated** — FEATURES.md, TODO_LIST.md, AGENTS.md don't reflect new features
6. **Coverage regression** — v2 dropped from 87.9% → 81.3% with new untested code

### Nice-to-Have

7. **validate tag documentation** — No docs for tag format or available validators
8. **TypeRegistry refactor** — 4+ duplicated switch statements across `flags.go`, `flags_parse.go`, `config_setfield.go`, `config_parsing.go`
9. **Custom type extensibility** — `RegisterCustomType[T]()` for user-defined flag types
10. **Fuzz tests** — For `config_parsing.go` and `flags_parse.go`

---

## F. Top #25 Things We Should Get Done Next

### P0 — Correctness & Quality (must-do)

| #   | Task                                    | Estimate | Why                                    |
| --- | --------------------------------------- | -------- | -------------------------------------- |
| 1   | Write comprehensive validate tag tests  | 30min    | 296 lines of untested production code  |
| 2   | Fix nilerr warnings in validators       | 10min    | Silent error swallowing is a bug smell |
| 3   | Fix err113 — wrap sentinel errors       | 10min    | Enables `errors.Is()` for callers      |
| 4   | Remove dead `wireHandler` function      | 2min     | Dead code is confusing                 |
| 5   | Run golangci-lint, fix all issues in v2 | 20min    | Zero warnings target                   |

### P1 — Documentation & Maintenance (should-do)

| #   | Task                                        | Estimate | Why                                                  |
| --- | ------------------------------------------- | -------- | ---------------------------------------------------- |
| 6   | Update FEATURES.md with new features        | 15min    | Lists middleware/groups/validate as "planned"        |
| 7   | Update TODO_LIST.md — mark completed        | 10min    | Result[T], groups, middleware, validate all done     |
| 8   | Update AGENTS.md with new API surface       | 15min    | Missing: middleware, groups, validate, Result[T]     |
| 9   | Document validate tag format in README/docs | 10min    | Users don't know this exists                         |
| 10  | Add godoc examples for new public APIs      | 20min    | RegisterValidator, WithMiddleware, WithGroup, Result |

### P2 — Feature Completion (good-to-do)

| #   | Task                                               | Estimate | Why                                        |
| --- | -------------------------------------------------- | -------- | ------------------------------------------ |
| 11  | Refactor custom type parsing to TypeRegistry       | 60min    | Eliminates 4+ duplicated switch statements |
| 12  | Allow user `RegisterCustomType[T]()` extensibility | 30min    | Registry pattern enables this              |
| 13  | Integration test: validate tags + CLI execution    | 15min    | End-to-end validation in real CLI flow     |
| 14  | Add `WithMiddlewareCommandFilter` option           | 20min    | Allow middleware only on specific commands |
| 15  | Add `Result.Try(f)` for wrapping panics            | 10min    | Natural extension of Result type           |

### P3 — Polish & DX (stretch)

| #   | Task                                      | Estimate | Why                          |
| --- | ----------------------------------------- | -------- | ---------------------------- |
| 16  | Add fuzz tests to config_parsing.go       | 20min    | Currently in TODO_LIST       |
| 17  | Add performance benchmarks for middleware | 15min    | Ensure zero-overhead promise |
| 18  | Shell completion helpers                  | 60min    | Listed in TODO_LIST v3.0     |
| 19  | Config file auto-loading with koanf       | 60min    | Listed in TODO_LIST v3.0     |
| 20  | Progress/Spinner Type (bubble tea)        | 60min    | Listed in TODO_LIST v3.0     |

### P4 — Release & Infrastructure (eventually)

| #   | Task                                         | Estimate | Why                              |
| --- | -------------------------------------------- | -------- | -------------------------------- |
| 21  | Fix pre-commit hooks (5 pre-existing errors) | 30min    | Currently requires `--no-verify` |
| 22  | Set up CI pipeline (lint + test + coverage)  | 30min    | No CI exists                     |
| 23  | Add codecov integration                      | 15min    | Track coverage trends            |
| 24  | Create v2.2.0 release tag and notes          | 20min    | Ship the new features            |
| 25  | Add benchmark regression detection           | 15min    | Performance guardrail            |

---

## G. Top #1 Question I Cannot Figure Out Myself

**Should the validate tag feature keep the current "parameter:value" encoding scheme (e.g., `minlen=3` becomes `validateMinLen("3:hello")`) or switch to a cleaner `TypedValidator func(reflect.Value) error` that receives the actual typed value?**

**Current approach (string-based):**

- ✅ Works for all types via `formatFieldValue`
- ✅ Simple for users — one function signature
- ❌ Loses type info, requires re-parsing in validators
- ❌ `min`/`max` parse floats from strings twice

**Alternative approach (typed):**

- ✅ Type-safe, no double-parsing
- ✅ More idiomatic Go
- ❌ Users must handle type assertions
- ❌ Two validator signatures increases API surface

**My recommendation:** Keep string-based for v2.2, add typed validators as `TypedFlagValidator func(reflect.Value) error` in v3.0 as a separate registration path. Best of both worlds, no breaking change.

---

## Project Stats

| Metric                 | Value                        |
| ---------------------- | ---------------------------- |
| v2 source files        | 28                           |
| v2 test files          | 55                           |
| v2 source lines        | 4,896                        |
| v2 test lines          | 9,797                        |
| v2 total lines         | 14,693                       |
| Total packages         | 14                           |
| Test pass rate         | 100% (14/14)                 |
| Race detection         | Clean                        |
| v2 coverage            | 81.3%                        |
| New code this session  | ~2,242 lines across 17 files |
| Correctness bugs fixed | 4                            |
| New features added     | 4                            |
| Public API functions   | ~60+                         |
| Public API types       | ~20+                         |

---

## Public API Surface (v2)

### CLI Construction

`NewCLI[T]`, `AddCommand[T,F]`, `Execute`, `ExecuteWithArgs`, `ExecuteAndExit`

### CLI Options

`WithCLIVersion`, `WithCLILong`, `WithCLIScope`, `WithSilenceErrors`, `WithSilenceUsage`, `WithColor`, `WithMiddleware`, `WithGroup`

### Dependency Injection

`NewScope`, `Provide`, `ProvideValue`, `ProvideNamed`, `Invoke`, `InvokeNamed`, `MustInvoke`, `MustInvokeNamed`, `ScopedProvider`, `RegisterInScope`, `Package`

### Types

`CLI[T]`, `Command[T,F]`, `Scope`, `NoFlags`, `Option[T]`, `Result[T]`, `Enum[T]`, `LogLevel`, `LogFormat`, `Duration`, `URL`, `Email`, `Port`, `FilePath`, `HostPort`, `BranchingFlowContext`

### Validation

`RegisterValidator`, `FlagValidator`, `ValidateConfig`, `EnsureValid`, `FlagTypeConstraint`

### Error Types

`CommandError`, `ServiceError`, `FlagError`, `ErrInvalidCommand`, `ErrMissingHandler`, `ErrDuplicateCommand`, `ErrInvalidConfig`

### Middleware

`Middleware[T]`, `TimingMiddleware`, `RecoveryMiddleware`, `CommandInfo`

---

_Report generated across 83 source and test files._
