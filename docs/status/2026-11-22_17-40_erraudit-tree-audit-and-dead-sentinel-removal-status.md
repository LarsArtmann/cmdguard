# Status: erraudit tree audit → dead-sentinel removal (2026-11-22 17:40)

**Scope line.** VERIFIED this session (fresh command output): `erraudit tree ./... --max`
(70→65 named errors, editor/nil nodes gone), `erraudit lint ./... --max --disable
generic_return` (0 real findings: no legacy_as/context_loss/silent_swallow/legacy_is),
`go build ./...` across all 8 modules, full `go test ./... -count=1 -race` (root +
flightrecorder) plus plain test runs of glamour/prompts/spinner/telemetry/lint/schema,
`golangci-lint run ./...` 0 issues, `docs-gates run .` all gates PASS,
`scripts/check-guide-sync.sh` OK, and the new tripwire probe-verified in BOTH directions
(plants a dead sentinel → FAILs; clean tree → PASSes; done twice). NOT verified:
`nix run .#check-all` as a single gate (components run individually instead), website
build/astro check/smoke (no website content changed — `changelog.mdx` deliberately
untouched), `erraudit fix` dry-run / `--type-aware` sweep across modules, `nolint-audit`,
fan-out guard `scripts/check-skill-fanout.sh` after the SKILL.md edit, and the git commit
state (auto-commit daemon owns it; harness forbids manual commits). This report covers
ONLY this session's run.

## Session summary

Input: `erraudit tree ./... --max` output — 70 named errors, 0 hierarchy relationships,
286 inline wrap sites, `[1 violations]` badge on every node, `unreferenced` on 9.

Investigation (erraudit source at `~/projects/erraudit`): `[N violations]` badges are
declaration-LINE violations attributed by (file,line) match
(`internal/analyzer/violation_attribution.go`); every sentinel earned exactly one
`stdlib_constructor` style finding from the opt-in `--enforce-go-error-family` policy
that `--max` enables; tree is a visualization and always exits 0; cmdguard has NO
erraudit step in CI/BuildFlow, so this was an audit surface, not a broken gate.

Triage of the 9 `unreferenced` sentinels: 3 are wired via plain `return ErrX` (invisible
to erraudit's wrap-site counting — `ErrAlreadyStarted`, `ErrNotEnabled`,
`errKoanfReadNotImplemented`); 5 were genuinely dead: `ErrEditorTempFile`, `ErrEditorWrite`,
`ErrEditorRun`, `ErrEditorRead`, `ErrNilValue` — declared + fact-tabled but returned by
NO code path anywhere. Git archaeology: v3's "cut dead weight" commit removed the editor
sentinels, v4.3.1 had none of the five, and the v5 scaffold auto-commit (`5fbdb20`, 405
files, heuristic) mechanically re-added them. No editor feature is planned in
ROADMAP/FEATURES/ADRs; the only live-surface reference was the generated
`docs/ERROR_REFERENCE.md`.

## a) FULLY DONE

1. **Five dead sentinels removed** — declarations (`errors.go`, `errors_flags.go`), 5
   facts rows (`errors_facts.go`), `editor` class dropped from `referenceClassOrder`
   (`errors_reference.go`); sentinel fact table now covers **57** sentinels.
2. **`docs/ERROR_REFERENCE.md` regenerated** via `TestErrorReferenceDoc -update-golden`
   — "57 sentinels" header, `nil_value` row and the whole editor section gone; drift
   test green.
3. **`TestNoDeadSentinels` tripwire added** (`pkg/cmdguard/v5/errors_dead_sentinel_test.go`)
   — go/parser AST walk: any package-level `Err*`/`err*` var initialized by a call (or
   typed `error`) that is never referenced outside its declaration and `errors_facts.go`
   fails the test. Probe-verified twice (planted sentinel → FAIL with the right message;
   removed → PASS). Prevents the exact bug class the audit surfaced.
4. **Disposition recorded**: `stdlib_constructor` badges on `errors.New` sentinels are
   EXPECTED and accepted — cmdguard's facts table/codes/hints ARE its error family;
   go-error-family adoption stays ROADMAP-parked. AGENTS.md now carries the full
   gotcha (badge semantics, unreferenced≠dead caveat, verify-with-rg rule).
