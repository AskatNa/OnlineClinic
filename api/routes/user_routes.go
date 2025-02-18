package routes

import (
	"github.com/AskatNa/OnlineClinic/api/controllers"
	"github.com/AskatNa/OnlineClinic/api/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {

	router.POST("/createUser", middlewares.AuthMiddleware(), controllers.CreateUser)
	router.GET("/userByID/:userId", controllers.GetUser)
	router.GET("/users", controllers.GetAllUsers)
	router.PUT("/updateUser/:userId", middlewares.AuthMiddleware(), controllers.UpdateUser)
	router.DELETE("/deleteUserByID/:userId", middlewares.AuthMiddleware(), controllers.DeleteUser)
}
