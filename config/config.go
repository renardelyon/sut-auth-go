package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port     string `mapstructure:"PORT"`
	DBUrl    string `mapstructure:"DB_URL"`
	JWTKey   string `mapstructure:"JWT_SECRET_KEY"`
	AdminKey string `mapstructure:"ADMIN_KEY"`
}

func loadConfig() (config Config, err error) {
	viper.AddConfigPath("sut-auth-go/config/env/dev.env")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadConfig()
	if err != nil {
		log.Println("Error: ", err.Error())
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Println("Error: ", err.Error())
		return
	}

	return
}
