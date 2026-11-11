# Status Report — Tier A Releases + v5 Scaffold (2026-11-08 03:56)

> **Session scope:** execute the Pareto plan
> ([2026-11-08_02-32](../planning/2026-11-08_02-32_pareto-post-sweep-execution-plan.md))
> top-down: Tier A (ship lint releases + BuildFlow bump + W1), then Tier B
> (the v5 break, scaffold phase). This report covers exactly that run.

**State at writing:** branch `v5-scaffold` (branched from `master` @
`2feec6f`), daemon-committed `5fbdb20` (405 files), a handful of post-daemon
edits uncommitted, 2 known red tests (both understood, fixes queued — see b).

---

## a) FULLY DONE (verified)

### Tier A — the 1% (lint releases + fleet)

1. **`lint/v0.2.0` shipped end-to-end.** Publish shape restored in
   `lint/go.mod` (real pin `cmdguard/v4 v4.2.1` + version-scoped
   `replace … => ../`; GOWORK=off build+vet green). New
   `lint/CHANGELOG.md` (0.1.0 backfilled from the GH release + 0.2.0). Root
   CHANGELOG `## [lint/v0.2.0]` section cut (Unreleased lint bullets moved
   in). Full lint suite + dogfood green (8 rules load, 0 findings on
   cmdguard itself). Annotated tag → pushed → proxy-verified via scratch
   consumer (`go get` OK) → [GH release published](https://github.com/LarsArtmann/cmdguard/releases/tag/lint/v0.2.0).
2. **Master CI un-red BEFORE the tag push (go-release Phase 4.4).** Three red
   workflows found on master (pre-dating this session): (a) 8 lint issues in
   the new sentinel-reference files — fixed for real (`fmt.Fprintf`, dead
   assignment removed, golines applied, reasoned inline nolints); (b) website
   lockfile drift — `pnpm-lock.yaml` regenerated, verified with
   `--frozen-lockfile` + full astro build (22 pages); (c) same lint job as
   (a). Commit `4557e93` + AGENTS.md ledger note for the two new inline
   nolints (sentinelFacts `funlen,maintidx`; `referenceClassOrder`
   `gochecknoglobals`).
3. **`lint/v0.2.1` shipped (bonus, the session's biggest catch).** Verifying
   v0.2.0 with a generic-shaped probe (`v4.NewCLI[Config](…)`) exposed that
   **every rule (CG001–CG008) and the cross-file tracing were blind to
   generic instantiations** — cmdguard's recommended entry points. Fixed via
   one `selectorExprOf` helper unwrapping `IndexExpr`/`IndexListExpr` at all
   four extraction sites; regression tests added
   (`generic_instantiation_test.go`: single/two type params, cross-file
   constructor union, non-generic forms). Dogfooding the fix immediately
   caught a real, previously invisible CG008 in `examples/config-explain`
   (fixed: `WithCLIVersion(version.Version)`). Tagged, pushed, proxy-verified
   (`PROXY_V021_OK`), [GH release published](https://github.com/LarsArtmann/cmdguard/releases/tag/lint/v0.2.1)
   with an honest note about the `go install` limitation (see d1).
4. **BuildFlow bumped to `lint/v0.2.1` and verified on the delivering
   layer.** `tools/go.mod` + root `go.mod` requires bumped; `execution` +
   `tools` tidied (the BuildFlow pre-commit hook hard-blocked on un-tidy
   modules first time — that gate works); flake inputs updated
   (`cmdguard` lock `f2240dd` → `2feec6f`, `go-finding`, `go-linter-sdk`);
   `vendorHash.nix` re-landed from the FOD mismatch (hash pasted once, from
   the error, per the fail-fast pattern); `nix build .#default` green twice
   (before and after the tidy); providers wiring suite green (~100s, 127
   providers, cmdguard-lint registered); `docs --check` "no drift"; fleet
   smoke with the v0.2.1 binary: **erraudit 0 findings, go-structure-linter
   0, timesheets 1×CG003 true positive** (recorded). All committed (daemon
   commits `b6ef3f608`, `136e53a49`).
5. **W1 done.** Retrospective BREAKING entry for `CLI.HealthCheck(ctx)`
   (silently shipped in v4.1.0) added to root CHANGELOG `[Unreleased] →
   Changed`, signature change verified at source (v4.0.2 `HealthCheck()`
   vs current `HealthCheck(ctx)`); TODO_LIST W1 → ✅ DONE with pointer.

### Tier B — v5 scaffold (B1 core mechanics)

6. **Branch `v5-scaffold` created** from green master; all v5 work is
   branch-local (v4 continues on master).
7. **Mechanical scaffold complete and building:** `git mv
   pkg/cmdguard/v4 → pkg/cmdguard/v5`; `package v4`/`v4_test` →
   `v5`/`v5_test` (166 files); every import path, alias (`v4 "…"` →
   `v5 "…"`), and qualified ref (`v4.X` → `v5.X`) rewritten repo-wide —
   **zero** `cmdguard/v4` references remain in any `.go` file; root module
   path → `github.com/larsartmann/cmdguard/v5`; all 7 sub-module `go.mod` files
   moved to `/v5` requires (v5-major pseudo-versions — Go rejects a v4
   pseudo-version under a /v5 path, caught and fixed); `lint` on the branch
   uses dev-shape (unversioned replace) until its v5-era release; CI
   `submodule-smoke.yml` refs → v5. Root-module test suite: **0 failures**
   (after regenerating `docs/ERROR_REFERENCE.md` — the only drift was the
   package name inside the generated comments). All 7 sub-modules build +
   test green (glamour, prompts, spinner, telemetry, flightrecorder, schema;
   lint pending — see b2).
8. **All five `TODO(v5)` markers executed (B1.7–B1.9), 0 markers remain:**
   - `CommandInfo` → `CommandMetadata` (repo-wide incl. flightrecorder,
     spinner).
   - `TypeHandler` family → `TypeCodec`: type, `TypeCodecFunc`,
     `RegisterTypeCodec` (package-level + FlagRegistry method),
     `PluginRegistrar.TypeCodec` method, files renamed
     `type_codec*.go`, tests `TestTypeCodec_*`, error-hint strings updated.
   - `SetConfig` **removed** (unsafe post-construction mutation) + its test
     deleted.
   - `RegisterInScope` → **generic `RegisterInScope[T]`** (one typed
     provider per call, stored under T's key — the any-key collision class
     is gone; nil-parent + nil-provider errors preserved; test rewritten and
     proves typed resolution via `do.Invoke[string]`).
   - `HuhRunner` → `HuhPrompter` (prompts sub-module impl; core
     `PromptRunner` interface untouched — stays huh-agnostic).

---

## b) PARTIALLY DONE

1. **`TestErrorReferenceDoc` fails on the branch** — the TypeCodec hint-text
   rename drifted the generated doc. Fix queued: one
   `go test ./pkg/cmdguard/v5 -run TestErrorReferenceDoc -update-golden` +
   review diff. (~2 min)
2. **`lint` sub-module tests fail on the branch** (typecheck error naming
   `package v5_test`) — almost certainly a lint test importing a renamed/
   removed v5 symbol (e.g. a fixture referencing `RegisterTypeHandler` or
   `SetConfig`). Needs ~10 min diagnosis. This also blocks the root
   golangci-lint pass (1 typecheck issue).
3. **`prompts_test.go` function names** still say `TestHuhRunner_*` (the
   `\b` sed can't match inside the `TestHuhRunner_` prefix) — rename to
   `TestHuhPrompter_*`. (~2 min)
4. **Nothing on the branch committed by me** — the daemon captured the bulk
   (`5fbdb20`, 405 files); post-daemon edits (CI yaml, ERROR_REFERENCE
   regen, rename stragglers) are uncommitted and should land as one
   reviewed scaffold commit.
5. **CurrentMajor bumped to `v5`** in `lint/walk.go` (B1.10) but the **CG002
   corpus re-verify** (erraudit/timesheets/go-structure-linter now flag as
   stale-major) is not re-run.

## c) NOT STARTED

- **B1.3** Builder type sketch (M12 deferred-error Builder + NewApp
  integration).
- **B1.4–B1.6** the v5 default flips (strict, signals+graceful, batteries)
  - fallout fixes; **B1.11** godox `TODO(v5)` exclusion removal in
    `.golangci.yml` (markers are gone; the config exclusion remains).
- **B2** codemod (import + rename rewriters) + `MIGRATION_v4_v5.md` +
  dry-runs on erraudit/timesheets.
- **B3** v5 flip test matrix (opt-outs green, batteries-on suite green).
- **C1–C10** consumer proof, QUICKSTART/TUTORIAL parity, koanf edge tests,
  dead v2 lint-exclusion cleanup, taskctl/examples into check-all + CI,
  samber-linter baseline + CI, HealthCheck cluster (W2/W3/W4),
  ERROR_REFERENCE CI wiring, cancellation tests (W5).
- **Tier D** bounded items (D18–D26, D33–D40) and the big bundles (website
  ports ×8, marketing/infra ~17, doc/CI guards ~15).
- **Final battery:** docs-gates 5/5, `nix flake check`, `check-all`, race,
  lint across all modules.

## d) TOTALLY FUCKED UP (honest incidents — all recovered or contained)

1. **`go install …/cmd/cmdguard-lint@v0.2.1` is rejected by the go tool**
   ("go.mod … contains replace directives"). This is **pre-existing**
   (reproduced identically on v0.1.0): the documented publish shape
   (filesystem `replace` for GOWORK=off builds) is incompatible with
   `go install pkg@version`. Consumers must `go get` + build, or use the
   BuildFlow provider. Documented in the release notes; needs an owner
   decision (see g2). Not a regression — but a real distribution gap.
2. **Invalid "pre-existing failure" claim (caught and redone).** My first
   check that BuildFlow's `tools/gomod` failure was pre-existing was a
   no-op (`git stash` of daemon-already-committed files) — the "baseline"
   run executed the bumped tree. Redone properly: the failure was
   environmental (`GOWORK=off` breaks the real-`go work use` integration
   test; green in a normal env). Lesson recorded in the BuildFlow commit
   body.
3. **BuildFlow pre-commit hook blocked my first bump commit** (untidy
   `execution`/`tools`) — fixed with `go mod tidy`, but the daemon then
   raced my re-commit, so the bump landed under heuristic
   "chore: auto-commit" messages instead of my detailed message. Content
   fully verified; history aesthetics suboptimal.
4. **Pseudo-version major mismatch** (`v4.0.0-00010101…` under a `/v5`
   module path) broke the whole workspace load — caught immediately by
   `go mod tidy`, fixed by regenerating pseudo-versions with the v5 major.
5. **Tag hygiene near-miss (self-caught):** the first `lint/v0.2.0` tag was
   cut before the master-CI fixes existed; caught by the go-release CI gate
   before push, deleted, re-cut on the fixed commit. Nothing poisoned (the
   proxy never saw it).

## e) WHAT WE SHOULD IMPROVE

1. **Verify consumers' delivering layer, not just the library API** — the
   generic-instantiation blind spot survived three releases because every
   test fixture used the non-generic call form. Fixtures should cover BOTH
   call shapes (now they do).
2. **Commit promptly after staging** to stop losing commit messages to the
   daemon race; for release-tag commits, stage+commit in one command.
3. **Sub-module test fixture sweep should ride the mechanical rename** —
   the lint typecheck failure is exactly the "one file missed by the sweep"
   class; next time run `go vet ./...` per sub-module inside the sweep.
4. **The 100+ remaining `cmdguard/v4` references in non-Go files**
   (docs/, website/, README, AGENTS.md) are intentionally untouched — they
   need the M16 docs pass, not the code sweep. Currently the tree is
   "code-v5, docs-v4" which is honest for a scaffold branch but must not
   ship.
5. **`go install` story for cmdguard-lint** needs a decision + fix (root
   `cmd/` shim, or replace-free publish shape, or GoReleaser binaries).
6. **Parallelize the long verifications** (sub-module tests + nix build +
   fleet smoke ran mostly serially).

## f) UP TO 50 THINGS TO DO NEXT (ranked)

**Immediate (unblock the branch):**

1. Regenerate `docs/ERROR_REFERENCE.md` on the branch (fix b1).
2. Diagnose + fix lint sub-module typecheck failure (b2).
3. Rename `TestHuhRunner_*` → `TestHuhPrompter_*` (b3).
   <!-- ✅ DONE 2026-11-08: renamed on v5-scaffold (TestHuhPrompter_*). -->
4. ~~Commit the scaffold as one reviewed commit (b4).~~ done — scaffold committed on v5-scaffold (bulk 9c7f05e + follow-ups); merged to master 359feb3
5. Re-run full branch battery: root tests, sub-module tests, `-race`,
   golangci-lint, GOWORK=off per-module builds.
6. CG002 corpus re-verify with `CurrentMajor = "v5"` (b5) + fix walk_test
   expectations if any.
7. B1.4 flip strict-validation default; fix fallout; opt-out test.
8. B1.5 flip signals+graceful default; fix fallout.
9. ~~B1.6 flip batteries default; fix fallout.~~ done — v5.0.0 ships the production batteries ON by default (CHANGELOG [5.0.0]; cli_defaults_matrix_test.go pins the matrix)

<!-- ✅ DONE 2026-11-08: shipped in v5.0.0 — the flip landed on the v5-scaffold branch, pinned by cli_defaults_matrix_test.go (Without* opt-outs shipped alongside); see the tagged v5.0.0 tree. -->

10. B3.1/B3.2 flip test matrix (each `Without*` proves green; batteries-on
    suite green; race + lint).
11. B1.3 Builder sketch (deferred-error Builder + NewApp integration).
12. B1.11 remove the now-dead `TODO(v5)` godox exclusion in
    `.golangci.yml`.
13. `pkg/version` check: make the version package report v5-era version on
    the branch (currently carries v4.2.1-era content).
14. Update `.agents/skills/cmdguard/SKILL.md` for v5 (import path, renames,
    flips) — it documents flag-scoping with source lines that rot.

**Tier B completion:**
15. B2.1 codemod skeleton: import rewrite `/v4`→`/v5`.
16. B2.2 codemod rename rewrites (CommandMetadata/TypeCodec/HuhPrompter,
SetConfig/RegisterInScope fixups).
17. B2.3 dry-run codemod on erraudit; diff review.
18. B2.4 dry-run codemod on timesheets; diff review.
19. B2.5 `MIGRATION_v4_v5.md`: flips + renames + codemod usage.
20. B2.6 migration guide worked example (taskctl before/after).
21. Docs sweep: the ~100 non-Go `cmdguard/v4` references (README, AGENTS,
docs/, website) — decide v4-vs-v5 split per file.
22. Root CHANGELOG `[Unreleased]` → v5.0.0 section skeleton.
23. Decide v5 lint release vehicle (lint/v0.3.0 with CurrentMajor=v5?) and
its BuildFlow bump.

**Tier C (unblocked, v4-side — can land on master):**
24. C2 QUICKSTART parity (batteries, describe, hints, runtime tag
semantics).
25. C2b TUTORIAL parity pass.
26. C3 koanf edge tests (TOML datetimes, YAML anchors, int-vs-float64,
empty file).
27. C4 dead `pkg/cmdguard/v2` lint-exclusion cleanup + art-dupl re-run.
28. C5/C6 taskctl + examples into `check-all` + CI leg.
29. C7 samber-linter baseline committed + flake/CI wiring (W6).
30. C8 HealthCheck cluster: dedupe `HealthCheckWithContext` (W2), Scope
split-brain (W3), `var _ do.HealthcheckerWithContext` guards (W4).
31. C10/W5 cancellation tests for the healthcheck path.
32. C9 ERROR_REFERENCE `-update-golden` verification wired into CI.
33. C1 consumer proof: erraudit + timesheets full suites vs HEAD.
34. W7 README/doc.go healthcheck docs refresh (D18).

**Tier D quick wins:**
35. ~~D19 rewrite/delete the no-op `TestPromptConfirm`.~~ done — the no-op is gone: real tests TestPromptConfirm_DelegatesToRunner/_WrapsError (coverage_improvement_test.go:281,297)

<!-- ✅ DONE 2026-11-08: real tests exist (TestPromptConfirm_DelegatesToRunner + _WrapsError). -->

36. D26 S2: re-check nixpkgs `go_1_27` ≥ 1.27.2; flip if available.

<!-- ✅ DONE 2026-11-08: nixpkgs go_1_27 = 1.27.3; flake already uses it — no flip needed. -->

37. ~~D33 `gopls check` infertypeargs verdict; D34 git-town/library-policy~~ done — gopls check clean on go-composable-business-types (no infertypeargs findings); git-town.toml main="master"; that repo carries no library-policy.yaml
    ~~values; record both.~~

<!-- ✅ DONE 2026-11-08: D33 — gopls check runs clean on go-composable-business-types (no infertypeargs findings); D34 — git-town.toml main="master" + github-connector=api; library-policy.yaml does not exist in that repo. -->

38. D35 file the docs-gates routed-verdict divergence upstream (repro
    exists).
39. D36 script to derive FEATURES/AGENTS dep tables + test counts from
    go.mod/go test.
40. D39 archive the 2026-11-02_00-11 HTML with its .md siblings.
41. ~~D38 eyeball re-rendered `architecture.svg` + examples/README loop.~~ done — docs/architecture.svg re-rendered 2026-11-08; examples/README covers taskctl/docs-generator/config-explain

<!-- ✅ DONE 2026-11-08: svg re-rendered 2026-11-08 01:34 (fresh); examples/README.md documents the test loop (go test ./examples/taskctl/... -count=1 -race). -->

42. D24 CT2 decision memo (schema sub-module release: ship v0.1.0 vs hold).
43. D25 config-explain into README + website.
44. D37 execute the §g1 thresholds decision (2026-10 batch archive vs keep).
45. D40 quarterly ROADMAP trim calendar note.
46. D20 hero-code single-source guard; D21 COMPARISON.md v4 prose audit;
    D22 GH release pages v4.0.0/v4.0.1; D23 lint dot/blank-import fixtures.
47. Tier D bundles: website ports ×8 (TUTORIAL first), marketing/infra
    (SECURITY.md first, 5 min), doc/CI guards (stale-ref guard first).
48. L6 (only after an Analyzer-bridge spike — do not half-ship). **→ TODO_LIST L6 — the Analyzer-bridge spike gating lives there**

<!-- fresh:skip — pull-based by design; stays in ROADMAP until the Analyzer-bridge spike is scheduled. -->

49. `go install` fix for cmdguard-lint (per g2 decision).
50. Final battery + a docs/status closing report with before/after numbers.

## g) THREE QUESTIONS (cannot answer myself)

1. **v4/v5 branch policy:** when v5.0.0 ships (root module path becomes
   `/v5` on the release), should `master` stay on v4 for v4.x patch
   releases (two-branch maintenance with the v5 branch merging to a
   release branch), or does v5 eventually REPLACE master (v4 frozen at
   v4.2.1)? This decides how the ~100 non-Go v4 references and the lint
   module's `CurrentMajor`/corpus are managed across branches.
2. **`go install` for cmdguard-lint:** keep the documented publish shape
   (filesystem replaces; `go get`-only consumption, BuildFlow as the
   distribution channel), or change the tag shape to be install-friendly
   (drop filesystem replaces at tag time; siblings resolved via real
   proxy pins — requires root's zero-pseudo sibling requires to become
   real pins first), or ship release binaries (GoReleaser)? I recommend
   the second (install-friendly tag shape) once root's sibling pins are
   real versions — but it changes the documented publish shape every
   sub-module tags with, so it is yours to call.
3. **May I push the BuildFlow repository commit** (the lint/v0.2.1 bump,
   currently local-only on its master) to origin now? The fleet wiring is
   verified locally; pushing makes CI pick it up but publishes the bump
   before the v5-era lint release exists.

---

_Verification evidence inline above; all commands and exit codes were run
live in this session. Nothing was tagged/pushed without its gate green._
