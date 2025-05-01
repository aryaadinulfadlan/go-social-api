package middleware

import (
	"fmt"
	"net/http"

	"github.com/aryaadinulfadlan/go-social-api/helpers"
	"github.com/aryaadinulfadlan/go-social-api/internal/config"
	"github.com/aryaadinulfadlan/go-social-api/internal/logger"
	ratelimiter "github.com/aryaadinulfadlan/go-social-api/internal/rate_limiter"
	"github.com/aryaadinulfadlan/go-social-api/internal/redis"
	"github.com/aryaadinulfadlan/go-social-api/internal/shared"
)

func RateLimiter() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.RateLimiter.Enabled {
				user := shared.GetUserFromContext(r)
				var key string
				if user != nil {
					key = fmt.Sprintf("rl:user:%s", user.Id)
				} else {
					ip := r.RemoteAddr
					key = fmt.Sprintf("rl:ip:%s", ip)
				}
				allowed, err := ratelimiter.AllowRequest(r.Context(), redis.RedisClient, key, config.RateLimiter.Max, config.RateLimiter.Duration)
				if err != nil {
					logger.Logger.Errorf("Rate limiter error: %v", err)
					helpers.InternalServerError(w, err.Error())
					return
				}
				if !allowed {
					helpers.RateLimitExceededResponse(w, r, "Rate limit exceeded.")
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
