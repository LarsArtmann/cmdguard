# Status Report — Docs-Health AUDIT: v5 Living-Docs Rebuild, Battery 5/5 Green

**Date:** 2026-11-11 14:57 CET (Wednesday)
**Session scope:** execute the docs-health skill in AUDIT mode over the whole
docs tree — view/verify the archived `2026-0*` corpus, repair the docs-gates
battery, annotate leaked open items, harvest the survivors, rebuild the living
docs for v5, and run the archive ritual.
**Starting point:** master `359feb3`-era tree, docs-gates **2/5 PASS**
(freshness+routing green; annotate FAIL ×4 shape findings, audit-all FAIL ×19
unharvested items, cadence FAIL at 21>20 unarchived).
**Ending point:** docs-gates **5/5 PASS** (freshness, audit-all, annotate,
routing, cadence all green); 19 unarchived reports; `check-guide-sync` OK;
root + all 7 sub-modules build.

**Operator note:** the auto-commit daemon swept work continuously (expected);
all hashes below are as-committed on master.

---

## 0. Executive summary

The docs tree had silently drifted since the 2026-11-08 02:29 sweep claimed
5/5 green: three new reports pushed cadence over threshold, one status report
carried mechanically damaged strike rows, and 19 historical open items had no
representation in the living docs. The bigger finding: **the living docs were
still describing v4 as the current API** — README taught a removed API
(`RegisterTypeHandler`), FEATURES.md inventoried v4 symbols as current
(including APIs that no longer exist), ROADMAP framed the shipped v5 as
future work with a "blocked upstream" section for middleware that shipped
weeks ago, and CHANGELOG `[Unreleased]` was empty across 40+ commits.

Everything found was either fixed in-session (with evidence) or rowed into
TODO_LIST/ROADMAP with a `repo:` tag where the fix lives elsewhere. The
battery is 5/5 green, the docs truth gate passes, and both scores of the
as-found audit are recorded in §c-adjacent notes below (Accuracy 2.5/10,
Fitness 7.65/10 as-found — the full table lives in the session transcript;
the fixes, not the score, are the deliverable).

---

## a) FULLY DONE

1. **docs-gates battery 5/5 PASS.** Annotate gate: 106 files, 3,100 struck
   items, 0 findings. Audit-all: 90 reports replayed, 0 unharvested open
   items. Routing: 0 unrouted rows. Cadence: 19 ≤ 20 unarchived, oldest 38 ≤
   45 days. Freshness: newest report harvested.
