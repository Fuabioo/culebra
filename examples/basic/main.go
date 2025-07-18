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
			fmt.Printf("App Name: %s\n", viper.GetString("app.name"))
			fmt.Printf("App Version: %s\n", viper.GetString("app.version"))
			fmt.Printf("Database Host: %s\n", viper.GetString("database.host"))
			port := viper.GetInt("database.port")
			if port == 0 {
				log.Fatal("Database port not set or invalid")
			}
			fmt.Printf("Database Port: %d\n", port)
			fmt.Printf("Debug Mode: %t\n", viper.GetBool("debug"))
		},
	}

	// Option 1: Full integration with auto-detection AND explicit --config flag
	culebra.UseWithCobra(rootCmd)

	// Option 2: Just auto-detection without --config flag (uncomment to try)
	// culebra.AutoLoadLua(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
