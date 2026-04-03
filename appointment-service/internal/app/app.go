package app

import (
	"appointment-service/internal/repository/http"
	"appointment-service/internal/repository/memory"
	repohttp "appointment-service/internal/transport/http"
	"appointment-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

func Run() {
	repo := memory.NewAppointmentRepo()
	doctorClient := http.NewDoctorClientHTTP("http://localhost:8080") // DoctorService URL
	uc := usecase.NewAppointmentUsecase(repo, doctorClient)
	handler := repohttp.NewAppointmentHandler(uc)

	r := gin.Default()
	r.POST("/appointments", handler.CreateAppointment)
	r.GET("/appointments/:id", handler.GetAppointment)
	r.GET("/appointments", handler.GetAllAppointments)
	r.PATCH("/appointments/:id/status", handler.UpdateStatus)
	r.Run(":8081")
}
