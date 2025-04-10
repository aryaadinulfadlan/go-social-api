package main

import (
	"net/http"

	"github.com/aryaadinulfadlan/go-social-api/helpers"
	"github.com/aryaadinulfadlan/go-social-api/internal"
)

func (app *Application) InternalServerError(w http.ResponseWriter, err string) {
	helpers.WriteErrorResponse(w, http.StatusInternalServerError, internal.StatusInternalServerError, err)
}
func (app *Application) NotFoundError(w http.ResponseWriter, err string) {
	helpers.WriteErrorResponse(w, http.StatusNotFound, internal.StatusNotFound, err)
}
func (app *Application) MethodNotAllowedError(w http.ResponseWriter, err string) {
	helpers.WriteErrorResponse(w, http.StatusMethodNotAllowed, internal.StatusMethodNotAllowed, err)
}
func (app *Application) BadRequestError(w http.ResponseWriter, err string) {
	helpers.WriteErrorResponse(w, http.StatusBadRequest, internal.StatusBadRequest, err)
}
