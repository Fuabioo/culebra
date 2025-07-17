package internal

import (
	"github.com/yuin/gopher-lua"
)

func LuaToGo(lv lua.LValue) any {
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
		return luaTableToGoMap(v)
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

func luaTableToGoMap(table *lua.LTable) map[string]any {
	result := make(map[string]any)
	table.ForEach(func(key, value lua.LValue) {
		result[key.String()] = LuaToGo(value)
	})
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
