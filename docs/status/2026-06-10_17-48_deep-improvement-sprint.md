# Status Report — cmdguard v2.4.0 Deep Improvement Sprint

**Date:** 2026-06-10 17:48 CEST
**Branch:** master (pushed to origin)
**Author:** Crush (AI-assisted), Lars Artmann
**Scope:** Two-session deep improvement sprint spanning formatFieldValue tests through configload.Auto() fix

---

## Executive Summary

cmdguard v2.4.0 is in **excellent shape**. This sprint pushed 12 commits fixing 5 bugs, adding 3 features, refactoring 2 dead code paths, and strengthening error chains across the type system. The project went from 367 tests / 82.9% coverage to **374 tests / 84.0% coverage** with zero lint issues and zero race conditions.

The most impactful finding was `configload.Auto()` being completely broken — it always parsed as JSON regardless of file format. This is now fixed.

---

## Metrics Dashboard

| Metric                            | Before Sprint | After Sprint | Delta |
| --------------------------------- | ------------- | ------------ | ----- |
| Tests                             | 367           | 374          | +7    |
| Coverage                          | 82.9%         | 84.0%        | +1.1% |
| Lint issues                       | 0             | 0            | —     |
| Race conditions                   | 0             | 0            | —     |
| Source lines (v2)                 | ~21,400       | 21,424       | +24   |
| Test files                        | 71            | 73           | +2    |
| Sentinels reachable via errors.Is | ~51/57        | ~54/57       | +3    |
| Bugs fixed                        | —             | 5            | —     |
| Commits pushed                    | —             | 12           | —     |

---

## A) FULLY DONE ✅

### Bugs Fixed

| #   | Bug                                                                  | Root Cause                                                                                            | Fix                                                     | Commit    |
| --- | -------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- | ------------------------------------------------------- | --------- |
| 1   | `configload.Auto()` always used JSON                                 | `autoLoader.Load()` delegated to `JSON()` unconditionally                                             | Try YAML → TOML → JSON sequentially                     | `6d832e3` |
| 2   | `ShutdownAll` double-wrapped `ErrServiceConstruction`                | `Shutdown()` wraps once, `ShutdownAll` wrapped again                                                  | Removed redundant wrapping                              | `3e4c5f5` |
| 3   | `getFieldValue` incomplete — missed struct types with `fmt.Stringer` | Only checked `fmt.Stringer` for primitive kinds, returned `("", false)` for all structs               | Replaced with `formatFieldValue`                        | `090b9f6` |
| 4   | `ErrLogLevel`/`ErrLogFormat` unreachable                             | `ParseLogLevel`/`ParseLogFormat` passed through `ParseEnum` error without wrapping their own sentinel | Added `fmt.Errorf("%w: %w", ErrLogLevel, err)` wrapping | `670e8ea` |
| 5   | 3 bare sentinel returns without context                              | `ErrConfigNil` (×2) and `ErrConfigFileNotFound` returned bare                                         | Added `fmt.Errorf` wrapping with diagnostic context     | `670e8ea` |

### Features Added

| #   | Feature                                                                           | Commit    |
| --- | --------------------------------------------------------------------------------- | --------- |
| 1   | `MustParseDuration`, `MustParseLogLevel`, `MustParseLogFormat`                    | `7fe2940` |
| 2   | `MustParseEnum` (custom signature, not generic MustParse[T])                      | `903c8c5` |
| 3   | `registerGoDurationHandler` now validates non-empty defaults at registration time | `7fe2940` |

### Architecture Improvements

| #   | Improvement                                                                                                 | Commit    |
| --- | ----------------------------------------------------------------------------------------------------------- | --------- |
| 1   | `validatorRegistry` threaded through `ValidateConfig` path (internal functions accept `*validatorRegistry`) | `37d994c` |
| 2   | Removed dead `getFieldValue` function (replaced by `formatFieldValue`)                                      | `090b9f6` |
| 3   | Enum Marshal/Unmarshal documented as intentionally hand-written                                             | `86bf188` |

### Test Coverage Added

| File                        | Tests Added | What's Covered                                                            |
| --------------------------- | ----------- | ------------------------------------------------------------------------- |
| `flag_helpers_test.go`      | 10 subtests | `formatFieldValue` for all kinds, pointers, fmt.Stringer                  |
| `version_test.go`           | 3 tests     | `MustVersionCommand` success, `GenerateVersionCommand`                    |
| `duration_test.go`          | 2 tests     | `MustParseDuration` success + panic                                       |
| `helpers_test.go`           | 6 subtests  | `MustParseLogLevel`, `MustParseLogFormat`, `errors.Is` chain verification |
| `type_handler_test.go`      | 1 test      | `time.Duration` invalid default returns error                             |
| `enum_test.go`              | 2 tests     | `MustParseEnum` success + panic                                           |
| `configload/loader_test.go` | 3 tests     | Auto() YAML/TOML/JSON detection + invalid data                            |

### Documentation Updated

- `AGENTS.md` — 44 gotchas (was 37), metrics updated to 374/84.0%
- `TODO_LIST.md` — metrics updated
- `configload/loader.go` — `Auto()` doc rewritten to match new behavior

