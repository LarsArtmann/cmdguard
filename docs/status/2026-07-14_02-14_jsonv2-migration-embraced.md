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

| #   | Item                                         | Details                                                                                                                                                                                 |
| --- | -------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | **`config_file.go` migrated to json/v2**     | `encoding/json/v2` + `jsontext.Value`. Added `json.MatchCaseInsensitiveNames(true)` to preserve v1-compatible struct field matching (critical — user config structs lack `json:` tags). |
| 2   | **`cli_errors_json.go` migrated to json/v2** | Fixed `enc.SetIndent("", "  ")` → `jsontext.WithIndent("  ")` (v2 API change). Fixed import ordering (gci lint).                                                                        |
| 3   | **Dependencies upgraded**                    | `go-output v0.30.1 → v0.30.4`, `samber-do-auditlog v0.4.0 → v0.5.0` across all 6 modules (root + 5 sub-modules).                                                                        |
| 4   | **`flake.nix` updated**                      | `GOEXPERIMENT = "jsonv2"` added to both `default` and `ci` devShells. `nix flake check` passes.                                                                                         |
| 5   | **`library-policy.yaml` updated**            | `encoding_json_v2_replacement` rule re-enabled with updated reason.                                                                                                                     |
| 6   | **`AGENTS.md` partially updated**            | Dependency versions, status line, GOEXPERIMENT section added.                                                                                                                           |
| 7   | **Root module verification**                 | Build PASS, tests PASS (race), lint 0 issues, coverage 87.6%.                                                                                                                           |
| 8   | **Sub-module verification**                  | All 5 sub-modules (glamour, manpage, prompts, spinner, telemetry) build and test clean.                                                                                                 |
| 9   | **Go workspace synced**                      | `go work sync` run after all dependency upgrades.                                                                                                                                       |

---

## b) PARTIALLY DONE

| #   | Item                  | What's done                                 | What's missing                                                                                                   |
| --- | --------------------- | ------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| 1   | **AGENTS.md update**  | Versions, GOEXPERIMENT section, status line | Go directive note, test commands need `GOEXPERIMENT=jsonv2` prefix, Quick Start section doesn't mention the flag |
| 2   | **Test command docs** | flake.nix sets it for nix develop           | AGENTS.md Quick Start commands still show bare `go test` without `GOEXPERIMENT=jsonv2`                           |
| 3   | **CHANGELOG.md**      | Not touched                                 | Should have an entry for the json/v2 migration + dependency upgrades                                             |

---

## c) NOT STARTED

| #   | Item                                        | Why it matters                                                                                                                                                                                                                                                        |
| --- | ------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | **Test files still use `encoding/json` v1** | 4 test files (`cli_errors_json_test.go`, `duration_test.go`, `enum_test.go`, `helpers_test.go`) import `"encoding/json"` (v1). For consistency, these should be migrated to v2. They work because v1 and v2 coexist under GOEXPERIMENT=jsonv2, but it's inconsistent. |
| 2   | **`go.mod` go directive**                   | Currently `go 1.26`. The `GOEXPERIMENT=jsonv2` requirement should be documented somewhere consumers can see it (e.g. a `toolchain` directive or a comment).                                                                                                           |
| 3   | **Sub-module lint issues**                  | Pre-existing lint issues in glamour (2 noinlineerr), manpage (3 exhaustruct/varnamelen/wrapcheck), prompts (2 staticcheck), spinner (4 issues), telemetry (3 staticcheck SA1019). These are NOT from this session but BuildFlow flagged them.                         |
| 4   | **`.golangci.yml` GOEXPERIMENT awareness**  | golangci-lint runs under GOEXPERIMENT=jsonv2 but the config doesn't document this requirement.                                                                                                                                                                        |
| 5   | **Consumer migration guide**                | Downstream consumers of cmdguard now need `GOEXPERIMENT=jsonv2` until Go 1.27. This is a **breaking change** for consumers — no README/CHANGELOG mention exists.                                                                                                      |
| 6   | **`go valid` / govalid-generate**           | BuildFlow's `govalid-generate` step was failing due to cascading markers. Not re-verified after the fix.                                                                                                                                                              |

---

## d) TOTALLY FUCKED UP

