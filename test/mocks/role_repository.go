package mocks

import (
	"context"

	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/stretchr/testify/mock"
)

type MockRoleRepo struct{ mock.Mock }

func (m *MockRoleRepo) GetRole(ctx context.Context, name string) (*db.Role, error) {
	return nil, nil
}
