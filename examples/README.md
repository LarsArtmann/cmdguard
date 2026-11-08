# cmdguard Examples

Three examples: the flagship **taskctl** plus two focused demos.

## taskctl — Task Manager CLI

A complete CLI application that shows all cmdguard capabilities in a realistic context:

```bash
# From the project root:
go run examples/taskctl/main.go --help
go run examples/taskctl/main.go list
go run examples/taskctl/main.go add --title "Fix bug" --priority high
go run examples/taskctl/main.go done --id 1
go run examples/taskctl/main.go stats --format json
```

See [examples/taskctl/README.md](taskctl/README.md) for the full feature matrix and usage guide.

## config-explain — Config Provenance Demo

Shows the `config explain` command (`WithConfigExplainCommand()`): resolves every config field to its source (flag → env → config file → default) with `secret:"true"` redaction.

```bash
go run examples/config-explain/main.go config explain
go run examples/config-explain/main.go config explain --output json
```

## docs-generator — GenerateDocs Demo

Minimal CLI exercising `cli.GenerateDocs(w)` — markdown documentation for the full command tree.

```bash
go run examples/docs-generator/main.go
```

## Running Tests

```bash
go test ./examples/taskctl/... -count=1 -race
```
