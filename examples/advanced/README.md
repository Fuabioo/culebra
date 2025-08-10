# Advanced Lua Configuration Example

This example demonstrates the real advantages of using Lua for configuration over static formats like JSON or YAML.

## Why Lua Over JSON/YAML?

### 1. **Environment-Dependent Configuration**
JSON/YAML require external templating or multiple files. Lua handles this natively:

```lua
-- Dynamic configuration based on environment
local environment = os.getenv("APP_ENV") or "development"
local is_production = environment == "production"

config.database = is_production and {
    url = os.getenv("DATABASE_URL"),
    pool_size = 20,
    ssl_mode = "require"
} or {
    url = "postgres://localhost:5432/myapp_dev", 
    pool_size = 5,
    ssl_mode = "disable"
}
```

### 2. **Computed Values & Logic**
Generate configuration based on conditions, calculations, or external data:

```lua
-- Service URLs computed dynamically
local function service_url(service, port)
    local host_suffix = is_production and ".prod.internal" or ".dev.local"
    return string.format("http://%s%s:%d", service, host_suffix, port)
end

services.user_service = service_url("user-service", 8001)
```

### 3. **Configuration Validation**
Built-in validation ensures configuration correctness:

```lua
local function validate_config()
    if is_production then
        local required_vars = {"DATABASE_URL", "SESSION_SECRET"}
        for _, var in ipairs(required_vars) do
            if not os.getenv(var) then
                error("Missing required environment variable: " .. var)
            end
        end
    end
    return true
end

validate_config()
```

### 4. **DRY (Don't Repeat Yourself)**
Helper functions eliminate repetition:

```lua
-- Helper function used multiple times
local function monitoring_url(path)
    local base = is_production and "https://monitoring.company.com" or "http://localhost:3000"
    return base .. path
end

-- Generate multiple endpoints
monitoring_endpoints = {
    monitoring_url("/health"),
    monitoring_url("/metrics"),
    monitoring_url("/dashboard")
}
```

### 5. **Time/Context-Aware Configuration**
Configuration can adapt based on runtime conditions:

```lua
-- Dynamic maintenance mode based on time
local current_hour = tonumber(os.date("%H"))
config.maintenance_mode = is_production and current_hour >= 2 and current_hour <= 4
```

## Running the Example

### Basic Usage
```bash
cd examples/advanced
go run main.go --config config.lua
```

### Different Environments
```bash
# Development (default)
go run main.go --config config.lua

# Staging
APP_ENV=staging go run main.go --config config.lua

# Production
APP_ENV=production DATABASE_URL=postgres://prod:5432/myapp SESSION_SECRET=supersecret go run main.go --config config.lua

# With experimental features
ENABLE_EXPERIMENTAL=true go run main.go --config config.lua
```

## Configuration Features Demonstrated

1. **Environment Detection**: Automatic environment detection from `APP_ENV`
2. **Conditional Logic**: Different settings per environment
3. **Computed Values**: Dynamic URL generation, port calculation
4. **Feature Flags**: Environment-based feature enabling
5. **Validation**: Runtime configuration validation
6. **Helper Functions**: Reusable configuration logic
7. **External Integration**: Reading environment variables
8. **Time-Based Logic**: Maintenance mode scheduling

## JSON/YAML Equivalent

To achieve the same functionality with JSON/YAML, you would need:

- External templating engine (Helm, Kustomize, etc.)
- Multiple configuration files per environment
- Environment variable substitution tools
- Separate validation scripts
- Manual DRY violation or complex inheritance

With Lua, all this logic is contained in a single, readable configuration file that adapts automatically to your environment and requirements.

## Key Takeaways

Lua configuration provides:

- ✅ **Programmability**: Full programming language features
- ✅ **Environment Awareness**: Native OS interaction
- ✅ **Logic & Computation**: Dynamic value generation  
- ✅ **Validation**: Built-in error checking
- ✅ **DRY Principle**: Reusable helper functions
- ✅ **Single Source**: One file for all environments
- ✅ **Maintainability**: Clear, readable configuration logic

This approach reduces configuration complexity while increasing flexibility and maintainability.