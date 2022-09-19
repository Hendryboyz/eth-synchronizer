package configs

import (
	"log"

	"github.com/spf13/viper"
)

var config *viper.Viper

func Init(environment string) {
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(environment)
	config.AddConfigPath("configs/")
	config.AddConfigPath("../configs/")
	if err := config.ReadInConfig(); err != nil {
		log.Fatal("Parse configuration fail")
	}
}

func GetConfig() *viper.Viper {
	return config
}
