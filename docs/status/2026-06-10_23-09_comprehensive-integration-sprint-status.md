# Comprehensive Status Report — cmdguard v2.5.0 → v2.6.0

**Date:** 2026-06-10 23:09
**Session:** Post-sprint cleanup + library integration sprint
**Commits this session:** 4 (ahead of origin by 4)

---

## Metrics

<!-- check-rows:ledger -->
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

<!-- check-rows:ledger -->
| # | Commit    | Description                                                                                                                      |
| - | --------- | -------------------------------------------------------------------------------------------------------------------------------- |
| 1 | `2f69674` | `WithCLIVersion` auto-pipes to `fang.WithVersion`; new `WithCLICommit` auto-pipes to `fang.WithCommit`                           |
| 2 | `75b6ff9` | ADR-001: fang integration strategy — why we skip `WithNotifySignal`/`WithoutManpage`, what gaps remain                           |
| 3 | `8fa843d` | 4 new output formats (JSONL, TOML, AsciiDoc, PlantUML) + `WithFangErrorHandler` + `WithFangColorScheme` + all docs updated 12→16 |

### Prior Sessions (already committed)

<!-- check-rows:ledger -->
| # | Description                                                                                                               |
| - | ------------------------------------------------------------------------------------------------------------------------- |
| 1 | samber/do v2 utilization sprint: `WithGracefulShutdown`, `Override[T]`, `CloneScope`, `NewScopeWithOpts`, `WithDILogging` |
| 2 | Post-sprint cleanup: dead code removal, `RootScope()`, DI benchmarks, library research                                    |
| 3 | Zero panics: all 16 Must\* functions removed                                                                              |
| 4 | NoFlags as distinct named type                                                                                            |
| 5 | NO_COLOR fix (env var restored after execution)                                                                           |
| 6 | Audit log integration with `samber-do-auditlog`                                                                           |

---

## B) PARTIALLY DONE ⚠️

### Documentation Gaps (new APIs not in doc.go or docs/API.md)

<!-- check-rows:ledger -->
| API                                  | doc.go     | docs/API.md | FEATURES.md    |
| ------------------------------------ | ---------- | ----------- | -------------- |
| `WithCLICommit`                      | ❌ MISSING | ❌ MISSING  | ✅             |
| `WithFangErrorHandler`               | ❌ MISSING | ❌ MISSING  | ✅             |
| `WithFangColorScheme`                | ❌ MISSING | ❌ MISSING  | ✅             |
| `FormatJSONL/AsciiDoc/TOML/PlantUML` | ❌ MISSING | ❌ MISSING  | ✅             |
| Fang version integration note        | ❌ MISSING | ❌ MISSING  | ✅ (AGENTS.md) |

**Impact:** Users browsing godoc or API.md won't discover these new options. FEATURES.md is the only place they're documented.

### Coverage (17 functions at 0%)

<!-- check-rows:ledger -->
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

<!-- check-rows:ledger -->
| Library          | Grade | Gap                                                                              |
| ---------------- | ----- | -------------------------------------------------------------------------------- |
| fang/v2          | A     | `WithNotifySignal` intentionally skipped, `WithoutManpage` intentionally skipped |
| go-output v0.8.0 | A     | All 16 formats exposed                                                           |
| samber/do v2     | B+    | 43% utilization; remaining features are server patterns                          |

---

## C) NOT STARTED 🔲

