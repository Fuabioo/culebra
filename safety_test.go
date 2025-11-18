package culebra

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// TestMaxDepthProtection tests that deeply nested structures are rejected
func TestMaxDepthProtection(t *testing.T) {
	// Create a Lua file with deeply nested tables (105 levels)
	content := "return "
	for i := 0; i < 105; i++ {
		content += "{"
	}
	content += "value = 1"
	for i := 0; i < 105; i++ {
		content += "}"
	}

	tmpfile, err := os.CreateTemp("", "deep-*.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Test with default depth limit (100)
	cfg := Config{
		FilePath: tmpfile.Name(),
	}
	_, err = Load(cfg)
	if err == nil {
		t.Error("Expected error for deeply nested structure, got nil")
	}
	if !strings.Contains(err.Error(), "maximum nesting depth exceeded") {
		t.Errorf("Expected depth limit error, got: %v", err)
	}
}

// TestMaxDepthCustomLimit tests custom depth limits
func TestMaxDepthCustomLimit(t *testing.T) {
	// Create a Lua file with 25 levels of nesting
	content := "return "
	for i := 0; i < 25; i++ {
		content += "{"
	}
	content += "value = 1"
	for i := 0; i < 25; i++ {
		content += "}"
	}

	tmpfile, err := os.CreateTemp("", "custom-depth-*.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Test with MaxDepth = 20 (should fail)
	cfg := Config{
		FilePath: tmpfile.Name(),
		MaxDepth: 20,
	}
	_, err = Load(cfg)
	if err == nil {
		t.Error("Expected error for structure exceeding custom depth limit")
	}

	// Test with MaxDepth = 30 (should succeed)
	cfg.MaxDepth = 30
	_, err = Load(cfg)
	if err != nil {
		t.Errorf("Expected success with MaxDepth=30, got error: %v", err)
	}
}

// TestMaxTableSizeProtection tests that oversized tables are rejected
func TestMaxTableSizeProtection(t *testing.T) {
	// Create a Lua file with a large table
	var content strings.Builder
	content.WriteString("return {\n")
	for i := 0; i < 1500; i++ {
		content.WriteString(fmt.Sprintf("  key_%d = %d,\n", i, i))
	}
	content.WriteString("}\n")

	tmpfile, err := os.CreateTemp("", "large-table-*.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content.String())); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Test with MaxTableSize = 1000 (should fail)
	cfg := Config{
		FilePath:     tmpfile.Name(),
		MaxTableSize: 1000,
	}
	_, err = Load(cfg)
	if err == nil {
		t.Error("Expected error for oversized table, got nil")
	}
	if !strings.Contains(err.Error(), "exceeds maximum allowed") {
		t.Errorf("Expected table size error, got: %v", err)
	}

	// Test with MaxTableSize = 2000 (should succeed)
	cfg.MaxTableSize = 2000
	_, err = Load(cfg)
	if err != nil {
		t.Errorf("Expected success with MaxTableSize=2000, got error: %v", err)
	}
}

// TestPanicProtection tests that panics are caught and converted to errors
func TestPanicProtection(t *testing.T) {
	// Test with a file that will cause Lua runtime error
	content := `
		error("Intentional Lua error for testing")
	`

	tmpfile, err := os.CreateTemp("", "panic-*.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	cfg := Config{FilePath: tmpfile.Name()}
	_, err = Load(cfg)

	// Should get an error, not a panic
	if err == nil {
		t.Error("Expected error from Lua error(), got nil")
	}
}

// TestDirectoryAsFilePath tests that directories are properly rejected
func TestDirectoryAsFilePath(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "culebra-test-dir-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	cfg := Config{FilePath: tmpdir}
	_, err = Load(cfg)

	if err == nil {
		t.Error("Expected error when FilePath is a directory, got nil")
	}
	if !strings.Contains(err.Error(), "directory") {
		t.Errorf("Expected 'directory' in error message, got: %v", err)
	}
}

// TestEmptyFilePath tests that empty filepath is rejected
func TestEmptyFilePath(t *testing.T) {
	cfg := Config{FilePath: ""}
	_, err := Load(cfg)

	if err == nil {
		t.Error("Expected error for empty FilePath, got nil")
	}
	if !strings.Contains(err.Error(), "empty") {
		t.Errorf("Expected 'empty' in error message, got: %v", err)
	}
}

// TestNonExistentFile tests that missing files are properly reported
func TestNonExistentFile(t *testing.T) {
	cfg := Config{FilePath: "/nonexistent/path/to/config.lua"}
	_, err := Load(cfg)

	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("Expected 'does not exist' in error message, got: %v", err)
	}
}

// TestGlobalVariableConflict tests behavior when Lua overwrites provided globals
func TestGlobalVariableConflict(t *testing.T) {
	content := `
		-- Override the global provided by Go
		myvar = "overridden by Lua"
		return {result = myvar}
	`

	tmpfile, err := os.CreateTemp("", "globals-*.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	cfg := Config{
		FilePath: tmpfile.Name(),
		Globals: map[string]any{
			"myvar": "from Go",
		},
	}

	result, err := Load(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Lua should be able to override the global
	resultMap := result["result"]
	if resultMap != "overridden by Lua" {
		t.Errorf("Expected Lua to override global, got: %v", resultMap)
	}
}

// TestEmptyTableBehavior tests how empty Lua tables are converted
func TestEmptyTableBehavior(t *testing.T) {
	content := `
		return {
			empty_table = {},
			nested_with_empty = {
				inner = {}
			}
		}
	`

	tmpfile, err := os.CreateTemp("", "empty-tables-*.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Test with ConvertArrays = false
	cfg := Config{
		FilePath:      tmpfile.Name(),
		ConvertArrays: false,
	}

	result, err := Load(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	emptyTable := result["empty_table"]
	if _, ok := emptyTable.(map[string]any); !ok {
		t.Errorf("Expected empty table to be map[string]any, got: %T", emptyTable)
	}

	// Test with ConvertArrays = true (empty tables should still be maps, not slices)
	cfg.ConvertArrays = true
	result, err = Load(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	emptyTable = result["empty_table"]
	if _, ok := emptyTable.(map[string]any); !ok {
		t.Errorf("Expected empty table to remain map[string]any even with ConvertArrays, got: %T", emptyTable)
	}
}

// TestConcurrentBindToViper tests concurrent access to BindToViper (should be documented as unsafe)
func TestConcurrentBindToViperSafety(t *testing.T) {
	// This test documents that concurrent BindToViper calls should NOT be done
	// It doesn't test concurrent access (which would be unsafe), just validates
	// that sequential calls work fine

	content := `return {test = "value"}`

	tmpfile, err := os.CreateTemp("", "concurrent-*.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	cfg := Config{FilePath: tmpfile.Name()}

	// Sequential calls should work fine
	_, err = Load(cfg)
	if err != nil {
		t.Errorf("First Load failed: %v", err)
	}

	_, err = Load(cfg)
	if err != nil {
		t.Errorf("Second Load failed: %v", err)
	}
}

// TestLuaArrayWithGaps tests Lua tables with non-sequential keys
func TestLuaArrayWithGaps(t *testing.T) {
	content := `
		return {
			array_with_gaps = {
				[1] = "first",
				[3] = "third",  -- Missing [2]
				[4] = "fourth"
			}
		}
	`

	tmpfile, err := os.CreateTemp("", "array-gaps-*.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	cfg := Config{
		FilePath:      tmpfile.Name(),
		ConvertArrays: true,
	}

	result, err := Load(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	arrayWithGaps := result["array_with_gaps"]

	// Should be a map, not a slice, because keys are not sequential
	if _, ok := arrayWithGaps.(map[string]any); !ok {
		t.Errorf("Expected array with gaps to be map[string]any, got: %T", arrayWithGaps)
	}
}

// TestVeryLargeFile tests performance and limits with large configs
func TestVeryLargeFile(t *testing.T) {
	var content strings.Builder
	content.WriteString("return {\n")

	// Create 5000 key-value pairs (reasonable large config)
	for i := 0; i < 5000; i++ {
		content.WriteString(fmt.Sprintf("  key_%d = %d,\n", i, i))
	}
	content.WriteString("}\n")

	tmpfile, err := os.CreateTemp("", "very-large-*.lua")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content.String())); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	cfg := Config{FilePath: tmpfile.Name()}

	// Should handle large files without issues (no MaxTableSize set)
	result, err := Load(cfg)
	if err != nil {
		t.Fatalf("Unexpected error loading large file: %v", err)
	}

	if len(result) != 5000 {
		t.Errorf("Expected 5000 entries, got: %d", len(result))
	}
}
