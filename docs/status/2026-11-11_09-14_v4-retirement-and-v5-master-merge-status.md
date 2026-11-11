# cmdguard — v4 Retirement, v5 Master Merge: Status & Self-Critique

**Date:** 2026-11-11 09:14 CET (Wednesday)
**Session scope:** decide whether to merge `v5-scaffold` into `master`, execute the
merge, retire the v4 major line with a final release, and leave the trunk green.
**Starting point:** `master` at `0ef84dd` (module path `/v4`), `v5-scaffold` 53
commits ahead / 0 behind, `v5.0.0` already tagged on the branch and published.
**Ending point:** `master` at `359feb3` (module path `/v5`), `release/v4` at
`d887363` with `v4.3.1` as the final v4 release.
**Operator note:** the auto-commit daemon was active throughout and committed
some of this session's changes itself (see §2.4).

---

## 0. Executive summary

The merge was the easy part: `master` was a strict ancestor of `v5-scaffold`, so
it was a clean fast-forward — no rebase, no conflicts, 72 commits and 417 files.

The hard part was that **the tree was not actually mergeable as found**. Four
independent defects had to be fixed before `master` could go green:

1. **v4 was uninstallable.** Every v4 tag since `v4.2.1` pinned the
   `flightrecorder`/`glamour`/`spinner` siblings to the zero pseudo-version
   (`v0.0.0-00010101000000-000000000000`), which no proxy can resolve.
2. **`v5-scaffold` did not build.** A dependency-bump commit had dropped four
   `go.sum` `/go.mod` checksum lines, so `go build ./...` failed outright.
3. **All seven sub-modules failed CI** with "updates to go.mod needed" after the
   same bump.
4. **The error-reference drift guard fought the markdown formatter**: dprint
   aligns tables, the guard compared bytes, so any formatting pass turned it red
   without a single fact changing.

All four are fixed, `nix run .#check-all` is green, and CI plus Sub-module Smoke
are green on `master`. One red check remains (website deploy, §4.1) and one
permanently red tag exists (`v4.3.0`, §4.2).

| Item                            | State                                      |
| ------------------------------- | ------------------------------------------ |
| `master` on v5, pushed          | DONE (`359feb3`)                           |
| v4 retired with final release   | DONE (`v4.3.1`)                            |
| Full quality gate green         | DONE (`check-all` → "All checks passed!")  |
| CI + Sub-module Smoke on master | DONE (both `success`)                      |
| Website deploy                  | **RED** (pre-existing + stale guard, §4.1) |
| `v4.3.0` tag release check      | **RED forever** (own workflow bug, §4.2)   |
| v5.1.0 release                  | NOT STARTED (42 commits past `v5.0.0`)     |

---

## 1. Why the merge was safe (evidence)

| Check               | Result                                                                                                |
| ------------------- | ----------------------------------------------------------------------------------------------------- |
| Divergence          | `v5-scaffold` 53 ahead, **0 behind** `master` → fast-forward, no merge commit                         |
| Local `v5.0.0` tag  | annotated, SSH-signed, identical object to remote (`79757cb` → `56fd69b`)                             |
| Tag ancestry        | `v5.0.0` is an ancestor of the branch                                                                 |
| Root `go.mod` at v5 | real sibling pins (`flightrecorder v0.2.0`, `glamour v0.3.0`, `spinner v0.3.0`) → consumer-resolvable |
| Consumer smoke      | `go get .../v5@v5.0.0` resolved from the proxy                                                        |
| Working tree        | clean; only `go.work`/`go.work.sum` untracked-by-design (gitignored)                                  |

The decisive argument for landing it on `master`: `ci.yml` triggers only on
`push`/`pull_request` to `master|main`, so **the branch had never run CI**. The
only place to get CI validation was `master` itself.

---

## 2. What landed

### 2.1 Topology after the session

| Ref                        | Commit    | Meaning                                          |
| -------------------------- | --------- | ------------------------------------------------ |
| `master` / `origin/master` | `359feb3` | v5 trunk (72 commits, 417 files since `0ef84dd`) |
| `v5-scaffold`              | `359feb3` | fast-forwarded to `master` (kept, not deleted)   |
| `release/v4`               | `d887363` | frozen v4 line (2 commits since `0ef84dd`)       |
| `v4.3.0`                   | `52c73af` | first installable v4 (superseded)                |
| `v4.3.1`                   | `d887363` | **final** v4 release, green release run          |
| `v5.0.0`                   | `56fd69b` | current v5 release, now GitHub "Latest"          |
| `release/v2.10`            | `1c61a0d` | untouched legacy v2 maintenance                  |