| #   | What                               | Impact                                                                                                                                                                                                                                                                                                      | Status                  |
| --- | ---------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------- |
| 1   | **First attempt: fought the tool** | I reverted `go-auto-upgrade`'s changes TWICE instead of embracing the upgrade. The user had to tell me "THEN UPGRADE!" — I should have recognized immediately that the tool wants this migration and the correct path was enabling `GOEXPERIMENT=jsonv2`, not fighting it. **Wasted an entire round-trip.** | Fixed in second attempt |
| 2   | **`enc.SetIndent` API mismatch**   | When I first "fixed" `cli_errors_json.go` in attempt 1, I didn't verify the `jsontext.Encoder` API. The `SetIndent` method doesn't exist on `jsontext.Encoder` — it's `jsontext.WithIndent()` passed as an option to `NewEncoder`. gopls caught it, but I should have checked the API before writing.       | Fixed                   |
| 3   | **Import ordering**                | After the migration, imports were out of order (`encoding/json/jsontext` after `os` instead of alphabetically first). golangci-lint's `gci` linter caught it. Should have run `gofumpt`/`goimports` before testing.                                                                                         | Fixed                   |

---

## e) WHAT WE SHOULD IMPROVE

### Process Improvements

1. **When a tool repeatedly applies a migration, embrace it — don't fight it.** The `go-auto-upgrade` tool was telling me the correct path. I should have investigated `GOEXPERIMENT=jsonv2` on the FIRST failure instead of reverting.

2. **Verify API differences before writing code.** The `jsontext.Encoder` API is different from `encoding/json.Encoder`. I should have checked the actual Go source (available at `/nix/store/.../share/go/src/encoding/json/jsontext/`) before editing.

3. **Run formatters after every edit.** The `gci` import ordering issue would have been caught by `gofumpt`/`goimports` before I even ran tests.

4. **Test files are part of the migration too.** I migrated source files but left 4 test files on `encoding/json` v1. A complete migration includes tests.

### Technical Improvements

5. **`MatchCaseInsensitiveNames(true)` is a band-aid.** The real fix is adding `json:"name"` tags to user-facing config struct examples and documentation. The case-insensitive match preserves backward compat but is a v2 anti-pattern.

6. **GOEXPERIMENT is a consumer-facing breaking change.** Anyone importing cmdguard v3 now needs `GOEXPERIMENT=jsonv2`. This needs to be in the README, CHANGELOG, and possibly a version bump (v3.1.0?).

7. **The sub-module lint issues are accumulating.** 14 pre-existing lint issues across 5 sub-modules. BuildFlow flags them every run. They should be fixed or the sub-modules need their own `.golangci.yml` with appropriate exclusions.

---

## f) Up to 50 Things We Should Get Done Next

### Critical (blocking consumers / BuildFlow)

1. Add `GOEXPERIMENT=jsonv2` requirement to **README.md** (consumer-facing breaking change)
2. Add **CHANGELOG.md** entry for json/v2 migration + dependency upgrades
3. Update **AGENTS.md Quick Start** — all `go test`/`go build` commands need `GOEXPERIMENT=jsonv2` prefix
4. Migrate 4 remaining **test files** to `encoding/json/v2` (`cli_errors_json_test.go`, `duration_test.go`, `enum_test.go`, `helpers_test.go`)
5. Verify **`govalid-generate`** passes now (was cascading-failing in BuildFlow)
6. Run full **BuildFlow** end-to-end and verify it passes

### High Priority (correctness / consistency)

7. Fix **sub-module lint issues** — 14 issues across glamour (2), manpage (3), prompts (2), spinner (4), telemetry (3)
8. Fix **telemetry** `trace.NewNoopTracerProvider` deprecation (SA1019) — use `noop.NewTracerProvider`
9. Fix **prompts** nil check that's never true (SA4031)
10. Fix **manpage** `wrapcheck` — wrap `mcobra.NewManPage` error
11. Fix **manpage** `varnamelen` — rename `w` parameter
12. Fix **manpage** `exhaustruct` — `cobra.Command` missing fields
13. Fix **glamour** `noinlineerr` — 2 inline error handlers
14. Fix **spinner** `exhaustruct` — `textSpinner` missing fields
15. Fix **spinner** `gochecknoglobals` — `defaultFrames`
16. Fix **spinner** `varnamelen` — `mw` variable name (2 occurrences)

### Medium Priority (polish / robustness)

