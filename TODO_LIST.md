# TODO List

> Short- and mid-term improvement tasks — actionable, bounded, with status.
> Derived from FEATURES.md partially-functional items, TODO(v5) markers,
> status report harvests, and code analysis. Updated as work completes.

---

## v5 Breaking Changes (Deferred)

These items have `TODO(v5)` markers in source code. They are public API
renames/removals from the 2026-07-18 naming review (plus follow-ups) that
cannot ship in v4.x without breaking downstream consumers.

<!-- check-rows:ledger -->

| # | Task | File | Status |
| - | ---- | ---- | ------ |

Further v5 candidates (Get[T] rename, Package[T] redesign) live in
ROADMAP.md §"Additional v5 candidates" — keep them there, not here.

> These must stay as `TODO` (not `NOTE`) so `grep TODO` finds them. The godox
> linter exclusion for `TODO(v5)` is deliberately narrow.

---

## Pareto Plan Tasks (2026-11-02)

Full ranked plan with fine-grained breakdown (123 tasks) and execution graph:
[docs/planning/2026-11-02_00-38_pareto-v5-release-execution-plan.md](docs/planning/2026-11-02_00-38_pareto-v5-release-execution-plan.md).
Only P0/P1 near-term items duplicated here; P2 lives in the plan.

<!-- check-rows:ledger -->

| #       | Task                                                | Notes                                                                                                                                                                                                                                                              | Priority | Status  |
| ------- | --------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | -------- | ------- |
| ~~P6~~  | ~~lint/v0.2.0 + BuildFlow bump + fleet baseline~~   | ~~CG007/CG008, baseline, options reach the fleet.~~ DONE 2026-11-08: lint/v0.2.0 + v0.2.1 + v0.3.0 all tagged and proxy-served; BuildFlow bumped to lint v0.3.0 (input pinned lint/v0.3.0 ref); CG002 CurrentMajor=v5 baseline shipping fleet-wide.                | P1       | ✅ DONE |
| ~~P10~~ | ~~Machine error codes + generated ERROR_REFERENCE~~ | ~~Stable `code` per sentinel; docs cannot drift.~~ DONE (shipped in v4.1.0/v4.2.x era — ErrorCode(err), envelope error_code, drift-guarded ERROR_REFERENCE regen; verified in the v5 tree 2026-11-08).                                                             | P1       | ✅ DONE |
| ~~P11~~ | ~~Validation-semantics hardening (M06–M08)~~        | ~~Env/config e2e, registration-time defaults, Duration bounds.~~ DONE (shipped in v4.2.0: ValidateFlags wired for env/default values, ErrInvalidFlagDefault at construction, min/max domain resolution incl. Duration bounds; verified in the v5 tree 2026-11-08). | P1       | ✅ DONE |

---

## v5 End-Goal Workstreams (pre-v5 slice)

Plan: [docs/planning/2026-11-01_v5-end-goal-paved-road.md](docs/planning/2026-11-01_v5-end-goal-paved-road.md).
These are the additive, non-breaking pieces shipped before the major so
consumers can adopt early and the v5 break stays small.

<!-- check-rows:ledger -->

| # | Task | Notes | Status |
| - | ---- | ----- | ------ |

---

## Config-Truth Follow-ups (2026-11-02)

From the config-truth execution (plan:
[docs/planning/2026-11-02_05-36_config-truth-pareto-execution-plan.md](docs/planning/2026-11-02_05-36_config-truth-pareto-execution-plan.md);
status report: `docs/status/2026-11-02_10-15_config-truth-execution-mid-flight.md`).

<!-- check-rows:ledger -->

| #   | Task                                                    | Notes                                                                                                                                                                                                                                                                                                                         | Priority | Status                                                                                                                                                    |
| --- | ------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
| CT1 | cmdguard-lint CG009 idea: example configs stay truthful | A rule that parses a repo's committed example configs (e.g. `config.example.yaml`) and rejects keys the CLI's config vocabulary does not accept. Needs cross-file type info (or a mode that shells `describe`/`ConfigKeys` from a built binary), so it is an IDEA, not a task — scope it in the lint backlog before building. | Low      | 💡 IDEA                                                                                                                                                   |
| CT2 | Decide the `schema` sub-module release                  | Ship as `schema/v0.1.0` alongside v4.3.0, or keep workspace-local until the API settles. Tag shape per the release playbook (real pin + version-scoped replace); README marks the surface experimental.                                                                                                                       | Medium   | 📝 DECIDED 2026-11-08: hold workspace-local until v5.0.0 ships; tag schema/v0.1.0 after, per the lint/v0.2.x playbook (real pin + version-scoped replace) |
| CT3 | Wire `examples/config-explain` into README + website    | The demo exists and runs; it is not linked from README.md or the website docs yet.                                                                                                                                                                                                                                            | Low      | 🔴 TODO                                                                                                                                                   |

