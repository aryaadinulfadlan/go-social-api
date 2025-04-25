package main

import (
	"errors"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/aryaadinulfadlan/go-social-api/helpers"
	"github.com/aryaadinulfadlan/go-social-api/internal"
	"github.com/aryaadinulfadlan/go-social-api/internal/store"
	"github.com/aryaadinulfadlan/go-social-api/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetUserFromContext(r *http.Request) *store.User {
	user := r.Context().Value(userCtx).(*store.User)
	return user
}

func (app *Application) ResendActivationHandler(w http.ResponseWriter, r *http.Request) {
	var payload model.ResendActivationPayload
	err := helpers.ReadFromRequestBody(r, &payload)
	if err != nil {
		app.BadRequestError(w, "Invalid JSON Body")
		return
	}
	if err := Validate.Struct(payload); err != nil {
		var validation_errors validator.ValidationErrors
		if errors.As(err, &validation_errors) {
			error_messages := make([]string, len(validation_errors))
			for idx, e := range validation_errors {
				message := GetValidationErrorMessage(e.Tag(), e.Field(), e.Param())
				error_messages[idx] = message
			}
			errorResponse := model.WebResponse{
				Code:   http.StatusBadRequest,
				Status: internal.StatusBadRequest,
				Data:   error_messages,
			}
			helpers.WriteToResponseBody(w, http.StatusBadRequest, errorResponse)
			return
		}
	}
	ctx := r.Context()
	user_data, user_err := app.Store.Users.GetExistingUser(ctx, "email", payload.Email)
	if user_err != nil {
		app.InternalServerError(w, user_err.Error())
		return
	}
	if user_data == nil {
		app.NotFoundError(w, "Invalid email")
		return
	}
	if user_data.IsActivated {
		app.BadRequestError(w, "Account is active")
		return
	}
	delete_err := app.Store.UserInvitations.DeleteUserInvitation(ctx, user_data.Id)
	if delete_err != nil {
		if errors.Is(delete_err, gorm.ErrRecordNotFound) {
			app.NotFoundError(w, delete_err.Error())
			return
		}
		app.InternalServerError(w, delete_err.Error())
		return
	}
	exp := time.Now().Add(app.Config.auth.tokenExp).UTC()
	token, token_err := app.authenticator.GenerateJWT(user_data.Id.String(), exp)
	if token_err != nil {
		app.InternalServerError(w, token_err.Error())
		return
	}
	user_invitation := store.UserInvitation{
		UserId:    user_data.Id,
		Token:     token,
		ExpiredAt: exp,
	}
	err = app.Store.UserInvitations.CreateUserInvitation(ctx, &user_invitation)
	if err != nil {
		app.InternalServerError(w, err.Error())
		return
	}
	web_response := model.WebResponse{
		Code:   http.StatusCreated,
		Status: internal.StatusCreated,
		Data: model.ResendActivationSuccess{
			Id:       user_data.Id,
			Name:     user_data.Name,
			Username: user_data.Username,
			Email:    user_data.Email,
			Token:    token,
		},
	}
	helpers.WriteToResponseBody(w, http.StatusCreated, web_response)
}

func (app *Application) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload model.LoginUserPayload
	err := helpers.ReadFromRequestBody(r, &payload)
	if err != nil {
		app.BadRequestError(w, "Invalid JSON Body")
		return
	}
	if err := Validate.Struct(payload); err != nil {
		var validation_errors validator.ValidationErrors
		if errors.As(err, &validation_errors) {
			error_messages := make([]string, len(validation_errors))
			for idx, e := range validation_errors {
				message := GetValidationErrorMessage(e.Tag(), e.Field(), e.Param())
				error_messages[idx] = message
			}
			errorResponse := model.WebResponse{
				Code:   http.StatusBadRequest,
				Status: internal.StatusBadRequest,
				Data:   error_messages,
			}
			helpers.WriteToResponseBody(w, http.StatusBadRequest, errorResponse)
			return
		}
	}
	ctx := r.Context()
	user_data, user_err := app.Store.Users.GetExistingUser(ctx, "email", payload.Email)
	if user_err != nil {
		app.InternalServerError(w, user_err.Error())
		return
	}
	if user_data == nil {
		app.UnauthorizedError(w, "Invalid email or password")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user_data.Password), []byte(payload.Password))
	if err != nil {
		app.UnauthorizedError(w, "Invalid email or password")
		return
	}
	if !user_data.IsActivated {
		app.BadRequestError(w, "Account is not activated")
		return
	}
	exp := time.Now().Add(app.Config.auth.tokenExp).UTC()
	token, token_err := app.authenticator.GenerateJWT(user_data.Id.String(), exp)
	if token_err != nil {
		app.InternalServerError(w, token_err.Error())
		return
	}
	web_response := model.WebResponse{
		Code:   http.StatusOK,
		Status: internal.StatusOK,
		Data: model.LoginSuccess{
			User: model.UserResponse{
				Id:       user_data.Id,
				Name:     user_data.Name,
				Username: user_data.Username,
				Email:    user_data.Email,
				Role:     user_data.Role.Name,
			},
			Token: token,
		},
	}
	helpers.WriteToResponseBody(w, http.StatusOK, web_response)
}

