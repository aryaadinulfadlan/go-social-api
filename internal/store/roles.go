package store

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	Id          uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Users       []User       `json:"users,omitempty"`
	Permissions []Permission `gorm:"many2many:role_permissions;joinForeignKey:role_id;joinReferences:permission_id" json:"permissions,omitempty"`
}

func (role *Role) BeforeCreate(db *gorm.DB) (err error) {
	role.Id = uuid.New()
	return
}

type RoleStore struct {
	db *gorm.DB
}

func (role_store *RoleStore) GetRole(ctx context.Context, name string) (*Role, error) {
	var role Role
	err := role_store.db.WithContext(ctx).Take(&role, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}
