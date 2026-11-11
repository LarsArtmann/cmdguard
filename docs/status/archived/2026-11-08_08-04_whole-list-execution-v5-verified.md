# Whole-List Execution — v5 Scaffold to Verified Green (2026-11-08 08:04)

> **Live remainder:** none — every item is struck with evidence below. §d residue and follow-on work live in TODO_LIST.md; open API-review questions live in ROADMAP.md.

> **What this is.** Closing report for the "DO THE WHOLE TODO LIST" run
> resuming from the 03:56 Tier A report. Tier A was already shipped (lint
> v0.2.0/v0.2.1 + BuildFlow + W1). This session completed Tier B on the
> `v5-scaffold` branch, Tier C, the bounded Tier D slices, and the full
> verification battery.

## a. What was done

**Red-test repair + battery (B1 close-out)**

- Fixed the 3 known reds: ERROR_REFERENCE regenerated (v5/TypeCodec drift),
  lint `walk_test` v5 expectations, `TestHuhRunner_*` → `TestHuhPrompter_*`.
- Full battery green: root + 7 sub-modules build/test, `-race`, lint 0 issues
  (root + 7), GOWORK=off per-module builds, `nix flake check`.

**B1.4–B1.6 + B3 — the default flips (M15)**

- v5 defaults: strict validation, signal handling + graceful DI shutdown,
  strict config keys, batteries bundle ON — each with a `Without*` opt-out
  (`WithoutStrictValidation`, `WithoutSignalHandling`, `WithoutGracefulShutdown`,
  `WithoutStrictConfigKeys`, `WithoutErrorHints`, `WithoutBatteries`;
  `WithBatteries()` restores the bundle).
- `cli_defaults_matrix_test.go` pins every default and every opt-out.
- Fallout (~35 `WithShort` call sites + helpers) fixed across 12 test files;
  describe golden regenerated (doctor command now in the contract).

**Two REAL bugs surfaced by the flips, fixed at root cause**

