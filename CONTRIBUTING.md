# Contributing to cmdguard

Thank you for your interest in contributing to cmdguard! This document provides guidelines for contributing.

## Development Setup

### With Nix (Recommended)

```bash
# Enter the dev shell (Go 1.26, gopls, golangci-lint)
nix develop

# Or with direnv for automatic activation
echo "use flake" > .envrc && direnv allow
```

### Without Nix

- Go 1.26
- [golangci-lint](https://golangci-lint.run/) v2.x

### Verify Everything Works

```bash
go build ./...
go test ./... -count=1 -timeout 120s -race
golangci-lint run ./...
nix flake check   # format check (requires Nix)
```

## Development Workflow

### 1. Create a Branch

```bash
git switch -c feature/your-feature-name
```

### 2. Make Changes

- Follow Go best practices
- Add tests for new functionality
- Ensure `golangci-lint run ./...` passes
- Run `nix fmt` to format Nix and Go files

### 3. Test Your Changes

```bash
# Run all tests with race detection
go test ./... -count=1 -timeout 120s -race

# Run with coverage
go test ./... -count=1 -timeout 120s -cover

# Run linting
golangci-lint run ./...

# Format check
nix fmt
```

### 4. Commit

We follow conventional commits:

```
feat: add new feature
fix: fix a bug
docs: update documentation
refactor: code refactoring
test: add or update tests
chore: maintenance tasks
```

### 5. Push and Create PR

```bash
git push -u origin feature/your-feature-name
```

Then create a pull request on GitHub.

## Code Standards

### Error Handling

Always check errors. Wrap with context:

```go
if err := doSomething(); err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
```

### Testing

- Table-driven tests preferred
- Test both success and error cases
- Use `t.Parallel()` in every test function and subtest
- Use `//nolint:paralleltest` for tests using `t.Setenv`
- Single internal test package (`v3`, accesses private helpers)
- Use `ExecuteWithArgs(ctx, args)` for integration tests that exercise the full CLI pipeline
- Fuzz targets use inline `f.Add()` for seed corpus; file-based seeds go in `testdata/fuzz/`
- `NoFlags` is a distinct named type — use `(NoFlags{})` with parens for comparisons
- Use `//nolint:fatcontext` at file level for test files with context in closures

### Documentation

- Add Go doc comments for all exported functions
- Keep README.md updated
- Update AGENTS.md for architecture changes

## Architecture Guidelines

### v3 Design Principles

- No panics in library code — all operations return errors
- Constructor pattern via `NewCommand`/`NewParentCommand`
- Functional options for configuration
- Composition over inheritance

When adding features:

1. Validate at construction time, return errors
2. Keep impossible states unrepresentable with strong types
3. Maintain >80% test coverage for new code
4. Add examples for significant features

## Questions?

- Open an issue for questions
- Check existing documentation first
- Be respectful and constructive

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
