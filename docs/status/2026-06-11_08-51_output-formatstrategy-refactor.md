# Output FormatStrategy Refactor — Status Report

**Date:** 2026-06-11 08:51
**Task:** Replace dual-registry (tableFormatRegistry + anyFormatRegistry) with FormatStrategy interface + unified formatRegistry
**Status:** COMPLETE — all tests pass, 0 lint issues, 0 race conditions

---

## a) FULLY DONE

### Core Refactor (output.go)

| What             | Before                                                                            | After                                                      |
| ---------------- | --------------------------------------------------------------------------------- | ---------------------------------------------------------- |
| Dispatch model   | Type switch → 2 separate registries                                               | Single `formatRegistry` lookup → `FormatStrategy.Render()` |
| Type definitions | `tableRenderer`, `anyRenderer`, `renderMarshalFunc`, `renderStringFunc` (4 types) | `FormatStrategy` interface + 5 strategy structs            |
| Helper functions | `renderAndWrite`, `marshalAndWrite`, `renderTableData`, `renderAny` (4 funcs)     | `unwrapTableData` (1 helper)                               |
| Registry entries | `tableFormatRegistry` (15) + `anyFormatRegistry` (3) = 18 entries                 | `formatRegistry` (16 entries, 3 dual)                      |
| Co-located logic | JSON/YAML/TOML table+any split across 2 registries                                | JSON/YAML/TOML each in 1 `dualStrategy` entry              |

### Strategy Types Introduced

| Type                         | Purpose                                                                                              |
| ---------------------------- | ---------------------------------------------------------------------------------------------------- |
| `FormatStrategy` (interface) | `Render(w io.Writer, data any) error` — single contract for all formats                              |
| `tableRenderStrategy`        | TableData-only formats (TSV, XML, Markdown, HTML, Tree, D2, Mermaid, DOT, JSONL, AsciiDoc, PlantUML) |
| `marshalStrategy`            | Any-data serialization (YAML, TOML any-path)                                                         |
| `dualStrategy`               | TableData + any-data dispatch (JSON, TOML, and future dual formats)                                  |
| `styledTableStrategy`        | Lipgloss-styled terminal table                                                                       |
| `csvStrategy`                | Streaming CSV via delimited writer                                                                   |

### Tests Added (output_test.go)

13 new test functions:

| Test                                              | Coverage                               |
| ------------------------------------------------- | -------------------------------------- |
| `TestUnwrapTableData`                             | Pointer, value, string, nil, int cases |
| `TestFormatRegistry_Completeness`                 | All 16 formats registered              |
| `TestTableRenderStrategy_RejectsNonTableData`     | `ErrFormatRequiresTypedData` chain     |
| `TestMarshalStrategy_RendersAnyData`              | Arbitrary struct marshaling            |
| `TestMarshalStrategy_RendersTableData`            | TableData passthrough                  |
| `TestDualStrategy_DelegatesToTable`               | TableData path selected                |
| `TestDualStrategy_DelegatesToAny`                 | Any-data path selected                 |
| `TestStyledTableStrategy_RejectsNonTableData`     | Error chain                            |
| `TestCSVStrategy_RejectsNonTableData`             | Error chain                            |
| `TestOutputResult_UnsupportedFormat`              | `ErrUnsupportedFormat` chain           |
| `TestOutputResult_AnyData_YAML/TOML/JSON`         | 3 dual formats with arbitrary data     |
| `TestOutputResult_TableOnlyFormats_RejectAnyData` | 13 formats reject non-TableData        |
| `TestFormatStrategy_Interface`                    | Compile-time interface compliance      |

### Metrics

| Metric                  | Before | After |
| ----------------------- | ------ | ----- |
| Tests passing           | 395+   | 410+  |
| Coverage                | 85.5%  | 85.5% |
| Lint issues             | 0      | 0     |
| Race conditions         | 0      | 0     |
| Lines in output.go      | 297    | 309   |
| Lines in output_test.go | 244    | 485   |

---

## b) PARTIALLY DONE

Nothing partially done. The refactor is complete and verified.

---

## c) NOT STARTED

See "Top #25 Next Steps" below.

---

## d) TOTALLY FUCKED UP

