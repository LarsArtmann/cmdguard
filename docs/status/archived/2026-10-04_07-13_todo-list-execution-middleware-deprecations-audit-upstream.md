# Status Report: TODO-List Execution — Context Middleware, Deprecations, Audit Upstream, Coverage

> **Live remainder:** none — every item carries an inline verdict; the v4-era work was retired at v4.3.1 and superseded by v5 (CHANGELOG.md).

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

<!-- check-rows:ledger -->

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

~~1. **F1 cmdguard side (BLOCKED, by design):** `AuditMiddleware[T]` is NOT wired.~~ done — F1 landed in v5: AuditMiddleware+WithAuditMiddleware via Plugin.RecordCommand (pkg/cmdguard/v5/auditlog.go:232,259)
It cannot ship until samber-do-auditlog is pushed + tagged ≥ v0.11.0 — a
committed call to `Plugin.RecordCommand` against go.mod's v0.10.0 would break
every non-local build (workspace `replace` doesn't travel). TODO_LIST carries
the exact unblock steps.
~~2. **CHANGELOG.md (cmdguard):** my Unreleased entries (Added/Deprecated/Changed +~~ done — CHANGELOG re-applied: [4.0.2] backfill at CHANGELOG.md:507; [4.1.0] 2026-10-04 section at :455
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

~~1. **CHANGELOG entries clobbered twice** (see b). I verified doc numbers late —~~ done — resolved later: [4.0.2] backfill + [4.1.0] section present (CHANGELOG.md:507,455)
caught it only because the coverage figure wasn't where I expected. Lesson:
with a concurrent session active, re-verify my doc edits after daemon commits.
~~2. **Flightrecorder singleton leak broke a test under shuffled order:** my first~~ done — taskctl wires explicit WithFlightRecorderRecorder lifecycle (examples/taskctl/main.go:133)
taskctl approach used `WithFlightRecorder` (lazy start, never stops). BuildFlow's
shuffled test-race run made `TestFR_Integration_CaptureOnCommandError` fail
(`rec.Start()` → "already enabled"). Fixed by explicit recorder lifecycle
(per-test `t.Cleanup(rec.Stop)`; `run()` stops the production recorder). My
local runs had passed because file order ran coverage_test.go first — I
initially shipped a test that was order-dependent and didn't know it.
~~3. **Coverage numbers churned in docs:** wrote 80.1% into three docs, refactor~~ **NOT-DO — historical incident note; numbers settled at 85.9%, now superseded by the CI badge (AGENTS.md:9).**
dropped real coverage to 76.9% → 78.5% → 85.9% final. Should have measured
once at the end; two docs carried a stale number briefly.
~~4. **`db status` in the store-failure table** — that command doesn't touch the~~ **NOT-DO — fixed in-session (db dropped from the table); no db store case in taskctl tests today.**
store; test failed once, fixed by dropping it. Read the handler before
assuming DI usage.
~~5. **Two edit-tool mid-air failures** (taskctl main.go tail, CHANGELOG §2) from~~ **NOT-DO — transient edit-tool mishap, recovered same session; no follow-up work.**
whitespace/context mismatch — recovered by re-reading; wasted a cycle each.
~~6. **`rg -rn` misuse** (replace flag) mangled one grep — harmless, caught.~~ **NOT-DO — transient grep mishap, caught immediately; no follow-up work.**
~~7. **BuildFlow full gate is red on environment, not code:** license-check~~ done — env gaps closed: go-licenses auto-defer (2026-10-07) + sandbox GOTOOLCHAIN=local pin (AGENTS.md:324,326)
(go-licenses not in PATH outside `nix develop`), nix-build-verify/treefmt
(sandbox DNS refused → go1.27.0 toolchain download fails). Pre-existing,
reproduced across 6+ runs before mine. Go build/test/race/lint steps green.
~~8. **lint/ sub-module is mid-refactor by a concurrent session** (working-tree~~ done — lint/ shipped: lint/v0.3.0 tag + CHANGELOG section (CHANGELOG.md:303)
changes I didn't author + a parse error in `lint/cmd/cmdguard-lint/main.go`

- their own status file written 2 min ago). Workspace-wide builds that fan
  into ./lint fail through no fault of this session's changes. Left untouched
  per "never revert changes you didn't author".

---

## e) WHAT WE SHOULD IMPROVE

~~1. TODO_LIST had 8 of 13 items stale — the harvest/verify loop runs too rarely;~~ done — harvest loop mechanized: docs-gates battery + TODO_LIST harvest sections (2026-11-08/11-11)
consider the docs-health drift alarm (verify-checklist reference) for this repo.
~~2. AGENTS.md project structure omits the `lint/` sub-module entirely (and says~~ done — AGENTS.md:98 structure tree documents the lint/ sub-module
"6 modules" while go.work now has 7 incl. ./lint + local go-output replaces
from the concurrent session). Structure section is drifting.
~~3. AGENTS.md exceeds the doctor's 377-line budget (386 pre-session → ~395 now~~ **Won't implement — 377-line doctor budget superseded; AGENTS.md is 456 lines after the v5 living-docs rebuild.**
despite my condensing). Needs a deliberate content split, not more trimming.
~~4. taskctl's `WithFlightRecorder` option (never stops the process-wide recorder)~~ done — Stop ownership documented: explicit Start/Stop pattern at flightrecorder/README.md:48-56
is a footgun for embedders/tests; consider documenting Stop ownership in the
sub-module README or providing an auto-stop variant.
~~5. CHANGELOG discipline: releases (v4.0.2) shipping without a CHANGELOG entry —~~ **Won't implement — never added; CONTRIBUTING.md has no CHANGELOG rules and v4.0.2-era discipline is moot post-v5.**
the 4.0.2 backfill I wrote from git evidence shouldn't have been needed.
~~6. BuildFlow env gaps (go-licenses PATH, offline nix toolchain download) make the~~ done — same fixes as §d.7: go-licenses auto-defer + GOTOOLCHAIN=local sandbox pin (AGENTS.md:324,326)
"full" gate permanently red locally — either fix the env or scope skip_steps
so the gate's signal is trustworthy.

---

## f) NEXT UP TO 50 (roughly ordered)

~~1. Re-apply CHANGELOG Unreleased entries + [4.0.2] backfill (content ready).~~ done — [4.0.2] backfill present (CHANGELOG.md:507) plus [4.1.0] section dated 2026-10-04 (CHANGELOG.md:455)
~~2. Re-run `buildflow --build-mode full` scoped to my modules after CHANGELOG fix.~~ **Won't implement — superseded: hard gates are now CI + `nix run .#check-all`; go-licenses blocker auto-deferred (AGENTS.md:326).**
~~3. Push + tag samber-do-auditlog v0.11.0 (needs Lars's go-ahead).~~ done — samber-do-auditlog v0.11.0 published (sibling repo tags v0.11.0 through v0.12.3)
~~4. After 3: bump cmdguard go.mod auditlog → v0.11.0.~~ done — cmdguard go.mod requires samber-do-auditlog v0.12.3 (go.mod:84)
~~5. After 4: implement `AuditMiddleware[T]` + `WithAuditMiddleware` + tests (F1 close).~~ done — v5 ships AuditMiddleware[T] + WithAuditMiddleware via Plugin.RecordCommand (pkg/cmdguard/v5/auditlog.go:232,259)
~~6. After 5: taskctl demo of AuditMiddleware; FEATURES row → FULLY_FUNCTIONAL.~~ done — taskctl wires WithAuditMiddleware (examples/taskctl/main.go:113); FEATURES.md:262 row FULLY_FUNCTIONAL
~~7. Migrate telemetry sub-module to ContextMiddleware (span ctx propagation).~~ **Won't implement — superseded: v5 Middleware[T] itself receives ctx (middleware.go:26); no ContextMiddleware migration needed.**
~~8. Add `telemetry.ContextTelemetryMiddleware` or fold into v0.3.0 of telemetry.~~ **Won't implement — ContextTelemetryMiddleware never added; moot once v5 Middleware signature carries ctx (telemetry.go:34).**
~~9. Consider spinner sub-module: nothing needed (no ctx), verify once.~~ **NOT-DO — no-op confirmed: spinner ships as plain middleware (spinner/v0.3.0); v5 signatures carry ctx.**
10. v5 prep: spec the Middleware/ContextMiddleware unification (breaking). **→ ROADMAP.md v6 candidate #8 — unification rejected for v5, open for v6**
~~11. v5 prep: `CommandInfo`→`CommandMetadata` mechanical-rename plan (T1).~~ done — v5 Middleware signature takes info CommandMetadata (pkg/cmdguard/v5/middleware.go:26)
~~12. v5 prep: `TypeHandler`→`TypeCodec` + deprecated alias plan (T2).~~ done — v5.0.0 ships TypeCodec/TypeCodecFunc (pkg/cmdguard/v5/type_codec.go:13; ROADMAP deferred #3 struck done)
~~13. v5 prep: `PromptRunner`→`HuhPrompter` (T3).~~ **Won't implement — v5 kept PromptRunner + SetPromptRunner (pkg/cmdguard/v5/prompts.go:15); HuhPrompter never introduced.**
~~14. v5: remove `SetConfig` (T4) — update doc comments referring to it.~~ done — SetConfig removed in v5; zero SetConfig symbols in pkg/cmdguard/v5 (grep clean)
~~15. v5: generic `RegisterInScope[T]` (T5) with variadic-provider design doc.~~ done — generic RegisterInScope[T any] ships in v5 (pkg/cmdguard/v5/scope.go:358)
~~16. v5: `Get[T]`→`GetService[T]` (candidate 6).~~ done — resolved differently: v5 has no Get[T]; getters are Invoke[T]/InvokeNamed[T] (ROADMAP candidate 6 struck MOOT)
17. v5: `Package[T]` redesign (candidate 7). **→ ROADMAP.md candidate #7 — Package[T](scope, ...) redesign still tracked open (scope.go:416)**
~~18. Document `lint/` sub-module in AGENTS.md structure + module count.~~ done — AGENTS.md:98 structure tree + sub-module table row (AGENTS.md:151) document lint/
~~19. Update AGENTS.md "6 modules" → current go.work truth (incl. go-output replaces).~~ done — AGENTS.md states 8 modules incl. lint (AGENTS.md:113,309); go.work confirms
~~20. AGENTS.md size diet to ≤377 lines (move detail to docs/).~~ **Won't implement — diet superseded by the v5 living-docs rebuild; AGENTS.md is now 456 lines.**
~~21. Investigate nix-build-verify 6/6 failures (offline toolchain fetch).~~ done — fixed: flake pins goPkg + GOTOOLCHAIN=local so sandboxed treefmt needs no network (AGENTS.md:324)
~~22. Install go-licenses in devShell/PATH so license-check can run locally.~~ done — superseded — BuildFlow auto-defers the go-licenses stdlib class since 2026-10-07 (AGENTS.md pre-commit note)
~~22. Decide skip_steps for sandbox-hostile nix checks (or make them hermetic).~~ done — made hermetic: goPkg ≥ go-directive rule + GOTOOLCHAIN=local documented (AGENTS.md:324)
~~23. taskctl: cover remaining buildCommands error branches via fault injection?~~ **Won't implement — v4-era ask; taskctl migrated to v5 and measures 93.9% coverage (AGENTS.md:70).**
~~24. taskctl: `db` parent + `env create --force` paths coverage.~~ **Won't implement — v4-era ask; taskctl migrated to v5 and measures 93.9% coverage (AGENTS.md:70).**
~~25. taskctl: prompt-tag path (`prompt:""` + WithPromptOnMissing) test with prompts stub.~~ **Won't implement — v4-era ask; taskctl migrated to v5 and measures 93.9% coverage (AGENTS.md:70).**
~~26. pkg/testutil: push 80.3% → 90% (failure_paths_test.go branch completion).~~ **NOT-DO — never landed: pkg/testutil measures 56.8% now (go test -cover, 2026-11-20); v4-era effort overtaken.**
27. Add fuzz corpus seeds for the 8 targets (ROADMAP deferred #7). **→ ROADMAP.md deferred #7 — fuzz corpus expansion still tracked open**
~~28. cmdguard-lint: once the concurrent refactor lands, add a rule for deprecated-API usage.~~ **Won't implement — no deprecated-API rule exists (cmdguard-lint rules CG001-CG004); v4 deprecated APIs removed in v5.**
~~29. docs/API.md: add ContextMiddleware/TimeoutMiddleware reference sections.~~ **Won't implement — docs/API.md no longer exists; middleware reference lives in FEATURES.md:228-232 and the website.**
~~30. README.md: feature list mention of context middleware.~~ **Won't implement — README never gained the mention; middleware docs live in FEATURES.md:228-232 and the website.**
~~31. docs/COBRA_FOOTGUNS.md: cross-link WithContextMiddleware for ctx threading.~~ **Won't implement — cross-link never added; COBRA_FOOTGUNS.md §3 covers signal-ctx only.**
~~32. ExportAuditLog: consider Writer-based test instead of read-only-dir (portability).~~ **Won't implement — taskctl kept file-based export; v5 has its own ExportAuditLog tests (cli_auditlog_test.go:261).**
~~33. Flightrecorder: document Stop-ownership semantics in README (embedders).~~ done — flightrecorder README:48-56 documents Stop ownership with defer rec.Stop pattern
~~34. Flightrecorder: optional auto-stop middleware variant (owns lifecycle).~~ **NOT-DO — auto-stop variant never built; README teaches explicit lifecycle instead (flightrecorder/README.md:48-56).**
~~36. Samber-do-auditlog: HTML viz "Command" color distinct from Registration (both --accent).~~ done — auditlog viz colors are centralized per-type lookup tables now (samber-do-auditlog daghtml_adapter.go status and provider maps, metadata_test.go event colors), so Command and Registration no longer share one accent in the graph layer
~~37. Samber-do-auditlog: consider dedicated `command_name` field over ServiceRef overload (schema 0.4.0).~~ **→ routed — ROADMAP Harvested, cmdguard-lint `command_name` field**
~~38. Samber-do-auditlog: cmdguard-lint-style rule or doc forbidding ServiceName pollution?~~ **→ ROADMAP Harvested (cmdguard-lint ServiceName-pollution rule)**
~~35. CHANGELOG entry for 4.0.2 should have existed at tag time — add release checklist item.~~ **Won't implement — checklist item never added; publish discipline moved to the lint/v0.2.1 release playbook (TODO_LIST F3).**
~~36. Benchmarks for ContextMiddleware overhead vs plain (benchmarks/ dir).~~ **NOT-DO — never landed; benchmarks/guard_bench_test.go covers CLI construct/execute, no middleware-overhead bench.**
~~37. Table-driven test consolidation in middleware_context_test.go (minor).~~ **Won't implement — v4 test file retired; v5 middleware_context_test.go (378 lines) is the current suite.**
~~38. TODO(v5) marker line numbers in ROADMAP drift already (I removed line pins — good); re-pin? no — keep unpinned, note convention.~~ **NOT-DO — self-answered in-item: keep markers unpinned; v4 retirement (v4.3.1) moots v4 marker bookkeeping.**
39. git-town/library-policy verification (ROADMAP deferred #10). **→ ROADMAP.md deferred #10 — git-town.toml + library-policy.yaml verification still tracked open**
~~40. Re-audit window per ROADMAP #12: after F1 lands.~~ done — window opened: v5 shipped 2026-11-08; audits ran (docs/status/2026-11-08_08-32, 2026-11-20_12-29)
~~41. Consider `WithMiddlewareBoth` ergonomic if users complain about two lists.~~ **NOT-DO — WithMiddlewareBoth never added; two lists kept deliberately, unification is ROADMAP v6 candidate #8.**
~~42. Verify fang + context middleware interplay (error display unaffected) — test exists indirectly via ExecuteWithArgs; add explicit case.~~ done — explicit fang coverage exists in v5 (cli_fang_on_test.go; WithFang cases in cli_core_new_test.go)
~~43. WithCleanup + ContextMiddleware ordering test (cleanup outside chain today).~~ **Won't implement — ordering test never added; v5 covers cleanup firing (cli_cleanup_test.go); v4 chain retired.**
~~44. `TimeoutMiddleware` jitter/backoff variant? (YAGNI — only on demand.)~~ **NOT-DO — YAGNI held: no jitter/backoff TimeoutMiddleware variant in v5 (middleware.go:136).**
~~45. Status-report this file → TODO_LIST harvest of items 1-23.~~ **Won't implement — superseded: TODO_LIST rebuilt for the v5 era (harvests 2026-11-08/11-11); 2026-10-04 items moot.**
~~46. Schedule docs-health AUDIT after F1 + v5 prep docs settle.~~ done — ran: docs-health sweep all-gates-green 2026-11-08 (docs/status/2026-11-08_02-29)

---

## g) QUESTIONS FOR LARS (cannot self-answer)

~~1. **Push permission:** may I `git push` + tag `samber-do-auditlog` v0.11.0~~ done — resolved: auditlog v0.11.0+ published (tags through v0.12.3); cmdguard consumes v0.12.3 (go.mod:84)
(and later cmdguard), or do you want to review the upstream API first
(`RecordCommand` overloads `ServiceRef.ServiceName` with the command name —
the alternative is a schema-0.4.0 `command_name` field, item 37)?
~~2. **Concurrent session on lint/:** another agent is mid-refactor of the~~ done — resolved: lint/ shipped (tags lint/v0.2.1, lint/v0.3.0; CHANGELOG.md:303)
`lint/` sub-module (working-tree changes + a parse error right now). Leave
it alone, or should I take over/finish it once you confirm it's abandoned?
~~3. **CHANGELOG collision:** my Unreleased entries were overwritten once by that~~ done — resolved: entries re-applied — [4.0.2] backfill present (CHANGELOG.md:507)
concurrent session. Re-apply immediately (risking another collision while
it's active), or wait for your go-ahead / its completion?

---

**Bottom line:** every actionable TODO item is either closed-and-verified or
blocked with the blocker named. The only session regression is the
twice-clobbered cmdguard CHANGELOG (content ready, 2-minute re-apply), and the
one engineering lesson worth keeping: test-order dependence around the
flightrecorder singleton — now shuffle-proofed.
