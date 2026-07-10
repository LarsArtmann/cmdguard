# Domain Language

A **Unified Language** for `cmdguard` — shared across contributors and AI assistants.

Inspired by Domain-Driven Design (DDD) Ubiquitous Language.

## Glossary

| Term          | Definition                                                          | Context                              |
| ------------- | ------------------------------------------------------------------- | ------------------------------------ |
| CLI           | A `CLI[T]` instance — the root application container                | Entry point for all cmdguard apps    |
| Command       | A `Command[T,F]` — a typed subcommand with config T and flags F     | Building block for command trees     |
| Config        | The root type parameter `T` on `CLI[T]` — typed application config  | Drives flag registration             |
| Flags         | The per-command type parameter `F` on `Command[T,F]`                | Typed flag struct with struct tags   |
| FlagRegistry  | Holds parsed flag tags and type/validator registries                | Copy-on-write, per-CLI instance      |
| FlagTag       | Parsed metadata from struct tags (`flag`, `default`, `help`, etc.)  | Drives flag registration and parsing |
| Scope         | DI container wrapping `samber/do/v2` injector                       | Manages service lifecycle            |
| Plugin        | Bundles custom type handlers + validators for one-step registration | Applied globally or per-instance     |
| TypeHandler   | Interface for registering, parsing, and defaulting a flag type      | Extensible via `RegisterTypeHandler` |
| FlagValidator | Named validation function applied to flag values                    | Registered via `RegisterValidator`   |
| FlowContext   | Tracks command execution path and branch state                      | Used by middleware and hooks         |
| SilenceUsage  | Suppresses cobra's usage-on-error footgun                           | True by default                      |

## Entities

Objects with identity and lifecycle.

| Term         | Definition                                                | Context                               |
| ------------ | --------------------------------------------------------- | ------------------------------------- |
| CLI[T]       | Root application — owns scope, root command, registry     | Created once per process              |
| Command[T,F] | Typed subcommand with RunE handler                        | Registered via `AddCommand`           |
| Scope        | DI container — manages service registration and lifecycle | Created by CLI or provided externally |

## Value Objects

Immutable objects defined by attributes.

| Term         | Definition                                         | Context                                 |
| ------------ | -------------------------------------------------- | --------------------------------------- |
| FlagTag      | Parsed struct tag metadata for a single flag field | Built by `ParseFlagTags`                |
| OutputFormat | Type-safe enum for output formats (table, json, …) | Aliased from go-output                  |
| Duration     | Validated duration string (e.g. "5s", "1h30m")     | Custom type with Parse/Default handlers |
| Port         | Validated network port (1-65535)                   | Custom type with Parse/Default handlers |
| Email        | RFC 5322 validated email address                   | Custom type with Parse/Default handlers |
| URL          | Validated URL string                               | Custom type with Parse/Default handlers |
| Enum         | Enumerated string value from a fixed set           | Custom type with Parse/Default handlers |

## Commands

Actions the system can perform.

| Term            | Definition                                       | Context                             |
| --------------- | ------------------------------------------------ | ----------------------------------- |
| NewCLI          | Creates a CLI application with typed config      | Entry point                         |
| AddCommand      | Registers a typed subcommand on the CLI          | Supports per-command flag types     |
| Execute         | Runs the CLI with os.Args                        | Returns error for exit-code mapping |
| ExecuteWithArgs | Runs the CLI with explicit args (testing)        | Used in tests                       |
| ExecuteAndExit  | Runs CLI and calls os.Exit with mapped exit code | Blessed entry point                 |

---

> **How to use this file:**
>
> - Keep terms concise — one clear sentence per definition
> - Update when new domain concepts emerge
> - Use these terms consistently in code, docs, and conversations
> - When in doubt about a word's meaning, check here first