5. **Consumer SKILL.md updated** (`.agents/skills/cmdguard/SKILL.md`): 62→57 in both the
   fact-table row and the prose bullet, derivation date refreshed, and the stored
   derivation pattern FIXED (the old `rg` pattern had drifted from the
   formatter-normalized table shape and matched nothing; corrected and re-derived = 57).
6. **CHANGELOG.md `[Unreleased] ### Removed`** entry documenting the five sentinels, the
   v3 precedent, the audit that caught them, and "no runtime behavior changed" (no
   caller could ever have matched them via `errors.Is`).
7. **Lint-zero restored** after the edits: trailing blank line removed from `errors.go`,
   nlreturn/golines formatting fixed in the new test file, and the now-unused
   `maintidx` half of the `sentinelFacts` nolint directive trimmed (the facts table got
   short enough that the linter flagged the directive as unused — `funlen` kept).
8. **Full verification battery green**: 8 modules build; full `-race` suite passes;
   golangci-lint 0 issues; docs-gates PASS (freshness/audit-all/annotate/routing/cadence);
   check-guide-sync OK; `erraudit tree` re-run shows 65 named errors with the editor/nil
   nodes gone; strict lint posture (`--max --disable generic_return`) shows zero real
   findings.

## b) PARTIALLY DONE

1. **Tripwire scope** — covers only the core `pkg/cmdguard/v5` package. `flightrecorder`
   (`ErrAlreadyStarted`, `ErrNotEnabled`) and `lint` (`errFindings`,
   `ErrBaselineNotFound`, `ErrBaselineVersionMismatch`) sentinels are unguarded.
2. **erraudit acceptance posture** — recorded as prose in AGENTS.md, but there is no
   machine-enforced artifact (no erraudit config file pinning `stdlib_constructor` /
   `generic_return` as accepted per module, no CI gate either way). A future session can
   re-litigate the decision because nothing pins it.
3. **tree-vs-lint divergence documented only shallowly** — `tree --max` runs the
   enforcement analyzers (hence the badges); `lint --max` drops them (0
   stdlib_constructor findings in lint mode). AGENTS.md says the badges are expected;
   the WHY (runner drops them in lint) is in the skill, not the repo.
4. **SKILL.md fan-out** — edited the in-repo SOURCE (correct per authoring norm), but
   did not run `scripts/check-skill-fanout.sh` to confirm the fan-out link still
   resolves through `~/.config/crush/skills/cmdguard`.
5. **`unreferenced` badge noise** — the 3 wired sentinels will keep showing `unreferenced`
   (plain returns are idiomatic; rewriting them to `%w` wraps just to satisfy a
   visualization was rejected). The noise is documented, not eliminated.

## c) NOT STARTED (noticed, deliberately out of scope)

