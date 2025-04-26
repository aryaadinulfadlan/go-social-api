package main

import (
	"context"
	"encoding/base64"
	"net/http"
	"slices"
	"strings"

	"github.com/aryaadinulfadlan/go-social-api/internal/store"
	"github.com/google/uuid"
)

func (app *Application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.UnauthorizedError(w, "Authorization header is missing")
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.UnauthorizedError(w, "Authorization header is malformed")
			return
		}
		claims_token, claims_err := app.Authenticator.ParseJWT(parts[1])
		if claims_err != nil {
			app.UnauthorizedError(w, claims_err.Error())
			return
		}
		userId, subject_err := claims_token.GetSubject()
		if subject_err != nil {
			app.UnauthorizedError(w, subject_err.Error())
			return
		}
		ctx := r.Context()
		user, err := app.GetUserCache(ctx, userId)
		if err != nil {
			app.UnauthorizedError(w, err.Error())
			return
		}
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *Application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.UnauthorizedError(w, "Authorization header is missing")
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.UnauthorizedError(w, "Authorization header is malformed")
				return
			}
			bytes, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.UnauthorizedError(w, err.Error())
				return
			}
			username := app.Config.Auth.Basic.User
			password := app.Config.Auth.Basic.Pass

			creds := strings.SplitN(string(bytes), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				app.UnauthorizedError(w, "Invalid Credentials")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (app *Application) RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUserFromContext(r)
			if !app.CheckUserPermission(user.Id, permission, w, r) {
				app.ForbiddenError(w, "You do not have permission to access this resource.")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (app *Application) CheckUserPermission(userID uuid.UUID, permission string, w http.ResponseWriter, r *http.Request) bool {
	user := GetUserFromContext(r)
	permissions, err := app.Store.Permissions.GetPermissionNamesByRoleId(r.Context(), user.RoleId)
	if err != nil {
		return false
	}
	if slices.Contains(permissions, permission) {
		return true
	}
	return false
}

func (app *Application) GetUserCache(ctx context.Context, userId string) (*store.User, error) {
	user, err := app.CacheStorage.Users.Get(ctx, userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		user, err = app.Store.Users.GetExistingUser(ctx, "id", userId)
		if err != nil {
			return nil, err
		}
		err = app.CacheStorage.Users.Set(ctx, user)
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}
