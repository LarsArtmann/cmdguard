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

| #  | Task                                                                                      | Evidence / notes                                                                                                                                  | Status  |
| -- | ----------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| F1 | Check `gh run` CI status for cmdguard v5 pushes; fix any red legs                         | all verification so far was local; `gh` could not reach api.github.com from the 2026-11-11 docs session — re-check when online                    | 🔴 TODO |
| F3 | Cut v5.1.0 carrying WithStrictPositionalArgs/WithArgsPassthrough + Builder mutex fix (untagged) | 42+ commits past the v5.0.0 tag; `[Unreleased]` is already fed; publish shapes per the lint/v0.2.1 playbook. Owner sign-off required first (see ROADMAP Open questions) | 🔴 TODO |
| F7 | Run cmdguard-lint v0.3.0 on the erraudit/timesheets PR branches; post real numbers        | PR bodies claim "0 findings" without an executed run on the branches                                                                              | 🔴 TODO |

---

## Website Follow-ups (harvested 2026-11-11 01:33 launch report)

The v5 docs site launched; these are its own report's open improvements.
Source: `docs/status/2026-11-11_01-33_website-v5-docs-launch-status.md` §f.

<!-- check-rows:ledger -->

| #  | Task                                                                        | Evidence / notes                                                                                                                       | Priority | Status  |
| -- | --------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------- |
| W1 | Live-domain smoke-test script after deploy                                  | assert 200 for `/`, all sidebar routes, `/sitemap-index.xml`, `/robots.txt`, `/404.html`; today only static link checks exist (`scripts/check-website-links.py`) | Medium | 🔴 TODO |
| W2 | Fragment-level type-check harness for docs snippets                          | complete Go programs already compile via `scripts/check-docs-api.py`; partial fragments are still not compiled                          | Medium   | 🔴 TODO |
| W3 | Add `llms.txt` (links + one-paragraph product summary)                       | `website/public/llms.txt` does not exist                                                                                               | Low      | 🔴 TODO |
| W4 | Validate JSON-LD (SoftwareApplication) with a structured-data test           | markup exists at `website/src/layouts/LandingLayout.astro:41`; nothing asserts its shape                                               | Low      | 🔴 TODO |
| W5 | Decide `hreflang`/locale policy                                              | add only if a second language is planned; otherwise document English-only                                                              | Low      | 🔴 TODO |
| W6 | Add a markdown/MDX linter                                                    | frontmatter required, heading order, code-fence languages; no remark/markdownlint config today                                         | Low      | 🔴 TODO |

---

## Fleet v5 Migration Follow-ups (harvested 2026-11-11 08:24 report)

Consumer-repo fixes discovered during the 25-repo v5 migration. Rows whose fix
lands outside this repo carry a `repo:` tag — tick them merge-by-reference from
that repo's history, not from greps here.
Source: `docs/status/2026-11-11_08-24_cmdguard-v5-fleet-consumer-migration-status.md` §f.

<!-- check-rows:ledger -->

| #   | Task                                                                     | Evidence / notes                                                                                                                                                  | Status  |
| --- | ------------------------------------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| FM1 | Bump the 25 consumers + `cv` to v5.0.1                                    | blocked on the v5.0.1 release (no tag yet; the v5.0.0 go.mod carries the leaked templ indirect floor) — see F3                                                    | 🔴 TODO |
| FM2 | Fix `TestLoadIdentityPortalsSection` (repo: CV)                           | fails at CV HEAD (A/B-verified pre-existing)                                                                                                                       | 🔴 TODO |
| FM3 | Fix `crm.checkpoint_path` missing defaultsYAML entry (repo: CV)           | `TestDefaultsYAMLContainsEveryConfigKey` fails at CV HEAD                                                                                                          | 🔴 TODO |
| FM4 | Fix health-probe panics — `Routes.Version` without `WithVersion` (repo: projects-management-automation) | fails at pre-migration commit (A/B-verified)                                                                                        | 🔴 TODO |
| FM5 | Fix/skip the two hanging `cqrs-lint` `go/packages` tests in offline CI (repo: go-cqrs-lite) | hang at pre-migration commit; network/toolchain-bound                                                                            | 🔴 TODO |
| FM6 | Add a `describe`-contract snapshot test per consumer                      | agents depend on the contract; one golden per consumer CLI catches accidental flag/tree changes                                                                    | 🔴 TODO |