func (app *Application) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload model.CreateUserPayload
	err := helpers.ReadFromRequestBody(r, &payload)
	if err != nil {
		app.BadRequestError(w, "Invalid JSON Body")
		return
	}
	if err := Validate.Struct(payload); err != nil {
		var validation_errors validator.ValidationErrors
		if errors.As(err, &validation_errors) {
			error_messages := make([]string, len(validation_errors))
			for idx, e := range validation_errors {
				message := GetValidationErrorMessage(e.Tag(), e.Field(), e.Param())
				error_messages[idx] = message
			}
			errorResponse := model.WebResponse{
				Code:   http.StatusBadRequest,
				Status: internal.StatusBadRequest,
				Data:   error_messages,
			}
			helpers.WriteToResponseBody(w, http.StatusBadRequest, errorResponse)
			return
		}
	}
	ctx := r.Context()
	username_count, username_err := app.Store.Users.IsUserExists(ctx, "username", payload.Username)
	if username_err != nil {
		app.InternalServerError(w, username_err.Error())
		return
	}
	if username_count > 0 {
		app.BadRequestError(w, "Username already exists")
		return
	}
	email_count, email_err := app.Store.Users.IsUserExists(ctx, "email", payload.Email)
	if email_err != nil {
		app.InternalServerError(w, email_err.Error())
		return
	}
	if email_count > 0 {
		app.BadRequestError(w, "Email already exists")
		return
	}
	bytes, hash_err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if hash_err != nil {
		app.InternalServerError(w, hash_err.Error())
		return
	}
	role, role_err := app.Store.Roles.GetRole(ctx, "user")
	if role_err != nil {
		if errors.Is(role_err, gorm.ErrRecordNotFound) {
			app.NotFoundError(w, role_err.Error())
			return
		}
		app.InternalServerError(w, role_err.Error())
		return
	}
	user := store.User{
		Id:       uuid.New(),
		RoleId:   role.Id,
		Name:     payload.Name,
		Username: payload.Username,
		Email:    payload.Email,
		Password: string(bytes),
	}
	exp := time.Now().Add(app.Config.auth.tokenExp).UTC()
	token, token_err := app.authenticator.GenerateJWT(user.Id.String(), exp)
	if token_err != nil {
		app.InternalServerError(w, token_err.Error())
		return
	}
	user_invitation := store.UserInvitation{
		UserId:    user.Id,
		Token:     token,
		ExpiredAt: exp,
	}
	err = app.Store.Users.CreateUserAndInvite(ctx, &user, &user_invitation)
	if err != nil {
		app.InternalServerError(w, err.Error())
		return
	}
	web_response := model.WebResponse{
		Code:   http.StatusCreated,
		Status: internal.StatusCreated,
		Data:   user,
	}
	helpers.WriteToResponseBody(w, http.StatusCreated, web_response)
}

func (app *Application) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userId, parse_err := uuid.Parse(chi.URLParam(r, "userId"))
	if parse_err != nil {
		app.BadRequestError(w, "Invalid User ID Parameters")
		return
	}
	ctx := r.Context()
	user_data, user_err := app.Store.Users.GetUserDetail(ctx, userId)
	if user_err != nil {
		if errors.Is(user_err, gorm.ErrRecordNotFound) {
			app.NotFoundError(w, user_err.Error())
			return
		}
		app.InternalServerError(w, user_err.Error())
		return
	}
	web_response := model.WebResponse{
		Code:   http.StatusOK,
		Status: internal.StatusOK,
		Data:   user_data,
	}
	helpers.WriteToResponseBody(w, http.StatusOK, web_response)
}

