package configuration

import "github.com/spf13/viper"

func apiConfigurations() {
	viper.SetDefault("PORT", "8888")
}
