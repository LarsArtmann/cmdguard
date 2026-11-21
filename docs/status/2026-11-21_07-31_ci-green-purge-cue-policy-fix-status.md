# Status: CI fully green, TODO purge executed, cue-policy reusable-call fix landed

**Session:** 2026-11-21, ~03:05–07:30 CET (continuation of the 2026-11-20 fleet-audit + Pareto sessions)
**Repo:** cmdguard master @ 6d206c9, in sync with origin; cue-policy @ c1482c8, pushed
**Battery:** docs-gates 5/5 PASS; check-all EXIT=0; guide-sync OK; CI 7/7 success (runs 44853943250, 44855848949, 44856412886, 44857994433); Sub-module Smoke success; Website green

## a) FULLY DONE

1. **First fully green CI runs (SU24/SU38)** — run 44853943250: all 7 jobs green (test, benchmark, nix-check, branded-ID gate, dogfood, lint, gowork-off). The badge bot commit c902c78 landed on master (job-scoped `permissions: contents: write` fixed the 403 from run 44845752228); retry guards (`git diff --cached --quiet ||`) in place. Subsequent runs 44855848949 (lint/rules.go fix), 44856412886, 44857994433 all 7/7.
2. **Coverage badge live** — raw URL resolves and renders (86.3% on the first refresh, 86.2% after the latest run; the badge is the live source, as SU36 stamped). README Coverage badge points at it.
3. **SU41 (sibling pushes) verified + struck** — SKILLS adf2e80 + 7d14347 and templ-components bc369053 confirmed on origin (the prior session's "daemons pushed" claim verified true for the SU41 commits; 2 unrelated later daemon commits in SKILLS remain unpushed — their daemon's business).
4. **SU39 fresh check-all green** — EXIT=0 after two real fixes (see d/e): lint/rules.go treeWalkingRule construction fix (exhaustruct_v5 toolchain drift) and AGENTS.md dep-budget 14→15 (go-branded-id direct since 7f7adfa for BD003).
5. **SU43 purge executed** — `docs-gates purge --apply` folded 47 struck rows into one pre-flighted CHANGELOG entry; battery green after honest source annotations (below).
6. **Dissolved-harvest repair (6 annotations)** — the purge dissolved accidental harvest bindings; items annotated done at their sources with recovered evidence: 2026-08-06 report item 40 (nix badge — covered by the CI badge, nix-check job), 2026-11-11_01-33 item 16 (implemented for real: check-docs-api.py check #5 fails any guide missing "## Where to go next"; all 20 guides comply; negative probe verified) + item 23 (machine-contract real describe bytes), 2026-11-11_14-57 items 21/32/33/34/35 (W1/FM2–FM5 verdicts recovered from pre-purge TODO_LIST).
7. **SU47 harvested** — report item 20 (SKILL deeper-truth site links, plan T20) is genuinely open; pre-flighted and inserted as TODO_LIST row SU47 (binds item 20).
8. **01-56 report §f items 1–6 annotated done** with this session's evidence (run IDs, commit hashes).
9. **T27/SU45 cue-policy fix LANDED AND PUSHED (c1482c8)** — job schemas are now a discriminated disjunction: run job (runs-on + steps required, job-level `uses` = bottom) or reusable-call job (`uses!` required, `steps` = bottom). Root-caused two CUE semantics first: `never` is not a predeclared identifier in CUE v0.17.1, and plain disjunctions with unsatisfied-required branches stay incomplete-and-unresolved — the bottom discriminators are load-bearing. Verified: cmdguard ci.yml + release.yml vet clean (incl. jobs.lint / call-lint call jobs), BuildFlow's ci.yml passes, a garbage job is still rejected. **The pre-commit hook now passes WITHOUT --no-verify** (the 166s hook run completed clean; the commit it belonged to died on a daemon race, not on gates).
10. **SU45 struck** (strike applied; see d for the shape damage).
11. Three consecutive pushed CI cycles green end-to-end; every verification gate that ran this session is green at report time.

## b) PARTIALLY DONE

1. **SU45 strike shape** — the whole-row strike was re-wrapped by the formatter (BuildFlow pre-commit runs with --fix; dprint re-aligns tables) into a PARTIAL row: first cell struck-empty (`~~ ~~`), task/notes/status cells bare. Row 158 reads open although SU45 is done with verdict text lost. Gates stay green (they don't classify this TODO_LIST table shape as PARTIAL — a gate gap). Fix path: `annotate.go repair-partial --dry-run` + set-cell verdict, then check-rows (see f1).
2. **Session commit messages** — the daemon swept nearly every change set into `chore: auto-commit` commits, and my one plain (hook-proven) commit attempt died on a HEAD lock race (expected abb8c59, daemon moved to ea5e87d during the 166s hook). The SU45 strike's detailed commit message never landed; evidence survives in the row verdict text and this report only. The SU43 strike commit (abb8c59) and the lint fix (7e0829e, daemon) carry theirs.
3. **Badge number reconciliation** — AGENTS.md says ~89.3% core-only (date-stamped, badge = live source); live badge is 86.2% root-total. Correct per D15 conventions, but the two numbers sitting 3.1pp apart invites misreading (see f12).

## c) NOT STARTED (unchanged open work)

- **SU40** (templ-components: fix or baseline the 55 pre-existing findings so it commits without --no-verify) — not touched.
- **SU42** (annotate.go per-row TSV validation with row-numbered errors; repo: docs-gates) — not touched.
- **SU46** (promote pre-push workflow smoke to a bounded task) — row exists; work not started.
- **SU47** (SKILL deeper-truth site links; plan T20/T20.1) — harvested this session; work not started.
- **Plan T8+** (T11 WithOutputFormat wiring, T9 release cut v5.1.1/v5.2.0, T21/T22 since-column + CG009 evaluation, T23 live-site smoke, T24 config-explain parity, …) — untouched; the plan file owns the order.

## d) TOTALLY FUCKED UP (honest)

1. **Formatter-damaged strike (b1)** — my own doing: I struck, battery'd, and then ran a formatter-carrying commit. The docs-health SKILL.md warns verbatim: "After any `dprint fmt`, re-run `check-rows` … the pass needed two repair cycles because the battery was trusted from before the format step." I repeated the documented failure mode on the same day I read the warning.
2. **Commit/daemon race** — launching a 166s hook commit while the auto-commit daemon is live guarantees a lock race. The commit failed (`cannot lock ref 'HEAD'`); the daemon's auto-commit carried the change with a throwaway message. Lesson: for changes that need a real message, either strike+format+check+commit in one tight window, or stage the message into the daemon sweep window knowingly.
3. **Purge evidence loss (minor)** — the purge's CHANGELOG entry carries row IDs only (`TODO_LIST purge (2026-11-21): SU24` …), not the struck verdicts (run IDs, hashes). Evidence survives in git history (TODO_LIST @ e5dbb9b) and in the annotated status reports — but a future reader of CHANGELOG alone gets IDs without proof. The tool's design; my miss was not pre-annotating the reports that carried the evidence BEFORE purging (I did it for the dissolved harvests, after the purge refused — that part worked as designed).

## e) WHAT WE SHOULD IMPROVE (incl. "what did I forget / do better")

1. **Run check-rows after EVERY formatter touch** — the formatter-order canon exists precisely because gates trusted pre-format state go stale. Concretely: strike → format → check-rows → battery → commit, never strike → battery → commit.
2. **check-rows coverage gap for TODO_LIST** — the battery's annotate gate passed while a struck row was half-unstruck by the formatter. A TODO_LIST-specific row-shape check (struck rows must have no bare cells) would have caught b1 red. New gate needs a red it would have caught — b1 IS that red.
3. **Long-hook + daemon racing** — either a pre-commit "check tree unchanged before writing HEAD" retry, or a session discipline: run `buildflow --build-mode pre-commit --staged-only` manually BEFORE `git commit`, so the hook's second run is cache-warm and fast (result cache hit ≈ instant), shrinking the race window to milliseconds.
4. **CUE policy versioning note** — the cue-policy fix relies on bottom-literal discriminators (`steps?: _|_`) because `never` doesn't exist in CUE v0.17.1; if a future CUE adds `never`, migrate for readability. A one-line comment in cue-policy (present) should also land in the cue-policy README/CHANGELOG (not done — repo has its own doc conventions; flagged here).
5. **Two CI truths for coverage** — badge (root-total 86.2%) vs AGENTS (~89.3% core-only). Both "correct"; the split is documented. Improvement: make the badge's SVG title or the README badge alt-text say "coverage (all modules)" so the tiering is self-describing.

## f) Up to 50 things we should get done next (ordered by impact)

1. Repair SU45's partial struck row (repair-partial + set-cell verdict `done (cue-policy c1482c8 …)`); re-run check-rows + battery (≤12 min).
2. Add the check-rows-shape gate for TODO_LIST struck rows (red = b1) to docs-gates (30–60 min; repo: docs-gates).
3. Re-land the SU45 commit message loss: fold the evidence into the next real commit's message or amend the row verdict (5 min; bundle with f1).
4. SU47: add website-guide links to SKILL.md "Where deeper truth lives" (plan T20.1, 10 min) — closes the newest TODO row.
5. SU44: describe-fixture byte guard for machine-contract.mdx (plan T7, 30 min) — prevents describe-excerpt rot.
6. SU42: annotate.go per-row TSV validation with row-numbered errors (plan T6; repo: docs-gates, 30–60 min).
7. SU40: templ-components — baseline or fix the 55 findings so commits stop needing --no-verify (half a session; repo: templ-components).
8. SU46: promote the pre-push workflow smoke (`go run <tool>@<pin> --help` per pinned action + job-permissions checklist) into a bounded script (plan-adjacent, 45 min).
9. Badge self-describing tiering: alt-text/title "coverage (all modules)" (5 min; f5).
10. README FLAGS_REFERENCE split plan (ROADMAP row; execute when the 900-line gate nears — check current line count first).
11. Live-site truth smoke after each push (ROADMAP row; scheduled-workflow candidate; pairs with plan T23).
12. Reconcile AGENTS coverage sentence with badge alt-text (f5's landing spot; 10 min).
13. Purge-entry evidence: adopt a docs-gates purge flag to carry per-row verdicts into the CHANGELOG entry (repo: docs-gates; prevents d3 recurrence).
14. cue-policy: add the bottom-discriminator rationale to its README/CHANGELOG (15 min).
15. cue-policy: sweep other fleet repos for workflows that now pass/fail under the disjunction (the SU45 audit clause — only BuildFlow spot-checked; go-output's missing `permissions:` is pre-existing and its own repo's business).
16. go-output: add the fleet-required top-level `permissions:` block to its ci.yml (30 min; repo: go-output — surfaced by the spot-check).
17. SKILLS: push or let the daemon push the 2 unrelated unpushed commits (stack-selection SKILL edit + CHANGELOG/check-skills) (5 min).
18. T9: release cut v5.1.1/v5.2.0 — fold Unreleased; re-verify SKILL claims; changelog since-tags (half a session; owner sign-off per ROADMAP).
19. T11: wire `WithOutputFormat(output.FormatJSON)` into examples/taskctl (plan T8/§f17, 30 min).
20. T21: draft the FEATURES since-column concept; T21.2 adopt/reject into ROADMAP (plan, 22 min total).
21. T22: CG009 doc-rule evaluation (new exported options without since-annotations) — scoped evaluation first (plan, 24 min).
22. T23: live-site guide-sync smoke post-deploy (plan, 45 min; pairs with f11).
23. T24: config-explain example README since-check parity sweep (plan, 30 min).
24. Add a docs-gates purge `--dry-run` warning when a struck row's verdict text carries distinctive tokens (early warning for dissolved harvests; repo: docs-gates).
25. Sub-module test counts: derive mechanically (ROADMAP row) and stamp AGENTS (30–45 min).
26. STATUS pagination: docs/status unarchived count is 22 (threshold 23) — the 2026-10-10_02-01 report ages out (45d) within days; schedule the D3 archive ritual this week.
27. SU22-line reports: annotate `2026-11-21_07-31` (this file) §f items when picked up (per-item, ongoing discipline).
28. Convert the session's "three open questions" from the 01-56 report into TODO_LIST rows or answer them (badge bot identity; badge label wording) — cue-policy question is now MOOT (T27 done).
29. Badge bot identity (f28): decide cmdguard-coverage-bot vs github-actions[bot] naming; it affects commit attribution on every master push.
30. Badge label wording (f28): "coverage" vs "coverage (all modules)" vs "cov"; owner preference pending.
31. Add `nix run .#check-all` runtime badge/check to CI parity with local (currently only local; CI has no check-all equivalent job — the CI jobs cover the same ground piecemeal; decide whether a consolidated job is worth 10 min of runner time).
32. Guard the badge-commit retry: it pulls --rebase with autostash; add a final `git diff --cached --quiet` guard after the second pull (already guarded — verify in a race drill) (15 min).
33. Fleet-audit leftover: items 17–45 of the 2026-11-20_12-29 §f — refine into bounded TODO rows when picked up (ROADMAP row; ongoing).
34. templ-components gate: after SU40, delete its `--no-verify` note from the pairing commit docs (5 min, after f7).
35. annotate.go: `unwrap` for prose rows (the docs-health SKILL lists a prose-side damage class with no writer; repo: docs-gates).
36. safe-strike: warn when an anchor regex contains unescaped parens/plus (refusal UX; repo: docs-gates; bit me once this session — caught by the tool's own refusal).
37. safe-strike: auto-escape or reject shell-substituting backticks in CLI args (bit me once; the tool refused safely both times).
38. docs-gates: expose `BUILDFLOW_NO_RESULT_CACHE`-style invalidation when the POLICY repo (cue-policy) changes — the cue-vet result cache served stale findings for one run (wasted one debug cycle).
39. cmdguard-lint: CG009 prototype if T22 adopts (repo: cmdguard/lint; post-T22).
40. AGENTS.md: record the CUE disjunction-semantics lesson (bottom discriminators load-bearing; `never` absent in v0.17.1) in the cue-policy bullet (10 min).
41. Record the formatter-order lesson in this repo's AGENTS.md gotchas (strike→fmt→check-rows→battery) (10 min).
42. Record the BuildFlow result-cache invalidation gap (policy repos) upstream in BuildFlow (repo: BuildFlow).
43. Consider splitting the CI test job's badge steps into a separate needs:-chained job so master pushes don't serialize the whole test job on git operations (design question; 30 min to spec).
44. Schedule the archive ritual for the 4 fleet-audit reports' unarchived siblings before cadence threshold trips (see f26).
45. ROADMAP: add "CUE policy: migrate `_|_` discriminators to `never` if/when CUE adds it" (2 min).
46. ROADMAP: add "docs-gates purge: verdict-preserving CHANGELOG entries" (2 min).
47. Sweep `--no-verify` mentions in fleet docs (AGENTS.md pre-commit section, lint pairing notes, templ-components skill) — several now claim the --no-verify era is current; it ended (30 min).
48. Verify the 86.3→86.2 badge drift cause (docs_tables_drift_test.go additions changed the root-module denominator; confirm and note in AGENTS if stable) (15 min).
49. Add a CI job (or extend dogfood) to run `scripts/check-guide-sync.sh` — currently a Quick Start command + check-all step; CI doesn't run it (verify claim, then 10 min).
50. Close the loop on the 01-56 report's item 11 (T11) after f19 lands — strike + annotate in one pass.

## g) Questions only the owner can answer

1. **Badge bot identity + label (f29/f30, carried over):** keep `cmdguard-coverage-bot` (self-documenting, but a "new" identity in the log) or switch to the stock `github-actions[bot]` identity — and should the badge label say plain `coverage` (current) or `coverage (all modules)`? Both are one-line changes; I won't touch them without the answer.
2. **SU40 templ-components strategy:** fix the 55 findings properly (hours, touches skill docs) or write a reasoned baseline (.cmdguard-lint-baseline.json pattern) so commits unblock now and fixes land incrementally? My recommendation: baseline now, fix incrementally — confirm or overrule.
3. **The purge CHANGELOG entry carries row IDs without verdicts (d3):** acceptable as-is (evidence lives in git + annotated reports), or should I retro-annotate the CHANGELOG entry with one-line verdicts per row (touches a frozen-ish section of CHANGELOG — your conventions, your call)?

---

*Verified at report time: docs-gates 5/5 PASS; CI 7/7 on 6d206c9; badge live at 86.2%; cue-policy c1482c8 pushed; TODO_LIST has SU45 (shape-damaged, content done), SU46/SU47 open; everything else in section a carries its own evidence line.*