| #  | Task                                                                                  | Priority | Effort                 |
| -- | ------------------------------------------------------------------------------------- | -------- | ---------------------- |
| ~~1~~  | ~~Document new APIs in `doc.go` and `docs/API.md`~~ **Superseded — rewritten repeatedly since; docs/API.md removed — options/format docs live in FEATURES/AGENTS.** | ~~HIGH~~ | ~~15m~~ |
| ~~2~~  | ~~Cover 17 functions at 0% (mainly validate, types, manpage)~~ done — core at 90.4% (FEATURES.md); named coverage tests shipped | ~~MEDIUM~~ | ~~30m~~ |
| ~~3~~  | ~~Deduplicate `validateEmail`/`validateURL` — delegate to `ParseEmail`/`ParseURL`~~ done — delegate to ParseEmail()/ParseURL() (AGENTS.md) | ~~MEDIUM~~ | ~~10m~~ |
| ~~4~~  | ~~Fix `HostPort.IsEmpty()` coupling — use `hp.port.IsEmpty()` instead of `hp.port.port`~~ done — fixed — hp.port.IsEmpty() (v4 types_hostport.go:69) | ~~LOW~~ | ~~2m~~ |
| ~~5~~  | ~~Remove unused `FlowContext` interface~~ done — removed — no FlowContext interface in v4 (flow_context_access.go) | ~~LOW~~ | ~~5m~~ |
| ~~6~~  | ~~Unified error types for all domain types (not just Enum/Duration)~~ **NOT-DO — NOT-DO — sentinel-based errors (ErrInvalidURL etc.) chosen in v4.** | ~~LOW~~ | ~~20m~~ |
| ~~7~~  | ~~Plugin system for custom validators~~ done — plugin.go — Plugin interface + RegisterPlugin/WithPlugin | ~~P1 (v3)~~ | ~~Large~~ |
| ~~8~~  | ~~Config file nested struct support~~ done — nested config structs supported (AGENTS.md) | ~~P1 (v3)~~ | ~~Large~~ |
| ~~9~~  | ~~Extract flag-related code to `flagtags` library~~ **Won't implement — idea dropped — not in 2026-11-02 ROADMAP restructure.** | ~~P1 (v3)~~ | ~~Large~~ |
| ~~10~~ | ~~CODECOV_TOKEN secret for GitHub~~ **Won't implement — tokenless Codecov upload works; CI coverage step is green.** | ~~P0~~ | ~~5m (needs repo access)~~ |

---

## D) TOTALLY FUCKED UP 💥

### Nothing catastrophic — but honest assessment:

1. ~~**doc.go and docs/API.md are stale** — We shipped `WithCLICommit`, `WithFangErrorHandler`, `WithFangColorScheme`, and 4 new format constants without updating the two most important documentation surfaces. This is the #1 gap right now.~~ **Superseded — rewritten repeatedly since; docs/API.md removed — FEATURES/AGENTS current.**

2. ~~**validateEmail/validateURL are duplicated** — `flags_validate.go` has standalone validators that reimplement the same logic as `ParseEmail`/`ParseURL` in `types_email.go`/`types_url.go`. This is a split brain.~~ done — delegate to ParseEmail()/ParseURL() (AGENTS.md)

3. ~~**Enum is NOT generic** — The doc comment says "Enum[T]" but it's a concrete struct. `LogLevel`/`LogFormat` are `type Xxx Enum` which means they lose all methods and must re-declare them. This is architectural debt that makes adding new enum types needlessly boilerplate-heavy.~~ **Won't implement — never adopted — Enum remains concrete struct in v4 (types_enum.go).**

4. ~~**TODO_LIST.md not updated** for the library integration sprint — still shows the old task list.~~ **Superseded — rewritten repeatedly since; TODO_LIST refreshed 2026-11.**

---

## E) WHAT WE SHOULD IMPROVE

### Architecture

1. ~~**Make Enum actually generic** — `Enum[T ~string]` would eliminate 4 method re-declarations per enum type. Adding a new enum would be a one-liner: `type Color Enum[string]`. This is the single biggest architecture improvement for the type system.~~ **Won't implement — never adopted — Enum remains concrete struct in v4 (types_enum.go).**

2. ~~**Deduplicate validators** — `flags_validate.go` validators should delegate to `ParseXxx()` functions in `types_*.go`. Currently email validation exists in two places with slightly different behavior.~~ done — delegate to ParseEmail()/ParseURL() (AGENTS.md)

3. ~~**Typed errors for all domain types** — Only `Enum` and `Duration` have typed error structs (`EnumError`, `DurationError`). URL, Email, Port, FilePath, HostPort just use `fmt.Errorf`. Inconsistent.~~ **NOT-DO — NOT-DO — sentinel-based errors chosen in v4.**

### Library Integration

4. ~~**go-output is at v0.8.0 but sub-packages at v0.7.0** — The main module is v0.8.0 but `serialization`, `markup`, `graph`, etc. are v0.7.0. Check if v0.7.0 sub-packages are compatible with v0.8.0 main module.~~ **Superseded — go-output external upstream (v0.30.4); version drift resolved upstream.**

5. ~~**samber-do-auditlog needs local replace** — Published v0.0.1 is missing `html_templ.go` (gitignored). Blocks clean `go get`.~~ done — published — samber-do-auditlog v0.12.1 from proxy; local replace removed

### Developer Experience

