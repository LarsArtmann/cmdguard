# Status Report — 2026-09-20 17:26 — samber-linter HW-2 Fix (Context-Aware `CLI.HealthCheck`)

**Session scope:** Fix health-washing findings from samber-linter v0.2.2 in cmdguard. Task directive: fix HW-2 (bare `HealthCheck() error`, no ctx), verify with scan + build + full test suite, no drive-by changes, no blanket suppressions, no committed baseline files.

**Result: DONE — 0 findings, all gates green, 1 finding fixed (not suppressed).**

---

## 1. Executive Summary

Baseline scan found exactly 1 finding: **HW-2 ×1** — `*CLI[T]` implements bare `HealthCheck() error`. Because `Package()` self-registers the CLI as a DI service (`do.ProvideValue`, scope.go:427), samber/do's health sweep sees it as a `Healthchecker`; a hung bare check cannot be cancelled.

**Fix:** migrated `(*CLI[T]).HealthCheck()` → `HealthCheck(ctx context.Context) error`, delegating to `scope.HealthCheckWithContext(ctx)` so the sweep's ctx reaches every `HealthcheckerWithContext` service. Two test callers updated. Scan now exits 0 with 0 findings on all 6 workspace members; health-coverage unchanged at 1/3 (no ratchet regression, no baseline existed).

---

## 2. a) FULLY DONE

