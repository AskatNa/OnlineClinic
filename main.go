package main

import (
	"fmt"
	"github.com/AskatNa/OnlineClinic/api/controllers"
	"github.com/AskatNa/OnlineClinic/api/routes"
	"github.com/AskatNa/OnlineClinic/config/configs"
	"github.com/gin-gonic/gin"
)

//var client *mongo.Client

func main() {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	//var err error
	//
	//client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	//if err != nil {
	//	log.Fatal("Error connecting to MongoDB:", err)
	//}
	//
	//err = client.Ping(ctx, nil)
	//if err != nil {
	//	log.Fatal("Couldn't connect to MongoDB:", err)
	//}
	//
	//fmt.Println("Connected to MongoDB")
	//userCollection := client.Database("online_clinic").Collection("users")
	//controllers.SetPatientCollection(userCollection)
	//controllers.SetUserCollection(userCollection)
	//router.Use(func(c *gin.Context) {
	//	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:9000")
	//	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	//	c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	//	c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	//	c.Writer.Header().Set("Pragma", "no-cache")
	//	c.Writer.Header().Set("Expires", "0")
	//	if c.Request.Method == "OPTIONS" {
	//		c.AbortWithStatus(204)
	//		return
	//	}
	//	c.Next()
	//	fmt.Println("⚠️ JWT Secret Key:", configs.JWTSecret)
	//})
	router := gin.Default()
	client := configs.ConnectDB()
	//controllers.SetPatientCollection(configs.GetCollection(client, "patients"))
	controllers.SetUserCollection(configs.GetCollection(client, "users"))
	//controllers.SetDoctorCollection(configs.GetCollection(client, "doctors"))

	router.LoadHTMLGlob("./ui/html/*")
	routes.SetupRoute(router)
	//routes.UnauthRoutes(router)
	routes.UserRoutes(router)
	routes.DoctorRoutes(router)
	fmt.Println("The server is running on port :9000")
	//log.Fatal(http.ListenAndServe(":9000", router))
	router.Run(":" + configs.EnvPort())
}
