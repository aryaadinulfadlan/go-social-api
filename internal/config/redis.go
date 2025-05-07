package config

import "github.com/aryaadinulfadlan/go-social-api/internal/env"

type RedisConfig struct {
	Addr string
	DB   int
}

func LoadRedisConfig() *RedisConfig {
	return &RedisConfig{
		Addr: env.Envs.REDIS_ADDR,
		DB:   env.Envs.REDIS_DB,
	}
}
