# cmdguard — Comprehensive Status Report

**Date:** 2026-06-08 09:46 CEST
**Version:** v2.4.0 (post-release maintenance, HEAD at 7d15f41)
**Reporter:** Crush (AI)
**Go:** 1.26.3 | **Platform:** linux

---

## Executive Summary

Since the last status report (2026-06-08 05:53), one full sprint was executed: 12 commits fixing a critical `flake.nix` regression, adding a `DoctorCommand[T]` convenience helper, DRY-ing the configload package, and closing all audit findings. A second audit found 7 additional issues (stale docs, misplaced sentinel error, incomplete references) — all now fixed.

**The codebase is clean.** 364 tests pass, 82.9% coverage on core, 87.5% on configload, 0 lint issues, 0 race conditions, `nix flake check` passes.

---

## a) FULLY DONE ✅

### Infrastructure

<!-- check-rows:ledger -->

| Item                                                          | Status              |
| ------------------------------------------------------------- | ------------------- |
| `flake.nix` infinite recursion fixed (`goPkg = pkgs.go_1_26`) | ✅ Commit `c6ee9c2` |
| `flake.nix` enhanced with gofumpt + goimports                 | ✅ Commit `48d46ff` |
| `nix flake check` passes                                      | ✅ Verified         |

### New Features

<!-- check-rows:ledger -->

| Item                                                                                                                         | Status              |
| ---------------------------------------------------------------------------------------------------------------------------- | ------------------- |
| `Scope.HealthCheckResults()` / `HealthCheckResultsWithContext(ctx)` returning `map[string]error`                             | ✅ Commit `7bf23b3` |
| `CLI.HealthCheckResults()` / `HealthCheckResultsWithContext(ctx)`                                                            | ✅ Commit `7bf23b3` |
| `DoctorCommand[T]` / `MustDoctorCommand[T]` with `WithDoctorCheck`, `WithDoctorShort`, `WithDoctorLong`, `WithDoctorGroupID` | ✅ Commit `7148aeb` |
| `ErrDoctorFailed` sentinel in `errors.go`                                                                                    | ✅ Commit `906767b` |

### Refactoring

<!-- check-rows:ledger -->

| Item                                                      | Status              |
| --------------------------------------------------------- | ------------------- |
| configload: 3 files → 1 `loader.go` with `genericLoader`  | ✅ Commit `2785467` |
| `command_suggest.go` consolidated into `flags_suggest.go` | ✅ Commit `135a868` |

### Testing

<!-- check-rows:ledger -->

| Package                      | Tests        | Coverage |
| ---------------------------- | ------------ | -------- |
| `pkg/cmdguard/v2`            | 364 total    | 82.9%    |
| `pkg/cmdguard/v2/configload` | 22 new tests | 87.5%    |
| `pkg/cmdguard/v2/testutil`   | existing     | 87.5%    |
| Doctor tests                 | 9 new tests  | —        |
| Example taskctl              | all pass     | 71.6%    |

### Documentation

<!-- check-rows:ledger -->

| Document                                                                           | Status |
| ---------------------------------------------------------------------------------- | ------ |
| All test counts updated (357/356 → 364)                                            | ✅     |
| All coverage percentages updated (82.8% → 82.9%)                                   | ✅     |
| doctor.go added to AGENTS.md project tree                                          | ✅     |
| DoctorCommand section added to AGENTS.md API ref                                   | ✅     |
| HealthCheckResults in AGENTS.md CLI methods table                                  | ✅     |
| FEATURES.md: Doctor Command + Health Check Results sections                        | ✅     |
| FEATURES.md: configload coverage row added                                         | ✅     |
| examples/taskctl: health → doctor in main.go, README.md, commands.go, main_test.go | ✅     |
| Gotchas #36-37 added to AGENTS.md                                                  | ✅     |
| TODO_LIST.md: Phase 10 completed                                                   | ✅     |

---

## b) PARTIALLY DONE ⚠️

Nothing is partially done. All items are either fully complete or not started.

---

## c) NOT STARTED 📝

