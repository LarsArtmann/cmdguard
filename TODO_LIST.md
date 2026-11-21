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

| #    | Task                                                                                                                                                                                                       | Evidence / notes                                                                                                                                              | Status                                                                                                                                                    |
| ---- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
~~| SU1 | Re-feed CHANGELOG `[Unreleased]` with post-5.1.0 work | pre-tag script, describe harness, retirement banner, llms.txt, fleet scripts, lint dot-import fix, website guard — all shipped after the tag, none listed | 🔴 TODO |~~ done (2026-11-21: [Unreleased] re-fed — release gate, website truth guards incl. describe excerpt pinning, fleet tooling, retirement banner, llms.txt; commit fbacf76)
~~| SU2 | FEATURES.md: add rows for the new tooling (describe harness, pre-tag gate, website smoke, fleet scripts, rulesets) | the 14-57 rebuild predates today's additions | 🔴 TODO |~~ done (2026-11-21: Repository Tooling section added with file evidence; commit fbacf76)
| SU3  | `.agents/skills/cmdguard/SKILL.md`: document the pre-tag gate + describe-harness adoption; re-run check-skill-fanout after                                                                                 | consumer skill fan-out verification per the docs-health procedure (B148)                                                                                      | 🔴 TODO                                                                                                                                                   |
| SU4  | Release the lint module (next tag) carrying the dot-import constructor tracing fix                                                                                                                         | behavior change is on master + the locally built binary; no published tag carries it yet; BuildFlow provider pin follows the tag                              | 🔴 TODO                                                                                                                                                   |
~~| SU5 | AGENTS.md Quick Start: add `scripts/pre-tag-check.sh` and `scripts/check-website-smoke.py` entries | both are now standing gates; AGENTS.md still documents only the older scripts | 🔴 TODO |~~ done (2026-11-21: pre-tag-check.sh + check-website-smoke.py added to Quick Start; commit f5602ee)
~~| SU6 | One FULL `check-all` run exercising the new vet + go-mod-verify legs end-to-end | legs verified standalone + flake evaluates; no complete run since the flake edit (pre-tag used `--fast`) | 🔴 TODO |~~ done (2026-11-21: full check-all EXIT=0 after treeWalkingRule exhaustruct fix + dep-budget line; SU39)
~~| SU7 | Unit test for `DescribeContractFromBinary` (process-capture helper) | shipped untested — the two pilots exercise the library path only | 🔴 TODO |~~ done (2026-11-21: TestDescribeContractFromBinary + testdata/describeapp fixture; env-prefix and greet assertions; commit bee6a57)
~~| SU8 | Gitignore audit: `.crush/shell-output/` + `website/src/styles/global.out.css` build artifacts | global.out.css churn got committed by the daemon mid-build today; verify both are ignored or deliberately tracked | 🔴 TODO |~~ done (2026-11-21: .crush/shell-output/ ignored, global.out.css untracked (1566 lines of generated Tailwind); commit f5602ee)
| SU9  | pingen-spike: fix its docs-gates red, then commit the staged v5.1.0 bump                                                                                                                                   | go.mod/go.sum staged; the repo's own hook refuses unrelated commits (no --no-verify)                                                                          | 🔴 TODO                                                                                                                                                   |
~~| SU10 | Merge Dependabot #21 (submodule tidy landed, checks green) after owner glance | all failing checks resolved; DontMergeMeYet label allows merge | 🔴 TODO |~~ done (2026-11-21: PR #21 state CLOSED on GitHub)
| SU11 | CV: pin render-goldens determinism (fixture order; stable ×3 runs)                                                                                                                                         | flapping order suspected; A43 in the 15-10 plan                                                                                                               | 🔴 TODO                                                                                                                                                   |
| SU12 | nsfw-classifier: `TestContainerHealthCheckJoinsAllFailures` single-cause `oops.Join` fix                                                                                                                   | single-cause join wraps as multi-cause; A44 in the 15-10 plan                                                                                                 | 🔴 TODO                                                                                                                                                   |
| SU13 | crm: restore the missing `createDeal` test helper + re-verify contact-notes fold after WIP lands                                                                                                           | httpapi tests affected; A46 in the 15-10 plan                                                                                                                 | 🔴 TODO                                                                                                                                                   |
~~| SU14 | Parameterize `scripts/fleet-bump-v5.1.0.sh` → `scripts/fleet-bump.sh <version>` for future cycles | one-release name today; add vendor + dirty-tree handling while at it (the two-pass discovery this cycle) | 🔴 TODO |~~ done (2026-11-21: git mv to scripts/fleet-bump.sh <version>, arg validation; commit f5602ee)
~~| SU15 | Pin firebase-tools in `website.yml` deploy step (exact version or nix binary) | CI ran 15.35.0 vs locally proven 15.28.1; unpinned `pnpm add -g` = repro + supply-chain gap (20-43 report §b1) | 🔴 TODO |~~ done (2026-11-21: firebase-tools@15.37.0 in website.yml deploy; commit f5602ee)
| SU16 | Migrate website deploy to Workload Identity Federation (GitHub OIDC), retire the long-lived SA JSON key + secret                                                                                           | deploy works today via `FIREBASE_SERVICE_ACCOUNT` secret; WIF kills key + rotation burden (20-43 report §b2)                                                  | 🔴 TODO                                                                                                                                                   |
| SU17 | Least-privilege the deploy SA: per-site deployer on `cmdguard` only                                                                                                                                        | SA is project-wide `firebasehosting.admin` over ~20 hosting sites on lars-software (20-43 report §b2)                                                         | 🔴 TODO                                                                                                                                                   |
| SU18 | Harden `scripts/check-website-smoke.py`: per-route timeout, bounded retries, clean FAIL summary                                                                                                            | first run today crashed with a raw TimeoutError traceback mid-run; rerun green instead of fix (20-43 report §b4)                                              | 🔴 TODO                                                                                                                                                   |
| SU19 | Post-deploy freshness assertion: served release must match the just-created Hosting release (id or hash), not just 200s                                                                                    | the ~60s stale/503 window would pass a 200-only check on old content (20-43 report §b3)                                                                       | 🔴 TODO                                                                                                                                                   |
| SU20 | Commit the GitHub-Actions secret-quoting lesson to crush-config `references/lessons.md` (repo: crush-config)                                                                                               | `echo "${{ secrets.X }}"` eats multiline JSON quotes; env + `printf '%s'` is the fix; proven locally (20-43 report §d1)                                       | 🔴 TODO                                                                                                                                                   |
| SU21 | Root-cause the ~60s post-deploy 503/stale window; eliminate or document it                                                                                                                                 | content-identical deploy, 60s brownout, "edge propagation" assumed never verified (20-43 report §d6)                                                          | 🔴 TODO                                                                                                                                                   |
| SU40 | templ-components: remediate or baseline the 55 pre-commit-blocking findings (file-size 50 + docs-gates 5) so the repo is committable (repo: templ-components)                                              | its pairing landed via --no-verify; the gate blocks every commit there (2026-11-20_20-15 §f3)                                                                 | 🔴 TODO                                                                                                                                                   |
| SU42 | annotate.go: per-row TSV validation (kind enum + raw-value rule) with row-numbered errors, extending adf2e80 (repo: docs-gates)                                                                            | a spec with kind `done` failed full-file and late (2026-11-20_20-15 §f5)                                                                                      | 🔴 TODO                                                                                                                                                   |
| SU44 | Promote the describe-fixture cross-reference to a bounded task: pin machine-contract.mdx excerpt bytes to describe_golden_test.go output                                                                   | the excerpt is hand-maintained and was fabricated once before the live-binary check caught it (2026-11-20_20-15 §f7)                                          | 🔴 TODO                                                                                                                                                   |
| SU46 | Promote the pre-push workflow smoke to a bounded task: `go run <tool>@<pin> --help` per pinned action dependency plus a job-level permissions checklist for any pushing step                               | today produced three authored CI failures (escaping, token permissions, retry); run 44845752228 log showed Contents: read only (2026-11-21_01-56 §f14)        | 🔴 TODO                                                                                                                                                   |
| SU48 | taskctl: store plain `Long`/`Example` text and style at help-render time — describe currently embeds ANSI-escaped styled help in the machine contract (found while pinning the SU44 excerpt to live bytes) | the contract `long` field carries raw `\u001b[38;5;…` escape sequences plus lipgloss padding; hostile for contract-driving agents (2026-11-21 excerpt repair) | 🔴 TODO                                                                                                                                                   |
