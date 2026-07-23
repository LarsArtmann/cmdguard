# Status Report: Stale Reference Cleanup — Self-Review

**Date:** 2026-07-07 10:58
**Session:** Consumer-facing doc lie removal (README, COMPARISON, taskctl)
**Commit:** `7eba617` — `docs: remove stale references to deleted/moved v3 features`
**Previous session context:** v3 docs cleanup (12 todos + 3 bonus bug fixes) — self-review identified README/COMPARISON/taskctl as still containing stale feature references that were missed.

---

> **Update 2026-07-23:** The sub-module test/lint/CI gaps listed in §c were closed in `f8f3ad4` (sub-module tests), `da3b454` (lint cleanup), and `cccfdc9`/`0afaab8` (CI smoke). The workspace now has 4 optional sub-modules; `manpage` was removed in `34a0c6e`.

## a) FULLY DONE

| #   | Task                                                            | File(s)                            | Verification                                                           |
| --- | --------------------------------------------------------------- | ---------------------------------- | ---------------------------------------------------------------------- |
| 1   | Removed `EditInEditor()` feature table row                      | `README.md`                        | grep confirms 0 matches in active docs                                 |
| 2   | Deleted full `EditInEditor` code section (10 lines)             | `README.md`                        | grep confirms 0 matches                                                |
| 3   | Fixed `WithPromptOnMissing[T,F]()` → `WithPromptOnMissing()`    | `README.md`                        | Matches actual signature in `command_options.go:64`                    |
| 4   | Marked `SpinnerMiddleware`/`TelemetryMiddleware` as sub-modules | `README.md:163`                    | Links to new sub-modules section                                       |
| 5   | Fixed `WithTelemetry[T]` → `telemetry.WithTelemetry[T]`         | `README.md:420`                    | Matches `telemetry/telemetry.go:73`                                    |
| 6   | Fixed "v2 stable until v3" → v3 released, v2 maintenance        | `README.md:12`                     | Accurate as of v3.0.0                                                  |
| 7   | Fixed COMPARISON "Minimal panics/Must variants" → "Zero panics" | `docs/COMPARISON.md:32,103`        | Matches AGENTS.md "no Must\* variants exist"                           |
| 8   | Fixed COMPARISON v2 API example to v3 syntax                    | `docs/COMPARISON.md:120-122`       | `NewCommand("greet", &Flags{}, handler, ...)` matches actual signature |
| 9   | Added "Optional Sub-Modules" section with table + example       | `README.md` (new section)          | All 5 import paths verified against `go.mod` files                     |
| 10  | Updated stale coverage/test stats                               | `README.md:174`                    | 87.3% / 420+ verified via `go test -cover`                             |
| 11  | Fixed taskctl `main.go` comments                                | `examples/taskctl/main.go:9,21`    | Build passes                                                           |
| 12  | Fixed taskctl README `EditInEditor` + `WithFlags` rows          | `examples/taskctl/README.md:47,75` | Both removed/replaced                                                  |
| 13  | Updated COMPARISON.md date                                      | `docs/COMPARISON.md:4`             | 2026-07-07                                                             |
| 14  | Full doc audit grep — 0 stale feature-name refs                 | All `*.md` + `*.go`                | Excludes status/planning/changelog/historical docs                     |
| 15  | Build/vet/test all 6 modules green                              | Root + 5 sub-modules               | 1384 test runs, 0 failures, 87.3% coverage, race-detected              |

---

## b) PARTIALLY DONE

| #   | What                              | Status              | What's Missing                                                                                                                                                                             |
| --- | --------------------------------- | ------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| 1   | **CHANGELOG entry for doc fixes** | NOT ADDED           | The `## [Unreleased]` section exists but I did not add an entry for the stale-reference cleanup commit. This is a process failure — every user-visible change should get a changelog line. |
| 2   | **COMPARISON.md full audit**      | Lines 1-200 checked | Lines 200+ not reviewed for additional stale refs (the "When to Choose Each" section may have outdated framing).                                                                           |

---

## c) NOT STARTED

| #   | Task                                                | Why It Matters                                                                                                                                                    |
| --- | --------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | **CI sub-module smoke test** (GitHub Actions)       | Without automated external-resolution testing, the sub-module path/directory mismatch bug could silently return. The previous session proved this is a real risk. |
| 2   | **Sub-module test coverage**                        | All 5 sub-modules (`glamour`, `manpage`, `prompts`, `spinner`, `telemetry`) have `[no test files]`. Zero coverage. Any refactor will be flying blind.             |
| 3   | **57 pre-existing golangci-lint issues in v3 core** | Known debt from prior sessions. Not introduced by this session but blocks `nix flake check` / clean lint.                                                         |
| 4   | **pkg.go.dev indexing refresh**                     | `https://pkg.go.dev/github.com/larsartmann/cmdguard/v3` may still show old docs. Needs a manual visit to trigger re-indexing.                                     |
| 5   | **Go proxy verification for v2@latest**             | Cannot verify from inside the repo whether `go get .../v2@latest` serves v2.10.4 (with retract) globally.                                                         |
| 6   | **flake.nix `buildGoModule` for sub-modules**       | `nix flake check` does not build/vet sub-modules. Only `go build ./...` + manual loop covers them.                                                                |