Nothing. The refactor is clean, all tests pass, no regressions.

---

## e) WHAT WE SHOULD IMPROVE — Self-Critique

### Critical Issues with Current Refactor

1. **`styledTableStrategy` and `csvStrategy` duplicate `tableRenderStrategy` pattern**
   Both do `unwrapTableData → error → delegate`. They should use `tableRenderStrategy` directly — `styledTableStrategy` = `tableRenderStrategy{label: "table", render: funcFrom(renderTableStyled)}`, `csvStrategy` = `tableRenderStrategy{label: "csv", render: funcFrom(renderTableCSV)}`. The standalone structs add zero value and add 2 types to the API surface for no reason.

2. **`marshalStrategy` silently accepts TableData but doesn't use TableData-aware marshalers**
   When YAML `marshalStrategy` receives TableData, it calls `serialization.MarshalYAML(data)` which will serialize the struct fields, NOT the tabular data. The old code had the same bug for YAML (it was in `anyFormatRegistry`), but now it's more visible because `marshalStrategy` is documented as handling "any data including TableData".

3. **`fmt.Fprintln(w, result)` adds a trailing newline unconditionally**
   Every strategy adds `\n` after the rendered output. If the rendered output already ends with `\n`, we get double newlines. This was inherited from the old code but should be fixed. The old code had the same bug.

4. **No `RegisterFormatStrategy()` API exposed**
   The whole point of the strategy pattern is extensibility, but `formatRegistry` is a private `var`. Users cannot add custom format strategies. We should expose `RegisterFormatStrategy(format OutputFormat, strategy FormatStrategy)` similar to how `RegisterTypeHandler` works.

5. **YAML `marshalStrategy` handles TableData as raw struct, not tabular data**
   YAML entry is `&marshalStrategy{label: "YAML", marshal: serialization.MarshalYAML}` — but `MarshalYAML` doesn't know about `TableData`. If someone passes TableData to YAML format, they get struct field serialization, not tabular YAML. It should be a `dualStrategy` like JSON and TOML.

6. **`OutputStyledTable` is a dead API — should be deprecated**
   It bypasses the registry entirely, writes directly to stdout, and is superseded by `OutputResult` with `FormatTable`. It's confusing to have two ways to render styled tables.

7. **`renderTableStyled` and `renderTableCSV` are still free functions**
   These should be methods on the strategy structs or private helpers, not package-level functions. They're only used by their respective strategies.

8. **No `RenderOptions` support**
   go-output v0.8.0 has `RenderOptions{Title, GraphID, Writer, ColorMode}` but our `FormatStrategy` interface only takes `(io.Writer, any)`. Title for HTML, GraphID for DOT — these are lost.

### Architectural Issues Worth Noting

9. **We're reinventing what go-output could provide**
   go-output already has `RenderTableData()` with a registry pattern (`RegisterTableDataMarshaler`). We could contribute more marshalers upstream and just call `output.RenderTableData()` instead of maintaining our own registry. The go-output registry only has markdown+tree registered, but it has the infrastructure for all formats.

10. **go-output `allDataMarshaler` concept missing**
    go-output has no registry for arbitrary data (non-TableData). We need this for JSON/YAML/TOML any-data paths. This could be a useful upstream contribution.

---

## f) Top #25 Things We Should Get Done Next

Sorted by Impact × Effort (Pareto order):

