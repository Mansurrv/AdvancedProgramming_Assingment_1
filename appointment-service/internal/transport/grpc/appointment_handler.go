package grpc

import (
	"context"
	"errors"
	"time"

	"appointment-service/internal/apperr"
	"appointment-service/internal/model"
	"appointment-service/internal/usecase"
	appointmentpb "appointment-service/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppointmentHandler struct {
	appointmentpb.UnimplementedAppointmentServiceServer
	uc *usecase.AppointmentUsecase
}

func NewAppointmentHandler(uc *usecase.AppointmentUsecase) *AppointmentHandler {
	return &AppointmentHandler{uc: uc}
}

func (h *AppointmentHandler) CreateAppointment(_ context.Context, req *appointmentpb.CreateAppointmentRequest) (*appointmentpb.AppointmentResponse, error) {
	appointment := model.Appointment{
		ID:          uuid.New().String(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		DoctorID:    req.GetDoctorId(),
	}

	if err := h.uc.CreateAppointment(appointment); err != nil {
		return nil, mapError(err)
	}

	stored, err := h.uc.GetAppointment(appointment.ID)
	if err != nil {
		return nil, mapError(err)
	}

	return toAppointmentResponse(*stored), nil
}

func (h *AppointmentHandler) GetAppointment(_ context.Context, req *appointmentpb.GetAppointmentRequest) (*appointmentpb.AppointmentResponse, error) {
	appointment, err := h.uc.GetAppointment(req.GetId())
	if err != nil {
		return nil, mapError(err)
	}

	return toAppointmentResponse(*appointment), nil
}

func (h *AppointmentHandler) ListAppointments(context.Context, *appointmentpb.ListAppointmentsRequest) (*appointmentpb.ListAppointmentsResponse, error) {
	appointments, err := h.uc.GetAllAppointments()
	if err != nil {
		return nil, mapError(err)
	}

	response := &appointmentpb.ListAppointmentsResponse{
		Appointments: make([]*appointmentpb.AppointmentResponse, 0, len(appointments)),
	}

	for _, appointment := range appointments {
		response.Appointments = append(response.Appointments, toAppointmentResponse(appointment))
	}

	return response, nil
}

func (h *AppointmentHandler) UpdateAppointmentStatus(_ context.Context, req *appointmentpb.UpdateStatusRequest) (*appointmentpb.AppointmentResponse, error) {
	if err := h.uc.UpdateStatus(req.GetId(), model.Status(req.GetStatus())); err != nil {
		return nil, mapError(err)
	}

	updated, err := h.uc.GetAppointment(req.GetId())
	if err != nil {
		return nil, mapError(err)
	}

	return toAppointmentResponse(*updated), nil
}

func toAppointmentResponse(appointment model.Appointment) *appointmentpb.AppointmentResponse {
	return &appointmentpb.AppointmentResponse{
		Id:          appointment.ID,
		Title:       appointment.Title,
		Description: appointment.Description,
		DoctorId:    appointment.DoctorID,
		Status:      string(appointment.Status),
		CreatedAt:   appointment.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   appointment.UpdatedAt.Format(time.RFC3339),
	}
}

func mapError(err error) error {
	switch {
	case errors.Is(err, apperr.ErrValidation):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, apperr.ErrDoctorNotFound):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, apperr.ErrAppointmentNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, apperr.ErrDoctorServiceUnavailable):
		return status.Error(codes.Unavailable, err.Error())
	case errors.Is(err, apperr.ErrDoctorServiceError):
		return status.Error(codes.Unavailable, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
