# Migration Guide: v4 → v5

> **TL;DR:** Run the codemod (`docs/migration_v4_v5/codemod.sh`), rename four
> identifiers, and review the new defaults. The batteries bundle (strict
> validation, signal handling + graceful shutdown, strict config keys) is now
> ON by default; every piece has a `Without*` opt-out. Construction is
> ladder-free with the new `Builder`.

**Source version:** `v4.2.1` (the last stable v4)
**Target version:** `v5.0.0` (`github.com/larsartmann/cmdguard/v5`)
**Why:** v5 executes the deferred v4 breaking changes: the batteries defaults
flip, five naming debts are paid, and command construction loses its
error-check ladder.

---

## 1. Run the codemod

From your repository root:

```bash
bash /path/to/cmdguard/docs/migration_v4_v5/codemod.sh
```

The codemod rewrites mechanically:

| Change        | Before                                                  | After                                             |
| ------------- | ------------------------------------------------------- | ------------------------------------------------- |
| Module path   | `github.com/larsartmann/cmdguard/v4`                    | `github.com/larsartmann/cmdguard/v5`              |
| Package path  | `.../v4/pkg/cmdguard/v4`                                | `.../v5/pkg/cmdguard/v5`                          |
| Import alias  | `v4 "..."` / `v4.NewCLI`                                | `v5 "..."` / `v5.NewCLI`                          |
| Metadata type | `CommandInfo`                                           | `CommandMetadata`                                 |
| Codec family  | `TypeHandler`, `RegisterTypeHandler`, `TypeHandlerFunc` | `TypeCodec`, `RegisterTypeCodec`, `TypeCodecFunc` |
| Prompts type  | `HuhRunner`                                             | `HuhPrompter`                                     |

It prints (but does not change) every call site that needs a human:

- `cli.SetConfig(...)` — **removed**. Set config fields in the struct you pass
  to `NewCLI`, or use `WithConfigValidation` for post-parse checks.
- `RegisterInScope(...)` — now generic:
  `RegisterInScope[T](parent *Scope, name string, provider func(do.Injector) (T, error))`.
  The provider returns the typed value directly; resolution uses
  `do.Invoke[T]`.
- `cli.HealthCheckWithContext(ctx)` — **removed**; `cli.HealthCheck(ctx)` was
  always the context-aware form. Rename the call. (`Scope` keeps both:
  `HealthCheck()` delegates to `HealthCheckWithContext(context.Background())`.)

Sub-module import paths are unchanged (they are unversioned:
`github.com/larsartmann/cmdguard/prompts`).

## 2. Update go.mod

```bash
go mod edit -require github.com/larsartmann/cmdguard/v5@latest
go mod tidy
```

Sub-module consumers: bump the submodule's `replace`/require pair the same
way (the `/v5` core pseudo-version is `v5.0.0-00010101000000-000000000000`
until the release tag lands).

## 3. Known traps (real incidents)

These three cost real migration time across the fleet — read before
`go build`, not after.

### Templ version floor (MVS vs `replace`)

cmdguard v5 transitively requires `github.com/a-h/templ` (via
`samber-do-auditlog`). Two failure shapes:

- **No pin:** MVS force-upgrades your templ to cmdguard's floor
  (`v0.3.1100`). Regenerate templ components if the newer codegen rejects
  them.
- **Pinned with a `replace`:** `replace` WINS over MVS, so a too-old pin
  silently downgrades templ below what auditlog needs — auditlog's generated
  templ code fails to compile with no cmdguard code involved. Align the pin
  with `go mod graph | grep a-h/templ`, and delete the `replace` once your
  floor matches.

### Nix builds need a lock step

A Nix-built consumer must refresh vendored sources after the bump, or the
build fails on stale `vendorHash`/flake lock checksums:

```bash
go mod tidy          # settle go.mod/go.sum first
nix flake update     # then refresh the lock
nix build            # vendorHash mismatch here means the order above was skipped
```

