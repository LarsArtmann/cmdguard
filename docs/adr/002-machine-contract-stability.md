# ADR-002: The Machine Contract (`Describe`) — Stability, Versioning, and Consumers

**Status:** ACCEPTED
**Date:** 2026-11-02
**Scope:** `pkg/cmdguard/v4/describe.go`, `config_explain.go`, `docgen.go`, `schema/`

## Context

cmdguard emits machine-readable descriptions of a built CLI:

- `CLI.Describe()` / `WriteDescription()` — the command-tree contract
  (`CLIDescription`, `schema_version` field)
- `ConfigExplanation` — the `config explain` JSON contract (own
  `schema_version`)
- the `schema` sub-module — the JSON Schema for config files
- generated markdown docs (`GenerateDocs`), which render from the same
  `flagsFor` source as the contract

Four consumers already exist: AI agents driving the binary, the docgen drift
gate (`docgen_drift_test.go`), `config explain`, and the schema emitter.
Each depends on field names and shapes staying stable. Nothing recorded the
rules. This ADR is that record; the field-level inventory lives in
`docs/design/2026-11-02_flag-representation-gap-matrix.md`.

## Decision

### 1. Additive-only within a major version

Contract fields may be added freely (`omitzero` keeps old consumers working);
they may never be removed, renamed, or change shape within `v4`. Breaking
changes ship with the next major version and are listed in the CHANGELOG's
upgrade notes. This is the same rule cmdguard applies to its Go API, extended
to the JSON surface.

### 2. `schema_version` is a compat string, not a changelog

`SchemaVersion` (`"1"` today, shared by the command-tree contract and
`ConfigExplanation`) changes only when the shape breaks: a field removed or
renamed, a semantic changed. Additive fields keep the version. Consumers gate
on equality, not ordering, so we never increment for new fields. If a future
v5 contract lands with breaking shape changes, the version string moves in
lockstep.

### 3. One source per fact

Every contract fact has exactly one producing function (`flagsFor`,
`flagMetadataFromTags`, `ConfigKeys`). Renderers (docgen, `config explain`
table, schema emitter) consume; they never re-derive from pflag or reflection.
The docgen drift test enforces the docgen half mechanically: a fact in
`FlagDescription` that stops appearing in the markdown fails the suite.

### 4. The contract reports validity facts, not wiring facts

Fields that affect whether an invocation is valid (required, values, env,
secret, prompt, validate, count) are contract fields. Fields that describe how
the CLI was assembled (`local:"true"` wiring, reflect paths) are not. The gap
matrix lists each field's layer explicitly.

### 5. Redaction is part of the contract, not the renderer

`secret:"true"` values are redacted to `***` at explanation-build time, and
`ResolvedFlag.Redacted` records that fact. A consumer can trust any JSON dump
of `config explain` not to leak a tagged secret, without knowing which
renderer produced it.

### 6. Alternatives considered

**CUE** — modeling the contract in CUE and generating Go + schema from it was
prototyped (see `docs/status/2026-11-02_04-52_cue-adoption-analysis-and-self-review.md`).
Rejected for v4: the struct tags already are the single declaration surface,
and a parallel CUE schema would be a second source of truth — the exact drift
this ADR exists to prevent. Revisit when a consumer needs schema _validation_
of the contract itself, not of config files.

**invopop/jsonschema** — reflection-only, no cmdguard tag knowledge, no env /
required / secret vocabulary. The `schema` sub-module is ~500 lines and owns
the mapping; importing a generic emitter would still require that mapping and
add a dependency for it.

### 7. Golden fixtures are load-bearing

`describe_golden_test.go` and `schema/testdata/config.schema.json` pin the
emitted bytes. A diff means either an intended additive change (regenerate and
review) or nondeterminism (bug). The schema emitter sorts properties
(`PropertyMap`) precisely so map iteration order cannot leak into the golden.
