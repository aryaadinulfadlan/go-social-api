package store

import (
	"context"
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

type PermissionStore struct {
	db *gorm.DB
}

func (permission_store *PermissionStore) GetPermissionNamesByRoleId(ctx context.Context, role_id uuid.UUID) ([]string, error) {
	var permission_names []string
	err := permission_store.db.WithContext(ctx).
		Model(&Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", role_id).
		Pluck("permissions.name", &permission_names).Error
	if err != nil {
		return nil, err
	}
	return permission_names, nil
}
