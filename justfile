# Culebra - Lua configuration loader for Go
# Available commands

# Default recipe - show available commands
default:
    @just --list

# Run the basic example
example:
    cd examples/basic && go run main.go --config config.lua
    cd examples/basic && go run main.go --config config-neovim-style.lua
    cd examples/autoload && go run main.go

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

# Build the example binary
build:
    cd examples/basic && go build -o example main.go

# Clean build artifacts
clean:
    rm -f examples/basic/example
    rm -f coverage.out coverage.html

# Install dependencies
deps:
    go mod download

# Run all checks (format, lint, test)
check: fmt lint test

# Development workflow - format, lint, test, and run example
dev: fmt lint test example
