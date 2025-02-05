package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AskatNa/OnlineClinic/api/controllers"
	"github.com/AskatNa/OnlineClinic/api/routes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Couldn't connect to MongoDB:", err)
	}
	fmt.Println("Connected to MongoDB")

	userCollection := client.Database("online_clinic").Collection("users")
	controllers.SetUserCollection(userCollection)

	router := gin.Default()

	routes.UnauthRoutes(router)
	routes.UserRoutes(router)

	fmt.Println("The server is running on port :9000")
	log.Fatal(http.ListenAndServe(":9000", router))
}
