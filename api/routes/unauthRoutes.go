package routes

import (
	"github.com/AskatNa/OnlineClinic/api/controllers"
	"github.com/gin-gonic/gin"
)

func UnauthRoutes(router *gin.Engine) {
	router.GET("/ping", controllers.Ping)
	//Home page
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
	//Registration
	router.POST("/register", controllers.RegisterUser)
	router.GET("/register", func(c *gin.Context) {
		c.HTML(200, "register.html", gin.H{}) // Renders register.html
	})

	//Login
	router.POST("/login", controllers.Login)
	router.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{}) // Renders login.html
	})
}