### Codemod renames are scoped to cmdguard files, with a dry run

`codemod.sh` sed-rewrites every file's IMPORT PATHS repo-wide (a
`cmdguard/v4` string anywhere IS cmdguard), but the identifier renames
(`TypeHandler` → `TypeCodec`, …) apply ONLY to files that reference
cmdguard — a consumer's OWN `TypeHandler` type, or a `v4 "..."` import
alias for a different package, is left alone (this bit two fleet repos
during the 2026-11 migrations; the scoping + `--dry-run` are the fix).

```bash
git status --short   # must be clean
git switch -c chore/cmdguard-v5-codemod
bash /path/to/cmdguard/docs/migration_v4_v5/codemod.sh --dry-run   # preview
bash /path/to/cmdguard/docs/migration_v4_v5/codemod.sh             # apply
git diff --stat      # review the diff; restore anything unexpected
```

Run it on a branch, preview with `--dry-run`, review the full diff.

## 4. Review the new defaults

v5 ships with the batteries ON. A bare `NewCLI` now means:

| Default              | Effect                                         | Opt-out                                     |
| -------------------- | ---------------------------------------------- | ------------------------------------------- |
| Strict validation    | every command needs `WithShort`                | `WithoutStrictValidation()`                 |
| Signal handling      | SIGINT/SIGTERM cancel the run                  | `WithoutSignalHandling()`                   |
| Graceful DI shutdown | services shut down on the signal path          | `WithoutGracefulShutdown()` (keeps signals) |
| Strict config keys   | unknown config keys fail the load              | `WithoutStrictConfigKeys()`                 |
| Error hints          | `Hint:` line under errors                      | `WithoutErrorHints()`                       |
| Batteries bundle     | auto-registered `doctor` (+`version` when set) | `WithoutBatteries()`                        |

Everything above stays individually opt-in/out-able; `WithBatteries()`
re-enables the whole bundle after any opt-out.

**Two semantic notes:**

- **Graceful shutdown fires on the signal path only.** A normally completing
  `Execute` no longer tears down the DI scope, so repeated
  `ExecuteWithArgs` is safe under the defaults. Explicit teardown:
  `cli.Shutdown(ctx)`.
- **Doctor reports the scope's built-in services.** Expect
  `2 passed, 0 failed` (FlagRegistry + your config) on a bare healthy CLI.

## 5. Commands without WithShort now fail

Strict validation is the default: `AddCommand` rejects commands without a
short description (`ErrMissingShort`). Add `WithShort("...")` — or opt out
with `WithoutStrictValidation()` while migrating.

## 6. Kill the construction ladder (optional, recommended)

v5 adds `Builder[T]`: one line per command, errors collected and reported
once at `Build`:

```go
b := v5.NewBuilder(cli)
v5.Add(b, "list", &ListFlags{}, listHandler, v5.WithShort("List tasks"))
v5.Add(b, "add", &AddFlags{}, addHandler, v5.WithShort("Add a task"))
child, err := v5.NewCommand("status", &DBFlags{}, statusHandler, v5.WithShort("Status"))
b.Check(err) // value consumed by the parent below; never queued twice
v5.AddParent(b, "db", "Database operations", &DBFlags{},
	v5.WithSubcommands(migrateCmd, seedCmd, statusCmd), v5.WithShort("DB"))
if err := b.Build(); err != nil {
	return err // every failure, joined; nothing attached on error
}
```

`examples/taskctl` migrated with this pattern — see its `buildCommands` for
the full before/after.

## 7. Worked example: taskctl

The flagship example moved from 14 `cmd, err := ...; if err != nil { return
err }` stanzas to a single `Builder` sequence (359 → ~310 lines in
`commands.go`, zero error-handling branches in the registration path). Its
test suite is unchanged — behavior parity is proven by the same tests
passing.

## 8. Verify

```bash
go build ./... && go test ./...
cmdguard-lint lint --dir .   # CG002 fires while you import v4 (frozen)
```
