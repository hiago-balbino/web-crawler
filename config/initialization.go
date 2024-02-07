package config

import (
	"github.com/spf13/viper"
)

func InitConfigurations() {
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	apiConfigurations()
	loggerConfigurations()
	mongoConfigurations()
}
