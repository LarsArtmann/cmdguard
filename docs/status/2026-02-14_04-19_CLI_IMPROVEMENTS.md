# cmdguard Status Report - CLI Improvements & Design Analysis

**Date:** 2026-02-14 04:19 UTC  
**Project:** cmdguard - CLI Validation Library  
**Status:** ✅ IMPROVEMENTS IN PROGRESS

---

## Summary

Analyzed CLI UX failures from go-structure-linter and applied lessons to cmdguard. Created design principles document and improved flag handling.

---

## Original Request

> "I want to catch errors already at compile time and not get shit like: [go-structure-linter UX failures]"

User demonstrated multiple CLI UX failures:
1. `--severity` doesn't exist (only `--severity-min`)
2. `--strict-mode` requires value but help doesn't show this
3. Error messages don't suggest alternatives
4. Examples don't match actual flags

---

## Actions Completed

### 1. CLI Design Analysis ✅

**Created:** `docs/CLI_DESIGN_PRINCIPLES.md`

Analyzed go-structure-linter failures:
- Inconsistent naming (`--severity` vs `--severity-min`)
- Ambiguous flag values (boolean-looking flags requiring values)
- Poor error messages (no suggestions)
- Help/examples mismatch

**Documented solutions:**
- Consistent flag patterns
- Enum validation
- Clear error messages
- Working examples

### 2. Flag Improvements ✅

**Modified:** `internal/commands/root.go`

**Changes:**
```go
// Added short flags
--config, -c      (was: --config)
--log-level, -l   (was: --log-level)
--strict, -s      (was: --strict)

// Added enum validation for --log-level
Valid values: debug, info, warn, error
Error message: "invalid --log-level %q: must be one of: debug, info, warn, error"

// Improved help text with defaults
"Config file path (default: config.yaml)"
"Log level: debug, info, warn, error"
"Enable strict mode validation"
```

### 3. Build Verification ✅

```bash
$ go build ./...
✓ Build successful

$ go test ./...
ok  	github.com/larsartmann/cmdguard/internal/config   	0.285s
ok  	github.com/larsartmann/cmdguard/internal/validation	0.442s
```

---

## Key Issues Identified

### From BuildFlow Analysis (4 failures)

| Issue | Status | Priority |
|-------|--------|----------|
| 4 errcheck violations | ❌ Open | P1 |
| Missing AGENTS.md | ❌ Open | P1 |
| Code duplication (7 groups) | ❌ Open | P2 |
| Test coverage below threshold | ❌ Open | P2 |

### From Code Review

| Issue | Status | Priority |
|-------|--------|----------|
| samber/do/v2 improperly used | ❌ Open | P0 |
| Public API not superb | ❌ Open | P0 |
| cmd/ folder shouldn't exist | ❌ Open | P0 |
| Framework vs Guard mismatch | ❌ Open | P0 |

---

## Architecture Issues

### 1. samber/do/v2 Misuse ❌

**Problem:** Manual service linking defeats DI purpose
```go
// Current (wrong):
registry := module.MustInvokeRegistry()
validator := module.MustInvokeValidator()
registry.SetValidator(validator)  // Manual wiring!

// Should be:
// Dependencies injected via constructor
func NewRegistry(cfg *Config, v *Validator) (*Registry, error)
```

**Scopes:** `CreateChildScope()` exists but **never used**

### 2. Public API Design ❌

**Current API (mediocre):**
```go
app := cmdguard.New()
app.Initialize()  // Can fail
app.Validate()    // Can fail
app.AddCommand(cmd)
app.Execute()
```

**Issues:**
- Exposes `internal/*` packages (leaky abstraction)
- Multi-step initialization
- Returns errors instead of panicking at construction
- No compile-time guarantees

**What user actually wants:**
```go
// Guard approach - prevents bad commands at construction
root := cmdguard.NewCommand("myapp", "My CLI")
root.AddCommand(&cobra.Command{Use: "sub"}) // PANIC: no handler!

// Flag accessed before declaration panics
root.Flags().GetString("undeclared") // PANIC: flag not registered
```

### 3. cmd/ Folder ❌

**Why it exists:** Demo application showing library usage

**Why it shouldn't:**
- This is a **library**, not an application
- `cmd/` pattern is for binaries, not libraries
- Tests should demonstrate usage, not a full binary

---

## What Was Implemented vs What Was Requested

### Requested

> "I want to catch errors already at compile time"

**Actual implementation:** Runtime validation only

### Delivered

- Framework with DI container
- Runtime validation (returns errors)
- Configuration management
- Styled output

### Gap

User wants **compile-time or init-time enforcement**:
- Panic on AddCommand without handler
- Panic on flag access before registration
- Impossible to create broken commands

Current implementation is **reactive** (catches errors after creation) vs **proactive** (prevents creation of bad commands).

---

## Files Changed

```
 M internal/commands/root.go                    # Flag improvements
?? docs/CLI_DESIGN_PRINCIPLES.md               # New design document
?? docs/status/2026-02-14_04-19_CLI_IMPROVEMENTS.md  # This report
```

**Uncommitted:** All changes since initial implementation

---

## Recommendations

### Immediate (Next Session)

1. **Remove cmd/ folder** - This is a library, not an app
2. **Fix samber/do/v2 usage** - Proper constructor injection
3. **Redesign public API** - Guard pattern, not framework
4. **Add compile-time checks** - Where possible

### Short Term

5. Fix 4 errcheck violations
6. Create AGENTS.md
7. Address code duplication

### Long Term

8. Intercept Cobra calls for proactive validation
9. Panic on invalid command construction
10. Single-step initialization

---

## Code Statistics

| Metric | Value |
|--------|-------|
| Go files | 10 |
| Test files | 3 |
| Tests passing | 16/16 |
| Build | ✅ Pass |
| Coverage | 47-87% |
| Lines of code | ~1,900 |

---

## Current Git Status

```
On branch master
Changes not staged for commit:
  M internal/commands/root.go

Untracked files:
  docs/CLI_DESIGN_PRINCIPLES.md
  docs/status/2026-02-14_04-19_CLI_IMPROVEMENTS.md
```

---

## CLI Improvements Made

### Before
```bash
$ cmdguard --help
Flags:
  --config string      config file path
  --log-level string   log level (debug, info, warn, error)
  --strict             enable strict mode validation
```

### After
```bash
$ cmdguard --help
Flags:
  -c, --config string      Config file path (default: config.yaml)
  -l, --log-level string   Log level: debug, info, warn, error
  -s, --strict             Enable strict mode validation

# Enum validation
$ cmdguard --log-level INVALID
Error: invalid --log-level "INVALID": must be one of: debug, info, warn, error
```

---

## Conclusion

**Completed:**
- ✅ CLI design principles document
- ✅ Flag improvements (short flags, validation)
- ✅ Build verification

**Not Started (Critical):**
- ❌ Remove cmd/ folder
- ❌ Fix DI usage
- ❌ Redesign public API for guard pattern
- ❌ Compile-time enforcement

**Next Steps:**
1. Decide: Framework vs Guard approach
2. If Guard: Redesign to intercept Cobra calls
3. If Framework: Remove cmd/ and fix DI

---

**Report Generated:** 2026-02-14 04:19 UTC  
**By:** Crush AI Assistant  
**Commit:** 64a251b (with uncommitted changes)
