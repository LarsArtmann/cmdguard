# Naming Review Report — cmdguard v2.4.0

## Executive Summary

- **50 source files** reviewed, **~250 identifiers** analyzed
- **0 honesty issues** (no lying names found)
- **3 clarity issues** (abbreviations, vagueness)
- **2 domain alignment issues** (inconsistent terminology)
- **1 consistency issue** (synonym drift)
- **0 implementation leakage** (no Impl/Base/I prefixes)

## Clarity Issues

| #     | File                                                                                                               | Line    | Identifier             | Issue                                                  | Better Name                  |
| ----- | ------------------------------------------------------------------------------------------------------------------ | ------- | ---------------------- | ------------------------------------------------------ | ---------------------------- |
| ~~1~~ | ~~errors.go~~ **Won't implement — name kept — labeledError survives unchanged (errors.go:65).**                    | ~~192~~ | ~~`labeledError()`~~   | ~~Not an error, formatting helper~~                    | ~~`formatLabeledMessage()`~~ |
| ~~2~~ | ~~config_parsing.go~~ **Won't implement — name kept — parseFieldFlag survives unchanged (config_parsing.go:124).** | ~~N/A~~ | ~~`parseFieldFlag()`~~ | ~~Ambiguous: parses flag on field or field for flag?~~ | ~~`parseFlagTag()`~~         |
| ~~3~~ | ~~flag_helpers.go~~ **Won't implement — case listing still present in flag_helpers.go.**                           | ~~39~~  | ~~Giant `case` block~~ | ~~Enumerates all Kinds when `default` suffices~~       | ~~Remove explicit listing~~  |

## Domain Alignment Issues

| #     | Files                                                                                                                                                | Concept                           | Current Names                                     | Canonical Name                              |
| ----- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------- | ------------------------------------------------- | ------------------------------------------- |
| ~~1~~ | ~~flags_validate.go, config_file.go~~ **Won't implement — never unified — fieldValueToString (config_file.go:123) + formatFieldValue both survive.** | ~~Convert field value to string~~ | ~~`formatFieldValue`, `fieldValueToString`~~      | ~~`formatFieldValue` (one canonical impl)~~ |
| ~~2~~ | ~~command.go, cli_command.go~~ done at `IsExecutable removed in v2.6.0 (CHANGELOG); HasHandler canonical (command.go:174)`                           | ~~Does command have a handler?~~  | ~~`HasHandler()`, `IsExecutable()` (deprecated)~~ | ~~`HasHandler()` — already canonical~~      |

## Consistency Issues

| #     | Operation                                                                         | Files                                        | Current Verbs                                           | Standardize To                                                                    |
| ----- | --------------------------------------------------------------------------------- | -------------------------------------------- | ------------------------------------------------------- | --------------------------------------------------------------------------------- |
| ~~1~~ | ~~Get flow context~~ done — both accessors documented (flow_context_access.go:33) | ~~flow_context_access.go, cli_accessors.go~~ | ~~`GetBranchingFlowContext(ctx)`, `cli.FlowContext()`~~ | ~~Document that `FlowContext()` returns nil before Execute; `Get*` from context~~ |

## Strengths (Good Naming)

- `CLI[T]`, `Command[T, F]` — honest, precise, domain-aligned
- `BranchingFlowContext` — unique, descriptive, clear purpose
- `FlagRegistry` — precise: it's a registry of flags, not a manager or handler
- `ValidationMode` with `Lenient`/`Strict`/`Draconian` — progressive spectrum, memorable
- `MustNewCommand`/`MustInvoke` — Go convention for panic variants
- `ErrMissingHandler`, `ErrDuplicateCommand` — honest sentinel names
- `OutputFormat` — clean type, not `OutputInfo` or `FormatData`
- `Duration`, `Email`, `Port`, `FilePath`, `HostPort`, `URL`, `LogLevel` — domain types named after what they validate
- `WithShort`, `WithLong`, `WithFlags`, `WithRunE` — consistent `With*` pattern for options
- `ExitCoder` interface — honest about what it does
- `Phase` enum (`PreRun`, `Run`, `PostRun`) — clear lifecycle stages
- `SpinnerConfig` — descriptive struct for spinner configuration
- `EditInEditor()` — honest about mechanism and purpose

## Notable Decisions

- `NoFlags` as `type NoFlags = struct{}` — type alias is intentional for ergonomics. Trade-off: loses type safety but gains API simplicity.
- `CommandOption[T, F]` — verbose due to Go generics, but honest. No alternative exists in current Go.
- `dispatchRegister`/`dispatchParse`/`dispatchDefault` — `dispatch` prefix is clear: runtime type dispatch. Acceptable for internal functions.

## Verdict

**9/10 naming quality.** This is one of the best-named Go codebases reviewed. The three minor issues (labeledError, parseFieldFlag, formatFieldValue duplication) are low-impact. No Manager/Handler/Processor/Helper/Util/Utility classes. No Impl/Base/I prefixes. No lying names. No split-brain terminology.
