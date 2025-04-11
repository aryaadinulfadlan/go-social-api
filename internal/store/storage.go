package store

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(context.Context, uuid.UUID) (*Post, error)
	}
	Users interface {
		Create(context.Context, *User) error
	}
	Comments interface {
	}
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
	}
}
