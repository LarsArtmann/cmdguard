# Feedback: `upd` Evaluation of `cmdguard`

**Date:** 2026-07-16
**From:** [`upd`](https://github.com/LarsArtmann/upd) — a focused single-command CLI that bumps NPM dependency versions in `package.json` while byte-preserving JSON formatting.
**Outcome:** Did **not** adopt. Documenting why and what could change that.

---

## Context

`upd` is a small, single-purpose CLI (~2,100 lines of non-test Go, 4 direct dependencies, 7 MB binary). It uses Go's stdlib `flag` package with ~100 lines of flag parsing. We evaluated `cmdguard` v3 as a potential replacement for our hand-rolled CLI layer.

---

## What Appealed to Us

These features caught our attention and are genuinely valuable:

1. **Type-safe flags via struct tags** — eliminates the `defineBoolFlag`/`defineStringFlag` boilerplate we wrote for short+long dual registration.
2. **Typo suggestions** — "did you mean --json?" — we don't have this and users would benefit.
3. **Env var support via `env:""` tag** — clean, declarative. We currently do `os.LookupEnv` manually.
4. **Auto-generated styled help** — our `PrintUsage()` is hand-maintained and drifts from actual flags.
5. **Shell completion** — free bash/zsh/fish completion is a nice-to-have.
6. **Zero-panics contract** — aligns with our own philosophy. Every function returns errors.
7. **Same author ecosystem** — `go-error-family` + `go-atomic-write` + `cmdguard` share design DNA. Conceptually a great fit.

---

## Why We Didn't Adopt

### 1. Dependency weight is the primary blocker

| Metric                 | `upd` (current) | With `cmdguard` |
| ---------------------- | --------------- | --------------- |
| Direct dependencies    | 4               | ~14             |
| Total modules          | 11              | ~92             |
| Binary size (stripped) | 7.1 MB          | ~15–17 MB       |
| go.sum lines           | 22              | ~184            |

For a tool that edits one JSON file, doubling the binary and adding 80+ transitive modules is a hard sell. `upd`'s leanness is an explicit design principle — users notice the binary size and startup speed of a focused tool.

### 2. Most features are unused for single-command CLIs

`upd` has no subcommands, no DI, no service lifecycle, no config files, no multi-format output. cmdguard's headline differentiators — `samber/do` integration, `BranchingFlowContext`, `DoctorCommand`, `WithGracefulShutdown`, 16 output formats via `go-output` — provide zero functional benefit here.

We needed maybe **10%** of what cmdguard offers. The other 90% is dead weight in our binary.

### 3. Error system overlap

`upd` already uses `go-error-family` for structured error classification (Rejection/Transient/Corruption/Conflict with exit codes). cmdguard brings its own error contract (`ExitCoder`, `NewCommandError`, `NewFlagError`). Having two overlapping error systems would create confusion about which owns the exit-code mapping.

### 4. Output is domain-specific

Our table renderer (`render.go`, 383 lines) produces domain-specific output: version arrows (`1.2.3 → 2.0.0`), state badges (`updated`, `kept`, `skipped`, `error`), verbose error chains. cmdguard's `go-output` table/JSON/CSV system doesn't replace this — it solves a different problem.

### 5. Signal handling already done

`signal.NotifyContext(SIGINT, SIGTERM)` is 2 lines in `main.go`. cmdguard's `WithSignalHandling()` is a wash.

### 6. Unique domain patterns don't map cleanly

- `upd` reads an `upd` field from `package.json` and prepends those args to CLI patterns. This is domain-specific logic that doesn't map to cmdguard's config file system.
- `--dry-run` is manually aliased to `--nop`. This kind of domain-specific flag aliasing works fine with stdlib `flag`.

---

## Improvement Suggestions

These changes would make `cmdguard` viable for CLIs like `upd`:

### Suggestion 1: A "lite" mode or stripped-down package

**Problem:** The full `cmdguard` package drags in Cobra, pflag, lipgloss, fang, glamour, goldmark, chroma, koanf, samber/do, go-output, and ~80 more modules.

**Proposal:** Consider a `cmdguard/lite` sub-module or build tag that provides:

- Struct-tag flags with short/long registration
- Typo suggestions
- Env var support
- Auto-generated help (plain text, no lipgloss)

...without Cobra, DI, koanf, go-output, or fang. Target: <10 transitive modules, stdlib `flag` compatible.

This would capture the "stdlib flag is getting tedious but I don't need a framework" market — which is large.

### Suggestion 2: Modular dependency boundaries

**Problem:** Even importing just the flag-parsing types pulls in the entire dependency graph.

**Proposal:** Split into smaller sub-modules (like you already did for glamour/prompts/telemetry/spinner):

- `cmdguard/flags` — type-safe flag structs, typo suggestions, env vars (no Cobra)
- `cmdguard/core` — the CLI struct, command registration, help generation
- `cmdguard/di` — samber/do integration (opt-in)
- `cmdguard/output` — go-output integration (opt-in)

Each module should have minimal transitive deps. Users import only what they need.

### Suggestion 3: Integrate with `go-error-family`

**Problem:** cmdguard defines `ExitCoder`, `NewCommandError`, `NewFlagError` — overlapping with `go-error-family`'s Family/ExitCode system.

**Proposal:** Since both are `larsartmann` projects, consider making cmdguard natively aware of `go-error-family`. If a returned error implements `errorfamily.Error`, use `Family.ExitCode()` for the exit code mapping. This would let consumers use one error system end-to-end.

### Suggestion 4: Single-command fast path

**Problem:** cmdguard assumes a multi-command structure (`NewCLI`, `AddCommand`, subcommand trees). For single-command CLIs this is ceremony.

**Proposal:** A `Run[T, F](name, desc, config, flags, handler, opts...)` function that sets up everything for a single-command CLI in one call — no `NewCLI` + `AddCommand` + `ExecuteAndExit` dance.

### Suggestion 5: Help generation from existing `flag.FlagSet`

**Problem:** Users with existing stdlib `flag` code can't incrementally adopt cmdguard's better help without migrating everything.

**Proposal:** A `cmdguard.PrintHelp(flagSet *flag.FlagSet, name, desc string)` function that auto-generates aligned help from an existing `flag.FlagSet`. This is a low-commitment entry point — adopters get better help without rewriting their flag definitions.

---

## Summary

`cmdguard` is an excellent framework that solves real problems for **multi-command, service-oriented CLIs**. For focused single-purpose tools like `upd`, the dependency cost vastly exceeds the value extracted.

The gap between "stdlib `flag`" and "full CLI framework" is currently unoccupied. A lightweight, modular cmdguard that lets users adopt incrementally — starting with just better flags or better help — would capture a much broader audience.

We'll revisit if a lite/modular variant ships. Until then, we're staying on stdlib `flag` and cherry-picking individual features (~50 lines of code for typo suggestions + env vars).
