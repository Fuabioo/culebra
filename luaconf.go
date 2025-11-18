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
	MaxDepth      int  // Maximum nesting depth for tables (0 = use default of 100)
	MaxTableSize  int  // Maximum number of entries in a single table (0 = unlimited)
}

// luaBuiltinGlobals contains all Lua built-in global names
var luaBuiltinGlobals = map[string]bool{
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

func Load(cfg Config) (result map[string]any, err error) {
	// Panic protection - convert any panics to errors
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("culebra.Load: panic during config loading: %v", r)
		}
	}()

	if cfg.FilePath == "" {
		return nil, fmt.Errorf("culebra.Load: FilePath is empty (required field)")
	}

	// Check if path exists and is not a directory
	info, err := os.Stat(cfg.FilePath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("culebra.Load: file does not exist: %s", cfg.FilePath)
	}
	if err != nil {
		return nil, fmt.Errorf("culebra.Load: cannot access file %s: %w", cfg.FilePath, err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("culebra.Load: path is a directory, not a file: %s", cfg.FilePath)
	}

	// Set default max depth if not specified
	maxDepth := cfg.MaxDepth
	if maxDepth == 0 {
		maxDepth = 100
	}

	L := lua.NewState()
	defer L.Close()

	for key, value := range cfg.Globals {
		L.SetGlobal(key, internal.GoToLua(L, value))
	}

	if err := L.DoFile(cfg.FilePath); err != nil {
		return nil, fmt.Errorf("culebra.Load: failed to execute lua config: %w", err)
	}

	// Create conversion config for depth tracking
	convCfg := internal.ConversionConfig{
		ConvertArrays: cfg.ConvertArrays,
		MaxDepth:      maxDepth,
		MaxTableSize:  cfg.MaxTableSize,
	}

	// Check if the Lua script returned a table
	if L.GetTop() > 0 {
		returnValue := L.Get(-1)
		if table, ok := returnValue.(*lua.LTable); ok {
			// Check top-level table size if configured
			if convCfg.MaxTableSize > 0 {
				size := 0
				table.ForEach(func(_, _ lua.LValue) {
					size++
				})
				if size > convCfg.MaxTableSize {
					return nil, fmt.Errorf("culebra.Load: table size (%d) exceeds maximum allowed (%d)", size, convCfg.MaxTableSize)
				}
			}

			result := make(map[string]any)
			var convErr error
			table.ForEach(func(key, value lua.LValue) {
				if convErr != nil {
					return // Skip if we already hit an error
				}
				converted, convertErr := internal.LuaToGoSafe(value, convCfg, 0)
				if convertErr != nil {
					convErr = convertErr
					return
				}
				result[key.String()] = converted
			})
			if convErr != nil {
				return nil, fmt.Errorf("culebra.Load: %w", convErr)
			}
			return result, nil
		}
	}

	// Fallback to global variables (traditional style)
	result = make(map[string]any)
	globalTable := L.Get(lua.GlobalsIndex).(*lua.LTable)

	// Check global table size if configured (count only non-builtin globals)
	if convCfg.MaxTableSize > 0 {
		size := 0
		globalTable.ForEach(func(key, _ lua.LValue) {
			if keyStr := key.String(); keyStr != "_G" && !isBuiltinGlobal(keyStr) {
				size++
			}
		})
		if size > convCfg.MaxTableSize {
			return nil, fmt.Errorf("culebra.Load: global table size (%d) exceeds maximum allowed (%d)", size, convCfg.MaxTableSize)
		}
	}

	var convErr error
	globalTable.ForEach(func(key, value lua.LValue) {
		if convErr != nil {
			return // Skip if we already hit an error
		}
		if keyStr := key.String(); keyStr != "_G" && !isBuiltinGlobal(keyStr) {
			converted, convertErr := internal.LuaToGoSafe(value, convCfg, 0)
			if convertErr != nil {
				convErr = convertErr
				return
			}
			result[keyStr] = converted
		}
	})
	if convErr != nil {
		return nil, fmt.Errorf("culebra.Load: %w", convErr)
	}

	return result, nil
}

func isBuiltinGlobal(key string) bool {
	return luaBuiltinGlobals[key]
}
