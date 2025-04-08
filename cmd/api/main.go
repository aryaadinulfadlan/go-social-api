package main

import (
	"log"

	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/aryaadinulfadlan/go-social-api/internal/env"
	"github.com/aryaadinulfadlan/go-social-api/internal/store"
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
	}
	db, err := db.OpenConnection(config.DB.DATABASE_URL, config.DB.MaxOpenConns, config.DB.MaxIdleConns, config.DB.MaxIdleTime)
	if err != nil {
		log.Fatal(err)
	}
	store := store.NewStorage(db)
	app := &Application{
		Config: config,
		Store:  *store,
	}
	mux := app.Mount()
	log.Fatal(app.Run(mux))
}
