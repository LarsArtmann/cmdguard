# Status Report: v3.0.0 Module Migration & Release Cleanup

**Date:** 2026-07-07 02:55
**Session Goal:** Fix the v2.11.0 mis-release — migrate breaking v3 redesign to proper /v3 module path
**Status:** CORE MIGRATION DONE — significant doc debt remains

---

> **Update 2026-07-23:** The mechanical migration was completed; the remaining doc debt listed in §b (README/AGENTS stale refs, GitHub Releases, sub-module tests) was closed in `98dd7d7` and later sessions. The workspace now has 4 optional sub-modules; `manpage` was removed in `34a0c6e`.

## Executive Summary

The v3 redesign (non-generic `CLIOption`/`CommandOption`, 5 extracted sub-modules) was mis-released as `v2.11.0` on a `/v2` module path — a Go semver violation. This session migrated the module to `github.com/larsartmann/cmdguard/v3`, created `v3.0.0` tag, deleted the bad `v2.11.0`, created a `release/v2.10` branch with a `retract` directive for the v2 line, and updated consumer-facing docs.

**Build passes. Tests pass (1831 runs, 0 failures). Race clean. Lint clean. v3.0.0 resolves on the Go module proxy.**

However, the migration was **mechanical** (sed-based find-replace) and several documents — particularly `AGENTS.md` — still describe the **old v2 architecture** with stale file listings, wrong dependency versions, and references to deleted features (`editor.go`, `result.go`, `WithGlamourHelp`, etc.).

---

## a) FULLY DONE

<!-- check-rows:ledger -->
| #  | What                                                                                             | Evidence                                               |
| -- | ------------------------------------------------------------------------------------------------ | ------------------------------------------------------ |
| 1  | Module path migrated `cmdguard/v2` → `cmdguard/v3`                                               | `go.mod` module line, all 731 internal refs, 159 files |
| 2  | Directory renamed `pkg/cmdguard/v2` → `pkg/cmdguard/v3`                                          | `git mv` — history preserved (98-99% similarity)       |
| 3  | Package declarations updated (`package v2` → `package v3`, `v2_test` → `v3_test`)                | All 140+ files                                         |
| 4  | Import aliases updated (`v2."..."` → `v3."..."`)                                                 | All 37 files using alias                               |
| 5  | `v2.X` qualified usages updated (`v2.NewCLI` → `v3.NewCLI`, etc.)                                | 731 occurrences                                        |
| 6  | All 5 sub-module `go.mod` files updated (module paths, version placeholders, replace directives) | glamour/manpage/prompts/spinner/telemetry              |
| 7  | `go.work` workspace file correct                                                                 | Uses all 6 modules                                     |
| 8  | `go mod tidy` run on all 6 modules                                                               | All exit 0                                             |
| 9  | `go build ./...` passes                                                                          | Exit 0                                                 |
| 10 | `go vet ./...` passes                                                                            | Exit 0                                                 |
| 11 | `go test ./... -count=1 -race` passes                                                            | 1831 runs, 0 failures                                  |
| 12 | `golangci-lint run ./...` passes                                                                 | Exit 0 (pre-existing warnings only)                    |
| 13 | Coverage: 87.3% core, 87.7% configload                                                           | Verified                                               |
| 14 | BuildFlow pre-commit hooks pass                                                                  | 26/26, 0 failures                                      |
| 15 | `v3.0.0` annotated tag created and pushed                                                        | `git ls-remote --tags origin` confirms                 |
| 16 | `v2.11.0` tag deleted locally and remotely                                                       | Both confirmed                                         |
| 17 | `release/v2.10` branch created at orphaned `v2.10.3` commit                                      | `git branch --contains v2.10.3` → `release/v2.10`      |
| 18 | `v2.10.4` tagged with `retract v2.11.0` directive                                                | Pushed to remote                                       |
| 19 | `v3.0.0` resolves on Go module proxy                                                             | `go list -m` confirms                                  |
| 20 | README.md fully updated (paths, code examples, API signatures, badge, features table)            | All v2 refs eliminated                                 |
| 21 | CHANGELOG.md updated (v3.0.0 section, retract warning, comparison links)                         | Links verified                                         |
| 22 | FEATURES.md updated (package path, coverage numbers)                                             | Partially — see section b                              |
| 23 | Consumer-facing docs updated (QUICKSTART, TUTORIAL, API, MIGRATION, ERROR_REFERENCE)             | All v2 path refs eliminated                            |
| 24 | Historical docs left as point-in-time snapshots (status reports, planning docs)                  | Intentional — 30+ files                                |

