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
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
		Protocol: config.Redis.Protocol,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		logger.Logger.Fatal("failed connect to Redis: " + err.Error())
	}
}
