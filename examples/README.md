# cmdguard Examples

A single, production-grade example demonstrating every cmdguard feature.

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

## Running Tests

```bash
go test ./examples/taskctl/... -count=1 -race
```
