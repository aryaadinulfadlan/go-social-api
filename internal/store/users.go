package store

import (
	"context"
	"time"

	"github.com/aryaadinulfadlan/go-social-api/helpers"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name      string
	Username  string `gorm:"type:citext;unique"`
	Email     string `gorm:"type:citext;unique"`
	Password  string `gorm:"type:bytea"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Posts     []Post
}

type UserStore struct {
	db *gorm.DB
}

func (user_store *UserStore) Create(ctx context.Context, user *User) error {
	err := user_store.db.WithContext(ctx).Create(&user).Error
	helpers.PanicIfError(err)
	return nil
}
