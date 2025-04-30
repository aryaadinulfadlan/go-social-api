package post

import (
	"context"

	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/aryaadinulfadlan/go-social-api/internal/shared"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, userId uuid.UUID, payload *CreatePostPayload) (*db.Post, error)
	GetDetail(ctx context.Context, postId uuid.UUID) (*db.Post, error)
	Update(ctx context.Context, postId uuid.UUID, userContext *db.User, payload *UpdatePostPayload) (*db.Post, error)
	Delete(ctx context.Context, postId uuid.UUID, userContext *db.User) error
}

type ServiceImplementation struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &ServiceImplementation{repository}
}

func (service *ServiceImplementation) Create(ctx context.Context, userId uuid.UUID, payload *CreatePostPayload) (*db.Post, error) {
	post := db.Post{
		UserId:  userId,
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
	}
	err := service.repository.Create(ctx, &post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}
func (service *ServiceImplementation) GetDetail(ctx context.Context, postId uuid.UUID) (*db.Post, error) {
	post, err := service.repository.GetDetail(ctx, postId)
	if err != nil {
		return nil, err
	}
	return post, nil
}
func (service *ServiceImplementation) Update(ctx context.Context, postId uuid.UUID, userContext *db.User, payload *UpdatePostPayload) (*db.Post, error) {
	post, err := service.repository.GetById(ctx, postId)
	if err != nil {
		return nil, err
	}
	if post.UserId != userContext.Id && userContext.Role.Name != "admin" {
		return nil, shared.ErrForbidden
	}
	post.Title = payload.Title
	post.Content = payload.Content
	post.Tags = payload.Tags
	post, err = service.repository.Update(ctx, post)
	if err != nil {
		return nil, err
	}
	return post, nil
}
func (service *ServiceImplementation) Delete(ctx context.Context, postId uuid.UUID, userContext *db.User) error {
	post, err := service.repository.GetById(ctx, postId)
	if err != nil {
		return err
	}
	if post.UserId != userContext.Id && userContext.Role.Name != "admin" {
		return shared.ErrForbidden
	}
	err = service.repository.Delete(ctx, post.Id)
	if err != nil {
		return err
	}
	return nil
}
