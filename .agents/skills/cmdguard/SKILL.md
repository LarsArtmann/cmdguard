---
name: cmdguard
description: Build, extend, or review Go CLIs built with cmdguard v4 (`github.com/larsartmann/cmdguard/v4`) — the LarsArtmann Cobra wrapper with typed struct-tag flags, per-command flag structs, DI, middleware, config files, and exit-code contracts. Use whenever a project imports cmdguard, when adding subcommands or flags (especially per-command or scoped flags), when wiring config/env precedence, output formats, middleware, prompts, spinners, or telemetry, or when any cmdguard API question arises — even when cmdguard is not named ("add a flag only to the X subcommand", "why is --foo accepted on every command", "unknown flag" surprises). Also covers the cmdguard-lint rules and the optional sub-modules. Not for developing cmdguard itself (that repo's AGENTS.md is authoritative for contributors).
---

# cmdguard — validated Cobra CLIs with type-safe flags and DI

cmdguard v4 wraps spf13/cobra with struct-tag flags, samber/do DI, middleware,
config-file loading, rich output, and a no-panics API. Everything returns errors.
Module: `github.com/larsartmann/cmdguard/v4`, package `pkg/cmdguard/v4`.
Published: v4.1.0 (2026-10-04, pkg.go.dev-verified).

## Verification status

All rows below verified against source in `~/projects/cmdguard` on 2026-10-16,
unless the Source column says otherwise.

| Claim                                                                        | Status      | Source                                                                                          |
| ---------------------------------------------------------------------------- | ----------- | ----------------------------------------------------------------------------------------------- |
| Per-command F-struct flags registered via `cmd.Flags()` (local, never inherited) | ✅ Verified | `cli_command.go:203` + `flags.go:69-71`; end-to-end regression test `cli_flag_scoping_test.go`   |
| Root T fields persistent by default; `local:"true"` → root-only                | ✅ Verified | `flags.go:85-104` (`RegisterScopedFlags`), wired at `cli.go:246`; `flags_scoped_test.go`         |
| `noop --depth` / `scan --build` fail with "Unknown flag"; root accepts `--build` | ✅ Verified | `cli_flag_scoping_test.go` (4 subtests, race-enabled)                                           |
| Local-flag default + env still apply on commands lacking the flag              | ✅ Verified | `flags_parse.go:26-38`                                                                          |
| Precedence: flag → env → default (default pre-updated by config file)          | ✅ Verified | `flags_parse.go:45-71` + `cli.go:230-235` (`updateTagDefaultsFromConfig`)                        |
| `cli.RegisterLocalCommandFlags(cmd)` opts a subcommand into root-local flags   | ✅ Verified | `cli_accessors.go:112`, `flags.go:115-137` (subcommand's own flag definition wins)               |
| Full tag set: flag, short, default, help, values, required, count, local, hidden, validate, env, prompt | ✅ Verified | `config_parsing.go:123-233`, `FlagTag` struct at `config.go:20-36` |
| Nested config structs recurse (except registered value types)                 | ✅ Verified | `config_parsing.go:100-118`                                                                      |
| `ArgsFromContext(ctx)` retrieves positional args                              | ✅ Verified | `cli_command.go:19`                                                                              |
| pkg.go.dev page for the package exists (v4.1.0)                               | ✅ Verified | Fetched live 2026-10-16                                                                          |
| 59 exported sentinel error vars; 87 test funcs in `examples/taskctl`           | ✅ Verified | `rg` over `errors*.go` / taskctl tests, 2026-10-16                                                |
| CG001-CG008 lint rule semantics                                               | 📖 Repo-documented | cmdguard AGENTS.md (lint section) — not independently re-verified                            |
| `Describe()` / hidden `describe` command contract                             | ✅ Verified | `pkg/cmdguard/v4/describe.go` + `describe_test.go`, 2026-11-01                                  |
| `AddCommands` atomic batch registration                                       | ✅ Verified | `pkg/cmdguard/v4/cli.go` + `describe_test.go` (atomicity, mixed F, name-keyed duplicates), 2026-11-01 |
| `ErrorHint` coverage of classified sentinels                                  | ✅ Verified | `errors_hints.go` + `errors_hints_test.go` (coverage test), 2026-11-01                          |

## The mental model: one type param on the CLI, a flag struct per command

```go
cli, err := cmdguard.NewCLI("app", "App description", AppConfig{}) // T = root config

scan, err := cmdguard.NewCommand("scan", ScanFlags{},               // F = THIS command's flags
    func(ctx context.Context, cfg *AppConfig, f ScanFlags) error { ... },
    cmdguard.WithShort("Scan things"),
)
err = cmdguard.AddCommand(cli, scan)

err = cli.ExecuteAndExit(ctx) // blessed entry point
```

- `CLI[T]` — `T` is the root config struct; its fields become the root's flags.
- `Command[T, F]` — `F` is this command's OWN flag struct, parsed into the handler's
  third parameter. Commands with no own flags use `cmdguard.NoFlags{}`.
- Type params are inferred at `NewCommand` — no explicit instantiation needed.
- Parents: `NewParentCommand[T, F]("db", "Long description (required)", NoFlags{},
  cmdguard.WithSubcommands(migrate, seed))`. Subcommands must share T and F.
- Handlers receive `(ctx, *T, F)` — no positional args. Get them with
  `cmdguard.ArgsFromContext(ctx)`; cmdguard stores them automatically.

## Machine-readable contract (`describe`)

Every CLI gets a hidden `describe` command at `Execute` time; it prints a
versioned JSON contract for agents and tooling:

```bash
app describe
# {"schema_version":"1","name":"app","version":"1.0.0","flags":[...],"commands":[{"name":"scan","path":"app scan","flags":[...]}]}
```

- Programmatic: `cli.Describe() *CLIDescription` and `cli.WriteDescription(w)`.
- Flag entries carry `type`, `default`, `env` (prefix already applied), `values`
  (enums), `required`, `count`, `prompt`, `validate`, `hidden`, `inherited`.
- The schema is additive within `schema_version`; names and meanings never change
  inside a major.
- A consumer command named `describe` wins over the built-in; `WithoutDescribeCommand()`
  opts out entirely (the programmatic API keeps working).

Prefer this over scraping `--help` when driving a CLI from an agent or a test.

## Registering many commands

`AddCommands(cli, a, b, c)` collapses the per-command error ladder into one check
and registers atomically: every command is validated and converted before any is
attached, so a failed batch leaves the CLI unchanged. Commands may have different
flags types; `T` is still type-checked.

```go
if err := cmdguard.AddCommands(cli, listCmd, addCmd, doneCmd); err != nil {
    return err
}
```

Duplicate detection keys on the cobra command name, so `"list [path]"` and
`"list"` correctly collide.

## Flag scoping — read this before claiming something is impossible

**Yes, you CAN scope a flag per-command. Three scopes exist:**

| Where the field lives                                    | Scope                                                       |
| --------------------------------------------------------- | ------------------------------------------------------------ |
| Root `T` struct (plain)                                   | Persistent: inherited by EVERY subcommand                    |
| Root `T` struct + `local:"true"`                           | Root only — not inherited by any subcommand                  |
| A command's own `F` struct (e.g. `ScanFlags`)              | That command only — never inherited, rejected elsewhere      |

Mechanics (all source-verified):

- Root `T` flags go through `RegisterScopedFlags` (`cli.go:246`): default persistent,
  `local:"true"` lands on `cmd.Flags()` (root-only).
- Every command's `F` flags go through `RegisterFlags` → `cmd.Flags()` (`cli_command.go:203`)
  — cobra-local, so they exist ONLY on that command. `noop --depth` fails with
  "Unknown flag: --depth" when `depth` lives in `ScanFlags`. That is the designed
  behavior, not a limitation. Pinned by `cli_flag_scoping_test.go`.
- A subcommand that needs a root-local flag opts in:
  `cli.RegisterLocalCommandFlags(cobraCmd)` registers the root's `local:"true"` flags
  on it. If the subcommand already defines a flag with the same name, its own
  definition wins.
- `hidden:"true"` hides a flag from `--help` anywhere, fully functional otherwise.

Subtleties worth knowing:

- A `local:"true"` flag's declared DEFAULT (and its `env` var, if tagged) still
  applies on commands where the flag is not registered — a default is a default
  regardless of which command runs. Only the explicit CLI value is impossible there
  (`flags_parse.go:26-38`).
- The bare root run (no subcommand) shows help and short-circuits BEFORE
  `PersistentPreRunE` parses flags into the config: root-local flags are ACCEPTED on
  the bare root, but their values only materialize in `cfg` when a command actually
  executes.

**What does NOT exist:** excluding one inherited root-persistent flag from one
specific subcommand ("persistent everywhere except `foo`"). The idiomatic fix is not
an exclusion mechanism — it is putting the flag where it belongs:

- Used by exactly one command? → that command's `F` struct.
- Used by root run + a few subcommands? → root `T` with `local:"true"` +
  `RegisterLocalCommandFlags` opt-in on each consumer.
- Genuinely global (verbosity, output, config path)? → root `T`, plain persistent.

The classic smell (seen in erraudit, 2026-10-16): stuffing every flag into the root
`AppConfig` makes them all persistent, so subcommands silently accept flags they
ignore ("dead flags"). That is a consumer-side modeling problem, not a cmdguard
limitation. Migration recipe:

1. Define the flag as a field in the consuming command's `F` struct (create the
   struct if the command used `NoFlags{}`), with the usual struct tags, e.g. a
   `DryRun bool` field tagged `flag:"dry-run"` plus `help:"..."`.
2. Move any `env`/`validate`/`values` tags along with it; keep `local:"true"` OFF —
   F-struct flags are already command-local.
3. Read the value from the handler's `flags` parameter instead of `cfg`.
4. Drop the field from root `T`. Any other command still passing the flag now gets
   the desired "Unknown flag" error instead of silent acceptance.
5. Config-file users of that key need the value re-homed (config file values feed
   root `T` defaults) — check `.yaml` consumers before migrating.

## Struct tag reference

A field becomes a flag only when it has a `flag` tag (`flag:"-"` skips it).

| Tag                        | Effect                                                                   |
| -------------------------- | ------------------------------------------------------------------------ |
| `flag:"name"`              | Flag name (required for the field to register at all)                     |
| `short:"v"`                | Single-character shorthand                                               |
| `default:"..."`            | Default value (zero value when absent)                                   |
| `help:"text"`              | Help text                                                                |
| `values:"a,b,c"`           | Enum constraint (validated → `ErrInvalidEnum`)                            |
| `required:"true"`          | Must be explicitly set (bool tag)                                         |
| `count:"true"`             | Counting flag (`-vvv` → 3); field must be `int` (bool tag)                |
| `local:"true"`             | Owning command only, not inherited (bool tag)                             |
| `hidden:"true"`            | Excluded from `--help`, fully functional (bool tag)                       |
| `validate:"min=1,max=100"` | Rule-based validation (`ErrValueOutOfRange` etc.)                         |
| `env:"VAR_NAME"`           | Environment fallback (prefixable via `WithEnvPrefix`)                     |
| `prompt:"Question?"`       | Interactive prompt when missing (needs `WithPromptOnMissing` + the `prompts` sub-module) |

Nested structs recurse automatically into flattened flags — EXCEPT registered value
types (`Duration`, `Port`, `Email`, `URL`, `HostPort`, `Enum`, `FilePath`, `LogLevel`,
`time.Duration`, `time.Time`), which stay atomic.

Precedence (source-verified, `flags_parse.go:45-71` + `cli.go:230-235`): explicit
flag → `env` → default, where the config file has already overwritten the defaults
before flags register. Effectively: flag → env → config file → tag default. Config
files (`WithConfigFile(paths...)`, JSON/YAML/TOML auto-detected, `$VAR`/`~` expanded,
missing files skipped) feed root `T` only.

## Error and exit contract

- cmdguard owns error display: the error prints exactly ONCE (fang when enabled,
  cobra when not). Never re-print the error returned by `Execute` — that double-prints.
- `ExecuteAndExit` is the blessed entry point. Need work after execution (flush,
  audit export, teardown)? Use `ExitCode(err)` (0 for nil) and exit yourself.
- `NewExitError(code, err)` returns an `*ExitError` for the 0-255 range;
  `ExecuteAndExit` honors `ExitCoder` implementations.
- `SilenceUsage` is true by default (no usage-on-error spam); `WithoutSilenceUsage()`
  re-enables it.
- 59 exported sentinel errors, all `fmt.Errorf("%w: ...", sentinel)`-wrapped — reach
  them with `errors.Is`, typed errors with `errors.AsType[*T]` (Go 1.26+). Full list:
  `docs/ERROR_REFERENCE.md`.
- `ErrorHint(err)` returns one actionable fix sentence per classified sentinel
  (empty when unknown). The JSON error envelope includes it as `hint`; use it in
  custom handlers instead of duplicating fix text.

## Gotchas that bite consumers

- **`PostRunE` does not run when `RunE` errors** (Cobra semantics). Hooks that must
  fire on BOTH success and error go through `WithCleanup[T](fn)` — the original error
  is never swallowed; cleanup errors join via `errors.Join`.
- `NoFlags` is a distinct named type, not an alias — compare/call with `(NoFlags{})`.
- `ConfigFromContext[T](ctx)` retrieves the resolved root config from any cobra
  command context — the escape hatch for raw cobra subcommands mixed into the tree.
- `WithPostFlagParse[T](fns...)` runs after flag parsing + config validation but
  before handlers — the replacement for hand-wrapped `PersistentPreRunE`.
- Signal handling: `WithSignalHandling()` cancels the context; `WithGracefulShutdown()`
  additionally shuts down DI services (implies the former).
- Integer flags `int8`/`int16`/`int32`/`uint8`/`uint16` are range-checked → `ErrValueOutOfRange`.
- Counting flags must not reuse flag names from the root config.

## Optional sub-modules (import only what you need)

Core has zero dependencies on these. Each lives at the repo root and imports as
`github.com/larsartmann/cmdguard/<name>` — directory placement is load-bearing
(moving one under `pkg/` breaks downstream `go get` while local builds still pass).

| Sub-module        | Purpose                                                          | Key API                                            |
| ----------------- | ---------------------------------------------------------------- | -------------------------------------------------- |
| `glamour`         | Markdown help rendering for `WithHelpTransform[T]`               | `RenderWithEnvironmentConfig` (`GLAMOUR_STYLE`)    |
| `prompts`         | huh/v2 interactive prompts for `prompt:"..."` tags               | wired via `SetPromptRunner()`                       |
| `spinner`         | Terminal spinner middleware (auto-skips non-TTY stderr)          | `SpinnerMiddleware[T]`                              |
| `telemetry`       | OpenTelemetry span per command                                   | `TelemetryMiddleware[T]`                            |
| `flightrecorder`  | Go 1.25+ `runtime/trace` flight recorder, zero external deps     | `WithFlightRecorder[T](cfg)`                        |
| `lint`            | cmdguard usage linter (go-finding; BuildFlow provider + CLI)     | `cmdguard-lint`, rules CG001-CG008                  |

cmdguard-lint rules worth knowing while writing consumers (repo-documented): CG001
execute bypass (`Execute` without exit handling), CG002 stale major version, CG003
constructor panic (ignoring `New*` errors), CG004 runtime `SetVersion`, CG005
duplicate version options, CG006 re-printing the execute error, CG007 missing
`WithShort` (blank help line), CG008 missing `WithCLIVersion`. Suppress with
`//cmdguard-lint:ignore <ID> <reason>`.

## Where deeper truth lives (cmdguard repo)

- `AGENTS.md` — contributor context, design principles, full gotcha list
- `docs/QUICKSTART.md`, `docs/TUTORIAL.md` — guided walkthroughs
- `docs/ERROR_REFERENCE.md` — the sentinel error catalog
- `docs/MIGRATION_FROM_COBRA.md`, `docs/MIGRATION_v3_v4.md` — migration paths
- `docs/COBRA_FOOTGUNS.md` — cobra behaviors cmdguard neutralizes
- `pkg/cmdguard/v4/flags.go`, `cli_command.go`, `cli.go`, `flags_parse.go` — flag
  scoping and precedence source of truth
- `pkg/cmdguard/v4/cli_flag_scoping_test.go`, `flags_scoped_test.go` — scoping
  regression tests (end-to-end and registry-level)
- `examples/taskctl/` — flagship example (87 test functions)
- pkg.go.dev/github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4 — godoc (verified live)
