package culebra

import (
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// UseWithCobra adds Lua config support to a Cobra command with automatic detection
func UseWithCobra(cmd *cobra.Command) {
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

// getViperConfigName safely attempts to get the config name from viper
// Returns empty string if unable to determine
func getViperConfigName() (result string) {
	// Try to use reflection with proper error handling
	defer func() {
		_ = recover() // Ignore panics from reflection
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

// getViperConfigPaths safely attempts to get the config paths from viper
// Returns empty slice if unable to determine
func getViperConfigPaths() (result []string) {
	// Try to use reflection with proper error handling
	defer func() {
		_ = recover() // Ignore panics from reflection
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

// AutoLoadLua automatically detects and loads .lua config files from Viper's config settings
func AutoLoadLua(cmd *cobra.Command) {
	cobra.OnInitialize(func() {
		// Check if Viper has a config file path configured
		if viperConfigFile := viper.ConfigFileUsed(); viperConfigFile != "" {
			tryLuaConfig(cmd, viperConfigFile)
			return
		}

		// Check Viper's config name and paths for .lua files
		configName := viper.GetString("config")
		if configName == "" {
			// Try common config names if none set
			for _, name := range []string{"config", cmd.Name()} {
				if tryLuaConfig(cmd, name) {
					return
				}
			}
		} else {
			tryLuaConfig(cmd, configName)
		}
	})
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
