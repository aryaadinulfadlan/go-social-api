package post

import (
	"net/http"

	"github.com/aryaadinulfadlan/go-social-api/helpers"
	"github.com/aryaadinulfadlan/go-social-api/internal"
	"github.com/aryaadinulfadlan/go-social-api/internal/shared"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetDetail(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
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
	var payload *CreatePostPayload
	err := helpers.ReadFromRequestBody(r, &payload)
	if err != nil {
		helpers.BadRequestError(w, "Invalid JSON Body")
		return
	}
	errMessages, err := helpers.ValidateStruct(payload)
	if err != nil {
		response := helpers.GenerateWebResponse(http.StatusBadRequest, internal.StatusBadRequest, errMessages)
		helpers.WriteToResponseBody(w, http.StatusBadRequest, response)
		return
	}
	ctx := r.Context()
	userContext := shared.GetUserFromContext(r)
	post, err := handler.service.Create(ctx, userContext.Id, payload)
	if err != nil {
		helpers.InternalServerError(w, err.Error())
		return
	}
	response := helpers.GenerateWebResponse(http.StatusCreated, internal.StatusCreated, post)
	helpers.WriteToResponseBody(w, http.StatusCreated, response)
}
func (handler *HandlerImplementation) GetDetail(w http.ResponseWriter, r *http.Request) {
	postId, parseErr := uuid.Parse(chi.URLParam(r, "postId"))
	if parseErr != nil {
		helpers.BadRequestError(w, "Invalid Post ID Parameters")
		return
	}
	ctx := r.Context()
	post, err := handler.service.GetDetail(ctx, postId)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			helpers.NotFoundError(w, err.Error())
		default:
			helpers.InternalServerError(w, err.Error())
		}
		return
	}
	response := helpers.GenerateWebResponse(http.StatusOK, internal.StatusOK, post)
	helpers.WriteToResponseBody(w, http.StatusOK, response)
}
func (handler *HandlerImplementation) Update(w http.ResponseWriter, r *http.Request) {
	postId, parseErr := uuid.Parse(chi.URLParam(r, "postId"))
	if parseErr != nil {
		helpers.BadRequestError(w, "Invalid Post ID Parameters")
		return
	}
	var payload UpdatePostPayload
	payload.Id = postId
	err := helpers.ReadFromRequestBody(r, &payload)
	if err != nil {
		helpers.BadRequestError(w, "Invalid JSON Body")
		return
	}
	errMessages, err := helpers.ValidateStruct(payload)
	if err != nil {
		response := helpers.GenerateWebResponse(http.StatusBadRequest, internal.StatusBadRequest, errMessages)
		helpers.WriteToResponseBody(w, http.StatusBadRequest, response)
		return
	}
	ctx := r.Context()
	userContext := shared.GetUserFromContext(r)
	post, err := handler.service.Update(ctx, postId, userContext, &payload)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			helpers.NotFoundError(w, err.Error())
		case internal.ErrForbidden:
			helpers.ForbiddenError(w, err.Error())
		default:
			helpers.InternalServerError(w, err.Error())
		}
		return
	}
	response := helpers.GenerateWebResponse(http.StatusOK, internal.StatusOK, post)
	helpers.WriteToResponseBody(w, http.StatusOK, response)
}
func (handler *HandlerImplementation) Delete(w http.ResponseWriter, r *http.Request) {
	postId, parseErr := uuid.Parse(chi.URLParam(r, "postId"))
	if parseErr != nil {
		helpers.BadRequestError(w, "Invalid Post ID Parameters")
		return
	}
	ctx := r.Context()
	userContext := shared.GetUserFromContext(r)
	err := handler.service.Delete(ctx, postId, userContext)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			helpers.NotFoundError(w, err.Error())
		case internal.ErrForbidden:
			helpers.ForbiddenError(w, err.Error())
		default:
			helpers.InternalServerError(w, err.Error())
		}
		return
	}
	response := helpers.GenerateWebResponse(http.StatusOK, internal.StatusOK, "Resource deleted successfully.")
	helpers.WriteToResponseBody(w, http.StatusOK, response)
}