---

## Harvested (2026-11-07 docs-health pass)

Items pulled forward from the 2026-06→2026-09 status reports during the
annotate/archive sweep. Everything else those reports contained was either
resolved (struck with evidence in the report itself) or routed to ROADMAP
"Harvested from 2026-0X status reports".

<!-- check-rows:ledger -->

| #      | Task                                                    | Notes                                                                                                                                                                                                                                                                                                                                                                             | Priority | Status  |
| ------ | ------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------- |
| H4     | Website hero-code sync guard                            | `website/src/data/hero-code.ts` (raw string) vs the inline highlighted copy in `HeroSection.astro` — no test or derive step keeps them in sync                                                                                                                                                                                                                                    | Low      | 🔴 TODO |
| H6     | Audit `docs/COMPARISON.md` prose for v4 accuracy        | Header still says "Last Updated: 2026-07-07" while claiming v4-era content; the full prose audit was left open by three 2026-08-06 reports                                                                                                                                                                                                                                        | Low      | 🔴 TODO |
| ~~H9~~ | ~~Fixture test for lint resolveImports dot/blank arms~~ | ~~`lint/walk.go:240,242` arms have 0 coverage~~ DONE 2026-11-08: `lint/crossfile_imports_test.go` covers both arms; uncovered a REAL gap — dot-imported bare constructors (`NewCLI(...)` without a selector) are untraced, pinned as current behavior and tracked as follow-up H9b (trace bare constructor calls when the file dot-imports cmdguard; non-idiomatic pattern, Low). | Low      | ✅ DONE |
| H9b    | Trace dot-imported bare constructors in cmdguard-lint   | `isSelectorCall` requires a selector, so `NewCLI(...)` under a dot import is invisible to CG003/CG004/CG006; pinned by `lint/crossfile_imports_test.go` (flips to expect findings when fixed). Non-idiomatic pattern; Low priority.                                                                                                                                               | Low      | 🔴 TODO |

---

## samber-linter HW-2 Follow-ups (2026-09-20 report)

The context-aware healthcheck FIX shipped, but its follow-up cluster is
genuinely open (verified 2026-11-07 by the docs-health sweep).

| # | Task | Notes | Priority | Status |
| - | ---- | ----- | -------- | ------ |

---

## Blocked on Upstream

<!-- check-rows:ledger -->

| # | Task | Blocker | Priority | Status |
| - | ---- | ------- | -------- | ------ |

---

## Lint Sub-Module (`lint/`)

Harvested from `docs/status/2026-10-04_07-11_cmdguard-lint-sub-module-build.md` §f
(full 50-item backlog lives there; these are the near-term actionable ones).

<!-- check-rows:ledger -->

| #  | Task                                                         | Notes                                                                                                                                                                                                                                                                                                                                                                                                                | Priority | Status      |
| -- | ------------------------------------------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ----------- |
| L6 | golangci-lint module-plugin distribution layer               | Scoped 2026-10-09: requires an `analysis.Analyzer` bridge (golangci/plugin-module-register, per go-humanize-linter pattern: `lint/plugin/` package + `.custom-gcl.yml` + `.golangci.custom.yml`). One pass = one package's files, so cross-file tracing carries over; checks need a pass.Files→project adapter. NOT a thin wiring task — do not half-ship. Also needs the custom-gcl build verified against corpora. | Low      | 🔴 TODO     |
| L8 | Bump `CurrentMajor` in `lint/walk.go` when cmdguard v5 ships | Standing duty; CG002 tracks the constant. Rule IDs CG001–CG008 are load-bearing (suppressions key on them) — never renumber                                                                                                                                                                                                                                                                                          | —        | 🔵 STANDING |

---

## Security

<!-- check-rows:ledger -->

| #  | Task                                        | Notes                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 | Priority | Status     |
| -- | ------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ---------- |
| S2 | Bump toolchain to go1.27.2 (6 stdlib vulns) | govulncheck (2026-10-09): GO-2026-6599/6600/6611/6613/6617/6603 reachable via stdlib `net/http`+`html/template` @go1.27.1; all fixed in go1.27.2. Module-side `x/net` already bumped to v0.60.0 (clears the x/net half of 6611/6617/6603). **Blocked:** nixpkgs `go_1_27` is still 1.27.1 on unstable — setting `toolchain go1.27.2` now breaks the sandboxed treefmt/nix checks (`GOTOOLCHAIN=local`, see AGENTS.md gotcha). Re-check `nix eval nixpkgs#go_1_27.version` periodically; when ≥1.27.2: set toolchain directive, run `nix flake check`. | Medium   | 🟡 BLOCKED |
