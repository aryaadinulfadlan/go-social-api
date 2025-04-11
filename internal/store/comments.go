package store

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct {
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	UserId    uuid.UUID `gorm:"type:uuid" json:"user_id"`
	PostId    uuid.UUID `gorm:"type:uuid" json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *User     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
	Post      *Post     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"post,omitempty"`
}

func (comment *Comment) BeforeCreate(db *gorm.DB) (err error) {
	comment.Id = uuid.New()
	return
}

type CommentStore struct {
	db *gorm.DB
}
