package internal

import (
	"fmt"

	"github.com/yuin/gopher-lua"
)

// ConversionConfig holds configuration for Lua to Go conversion with safety limits
type ConversionConfig struct {
	ConvertArrays bool
	MaxDepth      int
	MaxTableSize  int
}

// LuaToGo converts a Lua value to a Go value without array conversion
func LuaToGo(lv lua.LValue) any {
	return LuaToGoWithConfig(lv, false)
}

// LuaToGoWithConfig converts a Lua value to a Go value with optional array conversion
// This function does not enforce depth limits. Use LuaToGoSafe for production code.
func LuaToGoWithConfig(lv lua.LValue, convertArrays bool) any {
	switch v := lv.(type) {
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(v)
	case lua.LNumber:
		return float64(v)
	case lua.LString:
		return string(v)
	case *lua.LTable:
		if convertArrays && isLuaArray(v) {
			return luaTableToGoSlice(v, convertArrays)
		}
		return luaTableToGoMap(v, convertArrays)
	default:
		return v.String()
	}
}

// LuaToGoSafe converts a Lua value to a Go value with depth and size limits enforced
func LuaToGoSafe(lv lua.LValue, cfg ConversionConfig, currentDepth int) (result any, err error) {
	// Panic protection - convert panics to errors
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic during conversion: %v", r)
			result = nil
		}
	}()

	// Check depth limit
	if currentDepth > cfg.MaxDepth {
		return nil, fmt.Errorf("maximum nesting depth exceeded (%d levels)", cfg.MaxDepth)
	}

	switch v := lv.(type) {
	case *lua.LNilType:
		return nil, nil
	case lua.LBool:
		return bool(v), nil
	case lua.LNumber:
		return float64(v), nil
	case lua.LString:
		return string(v), nil
	case *lua.LTable:
		// Check table size limit if configured
		if cfg.MaxTableSize > 0 {
			size := 0
			v.ForEach(func(_, _ lua.LValue) {
				size++
			})
			if size > cfg.MaxTableSize {
				return nil, fmt.Errorf("table size (%d) exceeds maximum allowed (%d)", size, cfg.MaxTableSize)
			}
		}

		if cfg.ConvertArrays && isLuaArray(v) {
			return luaTableToGoSliceSafe(v, cfg, currentDepth+1)
		}
		return luaTableToGoMapSafe(v, cfg, currentDepth+1)
	default:
		return v.String(), nil
	}
}

func GoToLua(L *lua.LState, value any) lua.LValue {
	switch v := value.(type) {
	case nil:
		return lua.LNil
	case bool:
		return lua.LBool(v)
	case int:
		return lua.LNumber(v)
	case int8:
		return lua.LNumber(v)
	case int16:
		return lua.LNumber(v)
	case int32:
		return lua.LNumber(v)
	case int64:
		return lua.LNumber(v)
	case uint:
		return lua.LNumber(v)
	case uint8:
		return lua.LNumber(v)
	case uint16:
		return lua.LNumber(v)
	case uint32:
		return lua.LNumber(v)
	case uint64:
		return lua.LNumber(v)
	case float32:
		return lua.LNumber(v)
	case float64:
		return lua.LNumber(v)
	case string:
		return lua.LString(v)
	case []byte:
		return lua.LString(string(v))
	case map[string]any:
		return goMapToLuaTable(L, v)
	case []any:
		return goSliceToLuaTable(L, v)
	default:
		// Safe conversion using fmt.Sprintf to avoid panic
		return lua.LString(fmt.Sprintf("%v", v))
	}
}

func luaTableToGoMap(table *lua.LTable, convertArrays bool) map[string]any {
	result := make(map[string]any)
	table.ForEach(func(key, value lua.LValue) {
		result[key.String()] = LuaToGoWithConfig(value, convertArrays)
	})
	return result
}

// luaTableToGoMapSafe converts a Lua table to a Go map with safety limits
func luaTableToGoMapSafe(table *lua.LTable, cfg ConversionConfig, currentDepth int) (map[string]any, error) {
	result := make(map[string]any)
	var convErr error
	table.ForEach(func(key, value lua.LValue) {
		if convErr != nil {
			return // Skip if we already hit an error
		}
		converted, err := LuaToGoSafe(value, cfg, currentDepth)
		if err != nil {
			convErr = err
			return
		}
		result[key.String()] = converted
	})
	if convErr != nil {
		return nil, convErr
	}
	return result, nil
}

// isLuaArray checks if a Lua table is an array (sequential integer keys starting from 1)
func isLuaArray(table *lua.LTable) bool {
	length := table.Len()
	if length == 0 {
		return false
	}

	// Check if all keys from 1 to length exist
	if !hasSequentialKeys(table, length) {
		return false
	}

	// Check that there are no non-array keys
	return !hasNonArrayKeys(table, length)
}

// hasSequentialKeys checks if table has all keys from 1 to length
func hasSequentialKeys(table *lua.LTable, length int) bool {
	for i := 1; i <= length; i++ {
		if table.RawGetInt(i) == lua.LNil {
			return false
		}
	}
	return true
}

// hasNonArrayKeys checks if table has any keys that aren't valid array indices
func hasNonArrayKeys(table *lua.LTable, length int) bool {
	hasOtherKeys := false
	table.ForEach(func(key, _ lua.LValue) {
		if hasOtherKeys {
			return // Early exit if we already found non-array keys
		}

		// Check if key is a number
		if key.Type() != lua.LTNumber {
			hasOtherKeys = true
			return
		}

		// Check if the number is a valid array index
		if keyNum, ok := key.(lua.LNumber); ok {
			keyInt := int(keyNum)
			// Valid array indices are integers from 1 to length
			if keyInt < 1 || keyInt > length || float64(keyInt) != float64(keyNum) {
				hasOtherKeys = true
			}
		}
	})

	return hasOtherKeys
}

// luaTableToGoSlice converts a Lua array table to a Go slice
func luaTableToGoSlice(table *lua.LTable, convertArrays bool) []any {
	length := table.Len()
	result := make([]any, length)

	for i := 1; i <= length; i++ {
		value := table.RawGetInt(i)
		result[i-1] = LuaToGoWithConfig(value, convertArrays)
	}

	return result
}

// luaTableToGoSliceSafe converts a Lua array table to a Go slice with safety limits
func luaTableToGoSliceSafe(table *lua.LTable, cfg ConversionConfig, currentDepth int) ([]any, error) {
	length := table.Len()
	result := make([]any, length)

	for i := 1; i <= length; i++ {
		value := table.RawGetInt(i)
		converted, err := LuaToGoSafe(value, cfg, currentDepth)
		if err != nil {
			return nil, err
		}
		result[i-1] = converted
	}

	return result, nil
}

func goMapToLuaTable(L *lua.LState, m map[string]any) *lua.LTable {
	table := L.NewTable()
	for k, v := range m {
		table.RawSetString(k, GoToLua(L, v))
	}
	return table
}

func goSliceToLuaTable(L *lua.LState, s []any) *lua.LTable {
	table := L.NewTable()
	for i, v := range s {
		table.RawSetInt(i+1, GoToLua(L, v))
	}
	return table
}
