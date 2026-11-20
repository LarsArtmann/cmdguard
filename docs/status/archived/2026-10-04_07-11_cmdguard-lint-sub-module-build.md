# cmdguard-lint — New Sub-Module Build Status

**Date:** 2026-10-04 07:11 CEST
**Scope of this session:** Build a linter for the anti-patterns found in the timesheets cmdguard deep-dive (`timesheets/docs/research/2026-10-04_cmdguard-deep-dive.html`), on top of `go-finding` + `go-finding/toolsdk` + `go-linter-sdk`, as a new sub-module of this repo.
**Skills loaded:** linter-building, buildflow.

---

## a) FULLY DONE (verified)

### Research & design

- Read the full deep-dive report; extracted the 9 findings (F1–F9) and mapped each to lintable/unlintable with explicit rationale (absence-of-usage findings deliberately excluded — FP-prone, not file:line-attributable).
- Studied `go-finding` (Finding/Builder/Detector/Position/Suppression APIs, module layout), `toolsdk` (Spec contract, `Register`, `OnGoModule()`, ModuleFanOut, per-field semantics), `go-linter-sdk` v0.3.1 (Rule/RuleFunc/Registry, `DetectorFromRegistry`, `FilterRules`, exit codes) — all against **published tags**, not local HEADs (caught `CategoryBestPractice` existing only locally; used the SDK's documented custom-category mechanism instead).
- Studied this repo's sub-module pattern (root-level dir + own go.mod + version-pinned replace + go.work entry) and the erraudit provider as the reference provider implementation.
- Mechanism decision (per linter-building skill): pure **AST walker** (go/parser, no type info, no build required) — all six rules are syntax-shape + dataflow-lite.

### The `lint/` sub-module (`github.com/larsartmann/cmdguard/lint`)

- **6 rules**, each with origin citation, trust signals, severity, confidence:
  - **CG001 execute-bypass** (critical): `fang.Execute` receiving `RootCommand()` (direct call or ident), gated on project importing cmdguard — deep-dive F1.
  - **CG002 stale-major-import** (error): import of a cmdguard major < CurrentMajor — F2.
  - **CG003 panic-on-constructor** (error): `panic` in the `err != nil` branch of NewCLI/NewCommand/NewParentCommand/AddCommand errors — F6.
  - **CG004 set-version-runtime** (warning): `<cliVar>.SetVersion` where cliVar traces to NewCLI in-file — F7.
  - **CG005 duplicate-version-opts** (error): `WithCLIVersion` + `WithFangOptions(fang.WithVersion(...))` in one NewCLI call — ADR-001.
  - **CG006 execute-error-reprint** (warning): error from `<cliVar>.Execute` re-printed via fmt Print family (fmt.Errorf wrapping deliberately not flagged) — v4 error/exit contract.
- **Single-source-of-truth rule table** (`allRuleDefs()` in rules.go) feeding `AllRules()`, `NewRegistry()`, and the single-pass `Detect()` fast path. No package-level rule vars (they caused a real Go initialization cycle; gopls caught it).
- **Walker** (`walk.go`): skips vendor/testdata/hidden dirs, generated files ("Code generated… DO NOT EDIT"), unparseable files (recorded, not fatal); import alias resolution with known-paths table (fang/v2→"fang", cmdguard major segment convention); per-directory process-lifetime memo + `ClearCache()`.
- **Suppressions**: `//cmdguard-lint:ignore <ID> <reason>` on the finding line or the line above; reason required; wrong rule doesn't suppress.
- **Provider** (`lint/provider/`): toolsdk Spec registered (`cmdguard-lint`, `OnGoModule()`, `ModuleFanOut: true`, Detect via `WorkingDirFromContext`), all Spec fields explicit.
- **CLI** (`lint/cmd/cmdguard-lint/`): dogfoods cmdguard v4 itself (WithCLIVersion, WithSignalHandling, ExecuteAndExit, values-tag enum); `lint --dir --output text|json|sarif --enable --disable` + `rules`; findings → exit 1 with single-display error contract honored; `newApp()` split for ExecuteWithArgs testing.

### Tests (all green, `-race`)

- Walker tests (skips, alias/major resolution, majorOf, packageNameFor, cache clear).
- Per-rule positive + negative fixture tests (exact file:line assertions).
- Suppression tests (same-line, requires-reason, wrong-rule, parse table).
- Provider spec-contract test + WorkingDir-context detection test.
- CLI tests via `ExecuteWithArgs` (clean→pass, findings→fail with message, rules list).

### Validation (linter-author discipline)

- **Positive corpus (timesheets):** 18 findings — CG001 at `cmd/timesheet/main.go:26` (audit cites 26–29), 15× CG002 with **exact 15/15 match against `rg` ground truth (0 FP, 0 FN)**, CG003 at `cmdguard_helpers.go:21` panic, CG004 at `guard.go:50` (audit cites 50–51). Every location matches the hand audit.
- **Negative corpus:** cmdguard itself (0 after 1 legitimate self-suppression in `cli_lifecycle_test.go`), erraudit, go-structure-linter, branching-flow — all 0 findings.
- **Discrimination proof:** inverted CG002's condition → tests FAIL → restored → pass. The tests test.

### Quality gates

- `golangci-lint run ./...` in lint/: **0 issues** — via real fixes (rule-table refactor killing 12 gochecknoglobals findings, explicit struct literals, exhaustive Spec, Fprint writers for forbidigo, SplitSeq/Cut modernize, named constant for mnd, wrapper-error wrapcheck fixes, no named returns), not exclusion creep. Only nolints: the two documented process-lifetime cache globals + one deliberate-skip nilerr, each with reasons.
- `go build`, `go vet`, `go test -race` green; **GOWORK=off standalone build/test green** (replace directive works for downstream-shape builds).
- BuildFlow fast mode run: lint module contributes **0 findings**; the 2 failing steps (license-check go-licenses E1004, website eslint/vulns/TS) are **pre-existing**, unrelated to this work.

### Wiring & docs

- `go.work`: 7 modules now (added `use ./lint`).
- `flake.nix` check-all: lint added to build/test/lint/tidy loops; description 6→7 modules.
- `AGENTS.md`: project tree, sub-module dependency table, design principle #15/#20 counts, nested-modules gotcha, new full **lint** section (rule table, no-rule-vars lesson, CurrentMajor bump duty, suppression format, corpus evidence).
- `lint/README.md`: rules table, usage (CLI/BuildFlow/library), suppression format, add-a-rule guide, origin + non-goals.

---

## b) PARTIALLY DONE / rough edges

- **Vacuous test shipped:** `TestSuppressionLineAboveWithCmdguard`'s fixture passes `rootOfCLI()` (not a `RootCommand()` call) to fang.Execute, so CG001 is silent regardless of the directive — the test passes but does **not** prove line-above suppression of CG001. Needs a fixture where the call actually triggers CG001.
- **JSON/SARIF output paths implemented but untested** (only default text is covered by tests/manual runs).
- **Root-module tests not re-run** after adding the suppression comment to `pkg/cmdguard/v4/cli_lifecycle_test.go` (comment-only edit, near-zero risk, but unverified).
- **`nix flake check` / `nix fmt` never run** on the new files (treefmt gate unverified — my own final todo item, not executed).
- **Replace-directive inconsistency:** lint uses version-less `replace …/cmdguard/v4 => ../`; siblings pin `v4.0.2 => ../`. Works, but breaks the family pattern.
- **Rule-file doc comments uneven** after the mechanical meta-strip refactor (e.g. `rule_duplicate_version.go`'s comment line reads awkwardly/truncated).
- **CG001/CG006 dataflow is file-scoped:** `root := cli.RootCommand()` in file A + `fang.Execute(nil, root)` in file B is missed; same for Execute-error printed in another file.
- **Unexpected working-tree diffs I did not author** (left untouched per policy): `examples/taskctl/main.go` + `run_test.go` (a `run()` split refactor by someone else), `FEATURES.md`. Flagging for awareness.

## c) NOT STARTED (deliberate or deferred)

- Deep-dive findings F3 (sentinel classification), F4 (config loader duplication), F5 (DI scope/doctor), F8 (tabwriter/output layer), F9 (UX options/groups) — **deliberately not rules** (documented in README: absence-of-usage, not line-attributable without heavy type analysis).
- Release: no `lint/v0.1.0` tag, no CHANGELOG entry, no release-preflight — **external consumers cannot `go get` the module yet**.
- BuildFlow activation: BuildFlow repo does not yet import `lint/provider` (one import + release away).
- FEATURES.md / TODO_LIST.md / website docs for the linter.
- toolsdk Options (enable/disable/excludes at provider level), baseline/ratchet mode, golangci-lint plugin distribution, SARIF-in-CI example, fuzz targets, benchmarks, coverage report, LSP, dependabot entries.

## d) TOTALLY FUCKED UP (all recovered, but honest ledger)

~~1. **Initialization cycle** — rule vars referencing checks referencing rule vars (spec-legal cycle). gopls caught it; fixed with meta-passing table design. Cost: one rewrite of all six rule files.~~ done — rules.go allRuleDefs() table shipped; package-level rule vars gone (AGENTS.md warns against reintroducing)
~~2. **Assignment-shape bug in ALL four collectors** — `len(Lhs) != len(Rhs)` guard skipped the canonical `cmd, err := f()` form (2 LHS / 1 RHS), silently zeroing CG003/CG004/CG006. Found via debug tests; fixed with shared `visitCallAssignments`.~~ done — visitCallAssignments in ast_helpers.go; CG003/CG004/CG006 fixtures cover the 2-LHS/1-RHS form
~~3. **Package-name guessing bug** — generic major-suffix stripping mapped `cmdguard/v4/pkg/cmdguard/v4` to local name "cmdguard" instead of "v4"; fixed with `packageNameFor` special-case.~~ done — packageNameFor in walk.go handles the cmdguard major-segment case; walk_test.go covers it
~~4. **Accidentally created `lint/go.work`** (bad edit target) — would have broken the module; deleted immediately.~~ done — no lint/go.work exists; only root go.work with use ./lint
~~5. **sed deletion removed a closing paren** (nolint line was part of a multi-line `errors.New(` statement) — syntax break, fixed.~~ done — recovered in-session; lint module built and released through lint/v0.3.0
~~6. **Multiple junk artifacts in first drafts** (`localNameUnused` hack, `panic("unreachable")` placeholder, `filepath.ReadFile` which doesn't exist) — all removed; too many fix-rounds driven by linter output instead of writing it clean the first time.~~ done — grep clean, no localNameUnused / unreachable-placeholder artifacts remain in lint/
~~7. **Bulk sed/perl rewrites** created orphaned references needing 3+ repair rounds — mechanical refactors of Go code via regex were a net time loss vs. rewriting the six small files by hand.~~ done — recovered by hand-rewrites; rule files hand-maintained through lint/v0.3.0

## e) WHAT WE SHOULD IMPROVE (process takeaways)

- Write fixtures so positive tests fail for the right reason (the vacuous suppression test would have been caught by deleting the directive and asserting the finding reappears — "remove-one-thing" verification).
- Prefer hand-rewrites of small files over regex surgery; the mechanical pass was where every fuck-up cluster happened.
- Run the format gate (`nix fmt` / `nix flake check`) before declaring a module done — it was on the todo list and got skipped when BuildFlow "felt" sufficient.
- gopls diagnostics were screaming correctly the whole time; trust them earlier instead of assuming stale-cache noise.
- Root-module regression run after touching root files, even comment-only.

## f) NEXT — up to 50 items (Pareto-ranked)

~~1. Fix vacuous `TestSuppressionLineAboveWithCmdguard` (trigger CG001, then suppress it). ✅ DONE 2026-10-04 (control pass now asserts CG001 fires without the directive; suppression pass asserts it is silenced)~~ done — control + suppression passes assert CG001 fires without the directive then is silenced (suppress_test.go)
~~2. Add JSON + SARIF output tests (round-trip via `FindingsFromSARIF`). ✅ DONE 2026-10-04 (both green, SARIF round-trips CG002@stale.go:3)~~ done — TestLintCommandSARIFOutput + FindingsFromSARIF round-trip in main_test.go; JSON output tested
~~3. Run `nix flake check` + `nix fmt` on the new files. ✅ DONE 2026-10-04 (treefmt check green, 192 files 0 changed — after fixing `goPkg` go_1_26→go_1_27 + GOTOOLCHAIN=local goimports wrapper in flake.nix; the check had been broken repo-wide by the go.mod 1.27 bump)~~ done — flake.nix goPkg = pkgs.go_1_27; check-all loops include the lint module
~~4. Align replace directive with siblings (`v4.0.2 => ../`). ✅ DONE 2026-10-04~~ done — family-aligned: lint/go.mod carries cmdguard/v5 => ../ beside sibling => ../ replaces
~~5. Re-run root-module `go test ./... -race` (cli_lifecycle_test.go edit). ✅ DONE 2026-10-04 (TestCLISetVersion + full v4 package green)~~ done — suppression directive present in cli_lifecycle_test.go:43; inline note records v4 package green
~~6. Review/hand-off the non-authored taskctl/FEATURES diffs. ✅ RESOLVED (auto-commit daemon committed them; taskctl run() refactor untouched by lint work; FEATURES lint section added 2026-10-04)~~ done — FEATURES.md lint section present; diffs committed by the daemon per inline note
~~7. Tag `lint/v0.1.0` (go-release skill: sub-module tag, proxy check). 🟡 AWAIT-USER (TODO_LIST L1)~~ done — lint/v0.1.0 published 2026-10-04; v0.2.0/v0.2.1/v0.3.0 follow (git tag -l lint/*)
~~8. CHANGELOG.md entry for the sub-module. 🔴 TODO (ships with the v0.1.0 tag)~~ done — root CHANGELOG.md carries [lint/v0.1.0] 2026-10-04 plus v0.2.0/v0.2.1/v0.3.0 entries
~~9. Import `lint/provider` into BuildFlow (fleet activation). 🟡 AWAIT-USER (TODO_LIST L2)~~ done — blank import in BuildFlow tools/providers/sdk_imports.go; fleet provider count 118 (AGENTS.md)
~~10. FEATURES.md + TODO_LIST.md updates. ✅ DONE 2026-10-04 (lint sub-module section + L1–L8 harvest)~~ done — FEATURES.md lint rows + TODO_LIST Lint Sub-Module harvest section fed from this report §f
~~11. Doc-comment polish pass on rule files. ✅ DONE 2026-10-04 (stale `*Meta is the identity header` references replaced with check-function docs naming the rule ID)~~ done — rule files carry check-function docs naming rule IDs (e.g. CG007 header in rule_ux_baseline.go)
~~12. Cross-file CG001 dataflow (RootCommand ident across files).~~ **Won't implement — never pursued; crossfile.go serves CG003/CG004/CG006 only, CG001 stays file-scoped.**
~~13. Cross-file CG006 (Execute error printed elsewhere).~~ done — rule_error_reprint.go consults package-wide name tracing (crossfile.go)
~~14. toolsdk Options on the provider (enable/disable/exclude).~~ done — lint.FilterByIDs + provider enable/disable options via OptionsFromContext (provider.go)
~~15. Exclusion paths (beyond vendor/testdata/hidden) via option.~~ **Won't implement — never pursued through v0.3.0; options surface is enable/disable only, walker skips unchanged.**
~~16. Baseline/ratchet mode for incremental adoption.~~ done — baseline.go + CLI --write-baseline/--baseline shipped in lint/v0.2.0 (CHANGELOG)
17. golangci-lint custom-module plugin distribution. **→ **→ TODO_LIST.md L6 (golangci-lint plugin bridge; open TODO, scoped as not-a-thin-wiring-task)****
~~18. SARIF upload snippet in README (CI).~~ **Won't implement — never pursued; lint/README documents --output sarif for CI upload, no Actions snippet.**
~~19. Fuzz the walker + suppression parser on arbitrary input.~~ **Won't implement — never pursued through v0.3.0; no func Fuzz targets exist in lint/.**
~~20. Benchmarks (scan time on cmdguard repo as baseline).~~ **Won't implement — never pursued through v0.3.0; no func Benchmark targets exist in lint/.**
21. Coverage report; close gaps (dot-import branch, `fangLocalName`, skipped-file path). **→ **→ TODO_LIST.md H9b (dot-import bare-constructor tracing; dot/blank arms covered by crossfile_imports_test.go)****
~~22. Suppression via AST comments instead of raw line scan (directive-in-string-literal FP hardening).~~ **Won't implement — never pursued; suppress.go ships the line-scan regex, format stable at v0.3.0.**
~~23. Suppression expiry (`ExpiresAt`) support.~~ **Won't implement — never pursued; suppression format carries no expiry field (ignore <ID> <reason>).**
~~24. CG007 candidate: WithPostRunE used where WithCleanup is required (PostRunE not called on RunE error).~~ **NOT-DO — CG007 shipped in lint/v0.3.0 as missing-short-description; PostRunE/WithCleanup candidate superseded.**
~~25. CG008 candidate: `ExecuteAndExit` combined with custom printing.~~ **NOT-DO — CG008 shipped in lint/v0.3.0 as missing-CLI-version; ExecuteAndExit-printing candidate superseded.**
~~26. CG009 candidate: `WithFangOptions(fang.WithCommit(...))` + `WithCLICommit` duplicate (mirror of CG005).~~ **Won't implement — never pursued; rules stop at CG008 (rules.go); CG009 slot reused by a different idea (TODO_LIST CT1).**
~~27. `CurrentMajor` from a generated/versioned source instead of a hand-bumped constant.~~ **Won't implement — never pursued; CurrentMajor stays a hand-bumped const (walk.go:19) with documented bump duty.**
~~28. Ternary exit code (`--min-confidence` → exit 2 triage).~~ **Won't implement — never pursued; exit contract stays clean→0 / findings→1, no --min-confidence flag.**
~~29. Surface skipped-unparseable files as optional info findings.~~ **Won't implement — never pursued; unparseable files are deliberately excluded from findings (rules.go).**
~~30. Golden-file test for text output format.~~ **Won't implement — never pursued; text output covered by table-driven ExecuteWithArgs tests, no golden files.**
~~31. E2E binary test (build once in TestMain, run against fixture, assert exit code).~~ **Won't implement — never pursued; no TestMain binary-build E2E test in lint/.**
~~32. Property test: linter never panics on arbitrary Go-ish corpora.~~ **Won't implement — never pursued; no property/fuzz corpus tests in lint/.**
~~33. Dependabot entries for the new module (buildflow -s dependabot-auto-configure).~~ done — .github/dependabot.yml carries directory: /lint
~~34. website/ docs page + root README table row.~~ done — root README.md lint row (CG001-CG008) + website/src/content/docs/guides/lint.mdx exist
~~35. Dogfood gate: run cmdguard-lint on cmdguard in BuildFlow (self-consumer).~~ done — provider compiled into BuildFlow via sdk_imports.go blank import; fleet-wide since lint/v0.3.0 (AGENTS.md)
~~36. examples/ before-after fixtures in lint/README.~~ **Won't implement — never pursued; README carries rule table + usage; before/after fixtures live in *_test.go.**
~~37. Restore/replace the dead pre-commit hook (core.hooksPath → missing .githooks).~~ done — tracked .githooks/pre-commit exists; AGENTS.md documents the global core.hooksPath=.githooks convention
~~38. Fix pre-existing taskctl exhaustruct findings (repo's "0 lint issues" claim is stale for examples/).~~ done — .golangci.yml path-excludes examples/ under exhaustruct_v5
~~39. LSP integration via go-finding's lsp.go.~~ **Won't implement — never pursued through v0.3.0; no LSP surface in lint/.**
~~40. Per-rule parallelism via `DetectorsFromRegistry` in provider.~~ **Won't implement — never pursued; provider runs Detect single-pass, no per-rule DetectorsFromRegistry parallelism.**
~~41. Cache invalidation option (TTL or mtime fingerprint) for long-lived processes.~~ **Won't implement — never pursued; analyze() per-directory memo + manual ClearCache() is the shipped cache story.**
~~42. Colocate rule IDs with a machine-readable rule manifest for docs generation.~~ done — rules --markdown prints the README table; rules_table_test.go drift-guards it vs allRuleDefs()
~~43. Consider flagging `SetVersion` on the cmdguard CLI in v3 consumers too (alias any major).~~ **Won't implement — never pursued through v0.3.0; CG004 traces NewCLI-built vars without cross-major aliasing.**
~~44. Add `lint` to `docs/API.md` index.~~ **NOT-DO — docs/API.md no longer exists in the repo; index item moot.**
~~45. Verify provider under BuildFlow's ModuleFanOut wrapper (multi-module workspace fan-out test).~~ done — provider_test.go asserts spec.ModuleFanOut:true for multi-module workspaces; shipped fleet-wide at v0.3.0
~~46. Rule-messages review pass with a consumer (timesheets) for tone/actionability.~~ done — timesheets corpus validated 18/18 findings at hand-audited locations, 0 FP (AGENTS.md lint section)
~~47. Timesheets: adopt the linter in its CI once tagged (closes the loop with the audit).~~ **Won't implement — cannot confirm adoption; external repo, no cmdguard-lint in timesheets .github/workflows/ci.yml.**
~~48. Cross-repo corpus sweep across all 14+ cmdguard consumers for a measured FP rate.~~ **Won't implement — full 14+ sweep never ran; measured runs exist for timesheets+erraudit PR branches (TODO_LIST F7).**
~~49. Add `//cmdguard-lint:ignore` lint-of-itself (staleness audit like erraudit's nolint-audit).~~ **Won't implement — never pursued; no staleness audit of //cmdguard-lint:ignore directives in lint/.**
~~50. Write the assignment-shape lesson to crush-config lessons (collector must handle 2-LHS/1-RHS).~~ **Won't implement — never written; crush-config references/lessons.md carries no 2-LHS/1-RHS collector lesson.**

## g) QUESTIONS (cannot resolve myself)

~~1. **Fleet rollout:** should I wire `lint/provider` into the BuildFlow repo now (making cmdguard-lint a fleet tool for every Go repo), or hold until it's dogfooded here + tagged v0.1.0?~~ done — blank import in BuildFlow sdk_imports.go (fleet provider count 118); fleet-wide via lint/v0.3.0 bump
~~2. **Dead pre-commit hook:** `core.hooksPath=.githooks` points at a missing directory (gate silently dead). Restore via `buildflow precommit install`, or leave untouched? I won't change git config without your call.~~ done — tracked .githooks/pre-commit restored; global core.hooksPath=.githooks convention per AGENTS.md
~~3. **Tag now or review first:** tag `lint/v0.1.0` immediately so external consumers (timesheets) can adopt, or do you want to review rule messages/severities (e.g. CG004/CG006 as warnings) before the first tag freezes the IDs' reputations?~~ done — lint/v0.1.0 published 2026-10-04; v0.2.0/v0.2.1/v0.3.0 followed; rule IDs CG001-CG008 frozen
