package helpers

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

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
func ValidateStruct(payload any) ([]string, error) {
	err := Validate.Struct(payload)
	if err != nil {
		var validation_errors validator.ValidationErrors
		if errors.As(err, &validation_errors) {
			error_messages := make([]string, len(validation_errors))
			for idx, e := range validation_errors {
				message := GetValidationErrorMessage(e.Tag(), e.Field(), e.Param())
				error_messages[idx] = message
			}
			return error_messages, err
		}
	}
	return nil, nil
}
func InternalServerError(w http.ResponseWriter, err string) {
	WriteErrorResponse(w, http.StatusInternalServerError, internal.StatusInternalServerError, err)
}
func NotFoundError(w http.ResponseWriter, err string) {
	WriteErrorResponse(w, http.StatusNotFound, internal.StatusNotFound, err)
}
func MethodNotAllowedError(w http.ResponseWriter, err string) {
	WriteErrorResponse(w, http.StatusMethodNotAllowed, internal.StatusMethodNotAllowed, err)
}
func BadRequestError(w http.ResponseWriter, err string) {
	WriteErrorResponse(w, http.StatusBadRequest, internal.StatusBadRequest, err)
}
func UnauthorizedError(w http.ResponseWriter, err string) {
	WriteErrorResponse(w, http.StatusUnauthorized, internal.StatusUnauthorized, err)
}
func ForbiddenError(w http.ResponseWriter, err string) {
	WriteErrorResponse(w, http.StatusForbidden, internal.StatusForbidden, err)
}
func RateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	// app.Logger.Warn("rate limit exceeded", "method", r.Method, "path", r.URL.Path)
	w.Header().Set("Retry-After", retryAfter)
	WriteErrorResponse(w, http.StatusTooManyRequests, internal.StatusTooManyRequests, "Rate limit exceeded, retry after: "+retryAfter)
}
