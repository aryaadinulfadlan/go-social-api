package main

import (
	"log"

	"github.com/aryaadinulfadlan/go-social-api/internal/env"
	"github.com/aryaadinulfadlan/go-social-api/internal/store"
)

func main() {
	config := Config{
		Addr: env.Envs.ADDR,
	}
	store := store.NewStorage(nil)
	app := &Application{
		Config: config,
		Store:  *store,
	}
	mux := app.Mount()
	log.Fatal(app.Run(mux))
}
