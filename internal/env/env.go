package env

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ADDR              string
	DATABASE_URL      string
	DB_MAX_OPEN_CONNS int
	DB_MAX_IDLE_CONNS int
	DB_MAX_IDLE_TIME  string
	SECRET_KEY        string
}

var Envs = GetEnv()

func GetEnv() Config {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	config := viper.New()
	config.SetConfigFile(".env")
	config.AddConfigPath("../../")
	err := config.ReadInConfig()
	if err != nil {
		logger.Fatalln("Cannot load env file:", err)
	}
	return Config{
		ADDR:              config.GetString("ADDR"),
		DATABASE_URL:      config.GetString("DATABASE_URL"),
		DB_MAX_OPEN_CONNS: config.GetInt("DB_MAX_OPEN_CONNS"),
		DB_MAX_IDLE_CONNS: config.GetInt("DB_MAX_IDLE_CONNS"),
		DB_MAX_IDLE_TIME:  config.GetString("DB_MAX_IDLE_TIME"),
		SECRET_KEY:        config.GetString("SECRET_KEY"),
	}
}
