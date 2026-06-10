# Comprehensive Status Report — cmdguard v2.5.0 → v2.6.0

**Date:** 2026-06-10 23:09
**Session:** Post-sprint cleanup + library integration sprint
**Commits this session:** 4 (ahead of origin by 4)

---

## Metrics

| Metric                    | Value                                    |
| ------------------------- | ---------------------------------------- |
| **Version**               | v2.5.0 (heading to v2.6.0)               |
| **Coverage**              | 85.4% (statements), 84.7% (main package) |
| **Tests**                 | 395 passing                              |
| **Lint issues**           | 0                                        |
| **Race conditions**       | 0                                        |
| **Build errors**          | 0                                        |
| **Source files**          | 127 (22,950 lines)                       |
| **Test files**            | 78                                       |
| **0%-coverage functions** | 17                                       |
| **Output formats**        | 16 (was 12)                              |
| **Fang integration**      | 6/8 symbols (was 2/8)                    |
| **samber/do utilization** | 43% (23/54 symbols)                      |

---

## A) FULLY DONE ✅

### This Session (4 commits)

| #   | Commit    | Description                                                                                                                      |
| --- | --------- | -------------------------------------------------------------------------------------------------------------------------------- |
| 1   | `2f69674` | `WithCLIVersion` auto-pipes to `fang.WithVersion`; new `WithCLICommit` auto-pipes to `fang.WithCommit`                           |
| 2   | `75b6ff9` | ADR-001: fang integration strategy — why we skip `WithNotifySignal`/`WithoutManpage`, what gaps remain                           |
| 3   | `8fa843d` | 4 new output formats (JSONL, TOML, AsciiDoc, PlantUML) + `WithFangErrorHandler` + `WithFangColorScheme` + all docs updated 12→16 |

### Prior Sessions (already committed)

| #   | Description                                                                                                               |
| --- | ------------------------------------------------------------------------------------------------------------------------- |
| 1   | samber/do v2 utilization sprint: `WithGracefulShutdown`, `Override[T]`, `CloneScope`, `NewScopeWithOpts`, `WithDILogging` |
| 2   | Post-sprint cleanup: dead code removal, `RootScope()`, DI benchmarks, library research                                    |
| 3   | Zero panics: all 16 Must\* functions removed                                                                              |
| 4   | NoFlags as distinct named type                                                                                            |
| 5   | NO_COLOR fix (env var restored after execution)                                                                           |
| 6   | Audit log integration with `samber-do-auditlog`                                                                           |

---

## B) PARTIALLY DONE ⚠️

### Documentation Gaps (new APIs not in doc.go or docs/API.md)

| API                                  | doc.go     | docs/API.md | FEATURES.md    |
| ------------------------------------ | ---------- | ----------- | -------------- |
| `WithCLICommit`                      | ❌ MISSING | ❌ MISSING  | ✅             |
| `WithFangErrorHandler`               | ❌ MISSING | ❌ MISSING  | ✅             |
| `WithFangColorScheme`                | ❌ MISSING | ❌ MISSING  | ✅             |
| `FormatJSONL/AsciiDoc/TOML/PlantUML` | ❌ MISSING | ❌ MISSING  | ✅             |
| Fang version integration note        | ❌ MISSING | ❌ MISSING  | ✅ (AGENTS.md) |

**Impact:** Users browsing godoc or API.md won't discover these new options. FEATURES.md is the only place they're documented.

### Coverage (17 functions at 0%)

| Function                     | File                  | Why 0%                                                     |
| ---------------------------- | --------------------- | ---------------------------------------------------------- |
| `WithAuditLogGroupID`        | auditlog.go:43        | Option never tested                                        |
| `WithConfigFileLoader`       | config_file.go:178    | Advanced option, no test                                   |
| `WithDoctorLong`             | doctor.go:34          | Option never tested                                        |
| `RegisterValidator`          | flags_validate.go:81  | Global validator registration, only tested via struct tags |
| `validateEmail`              | flags_validate.go:155 | Duplicates ParseEmail logic                                |
| `validateURL`                | flags_validate.go:168 | Duplicates ParseURL logic                                  |
| `validateNonEmpty`           | flags_validate.go:300 | Only called via runValidateTagWithRegistry                 |
| `validateFieldByKind`        | flags_validate.go:309 | Internal dispatch, never directly tested                   |
| `runValidateTagWithRegistry` | flags_validate.go:320 | Internal, never directly tested                            |
| `NewManPage`                 | manpage.go:63         | Constructor never tested                                   |
| `PromptString`               | prompts.go:23         | Interactive, hard to test (huh)                            |
| `PromptSelect`               | prompts.go:37         | Interactive, hard to test (huh)                            |
| `PromptConfirm`              | prompts.go:57         | Interactive, hard to test (huh)                            |
| `IsEmpty` (Duration)         | types_duration.go:44  | Never called in tests                                      |
| `IsEmpty` (LogLevel)         | types_log.go:65       | Never called in tests                                      |
| `IsEmpty` (LogFormat)        | types_log.go:98       | Never called in tests                                      |
| `IsEmpty` (Port)             | types_port.go:99      | Never called in tests                                      |

