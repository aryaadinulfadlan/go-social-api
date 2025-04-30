package main

import (
	"errors"
	"net/http"

	"github.com/aryaadinulfadlan/go-social-api/helpers"
	"github.com/aryaadinulfadlan/go-social-api/internal/shared"
	"github.com/aryaadinulfadlan/go-social-api/internal/store"
	"github.com/aryaadinulfadlan/go-social-api/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (app *Application) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var payload model.CreatePostPayload
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
				Status: shared.StatusBadRequest,
				Data:   error_messages,
			}
			helpers.WriteToResponseBody(w, http.StatusBadRequest, errorResponse)
			return
		}
	}
	user := GetUserFromContext(r)
	post := store.Post{
		UserId:  user.Id,
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
	}
	ctx := r.Context()
	err = app.Store.Posts.CreatePost(ctx, &post)
	if err != nil {
		app.InternalServerError(w, err.Error())
		return
	}
	web_response := model.WebResponse{
		Code:   http.StatusCreated,
		Status: shared.StatusCreated,
		Data:   post,
	}
	helpers.WriteToResponseBody(w, http.StatusCreated, web_response)
}

func (app *Application) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	postId, parse_err := uuid.Parse(chi.URLParam(r, "postId"))
	if parse_err != nil {
		app.BadRequestError(w, "Invalid Post ID Parameters")
		return
	}
	ctx := r.Context()
	post_data, post_err := app.Store.Posts.GetPost(ctx, postId)
	if post_err != nil {
		if errors.Is(post_err, gorm.ErrRecordNotFound) {
			app.NotFoundError(w, post_err.Error())
			return
		}
		app.InternalServerError(w, post_err.Error())
		return
	}
	web_response := model.WebResponse{
		Code:   http.StatusOK,
		Status: shared.StatusOK,
		Data:   post_data,
	}
	helpers.WriteToResponseBody(w, http.StatusOK, web_response)
}

func (app *Application) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	postId, parse_err := uuid.Parse(chi.URLParam(r, "postId"))
	var payload model.UpdatePostPayload
	if parse_err != nil {
		app.BadRequestError(w, "Invalid Post ID Parameters")
		return
	}
	payload.Id = postId
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
				Status: shared.StatusBadRequest,
				Data:   error_messages,
			}
			helpers.WriteToResponseBody(w, http.StatusBadRequest, errorResponse)
			return
		}
	}
	ctx := r.Context()
	post, err := app.Store.Posts.CheckPostExists(ctx, postId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			app.NotFoundError(w, err.Error())
			return
		}
		app.InternalServerError(w, err.Error())
		return
	}
	if post.UserId != user.Id && user.Role.Name != "admin" {
		app.ForbiddenError(w, "You do not have permission to access this resource.")
		return
	}
	post.Title = payload.Title
	post.Content = payload.Content
	post.Tags = payload.Tags
	post_data, post_err := app.Store.Posts.UpdatePost(ctx, post)
	if post_err != nil {
		app.InternalServerError(w, post_err.Error())
		return
	}
	web_response := model.WebResponse{
		Code:   http.StatusOK,
		Status: shared.StatusOK,
		Data:   post_data,
	}
	helpers.WriteToResponseBody(w, http.StatusOK, web_response)
}

func (app *Application) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	postId, parse_err := uuid.Parse(chi.URLParam(r, "postId"))
	if parse_err != nil {
		app.BadRequestError(w, "Invalid Post ID Parameters")
		return
	}
	ctx := r.Context()
	post, err := app.Store.Posts.CheckPostExists(ctx, postId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			app.NotFoundError(w, err.Error())
			return
		}
		app.InternalServerError(w, err.Error())
		return
	}
	if post.UserId != user.Id && user.Role.Name != "admin" {
		app.ForbiddenError(w, "You do not have permission to access this resource.")
		return
	}
	err = app.Store.Posts.DeletePost(ctx, post.Id)
	if err != nil {
		app.InternalServerError(w, err.Error())
		return
	}
	web_response := model.WebResponse{
		Code:   http.StatusOK,
		Status: shared.StatusOK,
		Data:   "Resource deleted successfully.",
	}
	helpers.WriteToResponseBody(w, http.StatusOK, web_response)
}

func (app *Application) GetPostFeedHandler(w http.ResponseWriter, r *http.Request) {
	params := &model.PostParams{
		PerPage: 10,
		Page:    1,
		Sort:    "DESC",
	}
	params, err := params.Parse(r)
	if err != nil {
		app.BadRequestError(w, err.Error())
		return
	}
	if err := Validate.Struct(params); err != nil {
		var validation_errors validator.ValidationErrors
		if errors.As(err, &validation_errors) {
			error_messages := make([]string, len(validation_errors))
			for idx, e := range validation_errors {
				message := GetValidationErrorMessage(e.Tag(), e.Field(), e.Param())
				error_messages[idx] = message
			}
			errorResponse := model.WebResponse{
				Code:   http.StatusBadRequest,
				Status: shared.StatusBadRequest,
				Data:   error_messages,
			}
			helpers.WriteToResponseBody(w, http.StatusBadRequest, errorResponse)
			return
		}
	}
	ctx := r.Context()
	user := GetUserFromContext(r)
	feed, total, err := app.Store.Posts.GetPostFeed(ctx, user.Id, params)
	if err != nil {
		app.InternalServerError(w, err.Error())
		return
	}
	pagination := model.NewPaginationMeta(params.Page, params.PerPage, int(total))
	web_response := model.WebResponse{
		Code:   http.StatusOK,
		Status: shared.StatusOK,
		Data: model.PaginationResponse{
			Items:      feed,
			Pagination: pagination,
		},
	}
	helpers.WriteToResponseBody(w, http.StatusOK, web_response)
}
