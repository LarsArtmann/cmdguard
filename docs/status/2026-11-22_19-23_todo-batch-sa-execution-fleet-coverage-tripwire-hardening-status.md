# Status: TODO SA-batch execution — fleet coverage truth + tripwire hardening (2026-11-22 19:23)

**Scope line.** VERIFIED this session (fresh command output): full root `-race` suite green; all 7
sub-modules build + test green (incl. the new tripwires in `flightrecorder` and `lint`);
`golangci-lint run ./...` **0 issues** after fixing 35 findings in my own new code; `docs-gates run .`
5/5 PASS (run twice); fan-out guard green (54 symlinks, cmdguard skill resolves to this repo);
ci.yml YAML-validated and the NEW coverage+badge step simulated end-to-end locally (fleet-total
86.1%, sanity greps pass); true fleet coverage measured three ways (root-only 86.5%,
workspace-`./...` 86.6% which coincidentally matched the committed badge, explicit per-module
patterns 86.1% = the real fleet number). NOT verified: the new CI coverage step on REAL CI (only
local simulation — `go work init` + first workspace test run are unproven on a runner);
`nix run .#check-all` as a single gate (components run individually — same scope call as the
17-40 session); website gates (no website content changed); `erraudit fix`/`--type-aware` sweeps;
whether the next master push regenerates the badge SVG with the new alt-text (it is CI-generated
and deliberately never hand-edited).

## Session summary

Input: TODO_LIST "Sentinel & erraudit follow-ups" (SA1–SA10) harvested from the 17-40 report, plus
a generic execution directive. Biggest moment: **SA10's premise was a documented lie.** AGENTS.md
recorded (owner decision 2026-11-21) that the CI badge "deliberately reports the FLEET-TOTAL across
all 8 workspace modules" — but CI has no go.work and `./...` never descends into nested modules, so
the badge silently measured ROOT-MODULE-ONLY coverage. Verified with three local measurements and
`go list ./...` (root module has exactly 8 packages; zero sub-module lines in the profile). Fixed
the implementation to match the decision instead of writing alt-text that would have encoded the
lie deeper. Second moment: **the lint-module tripwire fired on a real judgment call** —
`ErrBaselineNotFound` flagged dead, `rg` proved it wired cross-package (`lint/cmd/cmdguard-lint/main.go:186`),
so the scanner grew an exported-sentinel qualified-reference rule (and the fixtures were updated to
pin both directions).

## a) FULLY DONE

1. **SA10 (CI + README)** — CI coverage job now creates the 8-module workspace
   (`go work init . ./flightrecorder …`) and patterns every sub-module explicitly; badge
   `aria-label`/`title` → "coverage (all modules)"; README badge alt-text updated
   (`.github/workflows/ci.yml:50-63`, `README.md:8`). Step simulated locally with the exact
   script body: fleet 86.1%, both sanity greps pass. Commit `81bcc5a`.
2. **SA2 (shared scanner + fixtures)** — `pkg/testutil.FindDeadSentinels` (AST scan, per-package
   reference scope, allow-file, test-file exclusion, exported-qualified-ref rule) +
   `dead_sentinels_test.go` with 12 temp-dir fixtures covering detection AND acceptance AND
   broken-source. The fixtures caught three of my own scanner bugs before they could land
   (§d1–d3). Commits `7a6d444`, `c71f028`, `1e5e706`.
3. **SA1 (sub-module tripwires)** — `TestNoDeadSentinels` now runs in `flightrecorder/` and
   `lint/` (covering `provider/` + `cmd/cmdguard-lint/` recursively); both green; the lint run
   surfaced and fixed the exported-sentinel cross-package rule (§summary). Commit `c71f028`.
4. **SA4 (SKILL↔reference count tripwire)** — `TestSkillSentinelCountMatchesReference`
   (`errors_reference_test.go`) compares EVERY sentinel count stated in SKILL.md against the LIVE
   `GenerateErrorReference()` output (immune to doc staleness); negative direction probe-verified
   in `/tmp` (doctored counts → drift detected). Commit `1e5e706`.
