# Fleet Findings — v5.1.0 Bump and Sweep (2026-11-11)

**Scope:** every clean cmdguard consumer across `~/projects`, following the
v5.1.0 release (tag → Release run success → proxy-verified → GitHub Latest).

## Bump results (FM1)

| Disposition | Repos |
| --- | --- |
| Bumped to v5.1.0, committed, pushed | 27 (AutoCart, InboxClean, Standup-Killer, Zlota44, code-duplicate-analyzer, crm, docs-organizer, erraudit, file-and-image-renamer, go-auto-upgrade, go-credentials-store, go-version-auto-configure, go-wizard-sdk, google-domain-detector, io-surface, lean-business-plan, media-atlas, mr-sync, nsfw-classifier, pim-sync, projects-management-automation, samber-linter, securitymd, terraform-diagrams-aggregator, timesheets, yt-history-intel, project-dependency-graph) |
| Bumped, commit blocked by the repo's OWN pre-existing docs-gates red | pingen-spike (staged `go.mod`/`go.sum`; its hook refuses unrelated commits — do NOT `--no-verify`) |
| NOT bumped (owner WIP / active concurrent session) | CV (10 dirty files), cv (forked fossil preserved on `local-fork-backup-20261111`; origin/master is 1302 commits ahead under another session and already carries the v5 API migration — it needs the v5.1.0 bump once its own `internal/ui` build break is fixed at HEAD) |
| Excluded by policy (D6 staged worktrees) | crm-stage-t13b, cv-stage-crmfix, BuildFlow-pr44 |

Vendored repos (media-atlas, mr-sync, pingen-spike, project-dependency-graph,
Zlota44, file-and-image-renamer) need `go mod vendor` (or `go work vendor`
under a `go.work`) AFTER tidy — the first sweep's five failures were all
inconsistent-vendoring, not code.

## Sweep verdicts (FM8)

`scripts/fleet-sweep.sh` — `go vet` + `golangci-lint run` per repo:
**26/26 repos vet=OK, lint issues=0.** The v5.1.0 fleet is clean.

## Consumer fix verdicts (FM2–FM7)

- **FM2 (CV portals test): STALE.** `TestLoadIdentityPortalsSection` exists
  and PASSES at CV HEAD; the recorded failure predates the current tree.
  CV's only red today is its OWN `TestKeywordLogoCoverage` (8 keywords
  without official icons) — a CV design decision, not cmdguard's.
- **FM3 (checkpoint_path defaults): FIXED** in cv — added to
  `defaultsCRMYAML` (empty = in-memory-only, per the struct doc). Verified:
  `TestDefaultsYAMLContainsEveryConfigKey` green. (Also present as the
  fleet report's "CV" finding; the code lives in the cv codebase.)
- **FM4 (pma health-probe panic): FIXED + PUSHED.** Root cause was go-health's
  `Routes.Version` mounted by `DefaultRoutes()` with an empty stamp; the fix
  clears `routes.Version` when no build stamp exists (option-gating alone was
  the wrong lever — first attempt, corrected same day).
- **FM5 (cqrs-lint offline hang): FIXED + PUSHED (go-cqrs-lite 57982df48).**
  Not a hang: the two multi-module integration tests each spent ~58s in
  go/packages type-checking, ~116s sequentially against a 120s per-binary
  timeout. A once-guarded shared `BuildContext` cuts the pair to ~24s.
- **FM6 (describe snapshots): HARNESS SHIPPED** in cmdguard
  (`testutil.AssertDescribeContract` + `-update-golden`), goldens committed
  for the two in-repo pilots (taskctl, docs-generator). Consumer adoption is
  one line per CLI.
- **FM7 (cv 20-file migration): SUPERSEDED.** cv's origin/master (another
  active session) already carries the v5 API migration; the local forked
  copy was preserved on `local-fork-backup-20261111` and NOT pushed. cv
  needs only the v5.1.0 bump once its HEAD builds again.

## Dependency verdicts

- **auditlog 0.12.1 → 0.12.3 (A38):** upstream delta is exactly one code
  commit — `b51d80a fix(live): regenerate templ output against templ
  v0.3.1100`. crm's pull was DELIBERATE (owner commit `7f101dd`, paired with
  an explicit contact-id conversion). No hidden behavior change.
- **templ pin verdict (A08): the four `replace` pins (CV, cv, nsfw-classifier,
  Zlota44) must STAY.** Evidence: without the pin, templ v0.3.1100 removes
  `templ.ResolveAttributeValue`, which templ-components v1.26.0's generated
  code calls (`svg_templ.go:57,70`) — the build breaks at GOWORK=off and in
  workspace mode. The plan's "remove pins once v5.0.1 exists" premise is
  disproven: v5.1.0 still floors templ at 0.3.1100 via samber-do-auditlog.
  Removal precondition: templ-components regenerating against ≥ v0.3.1100.

## Unresolved, upstream-blocked

- GO-2026-6859 (net/http MaxBytesReader race): fixed in go1.27.4; nixpkgs
  `go_1_27` still ships 1.27.3. Bump the flake when nixpkgs moves.
- pingen-spike commit: blocked by its own docs-gates red (pre-existing).
- cv: bump deferred until its active session greens the tree.