6. ~~**`VersionCommand` is a factory function, not a type** — The name sounds like a type. `NewVersionCommand` would follow the `NewCommand`/`NewParentCommand` convention. Breaking change for v3.~~ **Won't implement — not renamed — still VersionCommand factory in v4 (version.go:11).**

7. ~~**`GenerateManPageCommand` and `GenerateVersionCommand` are raw cobra helpers** — They break the typed Command[T,F] pattern. Should either be removed (in favor of typed versions) or clearly documented as escape hatches.~~ **Superseded — removed in v4 — man-page generation dropped (mango indirect only); VersionCommand is typed.**

---

## F) Top 25 Things We Should Get Done Next

Sorted by: impact × customer-value ÷ effort

| #  | Task                                                                                                  | Impact | Effort | Category     |
| -- | ----------------------------------------------------------------------------------------------------- | ------ | ------ | ------------ |
| ~~1~~  | ~~Document `WithCLICommit`, `WithFangErrorHandler`, `WithFangColorScheme` in `doc.go` and `docs/API.md`~~ **Superseded — rewritten repeatedly; docs/API.md removed — options documented in FEATURES/AGENTS (cli_options.go:28,36,96,146).** | ~~HIGH~~ | ~~10m~~ | ~~Docs~~ |
| ~~2~~  | ~~Document 4 new format constants in `doc.go` and `docs/API.md`~~ **Superseded — rewritten repeatedly — FEATURES.md documents all 16 formats.** | ~~HIGH~~ | ~~5m~~ | ~~Docs~~ |
| ~~3~~  | ~~Add fang version integration note to `doc.go` (WithCLIVersion auto-pipes to fang)~~ **Superseded — rewritten repeatedly — AGENTS/FEATURES carry the fang version notes.** | ~~MEDIUM~~ | ~~3m~~ | ~~Docs~~ |
| ~~4~~  | ~~Deduplicate `validateEmail` → delegate to `ParseEmail()`~~ done — delegates to ParseEmail() (flags_validate.go:238) | ~~MEDIUM~~ | ~~5m~~ | ~~Architecture~~ |
| ~~5~~  | ~~Deduplicate `validateURL` → delegate to `ParseURL()`~~ done — delegates to ParseURL() (flags_validate.go:251) | ~~MEDIUM~~ | ~~5m~~ | ~~Architecture~~ |
| ~~6~~  | ~~Fix `HostPort.IsEmpty()` — use `hp.port.IsEmpty()` not `hp.port.port`~~ done — fixed — hp.port.IsEmpty() (types_hostport.go:69) | ~~LOW~~ | ~~2m~~ | ~~Bug~~ |
| ~~7~~  | ~~Test `IsEmpty()` on Duration, LogLevel, LogFormat, Port (4 functions)~~ done — IsEmpty tests shipped — duration_test.go:123, types_port_test.go:94, helpers_test.go:64 | ~~MEDIUM~~ | ~~5m~~ | ~~Coverage~~ |
| ~~8~~  | ~~Test `RegisterValidator` (global validator registration)~~ done — flags_validate_test.go covers RegisterValidator | ~~MEDIUM~~ | ~~5m~~ | ~~Coverage~~ |
| ~~9~~  | ~~Test `WithAuditLogGroupID` option~~ **Superseded — symbol gone — auditlog options rewritten (WithAuditMiddleware in v4).** | ~~LOW~~ | ~~3m~~ | ~~Coverage~~ |
| ~~10~~ | ~~Test `WithConfigFileLoader` option~~ **Superseded — configload deleted; KoanfLoader is the single loader (AGENTS.md).** | ~~LOW~~ | ~~5m~~ | ~~Coverage~~ |
| ~~11~~ | ~~Test `WithDoctorLong` option~~ **Superseded — symbol rewritten in v4 (WithDoctorLong[T]); doctor covered by doctor_test.go.** | ~~LOW~~ | ~~3m~~ | ~~Coverage~~ |
| ~~12~~ | ~~Test `NewManPage` constructor~~ **Superseded — man-page generation removed in v4 (mango indirect only).** | ~~LOW~~ | ~~3m~~ | ~~Coverage~~ |
| ~~13~~ | ~~Cover `validateNonEmpty`, `validateFieldByKind`, `runValidateTagWithRegistry`~~ **Superseded — rewritten through v3/v4 — validateTagWithRegistry (config.go:88).** | ~~MEDIUM~~ | ~~10m~~ | ~~Coverage~~ |
| ~~14~~ | ~~Cover `validateMin`, `validateMax`, `validateMaxLen` (41-55%)~~ **Superseded — rewritten through v3/v4 — validateMin/validateMax/validateMaxLen reshaped.** | ~~LOW~~ | ~~5m~~ | ~~Coverage~~ |
| ~~15~~ | ~~Remove unused `FlowContext` interface from `flow_context_access.go`~~ done — removed — no FlowContext interface in v4 (flow_context_access.go) | ~~LOW~~ | ~~5m~~ | ~~Cleanup~~ |
| ~~16~~ | ~~Update `TODO_LIST.md` for current sprint~~ **Superseded — rewritten repeatedly since; TODO_LIST refreshed 2026-11.** | ~~MEDIUM~~ | ~~10m~~ | ~~Docs~~ |
| ~~17~~ | ~~Make `Enum` actually generic (`Enum[T ~string]`) — eliminate boilerplate for LogLevel/LogFormat~~ **Won't implement — never adopted — Enum concrete in v4 (types_enum.go).** | ~~HIGH~~ | ~~30m~~ | ~~Architecture~~ |
| ~~18~~ | ~~Add typed error structs for URL, Email, Port, FilePath, HostPort~~ **NOT-DO — NOT-DO — sentinel-based errors chosen in v4.** | ~~LOW~~ | ~~20m~~ | ~~Architecture~~ |
| ~~19~~ | ~~Check go-output sub-package version alignment (v0.7.0 vs v0.8.0)~~ **Superseded — go-output external upstream (v0.30.4); drift resolved upstream.** | ~~LOW~~ | ~~5m~~ | ~~Deps~~ |
| ~~20~~ | ~~Add `CODECOV_TOKEN` to GitHub repo settings~~ **Won't implement — tokenless Codecov upload works; CI coverage step is green.** | ~~LOW~~ | ~~5m~~ | ~~CI~~ |
| ~~21~~ | ~~Add structured JSON error output for `--output=json`~~ done — cli_errors_json.go structured JSON error output | ~~MEDIUM~~ | ~~20m~~ | ~~Feature~~ |
| ~~22~~ | ~~Add fuzz test corpus in `testdata/fuzz/`~~ done — testdata/fuzz seeds shipped; expansion tracked in ROADMAP deferred #7 | ~~LOW~~ | ~~15m~~ | ~~Testing~~ |
| ~~23~~ | ~~Update `docs/MIGRATION_FROM_COBRA.md` with fang version integration~~ done — docs/MIGRATION_FROM_COBRA.md exists, maintained through v4 doc passes | ~~LOW~~ | ~~5m~~ | ~~Docs~~ |
| ~~24~~ | ~~Rename `VersionCommand` → `NewVersionCommand` (breaking, defer to v3)~~ **Won't implement — not renamed — still VersionCommand factory (v4 version.go:11).** | ~~LOW~~ | ~~10m~~ | ~~Naming~~ |
| ~~25~~ | ~~Research `charm.land/log` as replacement for `WithDILogging` formatting~~ **Superseded — WithDILogging remains the DI logging hook (v4 cli_options.go:190); charm.land/log never adopted.** | ~~LOW~~ | ~~10m~~ | ~~Research~~ |

