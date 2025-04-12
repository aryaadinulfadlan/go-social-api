package main

import (
	"errors"
	"net/http"

	"github.com/aryaadinulfadlan/go-social-api/helpers"
	"github.com/aryaadinulfadlan/go-social-api/internal"
	"github.com/aryaadinulfadlan/go-social-api/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
