package controllers

import (
	"context"
	"github.com/AskatNa/OnlineClinic/api/responses"
	"github.com/AskatNa/OnlineClinic/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, responses.UserResponse{
		Status: http.StatusOK, Message: "Pong",
	})
}
func RegisterUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var data models.User
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status: http.StatusBadRequest, Message: "Invalid data format",
		})
		return
	}
	result, err := userCollection.InsertOne(ctx, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError,
			Message: "Error registering user"})
		return
	}
	c.JSON(http.StatusCreated, responses.UserResponse{
		Status: http.StatusCreated, Message: "User successfully registered",
		Data: map[string]interface{}{"id": result.InsertedID}})
}
func Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var data map[string]string
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest,
			Message: "Invalid login data format"})
		return
	}
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": data["email"]}).Decode(&user)
	if err != nil || user.Password != data["password"] {
		c.JSON(http.StatusUnauthorized, responses.UserResponse{Status: http.StatusUnauthorized,
			Message: "Invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, responses.UserResponse{
		Status: http.StatusOK, Message: "Login successful",
		Data: map[string]interface{}{"user": user}})
}
