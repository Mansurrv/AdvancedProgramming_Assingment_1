package apperr

import "errors"

var (
	ErrValidation               = errors.New("validation error")
	ErrDoctorNotFound           = errors.New("doctor not found")
	ErrDoctorServiceUnavailable = errors.New("doctor service unavailable")
	ErrDoctorServiceError       = errors.New("doctor service error")
	ErrAppointmentNotFound      = errors.New("appointment not found")
)
