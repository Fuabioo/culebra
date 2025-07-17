package culebra

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func UseWithCobra(cmd *cobra.Command) {
	var configFile string

	cmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (supports .lua)")

	cobra.OnInitialize(func() {
		if configFile != "" {
			cfg := Config{FilePath: configFile}
			if err := BindToViper(cfg, viper.GetViper()); err != nil {
				cmd.PrintErrf("Error loading config file: %v\n", err)
			}
		}
	})
}
