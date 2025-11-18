# 📦 Culebra

Culebra, meaning 'snake' in Costa Rica, is a library for loading **Lua scripts as configuration files** with seamless integration for **Viper** and **Cobra** CLI apps.

> [!IMPORTANT]
> Lua config works **alongside** existing config formats (.yml, .json) - it doesn't replace them.

## ⚠️ Security Notice

**Lua configurations have full access to the Lua standard library**, including file I/O (`io.*`), system commands (`os.*`), and environment variables. Only load configuration files from **trusted sources**.

Culebra provides safety limits:
- **MaxDepth**: Prevents stack overflow from deeply nested tables (default: 100 levels)
- **MaxTableSize**: Limits entries per table to prevent memory exhaustion (default: unlimited)
- **Panic Protection**: All public functions recover from panics and return errors

For untrusted configs, set appropriate limits:
```go
cfg := culebra.Config{
    FilePath:     "config.lua",
    MaxDepth:     50,      // Max 50 levels of nesting
    MaxTableSize: 10000,   // Max 10k entries per table
}
```

## 🛠️ Features

- ✅ **Plug-and-play Cobra/Viper integration** - Works like Viper's auto-loading for .json/.yml
- ✅ **Zero dependencies for core** - Cobra/Viper are optional
- ✅ **Two config styles** - Traditional (global vars) and Neovim-style (return statement)
- ✅ **Array conversion** - Lua arrays become Go slices
- ✅ **Global variables** - Pass Go values into Lua configs
- ✅ **Type conversion API** - Direct Lua ↔ Go conversion functions
- ✅ **Defensive programming** - Depth limits, panic protection, clear errors

## 📦 Installation

```bash
go get github.com/Fuabioo/culebra
```

## 🚀 Quick Start

Choose your starting point:

| Example | When to use | What it shows |
|---------|-------------|---------------|
| [basic/](examples/basic/) | **Start here** - First time using Culebra | Basic Lua config with Cobra/Viper integration |
| [autoload/](examples/autoload/) | You want automatic .lua file discovery | How Viper search paths work with .lua files |
| [arrays/](examples/arrays/) | Your config has lists/arrays | Array conversion and Viper integration |
| [advanced/](examples/advanced/) | Complex, environment-aware configs | Functions, conditionals, environment logic |
| [viper-showcase/](examples/viper-showcase/) | Using all Viper Get* functions | Comprehensive demo of Viper integration |

Run all examples: `just example`

## 📝 Configuration Styles

### Traditional Style (Global Variables)
```lua
-- config.lua
app = {
    name = "My App",
    version = "1.0.0"
}

database = {
    host = "localhost",
    port = 5432
}

debug = true

-- Environment-based logic
if os.getenv("ENVIRONMENT") == "production" then
    debug = false
    database.host = "prod-db.example.com"
end
```

### Neovim Style (Return Statement)
```lua
-- config-neovim-style.lua
local config = {}

local function setup_database()
    return {
        host = os.getenv("DB_HOST") or "localhost",
        port = tonumber(os.getenv("DB_PORT") or "5432"),
        ssl_mode = os.getenv("ENVIRONMENT") == "production" and "require" or "disable"
    }
end

config.app = {
    name = "My App",
    version = "1.0.0"
}

config.database = setup_database()
config.features = {
    debug_mode = os.getenv("ENVIRONMENT") ~= "production"
}

-- Validation
assert(config.database.host, "Database host is required")

return config
```

## 🔍 API Reference

### Core Functions

#### `Load(cfg Config) (map[string]any, error)`
Loads Lua configuration from a file with full safety protections.

```go
cfg := culebra.Config{
    FilePath:      "config.lua",
    ConvertArrays: true,        // Convert Lua arrays to Go slices
    MaxDepth:      100,          // Max nesting depth (0 = default 100)
    MaxTableSize:  10000,        // Max entries per table (0 = unlimited)
    Globals: map[string]any{    // Optional: Go values accessible in Lua
        "environment": "production",
    },
}
data, err := culebra.Load(cfg)
```

#### `BindToViper(cfg Config, v *viper.Viper) error`
Loads Lua config and injects it into Viper.

```go
cfg := culebra.Config{
    FilePath:      "config.lua",
    ConvertArrays: true,
}
err := culebra.BindToViper(cfg, viper.GetViper())
// Now use viper.GetString("database.host"), etc.
```

#### `EnableLuaConfig(cmd *cobra.Command)`
Adds seamless Lua config support to Cobra. Works like Viper's auto-loading for .json/.yml files.

```go
rootCmd := &cobra.Command{
    Use:   "myapp",
    Short: "My application",
}

// One-liner: adds --config flag + automatic .lua discovery
culebra.EnableLuaConfig(rootCmd)
```

### Type Conversion Functions

#### `LuaToGo(lv lua.LValue) any`
Converts Lua value to Go without array conversion.

```go
L := lua.NewState()
defer L.Close()
L.DoString(`return {name = "test", value = 42}`)
lv := L.Get(-1)
goValue := culebra.LuaToGo(lv)
// goValue is map[string]any{"name": "test", "value": 42.0}
```