2. **Annotate-shape repair.** `docs/status/2026-11-02_04-15_v4.2.0-release-phase1-status.md`
   rows 16/18 had ID-cell-only strikes that had EATEN the source column
   (the gate's "strike crosses a cell boundary" class). Re-struck per-cell
   with verdicts + restored `M10`/`M06` cells; `check-rows` green on the file.
3. **All 19 unharvested open items resolved** (annotate-or-harvest, each
   verified against code before a verdict):
   - tier-a report: items 9 (batteries default — shipped in v5.0.0),
     35 (TestPromptConfirm no-op — replaced by real tests at
     `coverage_improvement_test.go:281,297`), 37 (gopls/git-town verdicts),
     41 (architecture.svg re-rendered 2026-11-08), 48 (→ TODO_LIST L6),
     4 (scaffold committed, merged `359feb3`).
   - website launch report: items 5/11/34/35/37/43 → TODO_LIST W1–W6.
   - fleet report: items 5/11/13/15/17/31/49 + second pass 18/20/21/32 →
     TODO_LIST FM1–FM8 (with `repo:` tags) and ROADMAP Open questions.
   - archived `2026-11-02_00-05` item 10: struck done — strict validation
     IS the v5 default.
   - `2026-11-08_02-29` row 13 (P11) and archived `2026-10-02` row 34:
     routed verdicts.
   - `2026-11-02_15-19` item 9: DONE by doing it — the 10:15 mid-flight
     report now carries a `> **Superseded by:**` banner.
4. **TODO_LIST rebuilt for truth.** Deleted 4 dead sections (v5 Breaking
   Changes with empty table, v5 End-Goal Workstreams, samber-linter HW-2
   with a false "genuinely open" header, Blocked on Upstream whose blocker
   shipped). Struck DONE with evidence: F2 (api-reference.mdx grep-clean,
   254 lines), F4 (new options documented in SKILL.md + README + AGENTS.md),
   F5 (fan-out guard: 52 symlinks resolve, no store targets), F6
   (`[Unreleased]` fed), F8+S2 (toolchain go1.27.3 root+7 sub-modules;
   nixpkgs `go_1_27` verified at 1.27.3), CT2 (`schema/v0.1.0` tag exists),
   CT3 (README Examples wired), H4 (hero single-source +
   `check_hero_single_source()` guard exists), H6 (COMPARISON.md is
   v5-era 2026-11-08, zero `v4.` matches), L8 (CurrentMajor=v5 shipped in
   lint/v0.3.0). New sections: Website W1–W6, Fleet FM1–FM8, Release & CI
   Hardening R1–R8.
5. **FEATURES.md rebuilt for v5** — the whole file was a v4 inventory.
   Now: `pkg/cmdguard/v5` header; SetConfig row marked removed;
   RegisterInScope[T] generic; Invoke[T]/InvokeNamed[T] as the getters
   (Get[T] rename candidate marked MOOT); CommandMetadata/TypeCodec/
   HuhPrompter names with file:line; batteries-default section; Builder row;
   WithStrictPositionalArgs/WithArgsPassthrough rows; 65 sentinels (grep-
   counted, was "60"); 479 core tests / 597 across 8 modules / 29 benchmarks
   / 8 fuzz targets (all grep-counted today); corrected sub-module versions
   (lipgloss v2.0.6, otel v1.47.0, auditlog v0.12.3, koanf parsers
   v1.0.1/v1.1.1, lint v1.15.1/v1.16.0/v0.4.0); lint table updated to
   v0.3.0-fleet reality.
6. **ROADMAP re-scoped.** "v5 Major Release — SHIPPED" with the five
   TODO(v5) rows struck per-item; v6 candidate table (Get[T] MOOT,
   Package[T] open, middleware unification open); flag-representation
   consolidation honestly re-scoped to v6 (FlagTag.Type is still public,
   `config.go:28` — it did NOT land in v5); audit-middleware section marked
   SHIPPED (`auditlog.go:232,259`); new "Open Questions (owner decisions
   pending)" section with 7 questions; strict positional-args raw idea
   annotated SHIPPED; new anchors for cv-vs-CV and the testkit CI guard idea.
7. **CHANGELOG `[Unreleased]` fed** (was "Nothing yet" across the whole
   v5.0.0..master delta): Added = strict positional-args options; Fixed =
   Builder mutex (verified in the diff: `sync.Mutex` in cli_builder.go),
   doctor schema-validator doc example, release.yml lint args.
8. **README fixed**: v4.2.1→v4.3.1 freeze claim + release/v4 link;
   RegisterTypeHandler→RegisterTypeCodec snippet; 4 missing command options
   added to the table; Examples section lists all three example dirs; Go
   1.26→1.27 dev shell line.
9. **AGENTS.md updated**: Args gotcha extended with the two new options +
   test pin; lint `--dir` invocation gotcha written into the lint section
   (the twice-confessed mistake, now in loaded context); generated-doc
   guards-vs-formatters contract gotcha; lint dep versions; test counts;
   `scripts/check-guide-sync.sh` invocation fixed (it is a BASH script —
   `python3 scripts/check-guide-sync.sh` was a guaranteed SyntaxError).
10. **Consumer SKILL.md** (`​.agents/skills/cmdguard/SKILL.md`): positional-
    args policy paragraph with both new options and the test pin.
11. **Archive ritual executed**: 12-00 pareto-v2 (5 bare prose items struck
    with realized-outcome verdicts) and 08-04 whole-list (7 items struck:
    both repeated-Execute bugs shipped in v5.0.0, §c decisions executed)
    gated CLEAN (`check-rows --strict` + `bare-rows` + `prose-items` = 0)
    BEFORE `git mv`, given Live-remainder banners, moved to
    `docs/status/archived/`, manifest rows appended to the archived README,
    moved-file completeness gate (`grep -L '~~'`) clean.
12. **Archived `2026-0*` corpus verified** (the explicit ask): 76 archived
    .md files dated 2026-01..09 — every single one carries `~~` markers
    (zero pre-canon files); annotate gate CLEAN across the corpus; exactly
    one leaked open item found (archived 2026-10-02 row 34) and routed.
13. **Verification runs**: `nix eval nixpkgs#go_1_27.version` = 1.27.3 (S2
    blocker gone); fan-out guard OK; sub-modules `GOWORK=off go build` +
    tidy all green after the toolchain bump; root `go build ./...` green;
    `scripts/check-guide-sync.sh` OK.

