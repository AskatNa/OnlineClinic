package routes

import (
	"github.com/AskatNa/OnlineClinic/api/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.POST("/users", controllers.CreateUser)
	router.GET("/users/:userId", controllers.GetAUser)
	router.GET("/users", controllers.GetAllUsers)
	router.PUT("/users/:userId", controllers.UpdateUser)
	router.DELETE("/users/:userId", controllers.DeleteUser)
}
