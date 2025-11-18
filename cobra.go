package culebra

import (
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// EnableLuaConfig adds seamless Lua config support to a Cobra command.
// It works like Viper's automatic config loading for .json/.yml files, but for .lua files.
// Adds a --config flag and automatically searches for .lua configs in Viper's search paths.
func EnableLuaConfig(cmd *cobra.Command) {
	var configFile string

	cmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (supports .lua, .yml, .json)")

	cobra.OnInitialize(func() {
		initializeConfig(cmd, configFile)
	})
}

// initializeConfig handles the configuration loading logic
func initializeConfig(cmd *cobra.Command, configFile string) {
	// If config file is explicitly provided, use it
	if configFile != "" {
		loadConfig(cmd, configFile)
		return
	}

	// Check if Viper has a config file path configured
	if viperConfigFile := viper.ConfigFileUsed(); viperConfigFile != "" {
		loadConfig(cmd, viperConfigFile)
		return
	}

	// Try autoloading Lua configuration
	tryAutoloadLua(cmd)
}

// tryAutoloadLua attempts to automatically load Lua configuration files
func tryAutoloadLua(cmd *cobra.Command) {
	configName := getViperConfigName()
	if configName == "" {
		return
	}

	configPaths := getViperConfigPaths()
	// If no paths are configured, default to current directory
	if len(configPaths) == 0 {
		configPaths = []string{"."}
	}

	// Try to find .lua version in all configured paths
	for _, path := range configPaths {
		luaFile := filepath.Join(path, configName+".lua")
		if tryLuaConfig(cmd, luaFile) {
			return
		}
	}
}

// getViperConfigName safely attempts to get the config name from viper using reflection.
// Returns empty string if unable to determine.
//
// WARNING: This function uses reflection to access Viper's internal 'configName' field.
// It may break if Viper changes its internal structure in future versions.
// The reflection is protected with panic recovery to fail gracefully.
func getViperConfigName() (result string) {
	// Panic protection - reflection into private fields can panic
	defer func() {
		if r := recover(); r != nil {
			// Silent failure - autoload will fall back gracefully
			result = ""
		}
	}()

	v := viper.GetViper()
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return ""
	}

	field := rv.FieldByName("configName")
	if field.IsValid() && field.Kind() == reflect.String && field.CanInterface() {
		return field.String()
	}
	return ""
}

// getViperConfigPaths safely attempts to get the config paths from viper using reflection.
// Returns empty slice if unable to determine.
//
// WARNING: This function uses reflection to access Viper's internal 'configPaths' field.
// It may break if Viper changes its internal structure in future versions.
// The reflection is protected with panic recovery to fail gracefully.
func getViperConfigPaths() (result []string) {
	// Panic protection - reflection into private fields can panic
	defer func() {
		if r := recover(); r != nil {
			// Silent failure - autoload will fall back to current directory
			result = []string{}
		}
	}()

	result = []string{}
	v := viper.GetViper()
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return result
	}

	field := rv.FieldByName("configPaths")
	if field.IsValid() && field.Kind() == reflect.Slice && field.CanInterface() {
		paths := make([]string, field.Len())
		for i := 0; i < field.Len(); i++ {
			if elem := field.Index(i); elem.Kind() == reflect.String {
				paths[i] = elem.String()
			}
		}
		return paths
	}
	return result
}

func loadConfig(cmd *cobra.Command, configFile string) {
	ext := strings.ToLower(filepath.Ext(configFile))

	if ext == ".lua" {
		cfg := Config{
			FilePath:      configFile,
			ConvertArrays: true, // Enable array conversion for better Viper integration
		}
		if err := BindToViper(cfg, viper.GetViper()); err != nil {
			cmd.PrintErrf("Error loading config file %s: %v\n", configFile, err)
		}
	} else {
		// Let Viper handle non-Lua config files
		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err != nil {
			cmd.PrintErrf("Error loading config file %s: %v\n", configFile, err)
		}
	}
}

func tryLuaConfig(cmd *cobra.Command, basePath string) bool {
	// Remove extension if present
	nameWithoutExt := strings.TrimSuffix(basePath, filepath.Ext(basePath))
	luaFile := nameWithoutExt + ".lua"

	cfg := Config{
		FilePath:      luaFile,
		ConvertArrays: true, // Enable array conversion for better Viper integration
	}
	if _, err := Load(cfg); err == nil {
		if err := BindToViper(cfg, viper.GetViper()); err != nil {
			cmd.PrintErrf("Error loading lua config file %s: %v\n", luaFile, err)
		}
		return true
	}
	return false
}
