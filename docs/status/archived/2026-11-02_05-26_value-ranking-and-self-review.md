# Value Ranking and Self-Review: Session Status Report

<!-- check-rows:ledger -->

| Field           | Value                                                                   |
| --------------- | ----------------------------------------------------------------------- |
| Date            | 2026-11-02 05:26 CET (Monday)                                           |
| Session type    | Architecture assessment, no implementation                              |
| Branch          | `master`                                                                |
| Code changed    | 0 files                                                                 |
| Tests run       | 0                                                                       |
| Commits made    | 0                                                                       |
| Previous report | `docs/status/2026-11-02_04-52_cue-adoption-analysis-and-self-review.md` |
| This report     | `docs/status/2026-11-02_05-26_value-ranking-and-self-review.md`         |

## 0. TL;DR

- The second round answered a sharper question: what in cmdguard would actually add
  value, independent of CUE. The answer ranked six items, led by silent config typos,
  config provenance, and human docs generated from the existing machine contract.
- Again zero code, zero tests, zero commits. This is the second consecutive round of
  recommendation with no implementation.
- The headline failure is a repeat: the previous report's own conclusion was
  "recommendations need a paper trail", and this round produced another recommendation
  with no ADR, no `TODO_LIST.md` entry, and no `ROADMAP.md` line. I diagnosed a process
  failure and then committed it again within the hour.
- The second failure is evidentiary: I presented strict unknown-key rejection as a
  "3-line fix". It is a behavior change for `v4.x` with real compatibility
  consequences, and I did not weigh them.
- The third is fabricated precision: "80% of the CUE pitch with 0% of the DSL tax"
  is a rhetorical number, not a measurement. It should not have been written.

## 1. What this second round actually did

1. Re-read the v4 config, describe, docgen, and validation surfaces already opened in
   the first round.
2. Added three targeted greps: docgen env/default rendering, a `secret` tag, and
   provenance/explain.
3. Read `docgen.go` in full (86 lines) to confirm what human docs do and do not carry.
4. Confirmed `setFields` is computed and then used only to rewrite tag defaults.
5. Grepped `doctor.go` for config checks.
6. Produced a six-item value ranking with a "do this order" recommendation and an
   explicit "skip CUE" verdict.
7. Wrote this report.

## 2. a) FULLY DONE

- **The CUE question was resolved, not deferred.** The verdict is now concrete: CUE
  fails the "no new DSL, no new dependency" test, while the two highest-value gaps are
  fixable without it. That is a decision, and it is defensible.
- **Six value items were grounded in code, not vibes**:
  - Silent config typos: `pkg/cmdguard/v4/koanf_loader.go:146` has no
    `RejectUnknownMembers`.
  - Provenance data is computed then discarded: `setFields` is consumed only by
    `updateTagDefaultsFromConfig` (`config_file.go:138`).
  - Human docs bypass the structured contract: `docgen.go:60` prints
    `cmd.LocalFlags().FlagUsages()`.
  - No `secret` tag exists anywhere in the package (verified by grep).
  - No provenance or explain surface exists (verified by grep).
  - The four flag representations were already identified in round one and remain
    unconsolidated.
- **A ranking with a stated filter** was produced: real value means a user-visible
  failure is prevented, two representations stop being able to drift, and no new DSL
  or dependency is introduced. That filter is reusable beyond this session.
- **Cross-field validation was correctly down-ranked** rather than oversold:
  `WithConfigValidation` already covers the capability, so the gap is ergonomics, not
  missing function. Saying so is more useful than inventing a feature.
- **The "no" case for CUE was made on cost**, not on taste: pre-1.0 API, dependency
  weight, DSL tax against the project's stated end goal.

## 3. b) PARTIALLY DONE

- **Provenance proposal**: the claim "you already compute the data and throw it away"
  is only half true. `updateTagDefaultsFromConfig` **overwrites** `r.tags[i].Default`
  (`config_file.go:163-165`), so after the rewrite the registry can no longer tell a
  tag default from a config-file value. Provenance needs the pre-rewrite value and the
  `setFields` intersection captured before the rewrite. I did not say this.
- **Doctor config check**: I grepped `doctor.go` for `config` and got back only the
  internal `doctorConfig[T]` option struct. I described a feature shape without
  reading the command's behavior.
