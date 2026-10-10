# cmdguard-lint

Static-analysis rules that detect anti-patterns in codebases consuming
cmdguard — derived from real consumer audits (timesheets deep dive, 2026-10-04)
and cmdguard's documented execution contract. Built on
[go-finding](https://pkg.go.dev/github.com/larsartmann/go-finding) (finding
model) and [go-linter-sdk](https://pkg.go.dev/github.com/larsartmann/go-linter-sdk)
(rule/registry scaffolding); registers as a BuildFlow tool via
[go-finding/toolsdk](https://pkg.go.dev/github.com/larsartmann/go-finding/toolsdk).

This is an optional sub-module: `github.com/larsartmann/cmdguard/lint`. Import
it from BuildFlow (via `lint/provider`) or run the standalone CLI; the library
package depends only on go-finding and go-linter-sdk.

## Rules

<!-- BEGIN RULES TABLE: generated from allRuleDefs() — regenerate with `go run ./cmd/cmdguard-lint rules --markdown` -->

| ID    | Severity | Rule                       | Description                                                                                                                                              |
| ----- | -------- | -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------- |
| CG001 | critical | execute bypass             | fang.Execute runs the raw cobra tree, bypassing cli.Execute and its signal handling, graceful shutdown, cleanup hooks, and single-error-display contract |
| CG002 | error    | stale cmdguard major       | importing a frozen cmdguard major blocks all fixes and sub-module features; migrate to the current major                                                 |
| CG003 | error    | panic on constructor error | panicking on a cmdguard constructor error re-introduces the panic cmdguard removed by design; return the error instead                                   |
| CG004 | warning  | runtime SetVersion         | cli.SetVersion mutates after construction; pass WithCLIVersion to NewCLI (and use the VersionCommand helper) instead                                     |
| CG005 | error    | duplicate version options  | WithCLIVersion and WithFangOptions(fang.WithVersion(...)) in one NewCLI call pass duplicate fang version options                                         |
| CG006 | warning  | execute error reprinted    | the error returned by cli.Execute is already displayed by cmdguard; re-printing it double-reports the failure                                            |

<!-- END RULES TABLE -->

The table is generated (L7 drift guard): `rules_table_test.go` fails when the
section between the markers diverges from `allRuleDefs()`. After changing a
rule, regenerate with `go run ./cmd/cmdguard-lint rules --markdown` and
replace the marked section.

Detection is syntactic (AST walk, no type information, no build required):
multi-signal matching with dataflow-lite receiver tracing keeps false
positives near zero. Validated against the audited timesheets tree (18
findings, every location matching the hand-audit) and negative corpora
(cmdguard itself, erraudit, go-structure-linter, branching-flow: 0 findings).

## Usage

### Standalone CLI

```bash
go run ./cmd/cmdguard-lint lint --dir .            # text output, exit 1 on findings
go run ./cmd/cmdguard-lint lint --output sarif     # SARIF for CI upload
go run ./cmd/cmdguard-lint lint --disable CG004    # rule filtering
go run ./cmd/cmdguard-lint rules                   # list rules
go run ./cmd/cmdguard-lint rules --markdown        # rule table for this README
```

### Baseline / Ratchet Mode

Adopting on a codebase with pre-existing findings? Record them as the
baseline once, then only NEW findings gate the run:

```bash
go run ./cmd/cmdguard-lint lint --write-baseline   # snapshot current findings, exit 0
go run ./cmd/cmdguard-lint lint --dir .            # now fails only on findings beyond the baseline
```

The baseline defaults to `.cmdguard-lint-baseline.json` in the linted
directory (override with `--baseline <path>`). Matching is per rule+file with
line-drift tolerance: an edit that shifts a baselined finding's line keeps it
baselined; a surplus finding for the same rule+file is new and fails the run.
When findings get fixed, the run prints the stale entry count — rerun
`--write-baseline` to tighten the ratchet. Baselines never grow on their own.

### BuildFlow

Import `github.com/larsartmann/cmdguard/lint/provider` for the self-registering
toolsdk spec (`Name: "cmdguard-lint"`, `ModuleFanOut: true`); BuildFlow
discovers it via `toolsdk.All()`.

### Library

```go
findings, err := lint.Detect(ctx, ".")                     // single-pass, all rules
report, err := lint.NewRegistry().Run(ctx, ".")            // go-linter-sdk path
rules := lint.AllRules()                                   // for FilterRules enable/disable
```

Analysis results are memoized per directory for the process lifetime;
long-lived processes re-linting a mutated tree must call `lint.ClearCache()`.

## Suppressions

Reason-bearing, rule-scoped, in-source:

```go
//cmdguard-lint:ignore CG004 this test exercises the SetVersion API itself
cli.SetVersion("2.0.0")
```

Place the directive on the offending line or the line directly above it. A
directive without a reason does not suppress.

## Adding a rule

1. Add the stable ID constant and a row to the `allRuleDefs()` table in
   `rules.go` (identity header: ID, name, description, category, severity).
2. Write `check<Name>(proj *project, meta linter.RuleMeta) []finding.Finding`
   in `rule_<name>.go`, documenting the rule's origin and trust signals.
3. Add positive AND negative fixtures to `rules_test.go`; state the rule's
   false-positive budget in the doc comment.
4. Refresh the README rule table: `go run ./cmd/cmdguard-lint rules --markdown`
   (between the markers — `rules_table_test.go` fails the build on drift).
5. Sweep the corpora: the audited timesheets tree (positive) and cmdguard
   itself (negative). A rule that cannot hold its FP budget on the corpus
   ships disabled.
6. Bump `CurrentMajor` when cmdguard ships a new major so CG002 follows.

## Origin

The rule set encodes the mechanically detectable findings of the timesheets
cmdguard deep dive (2026-10-04): F1 execute bypass, F2 stale major, F6
must-wrapper panics, F7 version plumbing, plus the ADR-001 fang-option
misuse and the v4 error/exit display contract. The deep-dive's
absence-of-usage findings (hand-rolled config, unused sentinels, ungrouped
help) are deliberately not rules: they are not attributable to a file:line
without heavy type analysis and would drown adoption in triage noise.
