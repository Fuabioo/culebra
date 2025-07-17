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
