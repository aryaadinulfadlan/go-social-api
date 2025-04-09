package store

import (
	"context"
	"time"

	"github.com/aryaadinulfadlan/go-social-api/helpers"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Post struct {
	Id        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	UserId    uuid.UUID      `gorm:"type:uuid" json:"user_id"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Tags      pq.StringArray `gorm:"type:text[]" json:"tags"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	User      *User          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
}

func (post *Post) BeforeCreate(db *gorm.DB) (err error) {
	post.Id = uuid.New()
	return
}

type PostStore struct {
	db *gorm.DB
}

func (post_store *PostStore) Create(ctx context.Context, post *Post) error {
	err := post_store.db.WithContext(ctx).Create(&post).Error
	helpers.PanicIfError(err)
	return nil
}
