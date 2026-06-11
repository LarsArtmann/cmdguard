# Post-Sprint Threshold-21 Status Report

**Date:** 2026-06-11 07:54 CEST
**Sprint:** art-dupl threshold-21 inspection (final round of diminishing-returns triage)
**Previous Status:** `2026-06-11_05-39_post-dedup-helper-extraction.md` (Sprint 4, 13 → 12 groups at t=22)
**Branch:** `master` (in sync with `origin/master` at `ce54ad8`)
**Status:** ✅ **COMPLETE** — 1 group eliminated, 18 triaged as idiomatic at t=21; 12 at t=22

---

## Executive Summary

Lowered `art-dupl` threshold to 21 to surface every remaining structural duplication. Inspected all 19
clone groups individually. 18 are intentional Go idioms (handler signatures, middleware signatures,
Cobra completion signatures, accessors, table data rows, cross-package helper boundaries, exit-on-error
fprintf+Exit pattern) — accepted per the `deduplicate-code` skill's "zero HARMFUL duplication" principle.
One group (G16) was an assertion pattern that yielded to `testutil.AssertFieldEqString`. Net **4 lines
removed** from `configload/loader_test.go`, zero behavior change, zero lint regressions.

The campaign's diminishing-returns phase is now clear: the 12 remaining groups at t=22 are all
fundamentally untouchable (would require breaking idiomatic Go signatures, splitting per-file helper
shims, or merging data-distinct table rows). Further work should pivot to **architectural** targets
(`output.go` 17-clone format switch is the next high-leverage candidate) or **feature** work.

---

## a) FULLY DONE

### Sprint 5 (Just Completed)

| # | Refactor | File | Net Lines |
|---|----------|------|-----------|
| 1 | `if cfg.Name != tt.expect { t.Errorf(...) }` → `testutil.AssertFieldEqString(t, cfg.Name, tt.expect, "name")` (2 sites) | `pkg/cmdguard/v2/configload/loader_test.go` | -4 |
| | **Total** | 1 file modified | **-4 net** |

Commit: `ce54ad8 refactor,test: replace cfg.Name field checks with AssertFieldEqString, eliminate 1 clone group`

### Cumulative Campaign (5 Sprints)

| Sprint | Theme | Groups Eliminated | Net Lines | Date |
|--------|-------|-------------------|-----------|------|
| 1 | Test fixture consolidation | 159 → ~80 | -74 | 2026-06-10 |
| 2 | Shared test infrastructure | ~80 → 30 | -106 | 2026-06-10 |
| 3 | Cross-package helper dedup | 30 → 14 | -11 | 2026-06-10 |
| 4 | Sprint 4 helpers (`AssertOutputContains`, `AssertStringSlicesEqual`, `writeJSONConfigFile`) | 14 → 13 (at t=22) | +14 | 2026-06-11 |
| 5 | Threshold-21 triage | 13 → 12 (at t=22) | -4 | 2026-06-11 |
| **Total** | | **159 → 12 at t=22 (92.5% reduction)** | **-181 net** | |

### Quality Gates (Post-Sprint)

| Gate | Status | Details |
|------|--------|---------|
| `go build ./...` | ✅ PASS | Clean compile, zero errors |
| `go test ./... -race -count=1 -timeout 180s` | ✅ PASS | All packages green (385+ tests) |
| `golangci-lint run ./...` | ✅ PASS | **0 issues** |
| `golangci-lint fmt ./...` | ✅ PASS | No formatting changes needed |
| Coverage (v2 package) | ✅ 85.0% | Up from 84.8% (Sprint 4 baseline) |
| Coverage (configload) | ✅ 90.2% | Up from baseline |
| Coverage (testutil) | ✅ 55.2% | Helper-internal, low impact |
| Working tree | ✅ Clean | No uncommitted changes |
| Sync with origin | ✅ In sync | `ce54ad8 == origin/master` |

### Helper Inventory (Test Infrastructure)

Helpers available in `pkg/testutil` AND `pkg/cmdguard/v2/testutil` (cross-package mirror per Sprint 3 pattern):

