# CUE Adoption Analysis: Session Status and Self-Review

| Field        | Value                                                                   |
| ------------ | ----------------------------------------------------------------------- |
| Date         | 2026-11-02 04:52 CET (Monday)                                           |
| Session type | Research and architecture assessment, no implementation                 |
| Branch       | `master`                                                                |
| Code changed | 0 files                                                                 |
| Tests run    | 0                                                                       |
| Commits made | 0                                                                       |
| Report file  | `docs/status/2026-11-02_04-52_cue-adoption-analysis-and-self-review.md` |

## 0. TL;DR

- The session answered one question: how could cmdguard benefit from adopting CUE
  (`https://cuelang.org/`). It produced an architecture assessment with three ranked
  adoption paths, not code.
- The verdict delivered: CUE's constraint half fits config-file validation and the
  describe contract; its schema-as-source-of-truth half does not fit core.
  Recommended: build-time schema emitter first, optional sub-module second, never core.
- The biggest mistake: I recommended a generator and a new sub-module **before**
  checking the cheap non-CUE fix. The single highest-value gap I identified
  (config files silently ignore unknown keys) is one option away at
  `pkg/cmdguard/v4/koanf_loader.go:146`, verified during this report:
  `encoding/json/v2.RejectUnknownMembers(true)`.
- Evidence level: desk research only. No `cue` binary was run, no module probe was
  built, no dependency or binary-size delta was measured. Every quantitative claim
  in the answer is therefore an estimate, not a measurement. That is stated plainly
  here because it was not stated plainly there.

## 1. What this session actually did

1. Loaded the `how-to-golang` skill because the question is a Go library choice.
2. Read the v4 package surface: file layout, `config.go`, `describe.go`,
   `config_file.go`, `koanf_loader.go`, `flags_validate.go`, `go.mod`.
3. Grepped `ROADMAP.md`, `TODO_LIST.md`, `FEATURES.md`, and the repo for existing
   schema, codegen, or CUE work.
4. Grepped the `how-to-golang` references for CUE to establish whether policy bans
   or requires it.
5. Verified CUE facts against upstream docs: current version, license, package
   weight, supported file types, and the experimental status of JSON Schema output.
6. Wrote the assessment: four surface mappings, three adoption paths, cost list,
   recommendation.
7. Wrote this report, and during it verified the exact json/v2 option name and the
   import already present in `koanf_loader.go`.

## 2. a) FULLY DONE

- **The question was answered** with concrete, line-anchored mappings rather than
  generic CUE marketing:
  - Config loading: `pkg/cmdguard/v4/koanf_loader.go:146`
  - Field-level validation: `pkg/cmdguard/v4/config.go:54`,
    `pkg/cmdguard/v4/flags_validate.go:178`
  - Precedence rewriting: `pkg/cmdguard/v4/config_file.go:138`
  - Flag metadata split brain: `pkg/cmdguard/v4/describe.go:96`
- **Policy check done**: CUE appears nowhere in the `how-to-golang` references.
  It is neither banned nor required, so this is a greenfield dependency decision.
  CUE is also absent from `go.mod`.
- **Upstream facts verified**, not recalled:
  - `cuelang.org/go` is at v0.17.1, Apache-2.0, still pre-1.0.
  - The `cue` package alone imports 43 packages.
  - `cue help filetypes` lists cue, json, yaml, toml, jsonl, jsonschema, openapi,
    pb, textproto, proto, go, text, binary.
  - `jsonschema` is an **experimental output type** whose form may change between
    releases.
- **Three ranked adoption paths** with a cost story:
  1. Build-time schema emitter (no runtime dependency).
  2. Optional `cuecfg/` sub-module following the existing six-sub-module pattern.
  3. Canonical flag/config model, deferred to v5 because it is a breaking restructure.
- **Explicit "why not core" argument**, tied to the project's documented lean-core
  position (14 direct dependencies) and the stated end goal of minimal code
  requirements for developers and AI agents.
- **Named real risks**: pre-1.0 API churn, dependency weight, second-DSL learning
  cost, partial overlap with koanf, and CUE "incomplete value" errors needing a
  mapper into cmdguard's `ErrorHint` contract.

## 3. b) PARTIALLY DONE

- **"CUE could replace koanf"**: only true for the parsing half. CUE does not
  provide cmdguard's `$ENV` and `~` path expansion (`expandConfigPath`), its
  case-insensitive struct population, or its `collectKeysRecursive` set-field
  tracking used by `updateTagDefaultsFromConfig`. I overstated the replacement
  case in the answer.
