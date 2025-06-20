package config

import (
	"github.com/spf13/viper"
)

type ConfigEnv struct {
	DBDriver   string `mapstructure:"DB_DRIVER"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBName     string `mapstructure:"DB_NAME"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBSSLMode  string `mapstructure:"DB_SSL_MODE"`
}

var Config ConfigEnv

func NewConfig() (config ConfigEnv, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app.env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err == nil {
		err = viper.Unmarshal(&config)
		return
	}

	err = viper.Unmarshal(&config)
	return
}
