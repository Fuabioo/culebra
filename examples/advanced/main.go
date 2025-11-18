package main

import (
	"fmt"
	"log"

	"github.com/Fuabioo/culebra"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "advanced-example",
		Short: "Advanced Lua configuration example showcasing Lua's programmability",
		Run: func(cmd *cobra.Command, args []string) {
			// Print comprehensive configuration showcasing all Viper Get functions
			fmt.Printf("=== Advanced Configuration Showcase ===\n\n")
			
			// String functions
			fmt.Printf("🔤 String Values:\n")
			fmt.Printf("  Environment: %s\n", viper.GetString("environment"))
			fmt.Printf("  App Name: %s\n", viper.GetString("app.name"))
			fmt.Printf("  App Version: %s\n", viper.GetString("app.version"))
			
			// Boolean functions
			fmt.Printf("\n✅ Boolean Values:\n")
			fmt.Printf("  Debug enabled: %t\n", viper.GetBool("debug"))
			fmt.Printf("  TLS enabled: %t\n", viper.GetBool("server.tls.enabled"))
			fmt.Printf("  Maintenance mode: %t\n", viper.GetBool("maintenance_mode"))
			
			// Integer functions (various types)
			fmt.Printf("\n🔢 Integer Values:\n")
			fmt.Printf("  Server port (GetInt): %d\n", viper.GetInt("server.port"))
			fmt.Printf("  Server workers (GetInt32): %d\n", viper.GetInt32("server.workers"))
			fmt.Printf("  Rate limit (GetInt64): %d\n", viper.GetInt64("rate_limit"))
			fmt.Printf("  DB pool size (GetUint): %d\n", viper.GetUint("database.pool_size"))
			fmt.Printf("  DB timeout (GetUint32): %d\n", viper.GetUint32("database.timeout"))
			
			// Float functions
			fmt.Printf("\n📊 Float Values:\n")
			fmt.Printf("  Cache default TTL: %.1f seconds\n", viper.GetFloat64("cache.default_ttl"))
			
			// Duration functions
			fmt.Printf("\n⏰ Duration Values:\n")
			fmt.Printf("  Session duration: %v\n", viper.GetDuration("security.session.duration"))
			fmt.Printf("  Monitoring interval: %v\n", viper.GetDuration("monitoring.interval"))
			
			// Slice functions
			fmt.Printf("\n📋 Array Values:\n")
			features := viper.GetStringSlice("features")
			fmt.Printf("  Features (%d): %v\n", len(features), features)
			
			monitoring := viper.GetStringSlice("monitoring.endpoints")
			fmt.Printf("  Monitoring endpoints (%d): %v\n", len(monitoring), monitoring)
			
			// Map functions  
			fmt.Printf("\n🗺️  Map Values:\n")
			services := viper.GetStringMap("services")
			fmt.Printf("  Services (%d configured):\n", len(services))
			for name := range services {
				fmt.Printf("    - %s: %s\n", name, viper.GetString("services."+name))
			}
			
			corsOrigins := viper.GetStringSlice("security.cors.origins")
			fmt.Printf("  CORS origins (%d): %v\n", len(corsOrigins), corsOrigins)
			
			// Nested access demonstration
			fmt.Printf("\n🌳 Nested Access:\n")
			fmt.Printf("  Database host: %s\n", viper.GetString("database.host"))
			fmt.Printf("  Cache type: %s\n", viper.GetString("cache.type"))
			fmt.Printf("  Redis host: %s\n", viper.GetString("cache.redis.host"))
			fmt.Printf("  Worker queue max jobs: %d\n", viper.GetInt("workers.queue.max_jobs"))
			
			// Generic Get function
			fmt.Printf("\n🎯 Raw Values (using Get):\n")
			appInfo := viper.Get("app")
			fmt.Printf("  App info (raw): %v (type: %T)\n", appInfo, appInfo)
			
			// Edge cases and type conversions
			fmt.Printf("\n⚠️  Edge Cases:\n")
			fmt.Printf("  Port as string: %s\n", viper.GetString("server.port"))
			fmt.Printf("  Debug as string: %s\n", viper.GetString("debug"))
			fmt.Printf("  Non-existent key: '%s'\n", viper.GetString("does.not.exist"))
			fmt.Printf("  Non-existent int: %d\n", viper.GetInt("does.not.exist"))
		},
	}

	// Use Culebra for Lua configuration support
	culebra.EnableLuaConfig(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
