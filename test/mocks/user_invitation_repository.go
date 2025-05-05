package mocks

import (
	"context"

	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserInvitationRepo struct{ mock.Mock }

func (m *MockUserInvitationRepo) CreateUserInvitation(ctx context.Context, user_invitation *db.UserInvitation) error {
	return nil
}
func (m *MockUserInvitationRepo) DeleteUserInvitation(ctx context.Context, userId uuid.UUID) error {
	return nil
}
