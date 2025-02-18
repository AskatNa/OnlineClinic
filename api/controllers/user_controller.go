package controllers

import (
	"context"
	"fmt"
	"github.com/AskatNa/OnlineClinic/api/models"
	"net/http"
	"time"

	"github.com/AskatNa/OnlineClinic/api/responses"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

func SetUserCollection(collection *mongo.Collection) {
	userCollection = collection
}

func CreateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only the admin can create users"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid data format"})
		return
	}

	if user.Role != "doctor" && user.Role != "patient" {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Invalid role. Allowed roles: 'doctor', 'patient'"})
		return
	}

	insertResult, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error creating user"})
		return
	}

	c.JSON(http.StatusCreated, responses.UserResponse{
		Status:  http.StatusCreated,
		Message: "User successfully created",
		Data:    map[string]interface{}{"id": insertResult.InsertedID, "role": user.Role},
	})
}
func GetUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.Param("userId")
	objId, _ := primitive.ObjectIDFromHex(userId)

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, responses.UserResponse{
			Status:  http.StatusNotFound,
			Message: "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    map[string]interface{}{"user": user},
	})
}

func GetAllUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []models.User
	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error retrieving users",
		})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			continue
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    map[string]interface{}{"users": users},
	})
}
func UpdateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.Param("userId")
	objId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}
	isAdminInterface, exists := c.Get("isAdmin")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}

	isAdmin, ok := isAdminInterface.(bool)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid admin status"})
		c.Abort()
		return
	}
	currentUserEmail, _ := c.Get("email")

	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only admins can update users"})
		return
	}

	// Log to debug admin status
	fmt.Println("Admin Status:", isAdmin)
	fmt.Println("Current User Email:", currentUserEmail)

	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	if user.Email == currentUserEmail {
		c.JSON(http.StatusForbidden, gin.H{"message": "You cannot change your own role"})
		return
	}

	var updatedData struct {
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&updatedData); err != nil || (updatedData.Role != "doctor" && updatedData.Role != "patient") {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid role. Allowed: 'doctor', 'patient'"})
		return
	}

	_, err = userCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": bson.M{"role": updatedData.Role}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating user role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully"})
}

func DeleteUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	isAdminInterface, exists := c.Get("isAdmin")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}

	isAdmin, ok := isAdminInterface.(bool)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid admin status"})
		c.Abort()
		return
	}
	currentUserEmail, _ := c.Get("email")

	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only admins can delete users"})
		return
	}

	userId := c.Param("userId")
	objId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}

	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	if user.Email == currentUserEmail {
		c.JSON(http.StatusForbidden, gin.H{"message": "You cannot delete your own account"})
		return
	}

	result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil || result.DeletedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