| Helper | Purpose | Used At |
|--------|---------|---------|
| `AssertNoError` | `t.Fatalf` if err is non-nil | Most test files |
| `AssertErrorIs` | `errors.Is` chain check | Error wrapping tests |
| `AssertErrorContains` | Substring check on err.Error() | Error message tests |
| `AssertStderrContains` | Capture stderr and check substring | CLI output tests |
| `AssertFieldEq` | Generic equality | Field assertion tests |
| `AssertFieldEqString` | String equality with `%q` | Field assertion tests |
| `AssertFieldLen` | Slice length check | Slice tests |
| `AssertBoolTrue` / `AssertBoolFalse` | Boolean field checks | Flag tests |
| `AssertOutputContains` | Substring on captured output | Output buffer tests |
| `AssertStringSlicesEqual` | Length + per-element string slice equality | Order/middleware tests |
| `AssertEqual` / `AssertEqualf` | Generic equality with optional message | General |
| `AssertNotEqual` / `AssertNil` / `AssertNotNil` | Negation helpers | General |
| `AddCommand` | Fatal wrapper around `v2.AddCommand` | 40+ call sites |

File-local helpers extracted:
- `pkg/cmdguard/v2/test_helpers_test.go`: `registerFlags`
- `pkg/cmdguard/v2/testhelpers_test.go`: `newTestCLI`, `newTestCLIWithAuditLog`, `addCommand`, `addTestPlugin`
- `pkg/cmdguard/v2/cli_exec_test.go`: `newTestCmd`
- `pkg/cmdguard/v2/cli_superb_test.go`: `goodCommand`
- `pkg/cmdguard/v2/middleware_test.go`: `captureInfoMiddleware`, `captureNameMiddleware`, `beforeAfterMiddleware`
- `pkg/cmdguard/v2/config_file_test.go`: `writeJSONConfigFile`
- `tests/integration/v2_bdd_lifecycle_test.go`: `newLifecycleCLI`, `newLifecycleStrictCLI`, `newLifecycleDraconianCLI`, `recordLifecycleStep`, `lifecycleErrHandler`, `registerCommand`

---

## b) PARTIALLY DONE

Nothing currently in flight. The deduplication campaign is at a stable endpoint — all 12 remaining
groups at t=22 are individually triaged and accepted as idiomatic.

---

## c) NOT STARTED

### High-Impact Architectural Refactor Candidates

| # | Target | Duplication | Estimated Impact |
|---|--------|-------------|------------------|
| 1 | `pkg/cmdguard/v2/output.go` | 17-clone format-specific switch | Strategy pattern with `Formatter` interface, one impl per format registered in a map. Would eliminate 17+ near-identical switch cases. Medium-high risk, large refactor. |
| 2 | `pkg/cmdguard/v2/command_options.go` | Many similar `WithXxx` functional options | Could extract a generic `withString[T,F any](field *string) CommandOption[T,F]` helper. Currently 19 options follow the pattern manually. Medium risk. |

### Architectural / Documentation

| # | Item | Notes |
|---|------|-------|
| 3 | `docs/adr/002-formatter-strategy.md` | Document the output strategy pattern decision (if pursued) |
| 4 | `docs/DOMAIN_LANGUAGE.md` | Glossary doesn't exist yet — would help AI sessions interpret the v2 type system |
| 5 | `examples/` directory | Only 1 example (`taskctl`). Could add minimal/single-command examples for each major feature |

### Test Coverage

| # | Item | Current | Notes |
|---|------|---------|-------|
| 6 | `pkg/cmdguard/v2/testutil` coverage | 55.2% | Helpers are tested implicitly by callers; explicit unit tests would be nice |
| 7 | `pkg/testutil` coverage | 0.0% (no test files) | No `_test.go` exists; package only consumed by v2 internal tests |
| 8 | `examples/taskctl` coverage | 70.5% | Has tests; could push toward 85% to match v2 |

### Infrastructure

| # | Item | Notes |
|---|------|-------|
| 9 | `nix flake check` (full) | Currently has devShell + formatter + format check; could add `buildGoModule`, `go vet`, tests |
| 10 | CI workflow (`.github/workflows/`) | No CI observed; only local nix shell |
| 11 | Pre-commit hook `gomod-check` | Pre-existing failures documented in AGENTS.md; bypassed via `--no-verify` |
| 12 | `go mod tidy` | Stale `go.sum` entries (23 reported in AGENTS.md) — could clean up |

---

## d) TOTALLY FUCKED UP

**Nothing.** Working tree is clean, all commits are well-formed, all tests pass, lint is clean, no
panics in library code, no broken imports, no missing types. The only "issues" are non-blockers
documented in c) above (stale go.sum, missing CI, etc.).

### Minor Concerns (Not Blocking)

- **`gopls infertypeargs` info diagnostics** (~330 across the project) — Go 1.26 redundant type
  argument warnings. These are info-level, not errors. Would be a sweep task: `gofmt -r` or
  `goimports -d` to clean up. Not harmful.
