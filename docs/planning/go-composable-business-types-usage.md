# go-composable-business-types/id Integration Plan

**Date:** 2026-03-17  
**Library:** `github.com/larsartmann/go-composable-business-types/id`  
**Purpose:** Branded, strongly-typed identifiers for cmdguard applications  
**Status:** Planning Complete - Ready for Implementation

---

## Executive Summary

The `go-composable-business-types/id` package provides **phantom-type branded identifiers** that prevent mixing different entity IDs at compile time. This is a powerful pattern for CLI applications built with cmdguard that need to handle multiple entity types (users, orders, resources, etc.).

### Why Use This in cmdguard Applications?

```go
// Without branded IDs - stringly-typed, error-prone:
func GetUser(id string) error { ... }
func GetOrder(id string) error { ... }

GetOrder(userID)  // Compiles! Runtime bug.

// With branded IDs - compile-time safety:
type UserID = id.ID[UserBrand, string]
type OrderID = id.ID[OrderBrand, string]

func GetUser(id UserID) error { ... }
func GetOrder(id OrderID) error { ... }

GetOrder(userID)  // Compile error: type mismatch
```

---

## Integration Strategy

### 1. Dependency Addition

Add to `go.mod`:

```go
require (
    github.com/larsartmann/go-composable-business-types/id v0.1.0
)
```

### 2. Use Cases in cmdguard Applications

#### A. Command Handlers with Entity IDs

Commands that operate on specific entity types benefit from branded IDs:

```go
package main

import (
    "context"
    "fmt"

    "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
    "github.com/larsartmann/go-composable-business-types/id"
)

// Define brands (empty structs for compile-time distinctness)
type UserBrand struct{}
type OrganizationBrand struct{}

// Create branded ID types
type UserID = id.ID[UserBrand, string]
type OrganizationID = id.ID[OrganizationBrand, string]

// AppConfig uses branded IDs for configuration
type AppConfig struct {
    DefaultOrgID OrganizationID `flag:"default-org" help:"Default organization ID"`
    AdminUserID  UserID         `flag:"admin-user" help:"Admin user ID"`
}

// Command flags with branded IDs
type UserGetFlags struct {
    UserID UserID `flag:"user-id" short:"u" required:"true" help:"User ID to look up"`
}

func main() {
    ctx := context.Background()

    root, err := v2.New[AppConfig, v2.NoFlags]("myapp", "User management CLI", AppConfig{})
    if err != nil {
        panic(err)
    }

    // Add command with type-safe UserID flag
    err = root.AddCommand(v2.Command[AppConfig, UserGetFlags]{
        Use:   "get",
        Short: "Get user by ID",
        Flags: UserGetFlags{},
        RunE: func(ctx context.Context, cfg *AppConfig, flags UserGetFlags) error {
            // Type-safe: cannot accidentally pass wrong ID type
            user, err := fetchUser(flags.UserID)
            if err != nil {
                return err
            }
            fmt.Printf("User: %+v\n", user)
            return nil
        },
    })
    if err != nil {
        panic(err)
    }

    root.ExecuteAndExit(ctx)
}

// fetchUser only accepts UserID, cannot be called with OrganizationID
func fetchUser(id UserID) (*User, error) {
    // Implementation
    return &User{ID: id, Name: "John"}, nil
}

type User struct {
    ID   UserID
    Name string
}
```

#### B. DI Services with Typed IDs

Services registered in the DI scope can use branded IDs for type safety:

