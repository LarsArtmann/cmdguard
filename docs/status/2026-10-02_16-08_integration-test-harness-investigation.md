# Status Report — Integration-Test Harness for cmdguard Consumers

**Date:** 2026-10-02 16:08 CEST (Friday)
**Scope:** How to give apps a better, easier way to run integration tests — specifically for **the command surfaces cmdguard owns** (flag parsing, precedence, validation, error/exit codes, output formats, DI wiring, config loading).
**Session type:** Research / design only. **Zero code written. Zero commits.**
**Repo state at start:** `master`, clean (`git status --short` empty).
**Author:** Crush (this session).

> **Format override:** This report is Markdown at the user's explicit request (skill canonical format is styled HTML). Flagged per the `status-report` skill contract; the override does not change the skill default.

---

## Session Summary

This session investigated why cmdguard consumers must hand-roll integration-test scaffolding, and what a better harness would look like. The headline finding: **v2 shipped an in-process CLI test harness (`TestCLI[T]`) and v4 dropped it silently.** Only the *assertion* helpers survived into `pkg/testutil`; the harness itself was never ported. So every app that adopts cmdguard re-invents the same boilerplate.

Nothing was implemented. This report captures the findings so they are not lost to the conversation.

---

## Stat Cards

| Category | Count |
| --- | --- |
| a) Fully done (analysis deliverables) | 10 |
| b) Partially done | 5 |
| c) Not started | 11 |
| d) Totally fucked up | 5 |

---

## a) FULLY DONE

These are **analysis deliverables** (investigation complete), not code. No code was shipped this session.

| # | Finding | Evidence |
| --- | --- | --- |
| 1 | v4 has **no** in-process CLI test harness. | `ls pkg/cmdguard/` → only `v4`; no `testkit`/`harness`/`testutil` under `pkg/cmdguard/v4/`. |
| 2 | v2 **did** have one, and it was deleted. | `git show fdd6dca^:pkg/cmdguard/v2/testutil/testutil.go` → 127 lines (`TestCLI[T]`, `TestResult`, `NewTestCLI`). Test file 151 lines. |
| 3 | Deletion context identified. | Deleted as **collateral** in `fdd6dca` "feat,config: add koanf-based config loader with JSON/YAML support" (diffstat: 21 files, `testutil.go −127`, `testutil_test.go −151`). |
| 4 | What survived vs what was lost. | Assertions `AssertOutputContains`, `AssertStringSlicesEqual` migrated to `pkg/testutil/panic_test_helpers.go`; the `TestCLI[T]`/`TestResult` harness did **not**. |
| 5 | Current app boilerplate mapped. | `examples/taskctl/main_test.go` hand-rolls `newTestCLI`, `mustExec`, `mustExecOnCLI`, `runAndFailOnError`, `expectError`. |
| 6 | Public surface usable as harness base mapped. | `ExecuteWithArgs` (`cli.go:521`), `Execute` (`cli.go:447`), `ExecuteAndExit` (`cli.go:528`), public `ExitCode`, `RootCommand()`, `Scope()`, `Injector()`, `Config()`, `SetConfig`, `SetOutputFormat` (`cli_accessors.go`, `cli_output.go`). |
| 7 | Capture-seam flaw in the old harness identified. | v2 captured via `root.SetOut/SetErr` only — this misses `fang` styling, `glamour` help, and `spinner`, which write to `os.Stdout`/`os.Stderr` directly. |
| 8 | The missing harness is a known-and-ignored gap. | `docs/status/2026-08-06_14-41_...md:188,223` — "FR integration test in taskctl — deferred as 'needs test harness' without attempting". |
| 9 | Three implementation options proposed. | (A) core package `pkg/cmdguard/v4/testkit` — *recommended*; (B) root sub-module `github.com/larsartmann/cmdguard/testkit`; (C) extend `pkg/testutil`. |
| 10 | Harness API shape sketched. | `testkit.New(t, cli)`; `Run(args...) *Result` (`Stdout`/`Stderr`/`Err`/`ExitCode()`); `WithEnv`; `WithConfig`; assertions (`AssertNoError`, `AssertErrorIs`, `AssertExitCode`). |

---

## b) PARTIALLY DONE

| # | What works now | What remains open | Blocker | Effort |
| --- | --- | --- | --- | --- |
| 1 | Harness API sketch (option A). | Capture-seam design unresolved (fd-level redirect vs. injectable writer through `CLI.spec`). | Needs a design decision. | M |
| 2 | Home recommendation (core package). | Owner decision not obtained. | Awaiting user. | S |
| 3 | v2 harness source recovered. | Not yet evaluated for v2→v4 API drift (`v2.AddCommand` vs v4; constructor/option changes). | None. | S |
| 4 | Scope-isolation strategy named (clone → override → invoke). | Not designed into the API; no proof. | None. | M |
| 5 | `Result` concept defined. | Type not finalized; exit-code path from `ExecuteAndExit` (must not call `os.Exit`) not solved. | None. | S |

