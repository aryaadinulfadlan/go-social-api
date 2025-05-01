package config

import (
	"time"

	"github.com/aryaadinulfadlan/go-social-api/internal/env"
)

type RateLimiterConfig struct {
	Max      int
	Duration time.Duration
	Enabled  bool
}

func LoadRateLimiterConfig() *RateLimiterConfig {
	timeout, err := time.ParseDuration(env.Envs.RATE_LIMITER_DURATION)
	if err != nil {
		timeout = time.Second
	}
	return &RateLimiterConfig{
		Max:      env.Envs.RATE_LIMITER_MAX,
		Duration: timeout,
		Enabled:  env.Envs.RATE_LIMITER_ENABLED,
	}
}
