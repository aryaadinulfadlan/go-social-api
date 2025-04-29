package user

import "github.com/google/uuid"

type ResendActivationPayload struct {
	Email string `json:"email" validate:"required,email"`
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
	Role     string    `json:"role"`
}

type LoginSuccess struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}
type FollowUnfollowSuccess struct {
	SenderId uuid.UUID `json:"sender_id"`
	TargetId uuid.UUID `json:"target_id"`
	Message  string    `json:"message"`
}
type ResendActivationSuccess struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Token    string    `json:"token"`
}
