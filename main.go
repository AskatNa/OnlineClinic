package main

import (
	"fmt"
	"github.com/AskatNa/OnlineClinic/api/controllers"
	"github.com/AskatNa/OnlineClinic/api/routes"
	"github.com/AskatNa/OnlineClinic/config/configs"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	client := configs.ConnectDB()
	controllers.SetUserCollection(configs.GetCollection(client, "users"))

	router.LoadHTMLGlob("./ui/html/*")
	routes.SetupRoute(router)
	routes.UserRoutes(router)
	routes.DoctorRoutes(router)
	fmt.Println("The server is running on port :9000")
	router.Run(":" + configs.EnvPort())
}
