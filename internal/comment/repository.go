package comment

import (
	"context"

	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, comment *db.Comment) error
}

type RepositoryImplementation struct {
	gorm *gorm.DB
}

func NewRepository(gorm *gorm.DB) Repository {
	return &RepositoryImplementation{gorm: gorm}
}

func (repository *RepositoryImplementation) Create(ctx context.Context, comment *db.Comment) error {
	err := repository.gorm.WithContext(ctx).Create(&comment).Error
	if err != nil {
		return err
	}
	return nil
}
