package config

import "github.com/aryaadinulfadlan/go-social-api/internal/env"

type DBConfig struct {
	DATABASE_URL string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

func LoadDBConfig() *DBConfig {
	return &DBConfig{
		DATABASE_URL: env.Envs.DATABASE_URL,
		MaxOpenConns: env.Envs.DB_MAX_OPEN_CONNS,
		MaxIdleConns: env.Envs.DB_MAX_IDLE_CONNS,
		MaxIdleTime:  env.Envs.DB_MAX_IDLE_TIME,
	}
}