```go
package main

import (
    "context"

    "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
    "github.com/larsartmann/go-composable-business-types/id"
    "github.com/samber/do/v2"
)

type SessionBrand struct{}
type RequestBrand struct{}

type SessionID = id.ID[SessionBrand, string]
type RequestID = id.ID[RequestBrand, string]

// SessionService uses branded SessionID
type SessionService struct {
    currentSession SessionID
}

func NewSessionService(i do.Injector) (*SessionService, error) {
    return &SessionService{
        currentSession: id.NewID[SessionBrand]("session-123"),
    }, nil
}

func (s *SessionService) SessionID() SessionID {
    return s.currentSession
}

// RequestService uses branded RequestID
type RequestService struct {
    requestID RequestID
}

func NewRequestService(i do.Injector) (*RequestService, error) {
    return &RequestService{
        requestID: id.NewID[RequestBrand]("req-456"),
    }, nil
}

func main() {
    root, _ := v2.New[struct{}, v2.NoFlags]("app", "App", struct{}{})

    // Register services with typed IDs
    v2.Provide(root.ScopeStruct(), NewSessionService)
    v2.Provide(root.ScopeStruct(), NewRequestService)

    root.AddCommand(v2.Command[struct{}, v2.NoFlags]{
        Use:   "info",
        Short: "Show session info",
        RunE: func(ctx context.Context, cfg *struct{}, flags v2.NoFlags) error {
            sessionSvc, _ := v2.Invoke[*SessionService](root.ScopeStruct())
            requestSvc, _ := v2.Invoke[*RequestService](root.ScopeStruct())

            // Type-safe: these are different types
            sessionID := sessionSvc.SessionID()
            requestID := requestSvc.requestID

            // This would be a compile error:
            // useSession(requestID) // Error: cannot use RequestID as SessionID

            return useSession(sessionID)
        },
    })

    root.ExecuteAndExit(context.Background())
}

func useSession(id SessionID) error {
    // Only accepts SessionID
    return nil
}
```

#### C. Multi-Tenant CLI Applications

For CLIs that manage multiple tenants/organizations:

```go
package main

import (
    "context"
    "fmt"

    "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
    "github.com/larsartmann/go-composable-business-types/id"
)

type TenantBrand struct{}
type ProjectBrand struct{}
type EnvironmentBrand struct{}

type TenantID = id.ID[TenantBrand, string]
type ProjectID = id.ID[ProjectBrand, string]
type EnvironmentID = id.ID[EnvironmentBrand, string]

// Config with multiple ID types
type MultiTenantConfig struct {
    TenantID      TenantID      `flag:"tenant" short:"t" help:"Tenant ID"`
    ProjectID     ProjectID     `flag:"project" short:"p" help:"Project ID"`
    EnvironmentID EnvironmentID `flag:"env" short:"e" default:"prod" help:"Environment ID"`
}

// Commands that require specific ID combinations
type DeployFlags struct {
    ProjectID     ProjectID     `flag:"project" required:"true" help:"Project to deploy"`
    EnvironmentID EnvironmentID `flag:"env" default:"prod" help:"Target environment"`
}

type TenantListFlags struct {
    TenantID TenantID `flag:"tenant" help:"Filter by tenant"`
}

func main() {
    ctx := context.Background()

    root, _ := v2.New[MultiTenantConfig, v2.NoFlags]("deployer", "Deployment CLI", MultiTenantConfig{})

    // deploy command - requires project and environment
    root.AddCommand(v2.Command[MultiTenantConfig, DeployFlags]{
        Use:   "deploy",
        Short: "Deploy to environment",
        Flags: DeployFlags{},
        RunE: func(ctx context.Context, cfg *MultiTenantConfig, flags DeployFlags) error {
            // Type-safe: cannot mix up IDs
            return deploy(ctx, cfg.TenantID, flags.ProjectID, flags.EnvironmentID)
        },
    })

    // list command - filters by tenant
    root.AddCommand(v2.Command[MultiTenantConfig, TenantListFlags]{
        Use:   "list",
        Short: "List resources",
        Flags: TenantListFlags{},
        RunE: func(ctx context.Context, cfg *MultiTenantConfig, flags TenantListFlags) error {
            // Compile-time safety: these are distinct types
            // Cannot accidentally use ProjectID where TenantID is expected
            return listResources(ctx, flags.TenantID.Or(cfg.TenantID))
        },
    })

    root.ExecuteAndExit(ctx)
}

func deploy(ctx context.Context, tenant TenantID, project ProjectID, env EnvironmentID) error {
    fmt.Printf("Deploying project %s to %s for tenant %s\n", project, env, tenant)
    return nil
}

func listResources(ctx context.Context, tenant TenantID) error {
    fmt.Printf("Listing resources for tenant %s\n", tenant)
    return nil
}
```

#### D. With NanoId for Secure IDs

Use the `nanoid` package for cryptographically secure IDs:

