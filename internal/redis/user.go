package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/aryaadinulfadlan/go-social-api/internal/logger"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func GetUser(ctx context.Context, userId uuid.UUID) (*db.User, error) {
	cacheKey := fmt.Sprintf("user-%s", userId)
	value, err := RedisClient.Get(ctx, cacheKey).Result()
	if err == redis.Nil || value == "" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var user db.User
	if err := json.Unmarshal([]byte(value), &user); err != nil {
		logger.Logger.Warnf("Invalid user cache at key %s: %v", cacheKey, err)
		return nil, err
	}
	return &user, nil
}

func SetUser(ctx context.Context, user *db.User) error {
	cacheKey := fmt.Sprintf("user-%s", user.Id)
	json, err := json.Marshal(user)
	if err != nil {
		logger.Logger.Warnf("Invalid user cache at key %s: %v", cacheKey, err)
		return err
	}
	return RedisClient.SetEx(ctx, cacheKey, json, UserExpiredTime).Err()
}

func Delete(ctx context.Context, userId string) {
	cacheKey := fmt.Sprintf("user-%s", userId)
	RedisClient.Del(ctx, cacheKey)
}
