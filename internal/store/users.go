package store

import (
	"context"
	"time"

	"github.com/aryaadinulfadlan/go-social-api/helpers"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name      string    `json:"name"`
	Username  string    `gorm:"type:citext;unique" json:"username"`
	Email     string    `gorm:"type:citext;unique" json:"email"`
	Password  string    `gorm:"type:bytea" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Posts     []Post    `json:"posts,omitempty"`
}

type UserStore struct {
	db *gorm.DB
}

func (user_store *UserStore) Create(ctx context.Context, user *User) error {
	err := user_store.db.WithContext(ctx).Create(&user).Error
	helpers.PanicIfError(err)
	return nil
}
