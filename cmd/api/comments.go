package main

import (
	"errors"
	"net/http"

	"github.com/aryaadinulfadlan/go-social-api/helpers"
	"github.com/aryaadinulfadlan/go-social-api/internal"
	"github.com/aryaadinulfadlan/go-social-api/internal/store"
	"github.com/aryaadinulfadlan/go-social-api/model"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func (app *Application) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload model.CreateCommentPayload
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
	count, user_err := app.Store.Users.IsUserExists(ctx, "id", payload.UserId)
	if user_err != nil {
		app.InternalServerError(w, user_err.Error())
		return
	}
	if count == 0 {
		app.NotFoundError(w, "User data is not found")
		return
	}
	_, post_err := app.Store.Posts.CheckPostExists(ctx, payload.PostId)
	if post_err != nil {
		if errors.Is(post_err, gorm.ErrRecordNotFound) {
			app.NotFoundError(w, "Post data is not found")
			return
		}
		app.InternalServerError(w, post_err.Error())
		return
	}
	comment := store.Comment{
		UserId:  payload.UserId,
		PostId:  payload.PostId,
		Content: payload.Content,
	}
	err = app.Store.Comments.CreateComment(ctx, &comment)
	if err != nil {
		app.InternalServerError(w, err.Error())
		return
	}
	web_response := model.WebResponse{
		Code:   http.StatusCreated,
		Status: internal.StatusCreated,
		Data:   comment,
	}
	helpers.WriteToResponseBody(w, http.StatusCreated, web_response)
}
