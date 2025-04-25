package main

import (
	"time"

	"github.com/aryaadinulfadlan/go-social-api/internal/auth"
	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/aryaadinulfadlan/go-social-api/internal/env"
	"github.com/aryaadinulfadlan/go-social-api/internal/store"
	"github.com/sirupsen/logrus"
)

func main() {
	config := Config{
		Addr: env.Envs.ADDR,
		DB: DBConfig{
			DATABASE_URL: env.Envs.DATABASE_URL,
			MaxOpenConns: env.Envs.DB_MAX_OPEN_CONNS,
			MaxIdleConns: env.Envs.DB_MAX_IDLE_CONNS,
			MaxIdleTime:  env.Envs.DB_MAX_IDLE_TIME,
		},
		mail: MailConfig{
			exp: time.Hour * 48,
		},
		auth: AuthConfig{
			basic: AuthBasicConfig{
				user: env.Envs.AUTH_BASIC_USERNAME,
				pass: env.Envs.AUTH_BASIC_PASSWORD,
			},
			tokenExp: time.Hour * 2,
		},
	}
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	db, err := db.OpenConnection(config.DB.DATABASE_URL, config.DB.MaxOpenConns, config.DB.MaxIdleConns, config.DB.MaxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	store := store.NewStorage(db)
	jwtAuthenticator := auth.NewJWTAuthenticator(env.Envs.SECRET_KEY)
	app := &Application{
		Config:        config,
		Store:         *store,
		logger:        logger,
		authenticator: jwtAuthenticator,
	}
	mux := app.Mount()
	logger.Fatal(app.Run(mux))
}
