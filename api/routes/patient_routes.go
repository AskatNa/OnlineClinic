package routes

import (
	"github.com/AskatNa/OnlineClinic/api/controllers"
	"github.com/gin-gonic/gin"
)

func PatientRoutes(router *gin.Engine) {
	//router.GET("/patients", controllers.GetAllPatients)
	router.GET("/patients/:id", controllers.GetPatientById)
}
