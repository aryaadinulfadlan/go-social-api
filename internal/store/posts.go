package store

import (
	"context"
	"time"

	"github.com/aryaadinulfadlan/go-social-api/helpers"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	UserId    uuid.UUID `gorm:"type:uuid"`
	Title     string
	Content   string
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type PostStore struct {
	db *gorm.DB
}

func (post_store *PostStore) Create(ctx context.Context, post *Post) error {
	err := post_store.db.WithContext(ctx).Create(&post).Error
	helpers.PanicIfError(err)
	return nil
}
