package mocks

import (
	"context"
	"database/sql"

	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/aryaadinulfadlan/go-social-api/internal/post"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
	MockUser *db.User
}

func (m *MockUserRepo) CreateAndInvite(ctx context.Context, user *db.User, user_invitation *db.UserInvitation) error {
	args := m.Called(user)
	return args.Error(0)
	// return nil
}
func (m *MockUserRepo) GetDetail(ctx context.Context, userId uuid.UUID) (*db.User, error) {
	if userId == m.MockUser.Id {
		return m.MockUser, nil
	}
	return nil, sql.ErrNoRows
	// args := m.Called(userId)
	// if user, ok := args.Get(0).(*db.User); ok {
	// 	return user, args.Error(1)
	// }
	// return nil, args.Error(1)
}
func (m *MockUserRepo) GetByInvitation(ctx context.Context, token string) (*db.User, error) {
	return nil, nil
}
func (m *MockUserRepo) GetById(ctx context.Context, userId uuid.UUID) (*db.User, error) {
	return nil, nil
}
func (m *MockUserRepo) GetByUsernameEmail(ctx context.Context, username string, email string) (*db.User, error) {
	return nil, nil
}
func (m *MockUserRepo) FollowUnfollow(ctx context.Context, targetId uuid.UUID, senderId uuid.UUID) (string, error) {
	return "", nil
}
func (m *MockUserRepo) GetConnections(ctx context.Context, userId uuid.UUID, actionType string) ([]*db.User, error) {
	return nil, nil
}
func (m *MockUserRepo) Activate(ctx context.Context, user *db.User) (*db.User, error) {
	return nil, nil
}
func (m *MockUserRepo) Delete(ctx context.Context, userId uuid.UUID) error {
	return nil
}
func (m *MockUserRepo) GetFeeds(ctx context.Context, userId uuid.UUID, params *post.PostParams) ([]*db.PostWithMetadata, int64, error) {
	return nil, 0, nil
}
