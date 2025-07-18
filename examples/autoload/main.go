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
		Use:   "autoload-example",
		Short: "Example CLI app with automatic Lua config loading",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("App Name: %s\n", viper.GetString("app.name"))
			fmt.Printf("App Version: %s\n", viper.GetString("app.version"))
			fmt.Printf("Database Host: %s\n", viper.GetString("database.host"))
			fmt.Printf("Database Port: %d\n", viper.GetInt("database.port"))
			fmt.Printf("Debug Mode: %t\n", viper.GetBool("debug"))
			fmt.Printf("Environment: %s\n", viper.GetString("environment"))
		},
	}

	// Configure Viper for autoload - this triggers culebra's autoload mechanism
	// Setting config name enables autoload for .lua files
	viper.SetConfigName("example")
	// Adding config paths works with autoload - culebra will search these paths for .lua files
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("$HOME/.config")
	viper.AddConfigPath(".")

	// Enable Cobra integration (this handles --config flag and autoloading)
	culebra.UseWithCobra(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