### 2.2 v4 retirement content (`release/v4`)

- `go.mod`: `flightrecorder v0.1.0`, `glamour v0.2.0`, `spinner v0.2.0` (real
  pins; local `replace` directives stay unversioned, matching the documented
  release shape).
- `pkg/cmdguard/v4/doc.go`: `Deprecated:` retirement paragraph pointing at v5
  and at the migration guide.
- `README.md`: retirement banner at the top; "API Stability" rewritten; install
  command now points at v5.
- `CHANGELOG.md`: `[4.3.0]` and `[4.3.1]` sections, `[Unreleased]` reset.
- `AGENTS.md`: status line marks the line FINAL/retired.
- `.golangci.yml`: one scoped `staticcheck` / SA1019 exclusion for
  `(benchmarks|examples|tests)/` restricted to the retirement message text —
  without it the deprecation notice produced 13 self-inflicted findings.

### 2.3 Pre-merge repairs on the v5 trunk

| Fix                                           | Why it was blocking                                                                     | Commit                         |
| --------------------------------------------- | --------------------------------------------------------------------------------------- | ------------------------------ |
| `go.sum` restore (4 `/go.mod` lines)          | `go build` failed: "missing go.sum entry" — CI would have been red on the merged master | `0ce8112`                      |
| Drift guard normalization                     | formatter vs byte-equality split brain                                                  | `0ce8112`                      |
| 7 sub-modules workspace-aware `go mod tidy`   | "updates to go.mod needed" in the smoke job                                             | `143d7fe` (+ daemon `8568abd`) |
| `release.yml` lint args `run ./...` → `./...` | every `v*` tag since v4.1.0 had a red release run                                       | `2218a19`                      |
| Docs: retirement + topology + workflow trap   | AGENTS.md/CHANGELOG had no record                                                       | `78354b2`, `359feb3`           |

### 2.4 Commits the daemon made for this session

- `8b35230` — `docs/ERROR_REFERENCE.md` regeneration + `go.sum` + an untracked
  status report written by a concurrent session.
- `8568abd` — the `glamour` tidy (generic message; not amended).

Everything else was committed deliberately with a message explaining the why.

### 2.5 Release / CI evidence

| Run           | Workflow         | Commit/tag       | Result                          |
| ------------- | ---------------- | ---------------- | ------------------------------- |
| `42980794798` | Release          | `v4.3.0`         | **failure** (own `args` bug)    |
| `42981819212` | Release          | `v4.3.1`         | **success**                     |
| `42981877056` | Website          | master `2218a19` | **failure** (stale `/v4` guard) |
| `42982137576` | CI               | master `78354b2` | success                         |
| `42982137501` | Sub-module Smoke | master `78354b2` | success                         |
| `42982409522` | CI               | master `359feb3` | success                         |
| `42982409482` | Sub-module Smoke | master `359feb3` | success                         |

`nix run .#check-all` on `359feb3`: build (8 modules), race tests (8 modules),
lint (8 modules, 0 issues), format check, tidy check (8 modules), docs/guide
sync, dependency tables — **"All checks passed!"**.

---

## 3. Verification performed

- v4 release candidate, consumer shape: copy of the tree with the `replace`
  block stripped → `go build ./...` and `go test ./...` green, `go mod tidy`
  left the real pins untouched, no `00010101` anywhere.
- v4.3.1 through the proxy: scratch module, `go get .../v4@v4.3.1`, real program
  importing the library → built and ran.
- v5 consumer: `go get .../v5@v5.0.0` resolved.
- Retired-line lint: 0 issues after the scoped exclusion.
- Sub-modules: build + tests green for all 7 with `GOWORK=off` (the CI shape).
- GitHub: `v4.3.0` and `v4.3.1` release pages exist; `v5.0.0` marked _Latest_
  (it had been shadowed by the newly created v4.3.1).

---

## 4. (d) What is genuinely broken / my mistakes

