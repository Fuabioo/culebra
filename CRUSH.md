# Culebra Development Commands

## Build & Test Commands
```bash
# Run tests
just test

# Run tests with coverage
just test-coverage

# Lint code
just lint

# Format code  
just fmt

# Run all checks (format, lint, test)
just check
```

## Example Commands
```bash
# Run the basic example
just example

# Build example binary
just build
```

## Development Workflow
```bash
# Full development workflow
just dev
```

## Dependencies
```bash
# Install/update dependencies
just deps

# Tidy go modules
just tidy
```

## Project Structure
- `luaconf.go` - Core Lua config loading
- `viper.go` - Viper integration
- `cobra.go` - Cobra CLI integration  
- `internal/mapper.go` - Lua ↔ Go type conversion
- `examples/basic/` - Example CLI app
- `lua/stdlib.lua` - Optional Lua helpers

## Key Features
- Load `.lua` files as configuration
- Optional Viper integration with `BindToViper()`
- One-liner Cobra integration with `UseWithCobra()`
- Zero dependencies on Cobra/Viper for core functionality