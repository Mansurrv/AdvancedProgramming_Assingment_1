package usecase

import (
	"doctor-service/internal/model"
	"errors"
)

type DoctorUseCase struct {
	repo DoctorRepository
}

func NewDoctorUseCase(repo DoctorRepository) *DoctorUseCase {
	return &DoctorUseCase{repo: repo}
}

func (uc *DoctorUseCase) CreateDoctor(d model.Doctor) error {
	if d.FullName == "" {
		return errors.New("full name is required")
	}
	if d.Email == "" {
		return errors.New("email is required")
	}
	existing, _ := uc.repo.GetByEmail(d.Email)
	if existing != nil {
		return errors.New("doctor with this email already exists")
	}
	return uc.repo.Create(d)
}

func (uc *DoctorUseCase) GetDoctorByID(id string) (*model.Doctor, error) {
	return uc.repo.GetByID(id)
}

func (uc *DoctorUseCase) GetAllDoctors() ([]model.Doctor, error) {
	return uc.repo.GetAll()
}