```go
package main

import (
    "github.com/larsartmann/go-composable-business-types/id"
    "github.com/larsartmann/go-composable-business-types/nanoid"
)

type UserBrand struct{}

// Use NanoId as the underlying type for secure IDs
type UserID = id.ID[UserBrand, nanoid.NanoId]

func createUser() UserID {
    // Generate secure ID
    return id.NewID[UserBrand](nanoid.NewNanoId())
}

func parseUserID(s string) (UserID, error) {
    nano, err := nanoid.ParseNanoId(s)
    if err != nil {
        return UserID{}, err
    }
    return id.NewID[UserBrand](nano), nil
}
```

---

## Implementation Steps

### Step 1: Add Dependency

```bash
cd /Users/larsartmann/projects/cmdguard
go get github.com/larsartmann/go-composable-business-types/id
```

### Step 2: Create Example Application

Create `examples/branded-ids/` demonstrating:

- Basic ID usage with commands
- Config with branded IDs
- DI service integration
- NanoId integration

### Step 3: Add Tests

- Unit tests for ID serialization in config
- Integration tests for command handlers with IDs
- DI scope tests with typed services

### Step 4: Documentation

- Update README.md with ID usage examples
- Add to AGENTS.md integration patterns
- Document best practices

---

## Benefits of Integration

| Aspect            | Without Branded IDs | With Branded IDs       |
| ----------------- | ------------------- | ---------------------- |
| **Type Safety**   | Runtime errors      | Compile-time errors    |
| **Refactoring**   | Risky, needs tests  | Safe, compiler catches |
| **Documentation** | Comments needed     | Types document intent  |
| **IDE Support**   | Limited             | Autocomplete by type   |
| **Testing**       | Edge case heavy     | Focus on logic         |

---

## Compatibility Notes

### JSON Serialization

Branded IDs serialize seamlessly:

```go
type Config struct {
    UserID UserID `json:"user_id"`
}

// Serializes to: {"user_id": "user-123"}
// Zero value to: {"user_id": null}
```

### Flag Parsing

cmdguard's struct tags work with any type implementing `encoding.TextUnmarshaler`:

```go
type Flags struct {
    // id.ID implements TextUnmarshaler
    UserID UserID `flag:"user-id" required:"true"`
}
```

### DI Scope

Branded IDs work naturally with samber/do/v2:

```go
v2.Provide(scope, func(i do.Injector) (*UserService, error) {
    // Return service with typed ID
    return &UserService{ID: id.NewID[ServiceBrand]("svc-1")}, nil
})
```

---

## Migration Path

For existing cmdguard applications wanting to adopt branded IDs:

1. **Start with new entity types** - Don't refactor existing strings immediately
2. **Use type aliases** - `type UserID = id.ID[UserBrand, string]` keeps it simple
3. **Gradual adoption** - Convert functions one at a time
4. **Config first** - Add to config structs for new flags

---

## Recommendations

### DO

- ✅ Use for entity IDs (User, Order, Tenant, etc.)
- ✅ Combine with NanoId for secure identifiers
- ✅ Use in DI service constructors
- ✅ Apply to config values that are IDs
- ✅ Use type aliases for cleaner syntax

### DON'T

- ❌ Use for simple flags (enable, verbose, etc.)
- ❌ Create separate brands for similar concepts
- ❌ Use for internal IDs that never leave a function
- ❌ Over-engineer - start simple

---

## Conclusion

The `go-composable-business-types/id` library integrates naturally with cmdguard v2:

1. **Config fields** - Use branded IDs in config structs with flag tags
2. **Command flags** - Type-safe flags that prevent ID mix-ups
3. **DI services** - Services with typed IDs for compile-time safety
4. **JSON/YAML** - Automatic serialization support

This integration enables building CLI applications with the same level of type safety as the CLI framework itself.

**Next Steps:**

1. Add dependency to go.mod
2. Create example application
3. Add integration tests
4. Update documentation

## Resolution (2026-07-18)

**NOT IMPLEMENTED.** `go-composable-business-types` was never added to any module's `go.mod`, and no `examples/branded-ids/` exists. All code samples use the superseded `/v2` import path and the pre-v3 generic `Command[...]` API. Either execute the integration against `/v3` or archive this doc. Flagged for docs-health review (this file lives alongside dated historical plans but lacks a date stamp and reads as a living reference — its truthiness needs an explicit decision).
