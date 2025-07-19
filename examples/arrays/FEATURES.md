# Culebra Array Conversion Features

## 🆕 Enhanced Array Support

Culebra now provides seamless array conversion capabilities that bridge the gap between Lua's table-based arrays and Go's native slice types, enabling perfect integration with Viper and Go's type system.

## 🔄 Backward Compatibility

**✅ 100% Backward Compatible**: All existing code continues to work unchanged.

### Before (Traditional Behavior)
```lua
-- config.lua
database_hosts = {"localhost", "db1.example.com", "db2.example.com"}
```

```go
// Go code - arrays become maps with string keys
data, _ := culebra.Load(culebra.Config{FilePath: "config.lua"})
hosts := data["database_hosts"].(map[string]interface{})
// hosts = map["1":"localhost" "2":"db1.example.com" "3":"db2.example.com"]
firstHost := hosts["1"] // "localhost"
```

### After (Enhanced Behavior)
```lua
-- config.lua
return {
    database = {
        hosts = {"localhost", "db1.example.com", "db2.example.com"}
    }
}
```

```go
// Go code - arrays become proper slices
data, _ := culebra.LoadWithArrays("config.lua")
db := data["database"].(map[string]interface{})
hosts := db["hosts"].([]interface{})
// hosts = ["localhost", "db1.example.com", "db2.example.com"]
firstHost := hosts[0] // "localhost"
```

## 🚀 New API Functions

### Core Loading Functions
```go
// Load with array conversion enabled
culebra.LoadWithArrays(filePath string) (map[string]any, error)

// Load with global variables
culebra.LoadWithGlobals(filePath string, globals map[string]any) (map[string]any, error)

// Load with both features
culebra.LoadWithArraysAndGlobals(filePath string, globals map[string]any) (map[string]any, error)
```

### Enhanced Viper Integration
```go
// Automatic array conversion for Viper
culebra.AutoBindToViper(cfg Config, v *viper.Viper) error

// Convenience function for common use case
culebra.BindToViperWithArrays(filePath string, v *viper.Viper) error
```

### Configuration Control
```go
// Fine-grained control over array conversion
culebra.Load(culebra.Config{
    FilePath:      "config.lua",
    ConvertArrays: true,  // Enable array conversion
    Globals:       globals,
})
```

## 🎯 Perfect Viper Integration

### Before: Manual Array Handling Required
```go
// Complex manual parsing needed
data, _ := culebra.Load(culebra.Config{FilePath: "config.lua"})
if hostsMap, ok := data["hosts"].(map[string]interface{}); ok {
    hosts := make([]string, len(hostsMap))
    for i := 1; i <= len(hostsMap); i++ {
        hosts[i-1] = hostsMap[fmt.Sprintf("%d", i)].(string)
    }
    // Now you have a slice...
}
```

### After: Direct Viper Usage
```go
// One line with perfect integration
v := viper.New()
culebra.AutoBindToViper(culebra.Config{FilePath: "config.lua"}, v)
hosts := v.GetStringSlice("database.hosts") // Just works!
ports := v.GetIntSlice("database.ports")     // Multiple types supported
```

## 📦 Struct Unmarshaling Support

### Configuration Struct
```go
type Config struct {
    Database DatabaseConfig `mapstructure:"database"`
    Services []Service      `mapstructure:"services"`
    Features []string       `mapstructure:"features"`
}

type DatabaseConfig struct {
    Hosts    []string `mapstructure:"hosts"`
    Ports    []int    `mapstructure:"ports"`
    Timeouts []int    `mapstructure:"timeouts"`
}
```

### Seamless Unmarshaling
```go
v := viper.New()
culebra.AutoBindToViper(culebra.Config{FilePath: "config.lua"}, v)

var config Config
v.Unmarshal(&config) // Arrays automatically populate slice fields!

fmt.Printf("Hosts: %v\n", config.Database.Hosts)
fmt.Printf("Services: %d\n", len(config.Services))
```

## 🔍 Array Detection Logic

Culebra intelligently detects Lua arrays using these criteria:

1. **Sequential Integer Keys**: Keys must be consecutive integers starting from 1
2. **No Gaps**: All indices from 1 to length must exist  
3. **No Extra Keys**: No non-numeric or out-of-range keys
4. **Automatic Conversion**: When `ConvertArrays: true`, arrays become `[]interface{}`

### Examples

```lua
-- ✅ Detected as array
valid_array = {"a", "b", "c"}

-- ✅ Detected as array  
mixed_types = {1, "two", true, 4.5}

-- ❌ Not an array (gap in indices)
invalid_array = {[1] = "a", [3] = "c"}

-- ❌ Not an array (non-numeric keys)
object = {name = "value", count = 42}

-- ✅ Nested arrays work perfectly
complex = {
    servers = {"web1", "web2", "web3"},
    config = {
        timeouts = {30, 60, 120},
        features = {"ssl", "cache", "monitor"}
    }
}
```

## 🛡️ Migration Guide

### Existing Users (No Changes Required)
Your existing code continues to work exactly as before:

```go
// This still works identically
data, err := culebra.Load(culebra.Config{FilePath: "config.lua"})
// Arrays are still maps with string keys
```

### New Users (Recommended Approach)
Take advantage of enhanced array support:

```go
// For direct usage
data, err := culebra.LoadWithArrays("config.lua")

// For Viper integration (recommended)
v := viper.New()
err := culebra.AutoBindToViper(culebra.Config{FilePath: "config.lua"}, v)
```

### Gradual Migration
Enable array conversion selectively:

```go
// Test new behavior without breaking existing code
data, err := culebra.Load(culebra.Config{
    FilePath:      "config.lua",
    ConvertArrays: true, // Opt-in to new behavior
})
```

## ⚡ Performance Notes

- **Zero Overhead**: When `ConvertArrays: false` (default), no additional processing occurs
- **Efficient Detection**: Array detection uses O(n) time complexity where n is the number of keys
- **Memory Efficient**: Direct conversion without intermediate allocations
- **Lazy Evaluation**: Only converts arrays when `ConvertArrays: true`

## 🧪 Testing & Validation

The example includes comprehensive validation:

- **Backward Compatibility Tests**: Ensure existing behavior unchanged
- **Array Conversion Tests**: Verify proper slice creation
- **Viper Integration Tests**: Test seamless Viper usage
- **Struct Unmarshaling Tests**: Validate direct struct mapping
- **Complex Configuration Tests**: Test deeply nested structures

Run validation: `./validate.sh`

## 🎉 Summary

The enhanced Culebra provides:

✅ **Perfect Backward Compatibility** - No breaking changes  
✅ **Seamless Viper Integration** - Arrays work naturally with `GetStringSlice()`, etc.  
✅ **Direct Struct Unmarshaling** - No manual array conversion needed  
✅ **Intelligent Array Detection** - Automatic differentiation between arrays and objects  
✅ **Flexible Configuration** - Opt-in array conversion with fine-grained control  
✅ **Comprehensive Testing** - Full validation suite ensures reliability  

This enhancement makes Culebra a **first-class citizen** in the Cobra + Viper ecosystem! 🚀