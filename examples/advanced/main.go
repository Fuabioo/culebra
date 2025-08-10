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
			// Print the final configuration
			fmt.Printf("=== Final Configuration ===\n")

			// Environment-specific settings
			fmt.Printf("Environment: %s\n", viper.GetString("environment"))
			fmt.Printf("Debug enabled: %t\n", viper.GetBool("debug"))

			// Database configuration
			fmt.Printf("Database URL: %s\n", viper.GetString("database.url"))
			fmt.Printf("Connection pool size: %d\n", viper.GetInt("database.pool_size"))

			// Server configuration
			fmt.Printf("Server host: %s\n", viper.GetString("server.host"))
			fmt.Printf("Server port: %d\n", viper.GetInt("server.port"))

			// Features (computed based on environment)
			features := viper.Get("features")
			fmt.Printf("Enabled features: %v\n", features)

			// Service discovery
			services := viper.Get("services")
			fmt.Printf("Service endpoints: %v\n", services)

			// Rate limiting (environment-dependent)
			fmt.Printf("Rate limit: %d req/min\n", viper.GetInt("rate_limit"))

			// Monitoring endpoints (computed)
			monitoring := viper.Get("monitoring.endpoints")
			fmt.Printf("Monitoring endpoints: %v\n", monitoring)
		},
	}

	// Use Culebra for Lua configuration support
	culebra.UseWithCobra(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
