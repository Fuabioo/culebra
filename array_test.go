package culebra

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func TestArrayConversion(t *testing.T) {
	// Create a temporary Lua config file with arrays
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.lua")
	luaConfig := `return {
    simple_array = {"one", "two", "three"},
    mixed_array = {1, "two", true},
    nested_structure = {
        array_field = {"a", "b", "c"},
        map_field = {
            key1 = "value1",
            key2 = "value2"
        }
    },
    pure_map = {
        key1 = "value1",
        key2 = "value2"
    }
}`

	if err := os.WriteFile(configPath, []byte(luaConfig), 0644); err != nil {
		t.Fatal(err)
	}

	// Test without array conversion
	t.Run("WithoutArrayConversion", func(t *testing.T) {
		data, err := Load(Config{FilePath: configPath, ConvertArrays: false})
		if err != nil {
			t.Fatal(err)
		}

		// Should be maps with numeric string keys
		simpleArray, ok := data["simple_array"].(map[string]any)
		if !ok {
			t.Errorf("expected simple_array to be map[string]any, got %T", data["simple_array"])
		}

		if simpleArray["1"] != "one" || simpleArray["2"] != "two" || simpleArray["3"] != "three" {
			t.Errorf("unexpected simple_array values: %v", simpleArray)
		}
	})

	// Test with array conversion
	t.Run("WithArrayConversion", func(t *testing.T) {
		data, err := Load(Config{FilePath: configPath, ConvertArrays: true})
		if err != nil {
			t.Fatal(err)
		}

		// Should be actual slices
		simpleArray, ok := data["simple_array"].([]any)
		if !ok {
			t.Errorf("expected simple_array to be []any, got %T", data["simple_array"])
		}

		expectedSimple := []any{"one", "two", "three"}
		if !reflect.DeepEqual(simpleArray, expectedSimple) {
			t.Errorf("expected %v, got %v", expectedSimple, simpleArray)
		}

		// Test mixed array
		mixedArray, ok := data["mixed_array"].([]any)
		if !ok {
			t.Errorf("expected mixed_array to be []any, got %T", data["mixed_array"])
		}

		expectedMixed := []any{float64(1), "two", true}
		if !reflect.DeepEqual(mixedArray, expectedMixed) {
			t.Errorf("expected %v, got %v", expectedMixed, mixedArray)
		}

		// Test nested structure
		nested, ok := data["nested_structure"].(map[string]any)
		if !ok {
			t.Errorf("expected nested_structure to be map[string]any, got %T", data["nested_structure"])
		}

		arrayField, ok := nested["array_field"].([]any)
		if !ok {
			t.Errorf("expected array_field to be []any, got %T", nested["array_field"])
		}

		expectedArray := []any{"a", "b", "c"}
		if !reflect.DeepEqual(arrayField, expectedArray) {
			t.Errorf("expected %v, got %v", expectedArray, arrayField)
		}

		mapField, ok := nested["map_field"].(map[string]any)
		if !ok {
			t.Errorf("expected map_field to be map[string]any, got %T", nested["map_field"])
		}

		if mapField["key1"] != "value1" || mapField["key2"] != "value2" {
			t.Errorf("unexpected map_field values: %v", mapField)
		}
	})

	// Test convenience function
	t.Run("LoadWithArrays", func(t *testing.T) {
		data, err := LoadWithArrays(configPath)
		if err != nil {
			t.Fatal(err)
		}

		simpleArray, ok := data["simple_array"].([]any)
		if !ok {
			t.Errorf("expected simple_array to be []any, got %T", data["simple_array"])
		}

		expectedSimple := []any{"one", "two", "three"}
		if !reflect.DeepEqual(simpleArray, expectedSimple) {
			t.Errorf("expected %v, got %v", expectedSimple, simpleArray)
		}
	})
}

func TestViperIntegration(t *testing.T) {
	// Create a temporary Lua config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.lua")
	luaConfig := `return {
    database = {
        host = "localhost",
        port = 5432,
        names = {"db1", "db2", "db3"}
    },
    features = {"feature1", "feature2"}
}`

	if err := os.WriteFile(configPath, []byte(luaConfig), 0644); err != nil {
		t.Fatal(err)
	}

	// Test AutoBindToViper
	t.Run("AutoBindToViper", func(t *testing.T) {
		v := viper.New()
		err := AutoBindToViper(Config{FilePath: configPath}, v)
		if err != nil {
			t.Fatal(err)
		}

		// Test that arrays are properly converted for Viper
		features := v.GetStringSlice("features")
		expectedFeatures := []string{"feature1", "feature2"}
		if !reflect.DeepEqual(features, expectedFeatures) {
			t.Errorf("expected %v, got %v", expectedFeatures, features)
		}

		dbNames := v.GetStringSlice("database.names")
		expectedDBNames := []string{"db1", "db2", "db3"}
		if !reflect.DeepEqual(dbNames, expectedDBNames) {
			t.Errorf("expected %v, got %v", expectedDBNames, dbNames)
		}

		// Test scalar values still work
		if v.GetString("database.host") != "localhost" {
			t.Errorf("expected localhost, got %s", v.GetString("database.host"))
		}

		if v.GetInt("database.port") != 5432 {
			t.Errorf("expected 5432, got %d", v.GetInt("database.port"))
		}
	})

	// Test BindToViperWithArrays convenience function
	t.Run("BindToViperWithArrays", func(t *testing.T) {
		v := viper.New()
		err := BindToViperWithArrays(configPath, v)
		if err != nil {
			t.Fatal(err)
		}

		features := v.GetStringSlice("features")
		expectedFeatures := []string{"feature1", "feature2"}
		if !reflect.DeepEqual(features, expectedFeatures) {
			t.Errorf("expected %v, got %v", expectedFeatures, features)
		}
	})
}

// Test struct unmarshaling with Viper
func TestStructUnmarshaling(t *testing.T) {
	type DatabaseConfig struct {
		Host  string   `mapstructure:"host"`
		Port  int      `mapstructure:"port"`
		Names []string `mapstructure:"names"`
	}

	type AppConfig struct {
		Database DatabaseConfig `mapstructure:"database"`
		Features []string       `mapstructure:"features"`
	}

	// Create a temporary Lua config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.lua")
	luaConfig := `return {
    database = {
        host = "localhost",
        port = 5432,
        names = {"db1", "db2", "db3"}
    },
    features = {"feature1", "feature2"}
}`

	if err := os.WriteFile(configPath, []byte(luaConfig), 0644); err != nil {
		t.Fatal(err)
	}

	v := viper.New()
	err := AutoBindToViper(Config{FilePath: configPath}, v)
	if err != nil {
		t.Fatal(err)
	}

	var config AppConfig
	err = v.Unmarshal(&config)
	if err != nil {
		t.Fatal(err)
	}

	// Verify the configuration was properly unmarshaled
	if config.Database.Host != "localhost" {
		t.Errorf("expected localhost, got %s", config.Database.Host)
	}

	if config.Database.Port != 5432 {
		t.Errorf("expected 5432, got %d", config.Database.Port)
	}

	expectedDBNames := []string{"db1", "db2", "db3"}
	if !reflect.DeepEqual(config.Database.Names, expectedDBNames) {
		t.Errorf("expected %v, got %v", expectedDBNames, config.Database.Names)
	}

	expectedFeatures := []string{"feature1", "feature2"}
	if !reflect.DeepEqual(config.Features, expectedFeatures) {
		t.Errorf("expected %v, got %v", expectedFeatures, config.Features)
	}
}
