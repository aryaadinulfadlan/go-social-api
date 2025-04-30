package user

import (
	"net/http"
	"path"
	"strings"

	"github.com/aryaadinulfadlan/go-social-api/helpers"
	"github.com/aryaadinulfadlan/go-social-api/internal/auth"
	"github.com/aryaadinulfadlan/go-social-api/internal/post"
	"github.com/aryaadinulfadlan/go-social-api/internal/shared"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Handler interface {
	CreateAndInvite(w http.ResponseWriter, r *http.Request)
	GetDetail(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	ResendActivation(w http.ResponseWriter, r *http.Request)
	FollowUnfollow(w http.ResponseWriter, r *http.Request)
	GetConnections(w http.ResponseWriter, r *http.Request)
	Activate(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	GetFeeds(w http.ResponseWriter, r *http.Request)
}

type HandlerImplementation struct {
	authenticator auth.Authenticator
	service       Service
}

func NewHandler(authenticator auth.Authenticator, service Service) Handler {
	return &HandlerImplementation{
		authenticator: authenticator,
		service:       service,
	}
}

func (handler *HandlerImplementation) CreateAndInvite(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserPayload
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
	user, err := handler.service.CreateAndInvite(ctx, payload)
	if err != nil {
		switch err {
		case shared.ErrUserExists:
			helpers.BadRequestError(w, err.Error())
		default:
			helpers.InternalServerError(w, err.Error())
		}
		return
	}
	response := helpers.GenerateWebResponse(http.StatusCreated, shared.StatusCreated, user)
	helpers.WriteToResponseBody(w, http.StatusCreated, response)
}

func (handler *HandlerImplementation) GetDetail(w http.ResponseWriter, r *http.Request) {
	userId, parseErr := uuid.Parse(chi.URLParam(r, "userId"))
	if parseErr != nil {
		helpers.BadRequestError(w, "Invalid User ID Parameters")
		return
	}
	ctx := r.Context()
	user, err := handler.service.GetDetail(ctx, userId)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			helpers.NotFoundError(w, err.Error())
		default:
			helpers.InternalServerError(w, err.Error())
		}
		return
	}
	response := helpers.GenerateWebResponse(http.StatusOK, shared.StatusOK, user)
	helpers.WriteToResponseBody(w, http.StatusOK, response)
}

func (handler *HandlerImplementation) Login(w http.ResponseWriter, r *http.Request) {
	var payload LoginUserPayload
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
	loginSuccess, err := handler.service.Login(ctx, payload)
	if err != nil {
		switch err {
		case shared.ErrLoginInvalid:
			helpers.UnauthorizedError(w, err.Error())
		case shared.ErrAccountInactive:
			helpers.BadRequestError(w, err.Error())
		default:
			helpers.InternalServerError(w, err.Error())
		}
		return
	}
	response := helpers.GenerateWebResponse(http.StatusOK, shared.StatusOK, loginSuccess)
	helpers.WriteToResponseBody(w, http.StatusOK, response)
}

func (handler *HandlerImplementation) ResendActivation(w http.ResponseWriter, r *http.Request) {
	var payload ResendActivationPayload
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
	resendActivationSuccess, err := handler.service.ResendActivation(ctx, payload.Email)
	if err != nil {
		switch err {
		case shared.ErrEmailInvalid:
			helpers.NotFoundError(w, err.Error())
		case shared.ErrAccountActive:
			helpers.BadRequestError(w, err.Error())
		default:
			helpers.InternalServerError(w, err.Error())
		}
		return
	}
	response := helpers.GenerateWebResponse(http.StatusCreated, shared.StatusCreated, resendActivationSuccess)
	helpers.WriteToResponseBody(w, http.StatusCreated, response)
}

func (handler *HandlerImplementation) FollowUnfollow(w http.ResponseWriter, r *http.Request) {
	userId, parseErr := uuid.Parse(chi.URLParam(r, "userId"))
	if parseErr != nil {
		helpers.BadRequestError(w, "Invalid User ID Parameters")
		return
	}
	ctx := r.Context()
	userContext := shared.GetUserFromContext(r)
	followUnfollowSuccess, err := handler.service.FollowUnfollow(ctx, userContext, userId)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			helpers.NotFoundError(w, err.Error())
		default:
			helpers.InternalServerError(w, err.Error())
		}
		return
	}
	response := helpers.GenerateWebResponse(http.StatusOK, shared.StatusOK, followUnfollowSuccess)
	helpers.WriteToResponseBody(w, http.StatusOK, response)
}

