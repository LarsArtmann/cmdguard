# TODO List

> Short- and mid-term improvement tasks — actionable, bounded, with status.
> Derived from FEATURES.md partially-functional items, status report harvests,
> and code analysis. Updated as work completes. Done rows are struck with
> evidence; the full history lives in CHANGELOG.md.

---

## Delivery Follow-ups (harvested 2026-11-08 12:25 + 2026-11-11 09:14 reports)

Top items from the Pareto v2 execution's next-50 list plus the v4-retirement
session. The full 50-item lists live in
`docs/status/2026-11-08_12-25_self-critique-status-v2-execution.md` §f and
`docs/status/2026-11-11_09-14_v4-retirement-and-v5-master-merge-status.md` §9.

<!-- check-rows:ledger -->

| # | Task | Evidence / notes | Status |
| - | ---- | ---------------- | ------ |

---

## Website Follow-ups (harvested 2026-11-11 01:33 launch report)

The v5 docs site launched; these are its own report's open improvements.
Source: `docs/status/2026-11-11_01-33_website-v5-docs-launch-status.md` §f.

<!-- check-rows:ledger -->

| #  | Task                                                               | Evidence / notes                                                                                               | Priority | Status  |
| -- | ------------------------------------------------------------------ | -------------------------------------------------------------------------------------------------------------- | -------- | ------- |
| W2 | Fragment-level type-check harness for docs snippets                | complete Go programs already compile via `scripts/check-docs-api.py`; partial fragments are still not compiled | Medium   | 🔴 TODO |
| W4 | Validate JSON-LD (SoftwareApplication) with a structured-data test | markup exists at `website/src/layouts/LandingLayout.astro:41`; nothing asserts its shape                       | Low      | 🔴 TODO |
| W6 | Add a markdown/MDX linter                                          | frontmatter required, heading order, code-fence languages; no remark/markdownlint config today                 | Low      | 🔴 TODO |

---

## Fleet v5 Migration Follow-ups (harvested 2026-11-11 08:24 report)

Consumer-repo fixes discovered during the 25-repo v5 migration. Rows whose fix
lands outside this repo carry a `repo:` tag — tick them merge-by-reference from
that repo's history, not from greps here.
Source: `docs/status/2026-11-11_08-24_cmdguard-v5-fleet-consumer-migration-status.md` §f.

<!-- check-rows:ledger -->

| #   | Task                                               | Evidence / notes                                                                                                                                                               | Status  |
| --- | -------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------- |
| FM1 | Bump the 25 consumers + `cv` to the v5.1.0 release | UNBLOCKED 2026-11-11: v5.1.0 is Latest (F3 done); bump target is v5.1.0 (not v5.0.1 — it shipped first); the templ indirect floor story is documented in MIGRATION_v4_v5.md §3 | 🔴 TODO |

---

## Release & CI Hardening (harvested 2026-11-11 09:14 report)

From the v4-retirement session's §7/§9. Source:
`docs/status/2026-11-11_09-14_v4-retirement-and-v5-master-merge-status.md`.

<!-- check-rows:ledger -->

| # | Task | Evidence / notes | Priority | Status |
| - | ---- | ---------------- | -------- | ------ |

---

## Config-Truth Follow-ups (2026-11-02)

From the config-truth execution (plan:
[docs/planning/2026-11-02_05-36_config-truth-pareto-execution-plan.md](docs/planning/2026-11-02_05-36_config-truth-pareto-execution-plan.md)).

<!-- check-rows:ledger -->

| #   | Task                                                    | Notes                                                                                                                                                                                                                                                                                                                         | Priority | Status  |
| --- | ------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------- |
| CT1 | cmdguard-lint CG009 idea: example configs stay truthful | A rule that parses a repo's committed example configs (e.g. `config.example.yaml`) and rejects keys the CLI's config vocabulary does not accept. Needs cross-file type info (or a mode that shells `describe`/`ConfigKeys` from a built binary), so it is an IDEA, not a task — scope it in the lint backlog before building. | Low      | 💡 IDEA |

---

## Harvested (2026-11-07 docs-health pass)

Items pulled forward from the 2026-06→2026-09 status reports during the
annotate/archive sweep. Everything else those reports contained was either
resolved (struck with evidence in the report itself) or routed to ROADMAP
"Harvested from 2026-0X status reports".

<!-- check-rows:ledger -->

| #   | Task                                                  | Notes                                                                                                                                                                                                                               | Priority | Status  |
| --- | ----------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------- |
| H9b | Trace dot-imported bare constructors in cmdguard-lint | `isSelectorCall` requires a selector, so `NewCLI(...)` under a dot import is invisible to CG003/CG004/CG006; pinned by `lint/crossfile_imports_test.go` (flips to expect findings when fixed). Non-idiomatic pattern; Low priority. | Low      | 🔴 TODO |

---

