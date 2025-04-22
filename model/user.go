package model

import "github.com/google/uuid"

type FollowUnfollowPayload struct {
	UserSenderId uuid.UUID `json:"user_sender_id" validate:"required"`
}

type CreateUserPayload struct {
	Name     string `json:"name" validate:"required,min=4,max=20"`
	Username string `json:"username" validate:"required,min=6,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}
type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type UserResponse struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

type LoginSuccess struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}
