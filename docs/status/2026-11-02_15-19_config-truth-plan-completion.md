# Config-Truth Plan Completion — Full Status Report

**Date:** 2026-11-02 14:57–15:19 CET (Monday)
**Session scope:** Execute the entire 566-line plan
[docs/planning/2026-11-02_05-36_config-truth-pareto-execution-plan.md](../planning/2026-11-02_05-36_config-truth-pareto-execution-plan.md)
(28 coarse tasks, 149 fine tasks, workstreams A–F + verification gate), then
answer: what was forgotten, what could be better, what is still improvable.
**Predecessor report:** `docs/status/2026-11-02_10-15_config-truth-execution-mid-flight.md`
(workstreams A–E code state when the session was interrupted this morning).

---

## Headline

The plan is executed end-to-end and the full verification gate is green:
`nix run .#check-all` passes — build + test(-race) + lint (**0 issues across
all 8 modules**) + format + go-mod-tidy for core, glamour, prompts, spinner,
telemetry, flightrecorder, lint, schema. All six of the plan's §6
workstream-verification criteria are demonstrated, two of them by mutation
(deliberately breaking the thing to prove the guard fires).

Two real defects beyond the plan were found and fixed along the way:

1. **Schema array defect:** `[]string` fields emitted `"type": "array"` with
   **no `items`** (`buildConfigKey` only set `Items` when the element node was
   non-leaf). Scalar arrays are now complete.
2. **`secret` never reached the machine contract:** `FlagTag.Secret` existed,
   `config explain` and the schema consumed it, but `flagMetadata` →
   `FlagDescription` → docgen dropped it — so `describe` and generated docs
   could not report a secret. Threaded through the full chain with a new
   docgen **Secret column**, enforced by the drift gate.

Plus a consistency fix: `LogFormat` was missing from `shouldRecurseInto`'s
named-type exclusion list (worked only because its fields are unexported).

---

## Final verification evidence (plan §6, per workstream)

<!-- check-rows:ledger -->

| Workstream | Required verification                                    | Evidence                                                                                                                                                                                                                                                                                 |
| ---------- | -------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| A          | `-race` + manual typo repro with error text reviewed     | `-race` green (core 3.2 s); live repro via `examples/config-explain` with `bad.json`: `unknown config key "log-leve1" (did you mean "LogLevel"?` — text captured, format fixed in QUICKSTART to match reality                                                                            |
| B          | Unit tests for all four sources + e2e precedence run     | `config_explain_test.go` covers flag/env/config-file/default; live run printed `log-level config-file debug info DEMO_LOG_LEVEL`                                                                                                                                                         |
| C          | Drift test fails when a fact is removed + output diffed  | **Mutation-verified**: removed the Secret cell → `flag "secret" row does not document "\| yes \|"` FAIL; reverted → green. `examples/docs-generator` output eyeballed: 8-column tables, inherited section, root documented (no committed baseline existed to byte-diff against — see §d) |
| D          | Schema validates good AND bad config; golden byte-stable | `jsonschema` CLI: good config exit 0, bad config (`typo_key`) exit 1 against the golden; golden regenerated twice, byte-identical (deterministic via sorted `PropertyMap`)                                                                                                               |
| E          | Doctor pass/fail/partial with exit codes asserted        | `doctor_config_test.go` covers the cases; added the missing **exit-code assertion** (`errors.AsType[*ExitError]`, code == 1)                                                                                                                                                             |
| F          | Docs reviewed against code; links checked                | ADR-002 + gap matrix written FROM source; relative links checked (only false positives from Go code inside fences); every documented API name grepped against source before publishing                                                                                                   |

**Gate runs:** `nix flake check` green; `nix run .#check-all` green on the
final run (earlier runs failed — honestly documented in §d).

---

## a) FULLY DONE

### Workstreams (code + tests + docs)

- **A — Config strictness (all of A1–A6).** Opt-in decision recorded
  (`WithStrictConfigKeys()`, on inside `WithBatteries()`, v5 flips default);
  loader plumbing, `ErrConfigUnknownKey` + hint, 12-test matrix, QUICKSTART /
  CHANGELOG / AGENTS.md docs. Error carries dotted path + file + suggestion;
  nested/array paths (`servers[0].port`).
