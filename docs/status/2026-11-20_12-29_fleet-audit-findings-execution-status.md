# Status: fleet-audit 10-finding execution (docs / marketing / AI-skills)

**Scope of this report.** Verified in this session: all report-claim replays
(`rg` per the audit's appendix), live `taskctl` runs (describe, typo
suggestion, failure exit), a real JSON error envelope captured from a scratch
v5 CLI, the F7 verify-block commands executed verbatim, website build + astro
check + html-validate + internal-link check, `scripts/check-guide-sync.sh`,
docs-gates battery (per-gate), `nix fmt` (0 changes), and `go env
GOEXPERIMENT`. NOT verified: the new CI badge steps have never executed in a
real CI run (no push happened this session); the live deployed website was not
re-smoked after my content edits; the sibling repos' remote state was checked
only via `@{u}..HEAD` counts (0 unpushed each) — not by fetching the rendered
remotes; Codecov's backend state (why "unknown") is inferred, not diagnosed
server-side. Report written 2026-11-20 12:29 CET; session ran ~10:55-12:35 CET.

Session task: execute all 10 findings from
`docs/reviews/2026-11-17_16-14_cmdguard-docs-marketing-ai-skills-41-reviews.html`.

## a) FULLY DONE

1. **F1 — validate-rule registry documented.** Canonical 10-rule table (rule /
   parameter / checks / failure sentinel) plus the type-aware `min`/`max`
   rewriting story, written from `flags_validate.go:57-67` + `flag_bounds.go`
   (read in full first): README "Validate rules" section, website
   `type-safe-flags.mdx` ("Validation Rules" replaced the prose half-story),
   SKILL.md tag row now points at both instead of restating.
