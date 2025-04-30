package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/aryaadinulfadlan/go-social-api/helpers"
	"github.com/aryaadinulfadlan/go-social-api/internal/auth"
	"github.com/aryaadinulfadlan/go-social-api/internal/shared"
	"github.com/aryaadinulfadlan/go-social-api/internal/user"
	"github.com/google/uuid"
)

func AuthBearerMiddleware(authenticator auth.Authenticator, userRepository user.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				helpers.UnauthorizedError(w, "Authorization header is missing")
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				helpers.UnauthorizedError(w, "Authorization header is malformed")
				return
			}
			claims_token, claims_err := authenticator.ParseJWT(parts[1])
			if claims_err != nil {
				helpers.UnauthorizedError(w, claims_err.Error())
				return
			}
			userId, subject_err := claims_token.GetSubject()
			if subject_err != nil {
				helpers.UnauthorizedError(w, subject_err.Error())
				return
			}
			ctx := r.Context()
			userIdParsed, _ := uuid.Parse(userId)
			cachedUser, err := GetUserFromCache(ctx, userRepository, userIdParsed)
			if err != nil {
				helpers.UnauthorizedError(w, err.Error())
				return
			}
			ctx = context.WithValue(ctx, shared.UserCtx, cachedUser)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
