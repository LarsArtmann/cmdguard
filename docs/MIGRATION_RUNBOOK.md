# Migration Runbook: any consumer → cmdguard v5

> One doc for migrating ANY repository to `github.com/larsartmann/cmdguard/v5`.
> Proven three times on 2026-11-08: BuildFlow (in-repo), erraudit
> ([PR #15](https://github.com/LarsArtmann/erraudit/pull/15)), timesheets
> ([PR #3](https://github.com/LarsArtmann/timesheets/pull/3)) — plus the
> mechanical round-trip proof (`docs/migration_v4_v5/roundtrip.sh`).

**Prereq:** `v5.0.0` is proxy-served (`go get github.com/larsartmann/cmdguard/v5@v5.0.0`).

---

## 1. Branch

```bash
git switch -c cmdguard-v5
```

Verify the tree is clean first — a dirty tree makes the "was it the codemod?"
question unanswerable.

## 2. Run the codemod

```bash
bash /path/to/cmdguard/docs/migration_v4_v5/codemod.sh
```

Mechanical rewrites: module path `/v4` → `/v5`, import aliases `v4` → `v5`,
`CommandInfo` → `CommandMetadata`, `TypeHandler`/`RegisterTypeHandler`/`TypeHandlerFunc`
→ `TypeCodec` family, `HuhRunner` → `HuhPrompter`. It PRINTS (never edits)
manual steps: `SetConfig` removal, `RegisterInScope[T]` generic call sites.

**Known after-effect:** the codemod rewrites the go.mod module PATH but leaves
the old version string (`cmdguard/v5 v4.2.1` — an invalid pairing). Pin the
version yourself next; every real migration hit this.

## 3. Pin and tidy

```bash
sed -i 's|github.com/larsartmann/cmdguard/v5 v[0-9.]*|github.com/larsartmann/cmdguard/v5 v5.0.0|' go.mod
go get github.com/larsartmann/cmdguard/v5@v5.0.0
go mod tidy
go build ./...
```

Transitive note: cmdguard v5 requires `samber-do-auditlog v0.12.1+`. If your
repo pins an older auditlog and compiles `auditlog/live`, bump to `v0.12.2+`
(templ API in `live/fragments_templ.go` moved).

## 4. Review the v5 default flips

Batteries are ON by default now. Check each against your CLI's intent:

| Default (v5)               | Opt-out                     | Decide                                     |
| -------------------------- | --------------------------- | ------------------------------------------ |
| strict validation (shorts) | `WithoutStrictValidation()` | usually keep — CG007 hygiene pays off      |
| signals cancel the run     | `WithoutSignalHandling()`   | usually keep                               |
| graceful DI shutdown       | `WithoutGracefulShutdown()` | keep; fires ONLY on the signal path        |
| strict config keys         | `WithoutStrictConfigKeys()` | keep unless a shared config has extra keys |
| hints + doctor (+ version) | `WithoutBatteries()`        | master switch for the whole set            |

Two behavioral notes: graceful shutdown no longer fires on normal completion
(repeated `ExecuteWithArgs` is safe; explicit teardown is `cli.Shutdown(ctx)`),
and every `Execute` resets flag state + config to the constructor snapshot.

## 5. Verify

```bash
go test ./... -count=1
```

If a test fails, check whether it failed on the pre-migration commit too
(`git worktree add /tmp/pre master && cd /tmp/pre && go test ...`) — do not
let a stale spec block a mechanical migration; file it separately.

## 6. Lint gate

```bash
go run github.com/larsartmann/cmdguard/lint/cmd/cmdguard-lint@v0.3.0 lint --dir .
```

Expect **0 CG002** (v4 imports are the stale-major signal). CG007/CG008
(missing short/version) surface under v5's strict defaults — fix them, don't
suppress, unless the command is intentionally bare.

**Invocation gotcha:** the directory is the `--dir` FLAG, not a positional
(`lint --dir .`), and the dogfood CLI has subcommands (`lint`, `rules`,
`doctor`, `version`). A bare `lint .` errors with "Unknown command".

## 7. Open the PR

Evidence block: build green, suite counts (before/after), lint findings
(before/after), and any pre-existing failures ruled out via worktree.
Keep the PR mechanical — behavior opt-outs are a separate decision.
