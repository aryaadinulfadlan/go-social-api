package cache

import (
	"context"

	"github.com/aryaadinulfadlan/go-social-api/internal/store"
	"github.com/redis/go-redis/v9"
)

type CacheStorage struct {
	Users interface {
		Get(context.Context, string) (*store.User, error)
		Set(context.Context, *store.User) error
		Delete(context.Context, string)
	}
}

func NewCacheStorage(redisClient *redis.Client) *CacheStorage {
	return &CacheStorage{
		Users: &RedisUserStore{redisClient: redisClient},
	}
}
