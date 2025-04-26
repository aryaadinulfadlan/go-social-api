package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aryaadinulfadlan/go-social-api/internal/store"
	"github.com/redis/go-redis/v9"
)

type RedisUserStore struct {
	redisClient *redis.Client
}

func (redisUserStore *RedisUserStore) Get(ctx context.Context, userId string) (*store.User, error) {
	cacheKey := fmt.Sprintf("user-%s", userId)
	value, err := redisUserStore.redisClient.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var user store.User
	if value != "" {
		err := json.Unmarshal([]byte(value), &user)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func (redisUserStore *RedisUserStore) Set(ctx context.Context, user *store.User) error {
	cacheKey := fmt.Sprintf("user-%s", user.Id)
	json, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return redisUserStore.redisClient.SetEx(ctx, cacheKey, json, UserExpiredTime).Err()
}

func (redisUserStore *RedisUserStore) Delete(ctx context.Context, userId string) {
	cacheKey := fmt.Sprintf("user-%s", userId)
	redisUserStore.redisClient.Del(ctx, cacheKey)
}