### Library Utilization

| Library          | Grade | Gap                                                                              |
| ---------------- | ----- | -------------------------------------------------------------------------------- |
| fang/v2          | A     | `WithNotifySignal` intentionally skipped, `WithoutManpage` intentionally skipped |
| go-output v0.8.0 | A     | All 16 formats exposed                                                           |
| samber/do v2     | B+    | 43% utilization; remaining features are server patterns                          |

---

## C) NOT STARTED 🔲

| #   | Task                                                                                  | Priority | Effort                 |
| --- | ------------------------------------------------------------------------------------- | -------- | ---------------------- |
| 1   | Document new APIs in `doc.go` and `docs/API.md`                                       | HIGH     | 15m                    |
| 2   | Cover 17 functions at 0% (mainly validate, types, manpage)                            | MEDIUM   | 30m                    |
| 3   | Deduplicate `validateEmail`/`validateURL` — delegate to `ParseEmail`/`ParseURL`       | MEDIUM   | 10m                    |
| 4   | Fix `HostPort.IsEmpty()` coupling — use `hp.port.IsEmpty()` instead of `hp.port.port` | LOW      | 2m                     |
| 5   | Remove unused `FlowContext` interface                                                 | LOW      | 5m                     |
| 6   | Unified error types for all domain types (not just Enum/Duration)                     | LOW      | 20m                    |
| 7   | Plugin system for custom validators                                                   | P1 (v3)  | Large                  |
| 8   | Config file nested struct support                                                     | P1 (v3)  | Large                  |
| 9   | Extract flag-related code to `flagtags` library                                       | P1 (v3)  | Large                  |
| 10  | CODECOV_TOKEN secret for GitHub                                                       | P0       | 5m (needs repo access) |

---

## D) TOTALLY FUCKED UP 💥

### Nothing catastrophic — but honest assessment:

1. **doc.go and docs/API.md are stale** — We shipped `WithCLICommit`, `WithFangErrorHandler`, `WithFangColorScheme`, and 4 new format constants without updating the two most important documentation surfaces. This is the #1 gap right now.

2. **validateEmail/validateURL are duplicated** — `flags_validate.go` has standalone validators that reimplement the same logic as `ParseEmail`/`ParseURL` in `types_email.go`/`types_url.go`. This is a split brain.

3. **Enum is NOT generic** — The doc comment says "Enum[T]" but it's a concrete struct. `LogLevel`/`LogFormat` are `type Xxx Enum` which means they lose all methods and must re-declare them. This is architectural debt that makes adding new enum types needlessly boilerplate-heavy.

4. **TODO_LIST.md not updated** for the library integration sprint — still shows the old task list.

---

## E) WHAT WE SHOULD IMPROVE

### Architecture

1. **Make Enum actually generic** — `Enum[T ~string]` would eliminate 4 method re-declarations per enum type. Adding a new enum would be a one-liner: `type Color Enum[string]`. This is the single biggest architecture improvement for the type system.

2. **Deduplicate validators** — `flags_validate.go` validators should delegate to `ParseXxx()` functions in `types_*.go`. Currently email validation exists in two places with slightly different behavior.

3. **Typed errors for all domain types** — Only `Enum` and `Duration` have typed error structs (`EnumError`, `DurationError`). URL, Email, Port, FilePath, HostPort just use `fmt.Errorf`. Inconsistent.

### Library Integration

4. **go-output is at v0.8.0 but sub-packages at v0.7.0** — The main module is v0.8.0 but `serialization`, `markup`, `graph`, etc. are v0.7.0. Check if v0.7.0 sub-packages are compatible with v0.8.0 main module.

5. **samber-do-auditlog needs local replace** — Published v0.0.1 is missing `html_templ.go` (gitignored). Blocks clean `go get`.

### Developer Experience

6. **`VersionCommand` is a factory function, not a type** — The name sounds like a type. `NewVersionCommand` would follow the `NewCommand`/`NewParentCommand` convention. Breaking change for v3.