- **"Map `FlagTag` to a CUE schema"**: asserted, never proven. The hard types were
  not walked through: `time.Duration`, `Enum`, `URL`, `Email`, `Port`, slices,
  maps, nested structs via `FieldTag.Index`, `count` flags, `hidden`, `local`.
- **"Emit JSON Schema or OpenAPI"**: stated without the experimental caveat I had
  fetched minutes earlier. The emitter story is therefore less stable than
  presented.
- **"CUE error UX will need a mapper"**: asserted from general knowledge, not
  demonstrated with an actual CUE validation error from this codebase.
- **Build-time path distribution**: no story for where the generated `.cue` file
  lives, how it is versioned, or how a consumer pins it against a CLI version.
- **`describe` contract growth**: I proposed a new export surface without
  addressing that it would create a second contract surface needing its own
  versioning policy under the additive-only rule on `CLIDescription`.

## 4. c) NOT STARTED

- No spike. No `go get cuelang.org/go@v0.17.1`, no scratch module, no generated
  CUE file, no `cue vet` invocation.
- No measurement: dependency count delta, binary size delta, build time delta,
  startup latency.
- No prior-art check. I did not check `invopop/jsonschema` (Go struct to JSON
  Schema) or `cue exp gengotypes`, both of which could change the recommendation.
- No ADR or STACK decision record.
- No `TODO_LIST.md` or `ROADMAP.md` entry for any recommended path.
- No prototype of the cheap non-CUE fix.
- No change to code, tests, docs, or the consumer skill.

## 5. d) TOTALLY FUCKED UP

Nothing was broken, because nothing was built. That is the only mitigating factor.
Honest list of what was actually wrong:

1. **Hero-gardening a dependency.** I recommended building a generator and a new
   sub-module while the biggest concrete win I had just identified (silent config
   typos) is a one-option change to an existing `json.Unmarshal` call. I compared
   CUE against the status quo, never against the cheapest alternative. That is a
   recommendation shaped by the question's framing rather than by the problem.
2. **Presented desk research as a dependency decision.** A stack or dependency
   choice without a measured dependency, size, and build-time delta is an opinion
   dressed as an assessment. I did not label it as an estimate at the time.
3. **Used my own fetched evidence incompletely.** I fetched the filetypes page,
   read "jsonschema is currently experimental", and then omitted that caveat from
   the recommendation.
4. **Did not propose the boring fix.** `RejectUnknownMembers` plus a config check
   would close the typo gap, would cost about three lines, and would not add a
   single dependency. It should have been option zero in the answer.

## 6. What I forgot

1. To check whether `cue` is even available in the nix devShell (`flake.nix` does
   not declare it), which determines the cost of any spike.
2. To verify the `FlagTag` to CUE mapping against the real type registry instead
   of asserting it.
3. To check `docs/adr/` for an existing decision record covering the describe
   contract, and whether a schema export would require one.
4. To check for prior art before recommending a new generator.
5. To quantify anything: dependency count, binary size, startup cost.
6. To give an effort estimate or acceptance criteria for the recommended path.
7. To show a concrete before/after error message. For a "how could this benefit"
   question, the error text is the persuasive evidence.
8. To separate transitive dependency weight from the `cue` package's own 43
   imports. Those are different numbers and I conflated them.
9. To address `GOEXPERIMENT=jsonv2` interaction: CUE ships its own JSON encoder,
   so any schema output creates a second JSON story in a repo that just
   standardized on json/v2.
10. To consider the `lint` sub-module as an enforcement point (a rule that flags
    config structs with no schema, or config files with unknown keys).
11. To examine how `describe` behaves before `Execute` (the auto-registered
    command is absent), which matters for a schema emitter's lifecycle.
12. To check whether CUE can express cmdguard's string `min`/`max` meaning
    length while numeric `min`/`max` mean value. That semantic rewrite happens in
    `flags_validate.go` and must be preserved exactly or it is a breaking change.
13. To mention that `cue vet` errors are as strict as the schema is closed, so the
    typo win depends on emitting closed structs, not just a schema.
14. To note that a generated CUE file is a second artifact to regenerate, test, and
    keep in sync, which is exactly the kind of maintenance burden the project
    tries to eliminate.
15. To check `metadata`/`FEATURES.md` for an existing "machine-readable contract"
    roadmap item before implying this was unclaimed territory.
