# v5 "Paved Road" - End-Goal Plan

> Status: proposal, current as of 2026-11-01.
> Owner: Lars Artmann. Executor: whoever ships v5.
> Related: [`ROADMAP.md`](../../ROADMAP.md), [`TODO_LIST.md`](../../TODO_LIST.md),
> [`WHAT_THIS_PROJECT_IS_NOT.md`](../../WHAT_THIS_PROJECT_IS_NOT.md).

## The end goal (verbatim)

> Make it IMPOSSIBLE to NOT create SUPERB CLI apps with near ZERO-BUGs and
> fantastic User Experience while keeping Code requirements developers and AI
> Agents need to write to an absolute minimum!

## What that decomposes into

Four testable properties. A change that does not move one of them is not a v5
change, it is maintenance.

| # | Property | Question the release must answer "yes" to |
| - | -------- | ----------------------------------------- |
| P1 | Inevitable quality | Can a developer who ignores every doc still end up with a superb CLI? |
| P2 | Zero-bug by construction | Is every mistake caught at construction, lint, or first run, never in production? |
| P3 | Fantastic UX | Does every surface (help, errors, exit codes, docs, completion) teach the user what to do next? |
| P4 | Minimum code | Is handwritten code proportional to the domain, not to the framework? |

**Corollary:** v5 is not a rename release. The renames already on the roadmap
are hygiene; the release earns its major number by moving P1-P4.

## v4 reality check (measured, not asserted)

| Surface | v4 today | Gap to the end goal |
| ------- | -------- | ------------------- |
| Bootstrap | `NewCLI` + one error check per command + one `AddCommand` error check per command. `examples/taskctl/commands.go` has 10 `AddCommand` error blocks for 12 commands. | Every registration is a 3-line error stanza. P4 fails. |
| Safety defaults | `WithSignalHandling` / `WithGracefulShutdown` are opt-in. `Strict`/`Draconian` validation is opt-in. | An unopinionated developer gets no signal handling and no short-description checks. P1 fails. |
| Error output | Sentinel errors are precise and `errors.Is`-chainable, but the text says what failed, rarely what to do. | Humans re-read the docs; agents guess. P3 fails. |
| Machine readability | `--output json` affects data output; errors become a JSON envelope. The CLI itself has no schema. | An agent must scrape `--help` prose to drive an unknown binary. P4 fails for agents. |
| Guardrails | Eight lint rules: CG001-CG006 protect the execution contract, CG007/CG008 (shipped 2026-11-01) enforce the short-description and version baseline. | Nothing yet enforces dead-flag hygiene, and the baseline rules default to warnings rather than failures. P2/P3 gap. |
| Versioning of the contract | Sentinel errors are stable; the CLI description has no version. | Tooling cannot pin a schema. |

Evidence lives in `pkg/cmdguard/v4/cli.go`, `cli_command.go`,
`examples/taskctl/commands.go`, and `lint/rules.go`.

## The v5 thesis: "Paved Road"

Three sentences the release must make true:

1. **The default is the superb path.** A developer writes the domain handler;
   cmdguard supplies the lifecycle, the version command, the schema, the
   signals, the exit codes, and the error hints.
2. **Deviating is loud.** Wrong shapes fail at `NewCLI` / `AddCommand` time or
   are flagged by `cmdguard-lint`; they do not fail at 2am.
3. **The binary explains itself.** Humans get hints; agents get JSON. Neither
   needs the source.

## Workstreams

Each workstream is independently shippable and has an objective gate. Order
matters: W1 and W2 are the release's identity, W5 makes them safe, W7 pays for
the break.

### W1. Machine-readable CLI contract (P3, P4)

**Outcome:** every cmdguard CLI can describe itself.

- `CLIDescription` / `CommandDescription` / `FlagDescription` with
  `schema_version: "1"`.
- `(*CLI[T]).Describe()` and `(*CLI[T]).WriteDescription(io.Writer)`.
- A hidden `describe` command, registered automatically at `Execute` time,
  skippable with `WithoutDescribeCommand()`; a consumer-defined `describe`
  command wins.
- Flag metadata includes type, default, shorthand, usage, enum `values`,
  `env`, `required`, `count`, `hidden`, `inherited`. Root tags come from the
  registry; command tags ride on cobra annotations attached at wiring time.
- Stable, documented schema evolution rule: additive fields only within a
  schema major.

