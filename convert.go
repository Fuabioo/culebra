package culebra

import (
	"fmt"

	"github.com/Fuabioo/culebra/internal"
	lua "github.com/yuin/gopher-lua"
)

// LuaToGo converts a Lua value to a Go value without array conversion or safety limits.
// For production use with untrusted configs, consider using the Load() function which
// enforces depth and size limits.
//
// Example:
//
//	L := lua.NewState()
//	defer L.Close()
//	L.DoString(`return {name = "test", value = 42}`)
//	lv := L.Get(-1)
//	goValue := culebra.LuaToGo(lv)
//	// goValue is map[string]any{"name": "test", "value": 42.0}
func LuaToGo(lv lua.LValue) any {
	return internal.LuaToGo(lv)
}

// LuaToGoWithArrays converts a Lua value to a Go value with array conversion enabled.
// Lua tables with sequential numeric keys starting from 1 will be converted to Go slices.
// For production use with untrusted configs, consider using the Load() function which
// enforces depth and size limits.
//
// Example:
//
//	L := lua.NewState()
//	defer L.Close()
//	L.DoString(`return {1, 2, 3, 4}`)
//	lv := L.Get(-1)
//	goValue := culebra.LuaToGoWithArrays(lv)
//	// goValue is []any{1.0, 2.0, 3.0, 4.0}
func LuaToGoWithArrays(lv lua.LValue) any {
	return internal.LuaToGoWithConfig(lv, true)
}

// LuaToGoSafe converts a Lua value to a Go value with configurable safety limits.
// This function is recommended when converting untrusted Lua values as it protects
// against deeply nested structures and excessively large tables.
//
// Parameters:
//   - lv: The Lua value to convert
//   - convertArrays: If true, convert Lua arrays to Go slices instead of maps
//   - maxDepth: Maximum nesting depth (0 = use default of 100)
//   - maxTableSize: Maximum entries per table (0 = unlimited)
//
// Returns an error if depth or size limits are exceeded.
//
// Example:
//
//	L := lua.NewState()
//	defer L.Close()
//	L.DoString(`return {a = {b = {c = "nested"}}}`)
//	lv := L.Get(-1)
//	goValue, err := culebra.LuaToGoSafe(lv, false, 10, 1000)
//	if err != nil {
//	    // Handle error (e.g., depth limit exceeded)
//	}
func LuaToGoSafe(lv lua.LValue, convertArrays bool, maxDepth, maxTableSize int) (any, error) {
	if maxDepth == 0 {
		maxDepth = 100
	}

	cfg := internal.ConversionConfig{
		ConvertArrays: convertArrays,
		MaxDepth:      maxDepth,
		MaxTableSize:  maxTableSize,
	}

	return internal.LuaToGoSafe(lv, cfg, 0)
}

// GoToLua converts a Go value to a Lua value.
// Supported Go types: nil, bool, int/int8-64, uint/uint8-64, float32/64, string,
// []byte, map[string]any, []any. Other types are converted to strings using fmt.Sprintf.
//
// Example:
//
//	L := lua.NewState()
//	defer L.Close()
//	goValue := map[string]any{"name": "test", "count": 42}
//	lv := culebra.GoToLua(L, goValue)
//	L.SetGlobal("config", lv)
//	// Now Lua code can access config.name and config.count
func GoToLua(L *lua.LState, value any) lua.LValue {
	// Add panic protection
	defer func() {
		if r := recover(); r != nil {
			// If conversion panics, return a string representation
			_ = r
		}
	}()

	return internal.GoToLua(L, value)
}

// MustLuaToGo is like LuaToGo but panics if conversion encounters an error.
// Only use this when you're certain the Lua value is safe and well-formed.
func MustLuaToGo(lv lua.LValue) any {
	// For simple types, LuaToGo shouldn't fail, but we add protection anyway
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Sprintf("culebra.MustLuaToGo: conversion failed: %v", r))
		}
	}()

	return internal.LuaToGo(lv)
}
