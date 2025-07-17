package culebra

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test.lua")

	configContent := `
app = {
    name = "Test App",
    version = "1.0.0"
}

debug_mode = true
port = 8080
hosts = {"localhost", "127.0.0.1"}
`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg := Config{FilePath: configFile}
	result, err := Load(cfg)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if result["debug_mode"] != true {
		t.Errorf("Expected debug_mode=true, got %v", result["debug_mode"])
	}

	if result["port"] != float64(8080) {
		t.Errorf("Expected port=8080, got %v", result["port"])
	}

	app, ok := result["app"].(map[string]any)
	if !ok {
		t.Fatalf("Expected app to be a map, got %T", result["app"])
	}

	if app["name"] != "Test App" {
		t.Errorf("Expected app.name='Test App', got %v", app["name"])
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	cfg := Config{FilePath: "/nonexistent/file.lua"}
	_, err := Load(cfg)
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestLoadEmptyFilePath(t *testing.T) {
	cfg := Config{}
	_, err := Load(cfg)
	if err == nil {
		t.Error("Expected error for empty file path")
	}
}

func TestLoadWithGlobals(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test.lua")

	configContent := `
result = global_value * 2
`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg := Config{
		FilePath: configFile,
		Globals: map[string]any{
			"global_value": 21,
		},
	}

	result, err := Load(cfg)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if result["result"] != float64(42) {
		t.Errorf("Expected result=42, got %v", result["result"])
	}
}

func TestLoadWithReturnStatement(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test.lua")

	configContent := `
local config = {
    app = {
        name = "Returned App",
        version = "2.0.0"
    },
    debug_mode = true,
    port = 9090
}

return config
`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg := Config{FilePath: configFile}
	result, err := Load(cfg)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if result["debug_mode"] != true {
		t.Errorf("Expected debug_mode=true, got %v", result["debug_mode"])
	}

	if result["port"] != float64(9090) {
		t.Errorf("Expected port=9090, got %v", result["port"])
	}

	app, ok := result["app"].(map[string]any)
	if !ok {
		t.Fatalf("Expected app to be a map, got %T", result["app"])
	}

	if app["name"] != "Returned App" {
		t.Errorf("Expected app.name='Returned App', got %v", app["name"])
	}
}