**Gate:** a fixture CLI with nested commands, inherited flags, enums and env
tags produces byte-stable JSON; `go test -race` covers root-only, nested and
opt-out paths.

**Status:** shipped pre-v5 (additive) as `describe.go` so consumers can adopt
before the major bump.

### W2. Batteries: the one-call superb app (P1, P4)

**Outcome:** no decisions are required to get a production-grade CLI.

- `WithBatteries()` (name TBD at API review) turns on: describe command,
  version command (when a version is set), doctor wiring, completion,
  signal handling, graceful shutdown, strict validation, hint-rendered errors.
- `NewApp[T](name, short, defaults, build func(*Builder[T]) error, opts...)`
  is the one-call entry point: it constructs the CLI, lets `build` register
  commands with deferred errors, and validates the whole tree once.
- Any "battery" remains individually opt-out-able. Defaults flip in v5; the
  individual options stay.

**Gate:** the taskctl bootstrap shrinks demonstrably (target: under 25 lines
from `main` to a running CLI, excluding domain handlers) with no behavior
regression in `examples/taskctl/*_test.go`.

**Status:** partially shipped pre-v5 (`AddCommands` removes the per-command
error stanzas). `NewApp`/`WithBatteries` are v5.

### W3. Errors that teach (P3)

**Outcome:** every user-facing failure names the fix.

- `ErrorHint(err error) string` maps each sentinel to one actionable sentence;
  empty when no hint is known.
- JSON error envelope gains `hint`; the human path renders the same hint under
  the error (fang handler) without duplicating the message.
- Every sentinel in `docs/ERROR_REFERENCE.md` carries: cause, hint, exit code,
  and the doc anchor. The reference is generated from the code so it cannot
  drift.
- No error may be added without a hint (test-enforced table coverage).

**Gate:** a coverage test asserts every exported sentinel has a hint or an
explicit `no hint by design` entry; JSON envelope snapshot includes `hint`.

**Status:** partially shipped pre-v5 (`ErrorHint` + `hint` in the JSON
envelope). Fang rendering is v5.

### W4. Minimum-code ergonomics (P4)

**Outcome:** registration is declarative, not a sequence of error checks.

- `AddCommands[T](cli, cmds...)` batches registration atomically: all commands
  are prepared and validated before any is attached.
- Duplicate detection keys on the cobra command name, not the raw `Use` string,
  so `"list [path]"` and `"list"` collide correctly.
- `Builder[T]` collects construction errors and fails once at `Build()`,
  killing the `cmd, err := ...; if err != nil` ladder.
- Optional flag-less handler adapters (`NoFlags` dropped from the call site)
  only if the API review proves they do not create a second mental model.

**Gate:** taskctl's `buildCommands` loses the 10 error stanzas; duplicate-name
regression test covers `"list [path]"` vs `"list"`; a failed batch leaves the
CLI unchanged.

**Status:** `AddCommands` and name-keyed duplicate detection ship pre-v5; taskctl now registers 11 commands in one call (432 -> 402 LOC in `commands.go`, 10 error stanzas removed).

### W5. Guardrails: make the bad shape impossible (P1, P2)

**Outcome:** the linter enforces the superb baseline, not just the absence of
regressions.

- `CG007` missing short description: a command without `WithShort` renders as a
  blank line in `--help` and fails `WithStrictValidation`.
- `CG008` missing CLI version: no `WithCLIVersion` means no `--version`, no
  build identification, no operator support path.
- Candidate follow-ups (only with corpus evidence): dead root-persistent flags,
  unvalidated root config, `NoFlags` on a command that reads config flags.
- Defaults in v5 move from "lenient" to "strict"; `WithLenientValidation()`
  becomes the escape hatch. A consumer upgrade that breaks is a consumer with a
  blank help line, which is the point.
- Lint baselines are the migration path: the first v5 run reports the debt, the
  ratchet prevents new debt.

**Gate:** both rules pass positive/negative corpus tests, the generated README
table is in sync, and running the linter over the consumer corpus produces a
reviewed, baselined finding set.

**Status:** CG007/CG008 ship pre-v5 so consumers can clean up before the
default flips.

### W6. Contract honesty: the v5 break (P1)

**Outcome:** a single, documented, migration-guide-backed break.

- Execute the `TODO(v5)` items: `CommandInfo` -> `CommandMetadata`,
  `TypeHandler` -> `TypeCodec`, `PromptRunner` -> `HuhPrompter`, remove
  `SetConfig`, replace `RegisterInScope`.
