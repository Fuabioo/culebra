package culebra

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TestErrorConditions tests various error scenarios to ensure robust error handling
func TestErrorConditions(t *testing.T) {
	t.Run("LoadNonExistentFile", func(t *testing.T) {
		_, err := Load(Config{FilePath: "/nonexistent/path/config.lua"})
		if err == nil {
			t.Error("expected error for non-existent file")
		}
		// The error message should indicate file not found
		if err.Error() == "" {
			t.Error("expected non-empty error message")
		}
	})

	t.Run("LoadEmptyFilePath", func(t *testing.T) {
		_, err := Load(Config{FilePath: ""})
		if err == nil {
			t.Error("expected error for empty file path")
		}
	})

	t.Run("LoadDirectoryInsteadOfFile", func(t *testing.T) {
		tmpDir := t.TempDir()
		_, err := Load(Config{FilePath: tmpDir})
		if err == nil {
			t.Error("expected error when trying to load a directory")
		}
	})

	t.Run("LoadFileWithoutReadPermission", func(t *testing.T) {
		if os.Getuid() == 0 {
			t.Skip("Cannot test file permissions as root")
		}

		tmpDir := t.TempDir()
		restrictedFile := filepath.Join(tmpDir, "restricted.lua")

		// Create file and remove read permissions
		if err := os.WriteFile(restrictedFile, []byte("test = true"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(restrictedFile, 0000); err != nil {
			t.Fatal(err)
		}
		defer func() {
			_ = os.Chmod(restrictedFile, 0644) // Restore permissions for cleanup
		}()

		_, err := Load(Config{FilePath: restrictedFile})
		if err == nil {
			t.Error("expected error for file without read permission")
		}
	})

	t.Run("LoadInvalidLuaSyntax", func(t *testing.T) {
		tmpDir := t.TempDir()
		invalidFile := filepath.Join(tmpDir, "invalid.lua")
		invalidContent := `
			this is not valid lua
			function without_end(
			unclosed_table = {
		`

		if err := os.WriteFile(invalidFile, []byte(invalidContent), 0644); err != nil {
			t.Fatal(err)
		}

		_, err := Load(Config{FilePath: invalidFile})
		if err == nil {
			t.Error("expected error for invalid Lua syntax")
		}
	})

	t.Run("LoadLuaRuntimeError", func(t *testing.T) {
		tmpDir := t.TempDir()
		errorFile := filepath.Join(tmpDir, "runtime_error.lua")
		errorContent := `
			-- This will cause a runtime error
			error("intentional runtime error")
		`

		if err := os.WriteFile(errorFile, []byte(errorContent), 0644); err != nil {
			t.Fatal(err)
		}

		_, err := Load(Config{FilePath: errorFile})
		if err == nil {
			t.Error("expected error for Lua runtime error")
		}
	})

	t.Run("BindToViperWithInvalidFile", func(t *testing.T) {
		viper.Reset()

		cfg := Config{FilePath: "/nonexistent/file.lua"}
		err := BindToViper(cfg, viper.GetViper())
		if err == nil {
			t.Error("expected error when binding non-existent file to Viper")
		}
	})
}

// TestFileSystemEdgeCases tests edge cases related to file system operations
func TestFileSystemEdgeCases(t *testing.T) {
	t.Run("LoadEmptyFile", func(t *testing.T) {
		tmpDir := t.TempDir()
		emptyFile := filepath.Join(tmpDir, "empty.lua")

		if err := os.WriteFile(emptyFile, []byte(""), 0644); err != nil {
			t.Fatal(err)
		}

		data, err := Load(Config{FilePath: emptyFile})
		if err != nil {
			t.Fatalf("unexpected error for empty file: %v", err)
		}

		if len(data) != 0 {
			t.Errorf("expected empty data for empty file, got: %v", data)
		}
	})

	t.Run("LoadFileWithOnlyComments", func(t *testing.T) {
		tmpDir := t.TempDir()
		commentFile := filepath.Join(tmpDir, "comments.lua")
		commentContent := `
			-- This file only has comments
			-- No actual configuration
			--[[
				Multi-line comment
				with no code
			]]
		`

		if err := os.WriteFile(commentFile, []byte(commentContent), 0644); err != nil {
			t.Fatal(err)
		}

		data, err := Load(Config{FilePath: commentFile})
		if err != nil {
			t.Fatalf("unexpected error for comments-only file: %v", err)
		}

		if len(data) != 0 {
			t.Errorf("expected empty data for comments-only file, got: %v", data)
		}
	})

	t.Run("LoadVeryLargeFile", func(t *testing.T) {
		tmpDir := t.TempDir()
		largeFile := filepath.Join(tmpDir, "large.lua")

		// Create a large Lua configuration file
		var content string
		for i := 0; i < 1000; i++ {
			content += fmt.Sprintf("key_%d = \"value_%d\"\n", i, i)
		}

		if err := os.WriteFile(largeFile, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		data, err := Load(Config{FilePath: largeFile})
		if err != nil {
			t.Fatalf("unexpected error for large file: %v", err)
		}

		if len(data) != 1000 {
			t.Errorf("expected 1000 keys, got: %d", len(data))
		}
	})
}

// TestCobraErrorHandling tests error handling in Cobra integration
func TestCobraErrorHandling(t *testing.T) {
	t.Run("UseWithCobraInvalidConfig", func(t *testing.T) {
		tmpDir := t.TempDir()
		invalidFile := filepath.Join(tmpDir, "invalid.lua")
		invalidContent := `invalid lua syntax {{`

		if err := os.WriteFile(invalidFile, []byte(invalidContent), 0644); err != nil {
			t.Fatal(err)
		}

		viper.Reset()

		cmd := &cobra.Command{
			Use: "test",
			Run: func(cmd *cobra.Command, args []string) {
				// Command should not crash even with invalid config
			},
		}

		UseWithCobra(cmd)
		cmd.SetArgs([]string{"--config", invalidFile})

		// This should not panic, error should be handled gracefully
		err := cmd.Execute()
		// We don't expect Execute to fail, just the config loading to be handled internally
		_ = err
	})

	t.Run("TryLuaConfigWithPermissionError", func(t *testing.T) {
		if os.Getuid() == 0 {
			t.Skip("Cannot test file permissions as root")
		}

		tmpDir := t.TempDir()
		restrictedFile := filepath.Join(tmpDir, "restricted.lua")

		// Create file and remove read permissions
		if err := os.WriteFile(restrictedFile, []byte("test = true"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(restrictedFile, 0000); err != nil {
			t.Fatal(err)
		}
		defer func() {
			_ = os.Chmod(restrictedFile, 0644) // Restore permissions for cleanup
		}()

		viper.Reset()
		cmd := &cobra.Command{Use: "test"}

		result := tryLuaConfig(cmd, filepath.Join(tmpDir, "restricted"))
		if result {
			t.Error("expected false for file with no read permission")
		}
	})

	t.Run("LoadConfigWithNonExistentYamlFile", func(t *testing.T) {
		viper.Reset()
		cmd := &cobra.Command{Use: "test"}

		// Should handle non-existent YAML file gracefully
		loadConfig(cmd, "/nonexistent/file.yaml")

		// No assertion needed, just ensuring it doesn't panic
	})
}

// TestGlobalsErrorHandling tests error handling with global variables
func TestGlobalsErrorHandling(t *testing.T) {
	t.Run("LoadWithInvalidGlobalTypes", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "globals.lua")
		configContent := `
			-- Use the global variable
			result = complex_global
		`

		if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
			t.Fatal(err)
		}

		// Test with complex global types that should be handled safely
		globals := map[string]any{
			"complex_global": struct{ Name string }{Name: "test"}, // Struct type
		}

		data, err := Load(Config{FilePath: configFile, Globals: globals})
		if err != nil {
			t.Fatalf("unexpected error with complex globals: %v", err)
		}

		// Should convert unknown type to string representation
		result, ok := data["result"].(string)
		if !ok {
			t.Errorf("expected result to be string, got: %T", data["result"])
		}

		// Should contain some representation of the struct
		if result == "" {
			t.Error("expected non-empty string representation of struct")
		}
	})

	t.Run("LoadWithNilGlobals", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "nil_globals.lua")
		configContent := `
			-- Use globals that might be nil
			value1 = nil_global or "default1"
			value2 = another_nil or "default2"
		`

		if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
			t.Fatal(err)
		}

		globals := map[string]any{
			"nil_global":  nil,
			"another_nil": nil,
		}

		data, err := Load(Config{FilePath: configFile, Globals: globals})
		if err != nil {
			t.Fatalf("unexpected error with nil globals: %v", err)
		}

		// Should fall back to defaults
		if data["value1"] != "default1" || data["value2"] != "default2" {
			t.Errorf("unexpected values: %v", data)
		}
	})
}

