package mocks

import (
	"context"

	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockPostRepo struct{ mock.Mock }

func (m *MockPostRepo) Create(ctx context.Context, post *db.Post) error {
	return nil
}
func (m *MockPostRepo) GetDetail(ctx context.Context, postId uuid.UUID) (*db.Post, error) {
	return nil, nil
}
func (m *MockPostRepo) GetById(ctx context.Context, postId uuid.UUID) (*db.Post, error) {
	return nil, nil
}
func (m *MockPostRepo) Update(ctx context.Context, post *db.Post) (*db.Post, error) {
	return nil, nil
}
func (m *MockPostRepo) Delete(ctx context.Context, postId uuid.UUID) error {
	return nil
}