2. **F2 — sentinel-count contradiction closed.** SKILL.md:43 ("59 exported
   sentinel error vars") → "62 documented sentinel errors" with the derivation
   command in the Source column (`rg -c '^\|`[a-z_]+`\ +\|`Err'
   docs/ERROR_REFERENCE.md` = 62, matching the generated header); taskctl count
   corrected 86→87. Both SKILL.md mentions now agree (verified by rg).
3. **F3 — tag-vocabulary drift closed.** README Flag Tags Reference is the
   13-tag superset (adds `prompt`, `secret`); SKILL.md line-36 "Full tag set"
   fixed (was 12, missing `secret`); SKILL.md deeper-truth section names the
   README tables as canonical.
4. **F4 — released-vs-pinned gating.** "(since v5.1.0)" annotations on
   `WithStrictPositionalArgs`/`WithArgsPassthrough` in SKILL.md (incl. "do not
   exist at a v5.0.0 pin"), README features + Command Options tables, website
   shell-completion table+prose, api-reference.mdx.
5. **F5 — hand-maintained coverage claims replaced.** ci.yml now GENERATES the
   coverage SVG (`go tool cover -func` total → shields-style SVG) and commits
   it to `.github/badges/coverage.svg` on master pushes (`[skip ci]`,
   `pull --rebase` before push, `fetch-depth: 0`); README badge + Verified row
   link CI instead of restating a number. YAML parses; awk extraction verified
   against a real profile; SVG XML validity checked locally. Codecov badge
   abandoned deliberately: uploads succeed (OIDC, `conclusion=success` in run
   44135088092) but the badge serves "unknown" in every URL variant — project
   not activated on Codecov's side; recorded in AGENTS.md.
6. **F6 — "Replace your hand-rolled…" README section.** 7 verified rows
   (exit-code ladder→`ExecuteAndExit`, version command→`WithCLIVersion`,
   self-check→auto `doctor`, typo suggestions (real `taskctl lst`→`list`
   output), SIGINT teardown→`WithGracefulShutdown` (default), config loader
   precedence→`WithConfigFile`, usage-on-error→`SilenceUsage` default).
7. **F7 — "verify before trusting" block under the trinity table.** 4 runnable
   lines, each executed verbatim and passing (zero-`Must` rg exits 1 with the
   `!*_test.go` glob; `TestAddCommands_AtomicOnDuplicate` passes; samber/do in
   go.mod:85; taskctl describe prints the real contract). Env-portability fix
   applied afterwards: `GOEXPERIMENT=jsonv2` prefixes on the two go commands
   (`go env GOEXPERIMENT` = jsonv2 comes from the dev environment, not the
   toolchain default — a fresh machine would have failed without the prefix).
8. **F8 — cross-skill pollination.** cmdguard SKILL.md gains a "Works with"
   section (cqrs service + cmdguard ops CLI, go-aichat probe CLI,
   stack-selection routing); named cmdguard recipes added to go-cqrs-lite,
   cqrs-htmx, go-aichat SKILL.mds (in their repos, committed AND pushed — 0
   unpushed each). Fleet surfaces mentioning cmdguard: 3 → 7.
9. **F9 — agent loop walked end-to-end.** machine-contract.mdx "Driving a CLI
   from an agent" now shows: discover (describe) → invoke with `--output json`
   (`WithOutputFormat(output.FormatJSON)` wiring documented) → the REAL JSON
   envelope (captured bytes from a scratch CLI: `required_flag` + hint + code
   1. → an `error_code` recovery-decision table (all 5 codes verified present
      in the generated ERROR_REFERENCE) → human `Hint:` line + verified
      suggestion → `schema.FromCLI` tooling cross-link.
10. **F10 — headline tiering decision recorded.** AGENTS.md "Adoption-surface
    conventions (2026-11-20)": D15 (README=why-switch, website=what-it-is, do
    not converge) + the one-fact-one-home surface map (tags/validate rules →
    README; sentinel counts → generated ERROR_REFERENCE only; coverage →
    CI-committed badge; versions → since-annotations sourced from CHANGELOG).
    Also fixed AGENTS.md's docs-gates pin staleness (v0.3.0 → v0.3.1, the real
    ci.yml pin).
11. **Supporting work:** CHANGELOG [Unreleased] entry (pre-flighted SAFE);
    TODO_LIST SU23 (cadence-archive routing, pre-flighted SAFE); AGENTS.md
    decision block; docs-gates battery re-run green on all living-doc gates.

## b) PARTIALLY DONE

1. **F8 breadth.** 3 sibling skills done; sibling READMEs untouched; other
   plausible skills (samber-do-best-practices, templ-components, …) not
   evaluated — deferred to reciprocity pass 2 (SU30).
2. **F5 end-to-end.** Badge steps never executed in real CI (nothing was
   pushed this session); the badge SVG file does not exist yet, so the README
   image 404s until the first master CI run; push-vs-auto-commit-daemon race
   is designed for (pull --rebase) but unproven. Codecov step left in place,
   now redundant — removal/activation blocked on owner question (§g Q1).
3. **F9 depth.** The describe JSON in machine-contract.mdx is still the
   fictional "myapp" shape (real taskctl bytes were captured but not pasted);
   the optional MCP-facing recipe was not written (SU29).
4. **F10 cross-links.** The tiering decision text says each surface links the
   other — true via pre-existing links (README header → website; website hero
   → repo), but no NEW explicit link was added.
5. **Surface parity.** The website does not mirror the new "Replace your
   hand-rolled" / "How to verify" sections, and features.ts carries no
   since-column — README is ahead of the site (SU28).
6. **Verification depth.** Only fmt + targeted tests ran; the full
   build/test/lint battery did not (docs + one workflow file were the only
   functional changes — SU35).

## c) NOT STARTED

1. Drift guards for the new tables: nothing fails CI if someone edits the
   README validate-rules table or tag table away from source (SU25/SU26).
2. Build-time include so README and the website guide share ONE
   validate-rules source (the audit report's own "cross-surface plumbing"
   proposal; today the registry exists as two hand-kept copies, mitigated
   only by check-guide-sync checking API symbols, not table content) (SU27).
3. SU22 + SU23 archive batches (~190 prose verdicts total) — the docs-gates
   cadence gate remains red on master for this reason alone (pre-existing,
   documented as expected in AGENTS.md; now routed).
4. Quickstart/tutorial since-audit; AGENTS.md 89.3% stamps decision;
   resolution banner on the audit report itself; full buildflow battery
   (SU33-SU36).
5. Everything in §f items 15+ (ROADMAP umbrella — candidates, not commitments).

## d) TOTALLY FUCKED UP

1. **README multiedit double-anchor bug (worst of session).** Two edits in one
   multiedit anchored the same `hidden` tag row: the row was consumed and a
   stray prompt/secret pair landed after the validate table. Caught by grep
   sweep within one step, repaired, and — checked via `git log -S` — the
   daemon never committed a broken intermediate (single commit, repaired
   state). Lesson recorded in §e.
2. **Escaped-backslash mangling.** The SKILL.md count-row edit first wrote
   `'^\\| '` (double backslash) into the Source command, then my first fix
   attempt failed on a stale exact-match. Two wasted round-trips; final text
   verified by running the command it documents (produces 62).
3. **ci.yml anchor SHA typo.** Old_string carried `303a32d3…` for a
   `303a32d7…` pin — benign failure (refused to match), but it means I nearly
   batched a workflow edit without exact-copying the SHA.
4. **Documented-gate surprise.** I initially treated the red docs-gates
   cadence gate as a possible regression of mine and burned a diagnostic loop
   on it; AGENTS.md already documented exactly this red as the expected
   archive-ritual trigger. Read the gotchas FIRST next time.
5. **Env-dependent verify block (shipped, then fixed).** The F7 block's go
   commands ran fine locally only because `go env GOEXPERIMENT` = jsonv2 is
   set machine-wide by the dev environment; a fresh clone would have failed.
   Fixed this session (prefixes added), but it shipped broken for a window.

## e) WHAT WE SHOULD IMPROVE

1. **Multiedit anchor discipline:** never batch two edits that anchor the same
   line; verify anchor uniqueness first; grep the whole region immediately
   after any multi-insert instead of trusting the tool summary.
2. **Read AGENTS.md gotchas before diagnosing any gate/CI red** — this repo
   documents its expected reds (cadence ritual, workflow lint args, sandboxed
   treefmt).
3. **Workflow edits deserve execution-level verification** (act/`gh workflow
   run` on a scratch branch, or at minimum shellcheck of run-blocks) — YAML +
   bash-piece checks are not the same as a green run.
4. **Snapshot-then-write:** for reference tables, derive rows from source via
   script (like the sentinel count) instead of hand-typing; the two new
   tables are hand-kept copies until SU25-SU27 land.
5. **Portability check for any "runnable" doc command:** run it with the
   environment a fresh clone would have (`env -i`-ish mindset), not just the
   dev shell.
6. **Fleet review corpus counts rot:** the audit says 23/18 reviews; the
   filesystem holds 28/23 today — next audit should re-baseline (noted §f).

## f) NEXT — prioritized, each routed (harvest + ROADMAP umbrella)

1. Watch the first master CI run after push: badge Generate/Commit steps green, push lands, README badge image resolves → TODO_LIST SU24
2. Drift-guard test: README validate-rules table must match the `flags_validate.go` registry (docgen_drift_test pattern) → TODO_LIST SU25
3. Drift-guard test: 13-tag Flag Tags Reference must match the `FlagTag` struct → TODO_LIST SU26
4. Build-time include: README + website type-safe-flags share one validate-rules source → TODO_LIST SU27
5. Website mirrors: "Replace your hand-rolled", "How to verify", features.ts since-parity → TODO_LIST SU28
6. machine-contract.mdx: paste real taskctl describe bytes; add MCP-facing recipe → TODO_LIST SU29
7. Reciprocity pass 2: samber-do-best-practices / templ-components / remaining skills + sibling README mentions → TODO_LIST SU30
8. CI badge hardening: in-step SVG XML check; `pull --rebase --autostash` + retry-once → TODO_LIST SU31
9. SKILL.md: re-verify CG001-CG008 row against lint/rules.go (📖→✅); derive-or-drop the rotting taskctl-87 count → TODO_LIST SU32
10. Audit docs/QUICKSTART.md + docs/TUTORIAL.md for v5.1.0 since-gaps → TODO_LIST SU33
11. Non-destructive resolution banner on docs/reviews/2026-11-17 report (all 10 findings resolved 2026-11-20) → TODO_LIST SU34
12. Run the full buildflow battery (build + test + lint, all modules) — this session ran fmt + targeted tests only → TODO_LIST SU35
13. Decide AGENTS.md's hand-stamped 89.3% coverage mentions: date-stamp or link the badge → TODO_LIST SU36
14. Resolve Codecov: activate larsartmann/cmdguard there, or delete the now-redundant CI upload step → TODO_LIST SU37 (blocked on §g Q1)
15. Land SU22 (2026-11-08_02-29 archive batch, ~65 verdicts) → TODO_LIST SU22 (existing)
16. Land SU23 (2026-10-04 pair, 125 prose items) — greens the cadence gate → TODO_LIST SU23 (existing)
17. Wire `WithOutputFormat(output.FormatJSON)` into examples/taskctl so the flagship demonstrates the JSON error envelope live → ROADMAP candidates (SU38 umbrella)
18. Add the coverage badge to the website (hero strip or features data) → ROADMAP candidates
19. Mark `.github/badges/coverage.svg` linguist-generated in .gitattributes → ROADMAP candidates
20. README length watch: 801/900-line gate — split Flag Tags Reference + Validate rules into docs/FLAGS_REFERENCE.md when crossed → ROADMAP candidates
21. Add `nix run .#check-all` as the umbrella verify line in the README verify-block → ROADMAP candidates
22. Website recovery table: link each error_code row to its ERROR_REFERENCE anchor → ROADMAP candidates
23. Mirror the recovery-decision table into error-handling.mdx → ROADMAP candidates
24. External-link HEAD checker in CI (link checker is internal-only today; badge URL 404s until first CI run) → ROADMAP candidates
25. website changelog.mdx re-sync check at the next release cut → ROADMAP candidates
26. DOMAIN_LANGUAGE.md: add ValidateRule / coverage-badge domain terms → ROADMAP candidates
27. Scripted parity check: README trinity table vs website comparisons.ts → ROADMAP candidates
28. Re-baseline fleet-review corpus counts (audit said 23/18; filesystem has 28/23) at the next audit → ROADMAP candidates
29. SKILL.md deeper-truth: link the website guides (currently repo-docs only) → ROADMAP candidates
30. Website "Works with" section mirroring SKILL.md fleet pairings → ROADMAP candidates
31. CI: shellcheck the workflow run-blocks (badge heredoc) → ROADMAP candidates
32. Reference the describe golden-test fixture from machine-contract.mdx so doc bytes are pinned → ROADMAP candidates
33. FEATURES.md: since-column concept for versioned entries → ROADMAP candidates
34. Evaluate a cmdguard-lint doc-rule (CG009?) flagging new public options without since-annotations → ROADMAP candidates
35. Fleet-share the CI-committed-badge pattern to sibling repos → ROADMAP candidates
36. Fold CHANGELOG [Unreleased] into the next version cut (v5.1.1/v5.2.0) → ROADMAP candidates
37. Re-run check-guide-sync + link check against the LIVE site post-deploy (smoke exists; add guide-sync) → ROADMAP candidates
38. config-explain example README since-check → ROADMAP candidates
39. CI badge on release/**: currently master-only by design — confirm that stays intended → ROADMAP candidates
40. SKILL.md verification-table refresh cadence: quarterly date-stamp sweep → ROADMAP candidates
41. docs/CLI_DESIGN_PRINCIPLES.md staleness sweep (untouched this session) → ROADMAP candidates
42. Consider shields-style JSON endpoint badge later; today the committed SVG wins (recorded decision) → ROADMAP candidates
43. taskctl README: document `describe`/envelope demo once SU17-wiring (item 17) lands → ROADMAP candidates
44. docs-generator example: show GenerateDocs/describe parity → ROADMAP candidates
45. README Documentation list: link the machine-contract guide explicitly (agent story is README-absent today) → ROADMAP candidates

## g) QUESTIONS (cannot resolve myself)

1. **Codecov:** do you want larsartmann/cmdguard activated on Codecov (owner
   dashboard action — uploads already succeed via OIDC), or should I delete
   the CI upload step as dead weight now that the CI-committed badge is the
   source of truth? (SU37 is blocked on this.)
2. **Archive batches:** dedicate one session each to SU22 and SU23 (~190
   prose verdicts across three 40+ day-old reports), or batch both in a
   single docs-archaeology pass? The cadence gate stays red until both land.
3. **Badge commit traffic:** are CI commits straight to master (badge bot,
   `pull --rebase`, `[skip ci]`) acceptable next to the auto-commit daemon's
   push cadence, or do you prefer the badge hosted elsewhere (gh-pages /
   website deploy) to keep the bot off master?
