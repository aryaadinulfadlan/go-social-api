package main

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/aryaadinulfadlan/go-social-api/helpers"
	"github.com/aryaadinulfadlan/go-social-api/internal"
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
	Validate.RegisterValidation("notempty", func(fl validator.FieldLevel) bool {
		field := fl.Field()
		return field.Kind() == reflect.Slice && field.Len() > 0
	})
}
func GetValidationErrorMessage(tag string, field string, param string) string {
	var message string
	switch tag {
	case "required":
		message = fmt.Sprintf("%s is required.", field)
	case "max":
		message = fmt.Sprintf("%s must be at most %s characters.", field, param)
	case "min":
		message = fmt.Sprintf("%s must be at least %s characters.", field, param)
	case "len":
		message = fmt.Sprintf("%s must be exactly %s characters.", field, param)
	case "notempty":
		message = fmt.Sprintf("%s must not be empty.", field)
	default:
		message = fmt.Sprintf("%s is not valid.", field)
	}
	return message
}
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
func (app *Application) UnauthorizedError(w http.ResponseWriter, err string) {
	helpers.WriteErrorResponse(w, http.StatusUnauthorized, internal.StatusUnauthorized, err)
}
func (app *Application) ForbiddenError(w http.ResponseWriter, err string) {
	helpers.WriteErrorResponse(w, http.StatusForbidden, internal.StatusForbidden, err)
}