#### `LuaToGoWithArrays(lv lua.LValue) any`
Converts Lua value to Go with array conversion.

```go
L.DoString(`return {1, 2, 3, 4}`)
lv := L.Get(-1)
goValue := culebra.LuaToGoWithArrays(lv)
// goValue is []any{1.0, 2.0, 3.0, 4.0}
```

#### `LuaToGoSafe(lv lua.LValue, convertArrays bool, maxDepth, maxTableSize int) (any, error)`
Converts with configurable safety limits (recommended for untrusted data).

```go
goValue, err := culebra.LuaToGoSafe(lv, true, 50, 1000)
if err != nil {
    // Handle depth/size limit errors
}
```

#### `GoToLua(L *lua.LState, value any) lua.LValue`
Converts Go value to Lua.

```go
goValue := map[string]any{"name": "test", "count": 42}
lv := culebra.GoToLua(L, goValue)
L.SetGlobal("config", lv)
// Lua can now access config.name and config.count
```

### Config Struct

```go
type Config struct {
    FilePath      string         // Required: path to .lua file
    Globals       map[string]any // Optional: Go values accessible in Lua
    ConvertArrays bool           // Convert Lua arrays to Go slices
    MaxDepth      int            // Max nesting depth (0 = default 100)
    MaxTableSize  int            // Max entries per table (0 = unlimited)
}
```

## 🧪 Usage Examples

### Basic Configuration Loading
```go
cfg := culebra.Config{FilePath: "config.lua"}
data, err := culebra.Load(cfg)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("App name: %s\n", data["app"].(map[string]any)["name"])
```

### Array Support
```go
cfg := culebra.Config{
    FilePath:      "config.lua",
    ConvertArrays: true,
}
data, err := culebra.Load(cfg)
// Lua arrays become Go slices: []any{"item1", "item2", "item3"}
```

### Global Variables
```go
cfg := culebra.Config{
    FilePath: "config.lua",
    Globals: map[string]any{
        "environment":   "production",
        "debug_enabled": false,
    },
}
data, err := culebra.Load(cfg)
// Lua code can access 'environment' and 'debug_enabled' variables
```

### Viper Integration
```go
cfg := culebra.Config{
    FilePath:      "config.lua",
    ConvertArrays: true,
}
err := culebra.BindToViper(cfg, viper.GetViper())

// Access via Viper
dbHost := viper.GetString("database.host")
dbPort := viper.GetInt("database.port")
```

### Cobra CLI Integration
```go
rootCmd := &cobra.Command{
    Use:   "myapp",
    Short: "My application with Lua configuration",
    Run: func(cmd *cobra.Command, args []string) {
        // Access configuration via Viper
        fmt.Printf("Database: %s\n", viper.GetString("database.host"))
    },
}

// Adds --config flag with automatic .lua file discovery
culebra.EnableLuaConfig(rootCmd)
```

### Automatic Discovery
```go
// Configure Viper search paths
viper.SetConfigName("myapp")  // Will search for myapp.lua
viper.AddConfigPath("/etc/myapp")
viper.AddConfigPath("$HOME/.config/myapp")
viper.AddConfigPath(".")

// EnableLuaConfig auto-loads first .lua file found in search paths
culebra.EnableLuaConfig(rootCmd)
```

### With Safety Limits
```go
cfg := culebra.Config{
    FilePath:     "untrusted-config.lua",
    MaxDepth:     50,      // Prevent deeply nested structures
    MaxTableSize: 10000,   // Prevent huge tables
}
data, err := culebra.Load(cfg)
if err != nil {
    // Will return error if limits exceeded
    log.Printf("Config rejected: %v", err)
}
```

## 📦 Dependencies

- `github.com/yuin/gopher-lua` — **Required** for core functionality
- `github.com/spf13/viper` — **Optional** (only if using Viper integration)
- `github.com/spf13/cobra` — **Optional** (only if using Cobra integration)

**Total dependencies**: 1 required + ~148 transitive (if using Viper/Cobra)

The core `Load()` function has **zero dependencies** beyond gopher-lua.

## 🧭 Architecture Notes

### Thread Safety
- `Load()` creates a new Lua VM per call - safe for concurrent use
- `BindToViper()` uses Viper which is **not thread-safe**
- Don't call `BindToViper()` concurrently on the same Viper instance

### Performance
- Each `Load()` creates a new Lua VM (~1-2ms overhead)
- Designed for one-time config loading, not high-frequency operations
- No caching - configs are re-parsed on every `Load()` call

### Reflection in Cobra Integration
`EnableLuaConfig()` uses reflection to access Viper's private fields (`configName`, `configPaths`). This:
- May break if Viper changes its internals
- Is protected with panic recovery
- Fails gracefully (falls back to current directory)
- Is only needed for auto-discovery; explicit `--config` flag always works

## 📄 License

MIT

## 🗺️ Roadmap

See [ROADMAP.md](ROADMAP.md) for planned features and improvements.
