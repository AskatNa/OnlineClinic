package routes

import (
	"github.com/AskatNa/OnlineClinic/api/controllers"
	"github.com/AskatNa/OnlineClinic/api/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.POST("/createUser", middlewares.AuthMiddleware(), controllers.CreateUser)
	router.GET("/users/:userId", controllers.GetUser)
	router.GET("/users", controllers.GetAllUsers)
	router.PUT("/users/:userId", controllers.UpdateUser)
	router.DELETE("/users/:userId", middlewares.AuthMiddleware(), controllers.DeleteUser)
}