---

## b) PARTIALLY DONE

| # | What                                                                    | What's Missing                                                                                                                                                                                                                                                                                                                                                                                                                   |
| - | ----------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| ~~1~~ | ~~**AGENTS.md** — header, module path, status line, API reference updated~~ done — AGENTS.md fully rewritten through v4 (2026-11-02) | ~~Dependency table still shows `go-output v0.23.3` / `samber-do-auditlog v0.3.1` (actual: v0.30.1 / v0.4.0). File structure tree still lists `editor.go`, `result.go`, `telemetry.go`, `glamour.go`, `spinner.go`, `manpage.go` as core files — **all removed in v3**. Design principles list (`WithGlamourHelp[T]()`, `EditInEditor()`, `Result[T]` sum types, etc.) still describes deleted features. ~40% of the file is stale.~~ |
| ~~2~~ | ~~**FEATURES.md** — package path and coverage updated~~ done — FEATURES.md dep tables maintained through v4 releases | ~~Dependency table still shows `go-output v0.23.3` / `samber-do-auditlog v0.3.1`. References `huh/v2` and `glamour/v2` as core deps (now sub-module deps). Feature list still mentions `EditInEditor`, `Result[T]`, `WithGlamourHelp[T]`.~~ |
| ~~3~~ | ~~**CHANGELOG.md** — v3.0.0 section, retract note, link list updated~~ done — shipped — CHANGELOG [3.0.0] + [2.10.4] sections; before/after kept as migration contrast | ~~The "Before" code example still uses `v2.NewCommand[...]` syntax (intentional as migration contrast, but the alias should be contextualized). `[2.10.3]` link exists but no changelog entry for v2.10.3 or v2.10.4.~~ |
| ~~4~~ | ~~**Go module proxy** — `v3.0.0` resolves~~ done — GitHub Releases published (gh release list); retract shipped in v2.10.4 | ~~`v2.11.0` is **still cached on the proxy** and downloadable. The `retract` directive in `v2.10.4` will warn consumers, but won't prevent download. No GitHub Release created for v3.0.0 or v2.10.4.~~ |

---

## c) NOT STARTED

| # | What                                       | Why It Matters                                                                                          |
| - | ------------------------------------------ | ------------------------------------------------------------------------------------------------------- |
| ~~1~~ | ~~GitHub Releases for `v3.0.0` and `v2.10.4`~~ done — GitHub Releases published (gh release list) | ~~Tags exist but no release notes on GitHub — consumers browsing releases see nothing after v2.10.2~~ |
| ~~2~~ | ~~GitHub Release deletion/edit for `v2.11.0`~~ done — moot — no release ever existed (confirmed in-row) | ~~No GitHub release was ever created for it (confirmed), but worth double-checking the releases page~~ |
| ~~3~~ | ~~`pkg.go.dev` refresh for v3 module~~ done — indexed — README badge links pkg.go.dev/v4 | ~~Badge URL updated in README but pkg.go.dev hasn't indexed `/v3` yet — first `go get` triggers it~~ |
| ~~4~~ | ~~Migration guide for v2→v3 consumers~~ done — docs/MIGRATION_v2_v3.md exists | ~~No `docs/MIGRATION_v2_v3.md` exists — users upgrading from v2 to v3 have no upgrade path documented~~ |
| ~~5~~ | ~~Version bump in any Nix flake or CI config~~ done — flake.nix carries no version strings; CI clean | ~~`.github/` has no v2 refs (confirmed clean), but `flake.nix` wasn't checked for version strings~~ |
| ~~6~~ | ~~`TODO_LIST.md` sync~~ done — TODO_LIST.md maintained (v5 sections, 2026-11-02) | ~~Still references v2 work items~~ |
| ~~7~~ | ~~`ROADMAP.md` content audit~~ done — ROADMAP.md reflects v4/v5 direction (2026-11-02) | ~~Path refs updated but content may describe v2-era plans~~ |
| ~~8~~ | ~~Consumer smoke test~~ done — closed per Resolution (2026-07-23) | ~~Never verified `go get github.com/larsartmann/cmdguard/v3@v3.0.0` from a fresh module outside this repo~~ |

