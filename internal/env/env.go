package env

import (
	"github.com/spf13/viper"
)

type Config struct {
	ADDR         string
	DATABASE_URL string
}

var Envs = GetEnv()

func GetEnv() Config {
	config := viper.New()
	config.SetConfigFile(".env")
	config.AddConfigPath("../../")
	err := config.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return Config{
		ADDR:         config.GetString("ADDR"),
		DATABASE_URL: config.GetString("DATABASE_URL"),
	}
}
