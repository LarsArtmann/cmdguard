# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
Dates are in YYYY-MM-DD format (ISO 8601).

## [Unreleased]

### Added

- **Doc-table drift tripwires.** `docs_tables_drift_test.go` pins the root
  validate-verb table to the live `registerBuiltins()` registry, the 13-tag
  table and the `.agents` tag-set line to the `FlagTag` fields, and the
  website's validate-table copy to the same registry — hand-kept copies can no
  longer drift from code.
- **Release-readiness gate.** `scripts/pre-tag-check.sh` runs the
  release-blocking battery (build, tests, lint, contract, docs, website legs)
  in one command so a tag never ships unverified.
- **Website truth guards.** `scripts/check-guide-sync.sh` (stale-major scan,
  docs-symbol existence, compiled doc samples, per-guide "Where to go next",
  and describe-excerpt byte pinning against a real `taskctl describe` run),
  `scripts/check-website-links.py`, `scripts/check-website-smoke.py`, the
  v4 `RetirementBanner`, and agent-consumable `website/public/llms.txt`.
- **Fleet tooling.** `scripts/fleet-bump.sh <version>` (parameterized
  consumer-bump sweep with build+vet verification per repo) and
  `scripts/fleet-sweep.sh` (fleet evidence collection).

### Changed

- **Adoption-surface overhaul** from the 41-review fleet audit
  (`docs/reviews/2026-11-17_16-14_cmdguard-docs-marketing-ai-skills-41-reviews.html`,
  all 10 findings resolved): README gains a "Replace your hand-rolled"
  conversion table, a replayable "verify before trusting" block, the 13-tag
  Flag Tags Reference (adds `prompt`/`secret`), and the complete 10-rule
  validate registry; SKILL.md sentinel counts derive from the generated
  ERROR_REFERENCE (62, closing the 59-vs-62 split); v5.1.0-only flags carry
  since-annotations across README, website, and SKILL.md; the machine contract
  document walks the full agent loop with real output; CI generates and
  commits the README coverage badge (`.github/badges/coverage.svg`), retiring
  the hand-maintained static badge.

- TODO_LIST purge (2026-11-21): F1
- TODO_LIST purge (2026-11-21): F3
- TODO_LIST purge (2026-11-21): F7
- TODO_LIST purge (2026-11-21): W1
- TODO_LIST purge (2026-11-21): W3
- TODO_LIST purge (2026-11-21): W5
- TODO_LIST purge (2026-11-21): FM2
- TODO_LIST purge (2026-11-21): FM3
- TODO_LIST purge (2026-11-21): FM4
- TODO_LIST purge (2026-11-21): FM5
- TODO_LIST purge (2026-11-21): FM6
- TODO_LIST purge (2026-11-21): R1
- TODO_LIST purge (2026-11-21): R2
- TODO_LIST purge (2026-11-21): R3
- TODO_LIST purge (2026-11-21): R4
- TODO_LIST purge (2026-11-21): R5
- TODO_LIST purge (2026-11-21): R6
- TODO_LIST purge (2026-11-21): R7
- TODO_LIST purge (2026-11-21): R8
- TODO_LIST purge (2026-11-21): CT2
- TODO_LIST purge (2026-11-21): CT3
- TODO_LIST purge (2026-11-21): H4
- TODO_LIST purge (2026-11-21): H6
- TODO_LIST purge (2026-11-21): H9
- TODO_LIST purge (2026-11-21): L8
- TODO_LIST purge (2026-11-21): S2
- TODO_LIST purge (2026-11-21): FM7
- TODO_LIST purge (2026-11-21): FM8
- TODO_LIST purge (2026-11-21): | SU22 | Land the 2026-11-08_02-29 archive batch (~65 table verdicts; A35 work), return `max_unarchived` 23→20                                                                                                                          | cadence bumped 20→22→23 same day while the batch pends; this row is the path back down (20-43 report preamble)                                                                      | 🔴 TODO | — done 2026-11-20 — 2026-11-08_02-29 annotated (68 verdicts) and archived
- TODO_LIST purge (2026-11-21): | SU23 | Archive batch: annotate the two 2026-10-04 docs/status reports (60+65 bare prose items) inline per D3, then `git mv` to `docs/status/archived/` — greens the cadence gate (currently red on master: 23 unarchived, oldest 47d) | docs/status/2026-10-04_07-11_cmdguard-lint-sub-module-build.md + 2026-10-04_07-13_todo-list-execution-middleware-deprecations-audit-upstream.md; docs-gates cadence rc=1 2026-11-20 | 🔴 TODO | — done 2026-11-20 — both reports annotated (60+65) and archived; the cadence age bar additionally required the 2026-10-04_12-21 batch (landing same day)
- TODO_LIST purge (2026-11-21): | SU24 | Watch first master CI run: SVG Generate/Commit jobs green, push lands, served image resolves                                                                                                                                   | new ci.yml SVG jobs not yet exercised in real CI; the image 404s until then (2026-11-20_12-29 §f1)                                                                                  | 🔴 TODO | — done (CI 44853943250 7/7 green; badge commit c902c78 landed; raw SVG URL renders)
- TODO_LIST purge (2026-11-21): | SU25 | Drift tripwire: the ten-entry Validate table tracks flags_validate.go                                                                                                                                                          | docgen_drift pattern; hand-kept copy today (2026-11-20_12-29 §f2)                                                                                                                   | 🔴 TODO | — done 2026-11-20 — pkg/cmdguard/v5/docs_tables_drift_test.go TestValidateRulesTrackedInDocs
- TODO_LIST purge (2026-11-21): | SU26 | Drift tripwire: the 13-entry tag table tracks FlagTag                                                                                                                                                                          | FlagTag struct fields are the source; hand-kept copy (2026-11-20_12-29 §f3)                                                                                                         | 🔴 TODO | — done 2026-11-20 — TestFlagTagFieldsTrackedInDocs reflects FlagTag and pins README + .agents tag-set
- TODO_LIST purge (2026-11-21): | SU27 | Build-time include: one shared Validate-table source for docs-site + root docs                                                                                                                                                 | the fleet report's cross-surface plumbing proposal; two hand-kept copies (2026-11-20_12-29 §f4)                                                                                     | 🔴 TODO | — done 2026-11-20 — the same drift test pins the website copy to the one live registry
- TODO_LIST purge (2026-11-21): | SU28 | Site mirrors: hand-rolled + verify sections + since-parity in the features data                                                                                                                                                | the site is behind the 2026-11-20 execution pass (2026-11-20_12-29 §f5)                                                                                                             | 🔴 TODO | — done 2026-11-20 — migrating-from-cobra.mdx hand-rolled + verify sections; features.ts since-parity card
- TODO_LIST purge (2026-11-21): | SU29 | machine-contract.mdx: real taskctl describe bytes + MCP-facing recipe                                                                                                                                                          | real bytes captured 2026-11-20; fictional myapp JSON still shown (2026-11-20_12-29 §f6)                                                                                             | 🔴 TODO | — done 2026-11-20 — machine-contract.mdx real taskctl excerpt + MCP recipe
- TODO_LIST purge (2026-11-21): | SU30 | Reciprocity pass 2: samber-do-best-practices / templ-components entries + sibling-surface mentions                                                                                                                             | pass 1 landed go-cqrs-lite, cqrs-htmx, go-aichat (2026-11-20_12-29 §f7)                                                                                                             | 🔴 TODO | — done 2026-11-20 — SKILL entries + sibling commits (SKILLS 7d14347, templ-components bc369053)
- TODO_LIST purge (2026-11-21): | SU31 | SVG-commit hardening: in-run XML sanity; autostash + retry-once                                                                                                                                                                | push race with the auto-commit daemon designed-for but unproven (2026-11-20_12-29 §f8)                                                                                              | 🔴 TODO | — done 2026-11-20 — ci.yml SVG sanity greps + autostash + retry-once
- TODO_LIST purge (2026-11-21): | SU32 | Consumers claims table: re-verify CG001-CG008 against the lint defs; derive-or-drop the taskctl-87 count                                                                                                                       | repo-documented, not independently re-verified (2026-11-20_12-29 §f9)                                                                                                               | 🔴 TODO | — done 2026-11-20 — row Verified vs lint/rules.go:17-24; 87 derived from grep -c func Test
- TODO_LIST purge (2026-11-21): | SU33 | Sweep QUICKSTART.md + TUTORIAL.md for since-gaps                                                                                                                                                                               | the 2026-11-11 positional-args flags may be missing from the walkthroughs (2026-11-20_12-29 §f10)                                                                                   | 🔴 TODO | — done 2026-11-20 — since-v5.1.0 notes in QUICKSTART.md + TUTORIAL.md
- TODO_LIST purge (2026-11-21): | SU34 | Resolution banner on the 2026-11-17 fleet report: all ten findings resolved 2026-11-20                                                                                                                                         | historical report shows every finding Open (2026-11-20_12-29 §f11)                                                                                                                  | 🔴 TODO | — done 2026-11-20 — banner added to the fleet report
- TODO_LIST purge (2026-11-21): | SU35 | Run the complete build-lint-suite pass                                                                                                                                                                                         | the docs-session ran fmt + targeted checks only (2026-11-20_12-29 §f12)                                                                                                             | 🔴 TODO | — done 2026-11-20 — nix run .#check-all green + full golangci-lint 0 issues
- TODO_LIST purge (2026-11-21): | SU36 | Decide the two hand-maintained 89.3% numbers: date stamp or link the SVG                                                                                                                                                       | contributor-facing copies, unlike the now-generated SVG (2026-11-20_12-29 §f13)                                                                                                     | 🔴 TODO | — done 2026-11-20 — AGENTS.md mentions date-stamped, badge named live source
- TODO_LIST purge (2026-11-21): | SU37 | Resolve Codecov: activate the repo there, or delete the redundant upload job                                                                                                                                                   | uploads succeed via OIDC but the external SVG serves unknown; owner question (2026-11-20_12-29 §f14)                                                                                | 🔴 TODO | — done 2026-11-20 — decision: Codecov upload removed (repo not activated there; CI SVG is the source of truth)
- TODO_LIST purge (2026-11-21): | SU38 | Push master (or confirm daemon-only pushes) and watch the first green CI run: badge steps green, SVG committed, raw badge URL resolves                                                                                       | ci.yml fixes from runs 44768346938/44768346539 never executed; master ahead 10 (2026-11-20_20-15 §f1)                                                                               | 🔴 TODO | — done (CI 44853943250 7/7 green; c902c78 pushed; badge live at 86.3%)
- TODO_LIST purge (2026-11-21): | SU39 | Re-run `nix run .#check-all` after the final test-file format fix to close the full-suite caveat                                                                                                                            | check-all ran before the golines fix to docs_tables_drift_test.go (2026-11-20_20-15 §f2)                                                                                            | 🔴 TODO | — done (full check-all green 2026-11-21: build + 16-pkg race suite + lint 0-issue after treeWalkingRule construction fix + vet + tidy + verify + guide-sync + dep-tables)
- TODO_LIST purge (2026-11-21): | SU41 | Push the three sibling commits (SKILLS adf2e80 + 7d14347 pairing; templ-components bc369053) so fleet surfaces are visible fleet-wide (repo: SKILLS, repo: templ-components)                                                 | SKILLS ahead 2, templ-components ahead 1 (2026-11-20_20-15 §f4)                                                                                                                     | 🔴 TODO | — done (all three verified on origin 2026-11-21)

