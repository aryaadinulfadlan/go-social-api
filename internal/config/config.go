package config

import "github.com/aryaadinulfadlan/go-social-api/internal/env"

var (
	Addr        string
	SecretKey   string
	DB          *DBConfig
	Auth        *AuthConfig
	Redis       *RedisConfig
	RateLimiter *RateLimiterConfig
)

func Load() {
	Addr = env.Envs.ADDR
	SecretKey = env.Envs.SECRET_KEY
	DB = LoadDBConfig()
	Auth = LoadAuthConfig()
	Redis = LoadRedisConfig()
	RateLimiter = LoadRateLimiterConfig()
}