### 4.1 Website deploy is RED and has been since 2026-11-08 (pre-existing, now merged)

`website.yml` guards import paths with a hard-coded `/v4` check:

```bash
BAD=$(grep -n 'github.com/larsartmann/cmdguard"[^/]\|github.com/larsartmann/cmdguard/v3' "$file")
```

The v5 changelog page (`website/src/content/docs/changelog.mdx:63`) legitimately
quotes the historical `/v3` module path, so the guard fails. The guard was never
updated during the v5 migration, and the workflow is path-filtered to
`website/**`, so it only fires when website files change — it stayed red on the
Nov 8 master push, on Dependabot PR #24, and now on the merged commit.

Consequence: **the documentation website has not deployed since the v5 docs
landed.** This is a real user-facing regression caused by a stale gate, not by
the merge itself, but the merge made it visible on trunk.

### 4.2 `v4.3.0` has a permanently red release check — my error

I tagged `v4.3.0` after verifying code locally (build, tests, lint, fmt,
consumer-shape) but **without reading `release.yml` first**. The workflow passed
`args: run ./...` to `golangci-lint-action`, which already supplies `run`, so it
executed `golangci-lint run run ./...` and exited 7. Tag immutability means that
check is red forever.

Cost: one extra tag on the proxy (`v4.3.0`), one superseded tag (`v4.3.1`),
and a manual backfill release page. The module content is identical between the
two, and `v4.3.1` is green, but the correct move was to read the workflow before
tagging. This bug had also silently failed the `v4.1.0`, `v4.2.0`, and `v4.2.1`
release runs — I could have spotted the pattern earlier.

### 4.3 Smaller things I did not leave clean

- The daemon's `8568abd` (glamour tidy) kept a generic auto-commit message; I
  amended the `release.yml` one (`2218a19`) but not this one.
- I left `release/v4`'s **sub-module** `go.mod` files carrying the same
  pseudo-version pins the root module had. They are unpublished dev artifacts
  (the last published sub-module tags are `v0.2.0` and resolve their own pins),
  so nothing external is broken — but the frozen branch is not internally
  consistent with the hygiene standard I applied to the root module.
- `/tmp/fix-master` worktree metadata is still `prunable` (not mine; untouched).

---

## 5. (a) FULLY DONE

1. Merge decision backed by hard evidence (divergence, tags, proxy, consumer).
2. Fast-forward `master` → v5 tree, pushed (`0ef84dd..2218a19`, then docs).
3. v4 released as `v4.3.0` (pins + deprecation) and `v4.3.1` (green).
4. `release/v4` branch created, pushed, and frozen as the v4 home.
5. Root module's poison pinned to published sibling versions; verified by a
   stripped-`replace` consumer-shape build **and** a real `go get` + program run.
6. Deprecation signals: package doc, README banner, CHANGELOG, AGENTS.md.
7. Scoped SA1019 exclusion so the retirement does not fail the line's own lint.
8. `go.sum` repaired (4 missing `/go.mod` sums) so the trunk builds.
9. Error-reference drift guard made formatting-proof; doc regenerated.
10. All 7 sub-modules re-tidied and verified with `GOWORK=off`.
11. `release.yml` `args` bug fixed on both branches; v4.3.1 release run green.
12. `master` CI and Sub-module Smoke green on the final commit.
13. `nix run .#check-all` green (build/test/lint/fmt/tidy/docs/dep-tables).
14. GitHub releases: `v4.3.0` (backfilled, not latest), `v4.3.1`, `v5.0.0`
    marked _Latest_.