1. `erraudit fix` dry-run (`errors.As`→`AsType` sweep) across all 8 modules.
2. `erraudit lint ./... --type-aware` full advisory classification sweep per module.
3. `nolint-audit` per module (the skill's suppression-token audit).
4. go-error-family adoption spike (ROADMAP-parked by decision, not by omission).
5. Per-module generic_return disable posture in a runbook/CI config.
6. TODO_LIST/ROADMAP harvest of this report's section (f) — see the harvest note below.
7. Any website rebuild/deploy (nothing website-facing changed).

## d) TOTALLY FUCKED UP

Nothing is broken — every gate is green and the tree is clean. Honest process failures
that did NOT damage state:

1. **Probe hygiene** — verifying the tripwire by appending sentinel lines to real source
   files (`printf >> ... && sed -i delete`) briefly dirtied the working tree and once
   left a trailing blank line in `errors.go` (caught by golangci-lint, fixed). A scratch
   copy of the package in a temp dir would have been safe.
2. **Living-doc ritual violated by ordering** — AGENTS.md/CHANGELOG edits happened BEFORE
   `docs-freshness --pre-flight` (the ritual wants pre-flight first, battery before
   commit). Outcome was green, but the order was wrong.
3. **One wasted verification round** on stale gopls diagnostics (LSP flagged nlreturn on
   the already-fixed file; golangci-lint was ground truth).

## e) WHAT WE SHOULD IMPROVE

1. Pin the erraudit acceptance posture in a MACHINE artifact (config file or CI job),
   not just AGENTS.md prose.
2. Generalize the dead-sentinel tripwire (per-module copies or a shared `pkg/testutil`
   helper) so flightrecorder/lint get the same guard.
3. Adopt the pre-flight-first ritual strictly for any living-doc edit.
4. Use temp-dir package copies for negative-test probes instead of append+sed on source.
5. Add a cross-file drift tripwire: SKILL.md's derived sentinel count vs the generated
   ERROR_REFERENCE header (the stored derivation pattern rotted once already).
6. Decide the release vehicle for the API-surface correction (see question 1).
7. Run `nix run .#check-all` once as the closing gate instead of piecemeal equivalents.
8. After any SKILL.md edit, run `scripts/check-skill-fanout.sh` (verify the delivering
   layer, not just the commit).

## f) UP TO 50 THINGS TO GET DONE NEXT (brainstorm — most are ROADMAP fuel)

Sentinel & error system:
1. Owner decision: ship the 5-sentinel removal in v5.2.0 (minor + explicit
   API-surface-correction note) or hold for v6.
2. Extend `TestNoDeadSentinels` to the `flightrecorder` module.
3. Extend the tripwire to the `lint` module (`errFindings`, `ErrBaseline*`).
4. Extract the tripwire's scan logic into `pkg/testutil` for fleet-wide reuse.
5. Add a self-test fixture (temp-dir mini-package) inside the tripwire test so its
   detection logic is regression-tested without touching real source.
6. Confirm the tripwire handles future typed sentinels (`var X error = dashError(...)`)
   — parser path supports it; add the fixture to prove it.
7. Decide: wrap the 3 plain-return sentinels in `%w` for erraudit visibility, or accept
   the permanent `unreferenced` badges (current: accept).
8. Audit the 286 inline wrap sites for `errors.Is`-chain consistency (151 attributed —
   the rest wrap runtime values; spot-check the unattributed tail for lost sentinels).
9. Run `erraudit fix` dry-run across all 8 modules (`errors.As`→`AsType` sweep).
10. Run `erraudit lint ./... --type-aware` per module and disposition every `errors.Is`
    advisory (sentinel vs migrate).
11. Run `nolint-audit` per module; confirm no stale `//nolint:erraudit` tokens.
12. Consider carrying the envelope class into the JSON error envelope (currently only
    `error_code` + `hint`) — or reject as YAGNI explicitly.
13. Note the flat-by-design hierarchy in ERROR_REFERENCE intro prose (one sentence), so
    future `tree` audits don't re-derive it.
14. Keep AGENTS.md's loose "40+" sentinel phrasing vs tightening to 57 — decide per D15
    (counts live in ERROR_REFERENCE only; likely leave loose).

erraudit posture:
15. Decide: make erraudit a real CI gate (pinned `go run …@v1` like docs-gates) or keep
    it audit-only.
16. If gated: pin the exact invocation (`lint ./... --max --disable generic_return`,
    per module) in ci.yml.
17. If not gated: write the acceptance into an erraudit config file so findings are
    dispositioned by tooling, not prose.
18. Document in AGENTS.md WHY tree shows badges while lint mode shows none (runner
    drops enforcement analyzers in lint).
19. Track the erraudit binary version in use (`ee8a415`); bump consciously when the
    fleet standardizes a new one.
20. Add `stdlib_constructor` + `generic_return` to a repo-local "known accepted
    findings" inventory with dates and reasons (ratchet-style, like audit_max).

Docs & truth surfaces:
21. Harvest this report's section (f) into TODO_LIST/ROADMAP (docs-health HARVEST).
22. Update ROADMAP.md:194 (family-specific format constructors) to reference today's
    disposition so the row carries its decision context.
23. Schedule the website `changelog.mdx` update to ride the next release's changelog
    sync (historical v5.0.0 notes deliberately untouched).
24. Verify `website/public/llms.txt` regeneration picks up the 57-count on next build.
25. Add the guard from e) item 5 (SKILL count vs ERROR_REFERENCE header tripwire).
26. Re-derive the SKILL.md count again at the NEXT sentinel addition (process note in
    the SKILL row already dates it — keep dating on every change).
