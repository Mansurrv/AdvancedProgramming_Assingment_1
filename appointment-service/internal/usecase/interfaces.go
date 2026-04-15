package usecase

import (
	"appointment-service/internal/model"
)

type AppointmentRepository interface {
	Create(a model.Appointment) error
	GetByID(id string) (*model.Appointment, error)
	GetAll() ([]model.Appointment, error)
	UpdateStatus(id string, status model.Status) error
}

// DoctorClient is the interface the use case uses to communicate 
// with the Doctor Service, keeping it transport-agnostic.
type DoctorClient interface {
	DoctorExists(id string) (bool, error)
}