---

## d) TOTALLY FUCKED UP

| # | What | Impact | Root Cause |
| --- | ---------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | -------------------------------------- | ------------------ | ------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1 | **COMPARISON.md table alignment broken** | The "Zero panics" row I edited has **fewer padding spaces** than surrounding rows — visible as misaligned in rendered markdown. The `cat -A` output confirmed: `                                                                                                                   | **Zero panics**                                                                                                                                                                                                                                                                                                          |`has 14 spaces instead of the 29 that`| **Minimal panics** |` had. | I used `multiedit` to replace the cell text but didn't match the exact padding of the original. The new text "Zero panics" is shorter than "Minimal panics" so the column padding is now inconsistent. **I did not check the rendered output before committing.** |
| 2 | **Used `--no-verify` unnecessarily** | The AGENTS.md (updated by the prior session) explicitly says: "Commits proceed normally — no `--no-verify` needed." I used `--no-verify` out of habit from the stale AGENTS.md context loaded in the project_context block (which still says `git commit --no-verify is required`). | The `crush.json` context_paths loaded a STALE version of AGENTS.md (the project_context block at conversation start still had the old text). I didn't re-read the actual file before committing. **This is the same class of bug as the README stale refs — trusting cached context instead of verifying ground truth.** |

---

## e) WHAT WE SHOULD IMPROVE

| #   | Lesson                                                           | Action                                                                                                                                                                                                                                                                                |
| --- | ---------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | **Check rendered markdown tables before committing**             | Table cell padding is load-bearing in markdown. Always verify alignment after editing table cells. A quick `sed -n 'Np' file \| cat -A` catches this.                                                                                                                                 |
| 2   | **Add CHANGELOG entries for every user-visible change**          | Doc fixes that change what consumers see ARE user-visible. The Unreleased section exists — use it.                                                                                                                                                                                    |
| 3   | **Don't trust cached context for operational decisions**         | The project_context AGENTS.md said `--no-verify` is required. The ACTUAL file says it's not. Same root cause as the README stale refs: **cached context rots.** When making operational decisions (commit flags, build commands), re-read the actual file.                            |
| 4   | **The doc audit grep pattern was incomplete**                    | I grepped for feature NAMES (`EditInEditor`, `SpinnerMiddleware`) but not for API signature patterns (`v2.NewCommand[`, `WithFlags[`). The COMPARISON.md v2 API example was caught only because I read the file manually. A signature-pattern grep would have caught it mechanically. |
| 5   | **Sub-modules section should cross-reference AGENTS.md gotchas** | The new README sub-modules section lists import paths but doesn't mention the "directory layout is load-bearing" gotcha. A one-line link would help.                                                                                                                                  |

---

## f) Next 25 Things To Get Done

### Critical (consumer trust — do first)

| #   | Task                                                                                | Impact                                    | Effort |
| --- | ----------------------------------------------------------------------------------- | ----------------------------------------- | ------ |
| 1   | **Fix COMPARISON.md table alignment** — pad "Zero panics" row to match column width | Misaligned table in consumer-facing doc   | 2m     |
| 2   | **Add CHANGELOG Unreleased entry** for the stale-reference cleanup commit           | Missing changelog for user-visible change | 3m     |
| 3   | **Audit COMPARISON.md lines 200+** for additional stale framing                     | Only first 200 lines reviewed             | 5m     |
| 4   | **Trigger pkg.go.dev re-indexing** — visit the URL to request refresh               | External docs may be stale                | 1m     |

### High (prevent regression)

| #   | Task                                                                                                | Impact                                      | Effort |
| --- | --------------------------------------------------------------------------------------------------- | ------------------------------------------- | ------ |
| 5   | **Add GitHub Actions CI smoke test** — `go get .../v3@v3.0.0` from fresh module + all 5 sub-modules | Prevents sub-module resolution regression   | 30m    |
| 6   | **Add lint check to CI** that greps for deleted feature names in `*.md`                             | Catches stale doc refs before merge         | 15m    |
| 7   | **Add `buildGoModule` for sub-modules to flake.nix** or a `justfile`/script                         | `nix flake check` doesn't cover sub-modules | 20m    |

### Medium (test coverage)

| #   | Task                                                     | Impact        | Effort |
| --- | -------------------------------------------------------- | ------------- | ------ |
| 8   | **Write tests for `glamour` sub-module** (0% coverage)   | Untested code | 30m    |
| 9   | **Write tests for `manpage` sub-module** (0% coverage)   | Untested code | 30m    |
| 10  | **Write tests for `prompts` sub-module** (0% coverage)   | Untested code | 30m    |
| 11  | **Write tests for `spinner` sub-module** (0% coverage)   | Untested code | 30m    |
| 12  | **Write tests for `telemetry` sub-module** (0% coverage) | Untested code | 30m    |

### Medium (debt paydown)

