# Culebra - Lua configuration loader for Go
# Available commands

# Default recipe - show available commands
default:
    @just --list

# Run all examples
example:
    @echo "=== Basic Example ==="
    cd examples/basic && go run main.go --config config.lua
    @echo "=== Basic Example (Neovim style) ==="
    cd examples/basic && go run main.go --config config-neovim-style.lua
    @echo "=== Autoload Example ==="
    cd examples/autoload && go run main.go --config example.lua
    @echo "=== Arrays Example ==="
    cd examples/arrays && go run main.go
    @echo "=== Advanced Example ==="
    cd examples/advanced && go run main.go --config config.lua

# Run tests
test:
    go test -v ./...

# Run tests with coverage
test-coverage:
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    @echo "Coverage report generated: coverage.html"

# Lint the code using golangci-lint
lint:
    golangci-lint run

# Format code
fmt:
    go fmt ./...

# Tidy dependencies
tidy:
    go mod tidy

# Build all example binaries
build:
    cd examples/basic && go build -o basic-example main.go
    cd examples/autoload && go build -o autoload-example main.go 
    cd examples/arrays && go build -o arrays-example main.go
    cd examples/advanced && go build -o advanced-example main.go

# Clean build artifacts
clean:
    rm -f examples/basic/basic-example
    rm -f examples/autoload/autoload-example
    rm -f examples/arrays/arrays-example
    rm -f examples/advanced/advanced-example
    rm -f coverage.out coverage.html

# Install dependencies
deps:
    go mod download

# Run all checks (format, lint, test)
check: fmt lint test

# Development workflow - format, lint, test, and run example
dev: fmt lint test example
