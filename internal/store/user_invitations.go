package store

import (
	"context"
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

type UserInvitationStore struct {
	db *gorm.DB
}

func (user_invitation_store *UserInvitationStore) CreateUserInvitation(ctx context.Context, user_invitation *UserInvitation) error {
	err := user_invitation_store.db.WithContext(ctx).Create(&user_invitation).Error
	if err != nil {
		return err
	}
	return nil
}
func (user_invitation_store *UserInvitationStore) DeleteUserInvitation(ctx context.Context, userId uuid.UUID) error {
	err := user_invitation_store.db.WithContext(ctx).Where("user_id = ?", userId).Delete(&UserInvitation{}).Error
	if err != nil {
		return err
	}
	return nil
}
