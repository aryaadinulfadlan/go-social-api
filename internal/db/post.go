package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Post struct {
	Id        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	UserId    uuid.UUID      `gorm:"type:uuid" json:"user_id"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Tags      pq.StringArray `gorm:"type:varchar[]" json:"tags"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	User      *User          `gorm:"constraint:OnDelete:CASCADE;" json:"user,omitempty"`
	Comments  []Comment      `json:"comments,omitempty"`
}

type PostWithMetadata struct {
	Post
	Username      string `json:"username"`
	CommentsCount int64  `json:"comments_count"`
}

func (post *Post) BeforeCreate(db *gorm.DB) (err error) {
	post.Id = uuid.New()
	return
}
