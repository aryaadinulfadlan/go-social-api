package main

import (
	"time"

	"github.com/aryaadinulfadlan/go-social-api/internal/auth"
	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/aryaadinulfadlan/go-social-api/internal/env"
	ratelimiter "github.com/aryaadinulfadlan/go-social-api/internal/rate_limiter"
	"github.com/aryaadinulfadlan/go-social-api/internal/store"
	"github.com/aryaadinulfadlan/go-social-api/internal/store/cache"
	"github.com/sirupsen/logrus"
)

func main() {
	config := Config{
		Addr: env.Envs.ADDR,
		DB: DBConfig{
			DATABASE_URL: env.Envs.DATABASE_URL,
			MaxOpenConns: env.Envs.DB_MAX_OPEN_CONNS,
			MaxIdleConns: env.Envs.DB_MAX_IDLE_CONNS,
			MaxIdleTime:  env.Envs.DB_MAX_IDLE_TIME,
		},
		Mail: MailConfig{
			Exp: time.Hour * 48,
		},
		Auth: AuthConfig{
			Basic: AuthBasicConfig{
				User: env.Envs.AUTH_BASIC_USERNAME,
				Pass: env.Envs.AUTH_BASIC_PASSWORD,
			},
			TokenExp: time.Hour * 2,
		},
		Redis: RedisConfig{
			Addr: env.Envs.REDIS_ADDR,
		},
		RateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: env.Envs.RATE_LIMITER_REQUEST_COUNT,
			TimeFrame:            time.Second * 5,
			Enabled:              env.Envs.RATE_LIMITER_ENABLED,
		},
	}
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	db, err := db.OpenConnection(config.DB.DATABASE_URL, config.DB.MaxOpenConns, config.DB.MaxIdleConns, config.DB.MaxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	store := store.NewStorage(db)
	redisClient := cache.NewRedisClient(env.Envs.REDIS_ADDR)
	cacheStorage := cache.NewCacheStorage(redisClient)
	jwtAuthenticator := auth.NewJWTAuthenticator(env.Envs.SECRET_KEY)
	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		config.RateLimiter.RequestsPerTimeFrame,
		config.RateLimiter.TimeFrame,
	)
	app := &Application{
		Config:        config,
		Store:         *store,
		Logger:        logger,
		Authenticator: jwtAuthenticator,
		CacheStorage:  *cacheStorage,
		RateLimiter:   rateLimiter,
	}
	mux := app.Mount()
	logger.Fatal(app.Run(mux))
}
