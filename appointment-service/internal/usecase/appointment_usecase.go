package usecase

import (
	"fmt"
	"time"

	"appointment-service/internal/apperr"
	"appointment-service/internal/model"
)

type AppointmentUsecase struct {
	repo         AppointmentRepository
	doctorClient DoctorClient
}

func NewAppointmentUsecase(r AppointmentRepository, dc DoctorClient) *AppointmentUsecase {
	return &AppointmentUsecase{
		repo:         r,
		doctorClient: dc,
	}
}

func (uc *AppointmentUsecase) CreateAppointment(a model.Appointment) error {
	if a.Title == "" {
		return fmt.Errorf("%w: title is required", apperr.ErrValidation)
	}

	if a.DoctorID == "" {
		return fmt.Errorf("%w: doctor_id is required", apperr.ErrValidation)
	}

	// Check doctor exists via REST
	exists, err := uc.doctorClient.DoctorExists(a.DoctorID)
	if err != nil {
		return fmt.Errorf("failed to validate doctor: %w", err)
	}
	if !exists {
		return fmt.Errorf("%w: doctor does not exist", apperr.ErrDoctorNotFound)
	}

	a.Status = model.StatusNew
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	return uc.repo.Create(a)
}

func (uc *AppointmentUsecase) GetAppointment(id string) (*model.Appointment, error) {
	return uc.repo.GetByID(id)
}

func (uc *AppointmentUsecase) GetAllAppointments() ([]model.Appointment, error) {
	return uc.repo.GetAll()
}

func (uc *AppointmentUsecase) UpdateStatus(id string, status model.Status) error {
	if status != model.StatusNew && status != model.StatusInProgress && status != model.StatusDone {
		return fmt.Errorf("%w: invalid status", apperr.ErrValidation)
	}

	a, err := uc.repo.GetByID(id)
	if err != nil {
		return err
	}

	exists, err := uc.doctorClient.DoctorExists(a.DoctorID)
	if err != nil {
		return fmt.Errorf("failed to validate doctor: %w", err)
	}
	if !exists {
		return fmt.Errorf("%w: doctor does not exist", apperr.ErrDoctorNotFound)
	}

	// Simple status rule: cannot go from done back to new
	if a.Status == model.StatusDone && status == model.StatusNew {
		return fmt.Errorf("%w: cannot move status from done back to new", apperr.ErrValidation)
	}

	a.Status = status
	a.UpdatedAt = time.Now()

	return uc.repo.UpdateStatus(id, status)
}
