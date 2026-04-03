package memory

import (
	"sync"
	"time"

	"appointment-service/internal/apperr"
	"appointment-service/internal/model"
)

type AppointmentRepo struct {
	data map[string]model.Appointment
	mu   sync.RWMutex
}

func NewAppointmentRepo() *AppointmentRepo {
	return &AppointmentRepo{
		data: make(map[string]model.Appointment),
	}
}

func (r *AppointmentRepo) Create(a model.Appointment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[a.ID] = a
	return nil
}

func (r *AppointmentRepo) GetByID(id string) (*model.Appointment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	a, ok := r.data[id]
	if !ok {
		return nil, apperr.ErrAppointmentNotFound
	}

	return &a, nil
}

func (r *AppointmentRepo) GetAll() ([]model.Appointment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.Appointment
	for _, a := range r.data {
		result = append(result, a)
	}

	return result, nil
}

func (r *AppointmentRepo) UpdateStatus(id string, status model.Status) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	a, ok := r.data[id]
	if !ok {
		return apperr.ErrAppointmentNotFound
	}

	a.Status = status
	a.UpdatedAt = time.Now()
	r.data[id] = a
	return nil
}
