# Performance

> Benchmarks for cmdguard v2 on AMD Ryzen AI MAX+ 395 (Linux/amd64).
> **Last Updated:** 2026-05-27

---

## TL;DR

cmdguard adds **<10 µs** overhead for CLI creation and **<10 ns** for command validation. For typical CLIs with a few commands, the total cold-start overhead is well under 1 ms.

---

## Benchmarks

### CLI Lifecycle

| Operation | Time | Allocations | Memory |
| --- | --- | --- | --- |
| `NewCLI` (basic) | ~7.4 µs | 86 | ~8.0 KB |
| `NewCLI` (with long desc) | ~8.1 µs | 86 | ~8.0 KB |
| `AddCommand` | ~8.5 µs | 95 | ~9.5 KB |
| `Execute` (help) | ~919 µs | 9,175 | ~431 KB |
| `NewCommand` | ~76 ns | 1 | ~240 B |
| `Command.Validate` | ~8 ns | 0 | 0 B |

*Note: `Execute` with help is slower because fang renders styled output. Actual command execution is significantly faster.*

### Flag Parsing

| Operation | Time | Allocations | Memory |
| --- | --- | --- | --- |
| `ParseFlagTags` (4 fields) | ~1.5 µs | 9 | ~1.4 KB |
| `NewFlagRegistry` (2 fields) | ~1.8 µs | 21 | ~2.7 KB |
| `ParseDuration` | ~79 ns | 0 | 0 B |
| `ParseLogLevel` | ~67 ns | 0 | 0 B |
| `ParseURL` | ~523 ns | 6 | ~768 B |
| `ParseEmail` | ~942 ns | 25 | ~504 B |
| `ParsePort` (numeric) | ~77 ns | 0 | 0 B |
| `ParsePort` (named) | ~47 ns | 0 | 0 B |
| `ParseFilePath` | ~1.9 µs | 7 | ~586 B |
| `ParseHostPort` | ~132 ns | 0 | 0 B |

### Dependency Injection

| Operation | Time | Allocations | Memory |
| --- | --- | --- | --- |
| `NewScope` | ~543 ns | 16 | ~809 B |
| `Provide` | ~1.6 µs | 28 | ~1.7 KB |
| `Invoke` | ~179 ns | 5 | ~160 B |

---

## What This Means

### Startup Overhead

For a typical CLI with 5 commands:

- 1× `NewCLI`: ~7 µs
- 5× `AddCommand`: ~42 µs
- Total startup: **<50 µs**

This is negligible compared to Go runtime initialization (~1–5 ms).

### Per-Command Overhead

For a command with typed flags:

- Flag tag parsing: ~1.5 µs (once at registration)
- Flag registry creation: ~1.8 µs (once at registration)
- Command validation: ~8 ns (once at registration)

At execution time, cmdguard adds essentially zero overhead beyond what Cobra already does — flags are parsed by pflag just like in raw Cobra.

### DI Overhead

- Service registration (`Provide`): ~1.6 µs
- Service lookup (`Invoke`): ~179 ns

For a CLI that resolves 10 services per command:

- 10 × 179 ns = **1.8 µs** of DI overhead per execution

This is negligible for any CLI that does I/O.

---

## Reproducing

Run benchmarks locally:

```bash
go test ./benchmarks/ -bench=. -benchmem -count=5
```

For stable results, close other applications and run on a quiet machine.
