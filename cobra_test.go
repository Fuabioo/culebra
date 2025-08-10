package culebra

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TestUseWithCobraBasic tests the basic UseWithCobra functionality with simpler assertions
func TestUseWithCobraBasic(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test.lua")
	configContent := `test_value = "from_lua"`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("ConfigFlagIntegration", func(t *testing.T) {
		viper.Reset()

		cmd := &cobra.Command{
			Use: "test",
			Run: func(cmd *cobra.Command, args []string) {
				// Just verify the command can execute without panicking
			},
		}

		UseWithCobra(cmd)
		cmd.SetArgs([]string{"--config", configFile})

		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("NoConfigProvided", func(t *testing.T) {
		viper.Reset()

		cmd := &cobra.Command{
			Use: "test",
			Run: func(cmd *cobra.Command, args []string) {},
		}

		UseWithCobra(cmd)

		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}
	})
}

// TestDirectFunctionCalls tests the individual functions directly
func TestDirectFunctionCalls(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("LoadConfigDirectly", func(t *testing.T) {
		viper.Reset()

		luaFile := filepath.Join(tmpDir, "direct.lua")
		luaContent := `key = "direct_value"`
		if err := os.WriteFile(luaFile, []byte(luaContent), 0644); err != nil {
			t.Fatal(err)
		}

		cmd := &cobra.Command{Use: "test"}
		loadConfig(cmd, luaFile)

		// Check that the config was loaded into viper
		if viper.GetString("key") != "direct_value" {
			t.Errorf("expected key='direct_value', got '%s'", viper.GetString("key"))
		}
	})

	t.Run("TryLuaConfigDirectly", func(t *testing.T) {
		viper.Reset()

		luaFile := filepath.Join(tmpDir, "tryconfig.lua")
		luaContent := `found = true`
		if err := os.WriteFile(luaFile, []byte(luaContent), 0644); err != nil {
			t.Fatal(err)
		}

		cmd := &cobra.Command{Use: "test"}
		result := tryLuaConfig(cmd, filepath.Join(tmpDir, "tryconfig"))

		if !result {
			t.Error("tryLuaConfig should have returned true")
		}

		if !viper.GetBool("found") {
			t.Error("config should have been loaded")
		}
	})

	t.Run("InitializeConfigDirectly", func(t *testing.T) {
		viper.Reset()

		luaFile := filepath.Join(tmpDir, "init.lua")
		luaContent := `initialized = true`
		if err := os.WriteFile(luaFile, []byte(luaContent), 0644); err != nil {
			t.Fatal(err)
		}

		cmd := &cobra.Command{Use: "test"}
		initializeConfig(cmd, luaFile)

		if !viper.GetBool("initialized") {
			t.Error("config should have been initialized")
		}
	})
}

// TestReflectionFunctions tests the reflection-based functions
func TestReflectionFunctions(t *testing.T) {
	t.Run("GetViperConfigNameSafety", func(t *testing.T) {
		viper.Reset()

		// Should not panic
		name := getViperConfigName()
		_ = name // Result doesn't matter, just that it doesn't crash
	})

	t.Run("GetViperConfigPathsSafety", func(t *testing.T) {
		viper.Reset()

		// Should not panic
		paths := getViperConfigPaths()
		if paths == nil {
			t.Error("should return non-nil slice")
		}
	})
}

// TestErrorHandling tests error conditions
func TestErrorHandling(t *testing.T) {
	t.Run("InvalidLuaFileHandling", func(t *testing.T) {
		tmpDir := t.TempDir()
		invalidFile := filepath.Join(tmpDir, "invalid.lua")
		invalidContent := `invalid lua syntax {{`

		if err := os.WriteFile(invalidFile, []byte(invalidContent), 0644); err != nil {
			t.Fatal(err)
		}

		viper.Reset()
		cmd := &cobra.Command{Use: "test", Run: func(cmd *cobra.Command, args []string) {}}

		// This should not panic, just handle the error gracefully
		loadConfig(cmd, invalidFile)
	})

	t.Run("NonExistentFileHandling", func(t *testing.T) {
		viper.Reset()
		cmd := &cobra.Command{Use: "test", Run: func(cmd *cobra.Command, args []string) {}}

		// Should handle non-existent file gracefully
		result := tryLuaConfig(cmd, "/nonexistent/file")
		if result {
			t.Error("should return false for non-existent file")
		}
	})
}
