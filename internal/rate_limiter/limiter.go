package ratelimiter

import (
	"context"
	"time"

	"github.com/aryaadinulfadlan/go-social-api/internal/redis"
)

func AllowRequest(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	count, err := redis.RedisClient.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if count == 1 {
		_, err := redis.RedisClient.Expire(ctx, key, window).Result()
		if err != nil {
			return false, err
		}
	}
	return count <= int64(limit), nil
}
