# BDD Tests Review - cmdguard

**Date:** 2026-03-27
**Reviewer:** AI Assistant
**Status:** CRITICAL REVIEW REQUIRED

---

## Executive Summary

The cmdguard test suite consists of **56 test files** using standard Go testing (`testing` package). While the tests provide good coverage of the API surface, they are **NOT written in BDD style** and do not consistently test from the end-user perspective.

### Critical Finding: Ginkgo Not Used

Despite `AGENTS.md` stating:

> Use Ginkgo/Gomega for BDD-style tests

The project **does not use Ginkgo**. There are no Ginkgo imports, no `Describe`/`It`/`Context` blocks, and Ginkgo is not in `go.mod`.

---

## Test Inventory

| Package | Test Files | Focus Area |
|---------|------------|------------|
| `pkg/cmdguard/v2/` | 34 | Core v2 API |
| `pkg/cmdguard/` | 3 | v1 API |
| `internal/config/` | 4 | Configuration |
| `internal/logging/` | 4 | Logging |
| `tests/integration/` | 5 | E2E integration |
| `examples/*/` | 6 | Example tests |
| **Total** | **56** | |

---

## Current Test Patterns

### What's Working Well

| Aspect | Evidence | Quality |
|--------|----------|---------|
| Helper functions | `newTestCmd()`, `captureOutput()`, `runCLIWithArgs()` | Good |
| Table-driven tests | `TestParseLevel`, `TestSuggestFlag`, `TestEditDistance` | Good |
| Subprocess testing | `TestGuardedCommand_ExecuteAndExit` | Excellent |
| Edge case coverage | nil scope, empty names, duplicate commands | Good |
| Integration tests | `tests/integration/v2_mixed_flags_*` | Good |
| Parallel tests | `t.Parallel()` in validation tests | Good |

### Test Structure Analysis

```
Standard Pattern Used:
├── func TestFeature(t *testing.T)
│   ├── t.Run("subtest name", func(t *testing.T)
│   │   ├── Setup
│   │   ├── Execute
│   │   └── Assert
```

---

## BDD Assessment

### What's Missing (End-User Perspective)

| Category | Current State | BDD Requirement | Gap |
|----------|---------------|-----------------|-----|
| User scenarios | API contracts | "As a user, I want..." | Critical |
| Behavioral descriptions | Function names | Given/When/Then | Critical |
| Acceptance criteria | Implicit | Explicit | High |
| User journey tests | Per-function | End-to-end flows | Medium |
| Help text validation | Limited | Comprehensive | Medium |
| Error UX | Error exists | User-friendly messages | Medium |

### Examples of Missing BDD Scenarios

#### 1. User Discovers Commands
```go
// MISSING: As a user, I want to see what commands are available
func TestUserCanDiscoverCommands(t *testing.T) {
    // Given: A new user runs the CLI
    // When: They run "myapp --help"
    // Then: They see a list of available commands with descriptions
}
```

#### 2. User Gets Helpful Error Messages
```go
// MISSING: As a user, I want helpful suggestions when I mistype a flag
func TestUserGetsTypoSuggestions(t *testing.T) {
    // Given: A user runs "myapp greet --naem Alice"
    // When: The flag is unrecognized
    // Then: They see "Unknown flag: --naem. Did you mean --name?"
}
```

#### 3. User Understands Flag Defaults
```go
// MISSING: As a user, I want to know what the default values are
func TestUserCanSeeFlagDefaults(t *testing.T) {
    // Given: A user runs "myapp greet --help"
    // When: They view the help
    // Then: They see default values like "(default: World)"
}
```

#### 4. User Recovers from Errors
```go
// MISSING: As a user, I want clear guidance when validation fails
func TestUserGetsValidationGuidance(t *testing.T) {
    // Given: A user provides invalid input
    // When: Validation fails
    // Then: They see what was expected vs what was provided
}
```

---

## Quality Metrics

### Test Coverage by Category

| Category | Coverage | Notes |
|----------|----------|-------|
| Unit tests | High | Good API surface coverage |
| Integration tests | Medium | 5 files, good scenarios |
| E2E/Acceptance | Low | Missing user journeys |
| Error path tests | Medium | Error checks exist but shallow |
| Edge cases | High | Good boundary testing |

### Test Naming Convention

**Current:** `TestFunctionName_Scenario`
```go
TestNew_CreatesGuardedCommand
TestCommand_Validate_ErrorEmptyUseField
```

**BDD Style (Missing):** `TestFeature_Scenario_Outcome`
```go
TestUserCanCreateCLI_WithValidConfig_Succeeds
TestUserCanGreet_WithCustomName_ReturnsPersonalizedGreeting
```

