# ADR-003: Copy-on-Write Registry Pattern

**Status:** Accepted\
**Date:** 2026-07-11\
**Deciders:** Lars Artmann

## Context

cmdguard maintains two global registries: type handlers (`globalTypeRegistry`) and validators (`globalValidators`). These power the public `RegisterTypeHandler()` and `RegisterValidator()` APIs that allow users to extend flag parsing with custom types at the package level.

The problem: creating a new `CLI[T]` instance required deep-copying these registries, which was expensive (48% of `NewCLI` time). Every CLI instance got its own copy of all registered handlers and validators, even though most instances never customize them.

## Decision

### Copy-on-Write (COW) with lazy cloning

Instead of copying on every `NewCLI`, registries are **shared by default** and only **cloned on first write**:

1. **`share()`** — New `FlagRegistry` instances reference the global maps directly (zero copy)
2. **`register()`** — First write triggers a `sync.RWMutex`-guarded clone of the maps (guarded by a boolean `owned` flag); subsequent writes go to the instance-local copy
3. **Reads** always use the shared maps until the first write

### Public API implications

- `RegisterTypeHandler()` / `RegisterValidator()` — write to global defaults (visible to all instances that haven't cloned)
- `FlagRegistry.RegisterTypeHandler()` / `FlagRegistry.RegisterFlagValidator()` — trigger COW clone, write to instance-local maps only

### Global state justification

The `gochecknoglobals` linter flags these package-level variables. They are excluded in `.golangci.yml` because:

- Package-level registries are the COW pattern's foundation
- Injecting them would change the public API (no more `cmdguard.RegisterTypeHandler()`, instead `cli.Registry().RegisterTypeHandler()`)
- The COW pattern eliminates the performance concern (reads are shared, writes are local)

## Consequences

- `NewCLI` is ~48% faster (no deep copy)
- Global registration affects all new instances (expected behavior for a registry)
- Instance-level registration is isolated (COW guarantee)
- Three `gochecknoglobals` exclusions are permanent (documented in ADR-002)
