package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	secret string
}

func NewJWTAuthenticator(secret string) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret,
	}
}

func (jwt_authenticator *JWTAuthenticator) GenerateJWT(user_id string, exp time.Time) (string, error) {
	now := time.Now().UTC()
	if exp.Before(now) {
		return "", fmt.Errorf("expiration time must be in the future")
	}
	claims := jwt.MapClaims{
		"sub": user_id,
		"iat": now.Unix(),
		"exp": exp.Unix(),
		"nbf": now.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token_string, err := token.SignedString([]byte(jwt_authenticator.secret))
	if err != nil {
		return "", err
	}
	return token_string, nil
}

func (jwt_authenticator *JWTAuthenticator) ParseJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwt_authenticator.secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}