- TODO_LIST purge (2026-11-21): | SU43 | `docs-gates purge --apply` for the 15 struck rows into one pre-flighted CHANGELOG entry | strike ritual's optional closing step, not run (2026-11-20_20-15 §f6) | 🔴 TODO | — done (2026-11-21: 47 struck rows purged into the CHANGELOG entry; battery green after annotating the dissolved-harvest items at their sources)
- TODO_LIST purge (2026-11-21): | SU45 | cue-policy: exempt reusable-call jobs (`uses:` at job level) from required runs-on/steps in ci_workflow.cue + release_workflow.cue, then re-vet cmdguard workflows and audit other fleet repos for the same rejections (repo: cue-policy) | the pre-commit findings gate blocked every workflow-staging commit on 2 pre-existing cue-vet errors (2026-11-21_01-56 §f7,§f25) | 🔴 TODO | — done (2026-11-21: cue-policy c1482c8 discriminated-disjunction fix; cmdguard ci.yml + release.yml vet clean; pre-commit hook passes without --no-verify)

- TODO_LIST purge (2026-11-21): | SU47 | SKILL deeper-truth: link the website guides from SKILL.md "Where deeper truth lives" (plan T20/T20.1)                                                                        | one site link exists (type-safe-flags via SU30); the deeper-truth section still points repo-only                                                                                                                                          | 🔴 TODO                                                                                                                         | — done (2026-11-21: deeper-truth section links all six rendered mirrors — quick-start, error-handling, migrating-v4-to-v5, migrating-from-cobra, type-safe-flags, contributing, api-reference; fan-out link verified)

- TODO_LIST purge (2026-11-21): | SU1 | Re-feed CHANGELOG `[Unreleased]` with post-5.1.0 work | pre-tag script, describe harness, retirement banner, llms.txt, fleet scripts, lint dot-import fix, website guard — all shipped after the tag, none listed | 🔴 TODO | — done (2026-11-21: [Unreleased] re-fed — release gate, website truth guards incl. describe excerpt pinning, fleet tooling, retirement banner, llms.txt; commit fbacf76)
- TODO_LIST purge (2026-11-21): | SU2 | FEATURES.md: add rows for the new tooling (describe harness, pre-tag gate, website smoke, fleet scripts, rulesets) | the 14-57 rebuild predates today's additions | 🔴 TODO | — done (2026-11-21: Repository Tooling section added with file evidence; commit fbacf76)
- TODO_LIST purge (2026-11-21): | SU5 | AGENTS.md Quick Start: add `scripts/pre-tag-check.sh` and `scripts/check-website-smoke.py` entries | both are now standing gates; AGENTS.md still documents only the older scripts | 🔴 TODO | — done (2026-11-21: pre-tag-check.sh + check-website-smoke.py added to Quick Start; commit f5602ee)
- TODO_LIST purge (2026-11-21): | SU6 | One FULL `check-all` run exercising the new vet + go-mod-verify legs end-to-end | legs verified standalone + flake evaluates; no complete run since the flake edit (pre-tag used `--fast`) | 🔴 TODO | — done (2026-11-21: full check-all EXIT=0 after treeWalkingRule exhaustruct fix + dep-budget line; SU39)
- TODO_LIST purge (2026-11-21): | SU7 | Unit test for `DescribeContractFromBinary` (process-capture helper) | shipped untested — the two pilots exercise the library path only | 🔴 TODO | — done (2026-11-21: TestDescribeContractFromBinary + testdata/describeapp fixture; env-prefix and greet assertions; commit bee6a57)
- TODO_LIST purge (2026-11-21): | SU8 | Gitignore audit: `.crush/shell-output/` + `website/src/styles/global.out.css` build artifacts | global.out.css churn got committed by the daemon mid-build today; verify both are ignored or deliberately tracked | 🔴 TODO | — done (2026-11-21: .crush/shell-output/ ignored, global.out.css untracked (1566 lines of generated Tailwind); commit f5602ee)
- TODO_LIST purge (2026-11-21): | SU10 | Merge Dependabot #21 (submodule tidy landed, checks green) after owner glance | all failing checks resolved; DontMergeMeYet label allows merge | 🔴 TODO | — done (2026-11-21: PR #21 state CLOSED on GitHub)
- TODO_LIST purge (2026-11-21): | SU14 | Parameterize `scripts/fleet-bump-v5.1.0.sh` → `scripts/fleet-bump.sh <version>` for future cycles | one-release name today; add vendor + dirty-tree handling while at it (the two-pass discovery this cycle) | 🔴 TODO | — done (2026-11-21: git mv to scripts/fleet-bump.sh <version>, arg validation; commit f5602ee)
- TODO_LIST purge (2026-11-21): | SU15 | Pin firebase-tools in `website.yml` deploy step (exact version or nix binary) | CI ran 15.35.0 vs locally proven 15.28.1; unpinned `pnpm add -g` = repro + supply-chain gap (20-43 report §b1) | 🔴 TODO | — done (2026-11-21: firebase-tools@15.37.0 in website.yml deploy; commit f5602ee)

- TODO_LIST purge (2026-11-21): | SU48 | taskctl: store plain `Long`/`Example` text and style at help-render time — describe currently embeds ANSI-escaped styled help in the machine contract (found while pinning the SU44 excerpt to live bytes) | the contract `long` field carries raw `\u001b[38;5;…` escape sequences plus lipgloss padding; hostile for contract-driving agents (2026-11-21 excerpt repair) | 🔴 TODO | — done (2026-11-21: core fix instead of the row taskctl-side premise — glamour styles even plain prose, so plain text could not clean the contract; v5 now snapshots pristine Long/Example into the cmdguard.io/pristine-help annotation before transforms run and describe/GenerateDocs serve that text, terminal help unchanged; taskctl contract 5728 to 211 bytes on add, zero ANSI; pristine_help_test.go covers transform, repeat-Execute, docgen, passthrough)

- TODO_LIST purge (2026-11-21): | SU42 | annotate.go: per-row TSV validation (kind enum + raw-value rule) with row-numbered errors, extending adf2e80 (repo: docs-gates)                                              | a spec with kind `done` failed full-file and late (2026-11-20_20-15 §f5)                                                                               | 🔴 TODO | — done (2026-11-21: SKILLS/docs-health/assets/annotate.go c337064 — parse loop validates every row up front via validateRepairSpec, kind vocabulary h/v/p/w/n/r/s plus per-kind value requirements, errors carry the spec line number and fire before any target file is touched; f5 shape pinned in self-test, negative probe prints repair: spec line 3: bad kind done; docs-gates parity suite green)

