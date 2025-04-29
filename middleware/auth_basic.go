package middleware

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/aryaadinulfadlan/go-social-api/helpers"
	"github.com/aryaadinulfadlan/go-social-api/internal/config"
)

func AuthBasicMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				helpers.UnauthorizedError(w, "Authorization header is missing")
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				helpers.UnauthorizedError(w, "Authorization header is malformed")
				return
			}
			bytes, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				helpers.UnauthorizedError(w, err.Error())
				return
			}
			username := config.Auth.Basic.User
			password := config.Auth.Basic.Pass

			creds := strings.SplitN(string(bytes), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				helpers.UnauthorizedError(w, "Invalid Credentials")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