| Item                                                                                                                                                                   | Priority       | Notes                               |
| ---------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------- | ----------------------------------- |
| ~~Add `CODECOV_TOKEN` secret to GitHub~~ **Won't implement — tokenless Codecov upload works; CI coverage step is green.**                                              | ~~CI blocker~~ | ~~Requires repo admin access~~      |
| ~~Plugin system for custom validators/type handlers~~ done — shipped — plugin.go, RegisterPlugin/WithPlugin (FEATURES.md)                                              | ~~v3.0+~~      | ~~Future~~                          |
| ~~Config file nested struct support~~ done — nested config structs supported (AGENTS.md)                                                                               | ~~v3.0+~~      | ~~Future~~                          |
| ~~v3.0 API-breaking cleanup (NoFlags, Get/MustGet, RegisterInScope, Package)~~ **Superseded — shipped — v3.0.0 (2026-07-07) + v4.0.0 (2026-07-28) breaking cleanups.** | ~~v3.0+~~      | ~~Future~~                          |
| ~~Fuzz test seed corpus in `testdata/fuzz/`~~ done — testdata/fuzz seeds shipped; expansion tracked in ROADMAP deferred #7                                             | ~~Medium~~     | ~~7 fuzz targets exist, no corpus~~ |
| ~~Structured JSON error output for `--output=json`~~ done — cli_errors_json.go structured JSON error output                                                            | ~~Medium~~     | ~~Feature enhancement~~             |
| ~~Contribution guide (`CONTRIBUTING.md`)~~ done — CONTRIBUTING.md exists at repo root                                                                                  | ~~Low~~        | ~~Community~~                       |
| ~~Issue/PR templates (`.github/`)~~ done — .github/ISSUE_TEMPLATE/ + PULL_REQUEST_TEMPLATE.md exist                                                                    | ~~Low~~        | ~~Community~~                       |
| charmbracelet/log integration **→ ROADMAP raw ideas (charmbracelet/log integration)**                                                                                  | Low            | Logging                             |
| Metrics/hooks for custom observability **→ ROADMAP raw ideas (metrics/hooks observability)**                                                                           | Low            | Observability                       |

---

## d) TOTALLY FUCKED UP 💥

**Nothing is broken.** All checks pass:

- `go test ./... -race` → 364 PASS, 0 FAIL
- `golangci-lint run ./...` → 0 issues
- `go vet ./...` → clean
- `nix flake check` → all checks passed
- `go build ./...` → clean

The only pre-existing issue is in `pkg/testutil/panic_test_helpers.go` (type assertion on generic `T`), which predates this session.

---

## e) WHAT WE SHOULD IMPROVE 🔧

### High Priority

1. ~~**Fuzz test corpus** — 7 fuzz targets exist but have no seed corpus in `testdata/fuzz/`. Adding even minimal corpus files would significantly improve fuzz coverage.~~ done — testdata/fuzz seeds shipped; corpus expansion tracked in ROADMAP deferred #7
2. ~~**Structured JSON error output** — `--output=json` renders data but errors still go to stderr as plain text. A structured error format would help programmatic consumers.~~ done — cli_errors_json.go structured JSON error output

### Medium Priority

3. ~~**Test count for benchmarks** — benchmarks have 19 functions but the `benchmarks/` package reports `[no tests to run]`. Could add a smoke test.~~ **NOT-DO — accepted — benchmark packages have no unit tests by design; benchmarks run in CI via -bench.**
4. ~~**`checks.build` in flake.nix** — Tried but failed because sandboxed Nix can't download Go modules. Would need `buildGoModule` or a Go vendor directory.~~ **NOT-DO — AGENTS.md design decision — no build/vet checks in flake (sandbox cannot fetch modules).**
5. ~~**`WithColor` deprecation** — Still exported, still works. v3.0 removal.~~ **Superseded — WithColor removed in the v3.0.0 breaking cleanup (2026-07-07).**

### Low Priority

6. ~~**Contribution guide** — `CONTRIBUTING.md` with PR/issue templates.~~ done — CONTRIBUTING.md exists at repo root
7. **charmbracelet/log integration** — Replace any `fmt.Println` with structured logging. **→ ROADMAP raw ideas (charmbracelet/log integration)**
8. **Extract flag code** — `flagtags` standalone library for reuse outside cmdguard. **→ ROADMAP raw ideas (extract flagtags standalone library)**

