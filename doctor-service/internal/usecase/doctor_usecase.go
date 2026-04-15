package usecase

import (
	"doctor-service/internal/apperr"
	"doctor-service/internal/model"
	"fmt"
)

type DoctorUseCase struct {
	repo DoctorRepository
}

func NewDoctorUseCase(repo DoctorRepository) *DoctorUseCase {
	return &DoctorUseCase{repo: repo}
}

func (uc *DoctorUseCase) CreateDoctor(d model.Doctor) error {
	if d.FullName == "" {
		return fmt.Errorf("%w: full name is required", apperr.ErrValidation)
	}
	if d.Email == "" {
		return fmt.Errorf("%w: email is required", apperr.ErrValidation)
	}
	existing, _ := uc.repo.GetByEmail(d.Email)
	if existing != nil {
		return fmt.Errorf("%w: doctor with this email already exists", apperr.ErrEmailAlreadyExists)
	}
	return uc.repo.Create(d)
}

func (uc *DoctorUseCase) GetDoctorByID(id string) (*model.Doctor, error) {
	return uc.repo.GetByID(id)
}

func (uc *DoctorUseCase) GetAllDoctors() ([]model.Doctor, error) {
	return uc.repo.GetAll()
}