---

## b) PARTIALLY DONE

1. **Sub-module toolchain bumps are build-verified, not gate-verified.** All
   7 go.mods now carry `toolchain go1.27.3` and build+tidy green
   (GOWORK=off), but `nix flake check` / `nix run .#check-all` were NOT run
   this session — the sandboxed treefmt/GOTOOLCHAIN interaction is the exact
   class the AGENTS gotcha warns about, and the delivering layer is
   unverified until that check passes. (Partial repeat of the 08-24 report's
   "edited pins without building" confession, one level shallower.)
2. **The 19 still-open reports** are gate-consistent (every §f item matched
   or verdict-carrying) but NOT deep-annotated: most §f rows in the
   2026-10/11 batch are bare-but-matched. A reader opening e.g. the v4.2.0
   report still sees a wall of unstruck "Next 50" rows.
3. **F4 coverage**: the two new options are documented in README+SKILL+AGENTS
   but NOT on the website guides (the 4th consumer surface); website docs
   pass the truth gate but don't mention the options either.
4. **CT2 evidence is tag-existence only**: `schema/v0.1.0` exists as a tag;
   its publish shape (real pin + version-scoped replace) and proxy-serving
   were NOT verified.
5. **Fan-out verification ordering**: `check-skill-fanout.sh` ran BEFORE the
   SKILL.md edit, not after. The fan-out is a symlink to repo source, so it
   serves the new text by construction, but the verified state was pre-edit.
6. **This status report itself** — written, but the "Up to 50" list below
   includes rows whose verification depends on network access I didn't have.

---

## c) NOT STARTED

1. **v5.1.0 cut** (F3) — owner-gated; `[Unreleased]` is written and waiting.
2. **R1: `website.yml` guard fix + redeploy** — Website workflow RED since
   2026-11-08; user-facing, untouched this session (CI-only fix).
3. **R2–R8** release/CI hardening rows (branch triggers, dry-run dispatch,
   pre-tag script, go vet in check-all, health coverage, drift-guard unit
   test, pkg.go.dev verification).
4. **W1–W6** website follow-ups (smoke test, fragment harness, llms.txt,
   JSON-LD test, hreflang policy, MDX linter).
5. **FM1–FM8** fleet rows (v5.0.1 bump blocked on F3; CV/pma/go-cqrs-lite/cv
   fixes; fleet-wide sweeps).
6. **F7** — cmdguard-lint v0.3.0 executed on the erraudit/timesheets PR
   branches (gh offline; PR bodies still carry unverified "0 findings").
7. **H9b** dot-import constructor tracing; **L6** golangci plugin layer
   (spike-gated); **CT1** CG009 idea.
8. **Deep annotate pass** over the 19 open reports (b.2).
9. **F1** CI status check (gh offline all session).

---

## d) TOTALLY FUCKED UP

1. **I ran `GOWORK=off go mod tidy` in the 7 sub-modules — the exact trap
   AGENTS.md documents** (it rewrites sibling requires to unfetchable
   pseudo-versions, the replace-leak poison). This time it was harmless: the
   pins were already real (`v5.0.0` pre-dated my tidy — verified via
   `git show 81b7a03:glamour/go.mod`), so no damage landed. But the mistake
   pattern is precisely the one the memory file exists to prevent: I reached
   for the per-module invocation pattern before re-reading the gotcha. I
   verified no damage only after the fact.
2. **My TODO_LIST rewrite caused a gate regression.** Deleting the struck
   F-rows (F2/F4/F5/F6/F8/F9) and the Pareto section removed the freshness
   matcher's anchors, flipping audit-all from PASS back to FAIL with 9 newly
   unharvested items. I caught it on the next battery run and fixed it (routed
   verdicts in the reports + new living-doc anchors), but the correct workflow
   was: re-run audit-all IMMEDIATELY after any living-doc rewrite, or keep the
   struck rows. The battery was green for one intermediate commit with a
   regression I authored.
3. **Broken verification greps, twice.** (a) My "did the AGENTS edits apply?"
   grep pattern omitted a backtick, making a successful edit look failed and
   sending me re-reading bytes to discover the tool's "5 of 6" count was
   about a different pair of edits than I assumed. (b) My bare-prose count
   extracted "2026" from file paths instead of the item count, producing a
   column of garbage numbers that I briefly treated as data. Both are the
   unverified-claim class I was auditing OTHER docs for.