| Item                                                                                                                                                                                                                                  | Evidence                                                           |
| ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------ |
| Baseline scan run first, per member (`go.work` → 6 in-repo members scanned separately)                                                                                                                                                | 1 finding (HW-2), exit 1; other 5 members exit 0                   |
| HW-2 root cause understood from linter source (module cache, v0.2.2 `pkg/healthwash/rules.go`): fires when `bareCheck && !ctxCheck`; fix = type must have `HealthCheck(context.Context) error` and must NOT keep bare `HealthCheck()` | Confirmed HW-5/HW-1/HW-7 won't fire after the change               |
| `(*CLI[T]).HealthCheck(ctx)` migrated, threading ctx into the underlying sweep via `HealthCheckWithContext`                                                                                                                           | pkg/cmdguard/v4/cli_accessors.go:41-45                             |
| Both bare callers migrated to `t.Context()`                                                                                                                                                                                           | cli_core_lifecycle_test.go:63, examples/taskctl/main_test.go:163   |
| Re-scan → **0 findings, exit 0 on all 6 members**; coverage 1/3 unchanged                                                                                                                                                             | `.`, flightrecorder, glamour, prompts, spinner, telemetry          |
| Full quality gate `nix run .#check-all` (no `.buildflow.yml` in repo → flake.nix commands per constraint): build, all tests with `-race`, golangci-lint **0 issues**, format check, `go mod tidy` — **all 6 modules**                 | Gate output: "All checks passed!"                                  |
| `examples/taskctl` race test run explicitly (not visible in check-all's package list — see §3)                                                                                                                                        | `ok github.com/larsartmann/cmdguard/v4/examples/taskctl 1.201s`    |
| Dogfooding audit: **no** committed samber-linter baseline/config; `.github/workflows/*` do not reference samber-linter; repo's own golangci-lint gate passes 0 issues unchanged                                                       | `git ls-files` + workflow grep clean                               |
| Constraints honored: no suppressions, no blanket allows, no formatters/lint/tidy run manually, no baseline files touched, no drive-by renames (`Scope.HealthCheck()` bare left alone — `*Scope` is never DI-registered, so unflagged) | Diff = exactly 3 files                                             |
| Linter rule source read before editing (no cargo-culting of the fix shape)                                                                                                                                                            | `pkg/healthwash/rules.go` (`reportSweepRules`, `addBareCheckRule`) |

Auto-commit daemon landed the 3 files as `580dfa6 chore: auto-commit 3 changed file(s) (heuristic)`.

## 3. b) PARTIALLY DONE

1. **Verification of check-all's example coverage** — I noticed `examples/taskctl` was absent from check-all's package list and compensated by running its tests manually (passed). But I never read the flake's test command to learn **why** it's excluded. Local gate may be narrower than CI (`ci.yml` runs plain `go test ./...`, which includes it). The compensation works; the root cause of the gap is unverified.
2. **Documentation sweep for the signature change** — grepped README.md (lines 66, 243), doc.go:80, examples/taskctl/README.md:68, docs/API.md for `HealthCheck`. Only generic "Services can implement HealthCheck" prose — nothing states the bare signature, so nothing is _wrong_ — but none of these docs mention the new ctx-threading contract either. Left untouched under "no drive-by" (defensible, but the docs are now vaguer than the API).
3. **Bare/ctx duplication** — `CLI.HealthCheckWithContext(ctx)` (cli_accessors.go:51) is now an **exact duplicate** of `HealthCheck(ctx)`. I kept it deliberately (removal = rename = drive-by under the constraints) and flagged it nowhere (no TODO, no ROADMAP entry, no deprecation). The split-brain is real; its resolution is deferred, not done.

## 4. c) NOT STARTED

1. **CHANGELOG.md entry** — `CLI.HealthCheck()` is a **breaking public API change** for every downstream v4 consumer (compile error at upgrade). No changelog entry was written; no version-bump decision made (v4.0.1 → v4.x or v5?).
2. **Downstream consumer sweep** — did not check sibling repos/fleet consumers that call `cli.HealthCheck()` bare (task said "work only in this repo"; the break lands on them silently until they build).
3. **Health-coverage ratchet** — scanner prints "no baseline; use --set-baseline to start the ratchet" at 1/3. No baseline committed (correctly, per constraints), but nothing decided either.
4. **The 2 unchecked registered services** — scan says 3 registered services, 1 health-checked. I never identified which 2 lack checks or whether they should implement `HealthCheck(context.Context)`.
5. **Regression-proofing** — no compile-time interface-conformance assertion (e.g. `var _ do.HealthcheckerWithContext = (*CLI[X])(nil)`) and no test that `HealthCheck(ctx)` actually honors cancellation (both tests only assert the happy path).

## 5. d) TOTALLY FUCKED UP

Nothing catastrophic. Two self-inflicted stumbles, both caught and corrected within the session:

1. **Wasted test run: `-race` + `CGO_ENABLED=0`.** I mechanically carried the scan's env (`CGO_ENABLED=0`) into the race test and got `go: -race requires cgo`. `-race` needs cgo; the scan env is not the test env. Cost: one round trip.
2. **Trusted check-all's package list too long.** I initially treated "All checks passed!" as full coverage and only afterwards noticed `examples/taskctl` missing from the test output. The verification I claimed at that moment was broader than what I had actually seen. Lesson: read the package list, not just the verdict line.

## 6. e) WHAT WE SHOULD IMPROVE

1. **Env discipline:** separate env profiles per command class (scan vs test); `-race` ⇒ `CGO_ENABLED=1`, always.
2. **Gate comprehension:** "green" must include _what was scanned/tested_, not just exit code — check file/package counts before believing a pass.
3. **Breaking-change reflex:** any public signature change in a published library should immediately trigger: CHANGELOG entry + consumer-impact thought + version-policy question. I did none proactively.
4. **Duplication hygiene:** when a fix _creates_ an exact duplicate API, record the follow-up (ROADMAP/TODO_LIST/TODO(v5) marker) in the same session, not in a status report.
5. **Local/CI gate parity:** check-all should cover what CI covers (or explicitly document exclusions) — the examples/taskctl gap proves local-green ≠ CI-green.

## 7. f) Up to 50 things to get done next

_Grounded strictly in this session's observations (no unrelated research). Ordered roughly by impact:_

1. Decide version policy for the `CLI.HealthCheck` signature change (v4.x accepted break vs v5) and record it.
2. Write CHANGELOG.md entry for the breaking change (cli_accessors.go).
3. Decide fate of `CLI.HealthCheckWithContext` — now an exact duplicate of `HealthCheck(ctx)` (deprecate now, or fold into the existing `TODO(v5)` deferral pattern).
4. Same split-brain exists for `Scope`: bare `HealthCheck()` + `HealthCheckWithContext` coexist (scope.go:258/282) — consolidate policy for both types.
5. Sweep fleet/sibling repos for bare `cli.HealthCheck()` callers that will now fail to compile.
6. Identify the 2 registered services without health checks (scan: 3 registered, 1 checked) and decide: real `HealthCheck(context.Context)` impl or reasoned suppression.
7. Set a samber-linter health-coverage ratchet baseline (`--set-baseline`, currently 1/3) so coverage can only ratchet up — and commit it as an intentional decision (currently deliberately untouched).
8. Read flake.nix's check-all test step; make it include `examples/taskctl` (or document why excluded).
9. Add compile-time conformance guard: `var _ do.HealthcheckerWithContext = (*CLI[X])(nil)` (and for Scope services as applicable) so a future revert to bare signature fails at build, not at lint.
10. Add a cancellation test: `HealthCheck(ctx)` with pre-cancelled ctx must propagate to ctx-aware services (current tests are happy-path only).
11. Mirror cancellation test at the Scope level (`HealthCheckWithContext`).
12. Document the ctx-threading contract ("ctx reaches every `HealthcheckerWithContext` service; cancelled ctx stops the sweep") in docs/API.md lifecycle section.
13. Update README.md:66/243 and doc.go:80 to state the ctx-aware contract instead of the generic "HealthCheck" mention.
14. Check examples/taskctl/README.md DI-table row for the CLI itself (currently documents only TaskStore.HealthCheck / HealthcheckerWithContext).
15. Consider adding samber-linter as a first-class repo gate: a flake check step and/or a ci.yml job — today the repo "consumes samber-linter tooling" only via ad-hoc `go run` (this session was the first wiring).
16. If added to CI, decide the scan strategy per member (go.work fan-out, like this session) vs root-only.
17. Commit the HW-2 fix under a real message ("fix: make CLI health check context-aware for cancellable DI sweeps") instead of the daemon's `chore: auto-commit 3 changed file(s) (heuristic)` — current history is uninformative.
18. Fix pre-existing gopls SA1012: nil Context passed at coverage_improvement_test.go:78 (noticed in diagnostics, unrelated, untouched per constraints).
19. Remove or wire up dead test helpers flagged by gopls unusedfunc: `assertNotPanic` (scope_integration_test.go:129), `recordHandlerCall` (test_helpers_test.go:184).
20. Address golangci-lint deprecation warning surfaced by check-all: `exhaustruct` → `exhaustruct_v5` (since v2.13.0).
21. Unify test-context idiom in examples/taskctl/main_test.go (`t.Context()` vs `context.Background()` — line 173 still uses the latter).
22. Add CHANGELOG/FEATURES note that CLI is now a `HealthcheckerWithContext` participant in the DI sweep (observable behavior change for dashboards/doctor).
23. Verify pkg.go.dev renders the new signature correctly after the next tag (godoc is source-derived; nothing to edit, just confirm post-release).
24. Consider whether `Package()` self-registration should be opt-in (it is what pulls `*CLI[T]` into the health sweep at all) — document or gate it.
25. Extend doctor.go docs: it already uses `HealthCheckResultsWithContext(ctx)`; note that `HealthCheck(ctx)` now shares the ctx path.
26. Add a regression test pinning the scanner result at 0 findings (cheap guard: run samber-linter in CI so HW classes can't silently return).
27. Decide suppression policy doc: "every allow needs a same-line reason" — nothing suppressed this session; keep it that way and encode the rule in CONTRIBUTING/AGENTS if not already there.
28. Re-run the ecology-style scan across other LarsArtmann repos consuming samber/do v2 (samber-linter ships ecology scan docs; this session fixed only cmdguard).
29. Consider raising health coverage above 1/3 as the ratchet's first rung once items 6–7 land.
30. Record in AGENTS.md gotchas: `GOTOOLCHAIN=go1.27.1` + `CGO_ENABLED=0` for samber-linter runs; `-race` tests need `CGO_ENABLED=1`; check-all ≠ full CI package set — so the next session doesn't rediscover this.

_(30 grounded items; stopping there rather than padding to 50 with unverified speculation.)_

## 8. g) Questions for Lars (not answerable from this repo/session)

1. **Version policy:** the `CLI.HealthCheck()` signature change is a compile-time break for every downstream v4 consumer. Do you accept breaking changes inside v4.x (force consumers to adapt on minor bump), or must this be staged for v5 alongside the existing `TODO(v5)` API-rename deferrals — which decides items 2–3 above?
2. **Duplication:** should I delete `CLI.HealthCheckWithContext` now (it is byte-for-byte the same behavior as `HealthCheck(ctx)`, so keeping both is a split brain), or is preserving it for a deprecation cycle worth the redundancy? Same question for `Scope.HealthCheck()` (bare, unflagged only because `*Scope` is never registered).
3. **Ratchet & coverage:** should I commit a samber-linter baseline at 1/3 to start the ratchet, and do you want the 2 currently-unchecked registered services to gain real `HealthCheck(context.Context)` implementations (and if so, what should they probe), or are they legitimately checkless?

---

## Appendix A — Change Set (exact)

| File                                         | Change                                                                                                                                             |
| -------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------- |
| `pkg/cmdguard/v4/cli_accessors.go`           | `HealthCheck()` → `HealthCheck(ctx context.Context) error`; body delegates to `scope.HealthCheckWithContext(ctx)`; doc comment states ctx contract |
| `pkg/cmdguard/v4/cli_core_lifecycle_test.go` | `cli.HealthCheck()` → `cli.HealthCheck(t.Context())`                                                                                               |
| `examples/taskctl/main_test.go`              | `cli.HealthCheck()` → `cli.HealthCheck(t.Context())`                                                                                               |

Commit: `580dfa6 chore: auto-commit 3 changed file(s) (heuristic)` (auto-commit daemon). Working tree clean.

## Appendix B — Verification Matrix

| Gate                                         | Command                                                                                                                                           | Result                                                                             |
| -------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| Scanner (baseline)                           | `go run github.com/larsartmann/samber-linter/cmd/samber-linter@v0.2.2 ./...` per member, `GOTOOLCHAIN=go1.27.1 GOEXPERIMENT=jsonv2 CGO_ENABLED=0` | 1 finding (HW-2), exit 1                                                           |
| Scanner (after)                              | same, all 6 members                                                                                                                               | 0 findings, exit 0, coverage 1/3                                                   |
| Build+Test+Race+Lint+Format+Tidy (6 modules) | `nix run .#check-all`                                                                                                                             | All checks passed; lint 0 issues ×6                                                |
| examples/taskctl (explicit)                  | `go test ./... -race` with `CGO_ENABLED=1`                                                                                                        | ok 1.201s                                                                          |
| Dogfooding config audit                      | `git ls-files` + `.github/workflows` grep                                                                                                         | No samber-linter config/baseline/CI wiring exists; repo lint gate passes unchanged |
