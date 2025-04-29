package db

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id             uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	RoleId         uuid.UUID       `gorm:"type:uuid" json:"role_id"`
	Name           string          `json:"name"`
	Username       string          `gorm:"type:citext;unique" json:"username"`
	Email          string          `gorm:"type:citext;unique" json:"email"`
	Password       string          `json:"-"`
	IsActivated    bool            `json:"is_activated"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Posts          []Post          `json:"posts,omitempty"`
	UserInvitation *UserInvitation `json:"user_invitation,omitempty"`
	Role           *Role           `gorm:"constraint:OnDelete:CASCADE;" json:"role,omitempty"`
	Comments       []Comment       `json:"comments,omitempty"`
	Following      []*User         `gorm:"many2many:user_followers;joinForeignKey:follower_id;joinReferences:following_id" json:"following,omitempty"`
	Followers      []*User         `gorm:"many2many:user_followers;joinForeignKey:following_id;joinReferences:follower_id" json:"followers,omitempty"`
}