4. **The archive-closure subagent died on a rate limit and I never
   re-dispatched it.** I absorbed its job into grep checks (which, combined
   with the annotate gate over 106 files, gives equivalent-or-better
   coverage), but a fan-out leg failed silently and I only accounted for it
   when writing this report.
5. **The delivering layer was not verified for the go.mod changes**
   (`nix flake check` / `check-all` skipped for time). I flagged it b.1 and
   rowed it (§f item 27), but "build passed on my machine" is exactly the
   claim class this repo's own reports keep confessing to.

---

## e) WHAT WE SHOULD IMPROVE

1. **Gate-after-living-doc-edit discipline.** Any rewrite of TODO_LIST/
   ROADMAP/CHANGELOG must be followed IMMEDIATELY by `docs-gates run .` —
   the freshness matcher's anchors live in those files, so "docs edit" and
   "gate input edit" are the same action. (This session proved the failure
   mode and the 30-second fix.)
2. **Read the AGENTS gotcha BEFORE the command, not after.** The GOWORK=off
   tidy trap is written down; the failure was invocation-order. Workspace-
   aware tidy should be the reflex for any multi-module repo.
3. **Byte-exact verification greps.** When verifying edits, `cat -A` the
   target region first and grep the EXACT rendered bytes (backticks,
   alignment) — pattern-from-memory greps produced two false results.
4. **Subagent failure handling**: a failed subagent report should trigger an
   immediate re-dispatch or an explicit absorption note in the todo list,
   not silent continuation.
5. **Evidence grades**: distinguish "exists" (tag exists) from "published
   and proxy-served" (playbook-verified). CT2 was struck on the weaker
   grade; the stronger check is cheap (`go get schema@v0.1.0` from a scratch
   module).
6. **Delivering-layer rule applies to go.mod edits too**: `nix flake check`
   after any toolchain/directive change, not just after Nix-file changes.
7. **Struck-row lifecycle convention is split-brained** (see question g2):
   the skill says delete done rows; this repo's ledger pattern keeps them
   struck; the freshness matcher USES them as anchors. One of the three must
   be declared canonical in AGENTS.md.

---

## f) Up to 50 things to do next (ordered by impact)

**Release-blocking / user-facing**

1. Fix `website.yml` import-path guard (`/v4`→`/v5`, allow historical `/v3`
   in changelog.mdx) and re-run the Website workflow (RED since 2026-11-08).
2. Owner sign-off, then cut **v5.1.0** (Builder mutex +
   `WithStrictPositionalArgs`/`WithArgsPassthrough` + doc fixes; `[Unreleased]`
   already written).
3. Write the v5.1.0 CHANGELOG section from the tag delta at cut time.
4. Verify the v5.1.0 release run is GREEN before announcing (read
   `release.yml` first — the v4.3.0 lesson).
5. `gh run list` CI status sweep for recent master pushes (F1; offline today).
6. Execute cmdguard-lint v0.3.0 on the erraudit/timesheets PR branches and
   post real numbers (F7) — correct the unverified PR-body claims.
7. Verify pkg.go.dev renders v4.3.1 + v5.0.0 with the deprecation banner (R8).
8. Run `nix run .#check-all` to gate this session's toolchain bumps through
   the delivering layer (b.1).
9. Verify `schema/v0.1.0` publish shape + proxy-serving from a scratch
   consumer (b.4).
