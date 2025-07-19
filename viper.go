package culebra

import (
	"fmt"

	"github.com/spf13/viper"
)

func BindToViper(cfg Config, v *viper.Viper) error {
	data, err := Load(cfg)
	if err != nil {
		return fmt.Errorf("failed to load lua config: %w", err)
	}

	for key, value := range data {
		v.Set(key, value)
	}

	return nil
}

// BindToViperWithArrays loads a Lua config file with array conversion and binds to Viper
func BindToViperWithArrays(filePath string, v *viper.Viper) error {
	return BindToViper(Config{FilePath: filePath, ConvertArrays: true}, v)
}

// AutoBindToViper automatically determines the best configuration for Viper binding
func AutoBindToViper(cfg Config, v *viper.Viper) error {
	// For most use cases with Viper, we want array conversion enabled
	if !cfg.ConvertArrays {
		cfg.ConvertArrays = true
	}
	return BindToViper(cfg, v)
}
