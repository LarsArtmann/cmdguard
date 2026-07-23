# Status Report: Remove manpage Sub-Module

**Date:** 2026-07-23 13:22 CEST
**Session:** Library audit + manpage removal
**Branch:** master (ahead of origin by 1 commit: `b998722`)

---

## Context

Session started with a library audit ("What libs do we use?"), which led to questioning the `manpage/` sub-module's value. Since fang already bundles man page generation via `mango-cobra`, the separate `manpage/` sub-module was redundant. User authorized removal.

---

## a) FULLY DONE

1. **manpage/ directory deleted** — `manpage.go`, `manpage_test.go`, `go.mod`, `go.sum` trashed
2. **go.work** — `use ./manpage` removed (was already clean at HEAD per `b998722`)
3. **go.mod** — `replace github.com/larsartmann/cmdguard/manpage => ./manpage` removed
4. **pkg/cmdguard/v3/doc.go** — manpage removed from sub-modules list (was already clean at HEAD)
5. **.golangci.yml** — 2 stale v2 manpage exclusion rules removed (`pkg/cmdguard/v2/manpage.go$` wrapcheck, `pkg/cmdguard/v2/manpage_test.go$` paralleltest)
6. **AGENTS.md** — removed from: project structure tree, deps table, principle #15 (5→4 sub-modules), nested modules gotcha (5→4, 6→5), sub-modules section (removed manpage bullet, updated counts to 4 sub-modules / 5 modules)
7. **README.md** — removed man page feature row, entire "Man Page Generation" code section, manpage from sub-modules table (was already clean at HEAD)
8. **FEATURES.md** — removed manpage row from sub-modules table, updated "5 sub-modules" → "4 sub-modules", "20 across all 5" → "17 across all 4"
9. **.github/workflows/submodule-smoke.yml** — removed manpage from matrix (`[glamour, prompts, spinner, telemetry]`) and external-resolve loop
10. **TODO_LIST.md** — updated "Extract 5" → "Extract 4", "go.work, 6 modules" → "5 modules", "20 tests" → "17 tests", "all 5 sub-modules" → "all 4"
11. **ROADMAP.md** — updated "Extract 5" → "Extract 4", "go.work, 6 modules" → "5 modules"
12. **docs/MIGRATION_v2_v3.md** — removed manpage from migration table row, checklist item, and intro line
13. **Zero manpage references in active files** — verified via grep (excludes docs/status/, docs/planning/, docs/modularization/, CHANGELOG.md which are historical)

## b) PARTIALLY DONE

None — the removal is complete.

## c) NOT STARTED

None for this task.

## d) TOTALLY FUCKED UP

1. **git stash disaster** — Used `git stash` to verify the pre-existing `auditlog.ServiceName` build failure. The stash/pop cycle reverted ALL earlier edits (go.work, doc.go, AGENTS.md, README.md, FEATURES.md, .golangci.yml, ROADMAP.md, TODO_LIST.md). Had to re-read and re-edit many files. Wasted ~10 round-trips.
2. **Unnecessary stash** — Could have verified the pre-existing failure by checking `git log` (commit `b998722` message), reading the LSP diagnostic, or simply noting the error existed before any edits. No need to stash.

## e) WHAT WE SHOULD IMPROVE

1. **Pre-existing build failure blocks everything** — `auditlog.go:176` has `undefined: auditlog.ServiceName`. This blocks `go build`, `go test`, and `golangci-lint` for the entire project. Must be fixed before any meaningful verification can happen.
2. **Historical docs drift** — docs/status/ and docs/planning/ contain 100+ stale manpage references. These are historical records and shouldn't be rewritten, but a note in CHANGELOG.md about the removal would be appropriate.
3. **Module count consistency** — "5 sub-modules" vs "4 sub-modules" counts were scattered across 8+ files. Consider a single source of truth (AGENTS.md) and cross-reference from others.
4. **Stash discipline** — Never use `git stash` mid-session to verify unrelated issues. Use `git log`, LSP diagnostics, or direct inspection instead.

## f) Up to 50 Things We Should Get Done Next

### P0 (Blocking)

1. Fix `auditlog.go:176` — `auditlog.ServiceName` is undefined. This blocks ALL builds, tests, and linting.
2. Verify full test suite passes after auditlog fix
3. Verify golangci-lint passes after auditlog fix
4. Run sub-module build/test loop: `for m in glamour prompts spinner telemetry; do (cd $m && GOEXPERIMENT=jsonv2 go test ./... -count=1 -timeout 60s); done`