- **B — Config provenance (B1–B7).** `ValueSource` enum, `CLI.ExplainConfig()`,
  `config explain` command (table + JSON), `secret:"true"` redaction with
  `v4.RedactedValue` now exported, TUTORIAL Step 11, `examples/config-explain`
  (runs: greet + explain + strict failure all demonstrated live), runnable
  godoc `ExampleCLI_ExplainConfig`.
- **C — Docgen from the machine contract (C1–C5).** `docgen.go` renders from
  `flagsFor`; 8-column table (flag/type/env/default/values/required/secret/
  help); inherited-flags section; cell escaping (incl. backslashes now);
  root-only CLIs no longer empty; drift gate mutation-verified; docs-generator
  output verified.
- **D — Schema sub-module (D1–D4).** `schema/` module complete and lint-clean
  (0 issues): `FromConfig`/`FromCLI`/`Emit`/`Marshal`, deterministic output
  (sorted `PropertyMap`, because json/v2 marshals plain maps in random order —
  empirically verified), scalar-array items fix, golden byte-stable, README
  with CI + editor + regeneration instructions, `check-jsonschema`/`jsonschema`
  validation demonstrated on good+bad configs.
- **E — Doctor config checks (E1–E2).** `CLI.ConfigDiagnostics()` (path,
  parse, unknown-key, schema checks), `WithDoctorSchemaValidator`, exit code 1
  asserted, `DoctorCommand` refactored under the funlen limit via `doctorRun`.
- **F — Consolidation and hygiene (F1–F4, all four).**
  - F1: `docs/design/2026-11-02_flag-representation-gap-matrix.md` —
    field-by-field matrix of `FlagTag`/`flagMetadata`/`FlagDescription`/
    pflag, five documented gaps with accept/reject rationale — and the
    `secret` gap it describes was FIXED in code, not just written down.
  - F2: `docs/adr/002-machine-contract-stability.md` — additive-only rule,
    `schema_version` semantics, one-producer-per-fact, redaction-in-contract,
    CUE + invopop/jsonschema rejection reasoning, golden-fixture doctrine.
  - F3: v5 flag-representation consolidation scope in ROADMAP (target model,
    breaking changes, migration, 2–3 day estimate).
  - F4: TODO_LIST CT1–CT3 entries, ROADMAP CUE-revisit criteria, CG009
    lint-rule idea, dependency-budget line in AGENTS.md (core 14, glamour 3,
    prompts 2, spinner 3, telemetry 3, flightrecorder 0, lint 5, schema 0).

### Cross-cutting, completed and verified

- **flake.nix**: `schema` added to all four module lists; "all 7 modules" → 8.
- **`.golangci.yml`**: justified `schema/` exclusion (`exhaustruct_v5` +
  `tagliatelle`) with comment; exclusion-count line in AGENTS.md updated.
- **13 lint findings fixed** in code (funlen, makezero, misspell ×5, nlreturn,
  nolintlint ×3, nonamedreturns, prealloc, mnd, cyclop via split) — zero
  suppressions added for self-inflicted issues.
- **Two external-repo drift items settled**: `flightrecorder` + `lint`
  go.mod/go.sum were stale vs the koanf v2.4.1 / go-output v0.38.7 workspace
  (pre-existing, not mine); settled with `nix develop -c go mod tidy`
  (go1.27.3) and staged so the tidy gate reflects reality.
- **`lint/provider/provider.go`**: declared the new toolsdk-v1.16.0-required
  `Outputs` field (`[lint.DefaultBaselinePath]`) — unblocked a fleet-wide
  compile break someone else's dependency bump introduced.
- **Consumer skill updated** (`.agents/skills/cmdguard/SKILL.md`): secret tag,
  key-spelling vocabulary, strict keys, config explain, `schema` sub-module
  row; fan-out symlink verified live.
