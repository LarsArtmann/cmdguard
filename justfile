# Justfile for cmdguard
# https://github.com/casey/just

# Default recipe - show available commands
default:
    @just --list

# Build all packages
build:
    go build ./...

# Build with specific version (via ldflags)
build-version version:
    go build -ldflags "-X github.com/larsartmann/cmdguard/pkg/cmdguard.version={{version}}" ./...

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

# Run advanced-flags example
run-advanced:
    go run ./examples/advanced-flags/main.go env

# Run validation example
run-validation:
    go run ./examples/validation/main.go greet --name=World

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

# Dogfooding - Self-validate against HOW_TO_GOLANG policy
dogfood:
    @echo "🐕 Dogfooding cmdguard..."
    @echo ""
    @echo "=== Checking file size limits (250 lines max) ==="
    @find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*" -exec wc -l {} + | sort -rn | head -20 | awk '{if ($1 > 250) print "❌ OVER LIMIT: " $0; else print "✅ OK: " $0}'
    @echo ""
    @echo "=== Checking function size limits (30 lines max) ==="
    @find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*" -exec awk '/^func \(/,/^}/ {lines++} /^}/ {if (lines > 30) print "❌ OVER LIMIT: " FILENAME " (" lines " lines)"; lines=0}' {} \; | head -20
    @echo ""
    @echo "=== Checking for banned libraries ==="
    @grep -r "stretchr/testify" --include="*.go" . 2>/dev/null && echo "❌ FOUND: testify (banned)" || echo "✅ No testify found"
    @grep -r "pkg/errors" --include="*.go" . 2>/dev/null && echo "❌ FOUND: pkg/errors (banned)" || echo "✅ No pkg/errors found"
    @echo ""
    @echo "=== Running build ==="
    @just build
    @echo ""
    @echo "=== Running tests ==="
    @just test
    @echo ""
    @echo "🐕 Dogfooding complete!"