17. Consider bumping version to **v3.1.0** (breaking change for consumers — GOEXPERIMENT requirement)
18. Add a **`//go:build go1.26`** constraint comment or documentation noting the GOEXPERIMENT requirement
19. Document the **`MatchCaseInsensitiveNames(true)`** decision in an ADR
20. Consider adding `json:` tags to example config structs in `examples/taskctl/`
21. Update **FEATURES.md** with json/v2 migration status
22. Update **TODO_LIST.md** — move json/v2 migration from "not actionable" to "done"
23. Run **`go mod tidy`** on all modules one more time to ensure clean state
24. Verify **`go work sync`** didn't introduce any workspace inconsistencies
25. Check if **go-output v0.30.4** has other API changes beyond json/v2 that affect cmdguard
26. Check if **samber-do-auditlog v0.5.0** has breaking changes (was v0.4.0)
27. Run **fuzz tests** with GOEXPERIMENT=jsonv2 (7 fuzz targets exist)
28. Run **benchmarks** with GOEXPERIMENT=jsonv2 — json/v2 may have different perf characteristics

### Low Priority (future / cleanup)

29. When **Go 1.27** is released, remove `GOEXPERIMENT=jsonv2` from flake.nix
30. Consider adding a **`go generate`-based build tag** check that fails if GOEXPERIMENT is not set
31. Evaluate `json.MarshalEncode` vs `json.NewEncoder` for the error envelope — current code uses MarshalEncode which is fine
32. Consider migrating `output.go` (which uses go-output's JSON rendering) to verify it works with v0.30.4
33. Review if **koanf** JSON parser needs updating for v0.30.4 compatibility
34. Add a **CI guard** that verifies GOEXPERIMENT=jsonv2 is set
35. Consider adding a **`.go-version`** file or `toolchain` directive in go.mod
36. Review the **`configload.Auto()`** loader — it tries YAML → TOML → JSON; verify json/v2 doesn't change auto-detection behavior
37. Check if **`collectKeysRecursive`** works correctly with json/v2's `jsontext.Value` (it does, but worth a dedicated test)
38. Evaluate removing `MatchCaseInsensitiveNames(true)` once all example structs have proper `json:` tags
39. Consider a **migration script** for consumers to add GOEXPERIMENT to their environments
40. Document the **perf impact** of json/v2 vs v1 (json/v2 is generally faster, but worth measuring)
41. Review if **`auditlog.go`** export functions work correctly with v0.5.0 auditlog format changes
42. Add a **compatibility matrix** to README (Go version × GOEXPERIMENT × cmdguard version)
43. Consider whether the **`prompts` sub-module** needs json/v2 migration (it may use encoding/json internally via huh/v2)
44. Review **`docgen.go`** — does it use json for any output? Should it use v2?
45. Check if **`completion.go`** shell completion uses json internally
46. Review **`version.go`** and **`doctor.go`** for any json usage
47. Consider a **gopls configuration** update to suppress the `stdversion` warnings about go1.27
48. Evaluate if the **`errors.AsType` (Go 1.26)** usage in `cli_errors_json.go` is still the right approach with json/v2
49. Review all **`//nolint:`** comments — some may no longer be needed after json/v2 migration
50. Schedule a **full code review** to catch any other v1→v2 migration gaps

---

## g) Top 2 Questions I Cannot Answer Myself

### 1. Is the GOEXPERIMENT=jsonv2 requirement acceptable as a breaking change for consumers?

This is a **philosophical/business decision**: anyone importing cmdguard v3 now needs `GOEXPERIMENT=jsonv2` set in their environment until Go 1.27 ships. Options:

- **A) Ship as-is** — consumers must set GOEXPERIMENT. Breaking but forward-looking.
- **B) Pin go-output at v0.30.1** — avoid the dependency upgrade entirely, keep v1 json in source, disable the go-auto-upgrade migrator. Non-breaking but fights the tool.
- **C) Dual-version** — keep v1 compat with build tags, use v2 only when GOEXPERIMENT is set. Complex.

I chose (A) because the user said "THEN UPGRADE!" but this affects every downstream consumer.

### 2. Should `MatchCaseInsensitiveNames(true)` be the permanent design, or should we require `json:` tags on config structs?

The v2 default is case-sensitive matching. I added `MatchCaseInsensitiveNames(true)` to preserve v1 behavior, but this is explicitly called out in the json/v2 docs as a potential vector for duplicate name bugs. The alternative is requiring users to add `json:"fieldname"` tags to their config structs — cleaner but a breaking change for existing users. This is a design/architecture decision I can't make alone.

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
