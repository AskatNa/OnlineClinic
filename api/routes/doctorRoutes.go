package routes

//)
//
//func DoctorRoutes(router *gin.Engine) {
//	doctorGroup := router.Group("/doctor")
//	doctorGroup.Use(middlewares.AuthMiddleware())
//	{
//		doctorGroup.GET("/profile", controllers.GetAllDoctors)
//	}
//	router.GET("/doctor/dashboard", func(c *gin.Context) {
//		c.HTML(200, "doctorProfile.html", gin.H{})
//	})
//}
