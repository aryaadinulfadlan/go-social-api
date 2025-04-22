package store

import (
	"context"

	"github.com/aryaadinulfadlan/go-social-api/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Storage struct {
	Posts interface {
		CreatePost(context.Context, *Post) error
		GetPost(context.Context, uuid.UUID) (*Post, error)
		CheckPostExists(context.Context, uuid.UUID) (*Post, error)
		UpdatePost(context.Context, *Post) (*Post, error)
		DeletePost(context.Context, uuid.UUID) error
		GetPostFeed(context.Context, uuid.UUID, *model.PostParams) ([]*PostWithMetadata, int64, error)
	}
	Users interface {
		CreateUserAndInvite(context.Context, *User, *UserInvitation) error
		GetUser(context.Context, uuid.UUID) (*User, error)
		GetUserByInvitation(context.Context, string) (*User, error)
		GetExistingUser(context.Context, string, any) (*User, error)
		IsUserExists(context.Context, string, any) (int64, error)
		FollowUnfollowUser(context.Context, uuid.UUID, uuid.UUID) error
		GetConnections(context.Context, uuid.UUID, string) ([]*User, error)
		ActivateUser(context.Context, *User) (*User, error)
		DeleteUser(context.Context, uuid.UUID) error
	}
	Comments interface {
		CreateComment(context.Context, *Comment) error
	}
	UserInvitations interface {
		CreateUserInvitation(context.Context, *UserInvitation) error
		DeleteUserInvitation(context.Context, uuid.UUID) error
	}
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{
		Posts:           &PostStore{db},
		Users:           &UserStore{db},
		Comments:        &CommentStore{db},
		UserInvitations: &UserInvitationStore{db},
	}
}
