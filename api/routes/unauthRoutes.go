package routes

import (
	"fmt"
	"github.com/AskatNa/OnlineClinic/api/controllers"
	"github.com/AskatNa/OnlineClinic/api/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
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
		c.HTML(200, "register.html", gin.H{})
	})

	//Login
	router.POST("/login", controllers.Login)
	router.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{})
	})
	adminRoute := router.Group("/admin")
	adminRoute.Use(middlewares.AuthMiddleware())
	{
		adminRoute.GET("/", func(c *gin.Context) {
			email, _ := c.Get("email")
			isAdmin, _ := c.Get("isAdmin")

			fmt.Println("Extracted Email:", email, "Is Admin:", isAdmin)

			if email != "admin798@gmail.com" {
				c.JSON(http.StatusForbidden, gin.H{"message": "Access Denied"})
				return
			}
			c.HTML(http.StatusOK, "adminPanel.html", gin.H{})
		})
	}
}
