package config

import (
	"time"

	"github.com/aryaadinulfadlan/go-social-api/internal/env"
)

type AuthBasicConfig struct {
	User string
	Pass string
}
type AuthConfig struct {
	Basic    AuthBasicConfig
	TokenExp time.Duration
}

func LoadAuthConfig() *AuthConfig {
	return &AuthConfig{
		Basic: AuthBasicConfig{
			User: env.Envs.AUTH_BASIC_USERNAME,
			Pass: env.Envs.AUTH_BASIC_PASSWORD,
		},
		TokenExp: time.Hour * 2,
	}
}
