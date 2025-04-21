package main

import (
	"encoding/base64"
	"net/http"
	"strings"
)

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
			username := app.Config.auth.basic.user
			password := app.Config.auth.basic.pass

			creds := strings.SplitN(string(bytes), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				app.UnauthorizedError(w, "Invalid Credentials")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
