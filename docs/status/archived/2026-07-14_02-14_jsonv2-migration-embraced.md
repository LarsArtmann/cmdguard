# Status Report: encoding/json/v2 Migration — Embraced

**Date:** 2026-07-14 02:14
**Session focus:** Fixing recurring BuildFlow failures caused by `go-auto-upgrade` trying to migrate to `encoding/json/v2`

---

> **Update 2026-07-23:** Sub-module lint issues were fixed in `da3b454`; the four test files were migrated to `encoding/json/v2` in `2a673a4`; `GOEXPERIMENT=jsonv2` is documented in `AGENTS.md` and `flake.nix`. The CHANGELOG now records the migration. The workspace now has 4 optional sub-modules; `manpage` was removed in `34a0c6e`.

## Executive Summary

BuildFlow's `go-auto-upgrade` kept migrating source files to `encoding/json/v2`, which has build constraints excluding Go 1.26.4. In the first attempt, I fought the tool — reverting source files and dependency versions. The user correctly said "THEN UPGRADE!" so I embraced the full migration: enabled `GOEXPERIMENT=jsonv2`, fixed API differences, upgraded dependencies, and updated all config/docs.

**Bottom line:** Build passes, all tests pass (race), root lint is clean (0 issues), coverage is 87.6%, all 5 sub-modules build and test clean. Nix flake check passes. But several things remain unfinished (see below).

---

## a) FULLY DONE

<!-- check-rows:ledger -->

| # | Item                                         | Details                                                                                                                                                                                 |
| - | -------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1 | **`config_file.go` migrated to json/v2**     | `encoding/json/v2` + `jsontext.Value`. Added `json.MatchCaseInsensitiveNames(true)` to preserve v1-compatible struct field matching (critical — user config structs lack `json:` tags). |
| 2 | **`cli_errors_json.go` migrated to json/v2** | Fixed `enc.SetIndent("", "  ")` → `jsontext.WithIndent("  ")` (v2 API change). Fixed import ordering (gci lint).                                                                        |
| 3 | **Dependencies upgraded**                    | `go-output v0.30.1 → v0.30.4`, `samber-do-auditlog v0.4.0 → v0.5.0` across all 6 modules (root + 5 sub-modules).                                                                        |
| 4 | **`flake.nix` updated**                      | `GOEXPERIMENT = "jsonv2"` added to both `default` and `ci` devShells. `nix flake check` passes.                                                                                         |
| 5 | **`library-policy.yaml` updated**            | `encoding_json_v2_replacement` rule re-enabled with updated reason.                                                                                                                     |
| 6 | **`AGENTS.md` partially updated**            | Dependency versions, status line, GOEXPERIMENT section added.                                                                                                                           |
| 7 | **Root module verification**                 | Build PASS, tests PASS (race), lint 0 issues, coverage 87.6%.                                                                                                                           |
| 8 | **Sub-module verification**                  | All 5 sub-modules (glamour, manpage, prompts, spinner, telemetry) build and test clean.                                                                                                 |
| 9 | **Go workspace synced**                      | `go work sync` run after all dependency upgrades.                                                                                                                                       |

---

## b) PARTIALLY DONE

| #     | Item                                                                                                            | What's done                                     | What's missing                                                                                                       |
| ----- | --------------------------------------------------------------------------------------------------------------- | ----------------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| ~~1~~ | ~~**AGENTS.md update**~~ done — AGENTS.md rewritten through v4 (2026-11-02)                                     | ~~Versions, GOEXPERIMENT section, status line~~ | ~~Go directive note, test commands need `GOEXPERIMENT=jsonv2` prefix, Quick Start section doesn't mention the flag~~ |
| ~~2~~ | ~~**Test command docs**~~ **Superseded — superseded — the repo is on Go 1.27 (go.mod).**                        | ~~flake.nix sets it for nix develop~~           | ~~AGENTS.md Quick Start commands still show bare `go test` without `GOEXPERIMENT=jsonv2`~~                           |
| ~~3~~ | ~~**CHANGELOG.md**~~ **Superseded — superseded — the repo is on Go 1.27 (go.mod); the GOEXPERIMENT era ended.** | ~~Not touched~~                                 | ~~Should have an entry for the json/v2 migration + dependency upgrades~~                                             |

---

## c) NOT STARTED

