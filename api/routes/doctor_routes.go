package routes

import (
	"github.com/AskatNa/OnlineClinic/api/controllers"
	"github.com/gin-gonic/gin"
)

func DoctorRoutes(router *gin.Engine) {
	router.GET("/doctors", controllers.GetAllDoctors)
	router.GET("/doctors/available", controllers.GetAvailableDoctors)
	router.GET("/doctors/name/:name", controllers.GetDoctorByName)
	router.GET("/doctors/id/:id", controllers.GetDoctorById)
	router.GET("/doctors/schedule/:id", controllers.GetDoctorScheduleById)
}