- TODO_LIST purge (2026-11-21): | SU4  | Release the lint module (next tag) carrying the dot-import constructor tracing fix                                                                                           | behavior change is on master + the locally built binary; no published tag carries it yet; BuildFlow provider pin follows the tag                       | 🔴 TODO | — done (2026-11-21: lint/v0.4.0 tagged at ff6913d, 14/14 CI checks green on the release commit; proxy indexed, clean-dir go get resolves v0.4.0 with zero pseudo-versions; delta = dot-import constructor tracing + CoverageReporter filesScanned evidence + go-linter-sdk v0.6.1 + cmdguard/v5 v5.1.0 pin; lint/CHANGELOG.md carries the section; dev-shape replace restored at de05f31)

- TODO_LIST purge (2026-11-21): | SU16 | Migrate website deploy to Workload Identity Federation (GitHub OIDC), retire the long-lived SA JSON key + secret                                                             | deploy works today via `FIREBASE_SERVICE_ACCOUNT` secret; WIF kills key + rotation burden (20-43 report §b2)                                           | 🔴 TODO | — done 2026-11-21 (WIF live: pool github/provider cmdguard with repository condition, iamcredentials+STS APIs enabled, two dispatch deploys green; SA key deleted, FIREBASE_SERVICE_ACCOUNT secret removed)
- TODO_LIST purge (2026-11-21): | SU17 | Least-privilege the deploy SA: per-site deployer on `cmdguard` only                                                                                                          | SA is project-wide `firebasehosting.admin` over ~20 hosting sites on lars-software (20-43 report §b2)                                                  | 🔴 TODO | — done 2026-11-21 (amended mechanism: Firebase Hosting has no per-site IAM scope and no deployer role; custom projects/lars-software/roles/cmdguardHostingDeployer swapped in for firebasehosting.admin, no sites.create/sites.delete; dispatch deploy green)
- TODO_LIST purge (2026-11-21): | SU21 | Root-cause the ~60s post-deploy 503/stale window; eliminate or document it                                                                                                   | content-identical deploy, 60s brownout, "edge propagation" assumed never verified (20-43 report §d6)                                                   | 🔴 TODO | — done 2026-11-21 (documented: not reproducible; 5s polling across two live deploys saw the new release visible within <=6s, zero 503/stale samples; deploy-meta.json max-age=3600 is a browser directive, the edge invalidates on release)

### Removed

- **Five dead sentinel errors** that no code path ever returned:
  `ErrEditorTempFile`, `ErrEditorWrite`, `ErrEditorRun`, `ErrEditorRead`, and
  `ErrNilValue` (`nil_value` machine code). v3 had already cut the editor
  sentinels as dead weight; the v5 scaffold re-added them mechanically and
  they sat in the public API and the generated `docs/ERROR_REFERENCE.md` while
  being unreachable. Surfaces a strict `erraudit tree ./... --max` audit
  flagged as unreferenced (2026-11-22). The sentinel fact table now covers 57
  sentinels, the `editor` envelope class is gone, and a new
  `TestNoDeadSentinels` tripwire fails when a declared sentinel is never
  referenced outside its declaration and the facts table. No runtime behavior
  changed: no caller could ever have matched these via `errors.Is`.

### Fixed

- **Machine contract and generated docs no longer carry ANSI-styled help.**
  Help transforms (e.g. the glamour sub-module's markdown renderer) overwrite
  the live cobra `Long`/`Example` fields with terminal-styled bytes, and
  `describe`/`GenerateDocs` read those same fields — so the JSON contract
  embedded raw `\u001b[38;5;…` escapes and lipgloss padding (taskctl's
  `list` entry: 7151 styled bytes for 353 characters of source). Core now
  snapshots the author's original text into a cobra annotation before any
  transform runs and serves that pristine text from the contract and docs;
  terminal help rendering is unchanged. Repeated `Execute` calls never
  re-stash styled text.
- **Master CI red trio from 2026-11-20.** The coverage-stamp commit step
  aborts no more on a dirty index (`git pull --rebase --autostash` + one
  retry + an SVG sanity grep gate); the branded-ID gate pin
  `@linter/v0.5.1` was an invalid `go run` version string and now resolves
  as `@v0.5.1`; the Codecov upload was REMOVED (repo not activated there:
  tokenless uploads rejected, its external SVG served "unknown" — the
  CI-committed SVG stays the source of truth).
- **Website builds again.** A dependency-bump re-resolution paired
  `@emnapi/runtime@2.0.0-beta.2` with `@emnapi/core@1.11.3`, crashing
  astro's typechecker on boot ("this.bridge.setLastError is not a function");
  pnpm-workspace.yaml overrides now pin the emnapi family to 1.11.3.

## [5.1.0] - 2026-11-11

### Added

- **Strict positional-args control.** `WithStrictPositionalArgs()` makes a
  command reject undeclared positional args (inert when an explicit validator
  like `WithExactArgs` is already set); `WithArgsPassthrough()` explicitly
  accepts and ignores extras. Pinned by `cli_args_strict_test.go`.

### Fixed

- **`Builder` is safe for concurrent use.** `Add`/`AddParent`/`Try`/`Check`/
  `Build` serialize on a mutex, so queueing commands from multiple goroutines
  followed by one `Build()` no longer races on the internal slices.
- **`doctor`'s schema-validator doc example** referenced the removed
  `schema.ValidateFile` shape; it now shows `schema.FromConfig` + `Emit`.
- Release workflow lint invocation (`args: run ./...` → `./...`) — the action
  already supplies `run`, so the duplicated subcommand made every `v*` release
  check exit 7 since v4.1.0 (fixed on master and `release/v4`).

## [4.3.1] - 2026-11-11

Final release of the retired v4 major line, cut on the `release/v4` branch.
v4.3.0 pinned the root module's sibling requires to published versions (v4 had
been unresolvable through `go get` since v4.2.1, whose go.mod pointed at
unfetchable zero pseudo-versions) and marked the package deprecated; v4.3.1
carries the same API with the release workflow's lint invocation fixed. `master`
now tracks v5; v4 is frozen on `release/v4`.

## [5.0.0] - 2026-11-08

The paved-road major: the module path moves to `/v5`, production batteries
become the defaults, the per-command error ladder is replaced by a Builder,
and three identifiers get the names they should always have had.

### Changed

- **BREAKING: module path** — `github.com/larsartmann/cmdguard/v4` →
  `github.com/larsartmann/cmdguard/v5` (package `pkg/cmdguard/v4` →
  `pkg/cmdguard/v5`). For most consumers this is an import-path edit;
  `docs/migration_v4_v5/codemod.sh` performs it mechanically.
- **BREAKING: production batteries are ON by default.** `NewCLI` now enables:
  strict command validation (every command needs a short description),
  signal handling with graceful DI shutdown, strict config keys (unknown
  config-file keys are errors), actionable error hints, and the
  auto-registered `doctor` command (`version` appears when a CLI version is
  set). Every flip has an opt-out — `WithoutStrictValidation()`,
  `WithoutSignalHandling()`, `WithoutGracefulShutdown()`,
  `WithoutStrictConfigKeys()`, and `WithoutBatteries()` (the master switch).
  The full matrix is pinned by `cli_defaults_matrix_test.go`.
- **BREAKING: renames** — `CommandInfo` → `CommandMetadata`,
  `TypeHandler` → `TypeCodec` (and `TypeHandlerFunc` → `TypeCodecFunc`),
  `HuhRunner` → `HuhPrompter` (prompts sub-module). `CLI.SetConfig` is
  removed; `RegisterInScope` is the generic `RegisterInScope[T]`. The
  codemod rewrites all of these.
- **Graceful-shutdown semantics.** The automatic DI shutdown fires ONLY when
  the execution context was cancelled (the SIGINT/SIGTERM path); a normally
  completing `Execute` never tears down the scope, so repeated
  `ExecuteWithArgs` on one CLI is safe under the defaults. Explicit teardown
  is `CLI.Shutdown(ctx)`.
- **Repeated-Execute invariants.** Every run rebuilds the flow context and
  derives the command context from the current run (`execCtx`), fixing the
  cobra stale-subcommand-context shadowing that made every command after the
  first fail health checks with `context canceled`.
- `docs/QUICKSTART.md` and `docs/TUTORIAL.md` teach the v5 defaults.

### Added

- **Builder registration** — `NewBuilder[T](cli)` with free functions
  `Add`/`AddParent`/`Try` and methods `Check`/`Build`: constructor errors are
  collected, `Build()` attaches atomically or reports the first error.
  `AddCommand`/`AddCommands` remain for v4-style call sites.
- **Migration kit** — `docs/MIGRATION_v4_v5.md` (human guide) plus the
  additive-only codemod script; dry-run verified against erraudit (full
  suite) and timesheets.
- `cmdguard-lint` `CurrentMajor` is `v5`: CG002 now flags cmdguard v1–v4
  imports as stale majors (ships fleet-wide with `lint/v0.3.0`).
- samber-linter baseline ratchet (`.samber-linter-baseline.json`, 33% floor)
  wired into `check-all`.

### Fixed

