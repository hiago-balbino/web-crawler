package config

import "github.com/spf13/viper"

func mongoConfigurations() {
	viper.SetDefault("MONGODB_DATABASE", "crawler")
	viper.SetDefault("MONGODB_COLLECTION", "page")
	viper.SetDefault("MONGODB_PORT", "27017")
	viper.SetDefault("MONGODB_HOST", "localhost")
}
