# Autoload Example

This example demonstrates automatic Lua configuration loading using culebra's built-in auto-discovery mechanism.

## How it works

1. **Automatic Discovery**: When `viper.SetConfigName()` is called, culebra automatically searches for `{configname}.lua` in the configured search paths and loads it.

2. **Search Path Support**: `viper.AddConfigPath()` adds directories to search for `.lua` config files, just like Viper's standard config discovery.

3. **Environment-based Configuration**: The Lua config uses `os.getenv("APP_ENV")` to dynamically adjust settings based on the environment.

4. **No Explicit File Path**: You don't need to specify a config file - it's discovered automatically from the search paths.

## Running the example

```bash
# Run with default (development) environment
cd examples/autoload
go run main.go

# Run in development mode (explicit)
APP_ENV=development go run main.go

# Run in production mode  
APP_ENV=production go run main.go

# Override with explicit config file
go run main.go --config simple.lua
```

## Example output

**Development mode:**
```
App Name: Auto-loaded App
App Version: 2.0.0
Database Host: localhost
Database Port: 5433
Debug Mode: false
Environment: development
```

**Production mode:**
```
App Name: Auto-loaded App
App Version: 2.0.0
Database Host: prod-db.company.com
Database Port: 5432
Debug Mode: false
Environment: production
```

## Config file discovery

The example searches for `example.lua` in the following paths (in order):
1. `/etc/example.lua` (system-wide)
2. `$HOME/.config/example.lua` (user-specific)  
3. `./example.lua` (current directory)

The first `.lua` file found will be loaded automatically.

## Autoload behavior

### With SetConfigName (autoload enabled):
```go
viper.SetConfigName("example")
// Optional: Add search paths (defaults to current directory if none specified)
viper.AddConfigPath("/etc")
viper.AddConfigPath("$HOME/.config")
viper.AddConfigPath(".")
culebra.UseWithCobra(rootCmd)
```

### Without SetConfigName (no autoload):
```go
// No autoload - only explicit --config flag works
culebra.UseWithCobra(rootCmd)
```

## Config files in this example

- **`example.lua`**: Main config file with full features (auto-loaded by default)

## Features demonstrated

- **Automatic Discovery**: Finds `.lua` files in search paths when `SetConfigName` is used
- **Search Path Support**: Works with `AddConfigPath` for custom search locations
- **Environment variables**: Dynamic configuration based on `APP_ENV`
- **Lua logic**: Conditional configuration using `if/then` statements
- **Nested structures**: Complex configuration with tables and sub-tables
- **Type handling**: Strings, integers, booleans, and tables
- **Explicit override**: `--config` flag still works for manual file specification