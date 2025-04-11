package model

import "github.com/google/uuid"

type CreateCommentPayload struct {
	UserId  uuid.UUID `json:"user_id" validate:"required"`
	PostId  uuid.UUID `json:"post_id" validate:"required"`
	Content string    `json:"content" validate:"required,min=10,max=100"`
}
