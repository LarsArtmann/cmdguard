# Status: GitHub Issues/PRs Closure + Master Repair — 2026-11-02 08:40

Session directive: "Review all open GitHub issues and PRs and fix them aka get them done!"
Follow-up directive: "fix/git commit!" (repair the diverged/broken git state), then this report.

---

## Context at session start

- Open on GitHub: issues #11, #12, #13 (all filed 2026-10-14 by Lars against v4.1.0), PRs #16 and #18 (Dependabot go_modules bumps), plus 0 open Dependabot vulnerability alerts (the previous session's website dep refresh had already cleared the 19).
- A PARALLEL SESSION was actively working in the same repo on a "config-truth workstream" (uncommitted changes across config*.go, cli_accessors.go, cli_options.go, flags.go, koanf_loader.go, errors*.go; untracked config_keys.go, config_strict_test.go; local commit 692d786 pushed to master mid-session).
- Machine under heavy multi-session load (load average 146 at session start, 224 at report time; 54 users).

---

## a) FULLY DONE

1. **Issue #12 (required:"true" never validates on execution path) — CLOSED.** Verified already fixed by the parallel session in d3cdf93 (shipped in v4.2.0): `PersistentPreRunE` runs `registry.ValidateFlags` after `ParseFlags` (root flags, with `skipsFlagValidation` escape for describe/version/doctor/completion/help), and `prepareRunContext` validates per-command flags. Pinned by `TestExecute_RequiredFlagEnforced` + `TestValidateFlags_*` (all green locally). Closed with an explanatory comment referencing the shipped fix; no code change needed.
2. **Issue #11 (ExitCode: 0 from ExitCoder must mean "no override") — FIXED + CLOSED.** `extractExitCode` (cli_errors_json.go:78) now skips a zero code from an `ExitCoder` and falls back to 1, mirroring go-error-family's guard; `ExitCode(nil) == 0` unchanged; JSON error envelope `code` follows the same rule; `ExitCode` doc comment updated. Tests: `zeroExitError` stub (errorfamily's contract) in `TestExtractExitCode` + `TestExitCode`. Landed in 887dcae.
3. **Issue #13 (flag values leak across ExecuteWithArgs on reused CLI) — FIXED + CLOSED.** New `resetFlagStateForReuse` runs at every `Execute`: (1) `Changed=false` across the whole command tree, (2) values reset to defaults — with stringSlice flags rebuilt via a scratch-FlagSet Value transplant (pflag slice `Set` APPENDS after first write and cannot re-read its bracketed `"[a,b]"` default; re-registering in the real FlagSet PANICS on redefinition), (3) the typed config restored in place from a `pristineConfig` snapshot taken at the end of `initialize` (fields without tag defaults are never rewritten by `ParseFlags`, so the flag reset alone leaves them stale). 4 regression tests in `cli_flag_reuse_test.go`: --apply leak, fresh values, slice accumulation, required-flag re-validation. Landed in 887dcae.
4. **PR #16 (lint: go-finding 1.14.0→1.15.0, toolsdk 1.15.0→1.16.0) — MERGED.** Verified diff shape clean (external pins only) and CI fully green; rebase-merged as 52207c6 per repo convention (viewerDefaultMergeMethod REBASE, deleteBranchOnMerge true).
5. **PR #18 (root: koanf 2.3.8→2.4.1, go-output 0.38.4→0.38.7) — COMPLETED + MERGED.** Diagnosed the 4 failing `submodule-smoke` jobs: Dependabot bumped only root go.mod/go.sum; the sub-modules (telemetry/glamour/prompts/spinner) went stale in module mode ("updates to go.mod needed"; go.work is untracked so CI runs pure module mode). Completed the PR in an isolated worktree (/tmp/pr18): copied go.work, workspace-aware `go mod tidy` per sub-module via the main devShell, verified CI-style `GOWORK=off go build + go test` per module (all green), pushed 0549cd3 to the Dependabot branch, CI went fully green, rebase-merged as 98d89e3+e530204. Worktree removed.
6. **Docs + memory.** CHANGELOG [Unreleased] Fixed entries for #11/#13; consumer skill (.agents/skills/cmdguard/SKILL.md) exit contract documents the 0-means-no-override rule and CLI reuse safety; AGENTS.md Flag System gotcha records the full reset mechanism and both pflag traps. Pushed as 71f7ecf (skill file also picked up a dprint canonical reflow).
7. **Final GitHub state: 0 open issues, 0 open PRs, 0 open Dependabot alerts.**