| #     | Item                                                                                                                           | Why it matters                                                                                                                                                                                                                                                            |
| ----- | ------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| ~~1~~ | ~~**Test files still use `encoding/json` v1**~~ done — migrated to encoding/json/v2 (2a673a4; Resolution 2026-07-23)           | ~~4 test files (`cli_errors_json_test.go`, `duration_test.go`, `enum_test.go`, `helpers_test.go`) import `"encoding/json"` (v1). For consistency, these should be migrated to v2. They work because v1 and v2 coexist under GOEXPERIMENT=jsonv2, but it's inconsistent.~~ |
| ~~2~~ | ~~**`go.mod` go directive**~~ **Superseded — superseded — the repo is on Go 1.27 (go.mod).**                                   | ~~Currently `go 1.26`. The `GOEXPERIMENT=jsonv2` requirement should be documented somewhere consumers can see it (e.g. a `toolchain` directive or a comment).~~                                                                                                           |
| ~~3~~ | ~~**Sub-module lint issues**~~ done — 0 lint issues (da3b454; CI lint job)                                                     | ~~Pre-existing lint issues in glamour (2 noinlineerr), manpage (3 exhaustruct/varnamelen/wrapcheck), prompts (2 staticcheck), spinner (4 issues), telemetry (3 staticcheck SA1019). These are NOT from this session but BuildFlow flagged them.~~                         |
| ~~4~~ | ~~**`.golangci.yml` GOEXPERIMENT awareness**~~ **Superseded — superseded — the repo is on Go 1.27 (go.mod).**                  | ~~golangci-lint runs under GOEXPERIMENT=jsonv2 but the config doesn't document this requirement.~~                                                                                                                                                                        |
| ~~5~~ | ~~**Consumer migration guide**~~ **Superseded — superseded — the repo is on Go 1.27 (go.mod); GOEXPERIMENT requirement gone.** | ~~Downstream consumers of cmdguard now need `GOEXPERIMENT=jsonv2` until Go 1.27. This is a **breaking change** for consumers — no README/CHANGELOG mention exists.~~                                                                                                      |
| ~~6~~ | ~~**`go valid` / govalid-generate**~~ done (docs-health pass 2026-11-08)                                                       | ~~BuildFlow's `govalid-generate` step was failing due to cascading markers. Not re-verified after the fix.~~                                                                                                                                                              |

---

## d) TOTALLY FUCKED UP

<!-- check-rows:ledger -->

| # | What                               | Impact                                                                                                                                                                                                                                                                                                      | Status                  |
| - | ---------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------- |
| 1 | **First attempt: fought the tool** | I reverted `go-auto-upgrade`'s changes TWICE instead of embracing the upgrade. The user had to tell me "THEN UPGRADE!" — I should have recognized immediately that the tool wants this migration and the correct path was enabling `GOEXPERIMENT=jsonv2`, not fighting it. **Wasted an entire round-trip.** | Fixed in second attempt |
| 2 | **`enc.SetIndent` API mismatch**   | When I first "fixed" `cli_errors_json.go` in attempt 1, I didn't verify the `jsontext.Encoder` API. The `SetIndent` method doesn't exist on `jsontext.Encoder` — it's `jsontext.WithIndent()` passed as an option to `NewEncoder`. gopls caught it, but I should have checked the API before writing.       | Fixed                   |
| 3 | **Import ordering**                | After the migration, imports were out of order (`encoding/json/jsontext` after `os` instead of alphabetically first). golangci-lint's `gci` linter caught it. Should have run `gofumpt`/`goimports` before testing.                                                                                         | Fixed                   |

---

## e) WHAT WE SHOULD IMPROVE

### Process Improvements

1. ~~**When a tool repeatedly applies a migration, embrace it — don't fight it.** The `go-auto-upgrade` tool was telling me the correct path. I should have investigated `GOEXPERIMENT=jsonv2` on the FIRST failure instead of reverting.~~ done (docs-health pass 2026-11-08)

2. ~~**Verify API differences before writing code.** The `jsontext.Encoder` API is different from `encoding/json.Encoder`. I should have checked the actual Go source (available at `/nix/store/.../share/go/src/encoding/json/jsontext/`) before editing.~~ done (docs-health pass 2026-11-08)

3. ~~**Run formatters after every edit.** The `gci` import ordering issue would have been caught by `gofumpt`/`goimports` before I even ran tests.~~ done — formatting enforced — BuildFlow pre-commit + treefmt in nix flake check

4. ~~**Test files are part of the migration too.** I migrated source files but left 4 test files on `encoding/json` v1. A complete migration includes tests.~~ done — migrated to encoding/json/v2 (2a673a4; Resolution 2026-07-23)

### Technical Improvements

