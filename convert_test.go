package culebra

import (
	"fmt"
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestLuaToGo(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Test table conversion without array support
	if err := L.DoString(`return {name = "test", value = 42}`); err != nil {
		t.Fatal(err)
	}

	lv := L.Get(-1)
	result := LuaToGo(lv)

	m, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("Expected map[string]any, got %T", result)
	}

	if m["name"] != "test" {
		t.Errorf("Expected name='test', got %v", m["name"])
	}

	if m["value"] != 42.0 {
		t.Errorf("Expected value=42.0, got %v", m["value"])
	}
}

func TestLuaToGoWithArrays(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Test array conversion
	if err := L.DoString(`return {1, 2, 3, 4, 5}`); err != nil {
		t.Fatal(err)
	}

	lv := L.Get(-1)
	result := LuaToGoWithArrays(lv)

	arr, ok := result.([]any)
	if !ok {
		t.Fatalf("Expected []any, got %T", result)
	}

	if len(arr) != 5 {
		t.Errorf("Expected 5 elements, got %d", len(arr))
	}

	for i, expected := range []float64{1.0, 2.0, 3.0, 4.0, 5.0} {
		if arr[i] != expected {
			t.Errorf("Expected arr[%d]=%v, got %v", i, expected, arr[i])
		}
	}
}

func TestLuaToGoSafe(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Test with depth limit
	nestedLua := "return "
	for i := 0; i < 60; i++ {
		nestedLua += "{"
	}
	nestedLua += "1"
	for i := 0; i < 60; i++ {
		nestedLua += "}"
	}

	if err := L.DoString(nestedLua); err != nil {
		t.Fatal(err)
	}

	lv := L.Get(-1)

	// Should fail with maxDepth = 50
	_, err := LuaToGoSafe(lv, false, 50, 0)
	if err == nil {
		t.Error("Expected depth limit error, got nil")
	}
	if !strings.Contains(err.Error(), "maximum nesting depth exceeded") {
		t.Errorf("Expected depth error, got: %v", err)
	}

	// Should succeed with maxDepth = 100
	_, err = LuaToGoSafe(lv, false, 100, 0)
	if err != nil {
		t.Errorf("Expected success with maxDepth=100, got error: %v", err)
	}

	// Should succeed with maxDepth = 0 (uses default 100)
	_, err = LuaToGoSafe(lv, false, 0, 0)
	if err != nil {
		t.Errorf("Expected success with default maxDepth, got error: %v", err)
	}
}

func TestLuaToGoSafeTableSize(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Create a table with 100 entries
	luaCode := "return {"
	for i := 1; i <= 100; i++ {
		if i > 1 {
			luaCode += ", "
		}
		luaCode += fmt.Sprintf("key_%d = %d", i, i)
	}
	luaCode += "}"

	if err := L.DoString(luaCode); err != nil {
		t.Fatal(err)
	}

	lv := L.Get(-1)

	// Should fail with maxTableSize = 50
	_, err := LuaToGoSafe(lv, false, 100, 50)
	if err == nil {
		t.Error("Expected table size error, got nil")
	}
	if !strings.Contains(err.Error(), "exceeds maximum allowed") {
		t.Errorf("Expected table size error, got: %v", err)
	}

	// Should succeed with maxTableSize = 150
	_, err = LuaToGoSafe(lv, false, 100, 150)
	if err != nil {
		t.Errorf("Expected success with maxTableSize=150, got error: %v", err)
	}

	// Should succeed with maxTableSize = 0 (unlimited)
	_, err = LuaToGoSafe(lv, false, 100, 0)
	if err != nil {
		t.Errorf("Expected success with unlimited table size, got error: %v", err)
	}
}

func TestGoToLua(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tests := []struct {
		name     string
		goValue  any
		expected string
	}{
		{"nil", nil, "nil"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"int", 42, "42"},
		{"int64", int64(100), "100"},
		{"float64", 3.14, "3.14"},
		{"string", "hello", "hello"},
		{"byte slice", []byte("test"), "test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lv := GoToLua(L, tt.goValue)
			actual := lv.String()
			if !strings.Contains(actual, tt.expected) {
				t.Errorf("Expected %s to contain %s, got %s", tt.name, tt.expected, actual)
			}
		})
	}
}

