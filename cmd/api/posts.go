package main

import (
	"errors"
	"net/http"

	"github.com/aryaadinulfadlan/go-social-api/helpers"
	"github.com/aryaadinulfadlan/go-social-api/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=20"`
	Content string   `json:"content" validate:"required,max=100"`
	Tags    []string `json:"tags"`
}

func (app *Application) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	err := helpers.ReadFromRequestBody(r, &payload)
	if err != nil {
		app.BadRequestError(w, "Invalid JSON Body")
		return
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
	helpers.WriteToResponseBody(w, http.StatusCreated, post)
	// helpers.JSONFormatting(post)
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
	helpers.WriteToResponseBody(w, http.StatusOK, post_data)
	// helpers.JSONFormatting(post)
}
