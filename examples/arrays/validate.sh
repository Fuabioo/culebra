#!/bin/bash

# Validation script for Culebra Array Conversion Example
# This ensures all features work correctly and backward compatibility is maintained

echo "🧪 Running Culebra Array Conversion Validation"
echo "=============================================="

# Test 1: Check if the example runs without errors
echo "📋 Test 1: Running the main example..."
if go run main.go > /dev/null 2>&1; then
    echo "✅ Main example runs successfully"
else
    echo "❌ Main example failed to run"
    exit 1
fi

# Test 2: Check if all configuration files can be loaded by Culebra
echo "📋 Test 2: Validating configuration files with Culebra..."

configs=("config-traditional.lua" "config-neovim-style.lua" "config-complex.lua")
for config in "${configs[@]}"; do
    cat > test_config_load.go << EOF
package main
import (
    "fmt"
    "os"
    "github.com/Fuabioo/culebra"
)
func main() {
    _, err := culebra.Load(culebra.Config{FilePath: "$config"})
    if err != nil {
        fmt.Printf("Error loading $config: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("$config loaded successfully")
}
EOF
    
    if go run test_config_load.go > /dev/null 2>&1; then
        echo "✅ $config loads correctly"
    else
        echo "❌ $config failed to load"
        rm -f test_config_load.go
        exit 1
    fi
    rm -f test_config_load.go
done

# Test 3: Check if backward compatibility works
echo "📋 Test 3: Testing backward compatibility..."
cat > test_backward_compat.go << 'EOF'
package main

import (
    "fmt"
    "os"
    "github.com/Fuabioo/culebra"
)

func main() {
    // Test old behavior
    data, err := culebra.Load(culebra.Config{
        FilePath: "config-traditional.lua",
        ConvertArrays: false,
    })
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
    
    // Arrays should be maps
    if hosts, ok := data["database_hosts"].(map[string]interface{}); ok {
        if len(hosts) != 3 {
            fmt.Printf("Expected 3 hosts, got %d\n", len(hosts))
            os.Exit(1)
        }
        fmt.Println("Backward compatibility: OK")
    } else {
        fmt.Println("Backward compatibility: FAILED - database_hosts is not a map")
        os.Exit(1)
    }
    
    // Test new behavior
    data2, err := culebra.LoadWithArrays("config-neovim-style.lua")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
    
    // Arrays should be slices
    if db, ok := data2["database"].(map[string]interface{}); ok {
        if replicas, ok := db["replicas"].([]interface{}); ok {
            if len(replicas) != 3 {
                fmt.Printf("Expected 3 replicas, got %d\n", len(replicas))
                os.Exit(1)
            }
            fmt.Println("Array conversion: OK")
        } else {
            fmt.Println("Array conversion: FAILED - replicas is not a slice")
            os.Exit(1)
        }
    }
}
EOF

if go run test_backward_compat.go; then
    echo "✅ Backward compatibility test passed"
else
    echo "❌ Backward compatibility test failed"
    exit 1
fi

# Cleanup
rm test_backward_compat.go

# Test 4: Check if Viper integration works
echo "📋 Test 4: Testing Viper integration..."
cat > test_viper.go << 'EOF'
package main

import (
    "fmt"
    "os"
    "github.com/Fuabioo/culebra"
    "github.com/spf13/viper"
)

func main() {
    v := viper.New()
    err := culebra.AutoBindToViper(culebra.Config{
        FilePath: "config-neovim-style.lua",
    }, v)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
    
    // Test string slice access
    replicas := v.GetStringSlice("database.replicas")
    if len(replicas) != 3 {
        fmt.Printf("Expected 3 replicas, got %d\n", len(replicas))
        os.Exit(1)
    }
    
    // Test nested access
    appName := v.GetString("app.name")
    if appName != "ModernApp" {
        fmt.Printf("Expected 'ModernApp', got '%s'\n", appName)
        os.Exit(1)
    }
    
    fmt.Println("Viper integration: OK")
}
EOF

if go run test_viper.go; then
    echo "✅ Viper integration test passed"
else
    echo "❌ Viper integration test failed"
    exit 1
fi

# Cleanup
rm test_viper.go

echo ""
echo "🎉 All validation tests passed!"
echo "✨ Culebra array conversion is working perfectly with full backward compatibility!"