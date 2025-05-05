package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockPermissionRepo struct{ mock.Mock }

func (m *MockPermissionRepo) GetPermissionNamesByRoleId(ctx context.Context, role_id uuid.UUID) ([]string, error) {
	return nil, nil
}
