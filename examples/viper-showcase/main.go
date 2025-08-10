package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Fuabioo/culebra"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "viper-showcase",
		Short: "Comprehensive showcase of all Viper Get functions with Lua configuration",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("=== Comprehensive Viper Get Functions Showcase ===\n\n")
			
			// Basic Get Functions
			fmt.Printf("🔤 GetString Examples:\n")
			fmt.Printf("  App Name: '%s'\n", viper.GetString("app.name"))
			fmt.Printf("  Description: '%s'\n", viper.GetString("app.description"))
			fmt.Printf("  Empty String: '%s'\n", viper.GetString("empty_string"))
			fmt.Printf("  Non-existent: '%s'\n", viper.GetString("does.not.exist"))
			
			fmt.Printf("\n🔢 Integer Get Functions:\n")
			fmt.Printf("  Port (GetInt): %d\n", viper.GetInt("server.port"))
			fmt.Printf("  Port (GetInt32): %d\n", viper.GetInt32("server.port"))
			fmt.Printf("  Port (GetInt64): %d\n", viper.GetInt64("server.port"))
			fmt.Printf("  Port (GetUint): %d\n", viper.GetUint("server.port"))
			fmt.Printf("  Port (GetUint32): %d\n", viper.GetUint32("server.port"))
			fmt.Printf("  Port (GetUint64): %d\n", viper.GetUint64("server.port"))
			fmt.Printf("  Max Connections: %d\n", viper.GetInt("server.max_connections"))
			fmt.Printf("  Large Number (int64): %d\n", viper.GetInt64("numbers.large_int64"))
			
			fmt.Printf("\n✅ GetBool Examples:\n")
			fmt.Printf("  Debug Enabled: %t\n", viper.GetBool("debug.enabled"))
			fmt.Printf("  SSL Enabled: %t\n", viper.GetBool("ssl.enabled"))
			fmt.Printf("  Production Mode: %t\n", viper.GetBool("production_mode"))
			fmt.Printf("  Non-existent Bool: %t\n", viper.GetBool("does.not.exist"))
			
			fmt.Printf("\n📊 GetFloat64 Examples:\n")
			fmt.Printf("  Timeout: %.2f seconds\n", viper.GetFloat64("timeouts.request"))
			fmt.Printf("  Pi Value: %.5f\n", viper.GetFloat64("numbers.pi"))
			fmt.Printf("  Percentage: %.1f%%\n", viper.GetFloat64("metrics.success_rate"))
			fmt.Printf("  Integer as Float: %.1f\n", viper.GetFloat64("server.port"))
			
			fmt.Printf("\n📋 Slice Get Functions:\n")
			hosts := viper.GetStringSlice("database.hosts")
			fmt.Printf("  Database Hosts (%d): %v\n", len(hosts), hosts)
			
			ports := viper.GetIntSlice("allowed_ports")
			fmt.Printf("  Allowed Ports (%d): %v\n", len(ports), ports)
			
			tags := viper.GetStringSlice("tags")
			fmt.Printf("  Tags (%d): %v\n", len(tags), tags)
			
			fmt.Printf("\n⏰ Time & Duration Functions:\n")
			startTime := viper.GetTime("timestamps.startup")
			fmt.Printf("  Startup Time: %s\n", startTime.Format(time.RFC3339))
			
			cacheTTL := viper.GetDuration("cache.ttl")
			fmt.Printf("  Cache TTL: %v\n", cacheTTL)
			
			heartbeat := viper.GetDuration("monitoring.heartbeat")
			fmt.Printf("  Heartbeat Interval: %v\n", heartbeat)
			
			fmt.Printf("\n🗺️  Map Get Functions:\n")
			dbConfig := viper.GetStringMap("database.primary")
			fmt.Printf("  Database Config Keys: %v\n", getKeys(dbConfig))
			fmt.Printf("    Host: %v (type: %T)\n", dbConfig["host"], dbConfig["host"])
			fmt.Printf("    Port: %v (type: %T)\n", dbConfig["port"], dbConfig["port"])
			
			envVars := viper.GetStringMapString("environment")
			fmt.Printf("  Environment Variables (%d):\n", len(envVars))
			for key, value := range envVars {
				fmt.Printf("    %s=%s\n", key, value)
			}
			
			fmt.Printf("\n🌳 Nested Access Examples:\n")
			fmt.Printf("  Primary DB Host: %s\n", viper.GetString("database.primary.host"))
			fmt.Printf("  Primary DB Port: %d\n", viper.GetInt("database.primary.port"))
			fmt.Printf("  Primary DB SSL: %t\n", viper.GetBool("database.primary.ssl"))
			fmt.Printf("  Replica Count: %d\n", viper.GetInt("database.replicas.count"))
			fmt.Printf("  Log Level: %s\n", viper.GetString("logging.level"))
			fmt.Printf("  Log File: %s\n", viper.GetString("logging.file"))
			
			fmt.Printf("\n🎯 Generic Get Function:\n")
			appData := viper.Get("app")
			fmt.Printf("  App Data: %v (type: %T)\n", appData, appData)
			
			serverData := viper.Get("server")
			fmt.Printf("  Server Data: %v (type: %T)\n", serverData, serverData)
			
			numberData := viper.Get("numbers.pi")
			fmt.Printf("  Pi Value: %v (type: %T)\n", numberData, numberData)
			
			fmt.Printf("\n🔄 Type Conversion Examples:\n")
			fmt.Printf("  Port as String: '%s'\n", viper.GetString("server.port"))
			fmt.Printf("  Debug as String: '%s'\n", viper.GetString("debug.enabled"))
			fmt.Printf("  Float as Int: %d\n", viper.GetInt("timeouts.request"))
			fmt.Printf("  String Number as Int: %d\n", viper.GetInt("string_number"))
			fmt.Printf("  String Bool as Bool: %t\n", viper.GetBool("string_bool"))
			
			fmt.Printf("\n⚠️  Edge Cases & Error Handling:\n")
			fmt.Printf("  Zero Value: %d\n", viper.GetInt("numbers.zero"))
			fmt.Printf("  Negative Number: %d\n", viper.GetInt("numbers.negative"))
			fmt.Printf("  Empty Array Length: %d\n", len(viper.GetStringSlice("empty_array")))
			fmt.Printf("  Nil Value as String: '%s'\n", viper.GetString("nil_value"))
			fmt.Printf("  Non-existent Duration: %v\n", viper.GetDuration("does.not.exist"))
			
			fmt.Printf("\n🔍 Configuration Summary:\n")
			allSettings := viper.AllSettings()
			fmt.Printf("  Total Configuration Keys: %d\n", countKeys(allSettings))
			fmt.Printf("  Configuration File Used: %s\n", viper.ConfigFileUsed())
		},
	}

	// Use Culebra for Lua configuration support
	culebra.UseWithCobra(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// Helper function to get map keys
func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Helper function to count all keys recursively
func countKeys(m map[string]interface{}) int {
	count := 0
	for _, v := range m {
		count++
		if subMap, ok := v.(map[string]interface{}); ok {
			count += countKeys(subMap)
		}
	}
	return count
}