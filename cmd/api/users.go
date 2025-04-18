package main

import (
	"errors"
	"net/http"
	"path"
	"strings"

	"github.com/aryaadinulfadlan/go-social-api/helpers"
	"github.com/aryaadinulfadlan/go-social-api/internal"
	"github.com/aryaadinulfadlan/go-social-api/internal/store"
	"github.com/aryaadinulfadlan/go-social-api/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
	user := store.User{
		Name:     payload.Name,
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	}
	ctx := r.Context()
	err = app.Store.Users.CreateUser(ctx, &user)
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

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (app *Application) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userId, parse_err := uuid.Parse(chi.URLParam(r, "userId"))
	if parse_err != nil {
		app.BadRequestError(w, "Invalid User ID Parameters")
		return
	}
	ctx := r.Context()
	user_data, user_err := app.Store.Users.GetUser(ctx, userId)
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

func (app *Application) FollowUnfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	userId, parse_err := uuid.Parse(chi.URLParam(r, "userId"))
	var payload model.FollowUnfollowPayload
	if parse_err != nil {
		app.BadRequestError(w, "Invalid User ID Parameters")
		return
	}
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
	target_data, target_err := app.Store.Users.CheckUserExists(ctx, userId)
	if target_err != nil {
		if errors.Is(target_err, gorm.ErrRecordNotFound) {
			app.NotFoundError(w, "User Target is not found")
			return
		}
		app.InternalServerError(w, target_err.Error())
		return
	}
	sender_data, sender_err := app.Store.Users.CheckUserExists(ctx, payload.UserSenderId)
	if sender_err != nil {
		if errors.Is(sender_err, gorm.ErrRecordNotFound) {
			app.NotFoundError(w, "User Sender is not found")
			return
		}
		app.InternalServerError(w, sender_err.Error())
		return
	}
	err = app.Store.Users.FollowUnfollowUser(ctx, target_data.Id, sender_data.Id)
	if err != nil {
		app.InternalServerError(w, err.Error())
		return
	}
	web_response := model.WebResponse{
		Code:   http.StatusOK,
		Status: internal.StatusOK,
		Data: map[string]any{
			"senderId": sender_data.Id,
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
	user_data, user_err := app.Store.Users.CheckUserExists(ctx, userId)
	if user_err != nil {
		if errors.Is(user_err, gorm.ErrRecordNotFound) {
			app.NotFoundError(w, user_err.Error())
			return
		}
		app.InternalServerError(w, user_err.Error())
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
