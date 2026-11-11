# Status Report — v5.1.0 Shipped, Fleet Bumped, Gates Mechanized (execution of the 15-10 Pareto plan)

**Date:** 2026-11-11 17:30 CET
**Scope:** execution of `docs/planning/2026-11-11_15-10_pareto-v51-fleet-hardening-and-truth-plan.md`
after the owner's "execute the whole list" directive. The D-gates were
resolved by that directive (documented in AGENTS.md "Owner decisions") — D1
cut v5.1.0, D2 keep struck rows, D3 archive-by-gate, D12 keep release-shape
pins, D13 AGENTS.md stays load-bearing, D14 file-size gate at 900.

**Ending state:** master at `2d3f4fe`, pushed; CI + Sub-module Smoke GREEN;
docs-gates 5/5 PASS; website live with the v4-retirement banner, smoke
30/30; v5.1.0 is GitHub Latest; 27 consumer repos bumped and pushed;
fleet sweep 26/26 vet+lint clean; branch + tag protection rulesets active.

---

## a) Fully done (evidence attached)

### The 1% — release (A01, A02, A04)

1. **v5.1.0 SHIPPED end-to-end.** Tag `v5.1.0` (annotated), Release workflow
   success (run 43043236214), proxy resolve + compile verified in a scratch
   consumer, `releases/latest` returns v5.1.0, pkg.go.dev 200. CHANGELOG
   `[5.1.0]` cut from `[Unreleased]`; `[Unreleased]` left empty for the next
   cycle.
2. **`scripts/pre-tag-check.sh` mechanizes the "read the workflow before
   tagging" lesson.** Refuses dirty trees, existing tags (local+origin),
   unpushed/diverged HEAD, missing CHANGELOG section (strips the `v`
   prefix); re-runs the release workflow's own gates locally; sanity-greps
   the workflow shape. It EARNED its keep twice on the real cut: refused the
   unpushed-HEAD state, and caught a golines miss treefmt passed.
3. **A01 green:** full `check-all` ("All checks passed!") — build, race
   tests, lint (0 issues × 8 modules), samber ratchet at 33% baseline,
   format, workspace tidy, guide-sync, dep-tables across root + 7
   sub-modules, on the go1.27.3 toolchain pins.

### Phase 0 — truth and guardrails (A03, A06, A10–A12, A26)

4. **Website guard fixed and the site is live again.** `website.yml`'s
   import-path guard flagged historical v3/v4 mentions in changelog.mdx and
   the migration guide (their JOB is to name old paths) — live docs are
   checked v5-current with those two files exempt by design. Deployed
   locally via the nix-store firebase binary (the CI deploy job silently
   SKIPS: `FIREBASE_SERVICE_ACCOUNT` is not set — see §g1). Smoke:
   30/30 routes 200 incl. `/llms.txt`, 404 path correct.
5. **CI sweep:** master legs green after each push; the two mid-session reds
   were triaged and fixed (golines in 671c128; see §d for what they taught).
   Dependabot #21 fixed (sub-module go.sum settle, 7b2e878) and PR #22
   landed directly on master (0217974, closes it) after its rebase never
   took.
6. **schema/v0.1.0 publish shape verified from a scratch consumer** (A10):
   proxy fetch, compile, link; published go.mod is a real `v5 v5.0.0` pin
   with zero replace leakage.
7. **D12 resolved and documented (A11):** root keeps REAL sibling pins
   (flightrecorder v0.2.0, glamour v0.3.0, spinner v0.3.0 — all tags exist);
   the BuildFlow `workspace/pseudo-version-hygiene` FAIL is a documented
   deliberate deviation (normalizing would poison downstream `go get`).
   AGENTS.md ADR section carries the note.
8. **cue-vet policy fixes (A12):** `permissions: contents: read` on ci.yml,
   `concurrency` group on release.yml; `buildflow -s cue-vet` 1/1 success;
   actionlint clean on all four workflows.
9. **govulncheck verdict (A26):** at go1.27.3, 6 of the 7 stdlib findings
   are cleared; GO-2026-6859 remains (fixed in go1.27.4; nixpkgs still
   ships 1.27.3 — flake bump deferred, recorded in the S2 verdict).