## Lint Sub-Module (`lint/`)

Harvested from `docs/status/2026-10-04_07-11_cmdguard-lint-sub-module-build.md` §f
(full 50-item backlog lives there; these are the near-term actionable ones).

<!-- check-rows:ledger -->

| #  | Task                                           | Notes                                                                                                                                                                                                                                                                                                                                                                                                                | Priority | Status  |
| -- | ---------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------- |
| L6 | golangci-lint module-plugin distribution layer | Scoped 2026-10-09: requires an `analysis.Analyzer` bridge (golangci/plugin-module-register, per go-humanize-linter pattern: `lint/plugin/` package + `.custom-gcl.yml` + `.golangci.custom.yml`). One pass = one package's files, so cross-file tracing carries over; checks need a pass.Files→project adapter. NOT a thin wiring task — do not half-ship. Also needs the custom-gcl build verified against corpora. | Low      | 🔴 TODO |

---

## Security

<!-- check-rows:ledger -->

| # | Task | Notes | Priority | Status |
| - | ---- | ----- | -------- | ------ |

---

## Blocked on Upstream

Currently nothing. (The 2026-10 audit-middleware block resolved when
samber-do-auditlog v0.11.0+ shipped `Plugin.RecordCommand`; `AuditMiddleware[T]`
has been in core since — see FEATURES.md Audit Log table.)

### FM addenda (2026-11-11 audit-all pass)

<!-- check-rows:ledger -->

| # | Task | Evidence / notes | Status |
| - | ---- | ---------------- | ------ |

## Session follow-ups (2026-11-11 v5.1.0 execution pass)

| #    | Task                                                                                                                                                                         | Evidence / notes                                                                                                                                       | Status  |
| ---- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------ | ------- |
| SU3  | `.agents/skills/cmdguard/SKILL.md`: document the pre-tag gate + describe-harness adoption; re-run check-skill-fanout after                                                   | consumer skill fan-out verification per the docs-health procedure (B148)                                                                               | 🔴 TODO |
| SU9  | pingen-spike: fix its docs-gates red, then commit the staged v5.1.0 bump                                                                                                     | go.mod/go.sum staged; the repo's own hook refuses unrelated commits (no --no-verify)                                                                   | 🔴 TODO |
| SU11 | CV: pin render-goldens determinism (fixture order; stable ×3 runs)                                                                                                           | flapping order suspected; A43 in the 15-10 plan                                                                                                        | 🔴 TODO |
| SU12 | nsfw-classifier: `TestContainerHealthCheckJoinsAllFailures` single-cause `oops.Join` fix                                                                                     | single-cause join wraps as multi-cause; A44 in the 15-10 plan                                                                                          | 🔴 TODO |
| SU13 | crm: restore the missing `createDeal` test helper + re-verify contact-notes fold after WIP lands                                                                             | httpapi tests affected; A46 in the 15-10 plan                                                                                                          | 🔴 TODO |
| SU18 | Harden `scripts/check-website-smoke.py`: per-route timeout, bounded retries, clean FAIL summary                                                                              | first run today crashed with a raw TimeoutError traceback mid-run; rerun green instead of fix (20-43 report §b4)                                       | 🔴 TODO |
| SU19 | Post-deploy freshness assertion: served release must match the just-created Hosting release (id or hash), not just 200s                                                      | the ~60s stale/503 window would pass a 200-only check on old content (20-43 report §b3)                                                                | 🔴 TODO |
| SU20 | Commit the GitHub-Actions secret-quoting lesson to crush-config `references/lessons.md` (repo: crush-config)                                                                 | `echo "${{ secrets.X }}"` eats multiline JSON quotes; env + `printf '%s'` is the fix; proven locally (20-43 report §d1)                                | 🔴 TODO |
| SU40 | templ-components: remediate or baseline the 55 pre-commit-blocking findings (file-size 50 + docs-gates 5) so the repo is committable (repo: templ-components)                | its pairing landed via --no-verify; the gate blocks every commit there (2026-11-20_20-15 §f3)                                                          | 🔴 TODO |
| SU44 | Promote the describe-fixture cross-reference to a bounded task: pin machine-contract.mdx excerpt bytes to describe_golden_test.go output                                     | the excerpt is hand-maintained and was fabricated once before the live-binary check caught it (2026-11-20_20-15 §f7)                                   | 🔴 TODO |
| SU46 | Promote the pre-push workflow smoke to a bounded task: `go run <tool>@<pin> --help` per pinned action dependency plus a job-level permissions checklist for any pushing step | today produced three authored CI failures (escaping, token permissions, retry); run 44845752228 log showed Contents: read only (2026-11-21_01-56 §f14) | 🔴 TODO |

## Sentinel & erraudit follow-ups (harvested 2026-11-22 17-40 report)