---

## G) Top #1 Question I Cannot Figure Out Myself

**Should `Enum` become generic (`Enum[T ~string]`) in v2.x or v3.0?**

Making `Enum` generic would eliminate the boilerplate pattern where `LogLevel`/`LogFormat` are `type Xxx Enum` and must re-declare 4 methods each. But it's a breaking API change:

- Current: `var e v2.Enum` → users can create raw `Enum` values
- Generic: `var e v2.Enum[string]` → more explicit but changes the API surface
- Alternative: Keep `Enum` concrete, add a `DefineEnum(name, values)` helper that returns a typed wrapper

The question is whether the boilerplate reduction (4 methods × N enum types) justifies the API complexity of generics for what is fundamentally a CLI library, not a types library. This is a design judgment call that needs the project owner's input.

## Resolution (2026-07-18)

Superseded by v3.0.0 (2026-07-07). Current state: 6 Go modules (core + glamour/manpage/prompts/spinner/telemetry), go-output v0.30.4, 87.6% coverage, 1429 test runs, 58 sentinel errors. The "Enum generic" open question (§G) and removal of `ErrNoFlags`/`ErrTooFewArgs`/`ErrTooManyArgs` were resolved in the v3 line (see `2026-06-12_03-33` report, items 4 & 20). Tests 395 → 1429; coverage 85.4% → 87.6%.