- **Pre-commit hook `gomod-check`** — AGENTS.md documents `--no-verify` as established practice. Not
  a regression; pre-existing condition.
- **Cross-package helper duplication** (`pkg/testutil` vs `pkg/cmdguard/v2/testutil`) — intentional
  package boundary pattern (pkg/testutil imports v2, v2 cannot import it back). Mirrored ~5 helpers.
  Accepted trade-off.

---

## e) WHAT WE SHOULD IMPROVE

### Process

1. **Stop lowering the art-dupl threshold.** The 12 remaining groups at t=22 are idiomatic. Further
   lowering yields only false positives (signature fragments) or patterns that are correctly
   repeated. Diminishing returns reached.
2. **Pivot to architectural refactors.** `output.go` 17-clone format switch is the next high-leverage
   target. A `Formatter` interface + registry would replace 17+ switch cases with one map lookup per
   format invocation. This is a real, valuable refactor — not a stylistic cleanup.
3. **Add CI.** No automated quality gate exists. A `.github/workflows/ci.yml` running
   `go test -race -cover` + `golangci-lint` + `nix flake check` would prevent regressions. Today the
   only safety net is the developer's local environment.
4. **Document the helper inventory.** There's no `docs/TEST_HELPERS.md`. Future maintainers (and AI
   sessions) would benefit from a one-page reference of every helper, its signature, and example
   usage. The data is scattered across 7+ files today.

### Code Quality

5. **Generic options pattern in `command_options.go`.** 19 `WithXxx` options follow the pattern
   `func WithXxx[T, F any](val T) CommandOption[T, F] { return func(c *Command[T, F]) { c.field = val } }`.
   A generic `Setter[T, F, V any]` could eliminate boilerplate. Worth measuring LOC reduction before
   committing.
6. **Split `output.go`.** 17 format cases in one file is 600+ lines. Splitting by format group
   (text formats, structured formats, graph formats) would improve navigability even if no strategy
   pattern is introduced.
7. **Add `WithXxxTest` helpers for command options.** The `flag.Value.String() != "X"` pattern that
   Sprint 4 extracted in `prompts_test.go` likely has analogues in other test files at t=18-19 that
   were not surfaced at t=21.

### Repository Hygiene

8. **`nix develop` vs `go test` parity.** Verify the nix devShell produces the same test output as
   the bare Go toolchain. Drift here would cause "works on my machine" issues.
9. **`go mod tidy` cleanup.** 23 stale `go.sum` entries per AGENTS.md. Routine hygiene.
10. **Commit message discipline.** Recent commits are exemplary (`refactor,test: <verb> <object>,
    <impact>`). Maintain this — it's the project's most visible quality signal.

---

## f) Top #25 Things To Do Next

Priority-ordered by impact / effort ratio. Effort: L=Large, M=Medium, S=Small. Impact: H=High, M=Medium, L=Low.

