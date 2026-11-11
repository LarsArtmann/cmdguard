# cmdguard v5 — Fleet Consumer Migration: Status, Self-Review, and Next 50

**Date:** 2026-11-11 08:24 CET (Wednesday)
**Session scope:** migrate every direct cmdguard v4 consumer listed in the
`Consumers of cmdguard [v5.0.0]` inventory to `github.com/larsartmann/cmdguard/v5 v5.0.0`,
verify builds/tests, fix migration fallout, update living docs and Nix pins.
**Workspace:** `~/projects` (25 repositories + go-cqrs-lite's five `cmd/*/v4` modules).
**Operator note:** all work was done in-session against the local checkouts; the
auto-commit daemon committed most of it during the run.

---

## 0. Executive summary

The migration itself is **complete and verified at the Go level**: 30 modules
build, zero `cmdguard/v4` references remain in any migrated `go.mod`/`go.sum`,
and `cmdguard-lint` reports **0 CG002** across the whole migrated set. Tests are
green everywhere the failures were mine; the remaining red tests are provably
pre-existing (A/B-verified against pre-migration commits/worktrees) or owned by
a concurrent workstream in `crm`.

Two things keep this from being a clean "done":

1. **The Nix layer is only half-migrated.** I updated `flake.nix` inputs/keys to
   `/v5` in 9 repositories but did **not** run
   `nix flake lock --update-input cmdguard`, so every `flake.lock` still pins the
   old v4 revision. A `nix build` would now see a `/v5` module key mapped to a
   v4 source tree — the opposite of the Go layer, which is consistent.
2. **`cmdguard v5.0.0`'s own go.mod is poisoned** with a leaked
   `github.com/a-h/templ v0.3.1100 // indirect` requirement (and every sub-module
   go.mod carries the same class of leaked indirects, including `flightrecorder`,
   which is supposed to have zero dependencies). That floor forces templ
   `0.3.1100` on every consumer, which breaks `templ-components` ≤ v1.26
   (`ResolveAttributeValue`/`JoinURLErrs` were removed). I worked around it with
   per-consumer pins in 4 repos; the root fix requires a `v5.0.1` release.

A mis-scoped codemod run also happened mid-session (details in §d). It was
detected, contained, and reverted; the last artifact of it (cmdguard's own v4
roundtrip fixture) was found and restored **while writing this report**.

---

## 1. What was in scope

The inventory named 32 direct consumers; 7 were already on v5. The 25 v4
consumers were:

AutoCart · code-duplicate-analyzer · crm · CV · docs-organizer ·
file-and-image-renamer · go-auto-upgrade · go-cqrs-lite (5 `cmd/*/v4`
sub-modules) · go-credentials-store · go-version-auto-configure · go-wizard-sdk ·
io-surface · lean-business-plan · media-atlas · mr-sync · nsfw-classifier ·
project-dependency-graph · projects-management-automation · samber-linter ·
securitymd · Standup-Killer · terraform-diagrams-aggregator · timesheets ·
yt-history-intel · Zlota44

Additionally migrated: **`cv` (lowercase)** — a separate GitHub repo
(`LarsArtmann/cv`) that the inventory folds into `CV` and that the accidental
codemod had half-migrated.

Already-v5 consumers touched only to repair stale Nix pins: **erraudit**,
**BuildFlow**.

---

## 2. Method (what the recipe was)

Per repository:

1. Capture baseline: `go build ./...` + `go test ./... -count=1`.
2. Run the shipped codemod (`docs/migration_v4_v5/codemod.sh`).
3. Pin `github.com/larsartmann/cmdguard/v5 v5.0.0` in `go.mod`.
4. `go mod tidy`, `go build ./...`.
5. Fix compile/runtime fallout individually.
6. Re-run tests; compare failures against the baseline (worktree A/B where needed).
7. Update living docs (`AGENTS.md`, `README.md`, `FEATURES.md`, `TODO_LIST.md`,
   `ROADMAP.md`, `STACK.md`, `CONTRIBUTING.md`, `docs/dev/*`, `docs/agents/*`)
   and functional pins (`flake.nix`, `scripts/pin-guard.sh`, `go.work`).
8. Re-vendor where the repo commits/uses `vendor/`.
9. Verify with `cmdguard-lint` (0 CG002) and binary smoke (`--help`, `describe`).

---

## a) FULLY DONE

### a.1 Dependency/API migration (all 25 + `cv`)

- Module path `…/cmdguard/v4` → `…/cmdguard/v5`, package import to
  `…/v5/pkg/cmdguard/v5`, version pinned to `v5.0.0`.
- Codemod renames applied: `CommandInfo` → `CommandMetadata`,
  `TypeHandler`/`RegisterTypeHandler`/`TypeHandlerFunc` →
  `TypeCodec`/`RegisterTypeCodec`/`TypeCodecFunc`, `HuhRunner` → `HuhPrompter`.
- Removed-API audit: no `cli.SetConfig`, no `RegisterInScope`, no
  `cli.HealthCheckWithContext` in the migrated set (all remaining
  `HealthCheckWithContext` hits are samber/do scope/injector APIs, which still
  exist).
- Verification: **zero** `cmdguard/v4` in every migrated `go.mod`/`go.sum`.

### a.2 Builds

All 30 modules pass `go build ./...`:

| Group                   | Count | Result      |
| ----------------------- | ----- | ----------- |
| Listed consumer repos   | 24    | 24/24 build |
| `cv` (lowercase, extra) | 1     | build       |
| go-cqrs-lite `cmd/*/v4` | 5     | 5/5 build   |

### a.3 Lint (migration-specific rule)

`cmdguard-lint` built from the workspace, run `lint --dir <repo>`:

- **CG002 (stale major): 0 findings** across the 24 repositories, `cv`, and
  go-cqrs-lite.
- Remaining findings are pre-existing and reported per repo in §b/§d
  (CG003 in Standup-Killer ×2 and timesheets ×1; CG005 ×1 + CG006 ×1 in
  mr-sync; CG006 in go-credentials-store ×2 and securitymd ×2; CG008 in 8 repos).

### a.4 Tests (migration-caused failures resolved)

- Full `go test ./...` sweep ran over the 24 listed repositories; `cv` and the
  five cqrs modules were tested separately.
- Remaining failures are the verified pre-existing / concurrent-WIP set in §d.
  No suite that passed at baseline fails because of the cmdguard migration —
  with the caveat in §d about `crm`'s concurrent workstream.

### a.5 Migration fallout fixed (beyond mechanical rewrites)

- **crm**
  - 29 stale `replace … => ../go-cqrs-lite/<mod>` directives repaired to
    `…/<mod>/v4` (the sibling repo restructured; crm could not even load its
    go.mod before).
  - `internal/httpapi/deals.templ`: `contactCardLabel(contacts, d.ContactID)`
    → `string(d.ContactID)`; regenerated `deals_templ.go` with templ 0.3.1100
    (a branded-ID type change made the checked-in generated code uncompilable).
- **code-duplicate-analyzer** (was unbuildable before the migration)
  - `a-h/parse` downgraded to `v0.0.0-20250122154542-74294addb73e` (a too-new
    indirect pin made templ 0.3.1100's `internal/jstokenizer` fail to compile).
  - `errors.AsType[interface{ Unwrap() error }]` → add `error` to the interface
    (Go 1.27 generic constraint).
- **Standup-Killer** (was unbuildable before the migration)
  - 39 stale `use` paths in `go.work` repaired to `…/go-cqrs-lite/<mod>/v4`.
- **go-credentials-store**
  - Both mains now map exit codes through `errorfamily.ExitCode(err)` (the v5
    `ExitCode` no longer loses the failure — it falls back to 1 — but the
    family/sysexits mapping still belongs to errorfamily); removed the stale
    `max(…,1)`-era comments and the `//cmdguard-lint:ignore CG002` suppressions;
    updated AGENTS.md's v4 gotcha list to v5 semantics; split the TODO row into
    "cmdguard migration DONE" + "F032 catalog wiring still open".
- **CV / cv**
  - Exit path simplified to `os.Exit(v5.ExitCode(execErr))` (v5 guarantees a
    non-zero code for any non-nil error; the `max(…,1)` floor is now dead code).
  - `CV/docs/generated/events.md` regenerated (`go run ./scripts/eventmap`) —
    the freshness test was red on line-number drift.
  - `docs/agents/cli-cmdguard.md`: `v4.*` API references → `v5.*` (CV + cv).
- **Zlota44**
  - Repaired a templ three-way mismatch: its own generated code needed
    templ ≥ 0.3.1100 while `templ-components` v1.26 and `go-cqrs-lite/catalog`
    need ≤ 0.3.1070. Pinned templ to 0.3.1070 (`replace`), regenerated its 8
    `.templ` files with the 0.3.1070 generator, re-vendored. Build now passes.
  - `scripts/pin-guard.sh` updated to `/v5`; **61/61 pins consistent**.
- **projects-management-automation**
  - Renamed the stale `v2` import alias (24 files) to `cmdguardv5`; gofmt
    realigned tags; removed the resolved "import alias lie" FEATURES row.
- **media-atlas / mr-sync / file-and-image-renamer / project-dependency-graph**
  - Re-vendored (`go mod vendor` / `go work vendor`) so the committed/used
    vendor trees carry cmdguard v5.
- **Nix inputs/keys updated to v5** in: crm, file-and-image-renamer, mr-sync,
  projects-management-automation, Standup-Killer,
  terraform-diagrams-aggregator (was still on **v3**), docs-organizer,
  go-cqrs-lite, erraudit, BuildFlow (4 files).
- **Living docs:** 44 files across the listed repos + 1 in `cv` updated
  (`cmdguard v4`/`/v4` → v5), historical status/planning/review snapshots
  deliberately left untouched.

### a.6 Verification tooling

- Built `cmdguard-lint` from `cmdguard/lint/cmd/cmdguard-lint`.
- Binary smoke suite: built every `main` package and ran `--help` + `describe`
  for the real CLIs (all cmdguard CLIs construct under v5's default strict
  validation; no `ErrMissingShort` surfaced).

---

## b) PARTIALLY DONE

1. **Nix layer.** `flake.nix` now asks for v5 but `flake.lock` still pins the v4
   revision in all 9 repos (crm, file-and-image-renamer, mr-sync, pma,
   Standup-Killer, terraform-diagrams-aggregator, docs-organizer, go-cqrs-lite,
   erraudit). Must run `nix flake lock --update-input cmdguard` and then
   `nix build` to recompute `vendorHash`.
2. **`vendorHash`.** Not recomputed anywhere. Several repos changed their
   dependency graph, so their FOD hashes are now wrong until rebuilt.
3. **`cv` commits.** `cv`'s 20 migration files are uncommitted; the auto-commit
   daemon did not sweep that repo. Everything else was daemon-committed.
4. **Doc semantics.** I replaced version strings, but only hand-audited the
   statements I happened to see (exit codes, required enforcement, flag reset,
   pin-guard). Other prose may still describe v4 behavior in docs that were not
   on my keyword hits (e.g., ADRs intentionally untouched, but some living docs
   could still be stale in wording).
5. **Consumer lint/format beyond the new import lines.** `gofmt` checked on
   edited files/packages only; `golangci-lint` was not run per consumer, and
   `go vet ./...` was not run as a sweep.
6. **Race testing.** Consumer suites ran without `-race` except where their own
   CI config does it; no race sweep was done.
7. **CHANGELOG entries.** None of the 25 repos gained an "Unreleased: cmdguard
   v5" line.
8. **`describe`/smoke script correctness.** My smoke script's `describe` check
   had a `pipefail` + `grep -q` bug (SIGPIPE) that produced false
   DESCRIBE-MISS lines, and it flagged non-cmdguard helper binaries. The CLIs
   were verified manually, but the script is not CI-grade.

---

## c) NOT STARTED

1. `nix flake lock --update-input cmdguard` + `nix build`/`nix flake check` for
   the 9 Nix-managed repos.
2. Releasing **cmdguard v5.0.1** to remove the leaked indirect requires
   (templ and friends) from the root and sub-module go.mod files.
3. Removing the temporary top-level `replace github.com/a-h/templ` pins
   (CV, cv, nsfw-classifier, Zlota44) after that release.
4. Migrating out-of-scope v4 holders: `bank-sync`, `DiscordSync`,
   `Code-Quality-Agent`, `ci-siblings/go-cqrs-lite`, and the staged worktrees
   (`crm-stage-t13b`, `cv-stage-crmfix`, `BuildFlow-pr44`) — or archiving them.
5. Updating `cmdguard/docs/MIGRATION_v4_v5.md` with the traps discovered today
   (templ floor, Nix lock, codemod's unconditional identifier renames).
6. A fleet-wide CI-able verification script (what I ran ad hoc).
7. Consumer-side `golangci-lint`/`go vet`/`-race` sweeps.
8. CHANGELOG/docs of the "what changed for me" kind in each consumer.

---

## d) TOTALLY FUCKED UP (honest list)

1. **Mis-scoped codemod run.** My per-repo script took a repository _name_ and
   did `cd /home/lars/projects/$repo`. When I invoked it with an empty argument
   for the cqrs sub-modules, it `cd`'d to `/home/lars/projects` and ran the
   codemod across the **entire projects tree** for ~20–25 seconds before I
   killed it. Detection and containment:
   - Signature scan (`cmdguard/v5 v<4.x>` in go.mod) found exactly one
     half-migrated repo: **BuildFlow-pr44**. Its 8 touched files were reverted
     with a reverse patch of exactly that diff (the other 22 dirty files there
     were pre-existing and untouched).
   - **`cv`** was the only other repo the run partially rewrote; I completed
     its migration properly rather than reverting.
   - **`cmdguard`'s own `docs/migration_v4_v5/roundtrip/_fixture/main.go`** was
     also rewritten (it must stay a v4 fixture). I found it only while writing
     this report and restored it.
   - No other repository showed migrated paths or collateral identifier
     renames (`CommandMetadata`/`TypeCodec`/`HuhPrompter` scans were clean
     outside cmdguard consumers; the only `CommandMetadata` hits were CV's own
     pre-existing domain type and media-atlas's regenerated vendor).
   - Not assessed: `BuildFlow.vendor.bak` (a non-git dead backup directory).
     Root cause: no argument validation in my wrapper script. This was avoidable
     and is the single worst thing in the session.
2. **Codemod false positive inside a migrated repo.** Running the codemod over
   `go-cqrs-lite/cmd/cqrs-lint/v4` renamed that tool's **own**
   `analyzer.CommandInfo` type to `CommandMetadata`, which broke its API-surface
   golden. Found by `api-stability`, reverted (12 references), then re-verified.
   The shipped codemod should scope identifier renames to files that import
   cmdguard.
3. **No Nix verification at all.** I edited functional Nix pins without locking
   or building. That is a "claims not verified by the delivering layer" mistake.
4. **Trusting a bad smoke harness.** The `describe` false negatives (§b.8) meant
   I briefly believed several CLIs lacked `describe`; manual checks disproved it.
5. **Baseline truncation.** My first baseline script piped each test run through
   `tail -60`, so package results above the cut were lost. That cost time when
   classifying nsfw/crm/CV failures and is why I later re-verified with
   worktrees.
6. **Let parallel test runs muddy A/B conclusions.** Running heavy suites across
   7 repos at once produced contention-shaped failures (`crm` app BDD, CV
   chatstore) that I initially could not classify. Serial reruns resolved it.
7. **Almost mis-attributed `crm`'s failures.** The `internal/app` BDD suite is
   currently being edited by a concurrent workstream (contact-notes projection,
   `internal/app/projections.go` +67 lines, new `contact_notes_test.go`). Its 2
   failing specs reference that new code; they are not migration fallout.

**Pre-existing breakage confirmed by A/B (not mine, left or fixed as noted):**

| Repo                           | Failure                                                               | Baseline?                     | Status                                 |
| ------------------------------ | --------------------------------------------------------------------- | ----------------------------- | -------------------------------------- |
| CV                             | `TestLoadIdentityPortalsSection`                                      | fails at HEAD                 | left (pre-existing)                    |
| CV                             | render goldens (`seanthegeek` vs `crowdsecurity` order)               | fails at HEAD                 | left (non-deterministic fixture order) |
| CV                             | `TestDefaultsYAMLContainsEveryConfigKey` (`crm.checkpoint_path`)      | fails at HEAD                 | left                                   |
| nsfw-classifier                | `TestContainerHealthCheckJoinsAllFailures` (`oops.Join` single cause) | fails at pre-migration commit | left                                   |
| projects-management-automation | health-probe panics (`Routes.Version` without `WithVersion`) ×2       | fails at pre-migration commit | left                                   |
| go-cqrs-lite                   | `doc-check` path tests (README/recipes resolved from `cmd/`)          | fails at pre-migration commit | left                                   |
| go-cqrs-lite                   | 2 `cqrs-lint` tests hang in `go/packages` (network/toolchain)         | hangs at pre-migration commit | left                                   |
| crm                            | `internal/httpapi` test package: `undefined: createDeal`              | helper never existed at HEAD  | left                                   |
| code-duplicate-analyzer        | templ/parse build failure                                             | fails at HEAD                 | **fixed**                              |
| Standup-Killer                 | broken `go.work` (39 paths)                                           | fails at HEAD                 | **fixed**                              |
| mr-sync                        | inconsistent vendoring                                                | fails at HEAD                 | **fixed**                              |
| Zlota44                        | templ runtime/codegen mismatch + vendoring                            | fails at HEAD                 | **fixed**                              |
| crm                            | stale go-cqrs-lite replaces (29)                                      | fails at HEAD                 | **fixed**                              |

---

## e) WHAT WE SHOULD IMPROVE

### Process

1. **Every fan-out script must validate its arguments** (`[ -d "$dir" ]`,
   refuse empty/`/`/`$HOME`) and print the resolved absolute directory before
   mutating anything. A dry-run mode (`--print`) should be mandatory.
2. **Never run a tree-wide codemod from a parent directory.** Operate on an
   explicit list of repo roots; assert `go.mod` exists at each root and that
   the module path matches the target repo.
3. **Scope semantic renames by import.** `CommandInfo`/`TypeHandler` renames are
   only safe in files that import cmdguard. Make those rewrites conditional (or
   emit them as suggestions).
4. **Baseline capture must be lossless** — write full logs to files, only tail
   for the human summary; record build/test/lint/`go.mod` hash per repo.
5. **A/B every failure** with a worktree at the pre-migration commit before
   classifying it as pre-existing (this session proved its value 4×).
6. **Serial verification for stateful suites**, parallel only for independent
   builds.
7. **Verify the delivering layer**: after touching `flake.nix`, the lock file
   and a `nix build` are part of the change, not a follow-up.
8. **Check concurrency**: several repos (`crm`, CV) have other active
   workstreams; check `git status`/recent mtimes for edits that are not yours
   before blaming the migration.

### Library/ecosystem

9. `cmdguard v5.0.0`'s go.mod must not carry indirect requires its packages do
   not import. The leaked templ floor (and the dozens of leaked indirects in
   `flightrecorder`, a zero-dependency module) is an ecosystem-wide footgun.
10. The codemod should learn the two traps found today: Nix lock/vendorHash
    steps and templ/compatible-dependency pinning.
11. `cmdguard-lint` CG006 has false positives on `fmt.Fprintf(os.Stderr, "…: %v",
    err)` for a `newApp()` error (credentials-store lines 53/65). Worth an
    upstream rule refinement.

---

## f) Up to 50 things to do next (ordered by impact)

1. Run `nix flake lock --update-input cmdguard` in the 9 Nix repos, then commit
   the lock files.
2. Run `nix build` (or the repo's default package) in those repos; fix
   `vendorHash` mismatches.
3. Run `nix flake check` in those repos.
4. Release **cmdguard v5.0.1** with the leaked indirect requires removed from
   the root and all 7 sub-module go.mod files.
5. Bump the 25 consumers + `cv` to v5.0.1.
6. Remove the four `replace github.com/a-h/templ => v0.3.1070` pins after v5.0.1.
7. Decide the templ-components upgrade path (all to v1.27.0 + regenerate)
   instead of pinning templ forever.
8. Add the templ-pair biconditional (templ ≥ 0.1100 ⟺ templ-components ≥ 1.27.0)
   to each repo's bump tooling / AGENTS.
9. Fix crm's missing `createDeal` test helper so `internal/httpapi` tests
   compile.
10. Verify the crm contact-notes fold nil-pointer after the concurrent
    workstream lands (its BDD suite is red for that reason).
11. Fix CV `TestLoadIdentityPortalsSection`.
12. Make CV's render goldens deterministic (fixture ordering) or regenerate them
    deliberately.
13. Fix CV `crm.checkpoint_path` missing defaultsYAML entry.
14. Fix nsfw `TestContainerHealthCheckJoinsAllFailures` (single-failure join
    shape).
15. Fix pma's health-probe panics (`Routes.Version` without `WithVersion`).
16. Fix go-cqrs-lite `doc-check` path resolution (module lives at `cmd/…/v4`).
17. Fix/skip the two hanging `cqrs-lint` `go/packages` tests in offline CI.
18. Commit `cv`'s migration (20 files).
19. Run `golangci-lint run ./...` in all 26 repos; triage findings.
20. Run `go vet ./...` in all 26 repos.
21. Run `go test -race ./...` in all 26 repos.
22. Add "Unreleased — cmdguard v5 migration" lines to each CHANGELOG.
23. Add a short "cmdguard v5 notes" section to each consumer AGENTS.md
    (batteries defaults, runtime `required`, per-Execute flag reset, exit-code
    contract).
24. Audit remaining v4 mentions in ADRs/historical docs; annotate rather than
    rewrite.
25. Update cmdguard `docs/MIGRATION_v4_v5.md` with the templ trap, Nix lock step,
    and the codemod's unconditional-rename warning.
26. Harden `codemod.sh`: require a target directory, refuse unsafe roots, add
    `--dry-run`.
27. Scope the `CommandInfo`/`TypeHandler` renames to cmdguard-importing files.
28. Add a fleet verification script to cmdguard (build/test/lint/`go.mod` pin
    per consumer) and wire it into CI.
29. Add a CG002 gate to each consumer CI (io-surface TODO 63 is blocked on the
    lint module's replace directives).
30. Migrate `bank-sync`, `DiscordSync`, `Code-Quality-Agent`,
    `ci-siblings/go-cqrs-lite` or explicitly archive them.
31. Decide the fate of staged worktrees: `crm-stage-t13b`, `cv-stage-crmfix`,
    `BuildFlow-pr44`.
32. Reconcile `cv` vs `CV` (two repos, same codebase) — see question 1.
33. Add a `GOWORK=off go mod tidy` drift check to go-cqrs-lite CI (today's
    stale go.sums prove it is needed).
34. Investigate why `go mod tidy` pulled auditlog 0.12.1 → 0.12.3 in crm and
    whether that changed runtime behavior.
35. Refine cmdguard-lint CG006 to not flag non-`Execute` error prints.
36. Fix the CG003 findings (panic on constructor error) in Standup-Killer and
    timesheets.
37. Fix CG008 (missing `WithCLIVersion`) in AutoCart, code-duplicate-analyzer,
    docs-organizer, media-atlas, pma, samber-linter, Zlota44, go-cqrs-lite.
38. Remove the now-no-op `WithBatteries()` call in io-surface.
39. Replace remaining `AddCommand` ladders with `Builder` where it reduces
    error-branch noise (optional, per repo).
40. Evaluate `NewApp` adoption for repos still doing `NewCLI` + registration.
41. Fix my smoke harness (pipefail/grep and "skip non-cmdguard mains") and make
    `--help`/`describe` a CI test per consumer.
42. Maintain a single fleet inventory (repo → cmdguard version → sub-modules
    used) as the source of truth.
43. Add Renovate/Dependabot coverage for Go major bumps (Zlota44 documented the
    gap) or a scheduled manual sweep.
44. Pin Nix flake inputs to tags, not branches (`go-cqrs-lite`/pma use `master`).
45. Re-run `docs/migration_v4_v5/roundtrip.sh` to prove the codemod still works
    (the fixture is restored; this is now a real regression test).
46. Move/rename the roundtrip fixture so `find`-based scripts cannot reach it
    (underscore dirs are ignored by the Go toolchain, not by `find`).
47. Record today's lessons in cmdguard's AGENTS.md (templ floor; codemod scope;
    Nix lock as part of a flake change).
48. Verify `cmdguard/schema` and other sub-modules are not similarly
    over-requireing after the v5.0.1 tidy.
49. Add a `describe`-contract snapshot test per consumer (agents depend on it).
50. Schedule a follow-up fleet status report after the Nix layer is finished so
    the "done" claim is backed by `nix build`, not only `go build`.

---

## g) Questions I cannot answer myself

1. **`cv` vs `CV`:** both are real GitHub repos with the same codebase; I
   migrated both. Which is canonical, and should the other be archived (or kept
   as a staging clone)?
2. **cmdguard v5.0.1:** may I cut and tag the patch release that removes the
   leaked indirect requires (templ et al.) and then bump all consumers, or do
   you want to keep the consumer-side `replace` pins until a scheduled release?
3. **Nix lock updates:** may I run `nix flake lock --update-input cmdguard` +
   `nix build` in the 9 Nix-tracked consumers (network + possibly other input
   updates), or does your deployment process own the lock files?

---

## Appendix — evidence

- Pre-migration baselines: `/tmp/mig/baseline/*.log` (24 repos).
- Migration logs: `/tmp/mig/x-*.log`.
- Post-migration test sweep: `/tmp/mig/final-tests.log` (24 repos).
- Binary smoke: `/tmp/mig/smoke/*.log`.
- Linter: `/tmp/cmdguard-lint lint --dir <repo> --output text --no-color`
  (0 CG002 for all).
- Build sweep: 30/30 `go build ./...` OK.
- `Zlota44/scripts/pin-guard.sh`: 61/61 pins consistent with go.mod.
- Restored artifact: `cmdguard/docs/migration_v4_v5/roundtrip/_fixture/main.go`
  back to the v4 shape (the roundtrip script requires it).
