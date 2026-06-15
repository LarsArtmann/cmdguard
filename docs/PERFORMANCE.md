# Performance

> Benchmarks for cmdguard v2 on AMD Ryzen AI MAX+ 395 (Linux/amd64).
> **Last Updated:** 2026-06-14

---

## TL;DR

cmdguard adds **<6 µs** overhead for CLI creation and **~12 ns** for command validation. For typical CLIs with a few commands, the total cold-start overhead is well under 1 ms — less than 6% of Go runtime initialization.

Copy-on-write registries (v2.7.0+) reduce per-command allocations by **48%** and memory by **22%** compared to eager cloning.

---

## Optimizations Applied (v2.7.0)

| Optimization                          | Effect                                                  | Files                                              |
| ------------------------------------- | ------------------------------------------------------- | -------------------------------------------------- |
| Copy-on-write registries              | **-10 allocs, -1.9 KB per command**, -48% faster NewCLI | `type_handler.go`, `flags_validate.go`, `flags.go` |
| Cached `os.UserHomeDir()`             | Eliminates redundant syscalls for `~/` path expansion   | `config_file.go`                                   |
| Iterator-based traversal (`iter.Seq`) | Zero-allocation alternative to defensive copies         | `flags.go`, `flags_suggest.go`, `flow_context.go`  |
| Regex cache safety documentation      | Documents bounded usage of `sync.Map` regex cache       | `flags_validate.go`                                |

---

## Benchmarks

### CLI Lifecycle

| Operation                 | Time    | Allocations | Memory  |
| ------------------------- | ------- | ----------- | ------- |
| `NewCLI` (basic)          | ~5.8 µs | 81          | ~6.6 KB |
| `NewCLI` (with long desc) | ~6.1 µs | 81          | ~6.6 KB |
| `AddCommand`              | ~7.2 µs | 93          | ~8.4 KB |
| `Execute` (help)          | ~608 µs | 9,175       | ~423 KB |
| `NewCommand`              | ~100 ns | 1           | ~288 B  |
| `Command.Validate`        | ~12 ns  | 0           | 0 B     |

_Note: `Execute` with help is slower because fang renders styled output. Actual command execution is significantly faster._

### Flag Parsing

| Operation                    | Time    | Allocations | Memory  |
| ---------------------------- | ------- | ----------- | ------- |
| `ParseFlagTags` (4 fields)   | ~1.8 µs | 9           | ~1.5 KB |
| `NewFlagRegistry` (2 fields) | ~0.8 µs | 11          | ~912 B  |
| `ParseDuration`              | ~105 ns | 0           | 0 B     |
| `ParseLogLevel`              | ~83 ns  | 0           | 0 B     |
| `ParseURL`                   | ~657 ns | 6           | ~768 B  |
| `ParseEmail`                 | ~1.0 µs | 25          | ~504 B  |
| `ParsePort` (numeric)        | ~63 ns  | 0           | 0 B     |
| `ParsePort` (named)          | ~35 ns  | 0           | 0 B     |
| `ParseFilePath`              | ~1.2 µs | 7           | ~586 B  |
| `ParseHostPort`              | ~151 ns | 0           | 0 B     |

### Copy-on-Write Registry

| Operation                         | Time    | Allocations | Memory  | Notes                     |
| --------------------------------- | ------- | ----------- | ------- | ------------------------- |
| `NewFlagRegistry` (COW, no write) | ~0.8 µs | 11          | ~912 B  | Shares global maps        |
| `NewFlagRegistry` + 1 write       | ~0.9 µs | 19          | ~2.2 KB | Triggers lazy clone       |
| `TagsSeq()` (iterator)            | ~12 ns  | 0           | 0 B     | Zero-allocation traversal |
| `Tags()` (defensive copy)         | ~101 ns | 1           | ~384 B  | Legacy API                |

### Dependency Injection

| Operation  | Time    | Allocations | Memory  |
| ---------- | ------- | ----------- | ------- |
| `NewScope` | ~597 ns | 16          | ~809 B  |
| `Provide`  | ~2.1 µs | 30          | ~1.7 KB |
| `Invoke`   | ~186 ns | 5           | ~160 B  |

---

## What This Means

### Startup Overhead

For a typical CLI with 5 commands:

- 1× `NewCLI`: ~6 µs
- 5× `AddCommand`: ~36 µs
- Total startup: **<42 µs**

This is negligible compared to Go runtime initialization (~1–5 ms).

### Per-Command Overhead

For a command with typed flags:

- Flag tag parsing: ~1.8 µs (once at registration)
- Flag registry creation: ~0.8 µs (once at registration, COW)
- Command validation: ~12 ns (once at registration)

At execution time, cmdguard adds essentially zero overhead beyond what Cobra already does — flags are parsed by pflag just like in raw Cobra.

### DI Overhead

- Service registration (`Provide`): ~2.1 µs
- Service lookup (`Invoke`): ~186 ns

For a CLI that resolves 10 services per command:

- 10 × 186 ns = **1.9 µs** of DI overhead per execution

This is negligible for any CLI that does I/O.

---

## Copy-on-Write Registry Design

Starting in v2.7.0, `FlagRegistry` uses copy-on-write for its `typeRegistry` and
`validatorRegistry`. Instead of eagerly cloning the global registries at creation time
(~12 allocs, ~1.5 KB per command), the registries share the global maps and only clone
on the first write (`RegisterTypeHandler` / `RegisterFlagValidator`).

**Benefits:**

- 48% faster `NewCLI` (5.8 µs vs 11.0 µs)
- 10 fewer allocations per command
- 22% less memory per command
- Zero cost for commands that never customize their registries (the common case)

**Behavioral note:** With COW, global registrations via `RegisterTypeHandler()` /
`RegisterValidator()` are visible to `FlagRegistry` instances created before the
registration, as long as those instances haven't triggered a lazy clone. This is an
improvement over the previous behavior (where instances were snapshot at creation).

---

## Reproducing

Run benchmarks locally:

```bash
go test ./benchmarks/ -bench=. -benchmem -count=5
```

For stable results, close other applications and run on a quiet machine.