15. Dependency table drift in AGENTS.md fixed (v0.12.3).
16. AGENTS.md records the retirement topology and the workflow-args trap.
17. Dependabot rebase requested on all three open PRs (#21, #22, #24).
18. Temp worktree and scratch dirs removed; `v5-scaffold` fast-forwarded to
    `master` so no stale divergent branch remains.

---

## 6. (b) PARTIALLY DONE

1. **v4 line hygiene** — root fixed; sub-module pseudo pins on the frozen branch
   untouched (deliberate, but inconsistent).
2. **v4 branch gates** — build/test/lint/fmt verified; `check-all`-equivalent
   (dep tables, guide sync, tidy gate) never run on `release/v4`.
3. **Dependabot** — rebase requested, not yet observed/merged; PR #24 may still
   conflict after the `pnpm-lock.yaml` rewrite.
4. **Living docs** — AGENTS.md + CHANGELOG updated; `TODO_LIST.md`,
   `FEATURES.md`, `ROADMAP.md` untouched.
5. **pkg.go.dev** — not fetched/verified for `v4.3.1`; the deprecation banner was
   not visually confirmed.
6. **Website** — merged v5 docs are live in git but not deployed (see §4.1).
7. **CHANGELOG for internal fixes** — the `go.sum`, drift-guard, and sub-module
   sync work exists only in commit messages, not in `[Unreleased]`.

---

## 7. (c) NOT STARTED

1. `v5.1.0` — `master` is **42 commits** past `v5.0.0` (Builder mutex,
   `WithStrictPositionalArgs`/`WithArgsPassthrough`, doctor/glamour doc fixes).
   Nothing is published for them yet.
2. Website guard fix and redeploy.
3. CI coverage for `release/*` branches (the root cause of §4.2).
4. Any automated pre-tag gate (script or workflow dispatch).
5. Branch policy cleanup: delete/keep `v5-scaffold`, decide `release/v2.10`.
6. Status-report archival per `.docs-gates.json` cadence (20 unarchived / 45
   days) — now 21 unarchived.
7. `go vet` in the quality gate (not included in `check-all`).
8. Health-check coverage for the two eager DI services (samber-linter reports
   1/3 = 33 %, at baseline — no regression, no improvement).

---

## 8. (e) What we should improve (process lessons)

1. **Read the workflow before tagging.** Every release check is frozen into the
   tag. A 30-second read of `release.yml` would have saved a permanent red tag.
2. **Green CI must exist on the exact commit being tagged.** "Locally green" is
   not the same claim. Release branches had no CI at all — fix the trigger, not
   the discipline.
3. **A dependency bump is not done until the whole workspace is tidied.** One
   commit touched the root `go.mod` and left four `go.sum` lines missing plus
   seven sub-modules stale. Gate this in CI, not in a reviewer's memory.
4. **Generated artifacts and formatters must agree by contract.** The drift
   guard now normalizes table spacing; that contract belongs in AGENTS.md.
5. **Notice the pattern, not just the symptom.** Three earlier v4 tags had the
   same red release run — the failure was visible in `gh run list` from the
   start.
6. **Do not trust "the branch is fine" when the branch runs no CI.** Local
   gates covered code, not workflows.
7. **Prefer one clean final tag.** A workflow bug cost us an extra tag
   (`v4.3.0` + `v4.3.1`) on a line where one tag should have been enough.
8. **Close the loop on the website.** A red deploy is a user-facing regression;
   it should have been in scope the moment the merge triggered it.

---

## 9. (f) Up to 50 next steps

1. Fix `website.yml` import-path guard: `/v4` → `/v5`, and allow historical
   `/v3` mentions inside `changelog.mdx` (or match only import statements).
2. Re-run the Website workflow and confirm the docs site deploys.
3. Add a `Website` status note to AGENTS.md once green.
4. Add `release/**` (or "all branches") to `ci.yml` triggers so release
   branches are validated before tagging.
5. Add a `workflow_dispatch` dry-run path to `release.yml` that runs tests+lint
   without creating a release.
6. Add a repo script `scripts/pre-tag-check.sh` that reads the release workflow
   and runs `check-all`.
7. Cut `v5.1.0` (Builder concurrency fix + strict positional args options).
8. Write the `v5.1.0` CHANGELOG section from `v5.0.0..master` (42 commits).
9. Add `[Unreleased]` entries for the internal fixes landed in this merge.
10. Add a unit test for `normalizeTableSpacing` (aligned vs unaligned input).
11. Consider comparing the generated block against a formatter-normalized
    golden instead of only spacing normalization.
12. Decide `v5-scaffold`'s fate (delete local+remote, or repurpose as
    `next`).
13. Decide `release/v2.10`'s fate (retire v2 the same way?).
14. Prune `/tmp/fix-master` worktree metadata (confirm it is abandoned first).
15. Re-check Dependabot PRs #21/#22/#24 after the rebase request; merge or close.
16. Decide whether `v4`'s sub-module `go.mod` pseudo pins get a hygiene commit
    on the frozen branch.