---

## B) PARTIALLY DONE ⚠️

### validatorRegistry Threading

**Status:** Plumbing complete, public API unchanged.

`validateStructWithRegistry`, `validateTagWithRegistry`, `validateFieldByKind` now accept `*validatorRegistry`. But `ValidateConfig()` (the public API) still uses `globalValidators` because it has no instance context. To fully realize instance-scoped validators:

- `FlagRegistry.ValidateFlags` needs to pass its instance `validatorRegistry` down
- OR a new public `ValidateConfigWithRegistry` API is needed

**What's left:** Wire instance-scoped validators through `FlagRegistry.ValidateFlags` → `validateTagRules` path.

---

## C) NOT STARTED 📋

### Known Gaps (from audit, sorted by impact)

| #   | Gap                                                                               | Impact                                                                           | Effort                          |
| --- | --------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | ------------------------------- |
| 1   | `RegisterInScope` can't accept typed providers (`func(do.Injector) (*T, error)`)  | High — function is essentially unusable with the library's own `Provide` pattern | 2h (needs Go generics redesign) |
| 2   | `Package()` panics on CLI creation failure                                        | Medium — contradicts "no panics" principle                                       | 30min                           |
| 3   | Slice handler only supports `[]string`                                            | Medium — `[]int`, `[]Port`, etc. are common                                      | 4h                              |
| 4   | Config file nested struct support                                                 | Medium — flat-only is limiting                                                   | 1d                              |
| 5   | `validateTagRules` has 16.7% coverage                                             | Low                                                                              | 1h                              |
| 6   | `completion.go` has 0% coverage (thin cobra wiring)                               | Low                                                                              | 30min                           |
| 7   | `RegisterValidator()` public API has 0% coverage                                  | Low                                                                              | 30min                           |
| 8   | `manpage.go` — `GenerateManPageCommand` at 14.3%, `NewManPage` at 0%              | Low                                                                              | 1h                              |
| 9   | `flow_context_access.go` — `Get[T]` and `MustGet[T]` at 0%                        | Low                                                                              | 30min                           |
| 10  | `prompts.go` — `PromptString`/`PromptSelect`/`PromptConfirm` at 0% (requires TTY) | Low                                                                              | 2h                              |

### Future Work (from TODO_LIST.md)

| #     | Task                                                                                                                                        | Priority |
| ----- | ------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| 20    | Add `CODECOV_TOKEN` secret to GitHub repo settings                                                                                          | P4       |
| 21-28 | Plugin system, nested configs, doc gen, advanced types, koanf, JSON errors, extract flagtags, extract go-output                             | P5       |
| 29-35 | NoFlags distinct type, rename Get/MustGet, generic RegisterInScope, remove Package, safe SetConfig, remove WithColor, fix NO_COLOR mutation | P6       |

---

## D) TOTALLY FUCKED UP 💥

**Nothing is fucked up.** The codebase is clean:

- 0 lint issues across 96 linters
- 0 race conditions
- 0 build errors
- All 374 tests pass with `-race`
- No TODO/FIXME/HACK comments in source code

**Previous session mistakes that were caught and reverted:**

- Attempted to replace `getField` with `derefPointerToStruct` — `getField` has different semantics (requires pointer for mutability). Reverted before commit. Lesson learned.

---

## E) WHAT WE SHOULD IMPROVE 🎯

### Architecture

1. **`RegisterInScope` is unusable** — accepts `func(do.Injector) (any, error)` only, which excludes typed providers. This is the biggest API wart. Needs a generic redesign (`RegisterInScope[T]`) — API-breaking, defer to v3.

2. **`Package()` panics** — `scope.go:374` panics on CLI creation failure. Documented but contradicts "no panics" principle. Should return `(*Scope, error)`.

3. **ConfigFileLoader interface lacks path** — `Load(data []byte, cfg any)` has no path, so `Auto()` can't use file extension detection. Our try-each approach works but is less precise. Consider adding `LoadWithPath(data []byte, cfg any, path string)` as optional interface.

### Type Model

4. **Duration has both IsZero() and IsEmpty()** — Every other type only has `IsEmpty()`. The extra `IsZero()` breaks the pattern. Keep for now (API compat) but note for v3.

5. **Enum UnmarshalText bootstrapping** — When `allowed` is empty, `UnmarshalText` accepts any value and sets `allowed = []string{value}`. This creates Enums in a state where `ParseEnum` would accept only that single value. Design smell but intentional.

6. **`setStringField` Enum bypass** — `config_setfield.go:126` has a special case for `Enum` that bypasses the `TypeHandler` registry. If someone registers a custom handler for `Enum`, it's silently ignored.

### Coverage

7. **28 functions at 0% coverage** — See section C above. Most are thin wrappers or TTY-dependent prompts.

### Documentation

8. **`validate` tag not documented in README** — The `validate:"min=3,max=100"` tag is mentioned in AGENTS.md but not in README.md. Users may not know it exists.

---