## b) PARTIALLY DONE

1. **Master repair for the accidental WIP sweep (see d).** Root cause fixed in the /tmp/fix-master worktree: 4 surgical edits removing exactly the 5 foreign hunks from cli.go (strictConfigKeys, configExplainCmd, configSetFields+comment, configFilePath+comment, 3-var loadConfigFileOrSkip call + map fill); verified via `rg` that no other file on master references the removed identifiers. **NOT yet committed/pushed** — the worktree's build+test verification run has been CPU-starved for ~90 minutes (load 224; the go processes from other sessions dominate). Master CI is RED until this lands.
2. **Local git-state reconciliation (the "fix/git commit!" ask).** Local master diverged: 2 local commits (8babb0e, 41196e0) whose content is ALREADY on origin (rebase-pushed earlier as 887dcae/71f7ecf; git will skip them as cherry-picks) vs 5 origin commits. `git pull --rebase` is blocked by the parallel session's unstaged changes. Planned safe sequence (NOT executed): manual stash → rebase (equivalent commits skip) → stash pop → re-apply the parallel session's cli.go WIP hunks as an unstaged patch (`git diff origin/master 41196e0 -- pkg/cmdguard/v4/cli.go | git apply`) so their session is made whole again (the repair removes those hunks from master; without re-applying, their working tree would silently lose them).
3. **Full local gate (`nix run .#check-all`).** Never completed green this session: first run raced the parallel session's mid-edit tree (their config_keys.go/koanf_loader.go had transient compile errors), second run confirmed core broken by their WIP at that moment. Used scoped verification instead: core package build+tests green, `golangci-lint` clean on all my files (the 8 findings were in the parallel session's WIP files), sub-module module-mode builds/tests green in the PR worktree, CI green on every push (except the current red master from d).

## c) NOT STARTED

1. **v4.2.1 patch release** carrying #11/#13 (fixes sit in [Unreleased]; go-credentials-store still consumes v4.2.0 with the exit-0 bug and the max(ExitCode,1) workaround).
2. Restore of the parallel session's cli.go WIP hunks after the local rebase (planned, see b2).
3. /tmp/fix-master worktree cleanup (still exists, holding the uncommitted repair).
4. CHANGELOG entries for the two Dependabot bump merges (PR16/PR18) under [Unreleased] Changed.

## d) TOTALLY FUCKED UP (honest ledger)

1. **THE BIG ONE — I broke master by committing a co-edited file without hunk-level audit.** I ran `git add pkg/cmdguard/v4/cli.go` while the parallel session was mid-refactor on the SAME file: my commit 8babb0e→887dcae swept their 5 config-path hunks (caller side: 3-var `loadConfigFileOrSkip` call, `configFilePath`, `configSetFields`, `strictConfigKeys`, `configExplainCmd`) without the callee side (still uncommitted in config.go/config_file.go). Master CI went RED on build ("assignment mismatch: 3 variables but 2 values", cli.go:253) and stayed red ~90+ minutes. My audit was `git diff --stat` (117 lines "looked plausible") — I never opened the hunks of a file I KNEW was shared. This repeats the previous session's `git add -A` incident in a subtler form: I knew the rule and applied it only to file names, not to hunks.
2. **Wrong first reset designs (caught by tests, cost 3 iterations):** (a) blanket `Set(DefValue)` corrupts stringSlice values into a literal `"[]"` element (pflag CSV cannot re-read its own bracket form); (b) re-registering via pflag constructors PANICS on flag redefinition (pflag forbids it). Final scratch-FlagSet Value transplant is correct — and now documented.
3. **Small wasted cycles:** wrote `f.Short` (pflag field is `Shorthand`) — one build failure; wrote a nonexistent `cobraRunArgs` type in the first test draft — sed-fixed.
4. **Misattribution delay:** the first gate failure I read as "parallel session's WIP noise"; the final master red was MINE. I connected the dots only when master CI (not the local tree) failed on a line I did not remember writing.

