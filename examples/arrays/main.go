package main

import (
	"fmt"
	"log"
	"reflect"

	"github.com/Fuabioo/culebra"
	"github.com/spf13/viper"
)

// Example structs for demonstrating unmarshaling
type AppConfig struct {
	App      AppInfo     `mapstructure:"app"`
	Database Database    `mapstructure:"database"`
	Services []Service   `mapstructure:"services"`
	Logging  LogConfig   `mapstructure:"logging"`
}

type AppInfo struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
}

type Database struct {
	Primary        string    `mapstructure:"primary"`
	Replicas       []string  `mapstructure:"replicas"`
	ConnectionPool PoolConfig `mapstructure:"connection_pool"`
}

type PoolConfig struct {
	MinSize  int   `mapstructure:"min_size"`
	MaxSize  int   `mapstructure:"max_size"`
	Timeouts []int `mapstructure:"timeouts"`
}

type Service struct {
	Name      string   `mapstructure:"name"`
	Port      int      `mapstructure:"port"`
	Endpoints []string `mapstructure:"endpoints,omitempty"`
	Queues    []string `mapstructure:"queues,omitempty"`
	CronJobs  []string `mapstructure:"cron_jobs,omitempty"`
}

type LogConfig struct {
	Level      string          `mapstructure:"level"`
	Outputs    []string        `mapstructure:"outputs"`
	Formatters []LogFormatter  `mapstructure:"formatters"`
}

type LogFormatter struct {
	Name    string `mapstructure:"name"`
	Enabled bool   `mapstructure:"enabled"`
}

func main() {
	fmt.Println("🚀 Culebra Array Conversion Example")
	fmt.Println("====================================================")

	// Test 1: Backward Compatibility
	fmt.Println("\n1️⃣  Testing Backward Compatibility (without array conversion)")
	testBackwardCompatibility()

	// Test 2: Array Conversion
	fmt.Println("\n2️⃣  Testing Array Conversion (with ConvertArrays: true)")
	testArrayConversion()

	// Test 3: Viper Integration
	fmt.Println("\n3️⃣  Testing Viper Integration")
	testViperIntegration()

	// Test 4: Struct Unmarshaling
	fmt.Println("\n4️⃣  Testing Struct Unmarshaling")
	testStructUnmarshaling()

	// Test 5: Complex Configuration
	fmt.Println("\n5️⃣  Testing Complex Configuration")
	testComplexConfiguration()

	fmt.Println("\n✅ All tests completed successfully!")
	fmt.Println("🎯 Array conversion is working seamlessly with backward compatibility!")
}

func testBackwardCompatibility() {
	// Load without array conversion (traditional behavior)
	data, err := culebra.Load(culebra.Config{
		FilePath: "config-traditional.lua",
		ConvertArrays: false, // Explicit backward compatibility
	})
	if err != nil {
		log.Fatal("Failed to load traditional config:", err)
	}

	fmt.Printf("App Name: %s\n", data["app_name"])
	
	// Arrays should be maps with string keys
	if hosts, ok := data["database_hosts"].(map[string]interface{}); ok {
		fmt.Printf("Database Hosts (as map): %v\n", hosts)
		fmt.Printf("First host: %s\n", hosts["1"])
	}

	if features, ok := data["features"].(map[string]interface{}); ok {
		if auth, ok := features["authentication"].(map[string]interface{}); ok {
			fmt.Printf("Auth methods (as map): %v\n", auth)
		}
	}
}

func testArrayConversion() {
	// Load with array conversion (new behavior)
	data, err := culebra.LoadWithArrays("config-neovim-style.lua")
	if err != nil {
		log.Fatal("Failed to load neovim config:", err)
	}

	if app, ok := data["app"].(map[string]interface{}); ok {
		fmt.Printf("App Name: %s\n", app["name"])
	}

	if db, ok := data["database"].(map[string]interface{}); ok {
		if replicas, ok := db["replicas"].([]interface{}); ok {
			fmt.Printf("Database Replicas (as slice): %v\n", replicas)
			fmt.Printf("Number of replicas: %d\n", len(replicas))
		}
	}

	if services, ok := data["services"].([]interface{}); ok {
		fmt.Printf("Services (as slice): Found %d services\n", len(services))
		if len(services) > 0 {
			if service, ok := services[0].(map[string]interface{}); ok {
				fmt.Printf("First service: %s\n", service["name"])
			}
		}
	}
}