---

## c) NOT STARTED

| # | Item | Why not started | Still wanted? |
| --- | --- | --- | --- |
| 1 | Any `testkit` code. | Design not approved. | Yes |
| 2 | Capture seam implementation. | Unresolved design (b-1). | Yes |
| 3 | Env fixture helper (`WithEnv`). | Awaiting package. | Yes |
| 4 | Config-file fixture helper (`WithConfig`). | Awaiting package. | Yes |
| 5 | Exit-code / sentinel assertion helpers. | Awaiting package. | Yes |
| 6 | Golden-file / snapshot helpers. | Awaiting package. | Yes |
| 7 | Table-driven run helper. | Awaiting package. | Yes |
| 8 | Migrate `examples/taskctl` to the harness. | Blocked on package. | Yes |
| 9 | Docs: README / AGENTS.md / `docs/testing.md`. | Blocked on package. | Yes |
| 10 | Tests for the testkit itself. | Blocked on package. | Yes |
| 11 | Sub-module `go.mod` + `go.work` wiring (only if option B). | Decision-gated. | Maybe |

---

## d) TOTALLY FUCKED UP

| # | What is broken | Severity | Root cause | Mitigation |
| --- | --- | --- | --- | --- |
| 1 | **Capability regression:** v4.0.1 ships with no in-process test harness while v2 had one. Consumers hand-roll integration scaffolding. | **High** (consumer DX; every adopting app pays it) | Deleted as collateral in `fdd6dca` with no replacement. | None. Workaround is per-app boilerplate (`taskctl/main_test.go`). |
| 2 | Deletion appears **undocumented** — no deprecation or migration pointer found for the lost `TestCLI`. | Medium-High (trust + migration pain) | Bulk refactor commit swept it out. | None. *(Not yet verified against `CHANGELOG.md` / `MIGRATION_*` — see item d-5.)* |
| 3 | A missing-harness gap was already recorded and left to rot. | Medium (process) | `docs/status` 2026-08-06 flagged "needs test harness" but never spawned a TODO. | None. |
| 4 | **Self, this session:** did not run `go build`/`go test` even though it is cheap; config-loader-era regression state unverified. | Low (no code changed) | Scope discipline — investigated only. | Run the suite before any implementation. |
| 5 | **Self, this session:** asserted a documentation regression without checking `MIGRATION_v2_v3.md`, `MIGRATION_v3_v4.md`, or `CHANGELOG.md`. | Low (unverified claim) | Did not grep those files. | Verify d-2 before acting on it. |

---

## e) WHAT WE SHOULD IMPROVE

| # | Suboptimal pattern | Impact | Suggested fix |
| --- | --- | --- | --- |
| 1 | Consumer testing ergonomics treated as an afterthought. | Every consuming app re-invents the harness. | Restore/evolve `testkit` as a maintained product surface. |
| 2 | Silent removal of public test packages. | Consumers discover the loss at migration time. | Standing rule: removal requires a migration note + replacement (encode in AGENTS/RELEASE). |
| 3 | Capture relies on `root.SetOut/SetErr`. | fang/glamour/spinner output escapes capture. | Make capture injectable at the `CLI` level (writer carried through `spec`), fd-level fallback. |
| 4 | Integration tests share mutable DI scope. | State leaks; not parallel-safe. | `testkit` clones + overrides scope per run by default. |
| 5 | `docs/status` deferrals die in timestamped files. | Known gaps age silently (d-3). | HARVEST deferrals into `TODO_LIST.md` at write time. |
| 6 | This session's findings not yet persisted. | Lost when context ends. | Run `docs-health` HARVEST on section (f). |

---

## f) Top 35 things we should get done next

Ranked by impact. Effort: S <30min / M 30min–2h / L >2h.

