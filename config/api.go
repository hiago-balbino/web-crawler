package config

import "github.com/spf13/viper"

func apiConfigurations() {
	viper.SetDefault("API_PORT", "8888")
	viper.SetDefault("API_REQUEST_TIMEOUT", "60s")
}