- **FEATURES.md**: 8 new config-truth rows (legend-compliant statuses).
- **Plan doc updated**: findings F11 (key-spelling split brain) and F12
  (root-only docgen empty) recorded; §9 "Open questions **blocking**
  execution" → "RESOLVED during execution" with the three decisions.

### Session hygiene

- Working build recipe honored throughout (`GOCACHE=/home/lars/.cache/…`,
  `GOTMPDIR`, `GOEXPERIMENT=jsonv2`) — zero phantom compiler errors this
  session despite the shared-cache fragility.
- Stray artifacts cleaned (`config-explain` binary trashed, `bad.json`
  removed, scratch test deleted).

---

## b) PARTIALLY DONE

1. **`nix run .#check-all` — green, but only after staging.** The tidy gate
   compares worktree-vs-index, so the uncommitted flightrecorder/lint settles
   failed it twice before I understood the interaction. The settle files are
   STAGED, not committed. First daemon commit makes the gate green without
   the stage; until then the stage is what keeps it green.
2. **examples/docs-generator diff (C5).** Verified the NEW output is correct
   (8 columns, no lost facts vs the contract) but there is no committed
   pre-change baseline to byte-diff against, so "no fact is lost" is proven
   against the contract (drift gate), not against the old markdown.
3. **Secret in the docs vs code order.** The CHANGELOG claimed `describe`
   carries secret BEFORE the code did; the gap was then closed in code and
   the drift gate now enforces it — correct end state, wrong working order.
4. **`config explain` in `WithBatteries()` — deliberately NOT included**
   (it prints resolved values); documented in the option's godoc. v5 question
   recorded, undecided (see §g).
5. **Schema "experimental" status.** Marked `📦 SUB-MODULE` + experimental in
   FEATURES/README; the release decision (CT2) is deferred to you.
6. **Status reports.** The mid-flight report (10:15) is now partially
   outdated (its "Exact Next Steps" are done) but is left untouched as a
   historical snapshot; this report supersedes it. Per docs-health norms the
   old report could be annotated as done — not done.

## c) NOT STARTED

1. **Website (`website/`) content** — none of the new features (strict keys,
   config explain, secret, schema) are reflected there; CT3's README link for
   `examples/config-explain` also not done.
2. **`docs/ERROR_REFERENCE.md`** — `ErrConfigUnknownKey` is NOT in the
   sentinel catalog (checked: 7 `ErrConfig*` rows, no unknown-key row). The
   hints coverage test enforces hints in code, not the catalog.
3. **README.md** — the top-level sales page mentions none of the new
   config-truth features.
4. **lint/v0.2.0 release + BuildFlow bump (existing TODO P6)** — untouched;
   `lint/provider/provider.go` now additionally carries the `Outputs` fix that
   should ship with it.
5. **Consumer proof on erraudit + timesheets (P7)** — not run.
6. **v5 scaffold/renames (P8–P9), machine error codes (P10), validation
   hardening (P11), docs parity (P12)** — untouched (pre-existing backlog).
7. **Commit strategy decision** — nothing committed by me this session; the
   daemon's last sweep was `6873412` (this morning). Everything since is in
   the working tree (37 files: 26 modified + 11 untracked, incl. the 4 staged
   settles).

## d) TOTALLY FUCKED UP (honest list)

Ordered by severity; each caught and corrected, but they happened:

1. **Fabricated user-facing claims, caught by verification, not by care:**
   - QUICKSTART showed an error message format I invented
     (`in /home/me/...yaml: did you mean …?`); the real format is
     `loading config file: "<path>": unknown config key "x" (did you mean "y"?)`.
     Fixed only after running the binary and reading the real output.
   - TUTORIAL claimed a `--output json` runtime flag — no such flag exists;
     JSON is selected via `WithOutputFormat` at construction. Fixed.
   - I then wrote `v4.FormatJSON` — also wrong; cmdguard deliberately does
     not re-export Format* constants (AGENTS.md gotcha I should have
     remembered). Fixed to `output.FormatJSON` from go-output.
   - A TODO_LIST entry (CT4) claimed the website renders generated docs —
     unverified at write time; grep showed it doesn't. Removed.