| #  | Task                                                                                                                        | Impact | Effort | Category |
| -- | --------------------------------------------------------------------------------------------------------------------------- | ------ | ------ | -------- |
| 1  | Fix YAML to be `dualStrategy` (TableData+any) — currently `marshalStrategy` silently serializes struct fields for TableData | H      | S      | Bug fix  |
| 2  | Eliminate `styledTableStrategy`/`csvStrategy` — fold into `tableRenderStrategy`                                             | M      | S      | Cleanup  |
| 3  | Expose `RegisterFormatStrategy(format, strategy)` for user extensibility                                                    | H      | S      | Feature  |
| 4  | Add `UnregisterFormatStrategy(format)` or `ReplaceFormatStrategy(format, strategy)` for testing                             | M      | S      | Feature  |
| 5  | Fix trailing newline double-\n issue in all strategies                                                                      | L      | S      | Bug fix  |
| 6  | Deprecate `OutputStyledTable` — point users to `OutputResult(OutputConfig{Format: FormatTable}, data)`                      | M      | S      | Cleanup  |
| 7  | Add `FormatStrategy` to doc.go package docs as extension point                                                              | M      | S      | Docs     |
| 8  | Add `RegisterFormatStrategy` tests: custom strategy, override existing, unknown format                                      | H      | S      | Tests    |
| 9  | Fold `renderTableStyled`/`renderTableCSV` into their strategy structs as private helpers                                    | L      | S      | Cleanup  |
| 10 | Contribute TableDataMarshaler implementations upstream to go-output (CSV, JSON, YAML, XML, HTML, etc.)                      | H      | M      | Upstream |
| 11 | Propose `RegisterAnyDataMarshaler` to go-output for JSON/YAML/TOML any-data paths                                           | H      | M      | Upstream |
| 12 | Update AGENTS.md: FormatStrategy section, new gotchas, RegisterFormatStrategy                                               | M      | S      | Docs     |
| 13 | Update TODO_LIST.md: mark output refactor done, add remaining items                                                         | M      | S      | Docs     |
| 14 | Add `SupportedFormats()` function that returns `formatRegistry` keys for CLI help text                                      | M      | S      | Feature  |
| 15 | Add `IsFormatSupported(format) bool` helper                                                                                 | L      | S      | Feature  |
| 16 | Consider `RenderOptions` passthrough: Title for HTML, GraphID for DOT                                                       | M      | M      | Feature  |
| 17 | Benchmark: measure registry lookup vs old type-switch for 16 formats                                                        | L      | S      | Perf     |
| 18 | Consider `sync.RWMutex` for `formatRegistry` if `RegisterFormatStrategy` is called concurrently                             | M      | S      | Safety   |
| 19 | Example_test.go: add FormatStrategy example showing custom format registration                                              | M      | S      | Docs     |
| 20 | Audit all `label` strings in strategies for consistency (some use lowercase "markdown", others "JSON")                      | L      | S      | Cleanup  |
| 21 | Consider `Validate()` on `FormatStrategy` interface for registration-time validation                                        | L      | M      | Feature  |
| 22 | Update FEATURES.md: add FormatStrategy interface, RegisterFormatStrategy                                                    | M      | S      | Docs     |
| 23 | Consider typed generics: `TypedFormatStrategy[T any]` for type-safe renderers                                               | L      | M      | Future   |
| 24 | Investigate go-output v0.8.0 `shape.go` for structured data rendering beyond TableData                                      | M      | M      | Research |
| 25 | Consider `formatStrategyName(strategy) string` for error messages instead of label duplication                              | L      | S      | Cleanup  |

---

## g) Top #1 Question I Cannot Figure Out Myself

**Should we invest in upstreaming format dispatch to go-output, or keep our own registry?**

go-output v0.8.0 has the infrastructure (`RegisterTableDataMarshaler`, `RenderTableData`, `RenderOptions`) but only registers markdown+tree. If we contribute the remaining 14 TableData marshalers upstream AND propose a `RegisterAnyDataMarshaler` for JSON/YAML/TOML, cmdguard's output.go could shrink to ~50 lines (just wrapping go-output's dispatch + our `FormatStrategy` adapter). But this is a significant upstream contribution that requires coordination. **Is this direction desired, or should we keep the cmdguard-specific registry?**

---

## Commit Plan

1. Commit the output.go refactor + tests (current working state)
2. Then proceed with items #1-3 from the top 25 list

## Resolution (2026-07-18)

Superseded — this design was replaced the same day. The `FormatStrategy` interface + 5 strategy structs + cmdguard-side `formatRegistry` described here were discarded in favor of delegating to go-output's upstream `RenderTableData`/`RenderAnyData` registries (see `2026-06-11_09-15_output-registry-delegation-complete.md`). In v3.0.0, `output.go` no longer defines these types. The §f.10/11 upstream-contribution proposals shipped (`AnyDataMarshaler` registry and 16 `TableDataMarshaler` registrations were added to go-output).