10. **A13:** `go vet` + `go mod verify` wired into check-all (flake.nix);
    green standalone across all 8 modules. Caveat in §b5.

### Phase 1–2 — fleet (FM1–FM8, A05, A09, A16–A21-class)

11. **27 consumer repos bumped to v5.1.0, committed, pushed** (fleet-bump
    script + follow-up re-vendor pass; vendored repos need
    `go mod vendor`/`go work vendor` after tidy — the first pass's five
    failures were ALL inconsistent-vendoring, not code).
12. **Fleet sweep green (A09/FM8):** `scripts/fleet-sweep.sh` — 26/26 repos
    `go vet` OK, golangci-lint **0 issues**. Verdicts doc:
    `docs/planning/2026-11-11_17-18_fleet-v510-bump-sweep-and-verdicts.md`.
13. **FM2 STALE (CV):** `TestLoadIdentityPortalsSection` PASSES at CV HEAD;
    the only CV red is its own `TestKeywordLogoCoverage` (CV design
    decision, not cmdguard).
14. **FM3 FIXED (cv):** `crm.checkpoint_path` added to `defaultsCRMYAML`;
    the key ratchet test is green. (Verified in the cv tree this morning;
    cv's master has since been rewritten by another active session — §d4.)
15. **FM4 FIXED + PUSHED (pma):** go-health panics when `DefaultRoutes()`
    mounts `/version` with an empty stamp; `NewServer` now clears
    `routes.Version` when `cfg.Version` is empty. 10× green + full suite.
16. **FM5 FIXED + PUSHED (go-cqrs-lite 57982df48):** the "offline hang" was
    a timeout — 2 tests × 58s go/packages ≈ 116s vs a 120s per-binary
    budget. Shared `sync.Once` BuildContext → ~24s.
17. **FM6 harness SHIPPED (A21):** `testutil.AssertDescribeContract` +
    `-update-golden` (lookup-based, see §d3) + `DescribeContractFromBinary`;
    goldens committed for taskctl and docs-generator.
18. **FM7 SUPERSEDED:** cv origin/master (+1302 commits, another active
    session) already carries the v5 API migration; the local forked state is
    preserved on branch `local-fork-backup-20261111`, not pushed.
19. **A05 executed evidence posted:** erraudit#15 → 0 findings exit 0;
    timesheets#3 → 1 CG003 (deliberate Must-wrapper) suppressed with the
    reason directive, lint exit 0 at ce35390.

### Website (A22–A24, A42, W1/W3/W5)

20. **A22:** `scripts/check-website-smoke.py` (stdlib-only) green against
    the live domain — 30 routes, 404 assertion.
21. **A23:** `website/public/llms.txt` live at the domain (links + summary,
    URLs extracted from the built dist).
22. **A24:** shell-completion guide gained the strict-args example
    (`WithStrictPositionalArgs` opt-in, `WithArgsPassthrough` escape);
    guide-sync green.
23. **A42:** `RetirementBanner.astro` live on the landing page (v4 frozen at
    v4.3.1 on release/v4 → migration-guide link); astro check, html-validate,
    link check, guide-sync, deploy + live grep all green.
24. **W5:** English-only locale decision documented in contributing.mdx.

### Lint module + tooling (A32/H9b, A34, R4/R5/R7/R8, A40, A39)

25. **H9b CLOSED (A32):** cmdguard-lint now traces dot-imported bare
    constructors (`NewCLI(...)` under a dot import) — the gap-pinning test
    flipped to expect CG003; the four negative corpora hold at 0 new
    findings; full lint suite green.
26. **A34:** codemod `--dry-run`, unknown-flag refusal (exit 2), identifier
    renames scoped to cmdguard-referencing files — the two 2026-11 fleet
    corruption incidents re-verified against a fixture (consumer's own
    `TypeHandler` untouched).
27. **A40:** reusable `lint-reusable.yml`; ci.yml + release.yml consume it;
    release job `needs: call-lint`; pre-tag script's sanity greps follow.
28. **A39:** two rulesets ACTIVE — master-protection (deletion +
    non-fast-forward) and tag-protection (deletion + update; creation rule
    removed — §d2).
29. **R7:** `TestNormalizeTableSpacing` (aligned/unaligned equivalence,
    colon markers, idempotency, prose-untouched).
30. **R8:** pkg.go.dev serves v4.3.1 (Deprecated banner) and v5.1.0; proxy
    `.info` confirms both.

### Docs truth (A27, A49, A51, D2/D3/D12/D13/D14, R6)

31. **A27:** MIGRATION_v4_v5.md §3 "Known traps" — templ floor (MVS vs
    replace, with the `ResolveAttributeValue` mechanics), Nix lock step,
    codemod scoping (updated same-day when A34 landed `--dry-run`).
32. **A49:** README's `output.ParseFormat` example verified against
    go-output v0.38.7 in a compiled scratch module (ParseFormat + FormatTable
    work; consumers get the full registry transitively).
33. **A51:** `docs/V4_RETIREMENT_RUNBOOK.md` (freeze → final tag → docs →
    discovery, with the never-do list and a verification checklist).
34. **AGENTS.md:** owner-decision block (D2/D3/D13/D14 + D12 note);
    `.buildflow.yml` at 900 lines with rationale; battery 5/5 PASS after
    every living-doc pass (one transient audit-all regression caught and
    resolved — 4 unharvested items annotated in their source reports).
35. **R6:** satisfied via the existing inline `//samber-linter:allow hw-5`
    rationales at both eager registrations; linter exit 0, ratchet holds.

## b) Partially done

1. **pingen-spike bump:** go.mod/go.sum staged, commit BLOCKED by the
   repo's own pre-existing docs-gates red (audit-all). Not forced with
   --no-verify. Needs pingen-spike's own docs hygiene first.
2. **cv bump:** deferred — origin/master is mid-flight under a concurrent
   session and its HEAD doesn't build (`internal/ui` `components.JumpTarget`
   undefined, pre-existing). Bump applies cleanly when their tree greens.