- **Empty config files no longer crash startup.** A zero-byte or
  whitespace-only config file was parsed as an empty JSON document and failed
  hard; it is now treated uniformly as "absent" (skipped, like a missing
  file).
- **Flightrecorder snapshots survive `Execute` returning.** Async capture now
  runs on `context.WithoutCancel(ctx)` — under the v5 signal defaults the
  original context is cancelled at Execute return, which previously raced
  the in-flight snapshot write.
- **`doctor` is idempotent across repeated Execute** (stale cached flow
  context from run 1 leaked its cancelled NotifyContext into later runs).
- `.golangci.yml`: dead `pkg/cmdguard/v2` exclusions removed, duplicate
  `example_test.go` forbidigo entry removed, `TODO(v5)` godox exclusion
  removed (all five markers resolved by the renames above).

---

## [Backfill] - notes never sectioned before v5 (pending owner §g2)

The entries below landed across v4.1.0–v4.2.x but were never given a
release section. They are preserved verbatim until the owner decides
whether to fold them into the historical sections (see §g of
`docs/status/2026-11-08_08-32_whole-run-audit-and-next-50.md`).

### Changed

- **BREAKING (retrospective entry — the change landed in v4.1.0 without a
  CHANGELOG note).** `CLI.HealthCheck` now takes a `context.Context`:
  `HealthCheck(ctx context.Context) error` instead of `HealthCheck() error`,
  so services implementing `do.HealthcheckerWithContext` can honour
  cancellation. Callers must pass a context (usually `cmd.Context()`);
  compile-time breakage only, no behavior change for zero-arg callers beyond
  the signature. `CLI.HealthCheckResults` kept its no-context signature.
  Recorded here because the v4.1.0 release notes omitted it.

### Added

- **Stable machine error codes (`ErrorCode(err)`, envelope `error_code`).**
  Every sentinel now carries a stable lowercase snake_case code
  (`required_flag`, `config_unknown_key`, …) that survives Go identifier
  renames, so scripts and agents can switch on failures without parsing
  messages. The machine-readable error envelope (`--output json|yaml|toml`)
  gains an additive `error_code` field; the human error path is unchanged.
  Sentinel facts (identifier, envelope class, code, hint) live in ONE
  canonical table, so hint coverage, error classification, and codes cannot
  drift apart — and run-time errors that previously classified as `"unknown"`
  in the envelope (invalid URL/email/duration values, command timeouts,
  editor failures) now carry their real class.
- **Generated, drift-guarded error reference.** `docs/ERROR_REFERENCE.md`'s
  sentinel table is rendered from the same canonical fact table via
  `GenerateErrorReference()`; a test fails when the committed doc stops
  matching the code. Regenerate with
  `go test ./pkg/cmdguard/v5 -run TestErrorReferenceDoc -update-golden`.
- **Strict config keys, opt in (`WithStrictConfigKeys()`).** Unknown keys in a
  config file are now an error instead of a silent no-op, with the offending
  dotted key, the loaded file path, and a "did you mean" suggestion for likely
  typos (`ErrConfigUnknownKey`). Accepted spellings match what config loading
  already accepts (canonical field name, flag name, and case/hyphen/underscore
  variants), and nested/array keys are named in full (`servers[0].port`).
  Off by default in v4 because a shared config file may legitimately carry
  extra keys; **`WithBatteries()` turns it on** as part of the production
  defaults (v5 flips the global default). `doctor` reports config problems
  (missing file, parse failure, unknown keys) alongside DI health via the new
  `CLI.ConfigDiagnostics()`, and `WithDoctorSchemaValidator` attaches an
  external JSON Schema check.
- **Config provenance: the new `config explain` command (`WithConfigExplainCommand()`).**
  Answers "where did this value come from?" per field: flag, environment,
  config file, tag default, or unset — resolved with the same precedence the
  runtime uses. Table output for humans, `--output json` for scripts.
  `CLI.ExplainConfig()` exposes the same data programmatically. Fields tagged
  `secret:"true"` are redacted to `***` in every view while still reporting
  their source, so provenance can be audited without leaking credentials.
- **`secret:"true"` struct tag.** Marks a config value that must never be
  logged or committed. It appears in the machine contract (`describe`), the
  generated docs, the config schema, and drives `config explain` redaction.
- **`schema` sub-module** (`github.com/larsartmann/cmdguard/schema`): emits a
  JSON Schema draft 2020-12 document describing a config struct (`FromConfig`)
  or a whole CLI (`FromCLI`), including nested objects, arrays, enums,
  min/max constraints, and cmdguard hints (`flag`, `env`, format, `secret`)
  under `x-cmdguard`. Properties always reject unknown keys
  (`additionalProperties: false`), matching strict loading — a schema typo is
  the same typo the binary rejects. Output is deterministic and editor-friendly
  (`Emit`/`Marshal` write indented JSON). Experimental surface.
- **`v4.ConfigKeys(cfg)`**: the public vocabulary of a config struct — every
  accepted key with its flag, env var, type, format, constraints, default, and
  nested shape. It is the single source shared by strict-key checking and the
  schema emitter, so the two can never disagree.

### Fixed

- **`min`/`max` validate rules now resolve in the flag type's domain.**
  Duration-typed flags (`v4.Duration` and `time.Duration`) compare bounds as
  durations, so `validate:"min=10s,max=1m"` means ten seconds and one minute —
  previously the numeric validators parsed the value text as a plain number
  and failed with "invalid validator parameter" (or, for `time.Duration`,
  silently compared nanoseconds). Textual parsed types (URL, Email, HostPort,
  FilePath) fall back to length bounds instead of failing to parse as
  numbers. String kinds keep the existing length semantics, numeric kinds and
  Port keep numeric semantics. The same resolution now applies to
  `ValidateConfig`, which previously validated the raw tag text without the
  string length rewrite.
- **Impossible flag defaults now fail at construction, not on the first
  run.** `NewFlagRegistry` (and therefore `NewCommand`/`NewCLI`) validates
  every tag default against that tag's `values` and `validate` rules and
  rejects violations with the new `ErrInvalidFlagDefault` sentinel. Defaults
  whose `validate` tag names a validator registered after construction are
  skipped at construction time and remain enforced at run time, so the
  register-after-construction pattern keeps working.
- **Config values in camelCase, lowercase, or snake_case no longer lose the
  race against flag defaults.** Key matching accepted `{"LogLevel": "debug"}`
  (case-insensitive JSON decoding) while the "which fields did the file set?"
  bookkeeping matched only the exact flag spelling (`log-level`), so those
  fields were treated as unset and flag parsing overwrote the config value
  with the tag default. Every accepted spelling is now normalized to the
  canonical key before decoding, and the same accepted-spelling vocabulary is
  used for both the value assignment and the bookkeeping.
- **Generated markdown docs now carry exactly the facts the machine contract
  carries.** `GenerateDocs` previously rendered from pflag's `FlagUsages`,
  a second representation that silently drifted from `describe` (missing env
  vars, required markers, accepted values, counting and secret markers).
  It now renders from the same source as the JSON contract: an eight-column
  table (flag, type, env, default, values, required, secret, help), with inherited
  flags in their own section and table-escaping so pipes in help text cannot
  break the layout. Root-only CLIs (no subcommands) no longer produce empty
  output.