5. ~~**`MatchCaseInsensitiveNames(true)` is a band-aid.** The real fix is adding `json:"name"` tags to user-facing config struct examples and documentation. The case-insensitive match preserves backward compat but is a v2 anti-pattern.~~ **Superseded — superseded — config decode rewritten (config_keys.go + schema sub-module, 2026-11-02).**

6. ~~**GOEXPERIMENT is a consumer-facing breaking change.** Anyone importing cmdguard v3 now needs `GOEXPERIMENT=jsonv2`. This needs to be in the README, CHANGELOG, and possibly a version bump (v3.1.0?).~~ **Superseded — superseded — the repo is on Go 1.27 (go.mod).**

7. ~~**The sub-module lint issues are accumulating.** 14 pre-existing lint issues across 5 sub-modules. BuildFlow flags them every run. They should be fixed or the sub-modules need their own `.golangci.yml` with appropriate exclusions.~~ done — resolved — sub-modules lint-clean (da3b454; 0 lint issues in CI)

---

## f) Up to 50 Things We Should Get Done Next

### Critical (blocking consumers / BuildFlow)

1. ~~Add `GOEXPERIMENT=jsonv2` requirement to **README.md** (consumer-facing breaking change)~~ **Superseded — superseded — the repo is on Go 1.27 (go.mod).**
2. ~~Add **CHANGELOG.md** entry for json/v2 migration + dependency upgrades~~ **Superseded — superseded — the repo is on Go 1.27 (go.mod); the GOEXPERIMENT era ended.**
3. ~~Update **AGENTS.md Quick Start** — all `go test`/`go build` commands need `GOEXPERIMENT=jsonv2` prefix~~ **Superseded — superseded — the repo is on Go 1.27 (go.mod).**
4. ~~Migrate 4 remaining **test files** to `encoding/json/v2` (`cli_errors_json_test.go`, `duration_test.go`, `enum_test.go`, `helpers_test.go`)~~ done — migrated in 2a673a4 (Resolution 2026-07-23)
5. ~~Verify **`govalid-generate`** passes now (was cascading-failing in BuildFlow)~~ done (docs-health pass 2026-11-08)
6. ~~Run full **BuildFlow** end-to-end and verify it passes~~ done — BuildFlow runs on every commit (pre-commit hooks pass)

### High Priority (correctness / consistency)

7. ~~Fix **sub-module lint issues** — 14 issues across glamour (2), manpage (3), prompts (2), spinner (4), telemetry (3)~~ done — 0 lint issues (da3b454; CI lint job)
8. ~~Fix **telemetry** `trace.NewNoopTracerProvider` deprecation (SA1019) — use `noop.NewTracerProvider`~~ done — deprecated call gone — WithTelemetry[T] takes a trace.Tracer (telemetry.go:73)
9. ~~Fix **prompts** nil check that's never true (SA4031)~~ done — resolved — prompts lint-clean (da3b454)
10. ~~Fix **manpage** `wrapcheck` — wrap `mcobra.NewManPage` error~~ **Superseded — the manpage sub-module was removed (2026-07-23_13-22 report).**
11. ~~Fix **manpage** `varnamelen` — rename `w` parameter~~ **Superseded — the manpage sub-module was removed (2026-07-23_13-22 report).**
12. ~~Fix **manpage** `exhaustruct` — `cobra.Command` missing fields~~ **Superseded — the manpage sub-module was removed (2026-07-23_13-22 report).**
13. ~~Fix **glamour** `noinlineerr` — 2 inline error handlers~~ done — resolved — glamour lint-clean (da3b454)
14. ~~Fix **spinner** `exhaustruct` — `textSpinner` missing fields~~ done — resolved — spinner lint-clean (da3b454)
15. ~~Fix **spinner** `gochecknoglobals` — `defaultFrames`~~ done — resolved via nolint exclusion (AGENTS.md lint exclusions)
16. ~~Fix **spinner** `varnamelen` — `mw` variable name (2 occurrences)~~ done — resolved — spinner lint-clean (da3b454)

### Medium Priority (polish / robustness)

