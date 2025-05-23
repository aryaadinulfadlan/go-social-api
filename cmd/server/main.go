package main

import (
	"net/http"
	"time"

	"github.com/aryaadinulfadlan/go-social-api/entity/comment"
	"github.com/aryaadinulfadlan/go-social-api/entity/permission"
	"github.com/aryaadinulfadlan/go-social-api/entity/post"
	"github.com/aryaadinulfadlan/go-social-api/entity/role"
	"github.com/aryaadinulfadlan/go-social-api/entity/user"
	userinvitation "github.com/aryaadinulfadlan/go-social-api/entity/user_invitation"
	"github.com/aryaadinulfadlan/go-social-api/internal/auth"
	"github.com/aryaadinulfadlan/go-social-api/internal/config"
	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/aryaadinulfadlan/go-social-api/internal/logger"
	"github.com/aryaadinulfadlan/go-social-api/internal/redis"
	"github.com/aryaadinulfadlan/go-social-api/internal/router"
)

func main() {
	config.Load()
	db.Init()
	logger.Init()
	redis.Init()
	authenticator := auth.NewJWTAuthenticator(config.SecretKey)

	userRepo := user.NewRepository(db.DB)
	roleRepo := role.NewRepository(db.DB)
	permissionRepo := permission.NewRepository(db.DB)
	userInvitationRepo := userinvitation.NewRepository(db.DB)
	userService := user.NewService(authenticator, userRepo, roleRepo, userInvitationRepo)
	userHandler := user.NewHandler(authenticator, userService)

	postRepo := post.NewRepository(db.DB)
	postService := post.NewService(postRepo)
	postHandler := post.NewHandler(postService)

	commentRepo := comment.NewRepository(db.DB)
	commentService := comment.NewService(commentRepo, postRepo)
	commentHandler := comment.NewHandler(commentService)

	r := router.NewRouter(
		userHandler,
		authenticator,
		userRepo,
		permissionRepo,
		postHandler,
		commentHandler,
	)
	server := &http.Server{
		Addr:         config.Addr,
		Handler:      r,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	logger.Logger.Infof("Server has started at %s", config.Addr)
	if err := server.ListenAndServe(); err != nil {
		logger.Logger.Fatalf("Server failed to start: %v", err)
	}
}
