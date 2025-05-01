package ratelimiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func AllowRequest(ctx context.Context, redisClient *redis.Client, key string, limit int, window time.Duration) (bool, error) {
	count, err := redisClient.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if count == 1 {
		_, err := redisClient.Expire(ctx, key, window).Result()
		if err != nil {
			return false, err
		}
	}
	return count <= int64(limit), nil
}
