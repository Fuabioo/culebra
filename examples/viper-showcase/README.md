# Viper Functions Showcase

This example provides a comprehensive demonstration of **ALL** Viper `Get` functions working with Lua configuration loaded through Culebra.

## Viper Functions Demonstrated

### Basic Get Functions
- **`GetString(key)`** - String values, type conversion, empty strings, missing keys
- **`GetInt(key)`** - Integer values, type conversion from float/string
- **`GetBool(key)`** - Boolean values, string conversion ("true"/"false")  
- **`GetFloat64(key)`** - Float values, integer conversion to float
- **`Get(key)`** - Generic access returning interface{} with actual types

### Integer Type Variants
- **`GetInt32(key)`** - 32-bit signed integers
- **`GetInt64(key)`** - 64-bit signed integers (for large numbers)
- **`GetUint(key)`** - Unsigned integers
- **`GetUint32(key)`** - 32-bit unsigned integers  
- **`GetUint64(key)`** - 64-bit unsigned integers

### Array/Slice Functions
- **`GetStringSlice(key)`** - Arrays converted to string slices
- **`GetIntSlice(key)`** - Arrays converted to integer slices

### Time & Duration Functions
- **`GetDuration(key)`** - Parses duration strings ("5m30s", "1h", etc.)
- **`GetTime(key)`** - Parses time strings (RFC3339 format)

### Map Functions
- **`GetStringMap(key)`** - Nested objects as `map[string]interface{}`
- **`GetStringMapString(key)`** - Nested objects as `map[string]string`

## Key Features

### 🔧 **Type Conversion**
Lua's dynamic typing works seamlessly with Viper's type conversion:
```lua
port = 8080                    -- Number
debug = true                   -- Boolean  
timeout = "30s"               -- Duration string
hosts = {"a", "b", "c"}       -- Array → Slice
```

### 🎯 **Array Conversion**  
With `ConvertArrays: true`, Lua arrays become proper Go slices:
```lua
hosts = {"db1.com", "db2.com"}  -- Becomes []string
ports = {80, 443, 8080}        -- Becomes []int
```

### 🌳 **Nested Access**
Deep nested structures work with dot notation:
```lua
database = {
    primary = {
        host = "primary.db.com",
        port = 5432
    }
}
-- Access: viper.GetString("database.primary.host")
```

### ⚠️ **Edge Cases**
Handles edge cases gracefully:
- Missing keys return zero values
- Type mismatches convert appropriately  
- nil values handled correctly
- Empty arrays return empty slices

## Running the Example

```bash
cd examples/viper-showcase
go run main.go --config config.lua
```

This will demonstrate every Viper Get function with comprehensive examples showing:
- Basic type access
- Type conversions
- Array handling  
- Nested object access
- Duration/time parsing
- Map access patterns
- Edge case handling

## Configuration Coverage

The `config.lua` file includes:
- ✅ All primitive types (string, int, float, bool)
- ✅ Arrays of various types  
- ✅ Nested objects/maps
- ✅ Duration strings
- ✅ Timestamp strings
- ✅ Large numbers (int64 range)
- ✅ Edge cases (nil, zero, negative, empty)
- ✅ Type conversion scenarios

This example serves as both a demonstration and a test suite for Viper integration with Lua configuration.