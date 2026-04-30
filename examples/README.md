# cmdguard Examples

Working examples demonstrating cmdguard v2 features.

## Quick Start

```bash
# From the project root:
go run examples/basic/main.go hello
```

## Examples

| Example            | Feature                                         | Run                                                               |
| ------------------ | ----------------------------------------------- | ----------------------------------------------------------------- |
| **basic**          | Simple CLI with subcommands                     | `go run examples/basic/main.go hello`                             |
| **typed**          | DI and lifecycle hooks                          | `go run examples/typed/main.go`                                   |
| **di**             | Dependency injection patterns                   | `go run examples/di/main.go check`                                |
| **di-patterns**    | Service registration, invocation, health checks | `go run examples/di-patterns/main.go list`                        |
| **advanced-flags** | Custom flag types (Enum, Duration, LogLevel)    | `go run examples/advanced-flags/main.go server --port 8080`       |
| **validation**     | PreRunE validation, required flags              | `go run examples/validation/main.go greet --name=Alice`           |
| **env-tags**       | Environment variable binding, WithEnvPrefix     | `DB_HOST=db.example.com go run examples/env-tags/main.go show`    |
| **output**         | Rich output formats (table/json/csv/yaml)       | `go run examples/output/main.go users --format json`              |
| **counting**       | Counting flags (-v/-vv/-vvv)                    | `go run examples/counting/main.go greet -vvv`                     |
| **error-handling** | Sentinel errors, FlagError, suggestions         | `go run examples/error-handling/main.go fetch --url=invalid`      |
| **signals**        | Signal handling, graceful shutdown              | `go run examples/signals/main.go serve`                           |
| **subcommands**    | NewParentCommand, command groups                | `go run examples/subcommands/main.go migrate up --env=production` |

## Feature Matrix

| Feature             | basic | typed | di  | di-patterns | advanced-flags | validation | env-tags | output | counting | error-handling | signals | subcommands |
| ------------------- | ----- | ----- | --- | ----------- | -------------- | ---------- | -------- | ------ | -------- | -------------- | ------- | ----------- |
| NewCLI              | x     | x     | x   | x           | x              | x          | x        | x      | x        | x              | x       | x           |
| NewCommand          | x     | x     | x   | x           | x              | x          | x        | x      | x        | x              | x       | x           |
| WithFlags           |       | x     |     | x           | x              | x          | x        | x      | x        | x              | x       | x           |
| PreRunE             |       |       |     |             | x              | x          |          |        |          | x              |         |             |
| PostRunE            |       |       |     |             |                |            |          |        |          |                |         |             |
| DI (Provide/Invoke) |       | x     | x   | x           |                |            |          |        |          |                |         |             |
| env tags            |       |       |     |             |                |            | x        |        |          |                |         |             |
| WithEnvPrefix       |       |       |     |             |                |            | x        |        |          |                |         |             |
| count:"true"        |       |       |     |             |                |            |          |        | x        |                |         |             |
| OutputTable         |       |       |     |             |                |            |          | x      |          |                |         |             |
| Signal handling     |       |       |     |             |                |            |          |        |          |                | x       |             |
| Error wrapping      |       |       |     |             |                |            |          |        |          | x              |         |             |
| NewParentCommand    |       |       |     |             |                |            |          |        |          |                |         | x           |

## Adding a New Example

1. Create `examples/<name>/main.go`
2. Add a package comment with usage instructions
3. Import `v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"`
4. Use `examplesinternal.Execute(ctx, cli)` for consistent error handling
5. Update this README