## F) TOP 25 THINGS TO DO NEXT

### Tier 1: High Impact, Low Effort (Do Now)

| #   | Task                                                       | Effort | Impact               |
| --- | ---------------------------------------------------------- | ------ | -------------------- |
| 1   | Add `CODECOV_TOKEN` to GitHub repo settings                | 5min   | CI coverage tracking |
| 2   | Add tests for `RegisterValidator()` public API             | 30min  | Coverage gap         |
| 3   | Add tests for `Get[T]`/`MustGet[T]` flow context accessors | 30min  | Coverage gap         |
| 4   | Add tests for `WithCompletion`/`WithValidArgs`             | 30min  | Coverage gap         |
| 5   | Document `validate` tag in README.md                       | 30min  | Discoverability      |
| 6   | Add `WithArgs` test (currently 0% coverage)                | 20min  | Coverage gap         |

### Tier 2: Medium Impact, Medium Effort (Do Soon)

| #   | Task                                                         | Effort | Impact                   |
| --- | ------------------------------------------------------------ | ------ | ------------------------ |
| 7   | Wire instance-scoped validators through `ValidateFlags` path | 2h     | Architecture correctness |
| 8   | Add `validateTagRules` tests (16.7% → 80%+)                  | 1h     | Validation coverage      |
| 9   | Add manpage tests (`GenerateManPageCommand`, `NewManPage`)   | 1h     | Coverage gap             |
| 10  | Add `renderAndWrite` output tests                            | 1h     | Coverage gap             |
| 11  | Add `WithDoctorLong` test                                    | 20min  | Coverage gap             |
| 12  | Add `WithConfigFileLoader` test                              | 30min  | Coverage gap             |
| 13  | Add `HealthCheckResults` (Scope) test                        | 20min  | Coverage gap             |
| 14  | Fix `Package()` to return error instead of panic             | 30min  | API correctness          |

### Tier 3: High Impact, High Effort (Plan for v3)

| #   | Task                                                        | Effort | Impact                  |
| --- | ----------------------------------------------------------- | ------ | ----------------------- |
| 15  | Generic `RegisterInScope[T]` redesign                       | 2h     | API usability           |
| 16  | Slice handler: support `[]int`, `[]Port`, etc.              | 4h     | Common feature request  |
| 17  | Config file nested struct support                           | 1d     | Major feature gap       |
| 18  | `ConfigFileLoader` path-aware interface                     | 2h     | Better auto-detection   |
| 18  | Remove `SetConfig` or make it safe                          | 2h     | API correctness         |
| 19  | Fix `os.Setenv("NO_COLOR", "1")` process-wide mutation      | 1h     | Side-effect elimination |
| 20  | Make `NoFlags` distinct named type                          | 1h     | Type safety             |
| 21  | Structured JSON error output for `--output=json`            | 4h     | CLI best practice       |
| 22  | Config auto-loading with koanf integration                  | 4h     | Ecosystem integration   |
| 23  | Extract flag-related code to standalone `flagtags` library  | 1d     | Reusability             |
| 24  | Plugin system for custom validators and type handlers       | 2d     | Extensibility           |
| 25  | Documentation generation (GenerateDocs, markdown, API docs) | 2d     | Discoverability         |

---

## G) TOP #1 QUESTION

**Can we break API for v3.0 now, or should we remain fully backward-compatible?**

Specifically:

- `RegisterInScope` is broken by design (typed providers don't match `func(do.Injector) (any, error)`). Fixing it requires API changes.
- `Package()` panics — fixing requires returning `(*Scope, error)`.
- `NoFlags` as type alias means `type NoFlags = struct{}` — users can accidentally pass any `struct{}`. Making it a distinct type is breaking.
- `ConfigFileLoader.Load()` lacks path — adding it is a breaking interface change.

If v3 is on the table, several Tier 3 items collapse into a single coordinated API cleanup. If not, we continue with backward-compatible additions only.

---

## Commit History (This Sprint)

```
48acdab docs(cmdguard): update AGENTS.md with final metrics and new gotchas
86bf188 docs(cmdguard): clarify why Enum has hand-written Marshal/UnmarshalText
903c8c5 feat(cmdguard): add MustParseEnum for API consistency
3e4c5f5 fix(cmdguard): remove double-wrapping of ErrServiceConstruction in ShutdownAll
090b9f6 refactor(cmdguard): replace getFieldValue with formatFieldValue
6d832e3 fix(configload): Auto() now tries YAML/TOML/JSON instead of only JSON
59f7dc1 docs(cmdguard): update AGENTS.md with error chain and unused sentinel notes
670e8ea fix(cmdguard): wire ErrLogLevel/ErrLogFormat into parse chain, fix bare sentinel returns
7fe2940 feat(cmdguard): add MustParse for Duration/LogLevel/LogFormat, fix GoDuration validation
37d994c fix(cmdguard): thread validatorRegistry through ValidateConfig path
e3cd57b test(cmdguard): add tests for MustVersionCommand and GenerateVersionCommand
b62aed4 test(cmdguard): add comprehensive tests for formatFieldValue
```
