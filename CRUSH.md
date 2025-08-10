# Culebra Development Commands

## Build & Test Commands
```bash
# Run tests
just test

# Run tests with coverage and view html report
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
# Run all examples (basic, autoload, arrays, advanced, viper-showcase)
just example

# Build all example binaries
just build

# Clean build artifacts
just clean
```

## Development Workflow
```bash
# Full development workflow (format, lint, test, examples)
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
- `cobra.go` - Cobra CLI integration with array conversion
- `internal/mapper.go` - Lua ↔ Go type conversion with array support
- `examples/basic/` - Basic usage example
- `examples/autoload/` - Automatic config detection
- `examples/arrays/` - Array handling demonstration
- `examples/advanced/` - Complex Lua config with environment logic
- `examples/viper-showcase/` - Comprehensive Viper Get functions demo
- `lua/stdlib.lua` - Optional Lua helpers

## Key Features
- Load `.lua` files as configuration with full Lua programmability
- Optional Viper integration with `BindToViper()` and comprehensive Get function support
- One-liner Cobra integration with `UseWithCobra()` and automatic array conversion
- Zero dependencies on Cobra/Viper for core functionality
- Lua array to Go slice conversion with `ConvertArrays` flag
- Environment-specific configuration with conditional logic
- Comprehensive test coverage (80.5% overall, 98.6% internal package)

## Test Coverage
- **43+ tests** covering all functionality
- **viper_comprehensive_test.go** - Tests all 20+ Viper Get functions
- **Type conversion tests** - Lua ↔ Go type mapping
- **Array handling tests** - Array conversion scenarios
- **Edge case tests** - Error conditions and boundary cases
- **Integration tests** - Viper and Cobra integration
