package configuration

import "github.com/spf13/viper"

func apiConfigurations() {
	viper.SetDefault("API_PORT", "8888")
}
