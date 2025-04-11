package model

import "github.com/google/uuid"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,min=4,max=20"`
	Content string   `json:"content" validate:"required,min=8,max=100"`
	Tags    []string `json:"tags" validate:"required,notempty,dive,required,min=4,max=20"`
}
type UpdatePostPayload struct {
	Id      uuid.UUID `json:"id" validate:"required"`
	Title   string    `json:"title" validate:"required,min=4,max=20"`
	Content string    `json:"content" validate:"required,min=8,max=100"`
	Tags    []string  `json:"tags" validate:"required,notempty,dive,required,min=4,max=20"`
}
