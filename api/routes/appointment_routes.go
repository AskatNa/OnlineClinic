package routes

import (
	"github.com/AskatNa/OnlineClinic/api/controllers"
	"github.com/gin-gonic/gin"
)

func AppointmentRoutes(router *gin.Engine) {
	router.POST("/appointments/book", controllers.BookAppointmentSlot)
	router.POST("/appointments/cancel", controllers.CancelAppointmentSlot)
	router.POST("/appointments/details", controllers.ViewAppointmentDetails)
	router.GET("/appointments/:id/history", controllers.ViewPatientAppointmentsHistory)
}
