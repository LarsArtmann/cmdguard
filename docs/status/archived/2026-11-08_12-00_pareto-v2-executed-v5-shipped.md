# Status Report — Pareto Plan v2 Executed: v5.0.0 Shipped, Delivery Loop Closed

> **Live remainder:** none — every item is struck with evidence below. The owner-gated queue resolved in the 2026-11-11 v4-retirement/master-merge session; standing candidates live in ROADMAP.md (v6 candidates, Open questions).

**Date:** 2026-11-08 (second session of the day, executing the 08:42 plan)
**Branch:** `v5-scaffold` (pushed; v5.0.0 tagged on it)
**Verdict:** R1–R9 of the Pareto plan EXECUTED AND VERIFIED; Tier D pulled 4 slices; 2 fleet PRs open.

---

## What shipped

### R1 — v5.0.0 is LIVE (the 1%)

- **Tag `v5.0.0`** pushed and proxy-served (commit 56fd69b). Coordinated
  sub-module tags on the same commit: `flightrecorder/v0.2.0`,
  `glamour/v0.3.0`, `spinner/v0.3.0`, `prompts/v0.3.0`, `telemetry/v0.3.0`,
  `schema/v0.1.0` (first tag). All seven serve from `proxy.golang.org`.
- **Scratch consumer verified end-to-end from the tag**: build, `describe`
  contract, batteries ON (doctor present), strict default (missing `WithShort`
  errors at `Build`), `WithoutBatteries` opt-out (doctor absent).
- **GH release** `v5.0.0` published (migration guide as notes, `--latest`).
- `check-all` END-TO-END GREEN before tagging (build + race tests ×8 modules
  - lint 0 ×8 + samber ratchet + format + tidy ×8); coverage re-measured and
    pinned: **89.3% core / 97.6% flightrecorder / 93.9% taskctl**; 468 core
    test functions, 26 benchmarks, 7 fuzz targets.
- AGENTS.md rewritten for v5 (21 v4 refs → 5 legitimate historical ones),
  SKILL.md rewritten (Builder, renames, flips, re-anchored source lines —
  every line number re-verified against the v5 tree), CHANGELOG cut
  (`[5.0.0]` section; the orphaned v4.x notes parked under a labeled
  `[Backfill]` pending owner §g2).

### R2 — the delivery loop closed

- **`lint/v0.3.0`** tagged, proxy-served, GH release: `CurrentMajor="v5"`, so
  CG002 flags v1–v4 imports fleet-wide. Dogfooded on the v5 tree (0 findings)
  and erraudit (12 findings) before tagging.
- **BuildFlow migrated to v5** — the bump turned out to require migrating
  BuildFlow's OWN CLI (it is a cmdguard v4 consumer): codemod applied, root
  go.mod → v5.0.0 with coordinated sibling pins, flake input pinned to the
  `lint/v0.3.0` ref, vendorHash cycle done. internal/cli + gomod-checker +
  tools/providers suites green; `docs --check` no drift; nix binary builds
  and runs. Committed locally (push remains owner-gated per §g2).
- **Fleet smoke recorded**: timesheets 15×CG002+1×CG003 → migrated;
  erraudit 12×CG002 → migrated; go-structure-linter 0; BuildFlow 0
  (self-consistent after its own migration).
- **erraudit PR #15** and **timesheets PR #3** opened — both mechanical,
  both with evidence blocks. timesheets' `internal/parser` WIP had resolved,
  so its planned hold was lifted; full suite green on both branches (erraudit's
  one BDD failure is pre-existing on master — worktree-verified, documented
  in the PR).

### R3–R9 — the rest of the 20%

- **Truth-in-docs**: README stability note/coverage/batteries rows/RegisterTypeCodec/Builder
  example; website quick-start + installation rewritten to v5 (they were
  teaching v4 with pointer-flag samples that never compiled);
  `docs/QUICKSTART.md` gained the Builder section.
- **NEW GUARDS, both wired into `check-all`**:
  - `scripts/check-guide-sync.sh` — docs/website current-major consistency
    (module path, batteries story, Builder presence, hero-code.ts). It caught
    a real QUICKSTART gap on its first run.
  - `scripts/check-dep-tables.sh` — AGENTS.md dependency tables derived-vs-claimed.
