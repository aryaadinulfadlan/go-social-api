package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aryaadinulfadlan/go-social-api/internal/auth"
	"github.com/aryaadinulfadlan/go-social-api/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
)

type DBConfig struct {
	DATABASE_URL string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}
type MailConfig struct {
	exp time.Duration
}
type AuthBasicConfig struct {
	user string
	pass string
}
type AuthConfig struct {
	basic    AuthBasicConfig
	tokenExp time.Duration
}
type Config struct {
	Addr string
	DB   DBConfig
	mail MailConfig
	auth AuthConfig
}
type Application struct {
	Config
	Store         store.Storage
	logger        *logrus.Logger
	authenticator auth.Authenticator
}

type userKey string

const userCtx userKey = "user"

func (app *Application) Mount() *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

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
		r.With(app.BasicAuthMiddleware()).Get("/test", app.Test)
		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.With(app.RequirePermission("post:create")).Post("/", app.CreatePostHandler)
			r.With(app.RequirePermission("post:detail")).Get("/{postId}", app.GetPostHandler)
			r.With(app.RequirePermission("post:update")).Patch("/{postId}", app.UpdatePostHandler)
			r.With(app.RequirePermission("post:delete")).Delete("/{postId}", app.DeletePostHandler)
		})
		r.Route("/users", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.With(app.RequirePermission("user:detail")).Get("/{userId}", app.GetUserHandler)
			r.With(app.RequirePermission("user:follow")).Post("/{userId}/follow", app.FollowUnfollowUserHandler)
			r.With(app.RequirePermission("user:followers")).Get("/{userId}/followers", app.GetConnectionsHandler)
			r.With(app.RequirePermission("user:following")).Get("/{userId}/following", app.GetConnectionsHandler)
			r.With(app.RequirePermission("user:feed")).Get("/feed", app.GetPostFeedHandler)
			r.With(app.RequirePermission("user:delete")).Delete("/{userId}", app.DeleteUserHandler)
		})
		r.Route("/auth", func(r chi.Router) {
			r.Post("/sign-up", app.CreateUserHandler)
			r.Post("/sign-in", app.LoginUserHandler)
			r.Post("/resend-activation", app.ResendActivationHandler)
			r.Put("/activate/{token}", app.ActivateUserHandler)
		})
		r.Route("/comments", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.With(app.RequirePermission("comment:create")).Post("/{postId}", app.CreateCommentHandler)
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
