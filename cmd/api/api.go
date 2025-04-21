package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aryaadinulfadlan/go-social-api/internal/mailer"
	"github.com/aryaadinulfadlan/go-social-api/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

type DBConfig struct {
	DATABASE_URL string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}
type sendGridConfig struct {
	apiKey string
}
type mailTrapConfig struct {
	apiKey string
}
type MailConfig struct {
	exp       time.Duration
	fromEmail string
	sendGrid  sendGridConfig
	mailTrap  mailTrapConfig
}
type Config struct {
	Addr string
	DB   DBConfig
	mail MailConfig
}

type Application struct {
	Config
	Store  store.Storage
	logger *logrus.Logger
	mailer mailer.Client
}

func (app *Application) Mount() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(3 * time.Second))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.NotFoundError(w, fmt.Sprintf("Route %s %s is not exists", r.Method, r.URL.Path))
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		app.MethodNotAllowedError(w, fmt.Sprintf("%s %s is not valid", r.Method, r.URL.Path))
	})
	r.Route("/v1", func(r chi.Router) {
		r.Get("/test", app.Test)
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.CreatePostHandler)
			r.Get("/{postId}", app.GetPostHandler)
			r.Patch("/{postId}", app.UpdatePostHandler)
			r.Delete("/{postId}", app.DeletePostHandler)
		})
		r.Route("/users", func(r chi.Router) {
			r.Get("/{userId}", app.GetUserHandler)
			r.Post("/{userId}/follow", app.FollowUnfollowUserHandler)
			r.Get("/{userId}/followers", app.GetConnectionsHandler)
			r.Get("/{userId}/following", app.GetConnectionsHandler)
			r.Get("/feed", app.GetPostFeedHandler)
			r.Delete("/{userId}", app.DeleteUserHandler)
		})
		r.Route("/auth", func(r chi.Router) {
			r.Post("/", app.CreateUserHandler)
			r.Put("/activate/{token}", app.ActivateUserHandler)
		})
		r.Route("/comments", func(r chi.Router) {
			r.Post("/", app.CreateCommentHandler)
		})
	})
	return r
}
func (app *Application) Run(mux *chi.Mux) error {
	server := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	app.logger.Infof("Server has started at %s", app.Config.Addr)
	return server.ListenAndServe()
}