| # | Task | Effort | Impact | Notes |
|---|------|--------|--------|-------|
| 1 | `output.go` Formatter strategy refactor | L | H | 17-clone format switch → interface + registry. Real architectural win. |
| 2 | Add GitHub Actions CI (test + lint + format check) | M | H | First automated quality gate. |
| 3 | Generic `Setter`/`WithXxx` helper for `command_options.go` | M | M | 19 options follow same pattern; could halve LOC. |
| 4 | `docs/TEST_HELPERS.md` — inventory of all testutil helpers | S | M | Discoverability for future maintainers. |
| 5 | Clean stale `go.sum` entries (`go mod tidy`) | S | M | Pre-existing condition. |
| 6 | Test coverage: `pkg/cmdguard/v2/testutil` (55.2% → 80%+) | M | M | Helpers are tested implicitly; explicit unit tests would be nicer. |
| 7 | Split `output.go` by format group | M | L | 600+ line file → 3-4 thematic files. Even without strategy pattern. |
| 8 | `docs/DOMAIN_LANGUAGE.md` — DDD glossary | M | M | Bounded contexts, ubiquitous language, value objects. |
| 9 | `docs/adr/002-formatter-strategy.md` (if #1 pursued) | S | M | Document decision. |
| 10 | Test coverage: `examples/taskctl` (70.5% → 85%) | M | L | Example coverage is the user's first impression. |
| 11 | Extract `WithXxxTest` test helpers (t=18-19 patterns) | M | L | Sprint 4 covered some; more likely exist below t=21. |
| 12 | Extend `nix flake check` with `buildGoModule` + `go vet` | M | M | Stricter local gate before pushing. |
| 13 | Add `pkg/testutil` unit tests (currently 0% coverage) | S | L | Trivial assertions on trivial helpers. |
| 14 | Generic `withString[T,F]` refactor in `command_options.go` | S | L | Subset of #3. |
| 15 | Improve pre-commit hook to skip `gomod-check` cleanly | S | L | Remove `--no-verify` workaround. |
| 16 | Sweep `gopls infertypeargs` info diagnostics (~330) | M | L | Mechanical, but low-priority cleanup. |
| 17 | Add second minimal `examples/` command | M | L | Show minimum surface area; `taskctl` is large. |
| 18 | Document the cross-package helper mirroring pattern in AGENTS.md | S | M | Save future sessions the discovery cost. |
| 19 | `go test -coverprofile` + badge in README | S | M | Visible quality signal. |
| 20 | Investigate `nix fmt` vs `golangci-lint fmt` parity | S | L | Ensure they don't drift. |
| 21 | Add `WithGroup`/`WithFang`/`WithGlamour` to testutil helper for `NewCLI[testConfig]` (G7/G14 accept pattern) | S | L | Test boilerplate redux. |
| 22 | Benchmark: add `OutputTable` vs `OutputResult` perf data | M | L | Required for any future optimization claim. |
| 23 | Consider publishing v2.5.1 patch (just helpers + cleanup, no behavior change) | S | M | Make the dedup work visible to consumers. |
| 24 | Add `v2_test` package vs internal package split decision doc | S | M | Open question from Sprint 4 summary. |
| 25 | Add architecture diagram (D2) for v2 internal flow | M | M | Visual reference for new contributors. |

---

## g) Top #1 Question I Cannot Figure Out

**Is the `output.go` 17-clone format switch worth refactoring to a `Formatter` strategy, or is the
current switch a feature (not a bug) because each format's render logic is genuinely different and
the switch is the simplest possible dispatch?**

I can see two valid arguments:

- **Extract:** 17 cases that all do `case "json": return marshalAndWrite(...)` etc. is structurally
  similar. A `map[string]Formatter` would let users register new formats without modifying the core
  switch. The current code violates open-closed: adding a format requires editing `output.go`.

- **Accept:** Each format's render code is genuinely distinct (JSON needs encoding, Table needs
  column layout, D2 needs graph syntax, Markdown needs specific transforms). A strategy pattern would
  move the switch's complexity into N strategy implementations + 1 dispatch function. Net LOC may
  not decrease; readability may suffer. And the format set is bounded by `go-output` v0.8.0's
  registry — users can't truly register new formats anyway without forking.

I have searched for prior decisions, scanned the `go-output` API surface, and read every comment in
`output.go`. The answer depends on whether `go-output` (or its `anyFormatRegistry` extension point) is
something the project plans to depend on, or whether `cmdguard`'s role is just to be a thin wrapper.
**I cannot determine this from the code alone — it requires a product/architecture call about the
library's long-term scope.**

---

## Appendix A: All 19 Clone Groups at t=21 (Triaged)

