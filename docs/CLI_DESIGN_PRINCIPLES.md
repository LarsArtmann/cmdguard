# CLI Design Principles - Lessons from go-structure-linter

## The Problem

The go-structure-linter CLI demonstrates multiple UX failures:

```bash
# User tries these - ALL fail with confusing errors:
$ go-structure-linter . --strict-mode
ERROR: Flag needs an argument: --strict-mode

$ go-structure-linter . --severity ALL
ERROR: Unknown flag: --severity

$ go-structure-linter . --severity-min
ERROR: Flag needs an argument: --severity-min
```

## Root Causes

### 1. Inconsistent Naming Conventions

| What User Expected | What Exists | Problem |
|-------------------|-------------|---------|
| `--severity` | `--severity-min` | Different naming pattern |
| `--strict` | `--strict-mode` | Unnecessary suffix |
| `--fix` (no args) | `--fix` (works) | Inconsistent behavior |

**Rule:** Use consistent naming. If `--fix` is boolean, `--strict` should be too.

### 2. Ambiguous Flag Values

```go
// Help shows:
--strict-mode         Strict mode (disabled, enforced, warn)

// User thinks:
--strict-mode          // Just enable it

// Actually required:
--strict-mode=enforced
```

**Rule:** Boolean flags should be boolean. Enum flags need clear syntax.

### 3. Poor Error Messages

| Current Error | Better Error |
|--------------|--------------|
| `Unknown flag: --severity` | `Unknown flag: --severity. Did you mean --severity-min?` |
| `Flag needs an argument` | `--strict-mode requires one of: disabled, enforced, warn` |

**Rule:** Error messages should suggest fixes.

### 4. Help/Examples Mismatch

```
EXAMPLES
  go-structure-linter --severity-min CRITICAL  # Works
  go-structure-linter --severity ALL           # Doesn't work!
```

**Rule:** Every example in help must be copy-pasteable and work.

---

## cmdguard Solutions

### Principle 1: Consistent Flag Patterns

```go
// ✅ GOOD: Boolean flags use BoolP
cmd.PersistentFlags().BoolP("strict", "s", false, "Enable strict mode")

// ❌ BAD: String flag that looks like boolean
cmd.PersistentFlags().String("strict-mode", "disabled", "Strict mode (disabled, enforced, warn)")
```

### Principle 2: Enum Flags Use Validated Strings

```go
// ✅ GOOD: Pre-validate enum values
cmd.PersistentFlags().String("severity", "INFO", "Minimum severity (DEBUG, INFO, WARN, ERROR)")

// In validation:
func validateSeverity(value string) error {
    allowed := []string{"DEBUG", "INFO", "WARN", "ERROR"}
    if !slices.Contains(allowed, value) {
        return fmt.Errorf("invalid severity %q, must be one of: %v", value, allowed)
    }
    return nil
}
```

### Principle 3: Clear Error Messages

```go
// ✅ GOOD: Suggest alternatives
if err := cmd.Execute(); err != nil {
    if strings.Contains(err.Error(), "unknown flag") {
        return suggestSimilarFlags(err, cmd)
    }
    return err
}
```

### Principle 4: Help Examples Must Work

```go
// ✅ GOOD: Programmatically verify examples
type Example struct {
    Command string
    Valid   bool // Set to true only if tested
}

var examples = []Example{
    {Command: "cmdguard validate --strict", Valid: true},
    {Command: "cmdguard --log-level debug", Valid: true},
}
```

---

## Recommended Flag Design

### Severity Levels

```bash
# Option 1: Enum with clear values (RECOMMENDED)
--severity=INFO        # Default
--severity=DEBUG       # Most verbose
--severity=ERROR       # Only errors

# Option 2: Boolean flags (simpler)
--verbose              # DEBUG level
--quiet                # ERROR level
```

### Mode Flags

```bash
# ❌ BAD: String flag with confusing syntax
--strict-mode=enforced

// ✅ GOOD: Boolean flag
--strict              # Enable strict mode
--strict=false        # Explicitly disable
```

### List Values

```bash
# ❌ BAD: Space-separated
--exclude "foo bar"

// ✅ GOOD: Repeated flags
--exclude foo --exclude bar

// ✅ GOOD: Comma-separated with clear docs
--exclude "foo,bar"   # Docs: "Comma-separated list"
```

---

## Validation Checklist

Before releasing a CLI tool, verify:

- [ ] Every example in `--help` works when copy-pasted
- [ ] Boolean flags don't require values
- [ ] Enum values are validated and suggested on error
- [ ] Error messages suggest the correct flag name
- [ ] Short flags exist for common long flags
- [ ] Flag names follow consistent pattern (kebab-case)
- [ ] Default values are shown in help
- [ ] Required flags are marked clearly

---

## Implementation in cmdguard

### Current Issues

```go
// internal/commands/root.go:43
root.PersistentFlags().Bool("strict", false, "enable strict mode validation")

// ✅ This is correct! Boolean flag, clear name.
```

### Recommended Improvements

1. **Add short flags** for common options
2. **Validate enum values** for --log-level
3. **Show defaults** in help text
4. **Add flag suggestions** on unknown flag errors
