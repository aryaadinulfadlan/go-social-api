package main

import (
	"log"

	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/aryaadinulfadlan/go-social-api/internal/env"
	"github.com/aryaadinulfadlan/go-social-api/internal/store"
)

const version = "1.0.0"

//	@title			Go Social API
//	@version		1.0
//	@description	API for Go Social Backend Technology
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	config := Config{
		Addr:   env.Envs.ADDR,
		apiURL: env.Envs.EXTERNAL_URL,
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
