package test

import (
	"net/http"
	"time"

	"github.com/aryaadinulfadlan/go-social-api/internal/auth"
	"github.com/aryaadinulfadlan/go-social-api/internal/comment"
	"github.com/aryaadinulfadlan/go-social-api/internal/config"
	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/aryaadinulfadlan/go-social-api/internal/logger"
	"github.com/aryaadinulfadlan/go-social-api/internal/permission"
	"github.com/aryaadinulfadlan/go-social-api/internal/post"
	"github.com/aryaadinulfadlan/go-social-api/internal/redis"
	"github.com/aryaadinulfadlan/go-social-api/internal/role"
	"github.com/aryaadinulfadlan/go-social-api/internal/router"
	"github.com/aryaadinulfadlan/go-social-api/internal/user"
	userinvitation "github.com/aryaadinulfadlan/go-social-api/internal/user_invitation"
)

func SetupTest() http.Handler {
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
	return r
}

func GenerateJWT(userId string) (string, error) {
	authenticator := auth.NewJWTAuthenticator(config.SecretKey)
	exp := time.Now().Add(config.Auth.TokenExp).UTC()
	token, err := authenticator.GenerateJWT(userId, exp)
	if err != nil {
		return "", err
	}
	return token, nil
}