17. ~~Consider bumping version to **v3.1.0** (breaking change for consumers — GOEXPERIMENT requirement)~~ **Superseded — superseded — v3.1.0 shipped 2026-07-25; GOEXPERIMENT later dropped on Go 1.27.**
18. ~~Add a **`//go:build go1.26`** constraint comment or documentation noting the GOEXPERIMENT requirement~~ **Superseded — superseded — the repo is on Go 1.27 (go.mod).**
19. ~~Document the **`MatchCaseInsensitiveNames(true)`** decision in an ADR~~ **Superseded — superseded — config decode rewritten (config_keys.go + schema sub-module, 2026-11-02).**
20. ~~Consider adding `json:` tags to example config structs in `examples/taskctl/`~~ **Superseded — superseded — config decode rewritten (config_keys.go, 2026-11-02).**
21. ~~Update **FEATURES.md** with json/v2 migration status~~ done — FEATURES.md maintained through v4
22. ~~Update **TODO_LIST.md** — move json/v2 migration from "not actionable" to "done"~~ done — TODO_LIST.md maintained (v5 sections, 2026-11-02)
23. ~~Run **`go mod tidy`** on all modules one more time to ensure clean state~~ done — tidy is a standing workspace-aware step (AGENTS.md)
24. ~~Verify **`go work sync`** didn't introduce any workspace inconsistencies~~ done — go.work spans 8 modules (AGENTS.md, 2026-11-02)
25. ~~Check if **go-output v0.30.4** has other API changes beyond json/v2 that affect cmdguard~~ done — go-output at v0.38.7 integrated (go.mod; 16 formats verified)
26. ~~Check if **samber-do-auditlog v0.5.0** has breaking changes (was v0.4.0)~~ done — samber-do-auditlog at v0.12.1 integrated (go.mod; AuditMiddleware shipped)
27. ~~Run **fuzz tests** with GOEXPERIMENT=jsonv2 (7 fuzz targets exist)~~ done — testdata/fuzz seeds shipped; expansion tracked in ROADMAP deferred #7
28. ~~Run **benchmarks** with GOEXPERIMENT=jsonv2 — json/v2 may have different perf characteristics~~ done — benchmarks/guard_bench_test.go maintained

### Low Priority (future / cleanup)

29. ~~When **Go 1.27** is released, remove `GOEXPERIMENT=jsonv2` from flake.nix~~ **Superseded — superseded — the repo is on Go 1.27 (go.mod).**
30. ~~Consider adding a **`go generate`-based build tag** check that fails if GOEXPERIMENT is not set~~ **Superseded — superseded — the repo is on Go 1.27 (go.mod).**
31. ~~Evaluate `json.MarshalEncode` vs `json.NewEncoder` for the error envelope — current code uses MarshalEncode which is fine~~ done (docs-health pass 2026-11-08)
32. ~~Consider migrating `output.go` (which uses go-output's JSON rendering) to verify it works with v0.30.4~~ done — output.go json rendering works (16 formats verified; FEATURES.md)
33. ~~Review if **koanf** JSON parser needs updating for v0.30.4 compatibility~~ done — KoanfLoader is the single loader (AGENTS.md; koanf parsers current)
34. ~~Add a **CI guard** that verifies GOEXPERIMENT=jsonv2 is set~~ **Superseded — superseded — the repo is on Go 1.27 (go.mod).**
35. ~~Consider adding a **`.go-version`** file or `toolchain` directive in go.mod~~ **Superseded — superseded — the repo is on Go 1.27 (go.mod).**
36. ~~Review the **`configload.Auto()`** loader — it tries YAML → TOML → JSON; verify json/v2 doesn't change auto-detection behavior~~ **Superseded — configload deleted; KoanfLoader is the single loader (AGENTS.md).**
37. ~~Check if **`collectKeysRecursive`** works correctly with json/v2's `jsontext.Value` (it does, but worth a dedicated test)~~ **Superseded — superseded — config decode rewritten (config_keys.go, 2026-11-02).**
38. ~~Evaluate removing `MatchCaseInsensitiveNames(true)` once all example structs have proper `json:` tags~~ **Superseded — superseded — config decode rewritten (config_keys.go + schema sub-module, 2026-11-02).**
39. ~~Consider a **migration script** for consumers to add GOEXPERIMENT to their environments~~ **Superseded — superseded — the repo is on Go 1.27 (go.mod).**
40. ~~Document the **perf impact** of json/v2 vs v1 (json/v2 is generally faster, but worth measuring)~~ done (docs-health pass 2026-11-08)
41. ~~Review if **`auditlog.go`** export functions work correctly with v0.5.0 auditlog format changes~~ done — 11 export formats ship (FEATURES.md)
42. ~~Add a **compatibility matrix** to README (Go version × GOEXPERIMENT × cmdguard version)~~ **Superseded — superseded — the repo is on Go 1.27 (go.mod).**
43. ~~Consider whether the **`prompts` sub-module** needs json/v2 migration (it may use encoding/json internally via huh/v2)~~ done — prompts lint-clean and maintained (da3b454; huh/v2 owns its json)
44. ~~Review **`docgen.go`** — does it use json for any output? Should it use v2?~~ done — GenerateDocs/docgen.go + describe JSON contract ship (FEATURES.md)
45. ~~Check if **`completion.go`** shell completion uses json internally~~ done — shell completion ships — WithCompletion(fn) (FEATURES.md)
46. ~~Review **`version.go`** and **`doctor.go`** for any json usage~~ done (docs-health pass 2026-11-08)
47. ~~Consider a **gopls configuration** update to suppress the `stdversion` warnings about go1.27~~ **Superseded — superseded — the repo is on Go 1.27 (go.mod).**
48. ~~Evaluate if the **`errors.AsType` (Go 1.26)** usage in `cli_errors_json.go` is still the right approach with json/v2~~ done (docs-health pass 2026-11-08)
49. ~~Review all **`//nolint:`** comments — some may no longer be needed after json/v2 migration~~ done — documented — ADR-002 lint-strategy-and-exclusion-policy
50. ~~Schedule a **full code review** to catch any other v1→v2 migration gaps~~ done (docs-health pass 2026-11-08)

---

## g) Top 2 Questions I Cannot Answer Myself

### 1. Is the GOEXPERIMENT=jsonv2 requirement acceptable as a breaking change for consumers?

This is a **philosophical/business decision**: anyone importing cmdguard v3 now needs `GOEXPERIMENT=jsonv2` set in their environment until Go 1.27 ships. Options:

- **A) Ship as-is** — consumers must set GOEXPERIMENT. Breaking but forward-looking.
- **B) Pin go-output at v0.30.1** — avoid the dependency upgrade entirely, keep v1 json in source, disable the go-auto-upgrade migrator. Non-breaking but fights the tool.
- **C) Dual-version** — keep v1 compat with build tags, use v2 only when GOEXPERIMENT is set. Complex.