17. Run the `check-all` equivalents (dep tables, guide sync) on `release/v4`.
18. Add a `golangci-lint` args guard: grep workflows for `args: run` in CI.
19. Audit every workflow for other duplicated-subcommand arguments.
20. Move the shared lint step into a reusable workflow to prevent divergence.
21. Fetch pkg.go.dev for `v4.3.1` and `v5.0.0`; verify the deprecation banner.
22. Update `TODO_LIST.md` with the items in this report.
23. Update `FEATURES.md` to mark v4 as retired / v5 as current.
24. Update `ROADMAP.md` with v6 candidates (middleware signature unification).
25. Archive status reports per `.docs-gates.json` cadence (now 21 unarchived).
26. Add `go vet ./...` (and `go mod verify`) to `check-all`.
27. Decide on the 31 `file-size-check` findings (split files or raise the limit
    with rationale).
28. Raise samber-linter health coverage above 33 % (config + registry services).
29. Add a health check or an explicit `allow` rationale for both eager services.
30. Document the generated-doc vs formatter contract in AGENTS.md.
31. Document the v4 retirement runbook (branch, final tag, deprecation, release
    page) for future majors.
32. Add branch protection for `master` (required checks: CI, Sub-module Smoke,
    Website when applicable).
33. Add tag protection (no force-update, no deletion) — the v4.3.0 lesson.
34. Consider `CHANGELOG`-driven release notes instead of
    `generate_release_notes`.
35. Verify the `v4.3.1` release body reads correctly (auto-generated).
36. Add an explicit "latest" policy: the newest supported major should be
    marked Latest, not the newest tag.
37. Add a README link to `release/v4` for historical v4 users.
38. Consider a `docs/DEPRECATION_POLICY.md` covering majors and sub-modules.
39. Check that the website's v4→v5 migration guide resolves after deploy.
40. Add a website banner or page stating v4 is retired.
41. Verify `cmdguard-lint` CG002 still fires on v4 imports after the merge
    (the fleet baseline moved to v5).
42. Re-run the fleet consumer check if any project still pins `/v4`.
43. Ensure `go.work.sum` is not accidentally committed (it is gitignored).
44. Confirm no `.crush/shell-output` artifacts leak into git.
45. Re-run `git town` compatibility check (branch/trunk config unchanged).
46. Consider deleting the `release/v2.10` remote branch if v2 is fully retired.
47. Add a test that sub-module `go.mod` files are tidy (currently only manual).
48. Add a CI job that builds with `GOWORK=off` explicitly to mirror consumers.
49. Re-verify the `v5.0.0` release notes claim "check-all green" (true at tag).
50. Schedule the next status report / docs-gates harvest.

---

## 10. (g) Questions I cannot answer myself

1. **Cut `v5.1.0` now?** `master` carries 42 commits past `v5.0.0`, including a
   data-race fix in `Builder` and two new public options. Publishing is
   irreversible, so I will not tag without your word.
2. **Branch policy:** should `v5-scaffold` be deleted now that it equals
   `master`, and is `release/v2.10` still maintained or should v2 be retired the
   same way v4 just was?
3. **Website scope:** fix only the stale `/v4` guard, or also add a v4
   retirement notice / make v5 the single documented line on the site?

---

## Appendix A — Commands that define the end state

```bash
git log --oneline -1 master                 # 359feb3
git log --oneline -1 release/v4             # d887363
git tag --points-at release/v4              # v4.3.1
git describe --tags --always master         # lint/v0.3.0-41-g359feb3
gh run list --branch master --limit 2       # CI + Sub-module Smoke: success
nix run .#check-all                         # All checks passed!
```

## Appendix B — Known reds

| Check                    | Since             | Cause                                                                         | Owner action                  |
| ------------------------ | ----------------- | ----------------------------------------------------------------------------- | ----------------------------- |
| Website (master)         | 2026-11-08        | hard-coded `/v4` guard trips on a historical `/v3` mention in `changelog.mdx` | fix guard, redeploy           |
| Release for tag `v4.3.0` | 2026-11-11 07:37Z | `golangci-lint run run ./...` (workflow `args`)                               | none — superseded by `v4.3.1` |
