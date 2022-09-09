package configuration

import "github.com/spf13/viper"

func loggerConfigurations() {
	viper.SetDefault("LOG_LEVEL", "")
}
