package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aryaadinulfadlan/go-social-api/helpers"
	"github.com/aryaadinulfadlan/go-social-api/internal/auth"
	"github.com/aryaadinulfadlan/go-social-api/internal/comment"
	"github.com/aryaadinulfadlan/go-social-api/internal/permission"
	"github.com/aryaadinulfadlan/go-social-api/internal/post"
	"github.com/aryaadinulfadlan/go-social-api/internal/user"
	"github.com/aryaadinulfadlan/go-social-api/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(
	userHandler user.Handler,
	authenticator auth.Authenticator,
	userRepository user.Repository,
	permissionRepository permission.Repository,
	postHandler post.Handler,
	commentHandler comment.Handler,
) http.Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(3 * time.Second))
	r.Use(middleware.RateLimiter())
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		helpers.NotFoundError(w, fmt.Sprintf("Route %s %s is not exists", r.Method, r.URL.Path))
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		helpers.MethodNotAllowedError(w, fmt.Sprintf("%s %s is not valid", r.Method, r.URL.Path))
	})
	r.Route("/v1", func(r chi.Router) {
		r.With(middleware.AuthBasicMiddleware()).Get("/basic", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Authenticated as Basic Authentication"))
		})
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("PONG"))
		})
		r.Route("/posts", func(r chi.Router) {
			r.Use(middleware.AuthBearerMiddleware(authenticator, userRepository))
			r.With(middleware.RequirePermission(permissionRepository, "post:create")).Post("/", postHandler.Create)
			r.With(middleware.RequirePermission(permissionRepository, "post:detail")).Get("/{postId}", postHandler.GetDetail)
			r.With(middleware.RequirePermission(permissionRepository, "post:update")).Patch("/{postId}", postHandler.Update)
			r.With(middleware.RequirePermission(permissionRepository, "post:delete")).Delete("/{postId}", postHandler.Delete)
		})
		r.Route("/users", func(r chi.Router) {
			r.Use(middleware.AuthBearerMiddleware(authenticator, userRepository))
			r.With(middleware.RequirePermission(permissionRepository, "user:detail")).Get("/{userId}", userHandler.GetDetail)
			r.With(middleware.RequirePermission(permissionRepository, "user:follow")).Post("/{userId}/follow", userHandler.FollowUnfollow)
			r.With(middleware.RequirePermission(permissionRepository, "user:followers")).Get("/{userId}/followers", userHandler.GetConnections)
			r.With(middleware.RequirePermission(permissionRepository, "user:following")).Get("/{userId}/following", userHandler.GetConnections)
			r.With(middleware.RequirePermission(permissionRepository, "user:feed")).Get("/feed", userHandler.GetFeeds)
			r.With(middleware.RequirePermission(permissionRepository, "user:delete")).Delete("/{userId}", userHandler.Delete)
		})
		r.Route("/auth", func(r chi.Router) {
			r.Post("/sign-up", userHandler.CreateAndInvite)
			r.Post("/sign-in", userHandler.Login)
			r.Post("/resend-activation", userHandler.ResendActivation)
			r.Put("/activate/{token}", userHandler.Activate)
		})
		r.Route("/comments", func(r chi.Router) {
			r.Use(middleware.AuthBearerMiddleware(authenticator, userRepository))
			r.With(middleware.RequirePermission(permissionRepository, "comment:create")).Post("/{postId}", commentHandler.Create)
		})
	})
	return r
}
