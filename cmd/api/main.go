package main

import (
	"time"

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
			exp: time.Minute * 3,
		},
	}
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	db, err := db.OpenConnection(config.DB.DATABASE_URL, config.DB.MaxOpenConns, config.DB.MaxIdleConns, config.DB.MaxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	store := store.NewStorage(db)
	app := &Application{
		Config: config,
		Store:  *store,
		logger: logger,
	}
	mux := app.Mount()
	logger.Fatal(app.Run(mux))
}
