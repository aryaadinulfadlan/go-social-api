package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserInvitation struct {
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	UserId    uuid.UUID `gorm:"type:uuid" json:"user_id"`
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *User     `gorm:"constraint:OnDelete:CASCADE;" json:"user,omitempty"`
}

func (user_invitation *UserInvitation) BeforeCreate(db *gorm.DB) (err error) {
	user_invitation.Id = uuid.New()
	return
}