---

## f) Top 25 Things We Should Get Done Next

| #      | Item                                                                                                                                  | Impact                     | Effort      | Tier           |
| ------ | ------------------------------------------------------------------------------------------------------------------------------------- | -------------------------- | ----------- | -------------- |
| ~~1~~  | ~~Add `CODECOV_TOKEN` to repo~~ **Won't implement — tokenless Codecov upload works; CI coverage step is green.**                      | ~~CI completeness~~        | ~~2min~~    | ~~Now~~        |
| ~~2~~  | ~~Fuzz test corpus for 7 targets~~ done — testdata/fuzz seeds shipped; corpus expansion tracked in ROADMAP deferred #7                | ~~Robustness~~             | ~~Medium~~  | ~~This week~~  |
| ~~3~~  | ~~Structured JSON error output~~ done — cli_errors_json.go structured JSON error output                                               | ~~Feature polish~~         | ~~Medium~~  | ~~This week~~  |
| ~~4~~  | ~~Smoke test for benchmarks package~~ **NOT-DO — benchmark packages carry no unit tests by design; benchmarks run in CI via -bench.** | ~~Testing gap~~            | ~~Small~~   | ~~This week~~  |
| ~~5~~  | ~~`checks.build` via `buildGoModule` in flake.nix~~ **NOT-DO — AGENTS.md design decision — no build/vet checks in flake.**            | ~~CI completeness~~        | ~~Medium~~  | ~~This month~~ |
| 6      | charmbracelet/log integration **→ ROADMAP raw ideas (charmbracelet/log integration)**                                                 | Modern logging             | Small       | This month     |
| ~~7~~  | ~~CONTRIBUTING.md~~ done — CONTRIBUTING.md exists at repo root                                                                        | ~~Community~~              | ~~Small~~   | ~~This month~~ |
| ~~8~~  | ~~Issue/PR templates~~ done — .github/ISSUE_TEMPLATE/ + PULL_REQUEST_TEMPLATE.md exist                                                | ~~Community~~              | ~~Small~~   | ~~This month~~ |
| ~~9~~  | ~~Test all examples in CI~~ **Superseded — superseded — CI examples job re-tracked in TODO_LIST (2026-11-07 harvest).**               | ~~CI reliability~~         | ~~Small~~   | ~~This week~~  |
| 10     | `Result[T]` type for error handling **→ ROADMAP raw ideas (Result[T] monad)**                                                         | Type safety                | Medium      | Next quarter   |
| 11     | `Validated[T]` wrapper with validation **→ ROADMAP raw ideas (Validated[T] wrapper)**                                                 | Type safety                | Medium      | Next quarter   |
| ~~12~~ | ~~Plugin system design document~~ **Superseded — plugin system shipped directly (plugin.go), no separate design doc.**                | ~~Architecture~~           | ~~Large~~   | ~~v3.0~~       |
| ~~13~~ | ~~Config file nested struct support~~ done — nested config structs supported (AGENTS.md)                                              | ~~Most requested~~         | ~~Large~~   | ~~v3.0~~       |
| ~~14~~ | ~~v3.0 API design document~~ **Superseded — v3.0.0 shipped 2026-07-07 (docs/planning/2026-07-06_06-54).**                             | ~~Foundation~~             | ~~Large~~   | ~~v3.0~~       |
| ~~15~~ | ~~Make NoFlags distinct named type~~ done — NoFlags is a distinct named type in v4 (AGENTS.md gotcha)                                 | ~~Type safety (breaking)~~ | ~~Medium~~  | ~~v3.0~~       |
| 16     | Rename Get[T]/MustGet[T] **→ ROADMAP v5 candidate #6; MustGet no longer exists in v4**                                                | Clarity (breaking)         | Medium      | v3.0           |
| 17     | Make RegisterInScope generic **→ TODO_LIST T5 (TODO(v5), scope.go)**                                                                  | Type safety (breaking)     | Medium      | v3.0           |
| ~~18~~ | ~~Remove WithColor deprecation~~ **Superseded — WithColor removed in the v3.0.0 breaking cleanup.**                                   | ~~Cleanup (breaking)~~     | ~~Trivial~~ | ~~v3.0~~       |
| ~~19~~ | ~~Remove IsExecutable deprecation~~ **Superseded — IsExecutable removed in the v3.0.0 breaking cleanup.**                             | ~~Cleanup (breaking)~~     | ~~Trivial~~ | ~~v3.0~~       |
| 20     | Extract flagtags standalone library **→ ROADMAP raw ideas (extract flagtags standalone library)**                                     | Reusability                | Large       | Future         |
| ~~21~~ | ~~go-output as optional sub-package~~ **Superseded — heavy deps moved to sub-modules in v3; go-output stayed core by design.**        | ~~Dependency hygiene~~     | ~~Large~~   | ~~v3.0~~       |
| ~~22~~ | ~~Documentation generation from CLI~~ done — GenerateDocs/docgen.go ships (FEATURES.md)                                               | ~~Self-documenting~~       | ~~Large~~   | ~~Future~~     |
| 23     | Metrics/hooks for observability **→ ROADMAP raw ideas (metrics/hooks observability)**                                                 | Production readiness       | Medium      | Future         |
| ~~24~~ | ~~BDD tests with ginkgo~~ **NOT-DO — decided against — cmdguard standardizes on table-driven stdlib tests (AGENTS.md Testing).**      | ~~Test quality~~           | ~~Medium~~  | ~~Future~~     |
| ~~25~~ | ~~Deprecate v1 API timeline~~ **Superseded — the v1 package no longer exists — pkg/ contains only v4 (v4.0.0 module migration).**     | ~~Process~~                | ~~Trivial~~ | ~~v3.0~~       |

