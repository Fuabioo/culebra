package culebra

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadWithGlobalsFunction(t *testing.T) {
	// Create a temporary Lua config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "globals_test.lua")
	luaConfig := `
-- Use the provided global variable
app_name = global_app_name or "DefaultApp"
port = global_port or 8080
debug_enabled = global_debug
`

	if err := os.WriteFile(configPath, []byte(luaConfig), 0644); err != nil {
		t.Fatal(err)
	}

	// Test LoadWithGlobals function
	globals := map[string]any{
		"global_app_name": "CustomApp",
		"global_port":     9000,
		"global_debug":    true,
	}

	data, err := Load(Config{FilePath: configPath, Globals: globals})
	if err != nil {
		t.Fatalf("Failed to load config with globals: %v", err)
	}

	// Verify that globals were used
	if data["app_name"] != "CustomApp" {
		t.Errorf("expected app_name='CustomApp', got %v", data["app_name"])
	}

	if data["port"] != float64(9000) { // Lua numbers become float64
		t.Errorf("expected port=9000, got %v (type %T)", data["port"], data["port"])
	}

	if data["debug_enabled"] != true {
		t.Errorf("expected debug_enabled=true, got %v (type %T)", data["debug_enabled"], data["debug_enabled"])
		// Debug the actual global values received
		t.Logf("All data: %+v", data)
		t.Logf("Globals provided: %+v", globals)
	}
}

func TestLoadWithArraysAndGlobalsFunction(t *testing.T) {
	// Create a temporary Lua config file with arrays and globals
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "arrays_globals_test.lua")
	luaConfig := `
-- Use global variables and create arrays
prefix = global_prefix or "test"
items = {prefix .. "_item1", prefix .. "_item2", prefix .. "_item3"}
numbers = {1, 2, 3}
config = {
	name = prefix .. "_config",
	settings = {10, 20, 30}
}
`

	if err := os.WriteFile(configPath, []byte(luaConfig), 0644); err != nil {
		t.Fatal(err)
	}

	// Test LoadWithArraysAndGlobals function
	globals := map[string]any{
		"global_prefix": "custom",
	}

	data, err := Load(Config{FilePath: configPath, Globals: globals, ConvertArrays: true})
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check that arrays were converted to slices
	items, ok := data["items"].([]any)
	if !ok {
		t.Fatalf("expected items to be []any, got %T", data["items"])
	}

	expected := []any{"custom_item1", "custom_item2", "custom_item3"}
	if len(items) != 3 {
		t.Errorf("expected 3 items, got %d", len(items))
	}

	for i, item := range items {
		if item != expected[i] {
			t.Errorf("expected items[%d]='%v', got '%v'", i, expected[i], item)
		}
	}

	// Check numbers array
	numbers, ok := data["numbers"].([]any)
	if !ok {
		t.Fatalf("expected numbers to be []any, got %T", data["numbers"])
	}

	if len(numbers) != 3 || numbers[0] != float64(1) || numbers[1] != float64(2) || numbers[2] != float64(3) {
		t.Errorf("unexpected numbers array: %v", numbers)
	}

	// Check nested structure
	config, ok := data["config"].(map[string]any)
	if !ok {
		t.Fatalf("expected config to be map, got %T", data["config"])
	}

	if config["name"] != "custom_config" {
		t.Errorf("expected config.name='custom_config', got %v", config["name"])
	}

	settings, ok := config["settings"].([]any)
	if !ok {
		t.Fatalf("expected settings to be []any, got %T", config["settings"])
	}

	if len(settings) != 3 || settings[0] != float64(10) || settings[1] != float64(20) || settings[2] != float64(30) {
		t.Errorf("unexpected settings array: %v", settings)
	}
}

func TestLoadWithInvalidLuaSyntax(t *testing.T) {
	// Create a temporary Lua file with invalid syntax
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.lua")
	invalidLua := `
this is not valid lua syntax {{
missing = quotes
function without end
`

	if err := os.WriteFile(configPath, []byte(invalidLua), 0644); err != nil {
		t.Fatal(err)
	}

	// Should return an error for invalid syntax
	_, err := Load(Config{FilePath: configPath})
	if err == nil {
		t.Error("expected error for invalid Lua syntax, got nil")
	}

	// Test with Load using different config options
	_, err = Load(Config{FilePath: configPath, ConvertArrays: true})
	if err == nil {
		t.Error("Load with ConvertArrays should fail with invalid syntax")
	}

	_, err = Load(Config{FilePath: configPath, Globals: map[string]any{}})
	if err == nil {
		t.Error("Load with Globals should fail with invalid syntax")
	}

	_, err = Load(Config{FilePath: configPath, Globals: map[string]any{}, ConvertArrays: true})
	if err == nil {
		t.Error("Load with Globals and ConvertArrays should fail with invalid syntax")
	}
}

func TestLoadWithLuaRuntimeError(t *testing.T) {
	// Create a Lua file that will cause runtime error
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "runtime_error.lua")
	errorLua := `
-- This will cause a runtime error
-- Call a non-existent function
nonexistent_function()

value = "should not reach here"
`

	if err := os.WriteFile(configPath, []byte(errorLua), 0644); err != nil {
		t.Fatal(err)
	}

	// Should return an error for runtime error
	_, err := Load(Config{FilePath: configPath})
	if err == nil {
		t.Error("expected error for Lua runtime error, got nil")
	}
}

