package redis

import (
	"context"
	"time"

	"github.com/aryaadinulfadlan/go-social-api/internal/config"
	"github.com/aryaadinulfadlan/go-social-api/internal/logger"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

const UserExpiredTime = time.Hour * 2

func Init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: config.Redis.Addr,
		DB:   config.Redis.DB,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			logger.Logger.Infof("New Redis connection established: %s", cn.String())
			return nil
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		logger.Logger.Fatal("failed connect to Redis: " + err.Error())
	}
}
