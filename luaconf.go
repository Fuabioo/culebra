package culebra

import (
	"fmt"
	"os"

	"github.com/Fuabioo/culebra/internal"
	lua "github.com/yuin/gopher-lua"
)

type Config struct {
	FilePath      string
	Globals       map[string]any
	ConvertArrays bool // Convert Lua arrays to Go slices instead of maps
}

func Load(cfg Config) (map[string]any, error) {
	if cfg.FilePath == "" {
		return nil, fmt.Errorf("config file path is required")
	}

	if _, err := os.Stat(cfg.FilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", cfg.FilePath)
	}

	L := lua.NewState()
	defer L.Close()

	for key, value := range cfg.Globals {
		L.SetGlobal(key, internal.GoToLua(L, value))
	}

	if err := L.DoFile(cfg.FilePath); err != nil {
		return nil, fmt.Errorf("failed to execute lua config: %w", err)
	}

	// Check if the Lua script returned a table
	if L.GetTop() > 0 {
		returnValue := L.Get(-1)
		if table, ok := returnValue.(*lua.LTable); ok {
			result := make(map[string]any)
			table.ForEach(func(key, value lua.LValue) {
				result[key.String()] = internal.LuaToGoWithConfig(value, cfg.ConvertArrays)
			})
			return result, nil
		}
	}

	// Fallback to global variables (traditional style)
	result := make(map[string]any)
	globalTable := L.Get(lua.GlobalsIndex).(*lua.LTable)
	globalTable.ForEach(func(key, value lua.LValue) {
		if keyStr := key.String(); keyStr != "_G" && !isBuiltinGlobal(keyStr) {
			result[keyStr] = internal.LuaToGoWithConfig(value, cfg.ConvertArrays)
		}
	})

	return result, nil
}

// LoadWithArrays loads a Lua config file and converts arrays to Go slices
func LoadWithArrays(filePath string) (map[string]any, error) {
	return Load(Config{FilePath: filePath, ConvertArrays: true})
}

// LoadWithGlobals loads a Lua config file with predefined global variables
func LoadWithGlobals(filePath string, globals map[string]any) (map[string]any, error) {
	return Load(Config{FilePath: filePath, Globals: globals})
}

// LoadWithArraysAndGlobals loads a Lua config file with both array conversion and global variables
func LoadWithArraysAndGlobals(filePath string, globals map[string]any) (map[string]any, error) {
	return Load(Config{FilePath: filePath, Globals: globals, ConvertArrays: true})
}

func isBuiltinGlobal(key string) bool {
	builtins := map[string]bool{
		"_VERSION": true, "assert": true, "collectgarbage": true, "dofile": true,
		"error": true, "getfenv": true, "getmetatable": true, "ipairs": true,
		"load": true, "loadfile": true, "loadstring": true, "next": true,
		"pairs": true, "pcall": true, "print": true, "rawequal": true,
		"rawget": true, "rawset": true, "require": true, "select": true,
		"setfenv": true, "setmetatable": true, "tonumber": true, "tostring": true,
		"type": true, "unpack": true, "xpcall": true, "coroutine": true,
		"debug": true, "io": true, "math": true, "os": true, "package": true,
		"string": true, "table": true, "_GOPHER_LUA_VERSION": true,
		"_printregs": true, "channel": true, "module": true, "newproxy": true,
	}
	return builtins[key]
}