---

## Recommendations

### Priority 1: Adopt Ginkgo (As Documented)

The `AGENTS.md` explicitly states to use Ginkgo. Either:
1. **Add Ginkgo** and migrate tests, OR
2. **Update AGENTS.md** to reflect the actual testing approach

**If adopting Ginkgo:**
```go
package v2_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

var _ = Describe("CLI User", func() {
    Context("when creating a new CLI", func() {
        It("should succeed with valid configuration", func() {
            cli, err := v2.New[AppConfig, v2.NoFlags]("myapp", "My App", AppConfig{})
            Expect(err).ToNot(HaveOccurred())
            Expect(cli).ToNot(BeNil())
        })

        It("should fail with empty name", func() {
            _, err := v2.New[AppConfig, v2.NoFlags]("", "My App", AppConfig{})
            Expect(err).To(HaveOccurred())
            Expect(errors.Is(err, v2.ErrInvalidCommand)).To(BeTrue())
        })
    })
})
```

### Priority 2: Add User-Perspective Tests

Create `tests/acceptance/` directory with tests written from user's viewpoint:

```go
// tests/acceptance/user_can_greet_test.go
func TestUserCanGreetSomeone(t *testing.T) {
    // As a user
    // I want to greet someone by name
    // So that I can personalize my CLI interactions

    t.Run("user greets with default name", func(t *testing.T) {
        // Given: A CLI with a greet command
        cli := setupCLI(t)

        // When: User runs "myapp greet"
        output := executeCLI(t, cli, "greet")

        // Then: User sees "Hello, World!"
        if !strings.Contains(output, "Hello, World!") {
            t.Errorf("Expected greeting not found in output: %s", output)
        }
    })

    t.Run("user greets with custom name", func(t *testing.T) {
        // Given: A CLI with a greet command
        cli := setupCLI(t)

        // When: User runs "myapp greet --name Alice"
        output := executeCLI(t, cli, "greet", "--name", "Alice")

        // Then: User sees "Hello, Alice!"
        if !strings.Contains(output, "Hello, Alice!") {
            t.Errorf("Expected personalized greeting not found: %s", output)
        }
    })
})
```

### Priority 3: Add Help Text Tests

```go
func TestUserCanUnderstandHelp(t *testing.T) {
    t.Run("help shows all available flags", func(t *testing.T) {
        output := executeCLI(t, cli, "--help")

        // User should see all flags they can use
        mustContain := []string{"--verbose", "--level", "--help"}
        for _, flag := range mustContain {
            if !strings.Contains(output, flag) {
                t.Errorf("Help missing flag: %s", flag)
            }
        }
    })

    t.Run("help shows default values", func(t *testing.T) {
        output := executeCLI(t, cli, "greet", "--help")

        // User should understand what happens if they don't provide a flag
        if !strings.Contains(output, "default") {
            t.Error("Help should show default values")
        }
    })
})
```

### Priority 4: Add Error UX Tests

```go
func TestUserGetsHelpfulErrors(t *testing.T) {
    t.Run("typo suggests correct flag", func(t *testing.T) {
        output, err := executeCLIWithError(t, cli, "greet", "--naem", "Alice")

        // User should see what went wrong
        if !strings.Contains(err.Error(), "unknown flag") {
            t.Error("Error should identify the unknown flag")
        }

        // User should see a suggestion
        if !strings.Contains(err.Error(), "did you mean") {
            t.Error("Error should suggest correct flag name")
        }
    })
})
```

---

## Action Items

### Immediate (This Sprint)

- [ ] Decide: Use Ginkgo OR update AGENTS.md
- [ ] Add `tests/acceptance/` directory
- [ ] Write 3-5 user-journey acceptance tests

### Short Term (Next Sprint)

- [ ] Add help text validation tests
- [ ] Add error message UX tests
- [ ] Document test naming conventions

### Long Term

- [ ] Migrate existing tests to BDD style (if Ginkgo adopted)
- [ ] Add behavior documentation to tests
- [ ] Create test coverage dashboard for user scenarios

---

## Conclusion

The current test suite provides **good technical coverage** but **lacks end-user perspective**. Tests verify that functions work correctly but don't verify that users can accomplish their goals.

**Key Metrics:**
- Technical coverage: 7/10
- BDD compliance: 2/10
- User perspective: 3/10
- Ginkgo usage: 0/10 (not installed)

**Recommendation:** Either adopt Ginkgo as documented in AGENTS.md, or update the documentation to reflect the current standard Go testing approach. Regardless, add acceptance tests written from the user's perspective.

---

*Review completed. See action items above for next steps.*
