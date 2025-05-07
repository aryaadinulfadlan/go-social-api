package env

import (
	"github.com/spf13/viper"
)

type EnvConfig struct {
	ADDR                  string
	DATABASE_URL          string
	DB_MAX_OPEN_CONNS     int
	DB_MAX_IDLE_CONNS     int
	DB_MAX_IDLE_TIME      string
	SECRET_KEY            string
	AUTH_BASIC_USERNAME   string
	AUTH_BASIC_PASSWORD   string
	REDIS_ADDR            string
	REDIS_DB              int
	RATE_LIMITER_MAX      int
	RATE_LIMITER_DURATION string
	RATE_LIMITER_ENABLED  bool
}

var Envs = GetEnv()

func GetEnv() EnvConfig {
	viper.AutomaticEnv()
	return EnvConfig{
		ADDR:                  viper.GetString("ADDR"),
		DATABASE_URL:          viper.GetString("DATABASE_URL"),
		DB_MAX_OPEN_CONNS:     viper.GetInt("DB_MAX_OPEN_CONNS"),
		DB_MAX_IDLE_CONNS:     viper.GetInt("DB_MAX_IDLE_CONNS"),
		DB_MAX_IDLE_TIME:      viper.GetString("DB_MAX_IDLE_TIME"),
		SECRET_KEY:            viper.GetString("SECRET_KEY"),
		AUTH_BASIC_USERNAME:   viper.GetString("AUTH_BASIC_USERNAME"),
		AUTH_BASIC_PASSWORD:   viper.GetString("AUTH_BASIC_PASSWORD"),
		REDIS_ADDR:            viper.GetString("REDIS_ADDR"),
		REDIS_DB:              viper.GetInt("REDIS_DB"),
		RATE_LIMITER_MAX:      viper.GetInt("RATE_LIMITER_MAX"),
		RATE_LIMITER_DURATION: viper.GetString("RATE_LIMITER_DURATION"),
		RATE_LIMITER_ENABLED:  viper.GetBool("RATE_LIMITER_ENABLED"),
	}
}