- **`ExitCode` no longer maps a failing error to exit 0 when an `ExitCoder`
  reports 0 (#11).** Third-party error types (e.g. go-error-family) implement
  `ExitCoder` with 0 as their "no explicit override" sentinel; `extractExitCode`
  trusted it verbatim, so failed CLIs exited 0. A zero code from an `ExitCoder`
  is now treated as "no override" and falls back to the generic failure code 1
  (a nil error still yields 0). The JSON error envelope `code` field follows the
  same rule.
- **Flag values no longer leak across repeated `ExecuteWithArgs` calls on a
  reused CLI (#13).** pflag marks a flag `Changed` forever after its first parse
  and never restores values, so a second execute re-ran with the first parse's
  flags (including write-action flags like `--apply` that the new args never
  mentioned). Every `Execute` now resets per-flag state first: `Changed`
  markers, values (stringSlice values are rebuilt because pflag's append
  semantics make `Set(DefValue)` corrupt them), and the typed config, which is
  restored to its post-constructor snapshot (defaults plus config-file
  overlay). One-execute-per-process CLIs are unaffected. Caveat: raw cobra
  escape-hatch flags of pflag's other appending kinds (intSlice and friends)
  keep their residue value; only their `Changed` marker resets.

---


### Housekeeping

- **TODO_LIST purge (2026-11-08).** Struck done-items removed from
  `TODO_LIST.md` after on-disk verification; substance preserved here.
  Completed this cycle on the `v5-scaffold` branch:

  - Rename `CommandInfo` → `CommandMetadata` T1 — `middleware.go` DONE 2026-11-08: v5-scaffold: repo-wide rename + codemod rule.
  - Rename `TypeHandler` → `TypeCodec` T2 — `type_handler.go` DONE 2026-11-08: v5-scaffold: TypeCodec rename + codemod.
  - Rename `PromptRunner` → `HuhPrompter` (or similar) T3 — `prompts/prompts.go` DONE 2026-11-08: v5-scaffold: HuhPrompter.
  - Remove `SetConfig` (now `// Deprecated:`) T4 — `cli_accessors.go` DONE 2026-11-08: v5-scaffold: SetConfig removed.
  - Replace `RegisterInScope` with generic variant (now `// Deprecated:`) T5 — `scope.go` DONE 2026-11-08: v5-scaffold: generic RegisterInScope[T].
  - DONE 2026-11-02: v4.2.0 + migration note ships enforcement; `/v5` reserved for v5.0.0 (scaffold = plan M11); v5 flips strict + signals/graceful + batteries on, opt-outs stay. Recorded in ROADMAP. P1 — DONE 2026-11-02: v4.2.0 + migration note ships enforcement; `/v5` reserved for v5.0.0 (scaffold = plan M11); v5 flips strict + signals/graceful + batteries on, opt-outs stay. Recorded in ROADMAP.
  - DONE 2026-11-02: `.golangci.yml` codified (lll out, golines max-len 120 owns length; issue caps 0 = unlimited, comment added); go.mod/go.sum deltas = Go 1.27 workspace-aware tidy (sibling replaces added, canonical state verified by full build+test of all 7 modules); no flake.lock change. P2 — DONE 2026-11-02: `.golangci.yml` codified (lll out, golines max-len 120 owns length; issue caps 0 = unlimited, comment added); go.mod/go.sum deltas = Go 1.27 workspace-aware tidy (sibling replaces added, canonical state verified by full build+test of all 7 modules); no flake.lock change.
  - DONE 2026-11-02: `cli_fang_on_test.go` — hint via fang handler (printed exactly once), batteries version/doctor under fang, WithoutErrorHints suppression, batteries×Draconian built-ins pass, JSON envelope carries hint with no standalone line. -race green, 0 lint. P3 — DONE 2026-11-02: `cli_fang_on_test.go` — hint via fang handler (printed exactly once), batteries version/doctor under fang, WithoutErrorHints suppression, batteries×Draconian built-ins pass, JSON envelope carries hint with no standalone line. -race green, 0 lint.
  - DONE 2026-11-02: `describe_golden_test.go` + 3 fixtures in `testdata/` (full tree incl. hidden describe cmd, opt-out, consumer-wins). Byte-stable under -race; schema stability rule documented on CLIDescription; regenerate via `-update-golden`. P4 — DONE 2026-11-02: `describe_golden_test.go` + 3 fixtures in `testdata/` (full tree incl. hidden describe cmd, opt-out, consumer-wins). Byte-stable under -race; schema stability rule documented on CLIDescription; regenerate via `-update-golden`.
  - DONE 2026-11-02: tag v4.2.0 pushed on 847d2c2 (check-all green, CI green); proxy serves it; scratch consumer verified end-to-end (enforcement + hint + describe + batteries from the published tag); GitHub release published. Near-miss caught: GOWORK=off tidy had poisoned sibling requires to pseudo-versions — root restored to real pins, gate made workspace-aware. P5 — DONE 2026-11-02: tag v4.2.0 pushed on 847d2c2 (check-all green, CI green); proxy serves it; scratch consumer verified end-to-end (enforcement + hint + describe + batteries from the published tag); GitHub release published. Near-miss caught: GOWORK=off tidy had poisoned sibling requires to pseudo-versions — root restored to real pins, gate made workspace-aware.
  - Consumer proof: erraudit + timesheets P7 — P1 DONE 2026-11-08: 2026-11-08: erraudit full suite green on v5 (codemod dry-run/tmp), timesheets HEAD cmd green.
  - v5 scaffold + Builder + default flips + renames P8 — P1 DONE 2026-11-08: 2026-11-08: v5-scaffold branch — module /v5, flips ON with Without* opt-outs, Builder, 5 renames, execCtx fix.
  - Codemod + MIGRATION_v4_v5.md P9 — P1 DONE 2026-11-08: 2026-11-08: docs/migration_v4_v5/codemod.sh + docs/MIGRATION_v4_v5.md; dry-runs green.
  - QUICKSTART + TUTORIAL parity P12 — P1 DONE 2026-11-08: 2026-11-08: QUICKSTART + TUTORIAL on v5 defaults (batteries/describe/hints/runtime tags sections).
  - DONE 2026-11-02: `pnpm install` (astro 7.4.4, vite 8.3.5, ts 7.0.2, html-validate 11.18.0) + deep `pnpm update undici zod` → `pnpm audit --prod` clean; `astro build` green (22 pages). P13 — DONE 2026-11-02: `pnpm install` (astro 7.4.4, vite 8.3.5, ts 7.0.2, html-validate 11.18.0) + deep `pnpm update undici zod` → `pnpm audit --prod` clean; `astro build` green (22 pages).
  - `Describe()` / `WriteDescription()` (`describe.go`) plus an auto-registered hidden `describe` command; flag metadata (env, required, values, count, prompt, validate) rides on cobra annotations. Consumer-defined `describe` wins; `WithoutDescribeCommand()` opts out. 7 describe + 6 AddCommands tests. D1
  - `AddCommands[T](cli, cmds...)`: prepares and validates every command before attaching any, accepts mixed flags types, rejects nil entries, and keys duplicate detection on the cobra command name (`"list [path]"` vs `"list"`). `examples/taskctl` registers 11 commands in one call (432→402 LOC). D2
  - `ErrorHint(err)` returns one actionable fix per classified sentinel; the JSON error envelope gains `hint`. A coverage test fails when a classified sentinel has no hint. D3
  - `cmdguard-lint` CG007 (missing short description) and CG008 (missing CLI version); test files and spread option lists are skipped; generated README table refreshed. D4
  - DONE 2026-11-01: `WithBatteries()` (strict validation, signals + graceful shutdown, hints, auto `version`/`doctor`), `NewApp` construction wrapper, `WithoutErrorHints()`. 8 tests. D5 — DONE 2026-11-01: `WithBatteries()` (strict validation, signals + graceful shutdown, hints, auto `version`/`doctor`), `NewApp` construction wrapper, `WithoutErrorHints()`. 8 tests.
  - DONE 2026-11-01: `ValidateFlags` had no production caller — `required`/`values`/`validate` were never enforced. Now wired into root + command parsing, against resolved values (flag/env/default); string `min`/`max` mean length; meta commands exempt. See CHANGELOG "Fixed". D6 — DONE 2026-11-01: `ValidateFlags` had no production caller — `required`/`values`/`validate` were never enforced. Now wired into root + command parsing, against resolved values (flag/env/default); string `min`/`max` mean length; meta commands exempt. See CHANGELOG "Fixed".
  - v5 default flip D7 — Path decided 2026-11-02: module moves to `/v5` with v5.0.0; defaults flip = strict validation + signals/graceful + batteries (opt-outs remain). Codemod + migration guide per plan M13–M16. DONE 2026-11-08: 2026-11-08: flips landed on v5-scaffold (see P8).
  - Test all examples in CI H1 — Low DONE 2026-11-08: 2026-11-08: examples are root-module packages — CI go test ./... covers them (verified).
  - H2 Rewrite or delete the no-op `PromptConfirm` test
  - koanf→JSON round-trip edge-case tests H3 — Medium DONE 2026-11-08: 2026-11-08: koanf_edge_test.go — TOML datetime, YAML anchors/merge, int-vs-float64, empty files.
  - Remove dead `pkg/cmdguard/v2` lint exclusions H5 — Low DONE 2026-11-08: 2026-11-08: dead v2 exclusions removed with the v5 .golangci.yml pass.
  - Create GitHub release pages for v4.0.0 / v4.0.1 (+ v0.2.0 sub-modules) H7 — Low DONE 2026-11-08: 2026-11-08: gh release pages created for v4.0.0 + v4.0.1.
  - Include `examples/taskctl` in `check-all` H8 — Medium DONE 2026-11-08: 2026-11-08: taskctl is a root-module package; check-all + CI cover it.
  - CHANGELOG entry for the `CLI.HealthCheck(ctx)` signature change DONE 2026-11-08: The breaking change shipped silently in v4.1.0 — no entry through v4.2.1 (grep-verified).  2026-11-08: retrospective BREAKING entry added to CHANGELOG `[Unreleased]` → Changed (owner's §g2 answer can move it into `[4.1.0]` if preferred).. W1 — CHANGELOG entry for the `CLI.HealthCheck(ctx)` signature change DONE 2026-11-08: The breaking change shipped silently in v4.1.0 — no entry through v4.2.1 (grep-verified).  2026-11-08: retrospective BREAKING entry added to CHANGELOG `[Unreleased]` → Changed (owner's §g2 answer can move it into `[4.1.0]` if preferred)..
  - Deduplicate `CLI.HealthCheckWithContext` vs `CLI.HealthCheck` W2 — Medium DONE 2026-11-08: 2026-11-08: duplicate removed — CLI.HealthCheck is the ctx-aware form (v5 branch).
  - Resolve the `Scope.HealthCheck()` / `HealthCheckWithContext` split-brain W3 — Medium DONE 2026-11-08: 2026-11-08: Scope.HealthCheck/Results delegate to the WithContext variants.
  - Add `var _ do.HealthcheckerWithContext = (*CLI)(nil)` conformance guards W4 — Low DONE 2026-11-08: 2026-11-08: var _ do.HealthcheckerWithContext = (*CLI[X])(nil) in cli_healthcheck_ctx_test.go; Scope deliberately exposes the richer WithContext pair instead.
  - Context-cancellation tests for the healthcheck path W5 — Medium DONE 2026-11-08: 2026-11-08: cli_healthcheck_ctx_test.go — live/cancelled/background ctx tests.
  - Commit a samber-linter baseline + wire it into flake/CI W6 — Medium DONE 2026-11-08: 2026-11-08: baseline .samber-linter-baseline.json committed; check-all runs it with ratchet (33% floor).
  - Refresh README/doc.go healthcheck docs W7 — Low DONE 2026-11-08: 2026-11-08: README + doc.go describe HealthCheck(ctx) / do.HealthcheckerWithContext.
  - DONE 2026-10-09: upstream v0.11.0 hit the module proxy (tag pushed); cmdguard `go.mod` bumped to v0.11.0; `AuditMiddleware[T]` + `WithAuditMiddleware[T]` implemented in `auditlog.go` (PhaseBefore/PhaseAfter + duration ms + error, nil-plugin passthrough, FullPath-preferred naming) with 7 tests in `audit_middleware_test.go`; taskctl wires it in `buildApp` and a live `taskctl list` run exported command events (`event_type: command`, before/after, `duration_ms: 1.264`). FEATURES row → FULLY_FUNCTIONAL. F1 — DONE 2026-10-09: upstream v0.11.0 hit the module proxy (tag pushed); cmdguard `go.mod` bumped to v0.11.0; `AuditMiddleware[T]` + `WithAuditMiddleware[T]` implemented in `auditlog.go` (PhaseBefore/PhaseAfter + duration ms + error, nil-plugin passthrough, FullPath-preferred naming) with 7 tests in `audit_middleware_test.go`; taskctl wires it in `buildApp` and a live `taskctl list` run exported command events (`event_type: command`, before/after, `duration_ms: 1.264`). FEATURES row → FULLY_FUNCTIONAL.
  - DONE 2026-10-04: user approved; CHANGELOG entry, annotated tag pushed, proxy-verified (scratch consumer go get + build + run), GitHub release published L1 — DONE 2026-10-04: user approved; CHANGELOG entry, annotated tag pushed, proxy-verified (scratch consumer go get + build + run), GitHub release published
  - DONE 2026-10-04: user approved wire-now; blank import in sdk_imports.go, require v0.1.0 + pinned replace, vendor + vendorHash, wiring test + inventory guards, docs --check green (provider 118) L2 — DONE 2026-10-04: user approved wire-now; blank import in sdk_imports.go, require v0.1.0 + pinned replace, vendor + vendorHash, wiring test + inventory guards, docs --check green (provider 118)
  - DONE 2026-10-09: `crossfile.go` builds a per-directory-group index of package-level CLI vars / Execute errors / constructor errors (test/prod split, blank-ident skip, function-locals never leak); CG003/CG004/CG006 consult it. 6 cross-file tests; FP budget held (cmdguard + examples: 0 findings; sibling scan catches known corpus findings). L3 — DONE 2026-10-09: `crossfile.go` builds a per-directory-group index of package-level CLI vars / Execute errors / constructor errors (test/prod split, blank-ident skip, function-locals never leak); CG003/CG004/CG006 consult it. 6 cross-file tests; FP budget held (cmdguard + examples: 0 findings; sibling scan catches known corpus findings).
  - DONE 2026-10-09: `baseline.go` (Load/Write/ApplyBaseline, versioned schema, per rule+file matching with line-drift tolerance), CLI `--write-baseline` / `--baseline` (default `.cmdguard-lint-baseline.json` when present), stale-entry reporting. 11 unit tests + end-to-end ratchet cycle verified (write→pass→new-finding-fails→stale→tighten). L4 — DONE 2026-10-09: `baseline.go` (Load/Write/ApplyBaseline, versioned schema, per rule+file matching with line-drift tolerance), CLI `--write-baseline` / `--baseline` (default `.cmdguard-lint-baseline.json` when present), stale-entry reporting. 11 unit tests + end-to-end ratchet cycle verified (write→pass→new-finding-fails→stale→tighten).
  - DONE 2026-10-09: the blocker was stale — toolsdk v1.14.0 HAS the options surface. `lint.FilterByIDs` extracted (CLI + provider share semantics); provider declares `enable`/`disable` string options read via `OptionsFromContext`; ValidateOptions contract tested. Needs `lint/v0.2.0` release + BuildFlow `go get` bump to reach the fleet — release on user approval. L5 — DONE 2026-10-09: the blocker was stale — toolsdk v1.14.0 HAS the options surface. `lint.FilterByIDs` extracted (CLI + provider share semantics); provider declares `enable`/`disable` string options read via `OptionsFromContext`; ValidateOptions contract tested. Needs `lint/v0.2.0` release + BuildFlow `go get` bump to reach the fleet — release on user approval.
  - DONE 2026-10-09: `RuleTableMarkdown()` + `rules --markdown` CLI flag + `rules_table_test.go` byte-exact drift guard over the marker-bracketed README section (mutation-verified: edit→FAIL, regenerate→PASS). "Adding a rule" now lists the regenerate step. L7 — DONE 2026-10-09: `RuleTableMarkdown()` + `rules --markdown` CLI flag + `rules_table_test.go` byte-exact drift guard over the marker-bracketed README section (mutation-verified: edit→FAIL, regenerate→PASS). "Adding a rule" now lists the regenerate step.
  - DONE 2026-10-09: all 25 lifetime alerts (incl. the 16 reported 2026-10-04) show state=fixed on GitHub — 22 auto-resolved 2026-10-04 21:07–21:10 UTC by the dependency bumps pushed that day, 3 earlier, 0 dismissed, 0 open (`gh api .../dependabot/alerts` + `pnpm audit` verified). Residual `postcss-selector-parser <7.1.6` (GHSA-rj75-hqrm-r3gf, moderate, via starlight→expressive-code→postcss-nested) fixed by `postcss-nested: ^8.0.1` override in `website/pnpm-workspace.yaml`; `pnpm audit` clean, `astro build` green (22 pages). S1 — DONE 2026-10-09: all 25 lifetime alerts (incl. the 16 reported 2026-10-04) show state=fixed on GitHub — 22 auto-resolved 2026-10-04 21:07–21:10 UTC by the dependency bumps pushed that day, 3 earlier, 0 dismissed, 0 open (`gh api .../dependabot/alerts` + `pnpm audit` verified). Residual `postcss-selector-parser <7.1.6` (GHSA-rj75-hqrm-r3gf, moderate, via starlight→expressive-code→postcss-nested) fixed by `postcss-nested: ^8.0.1` override in `website/pnpm-workspace.yaml`; `pnpm audit` clean, `astro build` green (22 pages).

## [lint/v0.3.0] - 2026-11-08

The cmdguard usage linter ships its v5 baseline: `CurrentMajor` is `"v5"`, so
CG002 flags `cmdguard/v1`–`/v4` imports as stale majors and accepts `/v5`.
Tracks `cmdguard/v5 v5.0.0` (release-shape go.mod). Mirrors
`lint/CHANGELOG.md` `[0.3.0]`.

## [lint/v0.2.1] - 2026-11-08

Patch release of the `lint` sub-module: rules now see through generic
instantiations.

- **Fixed:** `v4.NewCLI[Config](...)`, `v4.NewCommand[C, F](...)`, and
  package-level `var cli, err = v4.NewCLI[T](...)` were invisible to every
  rule (CG001–CG008) and to the cross-file name tracing — the called selector
  is wrapped in type arguments (`IndexExpr`/`IndexListExpr`) and none of the
  four extraction sites unwrapped them. Consumers using cmdguard's
  recommended generic entry points were effectively unlinted. The new
  `selectorExprOf` helper fixes detection for all rules; regression tests
  cover single- and two-type-parameter calls, cross-file constructor-error
  tracing, and the non-generic forms.
- Dogfooding the fix immediately caught one previously invisible CG008
  violation (`examples/config-explain`), which gained `WithCLIVersion`.

---

## [lint/v0.2.0] - 2026-11-08

Second release of the `lint` sub-module: ships the baseline ratchet, cross-file
name tracing, provider options, and the CG007/CG008 end-goal baseline rules to
BuildFlow consumers.

- Baseline/ratchet mode: `--write-baseline` snapshots current findings into
  `.cmdguard-lint-baseline.json`; subsequent runs fail only on findings beyond
  the baseline (per rule+file, line-drift tolerant, stale entries reported so
  the ratchet can tighten). Library surface:
  `LoadBaseline`/`WriteBaseline`/`ApplyBaseline`
- Cross-file name tracing for CG003/CG004/CG006: package-level CLI variables,
  Execute errors, and constructor errors are unioned per directory group (test
  files never merge with prod files; function-local names never leak). FP
  budget held: cmdguard itself and examples stay at 0 findings
- Provider options: `enable` and `disable` (comma-separated rule IDs, same
  semantics as the CLI flags via the new `lint.FilterByIDs`)
- Generated README rule table (`rules --markdown` from `allRuleDefs()`;
  `rules_table_test.go` fails the build when the README drifts)
- CG007 missing short description, CG008 missing CLI version — the end-goal
  baseline rules (warnings, skip `_test.go` files and spread option lists,
  suppressible via `//cmdguard-lint:ignore`)
- Dependency updates: go-finding v1.15.1, toolsdk v1.16.0, go-linter-sdk
  v0.4.0, core pinned to cmdguard v4.2.1

---

## [4.2.1] - 2026-11-04

### Changed

- **Dependency updates across the workspace** (no API changes): koanf v2.4.1
  (parsers json v1.0.1, yaml v1.1.1), go-output v0.38.7,
  samber-do-auditlog v0.12.1, lipgloss v2.0.6, and matching sub-module bumps
  (lint: go-finding v1.15.1, toolsdk v1.16.0, go-linter-sdk v0.4.0).

---

## [4.2.0] - 2026-11-02

**Behavior changes to review before upgrading:**

- **Flag tags are enforced now.** `required`, `values`, and `validate` run at
  execution time against resolved values (flag, env, default). CLIs that
  silently accepted missing required flags or invalid enum/range values will
  start failing with typed errors. Discovery commands (`describe`, `version`,
  `doctor`, `completion`, `help`) are exempt from `required` enforcement, and
  `min`/`max` on string flags now mean length.
- **Validation error text changed.** Failures read
  `Command "add": flag "title": value too short: "ab" must be at least 3 characters.`
  — no more stacked `preparing run context / validating flags for / ...`
  prefixes. Anything matching on the old wording (tests, parsing) must update.
- **Errors teach.** Human-readable errors print one actionable `Hint:` line
  (disable with `WithoutErrorHints()`); machine-format error envelopes carry
  the same sentence as `hint`.

### Added

- **Command-level audit middleware.** `AuditMiddleware[T]` records a
  PhaseBefore event before each command runs and a PhaseAfter event after it
  completes, carrying the wall-clock duration in milliseconds and the
  command's error (nil on success), via samber-do-auditlog v0.11.0's
  `Plugin.RecordCommand`. Command events surface in `Report.Events` and the
  NDJSON stream (`Event.IsCommand()` distinguishes them from service events).
  `WithAuditMiddleware[T](plugin)` wires it into the middleware chain — pair
  it with `WithAuditLog(plugin)` using the same plugin. A nil plugin makes
  the middleware a transparent passthrough.
- **Machine-readable CLI contract.** `cli.Describe()` and
  `cli.WriteDescription(w)` emit a versioned JSON schema (`schema_version`)
  of the whole command tree, including flag type, default, env var (with prefix
  applied), enum values, required/count/hidden/inherited, and the command tree
  with aliases, groups and deprecations. A hidden `describe` command is
  registered automatically at `Execute` time so agents and tooling can drive a
  CLI without scraping `--help`; a consumer-defined `describe` command wins and
  `WithoutDescribeCommand()` opts out. Additive within a schema major.
- **Atomic batch command registration.** `AddCommands[T](cli, cmds...)`
  validates and converts every command before attaching any, accepts
  heterogeneous flags types through the sealed `AddableCommand[T]` interface,
  rejects nil entries, and keys duplicate detection on the cobra command name
  so `"list [path]"` and `"list"` collide. `AddCommand` shares the same
  prepare/attach path and no longer reserves a name when conversion fails.
- **Actionable error hints.** `ErrorHint(err)` maps every sentinel used by the
  JSON error classifier to one fix sentence, and the machine-readable error
  envelope now carries it as `hint`. Human-readable errors also get a `Hint:`
  line (fang's default view and the plain cobra path); `WithoutErrorHints()`
  turns it off. A coverage test fails when a newly classified sentinel lacks a
  hint.
- **Production batteries.** `WithBatteries()` opts into strict command
  validation, signal handling with graceful DI shutdown, error hints, and
  auto-registered `version` (when `WithCLIVersion` is set) and `doctor`
  commands; consumer-defined commands of the same name win. `NewApp` wraps
  construction plus command registration and discards the half-built CLI when
  the build callback fails.

### Fixed

- **`required`, `values`, and `validate` tags are now enforced at run time.**
  `FlagRegistry.ValidateFlags` had tests but no production caller, so the tags
  were parsed and documented but never checked: a missing required flag fell
  through to the handler, and an invalid enum or `min`/`max` value was accepted.
  Validation now runs after parsing for the root and every command registry.
  `min`/`max` on string fields mean length (previously the text was parsed as a
  number, producing "invalid validator parameter"). Environment- and
  default-sourced values are validated too; `describe`, `version`, `doctor`,
  `completion`, and `help` skip required enforcement so discovery always works.
- **Readable validation errors.** Removed the stacked
  "preparing run context / validating flags for / validating flag / required
  flag not set" prefixes; a failure now reads
  `Command "add": flag "title": value too short: "ab" must be at least 3 characters.`

### Changed

- **Dependency bumps:** `samber-do-auditlog` v0.10.0 → v0.11.0 (command
  execution events), `golang.org/x/net` v0.59.0 → v0.60.0 in core + glamour
  (GO-2026-6617/6611/6603, govulncheck-verified reachable).
- **website:** `postcss-nested` overridden to ^8.0.1 in
  `website/pnpm-workspace.yaml`, pulling patched `postcss-selector-parser`
  ≥7.1.6 (GHSA-rj75-hqrm-r3gf, moderate). `pnpm audit` clean; `astro build`
  green.
- **website:** dependency refresh clearing the 2026-11-02 Dependabot batch —
  astro 7.4.4, vite 8.3.5, typescript 7.0.2, html-validate 11.18.0, plus
  transitive `undici` ≥8.12.2, `zod` ≥4.6.6, `fast-uri` ≥3.1.9. Full
  `pnpm audit` clean (19 findings, 6 high → 0); site builds 22 pages.

---

## [4.1.0] - 2026-10-04

### Fixed

- **Nil `TypeHandlerFunc.RegisterFunc` no longer silently drops the flag.** A handler
  that registers nothing made the flag vanish from the CLI, surfacing far from the
  cause as "unknown flag" at invocation time (found dogfooding in nsfw-classifier).
  `dispatchRegister` now rejects such handlers at CLI-build time with an error naming
  the type and flag. `Parse`/`Default` with nil funcs return errors/nil instead of
  panicking. The low-level `Register` stays nil-tolerant for derived handlers.
- **Repeated `[]string` flags no longer parse pflag's bracketed render form.** A
  changed slice flag is read through `flag.Value.String()`, which renders `"[a,b]"`;
  splitting that string on commas handed downstream handlers elements like `"[100"`
  and `"101]"` instead of `"100"` and `"101"` (seen live with a consumer's repeatable
  `--document` ID flag). The slice parse now trims the brackets before splitting,
  with an end-to-end regression test through `ExecuteWithArgs`.

### Changed

- **stdlib `time.Duration` is now handled by default.** A fresh registry (and the
  global template) includes the duration handler previously available only via the
  opt-in `RegisterGoDurationHandler()`: native pflag duration registration, `5s`-style
  parsing, duration defaults. Explicit `RegisterTypeHandler` overrides still win.

### CI

- Workflows install Go 1.27 (matching the repo's `go 1.27` floor) and pin
  golangci-lint to v2.14.0, the version whose `config verify` passes the repo config.

---

## [lint/v0.1.0] - 2026-10-04

First stable release of the `lint` sub-module: a usage linter for cmdguard-based CLIs.

- 6 syntactic rules (CG001–CG006): execute bypass (`fang.Execute` around cmdguard),
  stale major import, constructor panic (`MustNew*`), runtime `SetVersion`,
  duplicate version options, execute-error reprint
- Pure `go/parser` AST walk — no type info, no build required
- BuildFlow provider (`lint/provider`, toolsdk spec `cmdguard-lint`, module fan-out)
- Dogfood CLI `cmd/cmdguard-lint` (built with cmdguard v4 itself) with text, JSON,
  and SARIF 2.1.0 output (round-trips via `finding.FindingsFromSARIF`)
- Suppressions: `//cmdguard-lint:ignore <ID> <reason>` (reason required) on the
  finding line or the line above
- Per-directory memoized analysis with `ClearCache()` for long-lived processes
- `CurrentMajor` constant tracks the latest cmdguard major for CG002
- Validated: timesheets corpus 18/18 findings at hand-audited locations, 0 false
  positives on cmdguard/erraudit/go-structure-linter/branching-flow
- 34 test functions; golangci-lint clean under the repo config

---

## [4.0.2] - 2026-08-06

### Changed

- Dependency alignment release: sub-module and core dependency bumps to align
  with the v4.0.1 sub-module releases (no API changes).

---

## [glamour/v0.2.0] - 2026-08-06

### Changed

- **BREAKING:** Migrated dependency from `cmdguard/v3` to `cmdguard/v4`. All public function return types changed from `v3.CLIOption` to `v4.CLIOption`. Update import aliases from `v3` to `v4`.

---

## [prompts/v0.2.0] - 2026-08-06

### Changed

- **BREAKING:** Migrated dependency from `cmdguard/v3` to `cmdguard/v4`. All public function return types changed from `v3.CLIOption` to `v4.CLIOption`. Update import aliases from `v3` to `v4`.

---

## [spinner/v0.2.0] - 2026-08-06

### Changed

- **BREAKING:** Migrated dependency from `cmdguard/v3` to `cmdguard/v4`. All public function return types changed from `v3.CLIOption` / `v3.Middleware[T]` to `v4.CLIOption` / `v4.Middleware[T]`. Update import aliases from `v3` to `v4`.

### Fixed

- Added `sync.Once` and `sync.Mutex` to `textSpinner` for concurrency-safe stop behavior.

---

## [telemetry/v0.2.0] - 2026-08-06

### Changed

- **BREAKING:** Migrated dependency from `cmdguard/v3` to `cmdguard/v4`. All public function return types changed from `v3.CLIOption` / `v3.Middleware[T]` to `v4.CLIOption` / `v4.Middleware[T]`. Update import aliases from `v3` to `v4`.

---


## [4.0.1] - 2026-08-06

### Fixed

- `NO_COLOR` env var restore — `applyNoColorIfSet` now uses `os.LookupEnv` instead of `os.Getenv` to correctly distinguish "unset" from "set to empty string" when restoring the environment after `--no-color` flag execution. Previously, a pre-existing `NO_COLOR=""` would be unset instead of restored.
- `doc.go` godoc example — fixed stale `v3` import alias and "v2 constructors" reference.
- README.md sub-modules code example — fixed missing `"time"` import (used `time.Millisecond` without importing the `time` package).

### Changed

- `go-output` upgraded v0.35.0 → v0.37.0.
- `koanf` patched: v2.3.5 → v2.3.6, parsers/json v1.0.0 → v1.0.1, parsers/yaml v1.1.0 → v1.1.1.
- **Documentation v3→v4 drift fix** — 36 user-facing and contributor-facing documents updated from stale v3 API references to v4 (QUICKSTART, TUTORIAL, WHAT_THIS_PROJECT_IS_ABOUT, WHAT_THIS_PROJECT_IS_NOT, MIGRATION_FROM_COBRA, COMPARISON, PERFORMANCE, doc.go, README, website guides). 23 website source files updated.
- AGENTS.md project structure corrected from stale v3 references to v4; flightrecorder sub-module documented (package table, lint strategy, design principles, gotchas).
- Nix flake inputs updated; go.mod `replace` directives consolidated into a single block.
- `TODO_LIST.md` and `ROADMAP.md` created.

---

## [flightrecorder/v0.1.0] - 2026-08-06

First stable release of the `flightrecorder` sub-module. Tagged alongside v4.0.0 core.

- Zero external dependencies (Go 1.25+ stdlib `runtime/trace`)
- Continuous in-memory trace buffering with configurable capacity
- Context-aware snapshot capture: `Recorder.Capture`, `Recorder.CaptureToWriter`
- Automatic slow-command (`CaptureOnSlow`+`SlowThreshold`) and error (`CaptureOnError`) capture middleware
- `WithFlightRecorder[T]` and `WithFlightRecorderRecorder[T]` CLI options
- Process-wide singleton with `sync.WaitGroup` coordination
- 48 test functions, 3 benchmarks, 1 fuzz target, 96.1% coverage

---

## [4.0.0] - 2026-07-28

### Breaking Changes

- **Module path:** `github.com/larsartmann/cmdguard/v3` → `github.com/larsartmann/cmdguard/v4`
- **Package directory:** `pkg/cmdguard/v3/` → `pkg/cmdguard/v4/`
- **Package name:** `v3` → `v4` (update all import aliases)
- **Config loading:** `configload` sub-package deleted; `WithConfigFile(paths...)` now auto-detects JSON/YAML/TOML via `KoanfLoader`
- `NewJSONLoader` deleted; `configload` sub-package deleted

### Changed

- Consolidated config loading into unified `KoanfLoader` path resolution
- go-output upgraded v0.31.1 → v0.35.0
- samber-do-auditlog upgraded v0.7.0 → v0.8.1

---

## [3.1.0] - 2026-07-25

### Fixed

- Auditlog diagram export build break

---

## [3.0.0] - 2026-07-07

### Breaking Changes — Breaking API Redesign

Corrects the v2.11.0 mis-release that put breaking changes on a `/v2` path.

- **Non-generic `CLIOption` / `CommandOption`** — type inference via positional flags eliminates the v2 "7 type params per command" explosion. Metadata options (`WithShort`, `WithLong`, `WithExample`, etc.) take zero type parameters. Only lifecycle hooks (`WithPreRunE`, `WithPostRunE`, `WithSubcommands`) remain generic.
- **Module path:** `github.com/larsartmann/cmdguard/v2` → `github.com/larsartmann/cmdguard/v3`
- `NewCommand` / `NewParentCommand` API refactored — flags passed positionally, no `WithFlags` option needed
- Sealed lifecycle hook interfaces — compile-time type safety restored

### Added

- 5 extracted optional sub-modules (`glamour`, `manpage`, `prompts`, `spinner`, `telemetry`) — core has zero dependencies on these
- `go.work` multi-module workspace for unified local builds
- Audit logging with `samber-do-auditlog` integration

---

## [2.10.4] - 2026-07-07

### Fixed

- Retracts mis-released v2.11.0 (breaking changes were incorrectly published on a `/v2` path; corrected in v3.0.0)

---

## [2.10.3] - 2026-07-07

### Changed

- go-output v0.30.1, samber-do-auditlog v0.4.0 (`d006757`). Entry added 2026-11-07 — the tag existed without a CHANGELOG entry (found by the docs-health sweep).

---

## [2.10.0] - 2026-06-28

### Added — Cobra-Correctness Contract + Escape-Hatch APIs

Refocus on the founding mission: "make consumers use Cobra correctly."

- `SilenceUsage = true` by default (cobra's #1 footgun, off by default)
- `ExitCode(err) int` — public exit-code mapping function
- Scoped flags (`local:"true"`) — root-only flags not inherited by subcommands
- `hidden:"true"` flag tag — exclude from `--help`, stay functional
- `ConfigFromContext[T]` — type-safe config retrieval for raw cobra subcommands
- `WithPostFlagParse[T]` — post-parse hook (DI init, session storage)
- `WithCleanup[T]` — post-RunE cleanup that fires even when RunE errors

### Changed

- Go directive bumped 1.26.3 → 1.26.4 (CVEs GO-2026-5037/5038/5039)
- 457 tests, 1430 runs, 26 benchmarks, 7 fuzz targets, 86.7% coverage

---

## [2.9.0] - 2026-06-22

### Added

- 4 new audit-log export formats (d2, plantuml, tree, htmltree) — total now 11
- samber-do-auditlog v0.1.0 → v0.3.0
- go-output v0.17.1 → v0.17.2

### Changed

- 430 test functions, 26 benchmarks, 7 fuzz targets, 86.6% coverage, 0 lint issues

---

## [2.8.0] - 2026-06-19

### Changed

- Release v2.8.0 — incremental improvements and dependency updates

---

## [2.7.0] - 2026-06-17

### Changed

- Release v2.7.0 — incremental improvements

---

## [2.6.0] - 2026-06-12

### Added

- `RegisteredFormats()` — dynamic format discovery from registered marshalers
- Shape-aware error messages in `OutputResult()`
- `OutputTable()` uses `AddRowChecked()` for fail-fast row validation
- Dynamic `--output` flag help from `RegisteredTableDataFormats()`
- go-output v0.9.0 with generic registries

### Removed

- `IsExecutable()` — use `HasHandler()`
- 16 `Format*` constant re-exports — use `output.Format*` directly
- `ParseOutputFormat()` — use `output.ParseFormat()`
- `SupportedFormats()` — use `output.AllFormats`
- `IsFormatSupported()` — use `format.IsValid()`
- `ErrNoFlags`, `ErrTooFewArgs`, `ErrTooManyArgs` sentinels

### Changed

- 407+ tests passing, 85.9% coverage, 0 lint issues, 0 race conditions

---

## [2.5.0] - 2026-06-05

### Fixed

- Module path fixed for Go major version compatibility

---

## [1.0.0] - 2026-04-30

### Changed

- Stable release — dependency updates and go.sum alignment

---

## [0.2.0] - 2026-04-08

### Changed

- Minor update

---

## [0.1.0] - 2026-02-14

### Added — Initial Release

First stable release of cmdguard with the Guard API.

- Single-step initialization with `cmdguard.New()`
- Compile-time validation (panic on invalid commands)
- Built-in version and validate commands
- Environment-based configuration
- Strict mode for RunE enforcement
- 94% coverage on config package, 66% coverage on cmdguard package
- 3 direct dependencies (minimal external footprint)
