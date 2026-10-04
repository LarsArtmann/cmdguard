# Status Report: TODO-List Execution — Context Middleware, Deprecations, Audit Upstream, Coverage

**Date:** 2026-10-04 07:13
**Session scope:** Execute the cmdguard TODO_LIST.md (all sections), verify each
item against source, implement what is real, sync all living docs.
**Outcome:** 11 of 13 TODO items closed (8 were already done — stale entries),
P5 delivered as a new v4 feature, P1/P3 deprecation-marked, F1 upstream-blocked
with the upstream API implemented in ../samber-do-auditlog, taskctl coverage
73.2% → 85.9%. One regression: my CHANGELOG entries were clobbered by a
concurrent session and the re-apply failed on a mid-air collision.

---

## a) FULLY DONE (verified this session)

### Stale TODO entries verified complete (removed from TODO_LIST.md)

All verified against source before removal — none were re-done blindly:

| #   | Item                                       | Evidence                                                                                  |
| --- | ------------------------------------------ | ----------------------------------------------------------------------------------------- |
| D3  | CONTRIBUTING.md v3→v4 drift                | Lines 101/115 already say v4; fixed in commit `ef694eb`                                   |
| D4  | ERROR_REFERENCE.md title v2→v4             | Line 1 already "cmdguard v4"                                                              |
| D5  | docs/MIGRATION_v3_v4.md missing            | File exists, complete (171 lines)                                                         |
| D6  | manpage note in MIGRATION_v2_v3.md §3      | Note present at lines 116-118                                                             |
| D7  | git corruption (fsck broken links, reflog) | `git fsck --no-dangling` clean; reflog entry `3e483b3b` gone                              |
| D8  | Tag flightrecorder v0.1.0                  | `flightrecorder/v0.1.0` tag exists                                                        |
| D9  | Lost flightrecorder godoc examples         | `ExampleRecorder_CaptureToWriter` (:62) + `ExampleWithFlightRecorderRecorder` (:95) exist |
| D10 | `go tool trace` parseability test          | `TestTraceSnapshot_IsParseableByGoToolTrace` (recorder_test.go:687)                       |
| D1  | pkg/testutil coverage                      | Measured 80.3% (was 49.6%) — already improved by a prior session                          |

### P5 — Middleware context propagation (NEW v4 feature, non-breaking)

- `ContextMiddleware[T]` — `next func(context.Context) error` threads derived
  contexts through the chain (`middleware.go`).
- `WithContextMiddleware[T](mw...)` — CLIOption wiring; context middleware run
  OUTSIDE plain middleware, derived context reaches plain middleware + handlers.
- `TimeoutMiddleware[T](d)` — reference implementation; deadline errors match
  both new sentinel `ErrCommandTimeout` and `context.DeadlineExceeded`.
- `buildContextChain` + lazy plain-chain rebuild inside `wireHandlerWithMiddleware`
  (`cli_command.go`); zero behavior change when unused.
- 9 new tests in `middleware_context_test.go` (propagation to handler, reach into
  plain middleware, full-chain flow, ordering, error propagation, timeout hit/pass,
  non-positive-d passthrough, sealed-interface type-mismatch) — all green with
  `-race`. Core suite green with `-race`.

### P1/P3 — Deprecations (v4-safe treatment of v5 items)

- `CLI.SetConfig` → `// Deprecated:` + `TODO(v5)` marker, replacement documented
  (`cli_accessors.go`). Removal itself is v5 (shipped API in v4.0.0-v4.0.2 —
  removal would break semver).
- `RegisterInScope` → `// Deprecated:` + `TODO(v5)` marker; documented why it is
  broken (everything registers as `any`, two providers collide) and the
  replacement pattern `parent.Child(name)` + `do.Provide(child.Injector(), ...)`
  (`scope.go`).
- P2 (`Get` rename) and P4 (`Package` redesign) confirmed as ROADMAP v5 items —
  no v4 action; removed from TODO_LIST (they were duplicates of ROADMAP §v5).

### D2 — taskctl coverage 73.2% → 85.9%

