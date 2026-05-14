# Execution Plan — cmdguard Modularization

**Date:** 2026-05-14
**Status:** Ready for execution
**Based on:** PROPOSAL.md (post self-review)

---

## Overview

12 tasks, ordered by dependency and Pareto impact. Each task leaves the project in a buildable, testable state and is independently revertable via a single commit.

**Total estimated effort:** 4–6 hours

---

## Impact Tiers

| Tier | Tasks | Impact |
|------|-------|--------|
| 1% → 51% (Foundational) | 1–4 | Create modules, extract types and output |
| 4% → 64% (High leverage) | 5–8 | Wire dependencies, update imports, re-exports |
| 20% → 80% (Broad value) | 9–10 | Update examples and integration tests |
| Remaining (Polish) | 11–12 | Dead code cleanup, final verification |

---

## Pre-Flight

- [ ] Ensure clean git state: `git status`
- [ ] Create branch: `git checkout -b modularize/split-types-output`
- [ ] Run baseline: `go test ./... -count=1 -timeout 120s -race` (confirm 199 tests pass)

---

## Task 1: Create go.work

**What:** Initialize Go workspace at repository root.

**Steps:**
1. Run `go work init`
2. Add root module: `go work use .`
3. Verify: `go work sync`

**Verification:**
```bash
go build ./...
go test ./... -count=1 -timeout 120s -race
```

**Rollback:** Delete `go.work`

**Effort:** 5 min

---

## Task 2: Extract testutil module

**What:** Move `pkg/testutil/` to `/testutil/` with its own go.mod.

**Rationale:** testutil is a leaf dependency with no internal deps. Extracting it first means other modules can depend on it without circular imports.

**Steps:**
1. Create `/testutil/go.mod`:
   ```
   module github.com/larsartmann/cmdguard/testutil
   go 1.26.2
   require github.com/spf13/cobra v1.10.2
   ```
2. Move `pkg/testutil/panic_test_helpers.go` → `/testutil/panic_test_helpers.go`
3. Update package declaration (already `testutil`)
4. Run `go mod tidy` in `/testutil/`
5. Add `go work use ./testutil` to go.work
6. Update all test files that import `github.com/larsartmann/cmdguard/pkg/testutil` → `github.com/larsartmann/cmdguard/testutil`
7. Run `go work sync`

**Verification:**
```bash
cd testutil && go build ./... && go vet ./...
go build ./...
go test ./... -count=1 -timeout 120s -race
```

**Rollback:** Revert commit, move files back

**Effort:** 15 min

---

## Task 3: Extract types module — create structure

**What:** Create `/types/` directory with go.mod and move value type files.

**Steps:**
1. Create `/types/go.mod`:
   ```
   module github.com/larsartmann/cmdguard/types
   go 1.26.2
   ```
   (Zero external dependencies)
2. Copy (not move yet) these files to `/types/`, renaming:
   - `types_duration.go` → `duration.go`
   - `types_email.go` → `email.go`
   - `types_enum.go` → `enum.go`
   - `types_filepath.go` → `filepath.go`
   - `types_hostport.go` → `hostport.go`
   - `types_log.go` → `log.go`
   - `types_port.go` → `port.go`
   - `types_url.go` → `url.go`
3. Copy `type_helpers.go` → `/types/helpers.go` (keep only `MustParse[T]`; remove dead `Ptr`, `ValueOrDefault`, `EnsureValid`)
4. Add error sentinels and types to `/types/errors.go`:
   - `ErrInvalidURL`, `ErrInvalidEmail`, `ErrInvalidPort`, `ErrInvalidFilePath`, `ErrInvalidHostPort`, `ErrInvalidEnum`, `ErrInvalidDuration`
   - `EnumError`, `DurationError` + constructors
5. Update all package declarations to `package types`
6. Update internal imports within types files (e.g., `NewEnumError` → local, `ParseEnum` → local)
7. Copy corresponding test files, update package to `types_test`
8. Run `go mod tidy` in `/types/`
9. Add `go work use ./types` to go.work
10. Verify types module builds independently: `cd types && go build ./... && go vet ./...`