| # | Task | Impact | Effort | Category |
| --- | --- | --- | --- | --- |
| 1 | Design and implement the output capture seam (injectable writer through `CLI.spec`), covering fang/glamour/spinner. | Critical | M | Feature |
| 2 | Create the `testkit` package (`pkg/cmdguard/v4/testkit`) with `New(t, cli)` + `Run(args...) *Result`. | Critical | M | Feature |
| 3 | Prove the harness against `examples/taskctl` before generalizing further. | Critical | M | Feature |
| 4 | Decide home: core package vs root sub-module vs `pkg/testutil` (see question g-1). | High | S | Decision |
| 5 | Port the recovered v2 `TestCLI[T]`/`TestResult` to the v4 API (adapt `AddCommand`). | High | M | Feature |
| 6 | Implement per-run DI scope isolation (clone → override → invoke) in `testkit`. | High | M | Feature |
| 7 | `Result.AssertNoError`, `AssertErrorIs`, `AssertExitCode(n)` helpers. | High | S | Feature |
| 8 | Ensure `ExecuteAndExit` can be exercised without `os.Exit` (inject exit func or use `Execute`). | High | S | Feature |
| 9 | Parallel-safety test for `testkit` under `-race`. | High | S | Quality |
| 10 | Unit + integration tests for `testkit` itself. | High | M | Quality |
| 11 | HARVEST section (f) into `TODO_LIST.md` / `ROADMAP.md`. | High | S | Process |
| 12 | Verify `CHANGELOG.md` / `MIGRATION_*` record the v2 harness loss; add note if missing. | Medium | S | Documentation |
| 13 | `WithEnv(map)` helper that avoids the `t.Setenv` + `t.Parallel` panic. | Medium | S | Feature |
| 14 | `WithConfig(json/TOML/YAML)` temp-file fixture helper. | Medium | S | Feature |
| 15 | Golden-file assertion helper for help/docs/output. | Medium | M | Feature |
| 16 | Table-driven `RunTable(t, cases...)` helper. | Medium | S | Feature |
| 17 | Precedence-matrix assertion (flag > env > config > default). | Medium | M | Feature |
| 18 | Output-format rendering assertions (json/csv/table). | Medium | S | Feature |
| 19 | Capture `--output=json` structured errors in `Result`. | Medium | S | Feature |
| 20 | Migrate `examples/taskctl` to `testkit`; delete hand-rolled helpers. | Medium | M | Cleanup |
| 21 | Assertion for `ConfigFromContext[T]` resolved config. | Medium | S | Feature |
| 22 | Add fd-level (`os.Pipe`) fallback capture for stdout-direct writers. | High | M | Feature |
| 23 | Document `testkit` in README + AGENTS.md + new `docs/testing.md`. | Medium | S | Documentation |
| 24 | Minimal example integration test shipped inside `testkit`. | Medium | S | Documentation |
| 25 | Raw-cobra escape-hatch test support (`RunCLI`). | Low | M | Feature |
| 26 | Snapshot test for `GenerateDocs(w)` output. | Low | S | Feature |
| 27 | Prompt-runner test double for `prompt:`-tagged flags. | Low | M | Feature |
| 28 | DI audit-log export assertion helper. | Low | S | Feature |
| 29 | Fuzz harness helper feeding args into `ExecuteWithArgs`. | Low | M | Quality |
| 30 | Benchmark harness reuse in `benchmarks/` for command execution. | Low | M | Quality |
| 31 | Wire `testkit` into `docs/QUICKSTART.md` + `docs/TUTORIAL.md`. | Low | S | Documentation |
| 32 | Measure boilerplate reduction in taskctl (LOC delta after migration). | Low | S | Quality |
| 33 | Add `testkit` to `FEATURES.md` once implemented. | Low | S | Documentation |
| 34 | CI guard (depguard/forbidigo) nudging examples toward `testkit`. | Low | S | Quality |
| 35 | Sub-module `go.mod` + `go.work` entry — only if option B chosen. | Medium | S | Feature |

> **Handoff:** Section (f) is HARVEST input. It must be pulled into `TODO_LIST.md` / `ROADMAP.md` by `docs-health` HARVEST, or it dies in this timestamped file.

---

## g) Questions I cannot answer myself (max 3)

1. **Where should `testkit` live?** I recommend (A) core package `pkg/cmdguard/v4/testkit` — no new `go.mod`, zero heavy deps, imported only from `_test.go` so it never touches consumer binaries. Alternatives: (B) root sub-module for independent versioning, (C) fold into `pkg/testutil`. I cannot pick the versioning/ownership tradeoff for you.
2. **Is silent removal of a test-only public package ever acceptable?** If not, do you want a standing rule (migration note + deprecation window) encoded into AGENTS.md and the release process? This determines how item d-2 is handled.
3. **Should the harness take `*testing.T` directly, or stay `testing`-free and return plain results?** Direct `*testing.T` is simplest and matches v2; a testing-free core would also be usable from benchmarks and other runners. Both are defensible; it affects every helper signature.

---

## Process Notes

- **Commit:** skipped. Crush forbids commits unless the user says "commit"; the auto-commit daemon will pick up this file.
- **HARVEST:** not run this session — the user issued an explicit "WAIT FOR INSTRUCTIONS". Pending.
- **Not researched:** anything outside the integration-test-harness question, per instruction.
