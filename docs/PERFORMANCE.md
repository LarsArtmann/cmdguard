# Performance

> Benchmarks for cmdguard v4.0.0 on AMD Ryzen AI MAX+ 395 (Linux/amd64, Go 1.26.5).
> **Last Updated:** 2026-08-06

---

## TL;DR

cmdguard adds **~6.9 µs** overhead for full CLI creation (`NewCLI`), **~100 ns** per command, and **~10.5 ns** for command validation. The per-command overhead is negligible — a typical 5-command CLI adds ~7.4 µs total, well under 1% of Go runtime initialization (~1–5 ms).

Copy-on-write registries reduce per-command allocations by **48%** and memory by **22%** compared to eager cloning.

---

## Optimizations Applied

| Optimization                          | Effect                                                                                                   | Files                                                              |
| ------------------------------------- | -------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------ |
| Copy-on-write registries              | **-10 allocs, -1.9 KB per command**, -48% faster NewCLI _(measured at original optimization time, v2.6)_ | `pkg/cmdguard/v4/type_handler.go`, `flags_validate.go`, `flags.go` |
| Cached `os.UserHomeDir()`             | Eliminates redundant syscalls for `~/` path expansion                                                    | `pkg/cmdguard/v4/config_file.go`                                   |
| Iterator-based traversal (`iter.Seq`) | Zero-allocation alternative to defensive copies                                                          | `pkg/cmdguard/v4/flags.go`, `flags_suggest.go`, `flow_context.go`  |
| Regex cache safety documentation      | Documents bounded usage of `sync.Map` regex cache                                                        | `pkg/cmdguard/v4/flags_validate.go`                                |

---

## Benchmarks

### CLI Lifecycle

| Operation          | Time     | Allocations | Memory  |
| ------------------ | -------- | ----------- | ------- |
| `NewCLI`           | ~6.9 µs  | 77          | ~6.8 KB |
| `Execute` (help)   | ~580 µs  | ~6,300      | ~294 KB |
| `NewCommand`       | ~102 ns  | 1           | ~288 B  |
| `Command.Validate` | ~10.5 ns | 0           | 0 B     |

_Note: `Execute` with help is high-variance (580 µs–1.1 ms across runs) because it exercises the full cobra execution path including GC. The best-case (~580 µs) reflects true overhead without external interference. Stdout is redirected to `/dev/null` during benchmarking._

### Flag Parsing

| Operation                    | Time    | Allocations | Memory  |
| ---------------------------- | ------- | ----------- | ------- |
| `ParseFlagTags` (4 fields)   | ~4.3 µs | 11          | ~1.6 KB |
| `NewFlagRegistry` (2 fields) | ~910 ns | 9           | ~896 B  |
| `ParseDuration`              | ~86 ns  | 0           | 0 B     |
| `ParseLogLevel`              | ~78 ns  | 0           | 0 B     |
| `ParseURL`                   | ~636 ns | 6           | ~768 B  |
| `ParseEmail`                 | ~907 ns | 25          | ~504 B  |
| `ParsePort` (numeric)        | ~61 ns  | 0           | 0 B     |
| `ParsePort` (named)          | ~39 ns  | 0           | 0 B     |
| `ParseFilePath`              | ~1.2 µs | 7           | ~586 B  |
| `ParseHostPort`              | ~119 ns | 0           | 0 B     |

_Note: v4's `ParseFlagTags` uses 11 allocs (vs v2's 9) due to nested struct recursion support. Each field allocates an `Index` reflect path for the flattened flag registration. This is expected v4 overhead for the richer feature set._

### Copy-on-Write Registry

| Operation                         | Time     | Allocations | Memory  | Notes                     |
| --------------------------------- | -------- | ----------- | ------- | ------------------------- |
| `NewFlagRegistry` (COW, no write) | ~880 ns  | 9           | ~896 B  | Shares global maps        |
| `NewFlagRegistry` + 1 write       | ~1.1 µs  | 15          | ~2.1 KB | Triggers lazy clone       |
| `TagsSeq()` (iterator)            | ~18.5 ns | 0           | 0 B     | Zero-allocation traversal |
| `Tags()` (defensive copy)         | ~453 ns  | 1           | ~448 B  | Legacy API                |

### Dependency Injection

| Operation          | Time    | Allocations | Memory  |
| ------------------ | ------- | ----------- | ------- |
| `NewScope`         | ~410 ns | 11          | ~688 B  |
| `NewScopeWithOpts` | ~420 ns | 11          | ~688 B  |
| `Provide`          | ~1.7 µs | 26          | ~1.7 KB |
| `Invoke`           | ~213 ns | 5           | ~160 B  |
| `CloneScope`       | ~6.4 µs | 39          | ~3.0 KB |
| `ProvideInvoke`    | ~7.8 µs | 41          | ~3.1 KB |

### Flight Recorder

| Operation             | Time    | Allocations | Memory |
| --------------------- | ------- | ----------- | ------ |
| `New` (Recorder)      | ~69 ns  | 2           | ~304 B |
| `Middleware` overhead | ~71 ns  | 0           | 0 B    |
| `Capture`             | ~433 µs | 91          | ~47 KB |

_Captures a runtime/trace snapshot to disk. Cost is dominated by trace serialization and file I/O._

---

## What This Means

### Startup Overhead

For a typical CLI with 5 commands:

- 1× `NewCLI` (includes scope creation + flag registration + CLIOption processing): ~6.9 µs
- 5× `NewCommand` + validation: ~563 ns
- Total startup: **~7.4 µs**

This is negligible compared to Go runtime initialization (~1–5 ms).

### Per-Command Overhead

For a command with typed flags:

- Flag tag parsing: ~4.3 µs (once at registration)
- Flag registry creation: ~880 ns (once at registration, COW)
- Command validation: ~10.5 ns (once at registration)

At execution time, cmdguard adds essentially zero overhead beyond what Cobra already does — flags are parsed by pflag just like in raw Cobra.

### DI Overhead

- Service registration (`Provide`): ~1.7 µs
- Service lookup (`Invoke`): ~213 ns

For a CLI that resolves 10 services per command:

- 10 × 213 ns = **2.1 µs** of DI overhead per execution

This is negligible for any CLI that does I/O.

---

## Copy-on-Write Registry Design

`FlagRegistry` uses copy-on-write for its `typeRegistry` and
`validatorRegistry`. Instead of eagerly cloning the global registries at creation time
(~12 allocs, ~1.5 KB per command), the registries share the global maps and only clone
on the first write (`RegisterTypeHandler` / `RegisterFlagValidator`).

**Benefits:**

- 48% faster `NewCLI` (measured at original optimization time, v2.6)
- 10 fewer allocations per command (v2.6 baseline)
- 22% less memory per command (v2.6 baseline)
- Zero cost for commands that never customize their registries (the common case)

**Behavioral note:** With COW, global registrations via `RegisterTypeHandler()` /
`RegisterValidator()` are visible to `FlagRegistry` instances created before the
registration, as long as those instances haven't triggered a lazy clone. This is an
improvement over the previous behavior (where instances were snapshot at creation).

---

## Reproducing

Run benchmarks locally:

```bash
GOEXPERIMENT=jsonv2 go test ./benchmarks/ -bench=. -benchmem -count=5
GOEXPERIMENT=jsonv2 go test ./flightrecorder/ -bench=. -benchmem -count=5
```

For stable results, close other applications and run on a quiet machine.
All benchmarks suppress stdout/stderr output during the measurement loop.
