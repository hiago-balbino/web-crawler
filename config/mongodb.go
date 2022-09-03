package config

import "github.com/spf13/viper"

func mongoConfigurations() {
	viper.SetDefault("MONGODB_USERNAME", "")
	viper.SetDefault("MONGODB_PASSWORD", "")
	viper.SetDefault("MONGODB_DATABASE", "")
	viper.SetDefault("MONGODB_PORT", "")
	viper.SetDefault("MONGODB_HOST", "")
}
