package internal

import (
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
	case int64:
		return lua.LNumber(v)
	case float64:
		return lua.LNumber(v)
	case string:
		return lua.LString(v)
	case map[string]any:
		return goMapToLuaTable(L, v)
	case []any:
		return goSliceToLuaTable(L, v)
	default:
		return lua.LString(v.(string))
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

	// Check if all keys from 1 to length exist and are the only keys
	for i := 1; i <= length; i++ {
		if table.RawGetInt(i) == lua.LNil {
			return false
		}
	}

	// Check that there are no other keys
	hasOtherKeys := false
	table.ForEach(func(key, value lua.LValue) {
		if keyType := key.Type(); keyType != lua.LTNumber {
			hasOtherKeys = true
			return
		}
		if keyNum, ok := key.(lua.LNumber); ok {
			keyInt := int(keyNum)
			if keyInt < 1 || keyInt > length || float64(keyInt) != float64(keyNum) {
				hasOtherKeys = true
				return
			}
		}
	})

	return !hasOtherKeys
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
