package mocks

import (
	"context"

	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/stretchr/testify/mock"
)

type MockCommentRepo struct{ mock.Mock }

func (m *MockCommentRepo) Create(ctx context.Context, comment *db.Comment) error {
	return nil
}