27. Annotate this report when superseded (docs-gates battery governs; do not rewrite).
28. Check whether FEATURES.md mentions sentinel counts anywhere (swept today: clean;
    keep it that way via e) item 5's guard).

Release & CI:
29. Watch the next master-push CI run: coverage badge job, dogfood docs-gates, and the
    brandid gate should all stay green with the smaller facts table.
30. Verify the lint module's `rules_table_test.go` and dogfood CLI are untouched by the
    core change (they are — confirmed by build+test — but re-check on release).
31. Confirm pkg.go.dev renders the updated godoc after the next tag.
32. Re-check `examples/taskctl` describe output shape on next release (contract file
    unchanged by sentinel removal — verified via full test suite).

Process & tooling:
33. Pre-flight ritual: `docs-freshness --pre-flight` before ANY living-doc edit, battery
    before commit — make it the default first step.
34. Replace source-append probes with temp-dir package copies everywhere.
35. Run `nix run .#check-all` as the single closing verification next session.
36. Run `scripts/check-skill-fanout.sh` after SKILL.md edits.
37. Consider letting BuildFlow auto-fix gci/golines/nlreturn so 8-finding lint fallout
    never needs manual rewriting (it auto-fixes some classes; nlreturn needed hands).
38. Keep the tripwire's parser failure mode fatal (current) — do not soften to skip.
39. Note multi-name var specs (`var a, b = f()`) are skipped by the tripwire — acceptable
    for sentinels; document in the test godoc if it ever matters.
40. When a sentinel is next added: declare + wire in the SAME change (tripwire enforces),
    and bump the SKILL derivation date in the same change (AGENTS gotcha now says both).

Longer-term / ROADMAP candidates:
41. go-error-family adoption spike (one module first, e.g. flightrecorder) — parked.
42. Evaluate whether the 0-relationship flat hierarchy deserves a packaged "hierarchy
    explainer" (describe-side or docs-side) for consumers auditing their own trees.
43. Consider a fleet-wide erraudit-inventory script (cqrs-htmx pattern:
    `--strict` raw-count visibility, advisory-first).
44. Sweep other fleet repos for the same v5-scaffold dead-sentinel pattern (the
    scaffold that re-added these five was mechanical — siblings may carry similar noise).
45. Add sentinel-declaration contract to the fleet linter (CG rule candidate: dead
    sentinel detection), feeding back the tripwire idea upstream to cmdguard-lint.
46. Evaluate `errors.AsType` usage coverage repo-wide once the fleet Go version settles
    past 1.26 (AGENTS already mandates it; verify no stragglers).
47. Consider extending the JSON envelope's `error_code` stability test to the generated
    reference (byte-stable codes already tested; cross-check doc codes vs runtime map).
48. Review `coverage_improvement_test.go` SA1012 (nil context) gopls warning —
    pre-existing, untouched this session.
49. Re-run the tree audit after the next few features land; keep the 65-count visible
    in the next status report.
50. Close the loop on this session's SKILL.md pattern fix: re-run the corrected `rg`
    derivation in the next session to confirm it still matches the formatter's shape.

## g) QUESTIONS I CANNOT FIGURE OUT MYSELF (3)

1. **Release vehicle for the sentinel removal** — these are exported identifiers
   (`v5.ErrEditorRead`, `v5.ErrNilValue`, …) that no code path ever returned. Ship in
   v5.2.0 as a minor with an explicit "API surface correction" note (compile-time break
   only for consumers who referenced the identifiers directly), or hold for v6 per strict
   semver? The project is semver-disciplined; this call is yours.
2. **erraudit as CI gate or audit surface?** — Should I wire a pinned
   `erraudit lint ./... --max --disable generic_return` job into ci.yml (per module,
   like the docs-gates dogfood job), or keep erraudit on-demand only and pin the
   acceptance posture in a config file instead? Gate-vs-audit is a risk-posture decision.
3. **Tripwire reach** — extend `TestNoDeadSentinels` to `flightrecorder` and `lint` now
   (small, catches the same rot class early), or keep it core-only until a sentinel
   actually goes stale in those modules? Effort-vs-value is your judgment.

---

*Basis: session of 2026-11-22 ~17:00–17:40. Every claim in sections a)–d) was verified
by a command run in that window; section f) is brainstorm, not commitment.*
