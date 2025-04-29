package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	Id          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Roles       []Role    `gorm:"many2many:role_permissions;joinForeignKey:permission_id;joinReferences:role_id" json:"roles,omitempty"`
}

func (permission *Permission) BeforeCreate(db *gorm.DB) (err error) {
	permission.Id = uuid.New()
	return
}