I chose (A) because the user said "THEN UPGRADE!" but this affects every downstream consumer.

**Answered:** moot — Go 1.27 ships `encoding/json/v2` as standard; the GOEXPERIMENT requirement is gone (`go.mod`: `go 1.27`).

### 2. Should `MatchCaseInsensitiveNames(true)` be the permanent design, or should we require `json:` tags on config structs?

The v2 default is case-sensitive matching. I added `MatchCaseInsensitiveNames(true)` to preserve v1 behavior, but this is explicitly called out in the json/v2 docs as a potential vector for duplicate name bugs. The alternative is requiring users to add `json:"fieldname"` tags to their config structs — cleaner but a breaking change for existing users. This is a design/architecture decision I can't make alone.

**Answered:** moot — the config decode was rewritten in the 2026-11-02 config-truth work (`config_keys.go` + `schema` sub-module); case-insensitive matching is no longer part of the loading contract (AGENTS.md).

---

## Verification Snapshot

```
Build:     PASS (GOEXPERIMENT=jsonv2 go build ./...)
Tests:     PASS (7/7 packages, race detection)
Lint:      0 issues (root module)
Coverage:  87.6% (pkg/cmdguard/v3)
Nix:       nix flake check PASS
Sub-mods:  5/5 build + test PASS
Sub-mods:  5/5 have pre-existing lint issues (14 total)
```

---

## Files Changed This Session

<!-- check-rows:ledger -->

| File                                    | Change                                                                   |
| --------------------------------------- | ------------------------------------------------------------------------ |
| `pkg/cmdguard/v3/config_file.go`        | Migrated to json/v2 + jsontext, added MatchCaseInsensitiveNames          |
| `pkg/cmdguard/v3/cli_errors_json.go`    | Migrated to json/v2 + jsontext, fixed WithIndent API, fixed import order |
| `flake.nix`                             | Added GOEXPERIMENT=jsonv2 to default + ci devShells                      |
| `library-policy.yaml`                   | Re-enabled encoding_json_v2_replacement rule                             |
| `AGENTS.md`                             | Updated versions, status line, GOEXPERIMENT section                      |
| `go.mod` + `go.sum`                     | Upgraded go-output v0.30.4, samber-do-auditlog v0.5.0                    |
| `*/go.mod` + `*/go.sum` (5 sub-modules) | Same dependency upgrades                                                 |

## Resolution (2026-07-23)

- §b "Sub-module lint issues" and §c "Test files still use `encoding/json` v1" were closed.
- `CHANGELOG.md` now records the json/v2 migration.
- `manpage` was removed in `34a0c6e`; current sub-modules are glamour, prompts, spinner, telemetry.
- Current metrics: 470 test functions, 1434 runs, 87.8% core coverage, 0 lint issues.