func TestLoadWithComplexGlobals(t *testing.T) {
	// Test with complex global types
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "complex_globals.lua")
	luaConfig := `
-- Use various global types
string_val = global_string
number_val = global_number
bool_val = global_bool
nil_val = global_nil
map_val = global_map
slice_val = global_slice
`

	if err := os.WriteFile(configPath, []byte(luaConfig), 0644); err != nil {
		t.Fatal(err)
	}

	globals := map[string]any{
		"global_string": "test_string",
		"global_number": 42.5,
		"global_bool":   true,
		"global_nil":    nil,
		"global_map":    map[string]any{"key": "value"},
		"global_slice":  []any{"a", "b", "c"},
	}

	data, err := Load(Config{FilePath: configPath, Globals: globals})
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify all global types were processed correctly
	if data["string_val"] != "test_string" {
		t.Error("String global not processed correctly")
	}

	if data["number_val"] != 42.5 {
		t.Error("Number global not processed correctly")
	}

	if data["bool_val"] != true {
		t.Error("Bool global not processed correctly")
	}

	if data["nil_val"] != nil {
		t.Error("Nil global not processed correctly")
	}

	// Map and slice globals should be converted back to Go types
	mapVal, ok := data["map_val"].(map[string]any)
	if !ok || mapVal["key"] != "value" {
		t.Errorf("Map global not processed correctly: %v (type %T)", data["map_val"], data["map_val"])
	}

	// Note: Without ConvertArrays=true, slices become maps with numeric string keys
	sliceVal := data["slice_val"]
	if sliceVal == nil {
		t.Error("Slice global should not be nil")
	}
}

func TestLoadWithReturnStatementAndArrays(t *testing.T) {
	// Test return statement with array conversion
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "return_arrays.lua")
	luaConfig := `
local config = {
	items = {"a", "b", "c"},
	numbers = {1, 2, 3},
	nested = {
		array = {"x", "y"},
		value = "test"
	}
}
return config`

	if err := os.WriteFile(configPath, []byte(luaConfig), 0644); err != nil {
		t.Fatal(err)
	}

	// Load with array conversion
	data, err := Load(Config{FilePath: configPath, ConvertArrays: true})
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check arrays were converted
	items, ok := data["items"].([]any)
	if !ok {
		t.Fatalf("expected items to be []any, got %T", data["items"])
	}

	if len(items) != 3 || items[0] != "a" || items[1] != "b" || items[2] != "c" {
		t.Errorf("unexpected items: %v", items)
	}

	// Check nested arrays
	nested, ok := data["nested"].(map[string]any)
	if !ok {
		t.Fatalf("expected nested to be map, got %T", data["nested"])
	}

	nestedArray, ok := nested["array"].([]any)
	if !ok {
		t.Fatalf("expected nested.array to be []any, got %T", nested["array"])
	}

	if len(nestedArray) != 2 || nestedArray[0] != "x" || nestedArray[1] != "y" {
		t.Errorf("unexpected nested array: %v", nestedArray)
	}
}

func TestLoadWithEmptyReturn(t *testing.T) {
	// Test return statement with empty table
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "empty_return.lua")
	luaConfig := `return {}`

	if err := os.WriteFile(configPath, []byte(luaConfig), 0644); err != nil {
		t.Fatal(err)
	}

	data, err := Load(Config{FilePath: configPath})
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(data) != 0 {
		t.Errorf("expected empty config, got %v", data)
	}
}

func TestLoadWithNonTableReturn(t *testing.T) {
	// Test return statement with non-table value
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "non_table_return.lua")
	luaConfig := `return "not a table"`

	if err := os.WriteFile(configPath, []byte(luaConfig), 0644); err != nil {
		t.Fatal(err)
	}

	// Should fall back to global variables (which would be empty)
	data, err := Load(Config{FilePath: configPath})
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Should be empty since no global variables are set
	if len(data) != 0 {
		t.Errorf("expected empty config for non-table return, got %v", data)
	}
}

func TestIsBuiltinGlobal(t *testing.T) {
	// Test the isBuiltinGlobal function with various built-in names
	builtins := []string{
		"_VERSION", "assert", "collectgarbage", "dofile", "error",
		"getfenv", "getmetatable", "ipairs", "load", "loadfile",
		"pairs", "pcall", "print", "require", "tostring", "type",
		"coroutine", "debug", "io", "math", "os", "package", "string", "table",
	}

	for _, builtin := range builtins {
		if !isBuiltinGlobal(builtin) {
			t.Errorf("expected %s to be recognized as builtin", builtin)
		}
	}

	// Test non-builtins
	nonBuiltins := []string{
		"myvar", "config", "app", "custom_function", "user_data",
	}

	for _, nonBuiltin := range nonBuiltins {
		if isBuiltinGlobal(nonBuiltin) {
			t.Errorf("expected %s to NOT be recognized as builtin", nonBuiltin)
		}
	}

	// Test edge cases
	if isBuiltinGlobal("") {
		t.Error("empty string should not be builtin")
	}

	if isBuiltinGlobal("_G") {
		t.Error("_G should be filtered out by Load function, not isBuiltinGlobal")
	}
}