7. **`GenerateManPageCommand` and `GenerateVersionCommand` are raw cobra helpers** — They break the typed Command[T,F] pattern. Should either be removed (in favor of typed versions) or clearly documented as escape hatches.

---

## F) Top 25 Things We Should Get Done Next

Sorted by: impact × customer-value ÷ effort

| #   | Task                                                                                                  | Impact | Effort | Category     |
| --- | ----------------------------------------------------------------------------------------------------- | ------ | ------ | ------------ |
| 1   | Document `WithCLICommit`, `WithFangErrorHandler`, `WithFangColorScheme` in `doc.go` and `docs/API.md` | HIGH   | 10m    | Docs         |
| 2   | Document 4 new format constants in `doc.go` and `docs/API.md`                                         | HIGH   | 5m     | Docs         |
| 3   | Add fang version integration note to `doc.go` (WithCLIVersion auto-pipes to fang)                     | MEDIUM | 3m     | Docs         |
| 4   | Deduplicate `validateEmail` → delegate to `ParseEmail()`                                              | MEDIUM | 5m     | Architecture |
| 5   | Deduplicate `validateURL` → delegate to `ParseURL()`                                                  | MEDIUM | 5m     | Architecture |
| 6   | Fix `HostPort.IsEmpty()` — use `hp.port.IsEmpty()` not `hp.port.port`                                 | LOW    | 2m     | Bug          |
| 7   | Test `IsEmpty()` on Duration, LogLevel, LogFormat, Port (4 functions)                                 | MEDIUM | 5m     | Coverage     |
| 8   | Test `RegisterValidator` (global validator registration)                                              | MEDIUM | 5m     | Coverage     |
| 9   | Test `WithAuditLogGroupID` option                                                                     | LOW    | 3m     | Coverage     |
| 10  | Test `WithConfigFileLoader` option                                                                    | LOW    | 5m     | Coverage     |
| 11  | Test `WithDoctorLong` option                                                                          | LOW    | 3m     | Coverage     |
| 12  | Test `NewManPage` constructor                                                                         | LOW    | 3m     | Coverage     |
| 13  | Cover `validateNonEmpty`, `validateFieldByKind`, `runValidateTagWithRegistry`                         | MEDIUM | 10m    | Coverage     |
| 14  | Cover `validateMin`, `validateMax`, `validateMaxLen` (41-55%)                                         | LOW    | 5m     | Coverage     |
| 15  | Remove unused `FlowContext` interface from `flow_context_access.go`                                   | LOW    | 5m     | Cleanup      |
| 16  | Update `TODO_LIST.md` for current sprint                                                              | MEDIUM | 10m    | Docs         |
| 17  | Make `Enum` actually generic (`Enum[T ~string]`) — eliminate boilerplate for LogLevel/LogFormat       | HIGH   | 30m    | Architecture |
| 18  | Add typed error structs for URL, Email, Port, FilePath, HostPort                                      | LOW    | 20m    | Architecture |
| 19  | Check go-output sub-package version alignment (v0.7.0 vs v0.8.0)                                      | LOW    | 5m     | Deps         |
| 20  | Add `CODECOV_TOKEN` to GitHub repo settings                                                           | LOW    | 5m     | CI           |
| 21  | Add structured JSON error output for `--output=json`                                                  | MEDIUM | 20m    | Feature      |
| 22  | Add fuzz test corpus in `testdata/fuzz/`                                                              | LOW    | 15m    | Testing      |
| 23  | Update `docs/MIGRATION_FROM_COBRA.md` with fang version integration                                   | LOW    | 5m     | Docs         |
| 24  | Rename `VersionCommand` → `NewVersionCommand` (breaking, defer to v3)                                 | LOW    | 10m    | Naming       |
| 25  | Research `charm.land/log` as replacement for `WithDILogging` formatting                               | LOW    | 10m    | Research     |

---

## G) Top #1 Question I Cannot Figure Out Myself

**Should `Enum` become generic (`Enum[T ~string]`) in v2.x or v3.0?**

Making `Enum` generic would eliminate the boilerplate pattern where `LogLevel`/`LogFormat` are `type Xxx Enum` and must re-declare 4 methods each. But it's a breaking API change:

- Current: `var e v2.Enum` → users can create raw `Enum` values
- Generic: `var e v2.Enum[string]` → more explicit but changes the API surface
- Alternative: Keep `Enum` concrete, add a `DefineEnum(name, values)` helper that returns a typed wrapper

The question is whether the boilerplate reduction (4 methods × N enum types) justifies the API complexity of generics for what is fundamentally a CLI library, not a types library. This is a design judgment call that needs the project owner's input.
