package db

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	UserId    uuid.UUID `gorm:"type:uuid" json:"user_id"`
	PostId    uuid.UUID `gorm:"type:uuid" json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *User     `gorm:"constraint:OnDelete:CASCADE;" json:"user,omitempty"`
	Post      *Post     `gorm:"constraint:OnDelete:CASCADE;" json:"post,omitempty"`
}