2. **CHANGELOG written before the code it describes.** The `secret` contract
   claim was false for ~an hour of work; I closed it in code (best outcome),
   but the write-order was backwards and only self-review caught it.
3. **Godoc example expected wrong output** — I didn't know unset secrets are
   ALSO redacted (`value=***`), wrote the expected output from imagination,
   test failed. The behavior is better than my assumption; my process (guess
   the Output comment) was wrong.
4. **Misdiagnosed check-all tidy failures twice.** First treated the
   `flightrecorder not tidy` FAIL as mysterious flakiness; it took three runs
   to understand (a) HEAD's go.mod was genuinely stale from the koanf bump,
   (b) the gate compares worktree-vs-INDEX so an uncommitted settle always
   fails, and (c) the script runs `git restore` on settle failure — which
   silently reverted my first fix. The `nix develop -c go mod tidy` settle
   path was documented in AGENTS.md the whole time.
5. **Forgotten mandates caught late:** consumer SKILL.md update (AGENTS.md
   requires same-change), FEATURES.md rows, ERROR_REFERENCE entry (STILL
   missing — see §c2), the "seven-column → eight-column" doc drift I created
   myself, and the plan's §6 verification table (should have driven each
   workstream's done-check from it; instead ran the whole matrix at the end).
6. **Edits fought the tooling:** three edit-tool failures from external
   daemon modifications between read and write (had to re-read / script the
   patch), one failed `python3` heredoc with a syntax error, one edit to
   `flags.go` that had to be redone via script.
7. **golangci-lint timeout in check-all** (cold-cache + `/mnt/buildcache`
   persistence failures): first full run died mid-lint; needed
   `GOLANGCI_LINT_TIMEOUT=10m` + private `XDG_CACHE_HOME`. Environmental, but
   I burned a full gate cycle before adding the env.
8. **Small known-code debt knowingly shipped** (documented, not silent):
   `configKeyTree` parses the struct separately from `ParseFlagTags` (double
   parse per CLI construction — measured cost unknown), and
   `buildConfigKeyFromNode` derives element types via
   `configKeyElementType` rather than sharing the parse tree.

## e) WHAT WE SHOULD IMPROVE

Process:

1. **Verify-then-write ordering as a hard rule** for every CHANGELOG/docs
   claim: run the API/command first, paste the real output, then write. Three
   of the §d-1 fabrications came from writing prose ahead of execution.
2. **A claims audit step before finishing docs:** grep the CHANGELOG/diff for
   every API name, flag, and output snippet, and mechanically confirm each
   against source (I did this once, at 80% completion; it should be 100%).
3. **Run `golangci-lint` on new files immediately**, not after ~1,500 lines
   accumulated; the 13-finding batch at the end was avoidable.
4. **Drive "done" from the plan's §6 table per workstream**, not from
   intuition; two criteria (mutation test, exit-code assertion) surfaced only
   in the end-of-session sweep.
5. **check-all literacy:** read the script once (it runs `git restore` on tidy
   failure, and compares worktree-vs-index) before interpreting failures;
   record the `GOLANGCI_LINT_TIMEOUT` + private-cache env as the standard
   recipe for full-gate runs.
6. **The daemon interleaving:** re-read before every edit (it modified files
   under me three times); consider asking it to skip `*.go` while a session
   is active.

Codebase:

7. The drift gate only covers docgen↔describe; describe↔FlagTag↔annotations
   hops rely on convention (ADR-002 documents the rule, nothing enforces the
   first hops) — a lint rule or reflection test could.
8. `v4.ConfigKey`/`ConfigKeys` has no package-level godoc example.
9. `PropertyMap.MarshalJSONTo` has a symmetric unmarshal need (consumers
   decoding schemas) that is untested.
10. `examples/config-explain` has no test and is not linked anywhere.

## f) Up to 50 things to get done next

Ranked roughly by leverage; existing TODO IDs preserved.

**Release-blocking / fleet**

1. Commit the working tree in coherent checkpoints (pending your call, §g1):
   settles (`flightrecorder`+`lint` go.mod/go.sum), schema module, core
   config-truth, docgen/secret chain, docs, examples.
