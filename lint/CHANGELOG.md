# Changelog

All notable changes to the `lint` sub-module (the cmdguard usage linter) are
documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
Dates are in YYYY-MM-DD format (ISO 8601).

## [Unreleased]

### Added

- Nothing yet.

## [0.2.0] - 2026-11-08

Second release of the `lint` sub-module: ships the baseline ratchet, cross-file
name tracing, provider options, and the CG007/CG008 end-goal baseline rules.

### Added

- **Baseline/ratchet mode.** `--write-baseline` snapshots current findings into
  `.cmdguard-lint-baseline.json`; subsequent runs fail only on findings beyond
  the baseline (per rule+file, line-drift tolerant, stale entries reported so
  the ratchet can tighten). Library surface:
  `LoadBaseline`/`WriteBaseline`/`ApplyBaseline`.
- **Cross-file name tracing for CG003/CG004/CG006.** Package-level CLI
  variables, Execute errors, and constructor errors are unioned per directory
  group (test files never merge with prod files; function-local names never
  leak), so a constructor in one file plus its misuse in another is caught.
  FP budget held: cmdguard itself and examples stay at 0 findings.
- **Provider options.** The BuildFlow provider declares `enable` and `disable`
  (comma-separated rule IDs, same semantics as the CLI flags via the new
  `lint.FilterByIDs`).
- **Generated README rule table.** `rules --markdown` emits the rule table from
  `allRuleDefs()`; `rules_table_test.go` fails the build when the README
  section drifts.
- **CG007 missing short description, CG008 missing CLI version.** The
  end-goal baseline: a command without `WithShort` renders as a blank help
  line, and a CLI without `WithCLIVersion` cannot identify its build. Both are
  warnings, skip `_test.go` files and calls with spread option lists, and are
  suppressible with the usual `//cmdguard-lint:ignore`.

### Changed

- Dependency updates: go-finding v1.15.1, go-finding/toolsdk v1.16.0,
  go-linter-sdk v0.4.0, cmdguard v4 core pinned to v4.2.1.

## [0.1.0] - 2026-10-04

First release of the `lint` sub-module: a usage linter for cmdguard-based CLIs.

### Added

- 6 syntactic rules (CG001–CG006): execute bypass (`fang.Execute` around
  cmdguard), stale major import, constructor panic (`MustNew*`), runtime
  `SetVersion`, duplicate version options, execute-error reprint
- Pure `go/parser` AST walk — no type info, no build required
- BuildFlow provider (`lint/provider`, toolsdk spec `cmdguard-lint`, module
  fan-out)
- Dogfood CLI `cmd/cmdguard-lint` (built with cmdguard v4 itself) with text,
  JSON, and SARIF 2.1.0 output (round-trips via `finding.FindingsFromSARIF`)
- Suppressions: `//cmdguard-lint:ignore <ID> <reason>` (reason required) on the
  finding line or the line above
- Per-directory memoized analysis with `ClearCache()` for long-lived processes
- `CurrentMajor` constant tracks the latest cmdguard major for CG002
- Validated: timesheets corpus 18/18 findings at hand-audited locations, 0
  false positives on cmdguard/erraudit/go-structure-linter/branching-flow
