package memory

import (
	"sync"

	"doctor-service/internal/apperr"
	"doctor-service/internal/model"
)

type DoctorRepo struct {
	data map[string]model.Doctor
	mu   sync.RWMutex
}

func NewDoctorRepo() *DoctorRepo {
	return &DoctorRepo{
		data: make(map[string]model.Doctor),
	}
}

func (r *DoctorRepo) Create(d model.Doctor) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[d.ID] = d
	return nil
}

func (r *DoctorRepo) GetByID(id string) (*model.Doctor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	d, ok := r.data[id]
	if !ok {
		return nil, apperr.ErrDoctorNotFound
	}

	return &d, nil
}

func (r *DoctorRepo) GetAll() ([]model.Doctor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.Doctor
	for _, d := range r.data {
		result = append(result, d)
	}

	return result, nil
}

func (r *DoctorRepo) GetByEmail(email string) (*model.Doctor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, d := range r.data {
		if d.Email == email {
			return &d, nil
		}
	}

	return nil, nil
}
