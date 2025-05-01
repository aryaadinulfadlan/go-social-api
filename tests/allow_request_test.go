package tests

import (
	"context"
	"testing"
	"time"

	"github.com/aryaadinulfadlan/go-social-api/internal/config"
	ratelimiter "github.com/aryaadinulfadlan/go-social-api/internal/rate_limiter"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

type redisConfig struct {
	Addr     string
	Password string
	DB       int
	Protocol int
}

func NewRedisClient(config redisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
		Protocol: config.Protocol,
	})
}

var redisClient *redis.Client
var ctx = context.Background()

func setupRedis() {
	config.Load()
	cfg := redisConfig{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
		Protocol: config.Redis.Protocol,
	}
	redisClient = NewRedisClient(cfg)
	redisClient.FlushAll(ctx)
}
func TestConnection(t *testing.T) {
	setupRedis()
	result, err := redisClient.Ping(ctx).Result()
	assert.Nil(t, err)
	assert.Equal(t, "PONG", result)
}

func TestAllowRequest(t *testing.T) {
	setupRedis()
	key := "test_key"
	limit := 3
	window := time.Second
	for i := 1; i <= limit; i++ {
		allowed, err := ratelimiter.AllowRequest(ctx, redisClient, key, limit, window)
		assert.Nil(t, err)
		assert.True(t, allowed, "Request %d should be allowed", i)
	}
	allowed, err := ratelimiter.AllowRequest(ctx, redisClient, key, limit, window)
	assert.Nil(t, err)
	assert.False(t, allowed, "Request %d should be denied (over limit)", limit+1)
	time.Sleep(window)
	allowed, err = ratelimiter.AllowRequest(ctx, redisClient, key, limit, window)
	assert.Nil(t, err)
	assert.True(t, allowed, "Request after window should be allowed again")
}
