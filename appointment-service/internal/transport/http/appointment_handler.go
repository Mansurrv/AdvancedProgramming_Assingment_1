package http

import (
	"errors"
	"net/http"

	"appointment-service/internal/apperr"
	"appointment-service/internal/model"
	"appointment-service/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AppointmentHandler struct {
	uc *usecase.AppointmentUsecase
}

func NewAppointmentHandler(uc *usecase.AppointmentUsecase) *AppointmentHandler {
	return &AppointmentHandler{uc: uc}
}

type createAppointmentRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DoctorID    string `json:"doctor_id"`
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

func writeError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, apperr.ErrValidation):
		status = http.StatusBadRequest
	case errors.Is(err, apperr.ErrDoctorNotFound):
		status = http.StatusBadRequest
	case errors.Is(err, apperr.ErrAppointmentNotFound):
		status = http.StatusNotFound
	case errors.Is(err, apperr.ErrDoctorServiceUnavailable):
		status = http.StatusServiceUnavailable
	case errors.Is(err, apperr.ErrDoctorServiceError):
		status = http.StatusBadGateway
	}

	c.JSON(status, gin.H{"error": err.Error()})
}

func (h *AppointmentHandler) CreateAppointment(c *gin.Context) {
	var req createAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a := model.Appointment{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		DoctorID:    req.DoctorID,
	}

	if err := h.uc.CreateAppointment(a); err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, a)
}

func (h *AppointmentHandler) GetAppointment(c *gin.Context) {
	id := c.Param("id")
	a, err := h.uc.GetAppointment(id)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, a)
}

func (h *AppointmentHandler) GetAllAppointments(c *gin.Context) {
	list, err := h.uc.GetAllAppointments()
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *AppointmentHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req updateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.UpdateStatus(id, model.Status(req.Status)); err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": req.Status})
}
