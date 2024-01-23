package config

import (
	"context"
	"log"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func Load(ctx context.Context) (*Application, error) {
	appConf := &Application{}
	err := viper.Unmarshal(appConf, func(dc *mapstructure.DecoderConfig) {
		// Prevent service to bootup if configuration is missing
		dc.ErrorUnset = true
		dc.ErrorUnused = false
	})

	if err != nil {
		log.Printf("failed to load config: %s", err)
		return nil, err
	}

	return appConf, nil
}
