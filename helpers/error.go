package helpers

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/aryaadinulfadlan/go-social-api/internal/logger"
	"github.com/aryaadinulfadlan/go-social-api/internal/shared"
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
func GenerateWebResponse(code int, status string, data any) *shared.WebResponse {
	return &shared.WebResponse{
		Code:   code,
		Status: status,
		Data:   data,
	}
}
func InternalServerError(w http.ResponseWriter, err string) {
	WriteErrorResponse(w, http.StatusInternalServerError, shared.StatusInternalServerError, err)
}
func NotFoundError(w http.ResponseWriter, err string) {
	WriteErrorResponse(w, http.StatusNotFound, shared.StatusNotFound, err)
}
func MethodNotAllowedError(w http.ResponseWriter, err string) {
	WriteErrorResponse(w, http.StatusMethodNotAllowed, shared.StatusMethodNotAllowed, err)
}
func BadRequestError(w http.ResponseWriter, err string) {
	WriteErrorResponse(w, http.StatusBadRequest, shared.StatusBadRequest, err)
}
func UnauthorizedError(w http.ResponseWriter, err string) {
	WriteErrorResponse(w, http.StatusUnauthorized, shared.StatusUnauthorized, err)
}
func ForbiddenError(w http.ResponseWriter, err string) {
	WriteErrorResponse(w, http.StatusForbidden, shared.StatusForbidden, err)
}
func RateLimitExceededResponse(w http.ResponseWriter, r *http.Request, err string) {
	logger.Logger.Warn("rate limit exceeded", "method", r.Method, "path", r.URL.Path)
	WriteErrorResponse(w, http.StatusTooManyRequests, shared.StatusTooManyRequests, err)
}