---

## Release & CI Hardening (harvested 2026-11-11 09:14 report)

From the v4-retirement session's §7/§9. Source:
`docs/status/2026-11-11_09-14_v4-retirement-and-v5-master-merge-status.md`.

<!-- check-rows:ledger -->

| #  | Task                                                                      | Evidence / notes                                                                                                                     | Priority | Status  |
| -- | -------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------- |
| R1 | Fix `website.yml` import-path guard + redeploy                             | hard-coded `/v4` guard trips on a historical `/v3` mention in `changelog.mdx`; Website RED since 2026-11-08 (user-facing regression)  | High     | 🔴 TODO |
| R2 | Validate release branches in CI                                            | add `release/**` (or all branches) to `ci.yml` triggers — the root cause of the red-forever `v4.3.0` tag                              | Medium   | 🔴 TODO |
| R3 | `workflow_dispatch` dry-run path in `release.yml`                          | run tests+lint without creating a release                                                                                            | Medium   | 🔴 TODO |
| R4 | `scripts/pre-tag-check.sh`                                                 | read the release workflow, run `check-all` — the "read the workflow before tagging" lesson, mechanized                                 | Medium   | 🔴 TODO |
| R5 | Add `go vet` (and `go mod verify`) to `check-all`                          | not in the gate today                                                                                                                 | Medium   | 🔴 TODO |
| R6 | Raise samber-linter health coverage above 33%                              | or add an explicit `allow` rationale for the two eager DI services                                                                    | Low      | 🔴 TODO |
| R7 | Unit test for `normalizeTableSpacing`                                      | aligned vs unaligned input; guards the drift-guard fix (`0ce8112`)                                                                    | Low      | 🔴 TODO |
| R8 | Verify pkg.go.dev renders v4.3.1 + v5.0.0                                  | confirm the v4 deprecation banner visually; fetch from the proxy                                                                      | Low      | 🔴 TODO |

---

## Config-Truth Follow-ups (2026-11-02)

From the config-truth execution (plan:
[docs/planning/2026-11-02_05-36_config-truth-pareto-execution-plan.md](docs/planning/2026-11-02_05-36_config-truth-pareto-execution-plan.md)).

<!-- check-rows:ledger -->

| #   | Task                                                    | Notes                                                                                                                                                                                                                                                         | Priority | Status    |
| --- | ------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | --------- |
| CT1 | cmdguard-lint CG009 idea: example configs stay truthful | A rule that parses a repo's committed example configs (e.g. `config.example.yaml`) and rejects keys the CLI's config vocabulary does not accept. Needs cross-file type info (or a mode that shells `describe`/`ConfigKeys` from a built binary), so it is an IDEA, not a task — scope it in the lint backlog before building. | Low      | 💡 IDEA   |
| ~~CT2~~ | ~~Decide the `schema` sub-module release~~          | ~~Ship as `schema/v0.1.0` alongside v4.3.0, or keep workspace-local until the API settles.~~ DONE 2026-11-08-era: held workspace-local through v5.0.0, then tagged `schema/v0.1.0` (tag exists in the repo).                                               | Medium   | ✅ DONE   |
| ~~CT3~~ | ~~Wire `examples/config-explain` into README + website~~ | ~~The demo exists and runs; it is not linked from README.md or the website docs yet.~~ DONE 2026-11-11: README Examples section now lists config-explain + docs-generator + taskctl (website page wiring rides the W-rows).                              | Low      | ✅ DONE   |

---

## Harvested (2026-11-07 docs-health pass)

Items pulled forward from the 2026-06→2026-09 status reports during the
annotate/archive sweep. Everything else those reports contained was either
resolved (struck with evidence in the report itself) or routed to ROADMAP
"Harvested from 2026-0X status reports".

<!-- check-rows:ledger -->