- `main.go` restructured for testability: `run(ctx, args)` (full production
  composition incl. recorder lifecycle) + `buildApp(rec)` + `newProductionRecorder()`
  - `exportAuditLog(cli)` + thin `main()`; taskctl now demonstrates the explicit
    `WithFlightRecorderRecorder` lifecycle (the sub-module's documented pattern)
    instead of the never-stopping `WithFlightRecorder`.
- `run_test.go`: composition test, audit export (html/json/bogus-format/no-plugin/
  write-failure), unknown command, store-construction failure across 5 commands
  (covers every `resolveStore` error branch), nil-injector guard, recorder config.
- Verified green 3× with `-race -shuffle=on` (shuffle-proofing mattered — see d).

### F1 — Command-level audit middleware (upstream half)

Implemented in ../samber-do-auditlog (local, Unreleased):

- `EventTypeCommand "command"` + label/color meta + `Event.IsCommand()`.
- `Recorder.RecordCommand(scopeID, scopeName, commandName, phase, durationMs, err)`
  and `Plugin.RecordCommand(...)` — command events enter the event stream only;
  NO service records created (verified by test against `Report.Services`).
- Design kept upstream-light: command name rides `ServiceRef.ServiceName`
  (documented overload), free-form JSON schema → no schema regen, `IsKnown()`-
  based load validation accepts the new type automatically.
- 5 tests incl. before/after pairing, error recording, no-service-records,
  OnEvent callback firing; full upstream package suite green.

### Docs sync (living docs)

- `TODO_LIST.md` — rewritten per docs-health rules: done items deleted (no
  "resolved" sections), v5 marker table extended to T1-T5, new "Blocked on
  Upstream" section for F1 with the exact unblock recipe.
- `FEATURES.md` — Custom middleware row → FULLY_FUNCTIONAL; new rows for
  `ContextMiddleware`/`WithContextMiddleware`/`TimeoutMiddleware`; SetConfig/
  RegisterInScope rows marked Deprecated; audit row points at upstream state;
  coverage table updated (testutil 80.3%, taskctl 85.9%).
- `ROADMAP.md` — ctx propagation marked delivered-additively (v5 = unify
  signatures, new candidate #8); audit direction marked blocked-on-release;
  v5 tables renumbered with T4/T5 carrying source markers now.
- `AGENTS.md` — v4.0.1→v4.0.2 status, auditlog dep v0.8.1→v0.10.0 (go.mod truth),
  3→5 TODO(v5) markers with locations, middleware section documents
  ContextMiddleware, telemetry gotcha updated, auditlog bullet documents the
  blocked wiring.
- CHANGELOG (sibling repo) — Unreleased "Added — Command execution events" entry.

---

## b) PARTIALLY DONE

1. **F1 cmdguard side (BLOCKED, by design):** `AuditMiddleware[T]` is NOT wired.
   It cannot ship until samber-do-auditlog is pushed + tagged ≥ v0.11.0 — a
   committed call to `Plugin.RecordCommand` against go.mod's v0.10.0 would break
   every non-local build (workspace `replace` doesn't travel). TODO_LIST carries
   the exact unblock steps.
2. **CHANGELOG.md (cmdguard):** my Unreleased entries (Added/Deprecated/Changed +
   backfilled [4.0.2]) were written, then lost to a concurrent-session overwrite,
   then my re-apply edit failed on a mid-air "file modified" collision. The file
   currently has NONE of my entries. Content is prepared; needs one clean re-apply
   (2-minute task) — deliberately not retried mid-report.

---

## c) NOT STARTED (deliberate, this session)

- v5 breaking changes T1-T5 (renames/removals) — deferred by project decision;
  source markers + ROADMAP rows maintained instead.
- `Get[T]` rename, `Package[T]` redesign — ROADMAP v5 items, no v4 work.
- F1 cmdguard wiring — blocked on upstream release (see b).
- telemetry sub-module migration to ContextMiddleware (span-context propagation)
  — noted in AGENTS.md as future work.

---

## d) WHAT WENT WRONG (honest ledger)

1. **CHANGELOG entries clobbered twice** (see b). I verified doc numbers late —
   caught it only because the coverage figure wasn't where I expected. Lesson:
   with a concurrent session active, re-verify my doc edits after daemon commits.
2. **Flightrecorder singleton leak broke a test under shuffled order:** my first
   taskctl approach used `WithFlightRecorder` (lazy start, never stops). BuildFlow's
   shuffled test-race run made `TestFR_Integration_CaptureOnCommandError` fail
   (`rec.Start()` → "already enabled"). Fixed by explicit recorder lifecycle
   (per-test `t.Cleanup(rec.Stop)`; `run()` stops the production recorder). My
   local runs had passed because file order ran coverage_test.go first — I
   initially shipped a test that was order-dependent and didn't know it.
3. **Coverage numbers churned in docs:** wrote 80.1% into three docs, refactor
   dropped real coverage to 76.9% → 78.5% → 85.9% final. Should have measured
   once at the end; two docs carried a stale number briefly.
4. **`db status` in the store-failure table** — that command doesn't touch the
   store; test failed once, fixed by dropping it. Read the handler before
   assuming DI usage.
5. **Two edit-tool mid-air failures** (taskctl main.go tail, CHANGELOG §2) from
   whitespace/context mismatch — recovered by re-reading; wasted a cycle each.
6. **`rg -rn` misuse** (replace flag) mangled one grep — harmless, caught.
7. **BuildFlow full gate is red on environment, not code:** license-check
   (go-licenses not in PATH outside `nix develop`), nix-build-verify/treefmt
   (sandbox DNS refused → go1.27.0 toolchain download fails). Pre-existing,
   reproduced across 6+ runs before mine. Go build/test/race/lint steps green.
8. **lint/ sub-module is mid-refactor by a concurrent session** (working-tree
   changes I didn't author + a parse error in `lint/cmd/cmdguard-lint/main.go`
   - their own status file written 2 min ago). Workspace-wide builds that fan
     into ./lint fail through no fault of this session's changes. Left untouched
     per "never revert changes you didn't author".

---

## e) WHAT WE SHOULD IMPROVE

1. TODO_LIST had 8 of 13 items stale — the harvest/verify loop runs too rarely;
   consider the docs-health drift alarm (verify-checklist reference) for this repo.
2. AGENTS.md project structure omits the `lint/` sub-module entirely (and says
   "6 modules" while go.work now has 7 incl. ./lint + local go-output replaces
   from the concurrent session). Structure section is drifting.
3. AGENTS.md exceeds the doctor's 377-line budget (386 pre-session → ~395 now
   despite my condensing). Needs a deliberate content split, not more trimming.
4. taskctl's `WithFlightRecorder` option (never stops the process-wide recorder)
   is a footgun for embedders/tests; consider documenting Stop ownership in the
   sub-module README or providing an auto-stop variant.
5. CHANGELOG discipline: releases (v4.0.2) shipping without a CHANGELOG entry —
   the 4.0.2 backfill I wrote from git evidence shouldn't have been needed.
6. BuildFlow env gaps (go-licenses PATH, offline nix toolchain download) make the
   "full" gate permanently red locally — either fix the env or scope skip_steps
   so the gate's signal is trustworthy.

---

## f) NEXT UP TO 50 (roughly ordered)

1. Re-apply CHANGELOG Unreleased entries + [4.0.2] backfill (content ready).
2. Re-run `buildflow --build-mode full` scoped to my modules after CHANGELOG fix.
3. Push + tag samber-do-auditlog v0.11.0 (needs Lars's go-ahead).
4. After 3: bump cmdguard go.mod auditlog → v0.11.0.
5. After 4: implement `AuditMiddleware[T]` + `WithAuditMiddleware` + tests (F1 close).
6. After 5: taskctl demo of AuditMiddleware; FEATURES row → FULLY_FUNCTIONAL.
7. Migrate telemetry sub-module to ContextMiddleware (span ctx propagation).
8. Add `telemetry.ContextTelemetryMiddleware` or fold into v0.3.0 of telemetry.
9. Consider spinner sub-module: nothing needed (no ctx), verify once.
10. v5 prep: spec the Middleware/ContextMiddleware unification (breaking).
11. v5 prep: `CommandInfo`→`CommandMetadata` mechanical-rename plan (T1).
12. v5 prep: `TypeHandler`→`TypeCodec` + deprecated alias plan (T2).
13. v5 prep: `PromptRunner`→`HuhPrompter` (T3).
14. v5: remove `SetConfig` (T4) — update doc comments referring to it.
15. v5: generic `RegisterInScope[T]` (T5) with variadic-provider design doc.
16. v5: `Get[T]`→`GetService[T]` (candidate 6).
17. v5: `Package[T]` redesign (candidate 7).
18. Document `lint/` sub-module in AGENTS.md structure + module count.
19. Update AGENTS.md "6 modules" → current go.work truth (incl. go-output replaces).
20. AGENTS.md size diet to ≤377 lines (move detail to docs/).
21. Investigate nix-build-verify 6/6 failures (offline toolchain fetch).
22. Install go-licenses in devShell/PATH so license-check can run locally.
23. Decide skip_steps for sandbox-hostile nix checks (or make them hermetic).
24. taskctl: cover remaining buildCommands error branches via fault injection?
25. taskctl: `db` parent + `env create --force` paths coverage.
26. taskctl: prompt-tag path (`prompt:""` + WithPromptOnMissing) test with prompts stub.
27. pkg/testutil: push 80.3% → 90% (failure_paths_test.go branch completion).
28. Add fuzz corpus seeds for the 8 targets (ROADMAP deferred #7).
29. cmdguard-lint: once the concurrent refactor lands, add a rule for deprecated-API usage.
30. docs/API.md: add ContextMiddleware/TimeoutMiddleware reference sections.
31. README.md: feature list mention of context middleware.
32. docs/COBRA_FOOTGUNS.md: cross-link WithContextMiddleware for ctx threading.
33. ExportAuditLog: consider Writer-based test instead of read-only-dir (portability).
34. Flightrecorder: document Stop-ownership semantics in README (embedders).
35. Flightrecorder: optional auto-stop middleware variant (owns lifecycle).
36. Samber-do-auditlog: HTML viz "Command" color distinct from Registration (both --accent).
37. Samber-do-auditlog: consider dedicated `command_name` field over ServiceRef overload (schema 0.4.0).
38. Samber-do-auditlog: cmdguard-lint-style rule or doc forbidding ServiceName pollution?
39. CHANGELOG entry for 4.0.2 should have existed at tag time — add release checklist item.
40. Benchmarks for ContextMiddleware overhead vs plain (benchmarks/ dir).
41. Table-driven test consolidation in middleware_context_test.go (minor).
42. TODO(v5) marker line numbers in ROADMAP drift already (I removed line pins — good); re-pin? no — keep unpinned, note convention.
43. git-town/library-policy verification (ROADMAP deferred #10).
44. Re-audit window per ROADMAP #12: after F1 lands.
45. Consider `WithMiddlewareBoth` ergonomic if users complain about two lists.
46. Verify fang + context middleware interplay (error display unaffected) — test exists indirectly via ExecuteWithArgs; add explicit case.
47. WithCleanup + ContextMiddleware ordering test (cleanup outside chain today).
48. `TimeoutMiddleware` jitter/backoff variant? (YAGNI — only on demand.)
49. Status-report this file → TODO_LIST harvest of items 1-23.
50. Schedule docs-health AUDIT after F1 + v5 prep docs settle.

---

## g) QUESTIONS FOR LARS (cannot self-answer)

1. **Push permission:** may I `git push` + tag `samber-do-auditlog` v0.11.0
   (and later cmdguard), or do you want to review the upstream API first
   (`RecordCommand` overloads `ServiceRef.ServiceName` with the command name —
   the alternative is a schema-0.4.0 `command_name` field, item 37)?
2. **Concurrent session on lint/:** another agent is mid-refactor of the
   `lint/` sub-module (working-tree changes + a parse error right now). Leave
   it alone, or should I take over/finish it once you confirm it's abandoned?
3. **CHANGELOG collision:** my Unreleased entries were overwritten once by that
   concurrent session. Re-apply immediately (risking another collision while
   it's active), or wait for your go-ahead / its completion?

---

**Bottom line:** every actionable TODO item is either closed-and-verified or
blocked with the blocker named. The only session regression is the
twice-clobbered cmdguard CHANGELOG (content ready, 2-minute re-apply), and the
one engineering lesson worth keeping: test-order dependence around the
flightrecorder singleton — now shuffle-proofed.
