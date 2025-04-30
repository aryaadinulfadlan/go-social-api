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

func (app *Application) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	postId, parse_err := uuid.Parse(chi.URLParam(r, "postId"))
	if parse_err != nil {
		app.BadRequestError(w, "Invalid Post ID Parameters")
		return
	}
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
				Status: shared.StatusBadRequest,
				Data:   error_messages,
			}
			helpers.WriteToResponseBody(w, http.StatusBadRequest, errorResponse)
			return
		}
	}
	ctx := r.Context()
	user := GetUserFromContext(r)
	_, post_err := app.Store.Posts.CheckPostExists(ctx, postId)
	if post_err != nil {
		if errors.Is(post_err, gorm.ErrRecordNotFound) {
			app.NotFoundError(w, "Post data is not found")
			return
		}
		app.InternalServerError(w, post_err.Error())
		return
	}
	comment := store.Comment{
		UserId:  user.Id,
		PostId:  postId,
		Content: payload.Content,
	}
	err = app.Store.Comments.CreateComment(ctx, &comment)
	if err != nil {
		app.InternalServerError(w, err.Error())
		return
	}
	web_response := model.WebResponse{
		Code:   http.StatusCreated,
		Status: shared.StatusCreated,
		Data:   comment,
	}
	helpers.WriteToResponseBody(w, http.StatusCreated, web_response)
}
