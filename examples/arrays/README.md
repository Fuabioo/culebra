# Array Conversion Example

This example demonstrates Culebra's enhanced array conversion capabilities, showing how Lua arrays can be seamlessly converted to Go slices for better integration with Viper and struct unmarshaling.

## Features Demonstrated

- **Backward Compatibility**: Traditional map-based array handling still works
- **Array Conversion**: New `ConvertArrays: true` option converts Lua arrays to Go slices
- **Viper Integration**: Seamless integration with Viper's `GetStringSlice()`, `GetIntSlice()`, etc.
- **Struct Unmarshaling**: Direct unmarshaling into Go structs with slice fields
- **Nested Structures**: Arrays within nested objects
- **Mixed Types**: Arrays containing different data types

## Configuration Files

- `config-traditional.lua` - Traditional style with global variables
- `config-neovim-style.lua` - Modern return-style configuration
- `config-complex.lua` - Complex nested structures with arrays

## Usage

```bash
go run main.go
```

This will demonstrate:
1. Loading configs without array conversion (backward compatibility)
2. Loading configs with array conversion (new feature)
3. Viper integration with automatic array handling
4. Struct unmarshaling with slice fields

## Key Benefits

- **Before**: Arrays became `map[string]interface{}` with keys "1", "2", "3"
- **After**: Arrays become proper `[]interface{}` slices
- **Result**: Perfect compatibility with Viper and Go's type system