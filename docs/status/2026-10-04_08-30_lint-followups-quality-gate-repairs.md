# Status Report: cmdguard-lint Follow-Ups + Repo Quality-Gate Repairs

**Session window:** 2026-10-04 ~07:26 → 08:30
**Scope:** Execute the previous session's follow-up list for the `lint/`
sub-module; everything else in this report was discovered and repaired because
a listed follow-up ("run `nix flake check`") turned out to be broken repo-wide.
**Baseline:** continuation of `2026-10-04_07-11_cmdguard-lint-sub-module-build.md`.

---

## a) FULLY DONE (verified this session)

### Follow-up list execution (items 1–5, 10, 11 from prior report §f)

1. **Vacuous test fixed** — `TestSuppressionLineAboveWithCmdguard`
   (lint/suppress_test.go) now has a control pass that asserts CG001 fires at
   `main.go:10` WITHOUT the directive, plus the suppression pass asserting it
   is silenced. Self-proving: if the rule breaks, the test fails loudly.
2. **Root-module tests re-run** — `TestCLISetVersion` (+ `updates_version`
   subtest, where the suppression comment lives) and the full
   `pkg/cmdguard/v4` package: green.
3. **JSON + SARIF output tests added** (`lint/cmd/cmdguard-lint/main_test.go`)
   — `captureStdout` pipe helper + `writeStaleFixture`; JSON test asserts tool
   name / findings[0].rule=CG002 / position stale.go:3; SARIF test asserts
   version 2.1.0, driver name, one CG002 result, AND round-trips through
   `finding.FindingsFromSARIF`. One `//nolint:tagliatelle` (SARIF camelCase
   keys are the standard, not a style choice).
4. **Replace directive aligned** — `lint/go.mod` now
   `replace github.com/larsartmann/cmdguard/v4 v4.0.2 => ../` (sibling style),
   unpinned variant dropped via `go mod edit`. Verified: GOWORK=off standalone
   build+test AND workspace build+test green.
5. **Doc-comment polish** — all 6 rule files carried stale
   "`*Meta` is the identity header" comments referencing vars deleted in the
   init-cycle refactor. Replaced with check-function docs that name the stable
   rule ID (CG001–CG006). `lint/` module: golangci-lint **0 issues**,
   `-race` tests green.

### Repo gate repairs (discovered while executing the list)

6. **AGENTS.md size gate repaired** — agent-config (BuildFlow) errored: 401
   lines vs 377 max (my prior lint section pushed it over; file was already
   20 over since 09-20). Compressed the 51-line per-file v4 tree listing to a
   10-line directory-level summary (content was "discoverable from code
   alone" per the global doc rules, and already stale: wrong module count,
   missing `website/`). Fixed stale counts (5→6 sub-modules, 6→7 go.work
   modules), Go 1.26→1.27 mentions. Now 356 lines, gate passes.
7. **`nix flake check` repaired and green** (was broken repo-wide, not by
   lint/):
   - `flake.nix` pinned `goPkg = pkgs.go_1_26` while go.mod demands `go 1.27`
     → sandboxed formatters attempted toolchain downloads → DNS-dead sandbox →
     "formatting failures detected" on every Go file. Bumped to `go_1_27`
     (1.27.1 in pinned nixpkgs).
   - That alone was NOT enough: treefmt's `goimports` (nixpkgs `gotools`)
     shells out to `go` resolved from treefmt's own PATH (nixpkgs default
     `pkgs.go`, still 1.26). Wrapped the goimports formatter with
     `GOTOOLCHAIN=local` + `goPkg` on PATH + GOCACHE fallback — the exact
     pattern httputil adopted for its 2026-09-23 incident (found by grepping
     sibling sources in /nix/store). Result: "192 files formatted (0
     changed)" — hermetic, and the lint/ files were already conformant.
   - Gotcha recorded in AGENTS.md (Testing & Build): goPkg major.minor must
     be ≥ the go directive.
8. **golangci-lint exclusion-rename bug fixed** (the big catch): the path
   exclusions for `pkg/cmdguard/`, `examples/`, `benchmarks/`, `tests/`
   referenced the dead linter name `exhaustruct` while the config enables
   `exhaustruct_v5` — the exclusions silently matched nothing, producing 41
   phantom exhaustruct findings repo-wide (check-all lint stage red; the
   "0 lint issues" claim in AGENTS.md false). Fixed the 4 exclusions, removed
   1 redundant dead entry in the `_test.go` block. Root module: **0 issues**
   (verified with both PATH and nixpkgs golangci-lint 2.14.0).
9. **Two real findings fixed properly** (both pre-existing, from the original
   43):
   - `exhaustive` in `writeFormattedError` (pkg/cmdguard/v4/cli_errors_json.go):
     the switch was a membership test mis-modeled as a type dispatch →
     rewritten with `slices.Contains` over a function-local format slice
     (function-local to respect gochecknoglobals; caught on first attempt).
   - `contextcheck` in examples/taskctl/main.go:63 (`buildApp(rec)` chain):
     false positive — cmdguard DI construction is context-free by design; ctx
     arrives at `ExecuteWithArgs`. Suppressed with a surgical
     `//nolint:contextcheck` + reason on the line (did NOT re-architect the
     fresh parallel-session refactor).