---

## d) TOTALLY FUCKED UP

<!-- check-rows:ledger -->
| # | What                                                                                                                                                                                                                                                                                                                                  | Impact                                                                                                                                                                      | Severity                                                         |
| - | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------- |
| 1 | **AGENTS.md is a lie** — I updated the header and deps table partially but left **6 deleted files** in the structure tree (`editor.go`, `result.go`, `telemetry.go`, `glamour.go`, `spinner.go`, `manpage.go`) and **6+ deleted features** in the Design Principles section (`EditInEditor`, `Result[T]`, `WithGlamourHelp[T]`, etc.) | Any AI session or contributor reading AGENTS.md will look for files and features that **don't exist**. This is the #1 context file for the project and I left it 40% stale. | **HIGH** — this is the file I should have been most careful with |
| 2 | **FEATURES.md deps table** still says `go-output v0.23.3` (actual v0.30.1) and `samber-do-auditlog v0.3.1` (actual v0.4.0)                                                                                                                                                                                                            | Consumers reading FEATURES.md see wrong dependency versions                                                                                                                 | MEDIUM                                                           |
| 3 | **AGENTS.md deps table** has the same wrong versions                                                                                                                                                                                                                                                                                  | Same problem in the #1 context file                                                                                                                                         | **HIGH**                                                         |
| 4 | **CHANGELOG has no v2.10.3/v2.10.4 entries** — the tags exist but no changelog content describes what they contain                                                                                                                                                                                                                    | v2 consumers browsing the changelog see a gap between v2.10.2 and v3.0.0                                                                                                    | MEDIUM                                                           |
| 5 | **`/v3` in `.gitignore`** — BuildFlow added this automatically; it ignores a root-level `/v3` directory. Harmless in practice (no such directory exists), but it's **wrong intent** — it was triggered by the module path change and should be removed or explained                                                                   | Minor confusion for anyone reading `.gitignore`                                                                                                                             | LOW                                                              |
| 6 | **v2.11.0 is still downloadable from the Go proxy** — tag deletion doesn't purge the proxy cache. The `retract` directive helps but is reactive (requires consumers to see v2.10.4 first)                                                                                                                                             | A consumer who pins `@v2.11.0` directly still gets the broken release                                                                                                       | MEDIUM — `GOPROXY` purge or waiting for TTL is the only real fix |

---

## e) WHAT WE SHOULD IMPROVE