| #      | Task                                                    | Notes                                                                                                                                                                                                                                                                                                                                                                             | Priority | Status  |
| ------ | ------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------- |
| ~~H4~~ | ~~Website hero-code sync guard~~                        | ~~`website/src/data/hero-code.ts` (raw string) vs the inline highlighted copy in `HeroSection.astro` — no test or derive step keeps them in sync~~ DONE: obsolete premise — `HeroSection.astro` imports and renders `heroCode` (single source, no hand-written copy) and `scripts/check-docs-api.py:81` `check_hero_single_source()` gates it. Verified 2026-11-11.                 | Low      | ✅ DONE |
| ~~H6~~ | ~~Audit `docs/COMPARISON.md` prose for v4 accuracy~~    | ~~Header still says "Last Updated: 2026-07-07" while claiming v4-era content; the full prose audit was left open by three 2026-08-06 reports~~ DONE: header now reads 2026-11-08 (v5.0.0 era, audit pass 1 of 3) and zero `v4.` matches remain (re-grepped 2026-11-11); remaining audit passes ride the website W-rows.                                                          | Low      | ✅ DONE |
| ~~H9~~ | ~~Fixture test for lint resolveImports dot/blank arms~~ | ~~`lint/walk.go:240,242` arms have 0 coverage~~ DONE 2026-11-08: `lint/crossfile_imports_test.go` covers both arms; uncovered a REAL gap — dot-imported bare constructors (`NewCLI(...)` without a selector) are untraced, pinned as current behavior and tracked as follow-up H9b (trace bare constructor calls when the file dot-imports cmdguard; non-idiomatic pattern, Low). | Low      | ✅ DONE |
| H9b    | Trace dot-imported bare constructors in cmdguard-lint   | `isSelectorCall` requires a selector, so `NewCLI(...)` under a dot import is invisible to CG003/CG004/CG006; pinned by `lint/crossfile_imports_test.go` (flips to expect findings when fixed). Non-idiomatic pattern; Low priority.                                                                                                                                               | Low      | 🔴 TODO |

---

## Lint Sub-Module (`lint/`)

Harvested from `docs/status/2026-10-04_07-11_cmdguard-lint-sub-module-build.md` §f
(full 50-item backlog lives there; these are the near-term actionable ones).

<!-- check-rows:ledger -->

| #  | Task                                                         | Notes                                                                                                                                                                                                                                                                                                                                                                                                                | Priority | Status      |
| -- | ------------------------------------------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ----------- |
| L6 | golangci-lint module-plugin distribution layer               | Scoped 2026-10-09: requires an `analysis.Analyzer` bridge (golangci/plugin-module-register, per go-humanize-linter pattern: `lint/plugin/` package + `.custom-gcl.yml` + `.golangci.custom.yml`). One pass = one package's files, so cross-file tracing carries over; checks need a pass.Files→project adapter. NOT a thin wiring task — do not half-ship. Also needs the custom-gcl build verified against corpora. | Low      | 🔴 TODO     |
| ~~L8~~ | ~~Bump `CurrentMajor` in `lint/walk.go` when cmdguard v5 ships~~ | ~~Standing duty; CG002 tracks the constant. Rule IDs CG001–CG008 are load-bearing (suppressions key on them) — never renumber~~ DONE 2026-11-08: `CurrentMajor` is `"v5"`, shipped fleet-wide in `lint/v0.3.0` + the BuildFlow bump (CG002 flags v1–v4 imports).                                                                                                                                                     | —        | ✅ DONE     |

---

## Security

<!-- check-rows:ledger -->

| #  | Task                                        | Notes                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 | Priority   | Status     |
| -- | ------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------- | ---------- |
| ~~S2~~ | ~~Bump toolchain to go1.27.2+ (6 stdlib vulns)~~ | ~~govulncheck (2026-10-09): GO-2026-6599/6600/6611/6613/6617/6603 reachable via stdlib `net/http`+`html/template` @go1.27.1; all fixed in go1.27.2. Blocked: nixpkgs `go_1_27` was still 1.27.1.~~ DONE 2026-11-11: blocker gone — `nix eval nixpkgs#go_1_27.version` = 1.27.3; root go.mod already pins `toolchain go1.27.3`; the 7 sub-module go.mods now carry `toolchain go1.27.3` too (this pass). Residual: re-run govulncheck to confirm the stdlib findings clear. | Medium     | ✅ DONE    |

---

## Blocked on Upstream

Currently nothing. (The 2026-10 audit-middleware block resolved when
samber-do-auditlog v0.11.0+ shipped `Plugin.RecordCommand`; `AuditMiddleware[T]`
has been in core since — see FEATURES.md Audit Log table.)