5. **SA5 (tree-vs-lint why)** — AGENTS.md erraudit bullet now states the mechanism (tree runner
   executes the enforcement analyzers `--max` enables; lint runner drops them) with the
   verified-zero note. Commit `c1e383a`.
6. **SA7 (fan-out guard + ritual)** — guard green (54 symlinks; note: the script lives in the
   crush-config checkout, NOT in this repo — the TODO row's in-repo path was wrong); SKILL.md
   fan-out check folded into AGENTS.md "Living-doc edit ritual". Commit `c1e383a`.
7. **Living-doc ritual followed this time** — `docs-freshness --pre-flight` ran SAFE on all
   candidate lines BEFORE any living-doc edit (17-40 §d2's ordering violation not repeated);
   6 strikes via safe-strike (unique anchors, gate CLEAN after each apply); SA11 row added for
   the newly-found ERROR_REFERENCE v4-intro drift; CHANGELOG `[Unreleased]` carries 2 Added +
   1 Changed bullets. Battery 5/5 PASS after the edits.
8. **Lint-zero restored** — all 35 findings in my new code fixed (scanner rewritten around a
   struct to bring cyclop 16→≤15, exhaustruct full literals, wrapcheck `%w` wrapping, wsl_v5/
   nlreturn whitespace, golines on the test file). Final `golangci-lint run ./...` exit 0.

## b) PARTIALLY DONE

1. **SA10's CI leg is simulation-verified only.** The `go work init` + workspace coverage path has
   never run on a real runner: first-use dependency downloads for sub-modules, `go.work.sum`
   materialization, and the 300s timeout with the non-parallel flightrecorder singleton are
   unproven on CI. The next master push is the proof; nothing currently watches for it (rowed as
   SA12).
2. **Release-shape debt created knowingly:** sub-module go.mod pins (`cmdguard/v5 v5.1.0`) do not
   contain `FindDeadSentinels` — until the next root tag + sub-module pin bump (the documented
   release dance), a GOWORK=off test build of the sub-modules against PUBLISHED pins would fail.
   Workspace/dev builds are unaffected. Release-time concern, not a dev blocker.
3. **Badge SVG itself** still carries the old alt-text until CI regenerates on the next master
   push (by design — the badge is never hand-edited).
4. **Cadence gate aftermath:** this report is the 24th unarchived report (threshold 23) — the
   archive ritual (D3) runs at the end of THIS session; not done at the moment of writing.
5. **One observed-not-chased discrepancy:** docs-cadence printed "0 struck DONE row(s) in
   TODO_LIST.md" immediately after I struck 6 rows. Gate PASS, but the counter and reality
   disagree; I did not chase what it counts (rowed as SA13).

## c) NOT STARTED (noticed, deliberately out of scope)

1. SA3 / SA8 (owner decisions, pending from the previous session) and SA6 / SA9 (sweep work not
   attempted).
2. `nolint-audit` per module; erraudit as a real CI gate (SA3-adjacent).
3. Website gates (`check-guide-sync`, astro check, smoke) — no website-facing content changed.
4. The 286-inline-wrap-site audit for `errors.Is`-chain consistency (17-40 §f8, untouched).
5. taskctl-audit.html writer (untracked + gitignored by a prior session; not re-checked).

## d) TOTALLY FUCKED UP

Nothing broke — every gate is green at writing time. Honest process failures:

1. **I almost trusted a coincidence.** My first workspace-mode run returned 86.6% — exactly the
   committed badge — which superficially CONFIRMED the "badge = fleet-total" doc claim. Only the
   GOWORK=off run (86.5%) and inspecting the profile (zero sub-module packages) exposed the truth.
   The near-miss class: a verification that accidentally validates the bug.
2. **False "LINT-EXIT=0".** I piped golangci-lint through `tail` and echoed `$?` — which is
   tail's exit code, not the linter's. I printed "LINT-EXIT=0" directly under a findings list.
   The repo shipped with 35 lint findings in my code until a deliberate re-run with true exit
   caught it. Pipe-masking an exit code is exactly how the fleet describes "pipeline masking"
   failures; I produced one live.
3. **My Go style needed 35 corrections to meet repo standards** (21 wsl_v5, 8 nlreturn, 4
   wrapcheck, 1 cyclop, 1 exhaustruct, 1 golines). I know this repo's whitespace discipline from
   its AGENTS.md and still wrote fleet-default style first. First-time-right would have saved a
   full rewrite pass.
4. **A careless multiedit silently deleted the `_test.go` skip AND the parse-error handling** in
   the scanner (old_string spanned lines I didn't intend to touch). Found by failing tests
   (`errInnerError` flagged from a test file), not by reviewing my own diff. Also: one mid-refactor
   state had `skipSentinelTree` with a signature its own call site didn't use, a bogus half-fix
   (`filepath.Clean(".")` comparison that breaks absolute roots), and a typo (`wira`→`wire`) —
   all caught in-loop, all mine.
5. **Two safe-strike usage errors** (`--dry-run` isn't a flag — dry-run is the default; single-
   strike takes anchor and verdict as separate args, not `ANCHOR=VERDICT`). Read the usage header
   the second time only after the first refusal.
6. **Repeated the 17-40 §d3 lesson within one session:** stale gopls diagnostics flagged
   `dead_sentinels.go:86:49` for hours after that line stopped existing; I kept treating the
   panel as ground truth until `go build` contradicted it. Documented failure mode, repeated
   same-day.
7. **Wasted probe round on errInnerError:** I first suspected the fixture/detection logic and
   wrote a probe test when the diff of my own multiedit would have shown the deleted skip
   immediately. `git diff` before `go test -v` next time.

## e) WHAT WE SHOULD IMPROVE

1. **Watch the next master-push CI run** for the new coverage step (SA12) — a CI-only failure
   mode (first workspace resolution on a runner) exists right now with no tripwire.
2. **Truth-audit owner-decision claims before recording them.** The fleet-total claim was
   recorded 2026-11-21 without verifying the CI implementation. Process rule: a "the tool now
   does X" sentence in AGENTS.md needs a cited command output, same standard as
   verify-external-claims for external tools.
3. **Verify the verify:** pipe-exit-codes, simulated-equal-real, "PASS" counters that disagree
   with observable state (SA13) — this session produced three near-misses of the same class:
   trusting a green-looking signal over ground truth.
4. **Write Go in repo style first time:** blank-line-before-return discipline (wsl_v5/nlreturn),
   wrapped external errors, full struct literals. A personal pre-flight checklist would have
   saved the rewrite.
5. **After any multiedit, `git diff` the file before testing** — the deleted-skip incident was
   visible in the diff instantly.
6. **Safe-strike ergonomics:** two usage refusals per session (also reported 2026-11-21 §3). A
   `--examples` flag or richer usage text in the script would stop the repeat.
7. **Run the archive ritual BEFORE the battery goes red** — I'm doing it reactively at session
   end; the D3 rule ("when the 24th report lands") is known, so the ritual should be scheduled
   the moment a report is drafted.
8. **Check-all as single closing gate** remains unrun two sessions in a row (17-40 and this one);
   either run it or retire the AGENTS.md suggestion of it as the closing gate.

## f) UP TO 50 THINGS TO GET DONE NEXT (brainstorm — most are ROADMAP fuel)

Sentinel & tripwire system:

1. SA12: verify the fleet-coverage CI run green on the next master push (badge regenerates with
   "coverage (all modules)" alt-text; `go work init` works on the runner).
2. SA13: chase the docs-cadence "0 struck DONE row(s)" counter discrepancy (6 rows struck,
   counter said 0; gate PASS — what does it actually count?).
3. SA11: rewrite ERROR_REFERENCE.md's ~95-line v4 intro prose (owner question pending: v5 rewrite
   vs dual-major wording — see §g1).
4. Consider testutil-izing the SKILL-count tripwire's regexes if another repo adopts the pattern
   (fleet-level idea; cmdguard-local today).
5. Extend `FindDeadSentinels` fixtures: a sentinel referenced ONLY inside the allow-file of a
   DIFFERENT package (currently only same-dir allow-file is pinned).
6. Decide whether `wiredSentinels`' double-parse (walk + reference scan) should cache parsed
   files for large trees (perf; irrelevant at current repo size — YAGNI candidate).
7. Fixture for a grouped var block where ONE var of the group is dead (partially-wired group).
8. Fixture proving build-tag-guarded files (`//go:build ignore`) are scanned (they are — no tag
   filtering — is that desired?).
9. SA3/SA17 (from 17-40): machine-artifact the erraudit acceptance posture.
10. SA6: the `erraudit fix` dry-run + `--type-aware` sweep across 8 modules.
11. SA8: owner decision on release vehicle for the 5 removed sentinels (v5.2.0 vs v6).
12. SA9: cmdguard-lint doc for the `var cli CLI = <constructor>` unsupported form.

CI & release:

13. After SA12's green run: decide whether `submodule-smoke.yml`'s test legs are now redundant
    with the fleet-coverage job (owner question §g2) — dedupe or keep as independent signal.
14. Release dance at next tag: bump sub-module pins to the root version carrying
    `FindDeadSentinels` BEFORE tagging any sub-module (GOWORK=off test builds would otherwise
    fail against published pins).
15. Consider pinning the erraudit invocation in ci.yml if it ever becomes a gate (17-40 §16).
16. `nix run .#check-all` as a single closing gate, twice deferred — run it once or remove the
    AGENTS.md implication that it's the standard.
17. Watch CI badge-bot race on the next push (the coverage job now runs longer; push conflicts
    with other bots get the documented retry, but the timing shifted).

Docs & truth surfaces:

18. SA11 execution (after §g1 answer).
19. Sweep for other "the CI does X" claims in AGENTS.md and verify each against the workflow
    source (the fleet-total lie was one of potentially several).
20. Confirm `website/public/llms.txt` picks up nothing wrong from this session (it embeds counts?
    17-40 §f24 flagged regeneration — re-check on next website build).
21. Verify pkg.go.dev renders the new testutil godoc after the next root tag.
22. Keep SA-row harvest cadence: route actionable §f items at row time (SA12/SA13 done this
    session), park the brainstorm with an explicit struck harvest note (done this session).

Testing:

23. Property-style fixture: generate random mini-packages (dead/wired mixes) and assert scanner
    agreement with a brute-force reference (overkill today; PARK).
24. Add the tripwire to the remaining 5 sub-modules (glamour/prompts/spinner/telemetry/schema)
    if they ever grow sentinels — the helper makes it a 10-line test; do it when the first
    sentinel appears, not speculatively (YAGNI guard).
25. Benchmark the scanner on the repo tree (currently ~1s; only matters if fleet-adopted).

Process:

26. Close the loop on safe-strike usage text (file upstream or fix — BuildFlow rule: never
    absorb).
27. Record the pipe-exit-code lesson in the project's AGENTS.md gotchas if it recurs (once is a
    typo, twice is a pattern).

Deliberately parked (not rows): items 3–8, 23–25 above are fixture/perf ideas worth keeping in
this §f as fuel, not TODO_LIST rows.

## g) Questions I cannot figure out myself

1. **SA11 semantics:** ERROR_REFERENCE.md's hand-written intro (title, facts path, usage examples)
   says **v4** while the generated table and the repo are v5. Rewrite the intro for v5, or keep
   v4 wording deliberately (e.g. because `release/v4` consumers read master docs)?
2. **CI redundancy policy:** once the fleet-coverage job is proven green, do you want
   `submodule-smoke.yml` trimmed (its test legs now run inside the coverage job) or kept as an
   independent, workspace-free signal?
3. **SA8 (carried):** release vehicle for the 5 removed sentinel identifiers — v5.2.0 minor with
   an explicit API-surface-correction note, or hold for v6?

---

_Point-in-time snapshot. Section (f) is HARVEST input — do not let it die here. Per harness
contract, not committed manually; the auto-commit daemon will pick this file up._