10. **Docs harvest completed**:
    - FEATURES.md: lint row in sub-module table (+5→6 count fix), new
      "Usage Linter (`lint`)" capability section (DONE vs PLANNED).
    - TODO_LIST.md: new "Lint Sub-Module" section L1–L8 (tag, BuildFlow
      wiring, cross-file dataflow, baseline mode, provider options, plugin
      layer, README generation, CurrentMajor standing duty).
    - Prior status report §f items 1–11 annotated inline with
      ✅ DONE / 🟡 AWAIT-USER / 🔴 TODO resolutions.
    - AGENTS.md: exclusion-rename gotcha recorded.

### Verification state at report time

- `lint/` module: build, vet, `go test -race`, golangci-lint — all green,
  0 issues.
- Root module: build, tests, golangci-lint — 0 issues.
- `nix flake check` — **passes** (exit 0).
- `buildflow format` — all formatter steps ✔ (only pre-existing failures
  remain, see b).
- All work auto-committed by the daemon (commits d9f3e5e…a950d0a, 07:31→08:02).

---

## b) PARTIALLY DONE

1. **`nix run .#check-all` final green not yet demonstrated.** The last run
   cleared build, race tests, and root lint, then reported 2 phantom
   `typecheck` "issues" caused by golangci-lint's persistent fact cache
   referencing evicted GOCACHE export files
   (`/mnt/buildcache/go-build/…-d: no such file or directory`). The code is
   clean — direct `golangci-lint run ./...` reports 0 issues. Remedy started:
   `golangci-lint cache clean` (still deleting at report time — large dir).
   Needs: re-run `check-all` once the cache clean finishes.
   > ✅ RESOLVED 2026-10-04 (later session): cache clean completed (needed a
   > retry on the /mnt/buildcache mount); `nix run .#check-all` exit 0 —
   > build, race tests (15 pkgs), lint 0 issues × 7 modules, flake check,
   > tidy × 7 all green.
