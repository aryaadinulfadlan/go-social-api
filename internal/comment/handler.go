package comment

import (
	"net/http"

	"github.com/aryaadinulfadlan/go-social-api/helpers"
	"github.com/aryaadinulfadlan/go-social-api/internal/shared"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
}

type HandlerImplementation struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &HandlerImplementation{
		service: service,
	}
}

func (handler *HandlerImplementation) Create(w http.ResponseWriter, r *http.Request) {
	postId, parseErr := uuid.Parse(chi.URLParam(r, "postId"))
	if parseErr != nil {
		helpers.BadRequestError(w, "Invalid Post ID Parameters")
		return
	}
	var payload *CreateCommentPayload
	err := helpers.ReadFromRequestBody(r, &payload)
	if err != nil {
		helpers.BadRequestError(w, "Invalid JSON Body")
		return
	}
	errMessages, err := helpers.ValidateStruct(payload)
	if err != nil {
		response := helpers.GenerateWebResponse(http.StatusBadRequest, shared.StatusBadRequest, errMessages)
		helpers.WriteToResponseBody(w, http.StatusBadRequest, response)
		return
	}
	ctx := r.Context()
	userContext := shared.GetUserFromContext(r)
	comment, err := handler.service.Create(ctx, payload, postId, userContext)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			helpers.NotFoundError(w, err.Error())
		default:
			helpers.InternalServerError(w, err.Error())
		}
		return
	}
	response := helpers.GenerateWebResponse(http.StatusCreated, shared.StatusCreated, comment)
	helpers.WriteToResponseBody(w, http.StatusCreated, response)
}