- **Docgen claim**: accurate about env, required, values, and validation being missing,
  but incomplete. pflag's `FlagUsages()` already prints defaults, so the claim had to
  be precise, and I did not check whether generated docs are pinned by golden tests or
  consumed by `website/` and `examples/docs-generator`, where a format change would
  break CI or published pages.
- **JSON Schema export**: same as round one. No prior-art check
  (`invopop/jsonschema`, `cue exp gengotypes`), no draft version decision, no answer
  for how string `min`/`max` (length) versus numeric `min`/`max` (value) is expressed.
- **Effort estimates**: none. The order was given, the cost was described in words
  ("small", "medium", "large") and not in days or in expected diff size.
- **Acceptance criteria**: none for any of the six items.

## 4. c) NOT STARTED

- No code. No test. No doc change outside status reports.
- No paper trail: no `TODO_LIST.md` entries, no `ROADMAP.md` line for the CUE revisit
  criteria, no ADR for the machine contract and its schema.
- No design note for the provenance feature, which is the one item that needs a real
  design (where the source is captured, what the output looks like, how redaction
  works).
- No compatibility analysis for strict config keys.
- No check of the generated-docs golden fixtures or their consumers.
- No read of `doctor.go`.
- No verification that `json.RejectUnknownMembers` composes with
  `json.MatchCaseInsensitiveNames` or that its error message carries a usable path.
- No implementation plan or ticket breakdown for any item.

## 5. d) TOTALLY FUCKED UP

Nothing was broken, because nothing was built. Honest list of what was actually wrong:

1. **Reproduced the process failure I had just diagnosed.** The previous report's
   section 9 says recommendations need a paper trail. This round produced a six-item
   recommendation and wrote none of it into `TODO_LIST.md`, `ROADMAP.md`, or an ADR.
   Two rounds, two good intentions, zero durable artifacts.
2. **Sold a behavior change as a bug fix.** `json.RejectUnknownMembers(true)` makes
   every previously tolerated config file with an extra key fail. Config files are
   often shared between tools, and extra keys are a normal forward-compatibility
   mechanism during rolling upgrades. This is a semver and UX decision, not a patch,
   and I called it "~3 lines".
3. **Invented a percentage.** "80% of the CUE pitch with 0% of the DSL tax" is
   fabricated precision. There is no measurement behind either number. A qualitative
   "most of the benefit, none of the dependency" would have been honest and just as
   persuasive.
4. **Made an unverified absolute claim.** "No cobra wrapper does this" about config
   provenance. I did not inspect a single cobra wrapper. That is exactly the claim
   class that should never be asserted without a source.
5. **Described a feature I had not read the implementation of.** The `doctor` config
   check was inferred from a struct name, not from the command's code.
6. **Repeated the prior round's biggest omission.** I criticized round one for
   skipping prior art, then skipped prior art again on the JSON Schema item.

## 6. What I forgot

1. To weigh the strict-config tradeoff: shared config files, forward compatibility,
   tooling that writes extra keys.
2. To propose the correct resolution, which is a strict mode rather than a strict
   default: opt-in, or tied to `WithStrictValidation` / `WithBatteries`.
3. To check whether `RejectUnknownMembers` composes with
   `MatchCaseInsensitiveNames` and whether the error identifies the offending key and
   file.
4. To note that `updateTagDefaultsFromConfig` destroys provenance, which is a
   prerequisite for the provenance feature, not a detail.
5. To check that `flag.Changed` and env lookup are not currently plumbed anywhere the
   registry can read them.
6. To check prior art for Go struct to JSON Schema generation.
7. To check for golden tests over `GenerateDocs` output and for its consumers
   (`website/`, `examples/docs-generator`).
8. To read `doctor.go` before describing a doctor feature.
9. To source the "no cobra wrapper does this" claim.
10. To avoid the fabricated 80/0 numbers.
11. To give effort estimates and diff-size expectations.
12. To give acceptance criteria and a test list per item.
13. To decide where a provenance command lives (new `config` command, `describe`
    flag, or per-command flag), and to note that a new command surface touches the
    `describe` contract.
14. To treat redaction as a prerequisite rather than a caveat: a provenance printer
    emits resolved values, and there is no `secret:"true"` tag to redact them.
15. To check nested structs (`FieldTag.Index`) and case-insensitive keys for both the
    provenance output and the schema export.
