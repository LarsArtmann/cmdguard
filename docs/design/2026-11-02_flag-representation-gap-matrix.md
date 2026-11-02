# The Four Flag Representations: Field-by-Field Gap Matrix

_Truth source for the v5 consolidation scope (ADR-002 companion). Written
2026-11-02 while closing the config-truth work; read `docs/adr/002-*.md` for
the contract stability rules._

cmdguard describes a flag four times. Each representation exists for a reason;
the risk is drift between them. This matrix records every field, where it
lives, and how it reaches the next layer.

## The four representations

| # | Representation    | File          | Purpose                                                                  |
| - | ----------------- | ------------- | ------------------------------------------------------------------------ |
| 1 | `FlagTag`         | `config.go`   | Parsed struct tags; the declaration-time truth                           |
| 2 | `flagMetadata`    | `describe.go` | The annotation-encoded subset Describe needs; rides on the cobra command |
| 3 | `FlagDescription` | `describe.go` | The public machine contract (`describe` JSON, docgen source)             |
| 4 | `pflag.Flag`      | cobra/pflag   | The runtime representation cobra actually parses                         |

## The matrix

Fields marked "derived" are computed, not declared, and are listed for
completeness.

| Field / fact           | 1 FlagTag           | 2 flagMetadata   | 3 FlagDescription         | 4 pflag.Flag                 | Notes                                                                          |
| ---------------------- | ------------------- | ---------------- | ------------------------- | ---------------------------- | ------------------------------------------------------------------------------ |
| Name (flag spelling)   | `Name`              | `Name`           | `Name`                    | `Name`                       | Flows 1→2→3→4 at registration                                                  |
| Shorthand              | `Short`             | —                | `Shorthand`               | `Shorthand`                  | Not in annotations; read from pflag                                            |
| Help text              | `Help`              | —                | `Usage`                   | `Usage`                      | Written into pflag usage                                                       |
| Default (as tag)       | `Default`           | —                | via `DefValue`            | `DefValue`                   | Config files rewrite the pflag default (`DeclaredDefaults` keeps the original) |
| Type                   | `Type reflect.Type` | —                | `Type` (pflag value type) | `Value.Type()`               | Contract reports pflag's type string, not the Go type                          |
| Enum values            | `Values`            | `Values`         | `Values`                  | —                            | Annotations carry it; pflag has no notion                                      |
| Required               | `Required`          | `Required`       | `Required`                | —                            | Enforced at run time, not by pflag                                             |
| Env var                | `Env`               | `Env` (prefixed) | `Env` (prefixed)          | —                            | Prefix applied before annotation write                                         |
| Counting               | `Count`             | `Count`          | `Count`                   | (value type "count")         |                                                                                |
| Prompt                 | `Prompt`            | `Prompt`         | `Prompt`                  | —                            |                                                                                |
| Validate rule          | `Validate`          | `Validate`       | `Validate`                | —                            | Raw tag text; interpreted by the validator registry                            |
| Hidden                 | `Hidden`            | —                | `Hidden`                  | `Hidden`                     | Set on pflag at registration                                                   |
| Deprecated             | —                   | —                | `Deprecated`              | `Deprecated`                 | pflag-only concept today                                                       |
| Inherited              | —                   | —                | `Inherited`               | (flag set membership)        | Derived by docgen/describe from set placement                                  |
| Local (no inherit)     | `Local`             | —                | —                         | (registration-time behavior) | Wiring decision; intentionally not in the contract                             |
| Secret                 | `Secret`            | `Secret`         | `Secret`                  | —                            | Full chain landed 2026-11-02 (was the gap that started this matrix)            |
| Config key (canonical) | `Key`               | —                | —                         | —                            | Config-file vocabulary lives in `ConfigKeys`, not in the flag contract         |
| Reflect path           | `Field`, `Index`    | —                | —                         | —                            | Declaration-time plumbing                                                      |
| Changed at runtime     | —                   | —                | —                         | `Changed`                    | Runtime state; reset per Execute                                               |
| NoOptDefVal            | —                   | —                | —                         | `NoOptDefVal`                | pflag detail (bool/count behavior)                                             |

## Gaps and why they are acceptable (or not)

1. **`secret` crossed only two layers until 2026-11-02.** `config explain` and
   `ConfigKeys` consumed `FlagTag.Secret`, but the annotation payload and the
   machine contract dropped it, so `describe` and generated docs could not
   report it. Fixed by threading `Secret` through `flagMetadata` →
   `FlagDescription` → the docgen Secret column, with the drift test
   (`docgen_drift_test.go`) enforcing the docgen half.
2. **`Shorthand` and `Usage` are absent from the annotations.** They survive in
   pflag, which Describe reads anyway, so the annotation stays minimal on
   purpose. Do not add fields to `flagMetadata` that pflag already carries.
3. **`Local` and `Key` are deliberately not in the contract.** `Local` is a
   wiring decision (which commands inherit), not a per-flag fact;
   `Key` belongs to the config-file vocabulary (`v4.ConfigKey`), where array
   paths (`servers[0].port`) already express what a flag cannot.
4. **Two default notions.** The tag default (`FlagTag.Default` /
   `DeclaredDefaults()`) is what the developer declared; the pflag default
   (`DefValue`) is what config loading rewrote it into. `config explain`
   reports both (Default = declared, Value = resolved). Any new representation
   must keep the two apart or provenance breaks again.
5. **The contract's `Type` is pflag's string** (`string`, `int`, `bool`,
   `stringSlice`, ...), not the Go field type. Consumers map types themselves;
   the schema emitter instead derives JSON types from the Go type. Unifying
   these is v5 scope.

## The rule going forward

A metadata field that affects a valid invocation (required, values, env,
secret, prompt, validate, count) must reach ALL of: `flagMetadata`
(annotation), `FlagDescription` (contract), and the docgen table — enforced by
`TestGenerateDocs_MatchesDescribeContract` for the last hop. Wiring-only fields
(`Local`) and vocabulary fields (`Key`) stay out of the contract by design.