**Verification:**
```bash
cd types && go build ./... && go test ./... -race
go work sync
```

**Rollback:** Delete `/types/` directory, revert go.work

**Effort:** 30 min

---

## Task 4: Extract types module — update core imports

**What:** Update root module to import types from the new module instead of local files.

**Steps:**
1. Add types module to root go.mod:
   ```
   require github.com/larsartmann/cmdguard/types v0.0.0
   ```
   (go.work resolves this locally)
2. Update `type_handler.go` imports: add `"github.com/larsartmann/cmdguard/types"`
3. Update `config_setfield.go` imports for `Enum`, `ParseEnum`, `FromDuration`
4. Update `flags_validate.go` imports for shared error sentinels
5. Update `flags.go` imports for shared error sentinels
6. Delete original `types_*.go` and `type_helpers.go` from `pkg/cmdguard/v2/`
7. Run `go mod tidy` in root
8. Run `go work sync`

**Verification:**
```bash
go build ./...
go test ./... -count=1 -timeout 120s -race
golangci-lint run ./...
```

**Rollback:** Revert commit, restore deleted files

**Effort:** 30 min

**Critical:** This is the highest-risk step. If imports don't resolve, check that:
- `go.work` includes both `.` and `./types`
- Types module's exported symbols match what core expects
- Package name is `types` (not `v2`)

---

## Task 5: Extract types module — move tests

**What:** Move value type test files from `pkg/cmdguard/v2/` to `/types/`.

**Steps:**
1. Move these test files to `/types/`:
   - `duration_test.go`
   - `email_test.go`
   - `enum_test.go`
   - `filepath_test.go`
   - `hostport_test.go`
   - `port_test.go`
   - `url_test.go`
2. Update package declarations to `types_test`
3. Update imports to use `github.com/larsartmann/cmdguard/types`
4. Update testutil imports to `github.com/larsartmann/cmdguard/testutil`
5. Verify types tests pass independently

**Verification:**
```bash
cd types && go test ./... -race
go test ./... -count=1 -timeout 120s -race
```

**Rollback:** Move test files back

**Effort:** 20 min

---

## Task 6: Extract types module — add re-exports in v2

**What:** Add re-export type aliases in `pkg/cmdguard/v2` for backward compatibility.

**Steps:**
1. Create `pkg/cmdguard/v2/types_reexport.go`:
   ```go
   package v2

   import (
       "github.com/larsartmann/cmdguard/types"
   )

   // Value type re-exports for backward compatibility.
   type Duration = types.Duration
   type Email = types.Email
   type Enum = types.Enum
   type FilePath = types.FilePath
   type HostPort = types.HostPort
   type LogLevel = types.LogLevel
   type LogFormat = types.LogFormat
   type Port = types.Port
   type URL = types.URL

   // Re-export constructors.
   var ParseDuration = types.ParseDuration
   var ParseEmail = types.ParseEmail
   var ParseEnum = types.ParseEnum
   var ParseFilePath = types.ParseFilePath
   var ParseHostPort = types.ParseHostPort
   var ParseLogLevel = types.ParseLogLevel
   var ParseLogFormat = types.ParseLogFormat
   var ParsePort = types.ParsePort
   var ParseURL = types.ParseURL
   var MustParseURL = types.MustParseURL
   var MustParseEmail = types.MustParseEmail
   var MustParsePort = types.MustParsePort
   var MustParseFilePath = types.MustParseFilePath
   var MustParseHostPort = types.MustParseHostPort
   var FromDuration = types.FromDuration
   var PortFromInt = types.PortFromInt
   var NewHostPort = types.NewHostPort

   // Re-export error sentinels.
   var ErrInvalidURL = types.ErrInvalidURL
   var ErrInvalidEmail = types.ErrInvalidEmail
   // ... etc
   ```