## e) WHAT WE SHOULD IMPROVE

1. **Hunk-audit rule for shared trees:** never `git add` a file another session co-edits without reviewing every hunk; prefer building commits in an isolated worktree (`git worktree add` from a clean ref + apply ONLY your patch). I did this for pushes but built the original commit in the contended tree — that is where the sweep happened.
2. **Worktree-first for every commit in multi-session repos** (not just pushes/PRs).
3. **Record pflag traps once, correctly** (done: AGENTS.md gotcha — append-on-Set, unreadable bracket defaults, redefinition panic, transplant pattern).
4. **State verification scope explicitly** when the full gate is impossible (shared tree, load): scoped build+test+lint + post-push CI, and say so.
5. **Batch heavy verification when load > ~150**; starved gates waste 90 minutes and produce no signal.

## f) Next tasks (ranked, not padded to 50)

1. Push the /tmp/fix-master repair (re-verify under lower load or trust the 5-line revert + CI) and confirm master CI green.
2. Reconcile local master (stash → rebase-skip → pop) and re-apply the parallel session's cli.go WIP as unstaged patch (see b2) so their session is whole.
3. Remove /tmp/fix-master worktree.
4. Cut v4.2.1 (CHANGELOG [4.2.1] from Unreleased Fixed; tag; proxy + scratch-consumer verify; GitHub release) per go-release skill.
5. After v4.2.1: bump go-credentials-store, drop its max(ExitCode,1) workaround.
6. Add CHANGELOG [Unreleased] entries for PR16/PR18 dependency bumps.
7. lint/v0.2.0 release train (pre-existing P6: restore real pin + version-scoped replace before tagging).
8. Pareto plan M09/M10/M11 scheduling (v5 scaffold at M11 per owner decision).
9. Test: --output flag across repeated Executes (value + spec.outputFormat self-heal path).
10. Test: hidden stringSlice flag reset preserves Hidden/Deprecated through the transplant.
11. Close the documented caveat: extend transplant reset to other appending kinds (intSlice etc.) for raw cobra escape-hatch flags.
12. Benchmark resetFlagStateForReuse overhead (tree walk per Execute).
13. Check docs/ERROR_REFERENCE.md covers the ExitCode-zero rule; extend if not.
14. Verify the parallel session's eventual commit of the config-truth workstream re-adds the cli.go hunks cleanly (no duplicate field declarations).
15. Confirm the untracked release status report (docs/status/2026-11-02_04-15_*) lands via the daemon.
16. Consider branch protection / required CI on master so direct pushes cannot skip verification (this session pushed to master from worktrees).
17. Revisit Dependabot config: root + 4 sub-modules needed coordinated settles; consider a post-merge auto-tidy workflow so submodule-smoke does not gate on manual completion.
18. Consider a tiny "shared-tree commit protocol" note in AGENTS.md (hunk-audit rule from e1) once the parallel session lands.

## g) Questions I cannot answer myself

1. **Push the master repair now, or wait for local verification?** Master is red either way; the repair is an exact 5-hunk revert of foreign lines (rg-verified unreferenced), but its build+test run is starved (load 224). Push-now lets CI verify; wait-now burns more red-master time. Which do you want?
2. **Who reconciles the parallel session's cli.go WIP?** I can re-apply their hunks as unstaged changes after the local rebase (b2), or leave the whole reconciliation to that session/daemon. Doing it myself risks colliding with their in-flight edits; leaving it risks their session silently missing hunks. Preference?
3. **Release vehicle for #11/#13: v4.2.1 patch or fold into a larger train?** #11 is a behavior change of a public function (arguably patch — it fixes exit-0-on-failure); consumers are bitten today. v4.2.1 now, or bundle with lint/v0.2.0 + v4.3.0 later?

---

_State at report time: origin/master RED at 71f7ecf (my sweep; repair staged in /tmp/fix-master uncommitted), local master diverged 2/5 with parallel session's WIP unstaged, all 5 GitHub items closed/merged, machine load 224._