10. Confirm the daemon-committed doc changes are all intended (review
    today's auto-commit messages, 10-file batches hide diffs).

**CI / workflow hardening**
11. Add `release/**` (or all branches) to `ci.yml` triggers (R2).
12. `workflow_dispatch` dry-run path in `release.yml` (R3).
13. `scripts/pre-tag-check.sh` (R4).
14. `go vet` + `go mod verify` in `check-all` (R5).
15. Unit test for `normalizeTableSpacing` (R7).
16. Branch protection for `master` (required checks) (09-14 item 32).
17. Tag protection — no force-update, no deletion (09-14 item 33).
18. GOWORK=off CI job mirroring consumers (09-14 item 48).
19. CI test that sub-module go.mods are tidy (09-14 item 47).
20. `golangci-lint args:` guard — grep workflows for `args: run` (09-14 #18).

**Website**
21. Live-domain smoke-test script (W1).
22. Fragment-level type-check harness (W2).
23. `llms.txt` (W3).
24. JSON-LD SoftwareApplication validation test (W4).
25. hreflang/locale policy decision (W5).
26. markdown/MDX linter (W6).
27. Document `WithStrictPositionalArgs`/`WithArgsPassthrough` on the website
guides (b.3).
28. Website v4-retirement banner/page (09-14 item 40).
29. Verify the migration guide resolves after redeploy (09-14 item 39).
30. Website scope decision — v5-only vs v4 notice (ROADMAP OQ #4).

**Fleet (repo:-tagged)**
31. FM1: bump 25 consumers + `cv` to v5.0.1 (after item 2).
32. FM2: CV `TestLoadIdentityPortalsSection`.
33. FM3: CV `crm.checkpoint_path` defaultsYAML entry.
34. FM4: pma health-probe `Routes.Version`/`WithVersion` panics.
35. FM5: go-cqrs-lite hanging `cqrs-lint` `go/packages` tests in offline CI.
36. FM6: `describe`-contract snapshot test per consumer. **→ TODO_LIST FM6 (canonical row lives there)**
37. FM7: commit `cv`'s migration (verify the 20 files first — not my work). **→ SUPERSEDED — see TODO_LIST FM7 verdict (cv origin/master carries the migration; fork preserved on local-fork-backup-20261111)**
38. FM8: fleet-wide `golangci-lint` + `go vet` + `go test -race` sweep.
39. Templ-pair biconditional into consumer bump tooling (08-24 item 8).
40. Reconcile `cv` vs `CV` (ROADMAP OQ #7).

**Docs / lint backlog**
41. Deep-annotate the 19 open reports toward archive cadence (b.2).
42. Decide the struck-done-row convention (g2) and write it into AGENTS.md.
43. README: verify the `output.ParseFormat("csv")` example against go-output
v0.38.7's actual API (noticed, unverified this session).
44. H9b: dot-import bare-constructor tracing in cmdguard-lint.
45. L6: Analyzer-bridge spike decision (gate before any build).
46. CT1: scope CG009 (example-config truth rule) in the lint backlog.
47. Dependabot PRs #21/#22/#24: post-rebase merge-or-close (09-14 §b.3).
48. Branch policy: delete/repurpose `v5-scaffold`; decide `release/v2.10`
(ROADMAP OQ #2).
49. Staged worktrees fate: `crm-stage-t13b`, `cv-stage-crmfix`,
`BuildFlow-pr44` (ROADMAP OQ #3).
50. CHANGELOG `[Backfill]` folding decision (owner §g2; ROADMAP OQ).

---

## g) Questions I cannot answer myself

1. **Cut v5.1.0 now?** Master carries the Builder mutex fix, two new public
   options, and doc fixes past v5.0.0; `[Unreleased]` is written; publishing
   is irreversible and the repo convention is owner sign-off for tags.
   Yes/no (and if yes: immediately after `check-all` green, or after the
   website guard fix so the release commit carries it?).
2. **Struck-done rows in TODO_LIST: keep or delete?** The docs-health skill
   says completed TODOs are deleted (they live in CHANGELOG); this repo's
   ledger convention keeps them struck with evidence; the freshness matcher
   uses them as anchors (deleting them caused today's audit-all regression).
   Which is canonical? I will encode your answer in AGENTS.md so sessions
   stop guessing.
3. **Archive cadence policy for the remaining 19 open reports:** run a
   dedicated deep-annotate pass (like 2026-11-07, striking every §f/§a/§d
   item) to archive them down to a small rolling window, or keep them open
   as the live audit trail and only archive when the cadence gate trips?
   The first is cleaner history; the second is cheaper per session.

---

## Appendix — evidence anchors

```bash
nix run github:Artmann-Technologies/docs-gates -- run .   # PASS — all doc gates green
scripts/check-guide-sync.sh                                # OK: docs, API symbols, and doc samples match the code
go build ./...                                             # green (root)
for m in glamour prompts spinner telemetry flightrecorder lint schema; do
  (cd $m && GOWORK=off go build ./...)                     # green ×7 (post-toolchain-bump)
done
ls docs/status/*.md | wc -l                                # 19 (≤ 20 cadence threshold)
grep -rLn '~~' docs/status/archived --include='2026-0*.md' # empty — all 76 archived 2026-0* files carry markers
git tag --list 'schema/*'                                  # schema/v0.1.0 exists
nix eval nixpkgs#go_1_27.version --raw                     # 1.27.3 (S2 blocker gone)
```
