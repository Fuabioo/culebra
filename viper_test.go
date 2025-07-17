package culebra

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestBindToViper(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test.lua")

	configContent := `
database = {
    host = "localhost",
    port = 5432
}

debug_mode = true
`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	v := viper.New()
	cfg := Config{FilePath: configFile}

	err := BindToViper(cfg, v)
	if err != nil {
		t.Fatalf("BindToViper failed: %v", err)
	}

	if v.GetString("database.host") != "localhost" {
		t.Errorf("Expected database.host='localhost', got %v", v.GetString("database.host"))
	}

	if v.GetInt("database.port") != 5432 {
		t.Errorf("Expected database.port=5432, got %v", v.GetInt("database.port"))
	}

	if !v.GetBool("debug_mode") {
		t.Errorf("Expected debug_mode=true, got %v", v.GetBool("debug_mode"))
	}
}

func TestBindToViperWithInvalidFile(t *testing.T) {
	v := viper.New()
	cfg := Config{FilePath: "/nonexistent/file.lua"}

	err := BindToViper(cfg, v)
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}
