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

- ✅ `Load(cfg Config) (map[string]any, error)` — Loads Lua configuration.
- ✅ `BindToViper(cfg Config, v *viper.Viper) error` — Injects configuration into Viper.
- ✅ `UseWithCobra(cmd *cobra.Command)` — Adds a `--config` flag that loads Lua into Viper.
- ✅ Includes basic error handling and logging.
- ✅ Comes with an example CLI app utilizing `cobra` and `viper` alongside Lua configurations.

## 📦 Dependencies

- `github.com/yuin/gopher-lua` — Lua VM in Go.
- `github.com/spf13/viper` (optional).
- `github.com/spf13/cobra` (optional).

## 🧪 Usage

```go
// Basic usage
cfg := culebra.Config{FilePath: "config.lua"}
data, err := culebra.Load(cfg)
```

```go
// With Viper
err = culebra.BindToViper(cfg, viper.GetViper())
```

```go
// With Cobra
rootCmd := &cobra.Command{
    Use: "myapp",
    Short: "A brief description of your application",
}

// Integrate with Culebra
culebra.UseWithCobra(rootCmd)
```

```go
// With Cobra and Viper using autoload
rootCmd := &cobra.Command{
    Use: "myapp",
    Short: "A brief description of your application",
}

// Configure Viper for autoload - this triggers culebra's autoload mechanism
// Setting config name enables autoload for .lua files
viper.SetConfigName("example")
// Adding config paths works with autoload - culebra will search this path for .lua files
viper.AddConfigPath("$HOME/.config/example")

// Enable Cobra integration to load automatically $HOME/.config/example/example.lua
culebra.UseWithCobra(rootCmd)
```

## 📝 License

MIT
