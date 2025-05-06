package role

import (
	"context"

	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"gorm.io/gorm"
)

type Repository interface {
	GetRole(context.Context, string) (*db.Role, error)
}
type RepositoryImplementation struct {
	gorm *gorm.DB
}

func NewRepository(gorm *gorm.DB) Repository {
	return &RepositoryImplementation{gorm: gorm}
}

func (repository *RepositoryImplementation) GetRole(ctx context.Context, name string) (*db.Role, error) {
	var role db.Role
	err := repository.gorm.WithContext(ctx).Take(&role, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}