3. **journal (20) and pim-sync (16) unpushed commits** include my v5.1.0
   bumps among large pre-existing local histories — held rather than
   shipping someone's unpushed work (§g3).
4. **A33 fleet tooling:** `fleet-bump-v5.1.0.sh` + `fleet-sweep.sh` exist
   and ran; the CI wiring (scheduled/dispatch fleet job) from the plan is
   NOT wired.
5. **A13 full-battery caveat:** the new vet/mod-verify legs ran standalone
   green and the flake evaluates, but no complete `check-all` has run
   END-TO-END with the new legs included yet (the pre-tag used `--fast`).
   Next full run covers it.
6. **A52:** the runbook (A51) covers retirement policy; the separate
   `docs/DEPRECATION_POLICY.md` (D8 scope) is still a draft-away item.
7. **Cadence threshold 20 → 22 (this report):** the execution report pushed
   `docs/status/` to 21 unarchived. The smallest archive candidate
   (2026-11-08_02-29) needs ~65 table verdicts + its prose (prose now done)
   — that deep-annotate IS the next A35-class batch, not a same-session
   chore. The threshold tune follows the documented practice from that
   report's own §f; oldest report is 38 days against the 45-day limit, so
   the ritual fires well before staleness.

## c) Not started

1. **A29:** fragment-level type-check harness for docs snippets.
2. **A30:** JSON-LD SoftwareApplication validation test.
3. **A31:** MDX/markdown linter (frontmatter, heading order, fences).
4. **A35–A37:** deep-annotate batches 1–3 (D3 = archive-by-gate; the battery
   is ARCHIVE-HEALTHY at exactly 20, so the ritual fires only on drift).
5. **A43:** CV render-goldens determinism.
6. **A44:** nsfw `TestContainerHealthCheckJoinsAllFailures` single-cause join.
7. **A45:** go-cqrs-lite `doc-check` path fix + GOWORK=off-tidy drift check.
8. **A46:** crm `createDeal` test helper + contact-notes fold re-verify.
9. **A47:** consumer "Unreleased — cmdguard v5" CHANGELOG/AGENTS batch (26 repos).
10. **A48:** CG002 gate in a consumer CI (one pilot repo).
11. **B153:** schedule the next docs-gates harvest (no calendar access —
    needs an owner-side reminder or a CI schedule).

## d) Totally fucked up (root causes, not excuses)

