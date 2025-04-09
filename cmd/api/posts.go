package main

import (
	"net/http"

	"github.com/aryaadinulfadlan/go-social-api/helpers"
	"github.com/aryaadinulfadlan/go-social-api/internal/store"
	"github.com/google/uuid"
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
		helpers.WriteErrorResponse(w, http.StatusBadRequest, "BAD REQUEST", "Invalid JSON Body")
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
		helpers.WriteErrorResponse(w, http.StatusInternalServerError, "INTERNAL SERVER ERROR", err.Error())
	}
	helpers.WriteToResponseBody(w, http.StatusCreated, post)
	helpers.JSONFormatting(post)
}
