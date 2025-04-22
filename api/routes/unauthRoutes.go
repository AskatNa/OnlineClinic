package routes

import (
	"fmt"
	"github.com/AskatNa/OnlineClinic/api/controllers"
	"github.com/AskatNa/OnlineClinic/api/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRoute(router *gin.Engine) {
	router.GET("/ping", controllers.Ping)

	public := router.Group("/")
	{
		public.GET("/", func(c *gin.Context) {
			c.HTML(200, "index.html", gin.H{})
		})
		public.POST("/register", controllers.RegisterUser)
		public.POST("/login", controllers.Login)

		public.GET("/register", func(c *gin.Context) {
			c.HTML(200, "register.html", gin.H{})
		})
		public.GET("/login", func(c *gin.Context) {
			c.HTML(200, "login.html", gin.H{})
		})
	}

	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/profile", func(c *gin.Context) {
			c.HTML(200, "profile.html", gin.H{})
		})
	}

	admin := router.Group("/admin")
	admin.Use(middlewares.AuthMiddleware())
	{
		admin.GET("/", func(c *gin.Context) {
			email, _ := c.Get("email")
			isAdmin, exists := c.Get("isAdmin")

			if !exists || isAdmin != true {
				c.JSON(403, gin.H{"message": "Access Denied"})
				return
			}
			c.HTML(200, "adminPanel.html", gin.H{
				"email": email,
			})
		})
	}
}

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
