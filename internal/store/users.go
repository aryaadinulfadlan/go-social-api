package store

import (
	"context"
	"time"

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
	Comments  []Comment `json:"comments,omitempty"`
}

func (user *User) BeforeCreate(db *gorm.DB) (err error) {
	user.Id = uuid.New()
	return
}

type UserStore struct {
	db *gorm.DB
}

func (user_store *UserStore) CreateUser(ctx context.Context, user *User) error {
	err := user_store.db.WithContext(ctx).Create(&user).Error
	return err
}
func (user_store *UserStore) GetUser(ctx context.Context, userId uuid.UUID) (*User, error) {
	var user User
	err := user_store.db.WithContext(ctx).Take(&user, "id = ?", userId).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (user_store *UserStore) CheckUserExists(ctx context.Context, userId uuid.UUID) (*User, error) {
	var user User
	err := user_store.db.WithContext(ctx).Take(&user, "id = ?", userId).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
