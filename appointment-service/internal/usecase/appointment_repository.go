package usecase

import "appointment-service/internal/model"

type AppointmentRepository interface {
	Create(a model.Appointment) error
	GetByID(id string) (*model.Appointment, error)
	GetAll() ([]model.Appointment, error)
	UpdateStatus(id string, status model.Status) error
}