2. Verify existing consumers can still import from `pkg/cmdguard/v2`

**Verification:**
```bash
go build ./...
go test ./... -count=1 -timeout 120s -race
```

**Rollback:** Delete re-export file

**Effort:** 15 min

---

## Task 7: Extract output module

**What:** Create `/output/` module with output rendering logic.

**Steps:**
1. Create `/output/go.mod`:
   ```
   module github.com/larsartmann/cmdguard/output
   go 1.26.2
   require github.com/larsartmann/go-output v0.2.0
   ```
2. Copy `output.go` to `/output/output.go`
3. Add error sentinels to `/output/errors.go`:
   - `ErrUnsupportedFormat`, `ErrFormatRequiresTypedData`
4. Update package declaration to `package output`
5. Update imports (remove v2-internal refs, keep go-output)
6. Copy output test files, update to `package output_test`
7. Run `go mod tidy` in `/output/`
8. Add `go work use ./output` to go.work
9. Update root go.mod to depend on output module
10. Update `cli_output.go` imports to use output module
11. Delete original `output.go` from `pkg/cmdguard/v2/`
12. Run `go work sync`

**Verification:**
```bash
cd output && go build ./... && go test ./... -race
go build ./...
go test ./... -count=1 -timeout 120s -race
```

**Rollback:** Revert commit, restore files

**Effort:** 25 min

**Note:** `cli_output.go` stays in core. Only `output.go` (rendering) moves.

---

## Task 8: Add output re-exports in v2

**What:** Re-export output types from core for backward compatibility.

**Steps:**
1. Create `pkg/cmdguard/v2/output_reexport.go`:
   ```go
   package v2

   import (
       "github.com/larsartmann/cmdguard/output"
   )

   type OutputFormat = output.OutputFormat
   type OutputConfig = output.OutputConfig

   var ParseOutputFormat = output.ParseOutputFormat
   var DefaultOutputConfig = output.DefaultOutputConfig
   var OutputResult = output.OutputResult
   var OutputTable = output.OutputTable
   var OutputStyledTable = output.OutputStyledTable

   // Format constants
   var FormatTable = output.FormatTable
   var FormatJSON = output.FormatJSON
   // ... etc
   ```

**Verification:**
```bash
go build ./...
go test ./... -count=1 -timeout 120s -race
```

**Rollback:** Delete re-export file

**Effort:** 10 min

---

## Task 9: Update examples

**What:** Update example imports if any reference types directly.

**Steps:**
1. Check which examples import types/output directly from v2
2. If any examples use value types directly, update imports to use either:
   - `github.com/larsartmann/cmdguard/pkg/cmdguard/v2` (re-exports) — **preferred for examples**
   - `github.com/larsartmann/cmdguard/types` — only if demonstrating standalone usage
3. Verify each example builds and runs

**Verification:**
```bash
go build ./examples/...
go test ./examples/... -count=1 -timeout 120s -race
```

**Rollback:** Revert commit

**Effort:** 15 min

**Note:** Most examples use `pkg/cmdguard/v2` which will still work via re-exports. Likely no changes needed.

---

## Task 10: Update integration tests

**What:** Update integration test imports.

**Steps:**
1. Check `tests/integration/*.go` for any direct type/output imports
2. Update if needed (likely no changes — they import `pkg/cmdguard/v2`)

**Verification:**
```bash
go test ./tests/... -count=1 -timeout 120s -race
```

**Rollback:** Revert commit

**Effort:** 10 min

---

## Task 11: Dead code cleanup

**What:** Remove dead code discovered during analysis.

**Steps:**
1. Delete `ErrLogLevel` and `ErrLogFormat` from `errors.go` (never used)
2. Delete `Ptr`, `ValueOrDefault`, `EnsureValid` from `type_helpers.go` (if still in core; otherwise already removed in types module)
3. Delete any orphaned test files that moved to types/output modules

