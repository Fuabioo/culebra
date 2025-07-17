package main

import (
	"fmt"
	"log"

	"culebra"
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
			fmt.Printf("Database Port: %d\n", viper.GetInt("database.port"))
			fmt.Printf("Debug Mode: %t\n", viper.GetBool("debug"))
		},
	}

	culebra.UseWithCobra(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
