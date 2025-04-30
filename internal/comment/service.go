package comment

import (
	"context"

	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/aryaadinulfadlan/go-social-api/internal/post"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, payload *CreateCommentPayload, postId uuid.UUID, userContext *db.User) (*db.Comment, error)
}

type ServiceImplementation struct {
	repository     Repository
	postRepository post.Repository
}

func NewService(repository Repository, postRepository post.Repository) Service {
	return &ServiceImplementation{repository: repository, postRepository: postRepository}
}

func (service *ServiceImplementation) Create(ctx context.Context, payload *CreateCommentPayload, postId uuid.UUID, userContext *db.User) (*db.Comment, error) {
	_, err := service.postRepository.GetById(ctx, postId)
	if err != nil {
		return nil, err
	}
	comment := &db.Comment{
		UserId:  userContext.Id,
		PostId:  postId,
		Content: payload.Content,
	}
	err = service.repository.Create(ctx, comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}