**Verification:**
```bash
go build ./...
go test ./... -count=1 -timeout 120s -race
golangci-lint run ./...
```

**Rollback:** Revert commit

**Effort:** 10 min

---

## Task 12: Final verification

**What:** Full verification pass across all modules.

**Steps:**
1. Run full test suite:
   ```bash
   go work sync
   go test ./... -count=1 -timeout 120s -race
   ```
2. Run per-module verification:
   ```bash
   cd types && go build ./... && go test ./... -race && go vet ./... && go mod tidy && cd ..
   cd output && go build ./... && go test ./... -race && go vet ./... && go mod tidy && cd ..
   cd testutil && go build ./... && go test ./... -race && go vet ./... && go mod tidy && cd ..
   go build ./... && go test ./... -race && go vet ./... && go mod tidy
   ```
3. Run lint:
   ```bash
   golangci-lint run ./...
   ```
4. Verify go.work is minimal:
   ```bash
   go work sync
   cat go.work
   ```
5. Verify each module's go.mod is clean (no replace directives):
   ```bash
   grep -r "replace" */go.mod go.mod
   ```

**Verification:** All commands pass with zero errors.

**Rollback:** N/A (verification only)

**Effort:** 15 min

---

## Dependency Graph Between Tasks

```
Task 1 (go.work)
  │
  ├──► Task 2 (testutil)
  │      │
  │      ├──► Task 3 (types structure)
  │      │      │
  │      │      └──► Task 4 (types core imports)
  │      │             │
  │      │             ├──► Task 5 (types tests)
  │      │             │
  │      │             └──► Task 6 (types re-exports)
  │      │
  │      └──► Task 7 (output module)
  │             │
  │             └──► Task 8 (output re-exports)
  │
  └──► Task 9 (update examples)
         │
         ├──► Task 10 (update integration tests)
         │
         ├──► Task 11 (dead code cleanup)
         │
         └──► Task 12 (final verification)
```

**Parallelizable:** Tasks 3–6 (types chain) and Tasks 7–8 (output chain) can be done in parallel.

---

## Per-Module go.mod Templates

### `/types/go.mod`
```
module github.com/larsartmann/cmdguard/types

go 1.26.2
```

### `/output/go.mod`
```
module github.com/larsartmann/cmdguard/output

go 1.26.2

require github.com/larsartmann/go-output v0.2.0
```

### `/testutil/go.mod`
```
module github.com/larsartmann/cmdguard/testutil

go 1.26.2

require github.com/spf13/cobra v1.10.2
```

### `/go.mod` (root, updated)
```
module github.com/larsartmann/cmdguard

go 1.26.2

require (
    charm.land/fang/v2 v2.0.1
    github.com/larsartmann/cmdguard/types v0.0.0
    github.com/larsartmann/cmdguard/output v0.0.0
    github.com/muesli/mango v0.2.0
    github.com/muesli/mango-cobra v1.3.0
    github.com/muesli/roff v0.1.0
    github.com/samber/do/v2 v2.0.0
    github.com/spf13/cobra v1.10.2
    github.com/spf13/pflag v1.0.10
)
```

### `go.work`
```
go 1.26.2

use (
    .
    ./types
    ./output
    ./testutil
)
```

---

## Commit Messages (Template)

```
Task 1: chore(workspace): initialize go.work for multi-module development

Task 2: refactor(testutil): extract testutil to standalone module

Task 3: refactor(types): create types module with value types

Task 4: refactor(types): update core to import from types module

Task 5: refactor(types): move type tests to types module

Task 6: feat(types): add backward-compatible re-exports in v2 package

Task 7: refactor(output): extract output rendering to standalone module

Task 8: feat(output): add backward-compatible re-exports in v2 package

Task 9: chore(examples): update imports for modularized structure

Task 10: chore(tests): update integration test imports

Task 11: chore(cleanup): remove dead code (ErrLogLevel, ErrLogFormat, unused helpers)

Task 12: chore(verify): final verification pass across all modules
```
