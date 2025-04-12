package model

import "github.com/google/uuid"

type FollowUnfollowPayload struct {
	UserSenderId uuid.UUID `json:"user_sender_id" validate:"required"`
}