Near-term rows from the erraudit tree audit that removed 5 dead sentinels and
added the `TestNoDeadSentinels` tripwire. The long-tail brainstorm stays in the
report's §f (`docs/status/2026-11-22_17-40_erraudit-tree-audit-and-dead-sentinel-removal-status.md`).

| #    | Task                                                                                                                          | Evidence / notes                                                                                                                      | Status  |
| ---- | ----------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------- | ------- |
~~| SA1  | Extend `TestNoDeadSentinels` to the `flightrecorder` and `lint` modules                                                       | core-only today; those modules carry `ErrAlreadyStarted`/`ErrNotEnabled`/`ErrBaseline*` (2026-11-22 §b1)                              | 🔴 TODO |~~ done at `c71f028` (testutil.FindDeadSentinels + tripwires in flightrecorder/ and lint/, both green)
~~| SA2  | Add a temp-dir self-test fixture inside the tripwire test so its detection logic is regression-tested without touching source | source-append probes dirtied the tree once today (2026-11-22 §d1)                                                                     | 🔴 TODO |~~ done at `7a6d444` (pkg/testutil/dead_sentinels_test.go fixtures; caught the self-decl-wiring, root-dir-skip and lost-test-file-skip bugs before landing)
| SA3  | Pin the erraudit acceptance posture in a machine artifact (config file or CI gate), not just AGENTS.md prose                  | stdlib_constructor + generic_return accepted by prose only; a future session can re-litigate (2026-11-22 §b2; owner question pending) | 🔴 TODO |
~~| SA4  | Cross-file tripwire: SKILL.md's derived sentinel count must equal the generated ERROR_REFERENCE header                        | the stored derivation `rg` pattern had already drifted from the formatter-normalized table shape once (2026-11-22 §e5)                | 🔴 TODO |~~ done at `1e5e706` (TestSkillSentinelCountMatchesReference in errors_reference_test.go: every SKILL.md count must equal the live generator's; drift direction probe-verified)
~~| SA5  | Document in AGENTS.md why `tree --max` shows stdlib_constructor badges while `lint --max` shows none                          | runner drops enforcement analyzers in lint mode; verified 0 findings there (2026-11-22 §b3)                                           | 🔴 TODO |~~ done at `c1e383a` (AGENTS.md erraudit bullet carries the why)
| SA6  | Run `erraudit fix` dry-run + `lint --type-aware` sweep across all 8 modules; disposition every finding                        | not attempted this session (2026-11-22 §c1/§c2)                                                                                       | 🔴 TODO |
~~| SA7  | Run `scripts/check-skill-fanout.sh` for today's SKILL.md edit and fold the check into the post-edit ritual                    | in-repo source edited without verifying the delivering layer (2026-11-22 §b4)                                                         | 🔴 TODO |~~ done at `c1e383a` (guard green: 54 symlinks resolve; ritual folded into AGENTS.md Living-doc edit ritual bullet)
| SA8  | Owner decision recorded here once made: release vehicle for the 5 removed sentinels (v5.2.0 minor vs v6)                      | exported identifiers removed; compile-time break only for direct references (2026-11-22 §f1)                                          | 🔴 TODO |
| SA9  | Cross-file: document (or trace) `var cli CLI = <constructor>` type-decl form as unsupported (cmdguard-lint docs)              | open since 2026-10-10 backlog §15 (audit-all replay 2026-11-22)                                                                       | 🔴 TODO |
~~| SA10 | Badge self-describing tiering: README/CI alt-text + title "coverage (all modules)"                                            | 5-min fix; open since 2026-11-21_07-31 §f9 (audit-all replay 2026-11-22)                                                              | 🔴 TODO |~~ done at `81bcc5a` (CI coverage creates the 8-module workspace, fleet-total verified 86.1% locally; badge alt-text/title + README alt-text say coverage (all modules))
| SA11 | Fix ERROR_REFERENCE.md's stale v4 intro prose: title "cmdguard v4", facts path `pkg/cmdguard/v4`, examples use `v4.ErrRequiredFlag`, while the generated table (and repo) are v5 | noticed 2026-11-22 during SA4 tripwire work; the generated section between the BEGIN/END markers is current — only the ~95-line hand-written head drifted | 🔴 TODO |
| SA12 | Verify the new fleet-coverage CI step on the next master push: `go work init` resolves on a real runner, coverage covers all 8 modules, badge regenerates with the "coverage (all modules)" alt-text | CI leg of SA10 is simulation-verified only (2026-11-22_19-23 §b1); first-use workspace resolution + sub-module dep downloads + 300s timeout with the flightrecorder singleton are unproven on CI | 🔴 TODO |
| SA13 | Chase the docs-cadence counter discrepancy: it printed "0 struck DONE row(s) in TODO_LIST.md" immediately after 6 rows were struck with verdicts (gate PASS) | observed 2026-11-22_19-23 §b5; counter vs reality disagreement unexplained — what does it actually count? | 🔴 TODO |
