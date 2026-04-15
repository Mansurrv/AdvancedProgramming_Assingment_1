package apperr

import "errors"

var (
	ErrValidation          = errors.New("validation error")
	ErrEmailAlreadyExists  = errors.New("email already in use")
	ErrDoctorNotFound      = errors.New("doctor not found")
)