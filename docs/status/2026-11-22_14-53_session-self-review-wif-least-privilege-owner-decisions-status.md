# Status Report — Session Self-Review: WIF Migration, Hosting Least-Privilege, Owner-Decision Batch (SU16/SU17/SU21, Badge, Push Policy, audit.html)

_Point-in-time: 2026-11-22 14:53 CET. Branch master @ `0bf6678`, synced with origin. Live site serves HEAD (deploy-meta.json matches). Battery 5/5 at last run. This report covers the session that ran 2026-11-21 ~19:00–23:25 CET and re-states what a fresh session needs._

## Self-Critique (what I forgot, what was suboptimal, what to improve)

- **I pushed a broken deploy path to master twice** (8e92bc5, 52be0a2 — CI deploy red twice) because I never pre-flighted the impersonation chain's API surface. `iamcredentials.googleapis.com` was disabled; firebase-tools masked the 403 behind its generic "have you run firebase login?" and I spent two blind retries on IAM-propagation theories before pushing a one-run debug probe that named the cause instantly. The debug probe should have been push #1, and the API pre-flight (`gcloud services list`) should have preceded ANY deploy-path push. The site itself was never affected — every failed run died before creating a release — but "master CI red for ~15 minutes" is a real cost I chose by sequencing diagnosis after deployment instead of before it.
- **I committed infrastructure changes late and lost my commit message to the daemon** — the audit.html untrack + badge-decision edits sat uncommitted for ~25 minutes while I did GCP recon; the daemon swept them into `736d067` and my authored message vanished. Known race class, known lesson; I repeated it anyway. Rule I should follow: authored commit lands within minutes of the edit, before starting any parallel work.
- **My own gcloud command flubbed twice** on provider creation: dropped `--issuer-uri` on a reorder, then hit an org-enforced attribute condition whose error message ("must reference one of the provider's claims" when NO condition was sent) is actively misleading. The `--log-http` proof (request body clean) was the right move but should have been step 1, not step 4.
- **The gcloud active account changed mid-session** (lars@helpless.ai → lartyhd@gmail.com, from a parallel session) and I ran the first GCP commands on identity-assumption instead of `gcloud config get-value account` first. It worked — both accounts are project admins — but on a deploy path I verify identity before I mutate, full stop.
- **The TTL poll script had a plumbing bug** (echoed to the background shell's stdout instead of a file), so the first measurement window produced nothing retrievable — and that window covered only a FAILED deploy anyway. The measurement I actually needed needed poll #3.
- **The 23-10 report went stale before it finished committing**: I wrote §f2 as "confirm the next push-to-master deploy is green" while the final push that confirmed it was already in flight; it went green minutes later. Minor, but reports should never describe in-flight verification as future work without a timestamp caveat.
- **Single-geo honesty gap**: the SU21 "brownout not reproducible" verdict rests on one machine (one EU PoP). The historical ~60s window could be geo-specific or era-specific. The struck row's verdict says "documented", which is honest, but a stronger verdict needed a second vantage point I didn't build.
- **WIF principal is branch-agnostic**: I mapped `attribute.repository` (and owner) but NOT `attribute.ref`, and the provider condition doesn't constrain refs. Any workflow in the repo that can obtain an OIDC token can impersonate the deploy SA on any ref. workflow_dispatch widening (deploy now runnable on demand) makes this tradeoff live. Conscious, documented, but unhardened.
- **AGENTS.md grew to 475 lines** (advisory cap 220, D13 keeps gotchas load-bearing). My +25-line WIF section is earned gotcha, but I added without trimming any prose elsewhere.

## a) FULLY DONE

- **The three owner questions from the 17-02 report were answered and executed to green:**
  1. **Badge = fleet-total** (owner pick, carried since the 13-10 report). Recorded in AGENTS.md (status line + one-fact-one-home bullet): the CI badge deliberately reports fleet-total; core-only headline numbers are deliberate tiering, not drift. No badge code change needed.
  2. **GCP authorized → executed and verified:**
     - **SU16 (WIF migration):** pool `github` + OIDC provider `cmdguard` in `lars-software` (project 132045829579) with `assertion.repository == 'LarsArtmann/cmdguard'` attribute condition; `roles/iam.workloadIdentityUser` on `cmdguard-website-deploy` bound to the repo-scoped principalSet; `iamcredentials.googleapis.com` + `sts.googleapis.com` enabled; `WIF_PROVIDER` repo secret set; website.yml deploy job on `google-github-actions/auth` (fleet hook-pinned v3 SHA) with job-level `id-token: write`. User-managed SA key `98b8b134…` (2028 expiry) **deleted**; `FIREBASE_SERVICE_ACCOUNT` secret **removed**. Verified green on two dispatch deploys AND the push-to-master path (44956333968).
     - **SU17 (least privilege):** per-site IAM condition **does not exist** for Firebase Hosting (firebasehosting absent from Google's IAM condition resource-attributes table; only `firebasehosting.admin`/`.viewer` predefined roles exist — no deployer role). Executed the achievable core: custom role `projects/lars-software/roles/cmdguardHostingDeployer` (deploy to existing sites; NO `sites.create`/`sites.delete`; also drops `firebase.clients.get` and `resourcemanager.projects.list`), swapped in for project-wide `firebasehosting.admin`. Verified by a green deploy under the new role.
     - **SU21 (CDN brownout, owner said chase):** measured, not assumed — 5s polling around a live deploy showed the new release visible at the edge within ≤6 seconds with zero 503/stale samples across two deploy events; `deploy-meta.json`'s `max-age=3600` is a browser directive (the edge invalidates on release). The historical ~60s window is not reproducible today. Documented in the struck-row verdict.
  3. **Fleet push policy = push at row completion:** pingen-spike (`67ce753`), nsfw-classifier (`037c06e0`), CV (32 commits rebased onto origin → `4f9b68019`) all pushed. crm was already in sync (its parallel session had pushed `eb32450`); crm has live parallel-session WIP and was left untouched.
  4. **taskctl-audit.html = untrack + gitignore:** writer identified (`examples/taskctl/run_test.go` auditlog export writes to CWD); `git rm --cached` + `.gitignore` entries for both `taskctl-audit.html` and `taskctl-audit.json`.
- **17-02 report landing + D3 archive ritual completed** (carried from mid-flight): 2026-10-17_00-00 report gates already clean → archived; later the 2026-10-17_00-19 report archived the same way when the new report tripped cadence again (24 > 23, exactly as predicted); battery 5/5 both times.
- **website.yml deploy job now runs on `workflow_dispatch`** — on-demand deploy verification without no-op commits; used three times this session and it is how the role swap was verified safely.
- **Strikes:** SU16, SU17, SU21 struck with honest verdicts and purged via `safe-strike --apply --purge --battery --batch` (pre-flight SAFE; battery 5/5 after).
- **Documentation:** AGENTS.md deploy-auth block rewritten (WIF shape, custom role, "no per-site IAM scope" warning) + a new "Website deploy auth (WIF)" gotcha section (attribute-condition enforcement, API requirements, `--debug` diagnosis). Cross-project lesson pushed to crush-config (`226849e`).
- **Final state:** master `0bf6678` pushed; CI green; Website push deploy green; live site serving HEAD; battery 5/5; report 23-10 landed with 19 §f items.

## b) PARTIALLY DONE

- **SU21 multi-vantage measurement:** struck as "documented" (honest per the row's own OR — "eliminate or document it"), but the measurement is single-geo (this machine's EU PoP). A second vantage or a Firebase releases-API check would make the "not reproducible" claim portable.
- **SU17 is capability-narrowed, not atomized:** the custom role's permission set is empirically verified for `firebase deploy`, not derived permission-by-permission from the CLI's calls. Edge paths (`hosting:channels:*`, `hosting:secrets:*`) will fail with the deployer role — fine for the deploy pipeline, unverified for humans.
- **WIF hardening:** repo condition exists; ref condition does not (any ref's workflow can impersonate). The dispatch widening makes the tradeoff live rather than theoretical.
- **The 23-10 report's §f2** ("confirm push-path deploy green") was confirmed by the very push that carried the report — the row is stale-on-arrival and will read as open work in any harvest.

## c) NOT STARTED

- The parked §f queue from the 17-02 and 23-10 reports: **BuildFlow lint provider pin bump to linter/v0.4.0** (fleet still lints v0.3.0 — the forgotten half of SU4, highest-value item), cmdguard v5.x release carrying the pristine-help fix, pkg.go.dev verification for lint/v0.4.0, SU44+SU46 strikes, SU40 templ-components, safe-strike CI wiring + the two safe-strike UX hardenings, pristine-help repeated-Execute terminal-view guard, fleet landed-but-unstruck sweep, nsfw gomoddirectives finding, nsfw BuildFlow step failure chase, crm contact-notes explicit re-verify, describe `help_style` decision, TTL-poll script placement, annotate.go linespec parse-time validation.
- Fleet sweep: whether OTHER repos still authenticate CI with long-lived cloud keys (cmdguard did until yesterday).

## d) TOTALLY FUCKED UP

- **Two red CI deploy runs on master** (8e92bc5, 52be0a2) from pushing before root-causing. Root cause: `iamcredentials.googleapis.com` disabled → impersonation 403 → firebase-tools generic auth error. Site never affected (failures pre-release), but the cost was self-inflicted and the diagnosis order was wrong.
- **Authored commit message lost to the daemon** (736d067 swept the audit.html/badge edits) — I left my own work uncommitted for ~25 minutes during GCP recon.
- **Provider creation burned two attempts on my own mistakes** (dropped `--issuer-uri` on a reorder; then the misleading org-enforced attribute-condition error). `--log-http` proved the request clean only after several rounds.
- **Ran the first GCP mutations under an unverified active account** (mid-session switch to lartyhd@gmail.com from a parallel session).
- **First TTL poll window wasted** by my own script plumbing (stdout not captured to a file) — and that window covered only a failed deploy, so no brownout data either way.

## e) WHAT WE SHOULD IMPROVE

- **WIF pre-flight checklist** (as a script or AGENTS checklist): APIs enabled (sts, iamcredentials) → pool/provider with repo+ref conditions → workloadIdentityUser binding → secret set → THEN push. Would have made both red runs impossible.
- **Debug-probe-first discipline:** opaque CI auth failure → first push carries the probe (echo credential env + `--debug`), never a blind retry. Instant failures are local (parsing/disabled API), not propagation.
- **Commit authored work immediately** — the daemon race is deterministic, not bad luck.
- **Identity check before privileged ops:** `gcloud config get-value account` (and `gh auth status`) before touching shared infrastructure; parallel sessions switch active accounts.
- **Multi-geo TTL probe** if SU21's verdict ever needs to be portable (second vantage, or Firebase releases API timestamps).
- **Ref-condition hardening** on the WIF provider (map `attribute.ref`, require `refs/heads/master`) if deploy-dispatch should stay master-only.
- **AGENTS.md trim pass** — prose, never gotchas (D13), before the file grows further past the advisory cap.
- **Verification items should be verified, not scheduled** — if the confirming push is already in flight, wait for it and write the report against the result.

## f) NEXT TASKS (up to 50 — brainstorm, not commitment; HARVEST fuel for `TODO_LIST.md`/`ROADMAP.md`)

Impact: Critical/High/Medium/Low · Effort: S (<30min)/M (30min–2h)/L (>2h)

<!-- check-rows:ledger -->

| #  | Task                                                                                                                                                        | Impact | Effort | Category |
| -- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ | ------ | -------- |
| 1  | Bump the BuildFlow lint provider pin to `linter/v0.4.0` (fleet still lints with v0.3.0; blank import + toolsdk pin in BuildFlow's sdk_imports)              | High   | S      | Release  |
| 2  | Cut a cmdguard v5.x release carrying the pristine-help contract fix (consumers on v5.1.0 still get ANSI in describe)                                        | High   | S      | Release  |
| 3  | Add `attribute.ref` mapping + ref condition to the WIF provider so only refs/heads/master (and optionally tags) can impersonate the deploy SA               | Medium | S      | Security |
| 4  | Decide `WIF_PROVIDER` secret → repo variable (the provider path is not sensitive; a variable is visible in workflow logs and less rotation-sensitive)       | Low    | S      | Decision |
| 5  | Fleet sweep: which other repos still authenticate CI with long-lived cloud keys — apply the WIF playbook where found                                        | High   | M      | Security |
| 6  | Verify pkg.go.dev docs for lint/v0.4.0 generated (the /fetch trigger 404'd; re-fetch until 200)                                                             | Low    | S      | Release  |
| 7  | Add pkg.go.dev verification (fetch endpoint, bounded retry) to pre-tag-check.sh or the go-release ritual                                                    | Low    | S      | Release  |
| 8  | Strike SU44 + SU46 rows with the evidence already gathered (excerpt byte-guard green; check-workflow-smoke.sh wired into pre-tag §7b)                       | Medium | S      | Docs     |
| 9  | SU40: templ-components remediate-or-baseline (owner-gated since 2026-11-20)                                                                                 | Medium | M      | Quality  |
| 10 | Pristine-help: guard live help fields against repeated-Execute double-transform (contract protected, terminal view is not)                                  | Medium | M      | Quality  |
| 11 | SKILL.md "(since vX.Y.Z)" annotation for pristine-help once the v5.x release lands (D15: one fact, one home)                                                | Low    | S      | Docs     |
| 12 | safe-strike: echo the row's NOTE cell in dry-run output so a verdict must address everything the row scopes (the SU4 half-done class)                       | Medium | S      | Tooling  |
| 13 | safe-strike: refuse whole-line strikes whose matched line is already struck (the `~~~~` double-strike garbage)                                              | Medium | S      | Tooling  |
| 14 | Wire safe-strike --purge into docs-gates CI (shell smoke test so the flag cannot silently rot)                                                              | Medium | S      | Tooling  |
| 15 | Fleet sweep: confirm no OTHER TODO rows have landed-but-unstruck work (SU44/SU46 pattern — grep TODO rows against FEATURES/CHANGELOG evidence)              | Medium | M      | Docs     |
| 16 | nsfw-classifier: resolve the `gomoddirectives` go.mod:161 templ-replace finding (pre-existing lint noise)                                                    | Low    | S      | Quality  |
| 17 | Investigate the nsfw BuildFlow step failure not yet chased (which step, why)                                                                                | Medium | S      | Process  |
| 18 | Re-verify the crm contact-notes fold explicitly (SU13's second half)                                                                                        | Low    | S      | Quality  |
| 19 | Multi-vantage SU21 check: one more measurement from a different network/geo, or Firebase releases-API timestamps, to make "not reproducible" portable      | Low    | S      | Quality  |
| 20 | Commit the TTL/deploy-window poll script into `scripts/` (file-writing version) or retire it (currently /tmp only)                                          | Low    | S      | Tooling  |
| 21 | Derive the deployer custom role permission-by-permission from firebase-tools' actual API calls (channels/secrets paths fail today; document the boundary)   | Low    | M      | Security |
| 22 | AGENTS.md prose trim pass (gotchas stay, D13) — file is at 475 lines vs 220 advisory                                                                        | Low    | S      | Docs     |
| 23 | Branch protection policy: should the Deploy Website job be a required check, and who may dispatch it? (owner call)                                          | Low    | S      | Decision |
| 24 | annotate.go: repair TSV spec — validate target line exists/range is sane at parse time (linespec is checked per file later)                                 | Low    | M      | Tooling  |
| 25 | WIF runbook: one doc (AGENTS gotcha section is the seed) covering pool/provider/binding/APIs/secret + the debug-probe diagnosis path for future migrations  | Medium | S      | Docs     |

Items 26–50 would be padding: this session surfaced nothing beyond the above.

## g) Questions I cannot figure out myself

1. **Am I now authorized to execute the parked §f queue autonomously** (starting with #1 BuildFlow pin bump and #2 the v5.x pristine-help release), or do you want per-item orders as before? The queue includes a release and a fleet-wide BuildFlow change, which is why I am asking rather than starting.
2. **WIF hardening (#3):** may I add the `attribute.ref` condition (restricts impersonation to refs/heads/master) — it touches the deploy path again and re-verification is one dispatch run — and should deploy-dispatch stay restricted to you, or is any repo-write collaborator dispatching deploys acceptable?
3. **Fleet sweep (#5):** do you want me to audit the other fleet repos for long-lived cloud keys in CI and migrate them to the WIF playbook, or is cmdguard the only repo you care to harden right now?

---

_Point-in-time snapshot. Section (f) is HARVEST input — do not let it die here. Per harness contract, not committed manually; the auto-commit daemon will pick this file up._