| # | Pattern | File(s) | Decision |
|---|---------|---------|----------|
| G1 (5×) | `func(ctx, cfg, flags) error { return nil }` | `coverage_test.go:44,154,166,178,190` | Accept (handler sig) |
| G2 (4×) | `func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil }` | `cli_superb_test.go:22,244,278`, `glamour_test.go:183` | Accept (handler sig) |
| G3 (3×) | `tag: FlagTag{...}` table rows | `type_handler_test.go:71,571,590` | Accept (data fixtures) |
| G4 (3×) | `addCommand` / `AddCommand` / `registerCommand` facades | `testhelpers_test.go:93`, `testutil.go:93`, `v2_bdd_lifecycle_test.go:128` | Accept (intentional cross-package façades) |
| G5 (3×) | `func(_ context.Context, _ *lifecycleConfig, _ v2.NoFlags) error { return nil }` | `v2_bdd_lifecycle_test.go:27,87,685` | Accept (handler sig) |
| G6 (3×) | `func(_ context.Context, _ *T, info CommandInfo, next func() error) error` | `middleware.go:70`, `middleware_test.go:57,67` | Accept (middleware sig) |
| G7 (3×) | `func (c Command[T, F]) Method() func(ctx context.Context, cfg *T, flags F) error { return c.field }` | `command.go:60,65,70` | Accept (accessor pattern) |
| G8 (2×) | `AssertStringSlicesEqual` (cross-package) | `testutil.go:113`, `panic_test_helpers.go:222` | Accept (Sprint 4 helper boundary) |
| G9 (2×) | `fmt.Fprintf(os.Stderr, "X: %v\n", err); os.Exit(1)` | `taskctl/main.go:82`, `cli_exec_test.go:141` | Accept (Go exit-on-error idiom) |
| G10 (2×) | `func(_ context.Context, _ *T, _ CommandInfo, next func() error) error` | `middleware_test.go:43`, `spinner.go:94` | Accept (middleware sig) |
| G11 (2×) | `func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective)` | `taskctl/commands.go:159`, `coverage_improvement_test.go:168` | Accept (Cobra completion sig) |
| G12 (2×) | `AssertOutputContains` (cross-package) | `testutil.go:103`, `panic_test_helpers.go:273` | Accept (Sprint 3 helper boundary) |
| G13 (2×) | `NewParentCommand[testConfig, NoFlags]("parent", ...)` | `cli_superb_test.go:248,314` | Accept (different sub-cmds) |
| G14 (2×) | `NewCLI[testConfig](..., WithGlamourHelpTheme("dark"), WithFang(false))` | `glamour_test.go:82,172` | Accept (different app names "test"/"testapp") |
| G15 (2×) | `v2.WithLong[AppConfig, *XFlags](`...`)` | `taskctl/commands.go:66,103` | Accept (data-distinct markdown) |
| **G16 (2×)** | `if cfg.Name != tt.expect { t.Errorf("expected name %q, got %q", ...) }` | `configload/loader_test.go:183,250` | **Extract → AssertFieldEqString** |
| G17 (2×) | `NewCLI[testConfig](..., WithGroup("X", "Y"), WithFang(false))` | `cli_groups_test.go:89,124` | Accept (different group names) |
| G18 (2×) | `return func(_ context.Context, _ *lifecycleConfig, _ v2.NoFlags) error` | `v2_bdd_lifecycle_test.go:72,84` | Accept (handler sig) |
| G19 (2×) | `name: "X default", tag: FlagTag{...}, expected: "..."` | `type_handler_test.go:285,295` | Accept (data fixtures) |

---

## Appendix B: Clone Group Counts by Threshold (Campaign Trajectory)

| Threshold | Groups | Trend |
|-----------|--------|-------|
| t=30 | 0 | (Sprint 2 milestone) |
| t=22 | 12 | Stable since Sprint 4 |
| t=21 | 18 | Was 19 before Sprint 5 |
| t=20 | 29 | Not yet processed |
| t=18 | 68 | Not yet processed |

The 12-group floor at t=22 represents the genuine idiomatic Go signature repetition (handler params,
middleware params, accessor returns, table data rows, helper mirroring) that cannot be removed
without breaking the language's conventions.

---

## Appendix C: Recent Commits (Last 15)

```
ce54ad8 refactor,test: replace cfg.Name field checks with AssertFieldEqString, eliminate 1 clone group
8618022 docs: improve table alignment in status reports and remove obsolete pre-commit hook
a2bdc48 refactor,test: extract AssertStringSlicesEqual and writeJSONConfigFile helpers, eliminate 3 clone groups
f580852 refactor,test: add AssertOutputContains helper, replace 20 hand-rolled assertions
9cc2a94 refactor,test: extract newLifecycleCLI helper, convert output subtests to table-driven, eliminate 2 clone groups
efab861 refactor,test: extract newTestCLI/newTestCLIWithAuditLog/registerFlags helpers, eliminate 3 high-impact clone groups
839bb1a refactor(test): extract newTestPlugin helper, eliminate 9-clone group
3e8f949 refactor(test): extract addCommand/AddCommand/registerCommand helpers, eliminate 20-clone group
16476ac refactor,test: extract shared test helpers, collapse table-driven tests, fix duplicate item numbering
9166218 refactor,test,docs: extract shared test helpers, refactor Package signature, deduplicate test infrastructure
cd933ef refactor,test,docs: post-dedup cleanup — extract shared helpers, collapse table-driven tests, fix table alignment
a69e0ea docs,test,refactor: post-sprint cleanup — table alignment, test deduplication, nolint fixes
182a7e5 refactor(test): dedup clone groups to ZERO at aggressive threshold 30
435a063 docs(status): threshold-30 deduplication complete status report
2b89410 docs,refactor,test: post-sprint cleanup — documentation, dedup, coverage, status report
```

---

**Next action:** Awaiting user direction. The campaign is at a stable, well-triage'd endpoint.
Recommended next step is the `output.go` Formatter strategy refactor (Appendix E #1) or adding CI
(#2). See question in section g) for the open architectural call.