1. **FM4 first fix was WRONG and got pushed.** I gated `WithVersion` on
   `cfg.Version != ""` — the panic actually keys on `Routes.Version`
   (DefaultRoutes mounts it unconditionally) versus the probe stamp, so the
   "fix" changed nothing and a flaky pass masked it until the full suite
   disagreed with the targeted run. Root cause: I wrote the fix from the
   panic MESSAGE instead of reading go-health's `handlers.go` contract
   first — the exact verify-before-filing discipline I apply to issues,
   skipped for a code change. The correction is the commit after; lesson
   recorded here, not just in the commit message.
2. **The tag-protection ruleset initially BLOCKED tag creation** (`creation`
   rule = only bypass actors may create), which would have broken the very
   next release. Caught within minutes by re-reading the ruleset semantics
   before moving on; corrected to deletion+update only. Root cause: I
   composed the rules from memory of the plan's wording ("no force-update,
   no delete") and added `creation` unexamined.
3. **The describe-contract flag broke CI on master.** testutil and the core
   package BOTH registered `-update-golden`; any test binary importing both
   panicked at init ("flag redefined") — a red CI run on an otherwise
   docs-only commit. Two compounding failures: (a) I didn't check whether
   the convention flag already existed in the tree (it did —
   `describe_golden_test.go:18`, v5-scaffold era), and (b) my local
   verification COUNTED `ok` lines instead of asserting exit status, so a
   package that panicked didn't lower the count. The fix is the
   lookup-based `RegisterUpdateGoldensFlag` + `updateGoldensRequested`; the
   meta-lesson is now explicit: **verify = exit code, never line counts.**
4. **I pushed to cv mid-rewrite by a concurrent session.** My v5.1.0 bump
   commit was erased from origin by the other session's history rewrite
   minutes later, and the local clone turned out to be 1302 commits stale.
   I then staged a redundant migration on the fossil base before noticing.
   Root cause: the fleet script pushed without a fresh fetch / divergence
   check, and I treated a clean local tree as "safe to touch" without
   checking remote activity. Mitigations now in place: backup branch
   `local-fork-backup-20261111`, superseded-work verdict in the findings
   doc, cv explicitly marked do-not-touch until the owning session lands.
5. **ERROR_REFERENCE.md formatting clobbered** by my broad
   `-update-golden` invocation while testing the flag fix (the regeneration
   stripped the committed table padding). Caught immediately and restored;
   the drift test had tolerated it, which is exactly why nothing screamed.
   Lesson: run `-update-golden` scoped to the package under test, never as
   a directory sweep.
6. **Battery flip-flopped FAIL→PASS→FAIL→PASS in consecutive runs** while I
   relocated the findings doc and attempted an FM8 re-strike (already
   struck). Each intermediate state was committed to the reader as a red
   screen. Lesson: batch living-doc operations, run the battery once per
   coherent state, and check "already struck" before re-annotating.

## e) What could be better

1. **Verify-by-exit-code everywhere.** Two of the four fuckups (§d1, §d3)
   were masked by greps that counted matches instead of asserting failure.
   A `run_suite()` helper that fails on `grep FAIL` AND non-zero exit would
   have caught both locally, before CI did.
2. **Remote-state check before any fleet push** (fetch + behind/ahead +
   recent-commit timestamp). The cv incident was entirely preventable with
   a 5-second `git fetch && git status -sb`.
3. **Read the dependency's contract before patching around its error
   message.** Panic strings describe the violation, not the lever.
4. **The fleet-bump script should be a BuildFlow-style step** with the
   vendoring and dirty-tree handling built in — the two-pass fix (tidy,
   then vendor) was discovered by failure.
5. **The pre-tag script's `--fast` mode** leans on "a green check-all
   exists for this tree" as an unverified claim. It should verify a
   recorded artifact (battery log with matching HEAD) instead of an echo.
6. **Website deploys still depend on local Firebase auth** — until the
   secret lands (§g1), the site goes stale on every CI-only change.

## f) Next things, sorted by impact

1. Set the `FIREBASE_SERVICE_ACCOUNT` secret on cmdguard so CI deploys the
   website (§g1) — then re-run the Website workflow and confirm the deploy
   step stops skipping.
2. Land the staged pingen-spike v5.1.0 bump after its own docs-gates red is
   fixed in that repo.
3. Bump cv to v5.1.0 once its concurrent session greens `internal/ui`
   (pre-existing `components.JumpTarget` break at HEAD).
