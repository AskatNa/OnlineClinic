package main

import (
	"log"

	"github.com/AskatNa/OnlineClinic/api/routes"
	"github.com/AskatNa/OnlineClinic/config/configs"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	configs.ConnectDB()

	router.LoadHTMLGlob("templates/*")

	routes.UnauthRoutes(router)
	routes.UserRoutes(router)
	routes.DoctorRoutes(router)
	routes.PatientRoutes(router)
	routes.AppointmentRoutes(router)

	if err := router.Run(":9000"); err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}