// TestReflectionErrorHandling tests error handling in reflection code
func TestReflectionErrorHandling(t *testing.T) {
	t.Run("GetViperConfigNameReflectionFailure", func(t *testing.T) {
		// This tests the panic recovery in getViperConfigName
		viper.Reset()

		// Call the function - it should not panic even if reflection fails
		name := getViperConfigName()

		// The result might be empty due to reflection limitations, but it shouldn't crash
		_ = name
	})

	t.Run("GetViperConfigPathsReflectionFailure", func(t *testing.T) {
		// This tests the panic recovery in getViperConfigPaths
		viper.Reset()

		// Call the function - it should not panic even if reflection fails
		paths := getViperConfigPaths()

		// The result might be empty due to reflection limitations, but it shouldn't crash
		if paths == nil {
			t.Error("expected non-nil slice, even if empty")
		}
	})
}

// TestConcurrentAccess tests thread safety (basic test)
func TestConcurrentAccess(t *testing.T) {
	t.Run("ConcurrentLoad", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "concurrent.lua")
		configContent := `
			value = "concurrent_test"
			counter = 42
		`

		if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
			t.Fatal(err)
		}

		// Run multiple loads concurrently
		results := make(chan error, 5)

		for i := 0; i < 5; i++ {
			go func() {
				_, err := Load(Config{FilePath: configFile})
				results <- err
			}()
		}

		// Wait for all to complete
		for i := 0; i < 5; i++ {
			if err := <-results; err != nil {
				t.Errorf("concurrent load failed: %v", err)
			}
		}
	})
}