2. Ship `lint/v0.2.0` (P6) — now also carries the toolsdk `Outputs` fix;
   restore the version-scoped replace shape first, tag, bump BuildFlow,
   fleet baseline.
3. Add `ErrConfigUnknownKey` (and `RedactedValue`, `ConfigKeys`) to
   `docs/ERROR_REFERENCE.md`.
4. Consumer proof (P7): build + full tests of erraudit and timesheets against
   master + the lint bump.
5. Decide schema release (§g2); if yes: `schema/v0.1.0` with the release
   playbook (real pin + version-scoped replace), proxy + pkg.go.dev check.

**Docs / DX**

6. Update README.md (config files, strict keys, config explain, secret tag,
   schema emitter).
7. Update website content for the config-truth feature set.
8. Link `examples/config-explain` from README + website (CT3).
9. ~~Annotate the 10:15 mid-flight report as superseded by this report.~~ done — this pass: the 10:15 mid-flight report now carries the supersession banner
10. Add schema emitter godoc examples (`ExampleFromConfig`, `ExampleEmit`).
11. QUICKSTART/TUTORIAL parity sweep (P12): batteries, describe, hints,
    runtime tag semantics — verify every snippet compiles (extract-and-build
    harness for doc code blocks).
12. `MIGRATION_v4_v5.md` skeleton (P9 groundwork): strict-keys flip,
    FlagTag.Type removal, renames T1–T5.
13. Add the strict-keys migration note to a visible "upcoming behavior
    changes" doc beyond CHANGELOG (website callout).

**Guardrails / lint**

14. Scope CG009 "example configs stay truthful" (CT1): cross-file design note
    in the lint backlog before building.
15. Enforcement for ADR-002's first hops: a test (or lint rule) that fails
    when a `FlagTag` validity fact skips `flagMetadata`.
16. `PropertyMap` decode test (+ `UnmarshalJSONFrom` if needed).
17. Cross-module CI wiring: ensure `schema/` is in the lint sub-module's own
    dogfood corpus (0 findings) and in BuildFlow's provider tests.
18. Fuzz the key-normalization fold (Unicode/empty/separator edge cases) —
    `config_keys.go` folding logic is hand-rolled.
19. Benchmark `configKeyTree` + explain resolution (cold CLI construction
    path); decide whether the double-parse (§d8) is worth removing now.
20. `-race` the doctor + explain e2e paths under fang (currently exercised
    with fang off).

**Features (pre-v5 additive)**

21. `Registry.Redact`-style helper for consumers printing resolved values
    (only `RedactedValue` const exported today).
22. `config explain` per-flag filter (`config explain log-level`) — table is
    long for big configs.
23. `config explain --redact=false` opt-out with a loud warning? (Decide
    against by default; provenance without values is the safe default.)
24. Schema emitter: `examples/` entry (a `schema-generator` command example).
25. Schema emitter: `$id` + title/description template knobs for consumers.
26. Schema emitter: expose `Dialect`/`ExtensionKey` in README's editor setup
    snippet (partially there).
27. `doctor --json` output (machines read doctor today via human lines).
28. `describe` gains per-command `config_keys` (link the flag contract to the
    file vocabulary per command, not just root).
29. `WithConfigExplainCommand` alias `WithProvenanceCommand`? (Naming review
    candidate — "explain" vs domain term "provenance".)
30. Consider `WithoutStrictConfigKeys()` opt-out companion now (name exists
    conceptually in v5 plan; shipping the opt-out early lets batteries users
    escape hatch without code changes).

**Hygiene / repo**

31. Update `AGENTS.md` module-tree comment: `examples/config-explain/` entry.
32. `flake.nix` check-all: set `GOLANGCI_LINT_TIMEOUT` inside the script so
    gate runs don't depend on caller env.
33. check-all tidy step: compare pre/post tidy snapshots instead of
    worktree-vs-index (immune to uncommitted-state artifacts; would not have
    needed the staging workaround).