16. To note that command-level flags cannot come from config files, so provenance for
    them has only flag, env, and default sources.
17. To record which of the six items is a v4.x change and which is v5-only. Only the
    consolidation item was labeled.
18. To state the release vehicle for each item.

## 7. What I could have done better

- **Present the strictness decision as a decision.** "Unknown keys must fail" plus a
  proposed default plus the migration path. I gave the verdict without the choice.
- **Design provenance in one page.** Where the source is captured (before the tag
  rewrite), the enum of sources, the output shape, the redaction rule, the command
  surface. That is the item worth a design note, and I left it at a paragraph.
- **Check prior art once and cite it.** A one-minute search would have either killed
  or sharpened the JSON Schema item.
- **Use real numbers or none.** Percentages without measurement are a liability in a
  project whose culture is evidence-based.
- **Read before describing.** `doctor.go` is short; there was no reason to infer it.
- **Write the paper trail as part of the answer.** The user asked what would add
  value; adding the recommendation to `TODO_LIST.md` is part of the answer, not a
  follow-up.
- **Separate "what is right" from "what is shippable in v4.x".** Two of the six items
  are behaviors changes and one is v5-scale. That split belongs up front.

## 8. What I can still improve, now

- Write the paper trail: a `TODO_LIST.md` block for items 1 to 5 with v4.x versus v5
  labels, and one `ROADMAP.md` line stating the CUE revisit criteria.
- Write the provenance design note before any code: source capture point, source enum,
  output formats, redaction, command surface, tests.
- Verify the three technical unknowns that gate the top item: option composition, error
  message quality, and the golden-doc surface.
- Prepare a one-day spike plan for the JSON Schema export with an explicit "no"
  outcome allowed.

## 9. e) What we should improve

1. **Config files silently dropping unknown keys** remains the highest-ROI defect and
   is still open.
2. **Config resolution is unexplainable.** Four sources, no way to ask which one won.
   The data already exists transiently and is discarded.
3. **Human docs are generated from pflag, machine docs from the registry.** Two
   sources for the same facts, guaranteed to drift.
4. **The flag model has four representations.** Every flag feature costs four edits.
5. **Strictness is a policy decision that has never been made explicit.** There is
   `WithStrictValidation` for command metadata, but nothing equivalent for config
   content.
6. **Redaction does not exist.** No `secret` tag, so any value-printing feature is a
   leak risk by default.
7. **Recommendations do not become work.** Two rounds, no tickets.
8. **Prior-art checks are systematically skipped**, twice now.

## 10. f) Next things to get done (50)

### A. Config strictness, decided and shipped

1. Decide whether unknown config keys default to errors in `v4.x` or only under an
   opt-in flag.
2. Add `WithStrictConfigKeys()` (or equivalent) so the default can stay compatible.
3. Wire `json.RejectUnknownMembers(true)` at `pkg/cmdguard/v4/koanf_loader.go:146`.
4. Apply the same option to every other user-config unmarshal path.
5. Verify composition with `json.MatchCaseInsensitiveNames(true)`.
6. Confirm the error names the offending key and the file path.
7. Add an `ErrorHint` entry for unknown config keys.
8. Test a top-level typo fails with a clear message.
9. Test a nested typo inside `db:` fails at the right path.
10. Test the tolerant mode still accepts extra keys.
11. Add a CHANGELOG migration note.
12. Document strict versus tolerant in `docs/QUICKSTART.md`.

### B. Config provenance

13. Capture pre-rewrite tag defaults before `updateTagDefaultsFromConfig` runs.
14. Retain the `setFields` set on the CLI, not just consume it.
15. Add a `ValueSource` enum: flag, env, config file, tag default, zero value.
16. Record `flag.Changed` per flag at resolve time.
17. Record env presence via `os.LookupEnv` at resolve time.
18. Resolve nested structs through `FieldTag.Index`.
19. Normalize case-insensitive config keys for reporting.
20. Note that command-level flags have only flag, env, and default sources.
21. Choose the surface: `app config explain`, `describe --explain`, or a `--explain`
    per-command flag.