- **Correctness (R4)**: `WithStrictPositionalArgs()`/`WithArgsPassthrough()`
  shipped as mechanism (the default flip itself stays owner-gated per
  guardrail #1 and is pinned by `TestDefaultLeafCommand_StaysPassthrough`);
  Builder made concurrency-safe (mutex) + race tests (same-name batch →
  duplicate detection, unique-name batch → all attach); sequential
  `ExecuteWithArgs` × cancelled runs race-clean; **concurrent Execute on ONE
  CLI documented as out of contract** (single-flight reset by design).
- **Codemod round-trip proof (R4.6/7)**: `docs/migration_v4_v5/roundtrip.sh`
  - `_fixture` — a v4.2.1 consumer builds, migrates, compiles, runs
    (batteries ON), lints CG002-clean against the published v5.0.0.
- **Hardening (R5)**: koanf empty/missing → `ErrConfigFileNotFound` with
  empty `LoadedPath` (contract named); `WithoutBatteries` describe tree
  pinned reduced; doctor zero-services all-green shape; flightrecorder
  error-snapshot survives Execute's cancelled-context return (real `.trace`
  file asserted, `Stop()` awaited).
- **Security (R7)**: `toolchain go1.27.3` pinned; **govulncheck under 1.27.3
  reports ZERO vulnerabilities** (was 6 stdlib vulns on 1.27.0).
- **Truth pass (R8)**: P6/P10/P11 struck with shipped-where evidence; H9
  closed with dot/blank-import fixtures — which **found a real gap**:
  dot-imported bare constructors are untraced by cmdguard-lint (pinned as
  current behavior, follow-up H9b opened). The stale 2026-10-02 report
  annotated with its v5 resolution.
- **Guards/misc (R9 + Tier D)**: hero-code.ts was still entirely v4 —
  rewritten and guarded; 11 missing sub-module release pages created;
  SECURITY.md supported-versions updated (v5 supported, v4 frozen);
  `docs/MIGRATION_RUNBOOK.md` written from today's three real migrations;
  COMPARISON.md first audit pass; config-explain discoverability (README +
  website); samber-linter hw-5 suppressions with reasons on the two
  framework-internal ProvideValue sites (coverage counter remains 1/3 by
  tool design — suppressions gate findings, not the ratchet count);
  HTFC verdicts recorded in go-composable-business-types' TODO_LIST.
- Sub-modules reset to the documented dev shape (unversioned replace +
  zero-pseudo pin) post-tag; GOWORK=off builds verified ×7.

## Final gates (all green)

- `nix run .#check-all` END-TO-END (now 7 gate groups incl. the two new doc guards)
- full `-race` battery ×8 modules (inside check-all)
- `docs-gates run .` 5/5 PASS (audit-all + annotate reds fixed: two 2026-11-02/
  03:56 items annotated as shipped in v5.0.0; two plan-file spacing findings
  repaired with the `go.mod` files pattern)
- `nix flake check` green; website `pnpm build` green (22 pages)

## Owner-gated queue (D44 / §g)

1. ~~**BuildFlow push**: local commits with the v5 migration (through~~ done — pushed with the v5-scaffold→master fast-forward (2026-11-11, 359feb3 on origin/master)
   ~~789625788) are NOT pushed — awaiting the §g2 answer.~~
2. ~~**Branch policy**: v5-scaffold carries the v5.0.0 tag; whether master~~ done — master fast-forwarded to v5-scaffold (359feb3); residual branch-fate question moved to ROADMAP Open questions
   ~~fast-forwards to it (or v5 development continues on this branch) is the~~
   ~~§g1 answer.~~
3. ~~**Strict positional-args DEFAULT flip**: mechanism shipped~~ done — shipped as opt-in WithStrictPositionalArgs/WithArgsPassthrough (command_options.go); the hard default flip remains a ROADMAP idea
   ~~(`WithStrictPositionalArgs`), flip is one line + one matrix row after an~~
   ~~owner decision (guardrail #1).~~
4. **CHANGELOG [Backfill]**: the orphaned v4.x notes section awaits the §g2 **→ CHANGELOG [Backfill] section — the owner §g2 decision is still pending and tracked there**
   decision on folding into historical sections.

## Verschlimmbesserung guardrails honored

No default flips beyond the shipped v5 set (strict-args shipped as opt-in
mechanism); publish shapes followed the lint/v0.2.1 + v4.2.0 precedents
exactly; no fleet in-repo migration before the proxy served v5.0.0; codemod
stayed additive-only; no published CHANGELOG section rewritten.