1. ~~**Stale cobra subcommand context across executions.** `SetContext` in~~ done — fixed in v5.0.0: execCtx + PersistentPreRunE derivation (CHANGELOG [5.0.0] Repeated-Execute invariants)
   ~~PersistentPreRunE shadows the root ctx FOREVER: run 2+ inherited run 1's~~
   ~~cancelled context (`DI: health check timeout: context canceled` on every~~
   ~~command after the first doctor run). Fixed with `CLI.execCtx` (current~~
   ~~run's context, used as the derive base) + per-Execute `flowCtx` rebuild.~~
   ~~New contract test proves shutdown fires ONLY on the signal-cancelled path~~
   ~~and the scope stays alive across `ExecuteWithArgs`.~~
2. ~~**Graceful shutdown fired on EVERY execution** (v4 semantics: deferred~~ done — fixed in v5.0.0: auto-shutdown fires only when the execution context was cancelled (CHANGELOG [5.0.0] Graceful-shutdown semantics)
   ~~`scope.Shutdown` per Execute — fine one-shot, fatal for repeated runs~~
   ~~under defaults). Redesigned: shutdown fires only when the execution~~
   ~~context was cancelled (the documented SIGINT/SIGTERM semantics). Documented~~
   ~~in AGENTS.md (shutdown semantics + repeated-Execute gotchas + tracked~~
   ~~`execCtx //nolint:containedctx`).~~

**B1.3 — Builder (M12)**

- `Builder[T]` + `Add`/`AddParent`/`Try`/`Check`/`Build`: one line per
  command, every construction AND prepare error collected and joined, atomic
  attach, duplicate detection. `examples/taskctl` migrated (359 lines, the
  14-stanza error ladder is gone), same tests green.
- `Add`/`AddParent` are free functions (Go methods cannot carry F).

**B1.11 + C4 — lint config truth**

- Dead `TODO(v5)` godox exclusion removed (all 5 markers executed).
- All dead `pkg/cmdguard/v2/*` exclusions removed; v4 paths → v5
  (`type_codec.go`, flags_validate, cli_command, example_test); ireturn
  allow-list `TypeHandler` → `TypeCodec`.

**B2 — the break is payable (M13)**

- `docs/migration_v4_v5/codemod.sh`: package-path + module-path + alias +
  rename rewrites (`CommandInfo`, `TypeCodec` family, `HuhPrompter`); prints
  manual follow-ups (`SetConfig` removed, generic `RegisterInScope`,
  removed `CLI.HealthCheckWithContext`).
- `docs/MIGRATION_v4_v5.md`: flip table, semantic notes (signal-path
  shutdown, doctor counts), worked example.
- Dry-runs: erraudit codemodded → FULL suite green against local v5
  (end-to-end consumer proof = C1); timesheets committed HEAD → cmd packages
  green (its `internal/parser` break is pre-existing owner WIP, not codemod
  damage).

**Tier C**

- C2: QUICKSTART + TUTORIAL rewritten to v5 (paths, defaults, batteries,
  describe contract, hints, runtime tag enforcement sections).
- C3: `koanf_edge_test.go` — TOML datetime round-trip (RFC 3339 text lands in
  string fields), YAML anchors + top-level merge (nested keys never leak),
  int-vs-float64 fidelity across YAML/TOML, empty files.
- **C3 found a real bug**: a 0-byte `.json` config failed startup
  ("unexpected end of JSON input") while YAML/TOML accepted empty. Fixed:
  empty files are skipped like missing ones (uniform "empty = absent").
- C5/C6: taskctl + examples are root-module packages — check-all and CI
  `go test ./...` cover them (verified, no wiring needed).
- C7: `.samber-linter-baseline.json` committed (coverage floor 33%); check-all
  runs the ratchet (skips gracefully when the binary is absent).
- C8/C10: `CLI.HealthCheckWithContext` removed (HealthCheck IS ctx-aware);
  `Scope.HealthCheck`/`HealthCheckResults` delegate to their WithContext
  variants; `var _ do.HealthcheckerWithContext = (*CLI[X])(nil)` guard; live/
  cancelled/background-context healthcheck tests.
- C9: ERROR_REFERENCE drift test runs inside CI's `go test ./...` (verified).
- C1: consumer proof delivered via the B2 dry-runs (stronger than asked).

**Tier D (bounded slices)**

- D18: README + doc.go describe `HealthCheck(ctx)` /
  `do.HealthcheckerWithContext` + signal-path shutdown.
- D19: PromptConfirm tests already assert real behavior (delegation + error
  wrapping) — H2 closed.
- D22: GitHub release pages created for v4.0.0 + v4.0.1 from CHANGELOG.
- D26/S2: nixpkgs `go_1_27` = 1.27.3 ≥ 1.27.2 — already satisfied.
- D33: gopls check runs clean on go-composable-business-types (no
  infertypeargs findings) — flag settled.
- D34: git-town.toml main="master", connector=api; library-policy.yaml does
  not exist in that repo — flag settled.
- D38: architecture.svg re-rendered 2026-11-08; examples/README documents the
  test loop — pass.
- D39: `2026-11-02_00-11_two-day-session-summary.html` (+ its .md sibling)
  archived.
- Docs purge ritual executed: 43 struck rows removed from TODO_LIST,
  substance preserved in CHANGELOG `[Unreleased] → Housekeeping`.

## b. Verification (final battery, all green)

- `go build ./...` root + 7 sub-modules, GOWORK=off per-module builds
- `go test ./... -count=1 -race` root + 7 sub-modules
- `golangci-lint run ./...` — 0 issues (root + 7)
- `nix flake check` — all checks passed
- `docs-gates run .` — 5/5 PASS (freshness, audit-all, annotate, routing,
  cadence)
- v5 core coverage: 89.3%
- CG002 corpus re-verify with CurrentMajor=v5: erraudit 12×CG002,
  timesheets 15×CG002 + 1×CG003 (known TP), gsl/bf clean, probe fires —
  zero false positives

## c. Decisions made (owner may override)

1. ~~**Scope keeps `HealthCheck()`/`HealthCheckWithContext(ctx)`** with the~~ done — decision held in v5.0.0: the no-arg forms delegate via context.Background() (scope.go)
   ~~no-arg forms delegating via `context.Background()`; `Scope` deliberately~~
   ~~does NOT implement `do.HealthcheckerWithContext` (that's CLI's contract).~~
2. ~~**Schema (CT2): hold workspace-local** until v5.0.0 ships, then tag~~ done — executed as decided: held through v5.0.0, then schema/v0.1.0 was tagged (tag exists)
   ~~`schema/v0.1.0` per the lint playbook.~~
3. ~~**NewApp signature unchanged** (build callback stays `func(*CLI[T]) error`);~~ done — held: NewApp's build callback is still func(*CLI[T]) error (cli_batteries.go); NewApp-on-Builder remains a ROADMAP raw idea
   ~~the Builder composes inside it. NewApp-on-Builder stays a v5 API-review~~
   ~~question (ROADMAP candidate).~~
4. ~~**Empty config file = absent** (uniform across YAML/TOML/JSON).~~ done — shipped in v5.0.0 (CHANGELOG [5.0.0] Fixed: empty config files no longer crash startup)
5. ~~**BuildFlow commit left unpushed** (needs owner permission, per 03:56 §g).~~ done — pushed with the master merge (2026-11-11)

## d. Open items (honest residue)

- The v5-scaffold branch is NOT tagged; v5.0.0 release (tag + GH release +
  proxy verify + BuildFlow consumer bump) is the next milestone. Publish-shape
  go.mod work per the 2026-11-02 playbook is required first.
- ~100 non-Go `v4` references in docs/website remain (M16 full docs sweep —
  pull-based; QUICKSTART/TUTORIAL/README done this session).
- Website ports, marketing bundle, FR enhancements: untouched (Tier D
  pull-based by plan).
- BuildFlow's local commit (`b6ef3f608`, `136e53a49`) still needs a push
  decision (§g of the 03:56 report).
- Timesheets `internal/parser` is broken in the owner's working tree
  (pre-existing WIP; not touched).

## e. Self-critique

- The default-flip fallout was bigger than the 03:56 report estimated
  (~35 call sites vs "fix fallout tests" hand-wave) — the mechanical part
  was unavoidable once strict became default; the two REAL bugs it exposed
  (stale cobra ctx, per-exec shutdown) justify the flips being done with the
  full suite rather than piecemeal.
- The first "CG002 corpus re-verify" was invalid (positional arg silently
  ignored by the CLI — it linted cmdguard itself). Caught by probing; the
  corrected run is the recorded one. Footnote: cmdguard CLIs could benefit
  from a strict-args default for subcommands (ROADMAP candidate).
- The docs purge was reactive (cadence gate went red at 42 strikes) rather
  than mid-session; next time, strike+purge in the same pass as the work.
- Tier D website/marketing bundles were not attempted beyond bounded slices
  — correct per the plan's pull-based rule, but the "WHOLE list" framing
  means they remain open by design, not omission.
