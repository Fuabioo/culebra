package internal

import (
	"testing"

	"github.com/yuin/gopher-lua"
)

func TestLuaToGo(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tests := []struct {
		name     string
		input    lua.LValue
		expected any
	}{
		{"nil", lua.LNil, nil},
		{"bool true", lua.LBool(true), true},
		{"bool false", lua.LBool(false), false},
		{"number", lua.LNumber(42.5), 42.5},
		{"string", lua.LString("hello"), "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LuaToGo(tt.input)
			if result != tt.expected {
				t.Errorf("LuaToGo(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGoToLua(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tests := []struct {
		name     string
		input    any
		expected lua.LValueType
	}{
		{"nil", nil, lua.LTNil},
		{"bool", true, lua.LTBool},
		{"int", 42, lua.LTNumber},
		{"int64", int64(42), lua.LTNumber},
		{"float64", 42.5, lua.LTNumber},
		{"string", "hello", lua.LTString},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GoToLua(L, tt.input)
			if result.Type() != tt.expected {
				t.Errorf("GoToLua(%v) type = %v, want %v", tt.input, result.Type(), tt.expected)
			}
		})
	}
}

func TestLuaTableToGoMap(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.NewTable()
	table.RawSetString("key1", lua.LString("value1"))
	table.RawSetString("key2", lua.LNumber(42))
	table.RawSetString("key3", lua.LBool(true))

	result := luaTableToGoMap(table, false)

	if result["key1"] != "value1" {
		t.Errorf("Expected key1='value1', got %v", result["key1"])
	}
	if result["key2"] != float64(42) {
		t.Errorf("Expected key2=42, got %v", result["key2"])
	}
	if result["key3"] != true {
		t.Errorf("Expected key3=true, got %v", result["key3"])
	}
}

func TestGoMapToLuaTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	input := map[string]any{
		"key1": "value1",
		"key2": 42,
		"key3": true,
	}

	result := goMapToLuaTable(L, input)

	if result.RawGetString("key1").String() != "value1" {
		t.Errorf("Expected key1='value1', got %v", result.RawGetString("key1"))
	}
	if float64(result.RawGetString("key2").(lua.LNumber)) != 42 {
		t.Errorf("Expected key2=42, got %v", result.RawGetString("key2"))
	}
	if bool(result.RawGetString("key3").(lua.LBool)) != true {
		t.Errorf("Expected key3=true, got %v", result.RawGetString("key3"))
	}
}
