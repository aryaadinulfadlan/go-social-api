package middleware

import (
	"context"

	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/aryaadinulfadlan/go-social-api/internal/redis"
	"github.com/aryaadinulfadlan/go-social-api/internal/user"
	"github.com/google/uuid"
)

func GetUserFromCache(ctx context.Context, repository user.Repository, userId uuid.UUID) (*db.User, error) {
	user, err := redis.GetUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		user, err = repository.GetById(ctx, userId)
		if err != nil {
			return nil, err
		}
		err = redis.SetUser(ctx, user)
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}