34. Diagnose/quiet the `/mnt/buildcache` golangci-lint persistence failures
    (warnings pollute every gate run).
35. Add `nix run .#check-all` to the pre-commit advisory list? (Currently
    BuildFlow pre-commit; check-all is the hard wall — document which is which
    for agents.)
36. Mark the 04-15 v4.2.0 phase1 status report as done (docs-health sweep
    found several stale reports last time).
37. TODO_LIST: fold this report's §f into the file's sections (CT entries +
    new items), keeping one source of truth.
38. `FEATURES.md`: add a "Config File Loading" row for `KoanfLoader`
    strict-keys interface (`StrictConfigKeyLoader`) so custom-loader authors
    see the seam.
39. Verify `examples/taskctl` still renders identical docs after the 8-column
    change (it has no committed docs — decide whether to commit them as the
    C5-style baseline this repo lacked).
40. ADR-001 index: add ADR-002 to any ADR list/README if one exists.

**v5 (pre-work)**

41. P8: v5 scaffold at `/v5` + Builder + default flips (strict keys included)
    - renames T1–T5.
42. P10: machine error codes + generated ERROR_REFERENCE (would fix the
    catalog-drift class of problem permanently).
43. P11: validation-semantics hardening (env/config e2e, registration-time
    defaults, Duration bounds).
44. Middleware unification (#8 in ROADMAP): collapse `Middleware[T]` +
    `ContextMiddleware[T]`.
45. Flag-representation consolidation per the new ROADMAP section (target
    model already scoped).
46. Codemod dry-run for `FlagTag.Type` removal on consumers.
47. Decide batteries×explain default (§g3) before the v5 defaults freeze.

**Small self-contained fixes**

48. `docs/design/2026-11-02_flag-representation-gap-matrix.md`: add the
    `RedactedValue` export and doctor's `config-schema` check to the matrix
    once shipped (kept out of this report's scope).
49. `schema/README.md`: add `check-jsonschema` (vs deprecated `jsonschema`
    CLI) as the primary CI example — the exact command I verified with was
    the deprecated wrapper.
50. The plan document still lists B1's "design note" as a deliverable that
    lived only in godoc + this session's reports — write the one-page
    provenance design note it asked for, or amend the plan.

## g) Questions I cannot answer myself (up to 3)

1. **Commit strategy:** the working tree holds the entire session (37 files,
   ~3,000 lines incl. docs) plus 4 staged settle files that keep check-all
   green. The auto-daemon's last sweep was this morning. Do you want explicit
   checkpoint commits now (I'd propose: settles → schema module → core
   config-truth → docgen/secret chain → docs/examples), or does the daemon
   own everything and I should not touch git?
2. **Schema release timing:** ship `schema/v0.1.0` now (experimental tag,
   README caveat, release playbook shape) so consumers can adopt the emitter
   before v4.3.0, or keep it workspace-local until the `PropertyMap` decode
   path and godoc examples (items 10/16) land?
3. **`config explain` defaults:** it is excluded from `WithBatteries()`
   because it prints resolved values. For v5 (where batteries becomes the
   default), should explain be (a) batteries-on with secrets-only redaction
   as today, (b) batteries-on with ALL values redacted unless a
   `--values` opt-in flag is passed, or (c) stay opt-out of batteries? This
   decides the security posture of the default CLI.

---

## Session metadata

- **Toolchain:** Go 1.27.3 (devShell), `GOEXPERIMENT=jsonv2`; private
  `GOCACHE`/`GOTMPDIR` used for every Go command (shared `/tmp` and
  `/mnt/buildcache` remain fragile).
- **Verification commands (all green at report time):**
  `go build ./...` (8 modules), `go test ./... -count=1` (root + 7
  sub-modules), `go test ./pkg/cmdguard/v4 ./schema -race`, `golangci-lint
  run ./...` (8×0 issues), `nix fmt` (0 changed), `nix flake check`,
  `nix run .#check-all` — plus `jsonschema` good/bad validation and the
  drift-gate mutation check.
- **Working tree at report time:** 26 modified + 11 untracked files; 4
  settle files staged; nothing committed by this session.
