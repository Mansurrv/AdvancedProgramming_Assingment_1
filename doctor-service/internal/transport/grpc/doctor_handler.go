package grpc

import (
	"context"
	"errors"

	"doctor-service/internal/apperr"
	"doctor-service/internal/model"
	"doctor-service/internal/usecase"
	doctorpb "doctor-service/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DoctorHandler struct {
	doctorpb.UnimplementedDoctorServiceServer
	uc *usecase.DoctorUseCase
}

func NewDoctorHandler(uc *usecase.DoctorUseCase) *DoctorHandler {
	return &DoctorHandler{uc: uc}
}

func (h *DoctorHandler) CreateDoctor(_ context.Context, req *doctorpb.CreateDoctorRequest) (*doctorpb.DoctorResponse, error) {
	doctor := model.Doctor{
		ID:             uuid.New().String(),
		FullName:       req.GetFullName(),
		Specialization: req.GetSpecialization(),
		Email:          req.GetEmail(),
	}

	if err := h.uc.CreateDoctor(doctor); err != nil {
		return nil, mapError(err)
	}

	return toDoctorResponse(doctor), nil
}

func (h *DoctorHandler) GetDoctor(_ context.Context, req *doctorpb.GetDoctorRequest) (*doctorpb.DoctorResponse, error) {
	doctor, err := h.uc.GetDoctorByID(req.GetId())
	if err != nil {
		return nil, mapError(err)
	}

	return toDoctorResponse(*doctor), nil
}

func (h *DoctorHandler) ListDoctors(context.Context, *doctorpb.ListDoctorsRequest) (*doctorpb.ListDoctorsResponse, error) {
	doctors, err := h.uc.GetAllDoctors()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &doctorpb.ListDoctorsResponse{
		Doctors: make([]*doctorpb.DoctorResponse, 0, len(doctors)),
	}

	for _, doctor := range doctors {
		response.Doctors = append(response.Doctors, toDoctorResponse(doctor))
	}

	return response, nil
}

func toDoctorResponse(doctor model.Doctor) *doctorpb.DoctorResponse {
	return &doctorpb.DoctorResponse{
		Id:             doctor.ID,
		FullName:       doctor.FullName,
		Specialization: doctor.Specialization,
		Email:          doctor.Email,
	}
}

func mapError(err error) error {
	switch {
	case errors.Is(err, apperr.ErrValidation):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, apperr.ErrEmailAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, apperr.ErrDoctorNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