| #  | Area                                                                                                                              | Current State                                                                                                                                                 | Should Be                                                                                                                                        |
| -- | --------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------ |
| 1  | **Mechanical sed migration without verifying semantic correctness** **→ ROADMAP raw ideas (post-rewrite semantic doc audit)** | I ran sed on 731 occurrences and checked for build/lint/test failures — but never verified that the resulting **docs** accurately describe the v3 API surface | After any bulk rename, do a **semantic audit**: read every doc file and verify it describes reality, not just that string substitutions happened |
| 2  | **AGENTS.md as afterthought** **→ ROADMAP raw ideas (post-rewrite semantic doc audit)** | I updated it with targeted edits (header, deps, paths) but didn't re-read the whole file to find stale sections                                               | AGENTS.md is the **most critical context file** — it should be read end-to-end after any architectural change, not spot-edited                   |
| 3  | **Deleted features not cleaned from docs** **→ ROADMAP raw ideas (post-rewrite semantic doc audit)** | The v3 redesign removed `editor.go`, `result.go`, `telemetry.go`, `glamour.go`, `spinner.go`, `manpage.go` from core — but docs still list them               | When files are deleted, grep for their names in ALL docs and update/remove references                                                            |
| ~~4~~  | ~~**No GitHub Releases**~~ done — GitHub Releases published (gh release list) | ~~Git tags are pushed but GitHub Releases page is empty~~ | ~~Create proper release notes — that's what consumers browse~~ |
| ~~5~~  | ~~**No migration guide**~~ done — docs/MIGRATION_v2_v3.md + MIGRATION_v3_v4.md exist | ~~v1→v2 migration guide exists (`docs/MIGRATION_v1_v2.md`) but no v2→v3~~ | ~~Breaking changes need migration docs — especially when the module path changes~~ |
| 6  | **Dependency version drift in docs** **→ ROADMAP raw ideas (derive doc dep tables from go.mod)** | AGENTS.md and FEATURES.md show `go-output v0.23.3` and `samber-do-auditlog v0.3.1` while the actual deps are v0.30.1/v0.4.0                                   | Dependency tables in docs should be derived from `go.mod`, not hand-maintained                                                                   |
| ~~7~~  | ~~**CHANGELOG [3.0.0] comparison link** points to `releases/tag/v3.0.0` but there's no GitHub Release there~~ done — fixed — GitHub Releases published | ~~Link is dead until a release is created~~ | ~~Create the release or change to `compare/v2.10.4...v3.0.0`~~ |
| ~~8~~  | ~~**Sub-module deps table in AGENTS.md** still lists `huh/v2`, `glamour/v2` as core deps~~ done — AGENTS.md Optional Sub-Module Dependencies table | ~~These are now sub-module deps~~ | ~~Move them to a separate sub-module deps table~~ |
| ~~9~~  | ~~**Design Principles section in AGENTS.md** has 22 numbered principles, ~8 of which reference deleted/moved features~~ done — AGENTS.md principles rewritten for v4 (2026-11-02) | ~~The list should be pruned to match v3 reality~~ | ~~Rewrite this section~~ |
| ~~10~~ | ~~**`Gotchas` section in AGENTS.md** has extensive notes on glamour/prompts/spinner/telemetry behavior — these moved to sub-modules~~ done — AGENTS.md gotchas document the sub-modules | ~~Either move the gotchas to sub-module docs or add a note that they apply to the sub-modules~~ | ~~~~ |

---

## f) Next 25 Things to Get Done

### Critical (do first)

