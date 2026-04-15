package apperr

import "errors"

var (
	ErrValidation               = errors.New("validation error")
	ErrAppointmentNotFound      = errors.New("appointment not found")
	ErrDoctorNotFound           = errors.New("doctor not found")
	ErrDoctorServiceUnavailable = errors.New("doctor service unavailable")
	ErrDoctorServiceError       = errors.New("doctor service internal error")
)