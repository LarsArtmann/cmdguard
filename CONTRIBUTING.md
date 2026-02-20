# Contributing to cmdguard

Thank you for your interest in contributing to cmdguard! This document provides guidelines for contributing.

## Development Setup

### Prerequisites

- Go 1.24 or later
- [just](https://github.com/casey/just) for running common tasks
- [golangci-lint](https://golangci-lint.run/) for linting

### Clone and Build

```bash
git clone https://github.com/larsartmann/cmdguard.git
cd cmdguard
go build ./...
```

### Verify Everything Works

```bash
just verify
```

This runs:

- Build check
- All tests
- Linting

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Changes

- Follow Go best practices
- Add tests for new functionality
- Ensure errcheck compliance (check all error returns)
- Run `go fmt ./...` to format code

### 3. Test Your Changes

```bash
# Run all tests
just test

# Run with coverage
just test-cover

# Run linting
just lint

# Full verification
just verify
```

### 4. Commit

We follow conventional commits:

```bash
feat: add new feature
fix: fix a bug
docs: update documentation
refactor: code refactoring
test: add or update tests
chore: maintenance tasks
```

Example:

```bash
git commit -m "feat: add JSON logging option

- Add LoggerFormat option to Config
- Support "json" and "text" formats
- Default remains "text" for backward compatibility

Assisted-by: Your Name <your.email@example.com>"
```

### 5. Push and Create PR

```bash
git push origin feature/your-feature-name
```

Then create a pull request on GitHub.

## Code Standards

### Error Handling

Always check errors. Use `_ = ` to explicitly ignore when appropriate:

```go
// Good - explicit ignore
_ = os.Setenv("key", "value")

// Good - handle error
if err := doSomething(); err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
```

### Testing

- Use `testify/assert` and `testify/require`
- Table-driven tests preferred
- Test both success and error cases
- Use `t.Parallel()` for parallel tests when possible

Example:

```go
func TestNewGuardedCommand(t *testing.T) {
    tests := []struct {
        name    string
        args    []string
        wantErr bool
    }{
        {"valid", []string{"test"}, false},
        {"empty", []string{}, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := New(tt.args...)
            if tt.wantErr {
                require.Error(t, err)
                return
            }
            require.NoError(t, err)
        })
    }
}
```

### Documentation

- Add Go doc comments for all exported functions
- Keep README.md updated
- Update AGENTS.md for architecture changes

## Architecture Guidelines

### Guard Philosophy

cmdguard follows a "fail fast" philosophy:

- Panic at construction time for invalid commands
- No runtime validation errors
- Invalid states should be unrepresentable

When adding features:

1. Validate at construction time when possible
2. Panic with clear error messages
3. Never allow invalid states at runtime

### Adding New Features

1. **Discuss First**: For major features, open an issue to discuss
2. **Keep It Simple**: Prefer simple solutions over complex ones
3. **Test Coverage**: Maintain >80% coverage for new code
4. **Examples**: Add examples for significant features

## Release Process

1. Update version in `cmd/version.go`
2. Update CHANGELOG.md
3. Tag release: `git tag -a v0.x.x -m "Release v0.x.x"`
4. Push tag: `git push origin v0.x.x`

## Questions?

- Open an issue for questions
- Check existing documentation first
- Be respectful and constructive

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