1. ~~**Fix AGENTS.md fully** — remove deleted files from structure tree (`editor.go`, `result.go`, `telemetry.go`, `glamour.go`, `spinner.go`, `manpage.go`), update deps table versions (`go-output v0.30.1`, `samber-do-auditlog v0.4.0`), move `huh/v2` + `glamour/v2` to sub-module deps, remove deleted features from Design Principles (#11 EditInEditor, #16 WithGlamourHelp, #20 Result[T] sum types, #17 TelemetryMiddleware), update all `[T]` type params on now-non-generic options~~ done — AGENTS.md fully rewritten through v4 (2026-11-02)
2. ~~**Fix FEATURES.md deps table** — update `go-output v0.23.3` → `v0.30.1`, `samber-do-auditlog v0.3.1` → `v0.4.0`, move sub-module deps to separate table~~ done — FEATURES.md dep tables maintained through v4 releases
3. ~~**Create GitHub Release for v3.0.0** — with migration notes and breaking changes summary~~ done — GitHub Releases published (gh release list)
4. ~~**Create GitHub Release for v2.10.4** — with retract explanation~~ done — GitHub Releases published (gh release list)
5. ~~**Write `docs/MIGRATION_v2_v3.md`** — breaking changes, import path migration, before/after code examples, sub-module adoption guide~~ done — docs/MIGRATION_v2_v3.md exists

### High Priority

6. ~~**Add CHANGELOG entries for v2.10.3 and v2.10.4** — currently no changelog content for these tags~~ done — CHANGELOG [2.10.4] shipped (2026-07-07); v2.10.3 left tag-only (v2 frozen)
7. ~~**Consumer smoke test** — create a throwaway module, `go get github.com/larsartmann/cmdguard/v3@v3.0.0`, write a minimal CLI, verify it compiles and runs~~ done — closed per Resolution (2026-07-23)
8. ~~**Remove `/v3` from `.gitignore`** — BuildFlow artifact, wrong intent~~ done — clean — no /v3 in .gitignore
9. ~~**Update `docs/API.md`** — verify all function signatures match v3 (non-generic `CommandOption`, positional flags in `NewCommand`, etc.)~~ **Superseded — docs/API.md was removed in the 2026-08-06 docs refresh; the contract lives in FEATURES + the describe JSON.**
10. ~~**Audit `Gotchas` section in AGENTS.md** — move sub-module-specific gotchas (glamour/prompts/spinner/telemetry) to their respective packages or add cross-references~~ done — AGENTS.md gotchas document the sub-modules

### Medium Priority

11. ~~**Sync `TODO_LIST.md`** — remove v2-era items, add v3 post-migration tasks~~ done — TODO_LIST.md maintained (v5 sections, 2026-11-02)
12. ~~**Audit `ROADMAP.md`** — verify plans are still relevant for v3 architecture~~ done — ROADMAP.md reflects v4/v5 direction (2026-11-02)
13. ~~**Update `examples/taskctl/README.md`** — verify all API examples match v3 signatures (path refs are clean, but code examples may be stale)~~ done — cleaned — examples maintained through v4
14. ~~**Create `docs/adr/002-v3-module-migration.md`** — ADR documenting why the module path changed and the v2.11.0 mis-release incident~~ **Superseded — superseded — migration documented in CHANGELOG [3.0.0] + docs/MIGRATION_v2_v3.md; docs/adr/ took a different numbering.**
15. ~~**Verify sub-module READMEs or docs** — glamour/manpage/prompts/spinner/telemetry packages may need import path updates in their godoc examples~~ done — sub-module docs maintained through v4 (manpage later removed)
16. ~~**Update `flake.nix`** if it contains any version strings or module path references~~ done — flake.nix carries no version strings
17. ~~**Add `// Deprecated:` comments** on any v2 compatibility shims if consumers need gradual migration~~ done — // Deprecated: markers on SetConfig/RegisterInScope (ROADMAP v5 table)
18. ~~**Trigger pkg.go.dev indexing** — visit `https://pkg.go.dev/github.com/larsartmann/cmdguard/v3` to request a refresh~~ done — indexed — README badge links pkg.go.dev/v4

### Lower Priority

19. ~~**Audit all `docs/status/` historical reports** — decide whether to add a header note clarifying they describe v2 architecture~~ done — status reports carry Update/Resolution headers (docs-health passes)
20. ~~**Update `docs/design/` design docs** — at minimum add a note that the v2 design doc describes the predecessor architecture~~ **Superseded — v3.0.0 shipped 2026-07-07 (docs/planning/2026-07-06_06-54); docs/MIGRATION_v2_v3.md exists.**
21. ~~**Consider a `go-output` version audit** — verify all 10 go-output sub-modules are at v0.30.1 across all 6 go.mod files~~ done — go-output v0.38.7 lockstep across go.mod files
22. ~~**Run `golangci-lint` with `--new-from-rev=v3.0.0`** to establish a clean baseline for future development~~ done (docs-health pass 2026-11-08)
23. ~~**Update benchmark comments** — `benchmarks/guard_bench_test.go` package comment was updated but internal comments may still reference v2 patterns~~ done — benchmarks/guard_bench_test.go maintained
24. ~~**Verify `go.work.sum`** doesn't exist or is correct for the workspace~~ done — go.work.sum committed at repo root
25. ~~**Consider adding a `CONTRIBUTING.md`** note about the v2→v3 split for contributors who forked v2~~ done — CONTRIBUTING.md exists at repo root

---

## g) Top #1 Question I Cannot Answer Myself

**Should the v2 module path (`github.com/larsartmann/cmdguard/v2`) continue to receive the `retract` directive on the main `master` branch, or is the `release/v2.10` branch sufficient?**

Here's the dilemma:

- The `retract v2.11.0` directive lives in `release/v2.10`'s `go.mod` — which means it's tagged `v2.10.4` and visible to the Go proxy.
- But `master` now has module path `/v3`. The old `/v2` module path's `latest` version is whatever the proxy considers the highest tag on the `/v2` import path.
- **I cannot verify from inside this repo whether `go get github.com/larsartmann/cmdguard/v2@latest` resolves to `v2.10.4` (retracted v2.11.0) or still serves `v2.11.0`.** The proxy's behavior depends on whether it has seen the `v2.10.4` tag's `go.mod` yet.
- Testing this requires a **clean external module** (no replace directives, no workspace) — which I cannot create from within this repo.

**The action I need from you:** Run this in a throwaway directory outside the cmdguard repo:

```bash
mkdir /tmp/smoke-test && cd /tmp/smoke-test
go mod init smoke
go get github.com/larsartmann/cmdguard/v2@latest
# Does it resolve to v2.10.4? Does it show a retraction warning?
go get github.com/larsartmann/cmdguard/v3@latest
# Does v3 resolve and compile?
```

This will confirm the proxy state and whether the retract is effective.

**Answered:** moot — the v2 line is frozen at v2.10.4 with the `retract v2.11.0` directive shipped; the main line moved to `/v3` (2026-07-07) then `/v4`; v2.11.0 stays proxy-cached with retract as the standing mitigation.

---

## Metrics Summary

<!-- check-rows:ledger -->
| Metric                            | Value                                                |
| --------------------------------- | ---------------------------------------------------- |
| Commits this session              | 4 (migration + docs + v2.10.4 retract + this report) |
| Files changed (migration)         | 159                                                  |
| Internal refs updated             | 731                                                  |
| Test runs                         | 1831 (all pass)                                      |
| Test functions                    | ~457 (1371 top-level RUN entries)                    |
| Benchmarks                        | 26                                                   |
| Fuzz targets                      | 7                                                    |
| Coverage (core)                   | 87.3%                                                |
| Coverage (configload)             | 87.7%                                                |
| Lint issues                       | 0 (exit 0)                                           |
| Race conditions                   | 0                                                    |
| Build errors                      | 0                                                    |
| Tags created                      | v3.0.0, v2.10.4                                      |
| Tags deleted                      | v2.11.0 (local + remote)                             |
| Branches created                  | release/v2.10                                        |
| Modules migrated                  | 6 (root + 5 sub-modules)                             |
| Module proxy resolves v3.0.0      | YES                                                  |
| Module proxy still serves v2.11.0 | YES (cached — retract is the mitigation)             |

---

## Honest Self-Assessment

**What went well:** The mechanical migration was clean — `git mv` preserved history, sed was precise (no false positives on version strings like `v2.0.0`), all 6 modules tidied, full test suite passes. The release strategy (retract + branch + tag) is correct per Go conventions.

**What went poorly:** I treated AGENTS.md and FEATURES.md as "path updates" when they needed **architectural rewrites**. The v3 redesign deleted 6 core files and moved 5 features to sub-modules — but I only updated import paths, not the content describing those features. A contributor or AI session reading AGENTS.md today will be misled about what the codebase actually contains.

**Lesson:** Bulk find-replace is the easy part. The hard part — and the part that actually matters — is verifying that documentation **describes reality** after the change, not just that strings were substituted.

## Resolution (2026-07-23)

- §b GitHub Releases and missing docs were shipped in the 2026-07-07 and 2026-07-10 sessions.
- §c "NOT STARTED" items (CI smoke test, `pkg.go.dev` refresh, sub-module READMEs) were closed in `f8f3ad4`, `cccfdc9`, and subsequent work.
- `manpage` was removed post-v3.0.0; current sub-modules are glamour, prompts, spinner, telemetry.
- Current metrics: 470 test functions, 1434 runs, 87.8% core coverage, 0 lint issues.