func (app *Application) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userId, parse_err := uuid.Parse(chi.URLParam(r, "userId"))
	if parse_err != nil {
		app.BadRequestError(w, "Invalid User ID Parameters")
		return
	}
	ctx := r.Context()
	count, user_err := app.Store.Users.IsUserExists(ctx, "id", userId)
	if user_err != nil {
		app.InternalServerError(w, user_err.Error())
		return
	}
	if count == 0 {
		app.NotFoundError(w, "User data is not found")
		return
	}
	err := app.Store.Users.DeleteUser(ctx, userId)
	if err != nil {
		app.InternalServerError(w, err.Error())
		return
	}
	web_response := model.WebResponse{
		Code:   http.StatusOK,
		Status: internal.StatusOK,
		Data:   "Resource deleted successfully.",
	}
	helpers.WriteToResponseBody(w, http.StatusOK, web_response)
}

func (app *Application) FollowUnfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	userId, parse_err := uuid.Parse(chi.URLParam(r, "userId"))
	if parse_err != nil {
		app.BadRequestError(w, "Invalid User ID Parameters")
		return
	}
	ctx := r.Context()
	user := GetUserFromContext(r)
	target_data, target_err := app.Store.Users.GetExistingUser(ctx, "id", userId)
	if target_err != nil {
		app.InternalServerError(w, target_err.Error())
		return
	}
	if target_data == nil {
		app.NotFoundError(w, "User Target is not found")
		return
	}
	err := app.Store.Users.FollowUnfollowUser(ctx, target_data.Id, user.Id)
	if err != nil {
		app.InternalServerError(w, err.Error())
		return
	}
	web_response := model.WebResponse{
		Code:   http.StatusOK,
		Status: internal.StatusOK,
		Data: map[string]any{
			"senderId": user.Id,
			"targetId": target_data.Id,
		},
	}
	helpers.WriteToResponseBody(w, http.StatusOK, web_response)
}

func (app *Application) GetConnectionsHandler(w http.ResponseWriter, r *http.Request) {
	url := path.Base(r.URL.Path)
	actionType := strings.ToUpper(url[:1]) + strings.ToLower(url[1:])
	userId, parse_err := uuid.Parse(chi.URLParam(r, "userId"))
	if parse_err != nil {
		app.BadRequestError(w, "Invalid User ID Parameters")
		return
	}
	ctx := r.Context()
	user_data, user_err := app.Store.Users.GetExistingUser(ctx, "id", userId)
	if user_err != nil {
		app.InternalServerError(w, user_err.Error())
		return
	}
	if user_data == nil {
		app.NotFoundError(w, "User data is not found")
		return
	}
	users, err := app.Store.Users.GetConnections(ctx, user_data.Id, actionType)
	if err != nil {
		app.InternalServerError(w, err.Error())
		return
	}
	web_response := model.WebResponse{
		Code:   http.StatusOK,
		Status: internal.StatusOK,
		Data:   users,
	}
	helpers.WriteToResponseBody(w, http.StatusOK, web_response)
}

func (app *Application) ActivateUserHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr := chi.URLParam(r, "token")
	_, claims_err := app.authenticator.ParseJWT(tokenStr)
	if claims_err != nil {
		app.BadRequestError(w, claims_err.Error())
		return
	}
	ctx := r.Context()
	user, user_err := app.Store.Users.GetUserByInvitation(ctx, tokenStr)
	if user_err != nil {
		if errors.Is(user_err, gorm.ErrRecordNotFound) {
			app.NotFoundError(w, user_err.Error())
			return
		}
		app.InternalServerError(w, user_err.Error())
		return
	}
	user.IsActivated = true
	updated_user, updated_err := app.Store.Users.ActivateUser(ctx, user)
	if updated_err != nil {
		app.InternalServerError(w, updated_err.Error())
		return
	}
	web_response := model.WebResponse{
		Code:   http.StatusOK,
		Status: internal.StatusOK,
		Data:   updated_user,
	}
	helpers.WriteToResponseBody(w, http.StatusOK, web_response)
}
