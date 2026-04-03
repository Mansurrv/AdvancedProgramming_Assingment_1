package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"doctor-service/internal/model"
	"doctor-service/internal/usecase"
)

type DoctorHandler struct {
	uc *usecase.DoctorUseCase
}

func NewDoctorHandler(uc *usecase.DoctorUseCase) *DoctorHandler {
	return &DoctorHandler{uc: uc}
}

type createDoctorRequest struct {
	FullName       string `json:"full_name"`
	Specialization string `json:"specialization"`
	Email          string `json:"email"`
}

func (h *DoctorHandler) CreateDoctor(c *gin.Context) {
	var req createDoctorRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doctor := model.Doctor{
		ID:             uuid.New().String(),
		FullName:       req.FullName,
		Specialization: req.Specialization,
		Email:          req.Email,
	}

	if err := h.uc.CreateDoctor(doctor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, doctor)
}

func (h *DoctorHandler) GetDoctor(c *gin.Context) {
	id := c.Param("id")

	doctor, err := h.uc.GetDoctorByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, doctor)
}

func (h *DoctorHandler) GetAllDoctors(c *gin.Context) {
	doctors, _ := h.uc.GetAllDoctors()
	c.JSON(http.StatusOK, doctors)
}
