package store

import (
	"context"

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
	}
	Users interface {
		CreateUser(context.Context, *User) error
		GetUser(context.Context, uuid.UUID) (*User, error)
		CheckUserExists(context.Context, uuid.UUID) (*User, error)
		FollowUser(context.Context, uuid.UUID, uuid.UUID) error
		UnfollowUser(context.Context, uuid.UUID, uuid.UUID) error
		GetUsersFollowing(context.Context, uuid.UUID) ([]*User, error)
		GetUsersFollowers(context.Context, uuid.UUID) ([]*User, error)
	}
	Comments interface {
		CreateComment(context.Context, *Comment) error
	}
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
	}
}
