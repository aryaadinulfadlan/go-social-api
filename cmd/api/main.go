package main

import (
	"log"

	"github.com/aryaadinulfadlan/go-social-api/internal/env"
)

func main() {
	config := Config{
		Addr: env.Envs.ADDR,
	}
	app := &Application{
		Config: config,
	}
	mux := app.Mount()
	log.Fatal(app.Run(mux))
}
