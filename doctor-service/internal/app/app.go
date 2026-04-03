package app

import (
	"github.com/gin-gonic/gin"

	"doctor-service/internal/repository/memory"
	"doctor-service/internal/transport/http"
	"doctor-service/internal/usecase"
)

func Run() {
	repo := memory.NewDoctorRepo()
	uc := usecase.NewDoctorUseCase(repo)
	handler := http.NewDoctorHandler(uc)

	r := gin.Default()

	r.POST("/doctors", handler.CreateDoctor)
	r.GET("/doctors/:id", handler.GetDoctor)
	r.GET("/doctors", handler.GetAllDoctors)

	r.Run(":8080")
}
