# 📦 Culebra

Load **Lua scripts as configuration files** and optionally bind them into **Viper** and integrate with **Cobra** CLI apps.

> IMPORTANT! The lua config will not replace existing cobra options like yml and json, it should work alongside them.

## 🏗️ Project Layout

```
culebra/
├── go.mod
├── README.md
├── luaconf.go # Core logic to load Lua configs
├── viper.go # Bind Lua configs into Viper
├── cobra.go # One-liner integration with Cobra
├── lua/
│ └── stdlib.lua # Optional Lua helpers
├── examples/
│ ├── basic/
│ │ ├── main.go
│ │ ├── config.lua # Traditional style
│ │ └── config-neovim-style.lua # Neovim-style with return
└── internal/
└── mapper.go # Lua table → Go map
```

## 🛠️ Features

- Load `.lua` files as dynamic config sources.
- Return as `map[string]any` or bind into Viper.
- Seamlessly plug into Cobra with one-liner (`culebra.UseWithCobra()`).
- Zero dependencies on Cobra/Viper — integration is optional.
- **Supports both traditional and Neovim-style configurations**

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

- ✅ `Load(cfg Config) (map[string]any, error)` — loads Lua config
- ✅ `BindToViper(cfg Config, v *viper.Viper) error` — injects config into Viper
- ✅ `UseWithCobra(cmd *cobra.Command)` — adds `--config` flag, loads Lua into Viper
- ✅ Basic error wrapping/logging
- ✅ Example CLI app using `cobra` and `viper` + Lua config

## 📦 Dependencies

- `github.com/yuin/gopher-lua` — Lua VM in Go
- `github.com/spf13/viper` (optional)
- `github.com/spf13/cobra` (optional)

## 🧪 Usage

```go
// Basic usage
cfg := culebra.Config{FilePath: "config.lua"}
data, err := culebra.Load(cfg)

// With Viper
err = culebra.BindToViper(cfg, viper.GetViper())

// With Cobra (one-liner)
culebra.UseWithCobra(rootCmd)
```

## 📝 License

MIT