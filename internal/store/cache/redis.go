package cache

import (
	"time"

	"github.com/redis/go-redis/v9"
)

const UserExpiredTime = time.Hour * 2

func NewRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
		Protocol: 2,
	})
}