2. **`buildflow format` exit=69 (pre-existing, environmental, not mine):**
   - 8 failed step executions, all `license-check` fan-out: `go-licenses`
     E1004 "Package net/mail does not have module info" — go-licenses is
     incompatible with the Go 1.27 toolchain (upstream issue #128). Fleet
     issue → belongs in the BuildFlow repo, not here.
   - lychee link errors: `cmdguard.lars.software` serves a certificate valid
     only for `*.firebaseapp.com` (DNS/hosting misconfiguration outside the
     repo).

---

## c) NOT STARTED (deliberate / gated — unchanged from prior session)

> Session-3 resolution: ALL four gates below were approved and executed on
> 2026-10-04 (later session) — see §g annotations.

- ✅ Tag `lint/v0.1.0` — DONE (TODO_LIST L1): CHANGELOG entry, annotated tag
  pushed, proxy-verified via scratch consumer (`go get` + build + run),
  GitHub release published.
- ✅ Wire `lint/provider` into the BuildFlow repo (L2) — DONE: blank import
  in `sdk_imports.go`, require v0.1.0 + pinned local replace, `go work
  vendor` + vendorHash, wiring regression test, inventory guards updated
  (ToolCmdguardLint, moduleScopedGoToolSpecs, edge snapshot 279→287),
  `docs --check` green (provider count 118).
- ✅ Restore dead pre-commit hook — DONE: root cause was the GLOBAL
  `core.hooksPath=.githooks` (fleet convention) while this repo lacked the
  directory. Shipped a tracked `.githooks/pre-commit` (fleet pattern) running
  BuildFlow pre-commit mode; advisory until the license-check Go 1.27 fleet
  fix lands. Verified live on its own commit.
- ✅ CHANGELOG.md entry for the sub-module (ships with the tag) — DONE.
- Prior report §f items 12–50: cross-file dataflow, baseline mode, provider
  options, golangci plugin layer, fuzzing, benchmarks, coverage gaps,
  CG007–CG009 candidates, suppression AST/expiry, ternary exit codes, etc.

---

## d) TOTALLY FUCKED UP (all recovered — honest ledger)

1. **First `go mod edit -replace` attempt used wrong syntax**
   (`…=v4.0.2=../` → "malformed import path: trailing slash"). Correct form:
   `-replace=mod@v4.0.2=../`. Also briefly created a double-replace block
   (pinned + unpinned) before dropping the unpinned one.
2. **Three edit-tool failures on daemon-touched files** ("modified since last
   read") — the auto-commit daemon rewrites files mid-session, invalidating
   read-state. Recovered by re-viewing. Lesson: with a live commit daemon,
   view-before-edit must be immediately fresh.
3. **Mangled a nix store path from memory** trying to read a build log
   (3 invented paths, all ENOENT). Recovered by rebuilding the derivation
   with `-L`. Lesson: copy store paths verbatim from output, never recall.
4. **First goimports fix (goPkg bump alone) did not work** — I pinned the Go
   the devShell uses, but goimports resolves `go` from treefmt's own PATH
   (nixpkgs default). Wasted one build cycle. Lesson: trace which binary a
   sandboxed tool actually resolves before pinning anything.
5. **First `machineReadableErrorFormats` was a package-level var** → tripped
   gochecknoglobals (predictable — the repo enumerates its globals). Moved
   function-local.
6. **Time burned enumerating "8 failed steps"** in buildflow output when a
   single `grep '✗'` showed the one real failing step (license-check); the 8
   was its module fan-out count. Lesson: count unique failing step names
   first, then read details.

---

## e) WHAT WE SHOULD IMPROVE (process takeaways)

- **Run the flake gate at module-addition time**, not as a trailing TODO. The
  treefmt check had been broken for everyone since the go.mod 1.27 bump and
  only surfaced because a follow-up list item forced it.
- **Linter renames void exclusions silently.** When renaming a linter in the
  enable list (exhaustruct → exhaustruct_v5), grep the exclusions for the old
  name in the same commit. Better: a config test asserting every exclusion's
  linter name is in the enabled set (see f #25).
- **golangci-lint fact-cache corruption produces phantom typecheck findings.**
  `golangci-lint cache clean` belongs in the triage runbook; check-all could
  pin a private `GOLANGCI_LINT_CACHE` to stay hermetic.
- **The pre-commit hook is dead repo-wide** (`core.hooksPath=.githooks`, dir
  missing). Right now NO quality gate runs on commit — quality rests on
  explicit runs plus the daemon's blind auto-commits. This is the single
  highest-leverage repair pending user approval.
- **Stale claims need mechanical teeth.** "0 lint issues" and "Go 1.26" were
  both false at session start; agent-config caught the file-size drift but
  nothing catches content drift. The lint module itself (CG002-style stale
  detection) is the pattern to follow fleet-wide.

---

## f) NEXT — ranked (top ~25 of the ~60 combined backlog; prior report §f 12–50 still apply)

1. Re-run `nix run .#check-all` after the cache clean finishes → confirm full
   green (finish of b1).
2. Tag `lint/v0.1.0` via go-release skill (gated: g).
3. Wire `lint/provider` into BuildFlow (gated: g).
4. Restore pre-commit hook: `buildflow precommit install` (gated: g).
5. CHANGELOG.md entry for the lint sub-module.
6. Cross-file dataflow CG001 (RootCommand ident used in another file).
7. Cross-file dataflow CG006 (Execute error printed elsewhere).
8. Cross-file dataflow CG003/CG004 (constructor in one file, use in another).
9. Baseline/ratchet mode for incremental adoption.
10. Provider options (enable/disable/exclude via BuildFlow config).
11. Exclusion paths beyond vendor/testdata/hidden, via option.
12. golangci-lint module-plugin distribution layer.
13. SARIF upload snippet (CI) in lint/README.
14. Fuzz walker + suppression parser on arbitrary input.
15. Benchmarks: scan time on the cmdguard repo as baseline.
16. Coverage gaps: dot-import branch, `fangLocalName`, skipped-file path.
17. Suppression via AST comments (directive-in-string-literal FP hardening).
18. Suppression expiry (`ExpiresAt`).
19. CG007–CG009 rule candidates (PostRunE-vs-Cleanup, ExecuteAndExit+printing,
    WithCommit duplicate).
20. `CurrentMajor` from generated/versioned source instead of hand-bumped
    constant.
21. Ternary exit code (`--min-confidence` → exit 2 triage).
22. license-check fleet fix: pin/patch go-licenses for Go 1.27 (BuildFlow
    repo).
23. Website SSL cert fix (firebaseapp cert served on cmdguard.lars.software).
24. Add go-licenses/govulncheck to flake devShells (kills the nix-run
    fallback warnings).
25. Config test: every `.golangci.yml` exclusion linter name must be in the
    enabled set (mechanizes e's rename lesson).

---

## g) QUESTIONS (cannot resolve myself — carried from prior session, still open)

1. **BuildFlow wiring:** import `lint/provider` into the BuildFlow repo now
   (fleet activation), or dogfood locally + tag first?
   > ✅ ANSWERED: wire now. Done — see §c.
2. **Pre-commit hook:** restore via `buildflow precommit install`? It changes
   `core.hooksPath` git config — I will not touch git config without approval.
   > ✅ ANSWERED: yes. Root cause turned out to be the global hooksPath
   > convention; fixed via tracked `.githooks/pre-commit` instead (no git
   > config change needed) — see §c.
3. **Tag timing:** cut `lint/v0.1.0` immediately so timesheets can adopt, or
   wait for your review of rule severities/messages first?
   > ✅ ANSWERED: tag now. Done — see §c.