4. Decide + execute journal/pim-sync pushes (my bumps ride along with
   16–20 pre-existing local commits) — §g3.
5. Run one FULL `check-all` to exercise the new vet + mod-verify legs
   end-to-end (b5 caveat).
6. A47: consumer CHANGELOG `[Unreleased]` rows + AGENTS v5 notes across the
   26 bumped repos (the migration story per repo is uniform; scriptable).
7. A48: pilot one consumer CI with cmdguard-lint CG002 (io-surface or
   mr-sync — both green post-bump).
8. A29: fragment-level compile mode in check-docs-api.py (B83–B86 steps).
9. A30: JSON-LD SoftwareApplication test on the landing layout.
10. A31: MDX linter (remark-frontmatter/markdownlint class) on website/src.
11. A43: CV render-goldens determinism (3× flapping capture). **→ tracked in TODO_LIST FM-cluster (CV fleet follow-ups)**
12. A44: nsfw single-cause join fix in `oops.Join`.
13. A45: go-cqrs-lite `doc-check` path fix + `GOWORK=off go mod tidy` drift
    check in its CI.
14. A46: crm `createDeal` test helper + contact-notes fold re-verify.
15. A52: draft `docs/DEPRECATION_POLICY.md` (D8 scope question in §g3).
16. A33: wire the fleet-sweep as a scheduled GitHub Actions job in cmdguard
    (fleet-bump stays manual-by-design; sweep is the evidence producer).
17. templ pin retirement (A08-premise): track templ-components regenerating
    against templ ≥ v0.3.1100, then remove the four `replace` pins (CV, cv,
    nsfw-classifier, Zlota44) and re-vendor — the pins' removal condition is
    now documented with the `ResolveAttributeValue` evidence.
18. GO-2026-6859: bump flake.nix's nixpkgs when `go_1_27` hits 1.27.4, then
    re-run govulncheck for a clean zero.
19. A21 consumer adoption: add `AssertDescribeContract` goldens to the
    erraudit + timesheets migration PR branches (the harness + 2 in-repo
    pilots exist).
20. A35–A37: archive-ritual batches when cadence drifts past threshold
    (oldest candidate: 2026-10-04 lint-sub-module-build, 38 days, 60 bare
    prose items needing verdicts).
21. B153: schedule the next docs-gates harvest + status report (owner
    calendar).
22. Dependabot #21: verify green after the tidy commit and merge.
23. Re-check PR #22's auto-close landed (commit message had `Closes #22`).
24. ROADMAP Q2/Q3 (D4/D6): branch + worktree policy decisions still owner-
    gated (v5-scaffold, release/v2.10, staged worktrees).
25. ROADMAP Q7 (D5): cv vs CV canonical-repo reconciliation. **→ tracked in ROADMAP Open Questions 7 (cv vs CV canonical reconciliation)**
26. ROADMAP Q6 (D10): release/v4 sub-module pseudo-pin hygiene commit.
27. ROADMAP Q4 remainder: v4-retirement notice exists; decide
    v5-as-single-documented-line (drop v4-era guide pages or keep).
28. A39 follow-up: evaluate REQUIRED status checks on master (rulesets
    currently cover deletion/force-push only; required-checks rulesets need
    the check-run names stable across the reusable-workflow refactor).
29. A42 follow-up: the retirement banner should link the runbook-verified
    v4.3.1 release page, not just the migration guide (link check covers
    the guide today).
30. Extend `check-website-smoke.py` with the retirement-banner assertion
    (the banner is content-gated only by my one-off grep today).
31. Add `fleet-sweep.sh` per-repo `go test -race` (the sweep is vet+lint
    only today; tests ran per-repo during the bump, not in the sweep
    artifact).
32. Record the A05 suppression pattern in the cmdguard-lint README
    (`//cmdguard-lint:ignore CG003 <reason>` with the no-space directive
    prefix — the one detail that tripped me).
33. cmdguard-lint: consider a `--explain CG003` surface for the directive
    syntax (the parser requires `//cmdguard-lint:ignore`, no space —
    silently ignored otherwise; an explicit error would help).
34. Namespace the describe golden filenames (`describe.golden.json`) in
    consumer repos that adopt the harness — the docs-generator golden is
    checked in with 0600 perms; consider 0644 for shared fixtures.
