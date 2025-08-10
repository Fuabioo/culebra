# 📦 Culebra

Culebra, meaning 'snake' in Costa Rica, is a library for loading **Lua scripts as configuration files** and optionally bind them into **Viper** and integrate with **Cobra** CLI apps.

> [!IMPORTANT]
> The lua config will not replace existing cobra options like yml and json, it should work alongside them.

## 🛠️ Features

- Load Lua scripts as dynamic configuration sources with support for both traditional and Neovim-style configurations.
- Return configurations as `map[string]any` or easily bind them into Viper.
- Integrate with Cobra through a simple one-liner: `culebra.UseWithCobra()`.
- Provide flexibility with zero dependencies on Cobra/Viper — their usage is optional.

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

## 🔍 API

- ✅ `Load(cfg Config) (map[string]any, error)` — Loads Lua configuration from file.
- ✅ `LoadWithArrays(filePath string) (map[string]any, error)` — Loads Lua config with array conversion.
- ✅ `LoadWithGlobals(filePath string, globals map[string]any) (map[string]any, error)` — Loads Lua config with global variables.
- ✅ `LoadWithArraysAndGlobals(filePath string, globals map[string]any) (map[string]any, error)` — Loads with both features.
- ✅ `BindToViper(cfg Config, v *viper.Viper) error` — Injects configuration into Viper.
- ✅ `UseWithCobra(cmd *cobra.Command)` — Adds `--config` flag with automatic Lua discovery.
- ✅ Includes comprehensive error handling and validation.
- ✅ Supports both traditional and Neovim-style configuration patterns.

## 📦 Dependencies

- `github.com/yuin/gopher-lua` — Lua VM in Go.
- `github.com/spf13/viper` (optional).
- `github.com/spf13/cobra` (optional).

## 🧪 Usage Examples

### Basic Configuration Loading
```go
// Simple Lua configuration loading
cfg := culebra.Config{FilePath: "config.lua"}
data, err := culebra.Load(cfg)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("App name: %s\n", data["app"].(map[string]any)["name"])
```

### Array Support
```go
// Load with automatic array conversion
data, err := culebra.LoadWithArrays("config.lua")
// Lua arrays become Go slices: []any{"item1", "item2", "item3"}
```

### Global Variables
```go
// Pass Go variables to Lua configuration
globals := map[string]any{
    "environment": "production",
    "debug_enabled": false,
}
data, err := culebra.LoadWithGlobals("config.lua", globals)
```

### Viper Integration
```go
cfg := culebra.Config{
    FilePath: "config.lua",
    ConvertArrays: true,
}
err := culebra.BindToViper(cfg, viper.GetViper())
```

### Cobra CLI Integration
```go
rootCmd := &cobra.Command{
    Use: "myapp",
    Short: "My application with Lua configuration",
    Run: func(cmd *cobra.Command, args []string) {
        // Access configuration via Viper
        fmt.Printf("Database: %s\n", viper.GetString("database.host"))
    },
}

// Adds --config flag with automatic .lua file discovery
culebra.UseWithCobra(rootCmd)
```

### Automatic Discovery
```go
// Configure automatic .lua file discovery
viper.SetConfigName("myapp")  // Searches for myapp.lua
viper.AddConfigPath("/etc/myapp")
viper.AddConfigPath("$HOME/.config/myapp")
viper.AddConfigPath(".")

culebra.UseWithCobra(rootCmd)  // Auto-loads first .lua file found
```

## 📂 Examples

The repository includes comprehensive examples showcasing different use cases:

- **[`examples/basic/`](examples/basic/)** - Basic Lua configuration with Cobra/Viper
- **[`examples/autoload/`](examples/autoload/)** - Automatic configuration discovery  
- **[`examples/arrays/`](examples/arrays/)** - Array conversion and Viper integration
- **[`examples/advanced/`](examples/advanced/)** - Complex, environment-aware configuration
- **[`examples/viper-showcase/`](examples/viper-showcase/)** - Comprehensive demo of ALL Viper Get functions

Run all examples: `just example`

## 🚀 Future Improvements

### HTTP/Network Support
Currently, Culebra supports local configuration logic but doesn't include HTTP capabilities. Future versions could add network support for advanced use cases:

```lua
-- Future capability: Remote configuration fetching
local feature_flags = http.get("https://api.company.com/feature-flags")
local service_discovery = http.get("https://consul.company.com/v1/catalog/services")

config.features = feature_flags.enabled
config.services = service_discovery.endpoints
```

This would enable:
- **Remote Feature Flags**: Dynamic feature toggling from external services
- **Service Discovery**: Automatic service endpoint configuration
- **External Configuration**: Loading partial config from remote sources
- **API Integration**: Direct integration with configuration management APIs

**Implementation Options**:
1. Custom HTTP module for gopher-lua
2. Integration with existing `gluahttp` library  
3. Go-side pre-fetching with Lua globals
4. Async configuration updates

### Other Potential Enhancements
- **Configuration Caching**: Cache remote configurations with TTL
- **Hot Reloading**: Watch for configuration changes and reload automatically  
- **Configuration Validation**: Schema validation for Lua configurations
- **Encrypted Configurations**: Support for encrypted configuration files
- **Configuration Templating**: Advanced templating beyond basic Lua logic

## 📝 License

MIT
