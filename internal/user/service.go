package user

import (
	"context"
	"time"

	"github.com/aryaadinulfadlan/go-social-api/internal"
	"github.com/aryaadinulfadlan/go-social-api/internal/auth"
	"github.com/aryaadinulfadlan/go-social-api/internal/config"
	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/aryaadinulfadlan/go-social-api/internal/post"
	"github.com/aryaadinulfadlan/go-social-api/internal/role"
	userinvitation "github.com/aryaadinulfadlan/go-social-api/internal/user_invitation"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateAndInvite(ctx context.Context, payload CreateUserPayload) (*db.User, error)
	GetDetail(ctx context.Context, userId uuid.UUID) (*db.User, error)
	Login(ctx context.Context, payload LoginUserPayload) (*LoginSuccess, error)
	ResendActivation(ctx context.Context, email string) (*ResendActivationSuccess, error)
	FollowUnfollow(ctx context.Context, userContext *db.User, userId uuid.UUID) (*FollowUnfollowSuccess, error)
	GetConnections(ctx context.Context, userId uuid.UUID, actionType string) ([]*db.User, error)
	Activate(ctx context.Context, tokenStr string) (*db.User, error)
	Delete(ctx context.Context, userId uuid.UUID) error
	GetFeeds(ctx context.Context, userId uuid.UUID, params post.PostParams) (*internal.PaginationResponse, error)
}

type ServiceImplementation struct {
	authenticator            auth.Authenticator
	repository               Repository
	roleRepository           role.Repository
	userInvitationRepository userinvitation.Repository
}

func NewService(authenticator auth.Authenticator, repository Repository, roleRepository role.Repository, userInvitationRepository userinvitation.Repository) Service {
	return &ServiceImplementation{
		authenticator,
		repository,
		roleRepository,
		userInvitationRepository,
	}
}

func (service *ServiceImplementation) CreateAndInvite(ctx context.Context, payload CreateUserPayload) (*db.User, error) {
	user, err := service.repository.GetByUsernameEmail(ctx, payload.Username, payload.Email)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, internal.ErrUserExists
	}
	bytes, hash_err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if hash_err != nil {
		return nil, err
	}
	role, role_err := service.roleRepository.GetRole(ctx, "user")
	if role_err != nil {
		return nil, err
	}
	user = &db.User{
		Id:       uuid.New(),
		RoleId:   role.Id,
		Name:     payload.Name,
		Username: payload.Username,
		Email:    payload.Email,
		Password: string(bytes),
	}
	exp := time.Now().Add(config.Auth.TokenExp).UTC()
	token, token_err := service.authenticator.GenerateJWT(user.Id.String(), exp)
	if token_err != nil {
		return nil, err
	}
	user_invitation := &db.UserInvitation{
		UserId:    user.Id,
		Token:     token,
		ExpiredAt: exp,
	}
	err = service.repository.CreateAndInvite(ctx, user, user_invitation)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *ServiceImplementation) GetDetail(ctx context.Context, userId uuid.UUID) (*db.User, error) {
	user, err := service.repository.GetDetail(ctx, userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *ServiceImplementation) Login(ctx context.Context, payload LoginUserPayload) (*LoginSuccess, error) {
	user, err := service.repository.GetByUsernameEmail(ctx, "", payload.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, internal.ErrLoginInvalid
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return nil, internal.ErrLoginInvalid
	}
	if !user.IsActivated {
		return nil, internal.ErrAccountInactive
	}
	exp := time.Now().Add(config.Auth.TokenExp).UTC()
	token, err := service.authenticator.GenerateJWT(user.Id.String(), exp)
	if err != nil {
		return nil, err
	}
	loginSuccess := LoginSuccess{
		User: UserResponse{
			Id:       user.Id,
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role.Name,
		},
		Token: token,
	}
	return &loginSuccess, nil
}

func (service *ServiceImplementation) ResendActivation(ctx context.Context, email string) (*ResendActivationSuccess, error) {
	user, err := service.repository.GetByUsernameEmail(ctx, "", email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, internal.ErrEmailInvalid
	}
	if user.IsActivated {
		return nil, internal.ErrAccountActive
	}
	err = service.userInvitationRepository.DeleteUserInvitation(ctx, user.Id)
	if err != nil {
		return nil, err
	}
	exp := time.Now().Add(config.Auth.TokenExp).UTC()
	token, err := service.authenticator.GenerateJWT(user.Id.String(), exp)
	if err != nil {
		return nil, err
	}
	user_invitation := db.UserInvitation{
		UserId:    user.Id,
		Token:     token,
		ExpiredAt: exp,
	}
	err = service.userInvitationRepository.CreateUserInvitation(ctx, &user_invitation)
	if err != nil {
		return nil, err
	}
	resendActivationSuccess := ResendActivationSuccess{
		Id:       user.Id,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}
	return &resendActivationSuccess, nil
}

func (service *ServiceImplementation) FollowUnfollow(ctx context.Context, userContext *db.User, userId uuid.UUID) (*FollowUnfollowSuccess, error) {
	userTarget, err := service.repository.GetById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if userTarget == nil {
		return nil, internal.ErrNotFound
	}
	message, err := service.repository.FollowUnfollow(ctx, userTarget.Id, userContext.Id)
	if err != nil {
		return nil, err
	}
	followUnfollowSuccess := &FollowUnfollowSuccess{
		SenderId: userContext.Id,
		TargetId: userTarget.Id,
		Message:  message,
	}
	return followUnfollowSuccess, nil
}

func (service *ServiceImplementation) GetConnections(ctx context.Context, userId uuid.UUID, actionType string) ([]*db.User, error) {
	user, err := service.repository.GetById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, internal.ErrNotFound
	}
	users, err := service.repository.GetConnections(ctx, user.Id, actionType)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (service *ServiceImplementation) Activate(ctx context.Context, tokenStr string) (*db.User, error) {
	user, err := service.repository.GetByInvitation(ctx, tokenStr)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, internal.ErrNotFound
	}
	user.IsActivated = true
	user, err = service.repository.Activate(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *ServiceImplementation) Delete(ctx context.Context, userId uuid.UUID) error {
	user, err := service.repository.GetById(ctx, userId)
	if err != nil {
		return err
	}
	if user == nil {
		return internal.ErrNotFound
	}
	err = service.repository.Delete(ctx, userId)
	if err != nil {
		return err
	}
	return nil
}

func (service *ServiceImplementation) GetFeeds(ctx context.Context, userId uuid.UUID, params post.PostParams) (*internal.PaginationResponse, error) {
	feeds, total, err := service.repository.GetFeeds(ctx, userId, &params)
	if err != nil {
		return nil, err
	}
	pagination := internal.NewPaginationMeta(params.Page, params.PerPage, int(total))
	paginatedFeeds := &internal.PaginationResponse{
		Items:      feeds,
		Pagination: pagination,
	}
	return paginatedFeeds, nil
}
