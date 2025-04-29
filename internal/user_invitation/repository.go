package userinvitation

import (
	"context"

	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUserInvitation(context.Context, *db.UserInvitation) error
	DeleteUserInvitation(context.Context, uuid.UUID) error
}

type RepositoryImplementation struct {
	gorm *gorm.DB
}

func NewRepository(gorm *gorm.DB) Repository {
	return &RepositoryImplementation{gorm: gorm}
}

func (repository *RepositoryImplementation) CreateUserInvitation(ctx context.Context, user_invitation *db.UserInvitation) error {
	err := repository.gorm.WithContext(ctx).Create(&user_invitation).Error
	if err != nil {
		return err
	}
	return nil
}

func (repository *RepositoryImplementation) DeleteUserInvitation(ctx context.Context, userId uuid.UUID) error {
	err := repository.gorm.WithContext(ctx).Where("user_id = ?", userId).Delete(&db.UserInvitation{}).Error
	if err != nil {
		return err
	}
	return nil
}
