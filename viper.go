package culebra

import (
	"fmt"

	"github.com/spf13/viper"
)

func BindToViper(cfg Config, v *viper.Viper) (err error) {
	// Panic protection
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("culebra.BindToViper: panic during binding: %v", r)
		}
	}()

	data, err := Load(cfg)
	if err != nil {
		return fmt.Errorf("culebra.BindToViper: %w", err)
	}

	for key, value := range data {
		v.Set(key, value)
	}

	return nil
}