func (handler *HandlerImplementation) GetConnections(w http.ResponseWriter, r *http.Request) {
	url := path.Base(r.URL.Path)
	actionType := strings.ToUpper(url[:1]) + strings.ToLower(url[1:])
	userId, parseErr := uuid.Parse(chi.URLParam(r, "userId"))
	if parseErr != nil {
		helpers.BadRequestError(w, "Invalid User ID Parameters")
		return
	}
	ctx := r.Context()
	users, err := handler.service.GetConnections(ctx, userId, actionType)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			helpers.NotFoundError(w, err.Error())
		default:
			helpers.InternalServerError(w, err.Error())
		}
		return
	}
	response := helpers.GenerateWebResponse(http.StatusOK, shared.StatusOK, users)
	helpers.WriteToResponseBody(w, http.StatusOK, response)
}

func (handler *HandlerImplementation) Activate(w http.ResponseWriter, r *http.Request) {
	tokenStr := chi.URLParam(r, "token")
	_, claims_err := handler.authenticator.ParseJWT(tokenStr)
	if claims_err != nil {
		helpers.BadRequestError(w, claims_err.Error())
		return
	}
	ctx := r.Context()
	user, err := handler.service.Activate(ctx, tokenStr)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			helpers.NotFoundError(w, err.Error())
		default:
			helpers.InternalServerError(w, err.Error())
		}
		return
	}
	response := helpers.GenerateWebResponse(http.StatusOK, shared.StatusOK, user)
	helpers.WriteToResponseBody(w, http.StatusOK, response)
}

func (handler *HandlerImplementation) Delete(w http.ResponseWriter, r *http.Request) {
	userId, parseErr := uuid.Parse(chi.URLParam(r, "userId"))
	if parseErr != nil {
		helpers.BadRequestError(w, "Invalid User ID Parameters")
		return
	}
	ctx := r.Context()
	err := handler.service.Delete(ctx, userId)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			helpers.NotFoundError(w, err.Error())
		default:
			helpers.InternalServerError(w, err.Error())
		}
		return
	}
	response := helpers.GenerateWebResponse(http.StatusOK, shared.StatusOK, "Resource deleted successfully.")
	helpers.WriteToResponseBody(w, http.StatusOK, response)
}

func (handler *HandlerImplementation) GetFeeds(w http.ResponseWriter, r *http.Request) {
	params := &post.PostParams{
		PerPage: 10,
		Page:    1,
		Sort:    "DESC",
	}
	params, err := params.Parse(r)
	if err != nil {
		helpers.BadRequestError(w, err.Error())
		return
	}
	errMessages, err := helpers.ValidateStruct(params)
	if err != nil {
		response := helpers.GenerateWebResponse(http.StatusBadRequest, shared.StatusBadRequest, errMessages)
		helpers.WriteToResponseBody(w, http.StatusBadRequest, response)
		return
	}
	ctx := r.Context()
	userContext := shared.GetUserFromContext(r)
	paginatedFeeds, err := handler.service.GetFeeds(ctx, userContext.Id, *params)
	if err != nil {
		helpers.InternalServerError(w, err.Error())
		return
	}
	response := helpers.GenerateWebResponse(http.StatusOK, shared.StatusOK, paginatedFeeds)
	helpers.WriteToResponseBody(w, http.StatusOK, response)
}
