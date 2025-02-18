package routes

import (
	"github.com/AskatNa/OnlineClinic/api/controllers"
	"github.com/AskatNa/OnlineClinic/api/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoute(router *gin.Engine) {
	router.GET("/ping", controllers.Ping)

	// Public routes
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

	// Protected routes (requires authentication)
	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/profile", func(c *gin.Context) {
			c.HTML(200, "profile.html", gin.H{})
		})
	}

	// Admin routes (requires authentication & admin privileges)
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

// import (
//
//	"github.com/AskatNa/OnlineClinic/api/controllers"
//	"github.com/AskatNa/OnlineClinic/api/middlewares"
//	"github.com/gin-gonic/gin"
//
// )
//
// //func UnauthRoutes(router *gin.Engine) {
// //	router.GET("/ping", controllers.Ping)
// //	//Home page
// //	router.GET("/", func(c *gin.Context) {
// //		c.HTML(200, "index.html", gin.H{})
// //	})
// //	//Registration
// //	router.POST("/register", controllers.RegisterUser)
// //	router.GET("/register", func(c *gin.Context) {
// //		c.HTML(200, "register.html", gin.H{})
// //	})
// //
// //	//Login
// //	router.POST("/login", controllers.Login)
// //	router.GET("/login", func(c *gin.Context) {
// //		c.HTML(200, "login.html", gin.H{})
// //	})
// //	//adminRoute := router.Group("/admin")
// //	//adminRoute.Use(middlewares.AuthMiddleware())
// //	//{
// //	//	adminRoute.GET("/", func(c *gin.Context) {
// //	//		email, _ := c.Get("email")
// //	//		isAdmin, _ := c.Get("isAdmin")
// //	//
// //	//		fmt.Println("Extracted Email:", email, "Is Admin:", isAdmin)
// //	//
// //	//		if email != "admin798@gmail.com" {
// //	//			c.JSON(http.StatusForbidden, gin.H{"message": "Access Denied"})
// //	//			return
// //	//		}
// //	//		c.HTML(http.StatusOK, "adminPanel.html", gin.H{})
// //	//	})
// //	//}
// //
// //}
//
// // SetupRoutes registers all routes and connects them to HTML pages
//
//	func SetupRoutes(router *gin.Engine) {
//		router.GET("/ping", controllers.Ping)
//
//		// Public routes
//		public := router.Group("/")
//		{
//			public.GET("/", func(c *gin.Context) {
//				c.HTML(200, "index.html", gin.H{})
//			})
//			public.POST("/register", controllers.RegisterUser)
//			public.POST("/login", controllers.Login)
//
//			// Serve registration page
//			public.GET("/register", func(c *gin.Context) {
//				c.HTML(200, "register.html", gin.H{})
//			})
//
//			// Serve login page
//			public.GET("/login", func(c *gin.Context) {
//				c.HTML(200, "login.html", gin.H{})
//			})
//		}
//
//		// Protected routes (requires authentication)
//		protected := router.Group("/")
//		protected.Use(middlewares.AuthMiddleware())
//		{
//			// Profile page for authenticated users
//			protected.GET("/profile", func(c *gin.Context) {
//				c.HTML(200, "profile.html", gin.H{})
//			})
//		}
//
//		// Admin routes (requires authentication & admin privileges)
//		admin := router.Group("/admin")
//		admin.Use(middlewares.AuthMiddleware())
//		{
//			// Admin dashboard page
//			admin.GET("/", func(c *gin.Context) {
//				email, _ := c.Get("email")
//				isAdmin, exists := c.Get("isAdmin")
//
//				if !exists || isAdmin != true {
//					c.JSON(403, gin.H{"message": "Access Denied"})
//					return
//				}
//
//				// Render the admin panel if the user is an admin
//				c.HTML(200, "adminPanel.html", gin.H{
//					"email": email,
//				})
//			})
//		}
//	}
