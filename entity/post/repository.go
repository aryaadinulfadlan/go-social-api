package post

import (
	"context"

	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	err := repository.gorm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existingPost db.Post
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", post.Id).
			Take(&existingPost).Error; err != nil {
			return err
		}
		if err := tx.Model(&db.Post{}).Where("id = ?", post.Id).
			Updates(db.Post{
				Title:   post.Title,
				Content: post.Content,
				Tags:    post.Tags,
			}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return post, nil
}
func (repository *RepositoryImplementation) Delete(ctx context.Context, postId uuid.UUID) error {
	err := repository.gorm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var post db.Post
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", postId).Take(&post).Error; err != nil {
			return err
		}
		if err := tx.Delete(&post).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
