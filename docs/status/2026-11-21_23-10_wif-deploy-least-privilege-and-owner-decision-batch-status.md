# Status Report — WIF Deploy Migration, Hosting Least-Privilege, Owner Decision Batch

_Point-in-time: 2026-11-21 ~23:15. Branch master. The three open §g questions from the 17-02 report were answered by the owner and executed to green this session._

## Self-Critique (what I forgot, what was suboptimal, what to improve)

- Two master pushes (8e92bc5, 52be0a2) carried a deploy job that FAILED on CI before the root cause (iamcredentials API disabled) was found. The site itself was never affected (each failed run died before creating a release), but master CI was red twice. The failing-loud-then-diagnosing loop cost ~15 minutes that a pre-flight API checklist (`gcloud services list` before touching the deploy path) would have saved.
- The first OIDC provider creation failed on a missing `--issuer-uri` and then — more interestingly — on an org-enforced attribute condition that the API reports with a misleading message ("must reference one of the provider's claims" when NO condition was sent). `--log-http` proved the request was clean; the empirical fix (send the repo condition) is now an AGENTS.md gotcha.
- The TTL poll script wrote to stdout captured by the background shell, not to a file — first "log empty" surprise was my own plumbing, not the CDN.
- firebase-tools' generic "Failed to authenticate" hides the real ADC error; `--debug` surfaced it in one run. The debug probe should have been the FIRST push, not the third.

## a) FULLY DONE

- **D3 archive ritual (carried from mid-flight):** 2026-10-17_00-00 report gates were already clean (resolved in a prior pass); `git mv` to `docs/status/archived/`; battery 5/5 with cadence back to 23/23; committed c43d112 and pushed; CI green.
- **Owner questions answered and executed (§g of 17-02):**
  1. **Badge = fleet-total** (owner pick). Recorded in AGENTS.md (status line + one-fact-one-home bullet): the CI badge deliberately reports fleet-total; core-only headline numbers are tiering, not drift. No badge code change needed.
  2. **GCP authorized + executed:**
     - SU16 WIF: pool `github` + OIDC provider `cmdguard` (project lars-software / 132045829579) with `assertion.repository == 'LarsArtmann/cmdguard'` attribute condition; `roles/iam.workloadIdentityUser` on `cmdguard-website-deploy` for the repo-scoped principalSet; `iamcredentials.googleapis.com` + `sts.googleapis.com` enabled; `WIF_PROVIDER` secret set; website.yml deploy job switched to `google-github-actions/auth` (hook-pinned v3 SHA) with `id-token: write`. **Verified by a live dispatch deploy: green.** SA user-managed key `98b8b134…` deleted; `FIREBASE_SERVICE_ACCOUNT` secret removed (only `WIF_PROVIDER` remains).
     - SU17 least-privilege: **per-site IAM condition is impossible** — firebasehosting is absent from Google's IAM condition resource-attributes table and only admin/viewer predefined roles exist. Executed the achievable core instead: custom role `projects/lars-software/roles/cmdguardHostingDeployer` (deploy to existing sites, NO `sites.create`/`sites.delete` — the admin role's damage vectors), swapped in for project-wide `firebasehosting.admin`. **Verified by a second dispatch deploy: green.**
     - SU21 chase: measured, not guessed — 5s polling around a live deploy showed the new release visible at the edge within ≤6 seconds, zero 503/stale samples; `deploy-meta.json`'s `max-age=3600` is a browser directive (the edge invalidates on release). The historical ~60s brownout is not reproducible today; documented in the struck row verdict.
  3. **Fleet push policy = push at row completion:** pingen-spike (67ce753), nsfw-classifier (037c06e0), CV (32 commits rebased onto origin, 4f9b68019) all pushed. crm was already in sync (its parallel session pushed eb32450); crm has live parallel-session WIP — left untouched.
  4. **taskctl-audit.html = untrack + gitignore:** writer identified (`examples/taskctl/run_test.go` auditlog export writes to CWD); `git rm --cached` + `.gitignore` entries for `taskctl-audit.html` and `taskctl-audit.json`.
- **website.yml:** deploy job now also runs on `workflow_dispatch` (on-demand deploy verification without no-op commits) — used three times this session.
- **Strikes:** SU16, SU17, SU21 struck with verdicts and purged into CHANGELOG via `safe-strike --apply --purge --battery --batch` (pre-flight SAFE; battery 5/5 after).
- **AGENTS.md:** deploy auth block rewritten (WIF shape, custom role, no-per-site-IAM warning) + a new "Website deploy auth (WIF)" gotcha section (attribute-condition enforcement, API requirements, `--debug` diagnosis).

## b) PARTIALLY DONE

- Nothing — everything started reached green or an honest verdict.

## c) NOT STARTED

- The remaining §f queue of the 17-02 report (BuildFlow pin bump, v5.x pristine-help release, pkg.go.dev verification, SU44/SU46 strikes, safe-strike CI wiring, …) — untouched this session; carried below.

## d) TOTALLY FUCKED UP

