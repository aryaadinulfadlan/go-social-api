package permission

import (
	"context"

	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetPermissionNamesByRoleId(context.Context, uuid.UUID) ([]string, error)
}
type RepositoryImplementation struct {
	gorm *gorm.DB
}

func NewRepository(gorm *gorm.DB) Repository {
	return &RepositoryImplementation{gorm: gorm}
}

func (repository *RepositoryImplementation) GetPermissionNamesByRoleId(ctx context.Context, role_id uuid.UUID) ([]string, error) {
	var permission_names []string
	err := repository.gorm.WithContext(ctx).
		Model(&db.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", role_id).
		Pluck("permissions.name", &permission_names).Error
	if err != nil {
		return nil, err
	}
	return permission_names, nil
}