- Additional candidates (decide, do not drift): `Get[T]` -> `GetService[T]`,
  `Package[T]` redesign, unify `Middleware[T]` and `ContextMiddleware[T]`.
- Module path moves to `/v5` (or the owner explicitly decides to keep `/v4`
  and treat v5 as a minor, which would contradict semantic import versioning;
  do not leave this undecided).
- `docs/MIGRATION_v4_v5.md` with a mechanical diff per rename and a codemod
  script for the trivial ones.

**Gate:** every rename has a before/after compile-tested example; migration
doc is validated by migrating the three consumer repos (cmdguard itself,
erraudit, timesheets) on branches.

### W7. Proof, not promises (P2, P3)

**Outcome:** the release is measured against the goal, not against its own
tests.

- Consumer corpus: migrate `erraudit`, `timesheets`, and `taskctl`; record the
  diff in LOC, error stanzas removed, lint findings before/after.
- Golden CLI fixtures: help output, describe JSON, error output, exit codes
  (including `--output json` envelopes) snapshotted per release.
- Runtime gates unchanged and enforced: `-race`, fuzz targets, 0 lint issues,
  coverage floors (core >= 87%, sub-modules >= 80%).
- Performance gate: `benchmarks/` before/after, no regression without a
  written justification.

**Gate:** a `docs/status/` report with the numbers, linked from CHANGELOG.

### W8. Teaching surface (P3, P4)

**Outcome:** a new developer (or agent) reaches a superb CLI in one screen.

- QUICKSTART becomes a single copy-pasteable program that is already superb
  (version, signals, describe, exit codes).
- `.agents/skills/cmdguard/SKILL.md` documents `Describe`, `AddCommands`,
  batteries, and the lint baseline; verification table refreshed.
- Website landing page leads with the end goal and a live `describe` output.

**Gate:** a fresh session, given only README + skill, builds a CLI that passes
`WithDraconianValidation()` and `cmdguard-lint` with zero findings on the first
attempt.

## Sequencing

| Order | Workstream | Why this order |
| ----- | ---------- | -------------- |
| 1 | W1 describe, W4 AddCommands, W3 hints, W5 CG007/CG008 (all additive) | Ships value now, de-risks v5, and gives consumers time to adopt |
| 2 | W2 batteries + `NewApp` (additive, defaults flipped in v5) | Needs W1/W4 primitives |
| 3 | W8 docs/skill refresh | Documents what actually exists |
| 4 | W6 the break + migration guide | Only after the additive surface is proven |
| 5 | W5 strict-by-default flip, W7 corpus proof, release gate | Last, so the corpus migrates on a stable API |

## Release gate (definition of done for v5)

- [ ] P1: `WithBatteries` + strict-by-default; a no-option CLI has signals,
      describe, version (when set), hints, and strict validation.
- [ ] P2: `cmdguard-lint` runs clean on the three-consumer corpus with a
      committed baseline; CG007/CG008 active; `-race` and fuzz green.
- [ ] P3: every sentinel has a hint and a generated ERROR_REFERENCE row;
      help/describe/error/exit-code golden fixtures pass.
- [ ] P4: taskctl bootstrap under 25 lines; 10 error stanzas gone; no
      framework concepts leak into handlers.
- [ ] W6 renames shipped with a validated migration guide and codemod.
- [ ] CHANGELOG + status report publish the measured deltas.

## Decision points for the owner

These cannot be guessed; each blocks part of the release.

1. **Module path:** `/v5` (semantic import versioning, heavy but correct) or a
   deliberate exception? No third option.
2. **Strict by default:** accept that some consumers fail to construct after
   upgrade (blank help lines are the failure), or keep lenient and lose P1.
3. **Batteries name and scope:** `WithBatteries`, `WithStandardCommands`,
   something else; which batteries are on by default.
4. **`NoFlags` fate:** keep the explicit type or add flag-less constructors.
   Two ways to write the same command is a split brain; pick one.
5. **Doc generation:** is `describe` (JSON) the canonical contract, with
   markdown docs derived from it? Deriving prevents two sources of truth.

## Non-goals (unchanged from v4)

Not an application runtime, not a code generator, not a replacement for Cobra
or Fang, not a logging/ORM/web framework. See
[`WHAT_THIS_PROJECT_IS_NOT.md`](../../WHAT_THIS_PROJECT_IS_NOT.md).

Every workstream above must be achievable without violating those boundaries.
If it is not, the workstream is wrong, not the boundary.