func testViperIntegration() {
	v := viper.New()
	
	// Use AutoBindToViper for seamless integration
	err := culebra.AutoBindToViper(culebra.Config{
		FilePath: "config-neovim-style.lua",
	}, v)
	if err != nil {
		log.Fatal("Failed to bind to viper:", err)
	}

	// Now we can use Viper's native array methods!
	replicas := v.GetStringSlice("database.replicas")
	fmt.Printf("Database replicas via Viper: %v\n", replicas)

	outputs := v.GetStringSlice("logging.outputs")
	fmt.Printf("Logging outputs via Viper: %v\n", outputs)

	timeouts := v.GetIntSlice("database.connection_pool.timeouts")
	fmt.Printf("Connection timeouts via Viper: %v\n", timeouts)

	// Test nested access
	appName := v.GetString("app.name")
	fmt.Printf("App name via Viper: %s\n", appName)
}

func testStructUnmarshaling() {
	v := viper.New()
	
	err := culebra.AutoBindToViper(culebra.Config{
		FilePath: "config-neovim-style.lua",
	}, v)
	if err != nil {
		log.Fatal("Failed to bind to viper:", err)
	}

	var config AppConfig
	err = v.Unmarshal(&config)
	if err != nil {
		log.Fatal("Failed to unmarshal config:", err)
	}

	fmt.Printf("Unmarshaled App: %s v%s (%s)\n", 
		config.App.Name, config.App.Version, config.App.Environment)
	
	fmt.Printf("Database Primary: %s\n", config.Database.Primary)
	fmt.Printf("Database Replicas: %v\n", config.Database.Replicas)
	fmt.Printf("Connection Pool Timeouts: %v\n", config.Database.ConnectionPool.Timeouts)
	
	fmt.Printf("Services (%d):\n", len(config.Services))
	for i, service := range config.Services {
		fmt.Printf("  %d. %s (port %d)\n", i+1, service.Name, service.Port)
		if len(service.Endpoints) > 0 {
			fmt.Printf("     Endpoints: %v\n", service.Endpoints)
		}
		if len(service.Queues) > 0 {
			fmt.Printf("     Queues: %v\n", service.Queues)
		}
		if len(service.CronJobs) > 0 {
			fmt.Printf("     Cron Jobs: %v\n", service.CronJobs)
		}
	}

	fmt.Printf("Logging Outputs: %v\n", config.Logging.Outputs)
	fmt.Printf("Log Formatters: %d configured\n", len(config.Logging.Formatters))
}

func testComplexConfiguration() {
	data, err := culebra.LoadWithArrays("config-complex.lua")
	if err != nil {
		log.Fatal("Failed to load complex config:", err)
	}

	// Test deeply nested arrays
	if microservices, ok := data["microservices"].([]interface{}); ok {
		fmt.Printf("Microservices: Found %d services\n", len(microservices))
		
		for i, ms := range microservices {
			if service, ok := ms.(map[string]interface{}); ok {
				name := service["name"]
				instances := service["instances"]
				
				if ports, ok := service["ports"].([]interface{}); ok {
					fmt.Printf("  %d. %s (%v instances) - Ports: %v\n", 
						i+1, name, instances, ports)
				}
				
				if deps, ok := service["dependencies"].([]interface{}); ok {
					fmt.Printf("     Dependencies: %v\n", deps)
				}
			}
		}
	}

	// Test nested structure arrays
	if infra, ok := data["infrastructure"].(map[string]interface{}); ok {
		if lbs, ok := infra["load_balancers"].([]interface{}); ok {
			fmt.Printf("Load Balancers: %d configured\n", len(lbs))
			
			if len(lbs) > 0 {
				if lb, ok := lbs[0].(map[string]interface{}); ok {
					if targets, ok := lb["targets"].([]interface{}); ok {
						fmt.Printf("  Targets: %v\n", targets)
					}
				}
			}
		}
	}

	// Test type verification
	fmt.Printf("\nType Verification:\n")
	printTypeInfo(data, "microservices")
	printTypeInfo(data, "infrastructure.databases")
	printTypeInfo(data, "monitoring.metrics.collectors")
}

func printTypeInfo(data map[string]interface{}, path string) {
	current := data
	
	// Simple path traversal for demonstration
	if path == "microservices" {
		if val, ok := current["microservices"]; ok {
			fmt.Printf("  %s: %s\n", path, reflect.TypeOf(val))
		}
	} else if path == "infrastructure.databases" {
		if infra, ok := current["infrastructure"].(map[string]interface{}); ok {
			if dbs, ok := infra["databases"]; ok {
				fmt.Printf("  %s: %s\n", path, reflect.TypeOf(dbs))
			}
		}
	} else if path == "monitoring.metrics.collectors" {
		if monitoring, ok := current["monitoring"].(map[string]interface{}); ok {
			if metrics, ok := monitoring["metrics"].(map[string]interface{}); ok {
				if collectors, ok := metrics["collectors"]; ok {
					fmt.Printf("  %s: %s\n", path, reflect.TypeOf(collectors))
				}
			}
		}
	}
}