# Full Status — TODO Batch Execution: 16 Rows Retired, Byte Guards Landed, CI Format Repair (2026-11-21 13-10)

**Repo:** cmdguard master @ 3f9bdc9, pushed; CI re-running (previous red root-caused, fix landed).
**Session window:** 2026-11-21 ~07:45–13:10 CET, continuing from `2026-11-21_07-31_ci-green-purge-cue-policy-fix-status.md`.
**Gates at handoff:** docs-gates 5/5 PASS; check-guide-sync OK (6 legs); check-rows complete; pre-commit hook clean (no `--no-verify` used anywhere this session).

---

## a) FULLY DONE (verified this session)

1. **SU43/SU45 shape repair + verdict recovery** — both formatter-mangled rows rebuilt into the whole-line strike form (the active-table canon: `~~| ID | ... |~~ done (...)` — invisible to the table scanner, which is why per-cell strikes fail with `mixedTable`). SU43's verdict recovered verbatim from abb8c59; SU45's from status report §d. Purged with full verdict text into CHANGELOG (the purge tool carries rows verbatim — owner question 3 from the 07-31 report is answered for good: verdict loss does not happen when rows are purged through `docs-gates purge`; it only happened in the earlier 47-row purge because THAT tool version carried IDs only).
2. **SU47 — SKILL deeper-truth site links** — the consumer SKILL's deeper-truth section now links all six rendered mirrors (quick-start, error-handling, migrating-v4-to-v5, migrating-from-cobra, type-safe-flags, contributing, api-reference). Live-site fetch verified `/contributing/` serves; fan-out link resolves (8 site refs visible through `~/.config/crush/skills/cmdguard`).
3. **SU44 — describe excerpt byte guard (the session's flagship)** — new `check_describe_excerpt()` leg in `scripts/check-docs-api.py`: runs `taskctl describe` and requires every non-elided excerpt line to appear verbatim and in order (elision only at the guide's own `/* … */` cut marks). **The guard immediately caught real drift**: the guide claimed "Real bytes" but root flags had been alphabetized since `admin-email`/`api-url` landed, `verbose` was not the last flag, the `values` array is multi-line, and inherited flags carry no `env` field. Excerpt rewritten against live bytes; **negative probe verified** (single-byte corruption → DRIFT, restore → OK). Wired as guide-sync leg #6; runs in flake check with `GOEXPERIMENT=jsonv2`.
4. **SU48 harvested** — the byte guard exposed that taskctl embeds ANSI-escaped, lipgloss-padded help text in the machine contract's `long`/`example` fields — hostile for the contract-driving agents the guide targets. New TODO row SU48 (pre-flighted SAFE after two rewordings).
5. **Nine rows executed, struck, purged with verdicts** (commit ea583f7): SU1 (CHANGELOG `[Unreleased]` re-fed: release gate, website truth guards, fleet tooling, retirement banner, llms.txt), SU2 (FEATURES "Repository Tooling" section with file evidence), SU5 (AGENTS Quick Start gains `pre-tag-check.sh` + `check-website-smoke.py`), SU6 (full check-all EXIT=0 evidence — SU39), SU7 (**`TestDescribeContractFromBinary` + `pkg/testutil/testdata/describeapp` fixture**: real process-boundary capture test asserting schema version, command tree, env-prefix application — lint 0 issues), SU8 (`.crush/shell-output/` ignored; `global.out.css` untracked — 1566 lines of generated Tailwind), SU10 (Dependabot #21 verified CLOSED via `gh`), SU14 (`scripts/fleet-bump.sh <version>` — parameterized, arg-validated, `git mv` preserves history), SU15 (`firebase-tools@15.37.0` pinned in website.yml).
6. **SU18+SU19 — website smoke hardened + deploy freshness** — `check-website-smoke.py` rewritten: bounded retries (transient 5xx/429/timeout only, 3 attempts, backoff), URLError handling (no more tracebacks), clean FAIL summary with expected-vs-got and a debug pointer, `--expect-commit <sha>` freshness assertion backed by a build stamp (`website/scripts/stamp-deploy-meta.mjs` → `public/deploy-meta.json`, gitignored, wired into the build script). The deploy job now asserts the CDN serves THIS deploy within a bounded 90s window instead of accepting any 200. Bandit B310 on the fixed-https fetch suppressed with a placed-correctly `# nosec B310` (first attempt's `--` comment broke bandit's parser — it ate every word as a test name; second attempt skipped clean).
7. **SU3 — SKILL adoption docs** — describe section now documents `AssertDescribeContract`, `DescribeContractFromBinary`, `RegisterUpdateGoldensFlag`, and the pre-tag gate; `check-skill-fanout.sh` green (54 symlinks, no store pins).
8. **SU46 — workflow smoke promoted to a bounded script** — `scripts/check-workflow-smoke.sh`: (1) probes every `go run <module>@<pin>` in the workflows (the pins dependabot cannot see) for resolve+build+usage; (2) requires any job whose steps can push (`git push`, `gh release`, `firebase deploy`, `docker push`) to run under an explicit `permissions:` block. Wired into `pre-tag-check.sh` §7b. **Negative probe verified**: a synthetic leaky job fails the smoke; restored green.
9. **Report corpus honesty** — annotated at source: 2026-11-08_08-32 item 48 (`H4 hero-code sync guard` = `check_hero_single_source`, dissolved binding); 2026-11-11_17-42 "What was FORGOTTEN" item 5 AND §f item 5 (both = the AGENTS Quick Start gap, now done at f5602ee). Both annotations were what un-blocked the purge pre-flight.
10. **CI red root-caused and fixed** — runs on 029d19a went red (CI `lint / lint`, `nix-check`, Sub-module Smoke): all three the SAME cause — the parallel session's 08:02 daemon sweep (`5a4c840`) landed a `go-output` import in `examples/taskctl/main.go` without a hook pass, breaking gci/gofumpt/goimports. Fixed formatter-only in 3f9bdc9 (import moved into the external group), pushed; CI re-running at report time.
11. **Pushed through the daemon** — ~20 commits landed on origin across the session; every authored message either survived or its evidence lives in row verdicts / CHANGELOG purge entries / this report.

## b) PARTIALLY DONE

1. **SU19 freshness live-loop** — all mechanics landed and verified locally (stamp writes, smoke asserts, 404-on-undeployed correctly fails), but the first real `--expect-commit` pass happens on the CI deploy of 3f9bdc9. Until that deploy lands, a manual `check-website-smoke.py` run is red on `/deploy-meta.json` by design.
2. **SU4 — lint release** — verified the unreleased delta exists on lint master (rules.go treeWalkingRule/exhaustruct fix from SU39 + `go-linter-sdk v0.4.1` bump); the v0.4.0 tag, lint CHANGELOG entry, and proxy verification are NOT done.
3. **SU48 — taskctl plain Long/Example** — row harvested, investigation started (`grep lipgloss examples/taskctl/main.go` returned nothing — the styling enters through a different wiring, likely the glamour help transform path or command construction; the real source is not yet located). No code change yet.
4. **Bandit on the repo's python surface** — now that bandit scans changed python, `check-docs-api.py` may carry its own findings when next staged (not yet observed).

## c) NOT STARTED

1. SU9 — pingen-spike docs-gates red + staged v5.1.0 bump (repo verified present, tree clean, daemon-active).
2. SU11 — CV render-goldens determinism pin.
3. SU12 — nsfw-classifier `TestContainerHealthCheckJoinsAllFailures` single-cause `oops.Join` fix.
4. SU13 — crm missing `createDeal` test helper + contact-notes fold re-verify.
5. SU20 — GitHub-Actions secret-quoting lesson → crush-config `references/lessons.md`.
6. SU40 — templ-components 55 pre-commit-blocking findings (remediate or baseline — owner question, §g1).
7. SU42 — annotate.go per-row TSV validation (kind enum + raw-value rule) in the docs-gates repo.
8. SU16/SU17/SU21 — WIF migration, deploy-SA least-privilege, ~60s stale-window root cause (gcloud available and authed; production-gated — §g3).

## d) TOTALLY FUCKED UP (honest list)

1. **The formatter-mangle recurrence — my fault, third occurrence.** I struck 9 rows and committed BEFORE purging; the hook's dprint `--fix` re-wrapped all 9 whole-line strikes into `| ~~   | SU1 |…` partial rows and stripped the verdicts, and the daemon swept the damage in. Recovery took a python rebuild keyed on the damaged shape plus the two report annotations. **The rule is now crisp and should be mechanical**: strike → purge in the working tree → ONE commit. A struck TODO_LIST.md must never touch a hook pass. This is a safe-strike/purge tooling gap (see §e1).
2. **Per-cell strike dead end** — my first repair of SU43/SU45 used `repair-partial` (per-cell strikes, the 2026-11-05 canon for ARCHIVE-bound tables) which tripped `mixedTable` on an active table and burned several battery cycles. The two canon halves are in tension: SKILL.md's formatter-order canon says per-cell, but check-rows' mixed-table rule makes per-cell unusable in an ACTIVE table. The working practice is whole-line + immediate purge; SKILL.md does not say this.
3. **cwd trap, twice in one chain** — `cd /tmp && … docs-gates purge . --apply` ran the purge against /tmp and FATAL'd; then a later chain repeated it. Cost: two failed commands, no damage.
4. **Pre-flight token whack-a-mole** — the SU48 row wording flipped historical items three times (`guard`, `sync`, `plus` all bind), and I appended the second draft to the spec file with `>>` instead of `>` so the gate judged both lines. The honest fix was rewording; the stupid part was not reading the culprit list carefully the first time.
5. **LSP diagnostics noise** — golangci/oxlint LSP warnings repeatedly referenced lines from before my rewrites (nilfunc on a deleted test, noctx on a replaced line). I nearly chased ghosts; the real gate runs were clean. Trust `golangci-lint run`, not the editor pane.
6. **AGENTS.md size** — now 463 lines against the fleet's 220-line advisory (D13 accepts this as deliberate, but this session added Quick Start lines and the drift continues; the preflight warning grows every session).

## e) WHAT WE SHOULD IMPROVE

1. **Make strike→purge→commit atomic**: a `safe-strike.sh --purge` mode (or a wrapper) that strikes, purges, and leaves ONE commit-shaped tree, so the formatter can never see a struck row. This kills the recurring damage class mechanically.
2. **Document the active-table strike canon**: SKILL.md's formatter-order canon should state explicitly that in ACTIVE (non-ledger) TODO tables, strikes are whole-line `~~…~~` + purge same-session; per-cell is for tables being archived wholesale.
3. **Daemon hook bypass** is the root of both CI incidents (5a4c840's import, and generally any sweep touching Go files): consider a CI-side formatter gate that runs even on daemon commits (it does — that's what caught it — but the local tree should never diverge from the hook's view; a periodic `buildflow format` cron or post-sweep hook would close it).
4. **Pre-flight before writing TODO rows**: run `docs-freshness --pre-flight` on the candidate line BEFORE appending it to TODO_LIST, treating UNSAFE as a rewrite trigger, not a post-hoc discovery.
5. **The two-canon tension (per-cell vs whole-line)** should be resolved in docs-gates itself: either mixedTable exempts tables with a "todo" marker, or safe-strike grows a purge step. Otherwise every strike stays a two-step dance with a footgun in the middle.
6. **Status-report §f item counts**: the 07-31 report's "50 next tasks" were consumed down to ~8 rows; the next report should re-derive the queue from the actual TODO_LIST, not carry stale plans forward.
7. **`go run` pin inventory** is now scriptable (SU46's smoke) — the same probe could run in CI cheaply instead of only at pre-tag time.

## f) NEXT (ranked; ≤12-min slices unless noted)

1. Watch the CI run on 3f9bdc9: expect 7/7 green (the format fix covers lint / lint, nix-check treefmt, and Sub-module Smoke).
2. Verify the website deploy of 3f9bdc9 ships `/deploy-meta.json` and the deploy job's `--expect-commit` freshness leg passes (SU19's first live loop).
3. SU4: cut lint v0.4.0 — write lint/CHANGELOG entry (treeWalkingRule exhaustruct fix, go-linter-sdk v0.4.1), verify release shape (real pin + version-scoped replace), tag, push, proxy check.
4. SU48: locate where taskctl's Long/Example pick up styling (check command construction and the glamour help-transform path), switch to plain text + style-at-render, re-run the byte guard (it must stay green — the contract bytes will change, so regenerate any goldens).
5. SU42: annotate.go per-row TSV validation (docs-gates repo): kind enum check + raw-value rule, row-numbered errors, self-test cases for the 2026-11-20_20-15 §f5 failure shape.
6. SU40: templ-components — if §g1 answer is "baseline": write `.cmdguard-lint-baseline.json` via `--write-baseline`, verify a clean commit lands WITHOUT `--no-verify`, document the ratchet in that repo's AGENTS.md.
7. SU20: write the secret-quoting lesson to crush-config `references/lessons.md` (the 2026-11-21_01-56 §f14 evidence: unquoted token in a `run:` block; commit there, let its fan-out guards verify).
8. SU9: pingen-spike — run docs-gates there, fix the red, land the staged v5.1.0 bump.
9. SU12: nsfw-classifier — reproduce the test failure, apply the single-cause `oops.Join` fix, run that repo's battery.
10. SU13: crm — restore the `createDeal` test helper, re-run the contact-notes fold tests.
11. SU11: CV — pin render-goldens determinism (fixture order, stable ×3 runs).
12. §e1: implement strike→purge→commit atomicity in safe-strike (docs-gates repo).
13. §e2: SKILL.md canon amendment (active-table strike form) — one paragraph, battery-checked.
14. §e3: propose the post-sweep formatter catch-up (crush-config or BuildFlow territory — owner-gated, §g3 adjacent).
15. SU16 assessment only: `gcloud iam` list the deploy SA's roles + key age; produce the WIF migration checklist; DO NOT rotate anything without §g3.
16. SU21 assessment only: pull Firebase Hosting release history + the deploy job timings around the ~60s window; produce the evidence table; fix only if the cause is config-side.
17. Re-derive TODO_LIST: after SU4/SU48/SU9/SU11/SU12/SU13/SU20/SU40/SU42 land, the active table is empty — harvest the next batch from docs/status (the 2026-11-08_08-32 §f remaining items are the natural source).
18. Badge question (§g2): if the owner wants the README badge to show core-only, switch the coverage job's badge payload to the core package total and hand-edit the label once; if fleet-total stays, no action.
19. Add `check-website-smoke.py --expect-commit` to the local deploy runbook line in AGENTS.md Quick Start (one line).
20. Consider promoting the workflow smoke into CI (a fast job after lint) so rotten pins surface on PRs, not only pre-tag.
21. docs-gates: the `docs-freshness` culprit output could name the REWRITE suggestion (it names tokens; a "try replacing X with Y" hint would save the whack-a-mole).
22. annotate.go: the `repair-partial` docstring still claims per-cell is the fix for PARTIAL rows — update it to the active-table truth once §e2 lands.
23. taskctl-audit.html: the parallel session's stray 1585-line audit export sits tracked in the repo root — ask the owner whether to untrack it (I did not touch it: not mine, possibly their evidence).
24. examples/taskctl: the `help` flag appears first in describe output while docs imply flag order follows the struct — decide whether describe should pin struct order or document the current order (contract-facing question).
25. KEEP-GOING candidate: the 2026-11-08_08-32 report's remaining §f items (H-series) are mostly ≤30-min slices; harvest them into TODO_LIST after the current rows drain.

## g) THREE QUESTIONS I CANNOT ANSWER MYSELF

1. **SU40 templ-components — baseline or remediate?** Baseline (`--write-baseline`, ratchet forward) unblocks commits today and defers 55 findings (50 file-size + 5 docs-gates); remediation is a multi-hour restructure of files that are large by design or not. Which do you want? (If baseline: should the ratchet pin be committed as a visible debt marker, or kept local-only?)
2. **Badge identity + label** — the coverage badge is committed by the badge bot (`github-actions[bot]`, `chore: refresh coverage badge from CI run [skip ci]`) and currently shows the root-module total (~86.2%), while AGENTS.md quotes core-only ~89.3% date-stamped. Keep fleet-total on the badge (my recommendation — one number, live) or switch the badge payload to core-only to match AGENTS? And is the label wording ("coverage | 86.2%") what you want, or "core coverage | 89.3%"?
3. **Production GCP authority for SU16/SU17/SU21** — gcloud is authed as lars@helpless.ai on this machine. The WIF migration (retiring the FIREBASE_SERVICE_ACCOUNT SA key) and deploy-SA least-privilege are irreversible-leaning IAM changes on the lars-software project. Do I proceed autonomously (I'd draft, apply in a maintenance window, and verify via the freshness check), or prepare-only with you executing? Related: is `taskctl-audit.html` in the repo root yours to keep (committed by the 08:02 daemon sweep of a parallel session), or should I untrack it?

---

*Verified at report time: docs-gates 5/5 PASS; check-guide-sync 6 legs OK; check-rows complete; fan-out 54 links OK; CI re-running on 3f9bdc9 with the format fix (previous red fully root-caused: parallel-session import, one file, three jobs).*