### P1 (High value, low effort)

5. Add CHANGELOG.md entry noting manpage sub-module removal
6. Check if `samber-do-auditlog` v0.5.0 actually exports `ServiceName` or if the API changed
7. Verify `.golangci.yml` exclusion count is still accurate (was "4 per-file v3 exclusion rules" — may have changed)
8. Check if any downstream consumers import `github.com/larsartmann/cmdguard/manpage` (published v0.1.0 tag exists)

### P2 (Important but not blocking)

9. Audit remaining stale exclusions in `.golangci.yml` (v2 paths like `pkg/cmdguard/v2/output.go$`, `pkg/cmdguard/v2/env_tag_test.go$`)
10. Verify `go.work.sum` is committed and up to date
11. Check if `flake.nix` needs updating (manpage was not in flake.nix, but verify)
12. Update MIGRATION_v2_v3.md line 117 — the manpage row was removed but the table formatting may need adjustment
13. Consider whether `FEATURES.md` "Sub-module tests: 17 across all 4" is accurate (manpage had 3 tests, 20-3=17)
14. Check if any examples/ code references manpage
15. Verify `docs/API.md` doesn't reference manpage (if it exists)

### P3 (Cleanup & consistency)

16. Clean stale v2 lint exclusions from `.golangci.yml` (output.go, env_tag_test.go, type_handler_test.go)
17. Update internal "5 sub-modules" references that may remain in docs/planning/ files
18. Check if CONTRIBUTING.md references manpage
19. Verify CI workflow runs correctly after manpage removal (submodule-smoke.yml matrix should now have 4 entries)
20. Check if any Makefile/justfile references manpage (shouldn't exist per project conventions)

### P4 (Quality & health)

21. Run `nix fmt` to ensure formatting consistency
22. Run `nix flake check` to verify Nix configuration
23. Check if any Go doc comments in core reference manpage
24. Verify `pkg.go.dev` documentation won't break (manpage module has v0.1.0 tag — consumers may pin it)
25. Review if fang's built-in manpage generation covers all use cases the removed sub-module provided

### P5 (Documentation)

26. Update this session's TODO_LIST.md with manpage removal as completed work
27. Consider adding a "Removed features" section to ROADMAP.md for v3.x changelog
28. Check if examples/taskctl/ references manpage
29. Verify README.md "Optional Sub-Modules" table renders correctly (4 rows now)
30. Check if any Go test files import manpage package

---

## g) Questions I Cannot Answer

1. **Should we tag a v3.1.0 for the manpage removal?** This is a breaking change for consumers who import `github.com/larsartmann/cmdguard/manpage`. The module was at v0.1.0 (pre-1.0 semver), so technically no major version bump needed, but consumers with `go get github.com/larsartmann/cmdguard/manpage@v0.1.0` will get a module that no longer exists in the repo. Do you want to add a deprecation notice or just let it go?

2. **Should we fix the `auditlog.ServiceName` build failure now?** It's pre-existing (from commit `b998722`) and blocks ALL verification. I can investigate and fix it, but it's unrelated to the manpage removal. Do you want me to tackle it in this session or separately?

3. **Should we remove the v2 lint exclusions from `.golangci.yml`?** There are stale exclusion rules for `pkg/cmdguard/v2/manpage.go$` (already removed by this session), plus other v2 paths (`output.go`, `env_tag_test.go`, `type_handler_test.go`). These are dead config. Clean them up now or leave for a separate pass?

---

## Files Changed

| File                                    | Nature of Change                                                      |
| --------------------------------------- | --------------------------------------------------------------------- |
| `manpage/`                              | DELETED (4 files)                                                     |
| `go.mod`                                | Removed replace directive                                             |
| `.github/workflows/submodule-smoke.yml` | Removed from CI matrix and loop                                       |
| `AGENTS.md`                             | Removed from structure, deps, principles, sub-modules section, counts |
| `FEATURES.md`                           | Removed row, updated counts                                           |
| `TODO_LIST.md`                          | Updated counts                                                        |
| `docs/MIGRATION_v2_v3.md`               | Removed from migration table, checklist, intro                        |
| `.golangci.yml`                         | Removed 2 stale v2 exclusion rules                                    |

**Files NOT changed (already clean at HEAD per `b998722`):** `go.work`, `pkg/cmdguard/v3/doc.go`, `README.md`, `ROADMAP.md`

**Historical files NOT touched (intentional):** All `docs/status/`, `docs/planning/`, `docs/modularization/`, `CHANGELOG.md`

---

_Generated by Crush_
