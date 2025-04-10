package main

import (
	"errors"
	"net/http"

	"github.com/aryaadinulfadlan/go-social-api/helpers"
	"github.com/aryaadinulfadlan/go-social-api/internal"
	"github.com/aryaadinulfadlan/go-social-api/internal/store"
	"github.com/aryaadinulfadlan/go-social-api/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,min=4,max=10"`
	Content string   `json:"content" validate:"required,min=8,max=20"`
	Tags    []string `json:"tags" validate:"required,notempty,dive,required,min=10,max=20"`
}

func (app *Application) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
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
	user_id, _ := uuid.Parse("030e656e-cc3e-47f3-813a-33a3d50b5373")
	post := store.Post{
		UserId:  user_id,
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
	}
	ctx := r.Context()
	err = app.Store.Posts.Create(ctx, &post)
	if err != nil {
		app.InternalServerError(w, err.Error())
		return
	}
	web_response := model.WebResponse{
		Code:   http.StatusCreated,
		Status: internal.StatusCreated,
		Data:   post,
	}
	helpers.WriteToResponseBody(w, http.StatusCreated, web_response)
}

func (app *Application) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	postId, parse_err := uuid.Parse(chi.URLParam(r, "postId"))
	if parse_err != nil {
		app.BadRequestError(w, "Invalid URL Parameters")
		return
	}
	ctx := r.Context()
	post_data, post_err := app.Store.Posts.GetById(ctx, postId)
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
		Status: internal.StatusOK,
		Data:   post_data,
	}
	helpers.WriteToResponseBody(w, http.StatusOK, web_response)
}
