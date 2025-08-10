package internal

import (
	"fmt"

	"github.com/yuin/gopher-lua"
)

func LuaToGo(lv lua.LValue) any {
	return LuaToGoWithConfig(lv, false)
}

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
