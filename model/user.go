package model

import "github.com/google/uuid"

type FollowUnfollowPayload struct {
	UserSenderId uuid.UUID `json:"user_sender_id" validate:"required"`
}

type CreateUserPayload struct {
	Name     string `json:"name" validate:"required,min=4,max=20"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
