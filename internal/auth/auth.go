package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Authenticator interface {
	GenerateJWT(user_id string, exp time.Time) (string, error)
	ParseJWT(tokenStr string) (jwt.MapClaims, error)
}