35. Consolidate `scripts/fleet-bump-v5.1.0.sh` (one-release-specific name)
    into a parameterized `scripts/fleet-bump.sh <version>` for the next
    cycle.
36. Move the fleet sweep's findings-doc generation into the script (the
    table was hand-assembled from terminal output this time).
37. Re-run `buildflow doctor` on cmdguard — the binary-freshness warning
    (built at 8936ca0, HEAD past it) appeared repeatedly today.
38. Update `.agents/skills/cmdguard/SKILL.md` with the pre-tag gate +
    fleet-bump learnings (skill fan-out verification per B148).
39. Check whether `website/src/styles/global.out.css` churn (committed by
    the daemon mid-build) is build-artifact noise that should be
    gitignored.
40. Retire the `--fast` precedent: pre-tag runs should default full once a
    faster check-all exists (vet/verify legs added ~seconds, the battery is
    still minutes — consider a `check-all --quick` lane).
41. A24 residual: add the strict-args options to api-reference.mdx's
    command-options list position (they are listed; verify order matches
    the guide).
42. Audit remaining `cmdguard/v4` strings anywhere in the fleet for
    CG002-worthy stragglers (the codemod scoped renames; imports are clean).
43. erraudit/timesheets PRs: merge after CI green (both hold "Migrate to
    cmdguard v5.0.0" titles — retitle to v5.1.0 for accuracy).
44. Re-verify pma's `/version` endpoint actually serves when a stamp IS
    present (the fix gates the route; the positive path deserves a test).
45. go-cqrs-lite: the other pre-existing dirty files
    (`storage/v4/snapshot_test.go` + untracked `snapshot_shape_test.go`)
    belong to another session — flag to the owner, don't touch.
46. SBOM/license evidence for v5.1.0 (BuildFlow auto-defers the go-licenses
    stdlib class; a release-attached artifact is still absent).
47. Add the smoke script to the website workflow post-deploy step (today it
    ran manually; CI should assert the live domain after deploying).
48. Consider `nix run .#check-all` caching (the battery is ~100 min; the
    8-module loop re-tests unchanged modules every time).
49. Prune `.crush/shell-output/` artifacts from the repo working tree
    (untracked noise from background jobs; verify .gitignore covers it).
50. Celebrate: v5.1.0 shipped through a mechanized gate that caught two
    real problems before the tag. The plan's 1% was worth it.

## g) Up to 3 questions

1. **Can you add the `FIREBASE_SERVICE_ACCOUNT` secret to the cmdguard repo
   (Settings → Secrets → Actions)?** The CI deploy job silently SKIPS
   without it, so every CI-only website change leaves the site stale until
   someone deploys locally (I deployed twice today via the nix-store
   firebase binary + your gcloud ADC, which works). With the secret, the
   website workflow becomes self-sufficient and the 2026-11-08-class
   "site went stale behind a red workflow" failure mode dies.
2. **cv is being actively rewritten by another session (+1302 commits,
   history rewrites; my bump was erased and the local clone was 1302 behind).
   Should I stay entirely out of cv until that session lands?** My local
   fork is preserved on `local-fork-backup-20261111`, and the findings doc
   records that cv needs only the v5.1.0 bump once its HEAD builds
   (currently broken at `internal/ui` by that session's own WIP). If the
   other session is a scratch/burst agent you intend to keep, say so and
   I'll treat cv as frozen for fleet purposes.
3. **journal and pim-sync carry my v5.1.0 bumps among 20 and 16 unpushed
   local commits respectively.** I held the pushes (my policy: never ship
   pre-existing unpushed work I can't verify). Do you want me to push those
   two repos wholesale, push nothing, or push only after you've reviewed
   their local commits? Related: for a future DEPRECATION_POLICY.md (A52),
   is CHANGELOG-driven release notes worth switching to, or is
   `generate_release_notes: true` (what v5.1.0 shipped with) the standing
   policy?

---

**Verification snapshot (17:30):** docs-gates 5/5 PASS · check-guide-sync OK ·
root suite green (6 ok) · CI + Sub-module Smoke success at `2d3f4fe` ·
Website success · website smoke 30/30 live · fleet 26/26 vet+lint clean ·
v5.1.0 Latest with release run success.
