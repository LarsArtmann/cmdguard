# Justfile for cmdguard
# https://github.com/casey/just

# Default recipe - show available commands
default:
    @just --list

# Build all packages
build:
    go build ./...

# Run all tests
test:
    go test ./...

# Run tests with verbose output
test-v:
    go test -v ./...

# Run tests with coverage
test-cover:
    go test -cover ./...

# Run tests with race detector
test-race:
    go test -race ./...

# Run linting (requires golangci-lint)
lint:
    golangci-lint run --enable=errcheck ./...

# Format code
fmt:
    go fmt ./...

# Tidy dependencies
tidy:
    go mod tidy

# Verify build and tests
verify: build test lint
    @echo "✅ All checks passed"

# Run basic example
run-basic:
    go run ./examples/basic/main.go hello

# Run advanced example
run-advanced:
    go run ./examples/advanced/main.go db migrate

# Run guarded example
run-guarded:
    go run ./examples/guarded/main.go validate

# Clean build artifacts
clean:
    go clean
    rm -f cmdguard

# Show dependency graph
deps:
    go mod graph

# List direct dependencies
deps-list:
    go list -m all | grep -v "^github.com/larsartmann"

# Update dependencies
update:
    go get -u ./...
    go mod tidy