func TestGoToLuaMap(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	goMap := map[string]any{
		"name":  "test",
		"value": 42,
		"flag":  true,
	}

	lv := GoToLua(L, goMap)
	table, ok := lv.(*lua.LTable)
	if !ok {
		t.Fatalf("Expected *lua.LTable, got %T", lv)
	}

	// Verify values
	name := table.RawGetString("name")
	if name.String() != "test" {
		t.Errorf("Expected name='test', got %v", name)
	}

	value := table.RawGetString("value")
	if value.String() != "42" {
		t.Errorf("Expected value='42', got %v", value)
	}

	flag := table.RawGetString("flag")
	if flag.String() != "true" {
		t.Errorf("Expected flag='true', got %v", flag)
	}
}

func TestGoToLuaSlice(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	goSlice := []any{1, 2, 3, 4, 5}

	lv := GoToLua(L, goSlice)
	table, ok := lv.(*lua.LTable)
	if !ok {
		t.Fatalf("Expected *lua.LTable, got %T", lv)
	}

	// Lua arrays are 1-indexed
	for i := 1; i <= 5; i++ {
		val := table.RawGetInt(i)
		expected := float64(i)
		if val.String() != string(rune('0'+i)) {
			t.Errorf("Expected table[%d]=%v, got %v", i, expected, val)
		}
	}
}

func TestMustLuaToGo(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	if err := L.DoString(`return {test = "value"}`); err != nil {
		t.Fatal(err)
	}

	lv := L.Get(-1)

	// Should not panic for valid value
	result := MustLuaToGo(lv)
	m, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("Expected map[string]any, got %T", result)
	}

	if m["test"] != "value" {
		t.Errorf("Expected test='value', got %v", m["test"])
	}
}

func TestGoToLuaPanicProtection(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Test with various Go types that might cause issues
	testValues := []any{
		make(chan int),         // Unsupported type
		func() {},              // Function
		struct{ X int }{X: 42}, // Struct (not map[string]any)
	}

	for i, val := range testValues {
		t.Run("panic protection", func(t *testing.T) {
			// Should not panic, should convert to string
			lv := GoToLua(L, val)
			if lv == nil {
				t.Errorf("Test %d: Expected non-nil LValue for unsupported type", i)
			}
		})
	}
}

func TestLuaToGoAllTypes(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	luaCode := `
		return {
			nil_value = nil,
			bool_true = true,
			bool_false = false,
			number_int = 42,
			number_float = 3.14,
			string_val = "hello world",
			array = {1, 2, 3},
			map = {key = "value"},
			nested = {
				inner = {
					deep = "value"
				}
			}
		}
	`

	if err := L.DoString(luaCode); err != nil {
		t.Fatal(err)
	}

	lv := L.Get(-1)
	result := LuaToGoWithArrays(lv)

	m, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("Expected map[string]any, got %T", result)
	}

	// Check nil
	if m["nil_value"] != nil {
		t.Errorf("Expected nil_value to be nil, got %v", m["nil_value"])
	}

	// Check bools
	if m["bool_true"] != true {
		t.Errorf("Expected bool_true to be true, got %v", m["bool_true"])
	}
	if m["bool_false"] != false {
		t.Errorf("Expected bool_false to be false, got %v", m["bool_false"])
	}

	// Check numbers
	if m["number_int"] != 42.0 {
		t.Errorf("Expected number_int=42.0, got %v", m["number_int"])
	}
	if m["number_float"] != 3.14 {
		t.Errorf("Expected number_float=3.14, got %v", m["number_float"])
	}

	// Check string
	if m["string_val"] != "hello world" {
		t.Errorf("Expected string_val='hello world', got %v", m["string_val"])
	}

	// Check array
	arr, ok := m["array"].([]any)
	if !ok {
		t.Errorf("Expected array to be []any, got %T", m["array"])
	} else if len(arr) != 3 {
		t.Errorf("Expected array length 3, got %d", len(arr))
	}

	// Check nested
	nested, ok := m["nested"].(map[string]any)
	if !ok {
		t.Fatalf("Expected nested to be map[string]any, got %T", m["nested"])
	}
	inner, ok := nested["inner"].(map[string]any)
	if !ok {
		t.Fatalf("Expected nested.inner to be map[string]any, got %T", nested["inner"])
	}
	if inner["deep"] != "value" {
		t.Errorf("Expected nested.inner.deep='value', got %v", inner["deep"])
	}
}
