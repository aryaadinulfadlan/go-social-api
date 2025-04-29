package config

import "github.com/aryaadinulfadlan/go-social-api/internal/env"

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	Protocol int
}

func LoadRedisConfig() *RedisConfig {
	return &RedisConfig{
		Addr:     env.Envs.REDIS_ADDR,
		Password: env.Envs.REDIS_PASSWORD,
		DB:       env.Envs.REDIS_DB,
		Protocol: env.Envs.REDIS_PROTOCOL,
	}
}
