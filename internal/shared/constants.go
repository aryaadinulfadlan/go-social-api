package shared

import "errors"

const (
	StatusBadRequest          = "BAD REQUEST"
	StatusInternalServerError = "INTERNAL SERVER ERROR"
	StatusNotFound            = "NOT FOUND"
	StatusUnauthorized        = "UNAUTHORIZED"
	StatusForbidden           = "FORBIDDEN"
	StatusMethodNotAllowed    = "METHOD NOT ALLOWED"
	StatusCreated             = "CREATED"
	StatusOK                  = "OK"
	StatusTooManyRequests     = "TOO MANY REQUESTS"
)

var (
	ErrNotFound        = errors.New("resource not found.")
	ErrBadRequest      = errors.New("bad request.")
	ErrUserExists      = errors.New("user already exists.")
	ErrLoginInvalid    = errors.New("invalid email or password.")
	ErrAccountInactive = errors.New("account is not active.")
	ErrAccountActive   = errors.New("account is active.")
	ErrEmailInvalid    = errors.New("invalid email.")
	ErrForbidden       = errors.New("you do not have permission to access this resource.")
)