16. To state the policy finding explicitly in the answer ("the Go policy is silent
    on CUE"), rather than leaving it in my head.

## 7. What I could have done better

- **Run a 30-minute spike and let it decide.** Scratch module outside `go.work`,
  `go get cuelang.org/go@v0.17.1`, translate one real config struct, `cue vet` one
  deliberately broken YAML, capture the error. One artifact would have replaced
  three paragraphs of speculation.
- **Lead with the cheap fix.** `RejectUnknownMembers` closes the typo gap for
  free. CUE should have been pitched as "what you get beyond that", not as the
  entry point.
- **Quantify the trade.** "Heavy dependency" is not a number. 14 direct deps
  becoming 15 is fine; 40 transitive modules is not. I measured neither.
- **Show the user experience delta.** A side-by-side of the current error and a
  CUE error would have made the value obvious or exposed it as worse.
- **Treat the second-DSL cost as a decision criterion**, not a footnote, because
  the project's stated end goal is minimal code requirements for developers and
  AI agents. A new DSL is a direct tax on that goal.
- **Name the alternative stack explicitly**: `invopop/jsonschema` for contract
  export plus stricter json/v2 decoding for config correctness would deliver most
  of path 1 and half of path 2 with zero DSL.
- **Record the outcome.** A recommendation that lives only in a chat message is
  lost. It needed an ADR stub or a TODO entry.

## 8. What I can still improve, now

- Write the three-line config correctness fix as a candidate change and let it
  prove the value before any CUE work starts.
- Add `docs/spikes/cue-config-schema/` with its own `go.mod` excluded from
  `go.work`, so the dependency never leaks into the build graph.
- Define acceptance criteria for the spike: mapping coverage of all registered
  types, error quality versus current errors, startup cost under 5 ms, and a
  dependency footprint the maintainer accepts.
- If the spike succeeds, gate adoption on an ADR, not on a status report.

## 9. e) What we should improve

1. **Config correctness first.** Unknown keys in config files are silently
   dropped today. This is a user-facing bug class, independent of CUE.
2. **Cross-field validation has no home.** `validateTagWithRegistry` is strictly
   field-local. Anything conditional needs a hand-written `WithConfigValidation`.
3. **The flag model is represented four times**: `FlagTag`, the persisted
   `flagMetadata` annotation, `FlagDescription`, and pflag itself. Every new flag
   feature must be added in four places.
4. **The describe contract has golden bytes but no machine-checkable schema**, so
   drift is caught only by byte comparison against hand-maintained fixtures.
5. **Dependency policy is silent on DSLs and schema languages.** There is no
   stated budget for what a sub-module may pull in.
6. **Error quality is the product.** Any adoption that worsens error messages
   should lose, regardless of capability.
7. **Recommendations need a paper trail.** This session produced a ranking that
   exists nowhere in the repo.
8. **Cheap fixes should be evaluated before architectural ones**, as a rule, not
   as an afterthought.

## 10. f) Next things to get done (50)

### Config correctness, no new dependency (P0)

1. Add `json.RejectUnknownMembers(true)` to the config unmarshal path at
   `pkg/cmdguard/v4/koanf_loader.go:146`.
2. Audit every other `json.Unmarshal` of user-supplied config for the same gap.
3. Decide whether unknown-key rejection is default or opt-in; document the choice.
4. Add a test proving a typo'd key fails with a field path in the error.
5. Ensure the error names the file and the offending key.
6. Add an error hint for unknown config keys under the `ErrorHint` contract.
7. Verify nested objects and arrays are handled (unknown member inside `db:`).
8. Check the case-insensitive plus reject-unknown interaction works as intended.

### Spike and evidence (P0 before any CUE commitment)

9. Add `cue` to the devShell only, or document `nix run nixpkgs#cue`.
10. Create `docs/spikes/cue-config-schema/go.mod` excluded from `go.work`.
11. Translate one real config struct into a CUE schema by hand.
12. `cue vet` a valid and an invalid YAML; capture both error texts.
13. Compare those errors against current cmdguard errors, side by side.
14. Measure dependency count delta (`go list -m all`).
15. Measure binary size delta on a hello-world CLI.
16. Measure startup latency delta for config-validating runs.
17. Enumerate which `FlagTag` types map cleanly and which do not.
18. Write the spike result into `docs/spikes/`, including a "no" verdict if earned.

### Path 1: build-time schema emitter

19. Decide the output: `.cue`, JSON Schema, or both.
20. Decide the command surface: `describe --format cue` versus a separate
    `cmdguard-schema` binary.
21. Generate closed structs so typos are actually rejected.
22. Preserve string `min`/`max` as rune length, numeric as value.
23. Map `values` to a CUE disjunction.
24. Map `required` to required fields, not optional ones.
25. Map `env` and `default` as metadata, not constraints.
26. Map `count` flags to integer bounds.
27. Emit a `schema_version` matching `descriptionSchemaVersion`.
28. Add a golden test for a generated schema.
29. Document the regeneration command in AGENTS.md.
30. Add the generated artifact policy: committed or built on demand.
31. Test that a consumer can `cue vet` a config without running the binary.
32. Verify the jsonschema output path and note its experimental status in docs.

### Path 2: optional `cuecfg/` sub-module (only if the spike says yes)

33. Create `cuecfg/go.mod` with the unversioned replace to `../`.
34. Implement a `ConfigValidator` hook that runs before unmarshal.
35. Keep CUE out of core `go.mod`.
36. Add `WithConfigSchema[T](path)` as an opt-in option.
37. Translate CUE errors into `ErrorHint`-compatible messages.
38. Add a benchmark for schema compile amortization.
39. Cover cross-field rules with a real example (`tls` requires cert and key).
40. Document when to use `cuecfg` versus `WithConfigValidation`.

### Contract and consolidation

41. Draft an ADR: "machine-readable contract and its schema".
42. Decide whether `CLIDescription` gets a CUE or JSON Schema definition.
43. Add a test that the describe fixture validates against that schema.
44. Inventory every field duplicated across `FlagTag`, `flagMetadata`,
    `FlagDescription`, and pflag; record the drift risk.
45. Scope a v5 consolidation of those four representations.
46. Add `ROADMAP.md` entries for CUE paths 1, 2, and 3.

### Guardrails and process

47. Add a lint rule idea: config struct without a schema, or unknown config key.
48. Add a dependency-budget line to AGENTS.md for sub-modules.
49. Add "cheapest alternative first" to the recommendation checklist.
50. Re-run this assessment after the spike and either kill or commit to it.

## 11. g) Three questions I cannot answer myself

1. **Dependency and size budget.** Is a pre-1.0 module with a 43-import package
   acceptable inside an optional sub-module, or must CUE stay strictly build-time
   with zero runtime footprint? This is your risk tolerance, not something the
   code can tell me.
2. **Which surface matters more.** Runtime config validation for end users, or the
   machine-readable contract for AI agents and tooling? That picks between path 2
   and path 1 as the first investment, and I cannot infer it from the code.
3. **v5 appetite.** Is a breaking consolidation of the four flag representations
   on the table for v5, or must v4.x stay additive-only? The whole path 3 depends
   on that answer.

## 12. Evidence log

| Check                                                    | Result                                                               |
| -------------------------------------------------------- | -------------------------------------------------------------------- |
| `grep -rni cue` over `how-to-golang/references`          | no matches, policy is silent                                         |
| `grep -rni "json.schema\|jsonschema\|cuelang"` over repo | no implementation, one CHANGELOG mention of the describe contract    |
| `grep cuelang go.mod`                                    | absent                                                               |
| pkg.go.dev `cuelang.org/go/cue`                          | v0.17.1, Apache-2.0, 43 imports                                      |
| `cue help filetypes`                                     | jsonschema listed, marked experimental                               |
| `go doc encoding/json/v2.RejectUnknownMembers`           | exists, `func RejectUnknownMembers(v bool) Options`                  |
| `pkg/cmdguard/v4/koanf_loader.go:146`                    | already imports `encoding/json/v2`, so the option is directly usable |

## 13. Files read this session

- `pkg/cmdguard/v4/config.go`
- `pkg/cmdguard/v4/describe.go`
- `pkg/cmdguard/v4/config_file.go`
- `pkg/cmdguard/v4/koanf_loader.go`
- `pkg/cmdguard/v4/flags_validate.go`
- `go.mod`
- `CHANGELOG.md`, `FEATURES.md`, `TODO_LIST.md` (schema-related greps)
- `~/.config/crush/skills/how-to-golang/SKILL.md`

## 14. Bottom line

The analysis was directionally right and evidentially thin, and its ordering was
wrong. The next move is not a CUE adoption; it is proving the cheap config fix and
running a bounded spike that can say no. The report above is the paper trail this
session was missing.
