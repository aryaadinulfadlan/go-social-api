package test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ctx = context.Background()

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

func SetupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		logger.Logger.Fatalf("An error occured when opening Mock DB Connection: %s", err.Error())
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		logger.Logger.Fatalf("An error occured when opening Mock DB Connection: %s", err.Error())
	}
	return gormDB, mock
}