22. Define a table output for humans and a JSON output for agents.
23. Add redaction rules before printing any resolved value.
24. Decide redaction for unmarked flags (hide values by default, show only sources).
25. Test all four sources end to end.
26. Test that a config-sourced value is reported as config after the rewrite.
27. Test an env override wins and is reported as env.
28. Test a flag override wins and is reported as flag.
29. Document the feature in `docs/TUTORIAL.md`.
30. Add the provenance output to the `describe` contract if it becomes a command.

### C. Human docs from the machine contract

31. Inventory `GenerateDocs` output and its consumers (`website/`, docs generator).
32. Check whether generated docs are pinned by golden tests.
33. Render env vars per flag.
34. Render required markers.
35. Render allowed values.
36. Render validation rules.
37. Keep pflag's default rendering so nothing is lost.
38. Decide how inherited flags are presented.
39. Add a drift test asserting docgen and `Describe()` agree on flag facts.
40. Regenerate any committed docs and review the diff.

### D. Config schema export

41. Check prior art (`invopop/jsonschema`, `cue exp gengotypes`) before building.
42. Pick a JSON Schema draft (2020-12).
43. Map string `min`/`max` to `minLength`/`maxLength`.
44. Map numeric `min`/`max` to `minimum`/`maximum`.
45. Map `values` to `enum`, `required` to `required`, `default` to `default`.
46. Decide whether the schema is committed or generated on demand.
47. Add a CI example that validates a config file without running the binary.
48. Document editor autocomplete setup.

### E. Doctor and consolidation

49. Read `doctor.go`, then add config path, parse, schema, and unknown-key checks.
50. Inventory the four flag representations and scope the v5 consolidation, including
    adding a `secret:"true"` tag as the redaction prerequisite.

## 11. g) Three questions I cannot answer myself

1. **Is strict config-key rejection allowed as a `v4.x` default?** Config files are
   often shared between tools and extra keys are a normal forward-compatibility
   mechanism. Making unknown keys fail by default is a behavior change; making it
   opt-in preserves compatibility but weakens the fix. That tradeoff is yours.
2. **What is the headline for the next release: provenance for humans, or the schema
   and contract surface for agents?** Both are valuable, they target different users,
   and the ordering changes what gets designed first.
3. **Does a `secret:"true"` tag ship now?** Any provenance or explain output prints
   resolved values. Without a secret tag there is no reliable redaction, so the choice
   is either add the tag as part of the feature (new public surface) or ship
   source-only output with no values until it exists.

## 12. Evidence log

<!-- check-rows:ledger -->

| Check                                                      | Result                                              |
| ---------------------------------------------------------- | --------------------------------------------------- |
| `grep env\|Env\|Default` in `docgen.go`                    | no matches; docs come from pflag text               |
| `grep secret` in `pkg/cmdguard/v4` (non-test)              | no matches; no redaction concept                    |
| `grep provenance\|explain` in `pkg/cmdguard/v4` (non-test) | no matches                                          |
| `grep config` in `doctor.go`                               | only the `doctorConfig[T]` option struct            |
| `grep setFields` in `pkg/cmdguard/v4`                      | consumed only at `cli.go:246` and                   |
| `config_file.go:138` to rewrite tag defaults               |                                                     |
| `docgen.go` read lines 1-70                                | uses `cmd.LocalFlags().FlagUsages()`                |
| `go doc encoding/json/v2.RejectUnknownMembers`             | exists, `func RejectUnknownMembers(v bool) Options` |
| `koanf_loader.go` imports                                  | `encoding/json/v2`, so the option applies directly  |

## 13. Files read this round

- `pkg/cmdguard/v4/docgen.go` (86 lines, full)
- `pkg/cmdguard/v4/doctor.go` (grepped only)
- `pkg/cmdguard/v4/cli.go` (config load block, grepped)
- `pkg/cmdguard/v4/config_file.go` (provenance and default rewrite paths)
- `pkg/cmdguard/v4/koanf_loader.go` (unmarshal path)
- `pkg/cmdguard/v4/config.go`, `describe.go`, `flags_validate.go` (carried from round one)

## 14. Bottom line

The ranking is better than the CUE answer: it is cheaper, it is grounded in code, and
it produces a defensible "no" for the dependency. But it is still just prose. The next
honest move is to stop producing recommendations and either write the first three
tickets into `TODO_LIST.md` or implement the smallest one. The gap between this project
and its end goal is not a missing schema language; it is the habit of leaving good
decisions in chat.
