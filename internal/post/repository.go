package post

import (
	"context"

	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, post *db.Post) error
	GetDetail(ctx context.Context, postId uuid.UUID) (*db.Post, error)
	GetById(ctx context.Context, postId uuid.UUID) (*db.Post, error)
	Update(ctx context.Context, post *db.Post) (*db.Post, error)
	Delete(ctx context.Context, postId uuid.UUID) error
}

type RepositoryImplementation struct {
	gorm *gorm.DB
}

func NewRepository(gorm *gorm.DB) Repository {
	return &RepositoryImplementation{gorm: gorm}
}

func (repository *RepositoryImplementation) Create(ctx context.Context, post *db.Post) error {
	err := repository.gorm.WithContext(ctx).Create(&post).Error
	if err != nil {
		return err
	}
	return nil
}
func (repository *RepositoryImplementation) GetDetail(ctx context.Context, postId uuid.UUID) (*db.Post, error) {
	var post db.Post
	err := repository.gorm.WithContext(ctx).Preload("User").Preload("Comments").Take(&post, "id = ?", postId).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}
func (repository *RepositoryImplementation) GetById(ctx context.Context, postId uuid.UUID) (*db.Post, error) {
	var post db.Post
	err := repository.gorm.WithContext(ctx).Take(&post, "id = ?", postId).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}
func (repository *RepositoryImplementation) Update(ctx context.Context, post *db.Post) (*db.Post, error) {
	err := repository.gorm.WithContext(ctx).Model(&db.Post{}).
		Where("id = ?", post.Id).
		Updates(db.Post{
			Title:   post.Title,
			Content: post.Content,
			Tags:    post.Tags,
		}).Error
	if err != nil {
		return nil, err
	}
	return post, nil
}
func (repository *RepositoryImplementation) Delete(ctx context.Context, postId uuid.UUID) error {
	var post db.Post
	err := repository.gorm.WithContext(ctx).Where("id = ?", postId).Delete(&post).Error
	if err != nil {
		return err
	}
	return nil
}