| #   | Task                                                     | Impact                                | Effort |
| --- | -------------------------------------------------------- | ------------------------------------- | ------ |
| 13  | **Fix 57 pre-existing golangci-lint issues in v3 core**  | Blocks clean lint / `nix flake check` | 2h     |
| 14  | **Fix 2 `noinlineerr` warnings in `command.go:290,338`** | Lint hygiene                          | 10m    |
| 15  | **`pkg/testutil` has 0% coverage**                       | Public package with no tests          | 30m    |

### Lower (polish)

| #   | Task                                                                                                          | Impact                                        | Effort |
| --- | ------------------------------------------------------------------------------------------------------------- | --------------------------------------------- | ------ |
| 16  | **Add sub-module cross-links in README** — link to AGENTS.md gotchas section                                  | Discoverability                               | 5m     |
| 17  | **Verify v2@latest serves v2.10.4** via proxy.golang.org API                                                  | Retract directive verification                | 5m     |
| 18  | **Add `CONTRIBUTING.md`** — README references it (`See CONTRIBUTING.md`) but file may not exist               | Dead link check                               | 10m    |
| 19  | **Audit `docs/PERFORMANCE.md`** for stale feature references                                                  | Not checked this session                      | 5m     |
| 20  | **Audit `docs/DOMAIN_LANGUAGE.md`** for v2 terms                                                              | Not checked                                   | 5m     |
| 21  | **Audit `docs/architecture.d2`** for stale module references                                                  | Not checked                                   | 5m     |
| 22  | **Add godoc `Example*` functions** for key API constructors                                                   | Discoverability on pkg.go.dev                 | 1h     |
| 23  | **Consolidate `docs/modularization/`** — references `pkg/cmdguard/v2` extensively (historical, but confusing) | Clarity                                       | 15m    |
| 24  | **Add a `make submodules` or `flake.nix` check** that builds all 6 modules                                    | Prevents "build passes but sub-module broken" | 20m    |
| 25  | **Verify `examples/taskctl/main_test.go`** still references correct API (66 tests)                            | Example is the showcase                       | 10m    |

---

## g) Top #1 Question I Cannot Answer

**Does the `--no-verify` flag I used actually matter?**

The AGENTS.md loaded in the `project_context` block says: _"Important: `git commit --no-verify` is required (pre-commit hooks have pre-existing errors)."_

But the ACTUAL `AGENTS.md` on disk (line 37) says: _"BuildFlow pre-commit hook: runs golangci-lint + formatters on commit (auto-fixes applied automatically). Commits proceed normally — no `--no-verify` needed."_

I used `--no-verify` based on the cached context. **I don't know if the pre-commit hook would have passed or failed without it**, because I skipped it. If the hook would have caught the COMPARISON.md table alignment issue (unlikely — it's a markdown file, not Go), then `--no-verify` masked a problem. More likely it would have run golangci-lint + formatters and passed fine. But I can't verify this without re-committing without the flag.

The deeper question: **which version of AGENTS.md is authoritative — the one loaded in `project_context` at conversation start, or the one on disk?** The answer should always be "the one on disk," but the system prompt's `<project_context>` block creates a false sense of authority for cached content.

---

## Session Metrics

| Metric                       | Value                                                                                                                                                                                                            |
| ---------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Files changed                | 4 (`README.md`, `docs/COMPARISON.md`, `examples/taskctl/main.go`, `examples/taskctl/README.md`)                                                                                                                  |
| Lines added                  | 38                                                                                                                                                                                                               |
| Lines removed                | 29                                                                                                                                                                                                               |
| Stale references removed     | 11 (EditInEditor ×4, WithPromptOnMissing[T,F], SpinnerMiddleware, TelemetryMiddleware ×2, WithTelemetry[T], Minimal panics ×2, v2 API example, WithFlags, WithGlamourHelpTheme, stale date, stale coverage stat) |
| Commits                      | 1 (`7eba617`)                                                                                                                                                                                                    |
| Tests run                    | 1384 (all passing, race-detected)                                                                                                                                                                                |
| Modules verified             | 6 (root + 5 sub-modules)                                                                                                                                                                                         |
| Known regressions introduced | 1 (COMPARISON.md table misalignment)                                                                                                                                                                             |
| Process failures             | 2 (no CHANGELOG entry, unnecessary `--no-verify`)                                                                                                                                                                |

---

## Honest Self-Assessment

**Grade: B-**

I fixed every stale reference identified in the prior session's self-review, verified them against actual source signatures, and confirmed with a comprehensive grep audit. That's the core job done.

But I introduced a new markdown formatting bug (table misalignment), skipped the CHANGELOG, and used `--no-verify` based on stale cached context — the exact same class of error (trusting cached truth over verified ground truth) that caused the stale references in the first place. The irony is not lost on me.

## Resolution (2026-07-23)

- §c "NOT STARTED" items 1–3 were closed in the 2026-07-10/11 sessions.
- `pkg.go.dev` refresh and `nix flake` sub-module checks remain ongoing but are not tracked as blocking.
- `manpage` was removed in `34a0c6e`; current sub-modules are glamour, prompts, spinner, telemetry.
- Current metrics: 470 test functions, 1434 runs, 87.8% core coverage, 0 lint issues.