---

## g) Top #1 Question I Cannot Figure Out Myself

**Should `DoctorCommand` output use the `output.go` system (OutputTable/OutputResult) instead of raw `fmt.Fprintf`?**

Currently, `doctor` writes ✓/✗ lines directly via `fmt.Fprintf`. This means it can't benefit from `--output=json`, `--output=yaml`, etc. But doctor output is fundamentally a list of status checks, which maps naturally to a table. If we wired it through `OutputTable`, users could get structured health data in any of the 12 supported formats.

The tradeoff: this would couple the doctor command to the output system (which brings in the `go-output` dependency). Currently doctor.go has zero imports beyond stdlib. Is the structured output worth the coupling?

**Routed to ROADMAP** (raw ideas, 2026-11-07): doctor still writes plain text in v4; the coupling question was never decided — parked as an idea, revisit on demand.

---

## Metrics Summary

<!-- check-rows:ledger -->

| Metric                        | Value                                              |
| ----------------------------- | -------------------------------------------------- |
| **Version**                   | v2.4.0                                             |
| **Source files**              | 48 (core v2)                                       |
| **Test files**                | 74 (core v2 + configload)                          |
| **Total LOC**                 | ~20,500 (core v2, all .go)                         |
| **Tests passing**             | 364 (0 failures)                                   |
| **Test coverage**             | 82.9% (core), 87.5% (configload), 87.5% (testutil) |
| **Fuzz targets**              | 7                                                  |
| **Benchmarks**                | 19                                                 |
| **Lint issues**               | 0                                                  |
| **Race conditions**           | 0                                                  |
| **Sentinel errors**           | 36+                                                |
| **Value types**               | 9                                                  |
| **Output formats**            | 12                                                 |
| **CLI options**               | 17                                                 |
| **Command options**           | 21                                                 |
| **Dependencies**              | 8 direct                                           |
| **Build status**              | ✅ Pass                                            |
| **Nix flake check**           | ✅ Pass                                            |
| **Commits since last report** | 13                                                 |

---

## Verdict

**The library is production-ready and spotless.** All critical issues from the first sprint are fixed. All audit findings from the second pass are resolved. Zero lint, zero race, zero build errors.

**Overall health: 9.8/10** (docked for pre-existing `pkg/testutil` issue and missing CODECOV_TOKEN).

## Resolution (2026-07-18)

Many P6 "v3 breaking change" items deferred here shipped in v3.0.0 (2026-07-07): `NoFlags` is now a distinct named type, `WithColor` was removed. The DoctorCommand and configload refactors landed on `/v2` and were carried into `/v3`. Coverage 82.9% → 87.6%.
