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
		Use:   "example",
		Short: "Example CLI app using culebra",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("=== Basic Viper Functions Showcase ===\n\n")

			// String functions
			fmt.Printf("📝 String Access:\n")
			fmt.Printf("  App Name: %s\n", viper.GetString("app.name"))
			fmt.Printf("  App Version: %s\n", viper.GetString("app.version"))
			fmt.Printf("  Database Host: %s\n", viper.GetString("database.host"))

			// Integer functions
			fmt.Printf("\n🔢 Integer Access:\n")
			port := viper.GetInt("database.port")
			if port == 0 {
				log.Fatal("Database port not set or invalid")
			}
			fmt.Printf("  Database Port (GetInt): %d\n", port)
			fmt.Printf("  Database Port (GetInt32): %d\n", viper.GetInt32("database.port"))
			fmt.Printf("  Database Port (GetUint): %d\n", viper.GetUint("database.port"))

			// Boolean functions
			fmt.Printf("\n✅ Boolean Access:\n")
			fmt.Printf("  Debug Mode: %t\n", viper.GetBool("debug"))

			// Type conversion examples
			fmt.Printf("\n🔄 Type Conversions:\n")
			fmt.Printf("  Port as String: %s\n", viper.GetString("database.port"))
			fmt.Printf("  Debug as String: %s\n", viper.GetString("debug"))

			// Generic Get
			fmt.Printf("\n🎯 Generic Get:\n")
			appData := viper.Get("app")
			fmt.Printf("  App Data: %v (type: %T)\n", appData, appData)

			// Default values for missing keys
			fmt.Printf("\n⚙️  Default Values:\n")
			fmt.Printf("  Missing String: '%s'\n", viper.GetString("missing.key"))
			fmt.Printf("  Missing Int: %d\n", viper.GetInt("missing.key"))
			fmt.Printf("  Missing Bool: %t\n", viper.GetBool("missing.key"))
		},
	}

	// Integrate with Culebra for Lua configuration support
	culebra.EnableLuaConfig(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