- Two red CI deploy runs on master (8e92bc5, 52be0a2) before enabling `iamcredentials.googleapis.com`. Root-caused with a one-run probe (echo GAC + `--debug`); site availability never at risk. Lesson recorded: enumerate the API surface BEFORE touching the deploy path.
- The daemon swept the audit.html/badge-decision working-tree changes into 736d067 before my authored commit landed (commit message lost; content intact and pushed as part of a0bc9be's ancestry). Known race class.

## e) WHAT WE SHOULD IMPROVE

- WIF pre-flight checklist (pools/providers/APIs/bindings/secret) as a script or AGENTS checklist before any deploy-auth change — would have made both red runs impossible.
- Keep the `--debug` probe pattern: when CI auth fails opaquely, push the probe FIRST; one run beats two blind retries.
- The TTL poll could live in `scripts/` (deploy-window measurement on demand) instead of /tmp — low priority.

## f) NEXT TASKS (up to 50 — brainstorm, not commitment; HARVEST fuel for `TODO_LIST.md`/`ROADMAP.md`)

Impact: Critical/High/Medium/Low · Effort: S (<30min)/M (30min–2h)/L (>2h)

<!-- check-rows:ledger -->

| #  | Task                                                                                                                                                                | Impact | Effort | Category |
| -- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ | ------ | -------- |
| 1  | Bump the BuildFlow lint provider pin to `linter/v0.4.0` (fleet still lints with v0.3.0; blank import + toolsdk pin in BuildFlow's sdk_imports)                      | High   | S      | Release  |
| 2  | Confirm the next push-to-master deploy is green under WIF + custom role (the two green runs were dispatches; the push path is identical but deserves one green run) | Low    | S      | Verify   |
| 3  | Cut a cmdguard v5.x release carrying the pristine-help contract fix (consumers on v5.1.0 still get ANSI in describe)                                                | High   | S      | Release  |
| 4  | Verify pkg.go.dev docs for lint/v0.4.0 generated (the /fetch trigger 404'd; re-fetch until 200)                                                                     | Low    | S      | Release  |
| 5  | Strike SU44 + SU46 rows with the evidence already gathered (excerpt byte-guard green; check-workflow-smoke.sh wired into pre-tag §7b)                               | Medium | S      | Docs     |
| 6  | safe-strike: echo the row's NOTE cell in dry-run output so a verdict must address everything the row scopes (the SU4 half-done class)                               | Medium | S      | Tooling  |
| 7  | safe-strike: refuse whole-line strikes whose matched line is already struck (the `~~~~` double-strike garbage)                                                      | Medium | S      | Tooling  |
| 8  | Wire safe-strike --purge into docs-gates CI (shell smoke test so the flag cannot silently rot)                                                                      | Medium | S      | Tooling  |
| 9  | SU40: templ-components remediate-or-baseline (owner-gated since 2026-11-20)                                                                                         | Medium | M      | Quality  |
| 10 | Pristine-help: guard live help fields against repeated-Execute double-transform (contract protected, terminal view is not)                                          | Medium | M      | Quality  |
| 11 | SKILL.md "(since vX.Y.Z)" annotation for pristine-help once the v5.x release lands (D15: one fact, one home)                                                        | Low    | S      | Docs     |
| 12 | Fleet sweep: confirm no OTHER TODO rows have landed-but-unstruck work (SU44/SU46 pattern)                                                                           | Medium | M      | Docs     |
| 13 | nsfw-classifier: resolve the `gomoddirectives` go.mod:161 templ-replace finding (pre-existing lint noise)                                                           | Low    | S      | Quality  |
| 14 | Investigate the nsfw BuildFlow step failure not yet chased (which step, why)                                                                                        | Medium | S      | Process  |
| 15 | Re-verify the crm contact-notes fold explicitly (SU13's second half)                                                                                                | Low    | S      | Quality  |
| 16 | Add pkg.go.dev verification (fetch endpoint, bounded retry) to pre-tag-check.sh or the go-release ritual                                                            | Low    | S      | Release  |
| 17 | Decide whether `describe` should carry an explicit `help_style: "plain"` marker (additive schema field, if wanted)                                                  | Low    | S      | Decision |
| 18 | Commit the TTL/deploy-window poll script into `scripts/` or retire it (currently /tmp only)                                                                         | Low    | S      | Tooling  |
| 19 | annotate.go: repair TSV spec — validate target line exists/range is sane at parse time (linespec is checked per file later)                                         | Low    | M      | Tooling  |

Items 20–50 would be padding: this session surfaced nothing beyond the above.

## g) Questions I cannot figure out myself

None blocking — all three open questions from the 17-02 report were answered by the owner this session and executed (badge = fleet-total; GCP authorized → SU16/SU17 done, SU21 chased-to-documentation; push policy active; audit.html untracked). New surface decisions (e.g. promoting `WIF_PROVIDER` from secret to repo variable, or branch protection requiring the deploy job) are parked as §f items and not worth an owner round-trip.

---

_Point-in-time snapshot. Section (f) is HARVEST input — do not let it die here. Per harness contract, not committed manually; the auto-commit daemon will pick this file up._
